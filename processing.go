package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
	"time"
)

type Participant struct {
	id          string
	sessiondate string
	trials      []*Trial
	status      string
}
type Trial struct {
	id            string
	procedure     string
	trialnumber   string
	stimulus      string
	accuracy      float64
	correctanswer string
	response      string
	rt            float64
	stimuli       string
	clockcheck    float64
	falsealarm    float64
	pm3           float64
	pm6           float64
	pm9           float64
}

// Log file search criteria
var logframeend = "*** LogFrame End ***"
var logframestart = "*** LogFrame Start ***"
var measurevartypes = map[string]string{
	"Procedure:":        "procedure",
	"TBTrialList:":      "trialnumber",
	"LDTrialList:":      "trialnumber",
	"PracTrialList:":    "trialnumber",
	"Stimulus.RT:":      "rt",
	"Stimulus2.RT:":     "rt",
	"Stimulus2.ACC:":    "accuracy",
	"Stimulus.ACC:":     "accuracy",
	"Stimulus.RESP:":    "response",
	"Stimulus2.RESP:":   "response",
	"Stimulus.CRESP:":   "correctanswer",
	"Stimulus2.CRESP:":  "correctanswer",
	"Stimuli:":          "stimuli",
	"TimeOfClockCheck":  "timeofclockcheck",
	"TimeOfPMResponse:": "pmresponseclock",
}

// processFile opens a single file, processes all lines, and returns a participant
func processFile(file string) (subj *Participant, err error) {
	// Open file
	data, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer data.Close()

	// Read lines
	scanner := bufio.NewScanner(data)
	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "reading inputfile:", err)
	}

	// Initiate pointer to a new participant and (first) trial
	thissubject := &Participant{}
	thistrial := &Trial{}

	// For line in file
	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, "Subject:") && thissubject.id == "" {
			_, thissubject.id = ExtractKeyValue(line)
		}
		if strings.Contains(line, "SessionDate:") && thissubject.sessiondate == "" {
			_, thissubject.sessiondate = ExtractKeyValue(line)
		}

		if strings.Contains(line, logframestart) {
			thistrial = &Trial{}
		}

		if isMeasureVariable(line) {
			processMeasure(line, thistrial)
		}

		if strings.Contains(line, logframeend) {
			thissubject.trials = append(thissubject.trials, thistrial)
		}

	}

	return thissubject, nil
}

// isMeasureVariable checks whether the given measure var is in the lookups types
func isMeasureVariable(line string) bool {
	for k := range measurevartypes {
		if strings.Contains(line, k) {
			return true
		}
	}
	return false
}

// processMeasure Takes a single output line and a pointer to Trial
// Then calls utility functions to update point to Trial
func processMeasure(line string, results *Trial) {
	// Takes a line and a trial, processes trial
	if strings.Contains(line, ":") {
		measuretype, _ := GetMeasureType(line)
		_, value := ExtractKeyValue(line)
		ProcessKeyValue(measuretype, value, results)
	}
}

// GetMeasureType Searches line for variable in measure vartpes. Returns type if found.
func GetMeasureType(line string) (measuretype string, err error) {
	for k, v := range measurevartypes {
		if strings.Contains(line, k) {
			measuretype, err = v, nil
		}
	}
	return measuretype, err
}

func ExtractKeyValue(line string) (value string, key string) {
	//Extracts value from key, value pair seperated by `:`
	split := strings.Fields(line)
	if len(split) < 2 {
		key = split[0]
		value = "NA"
	} else {
		key = split[0]
		value = split[len(split)-1]
	}
	return key, value
}

func ProcessKeyValue(measuretype string, value string, trial *Trial) {
	// Process rt
	if measuretype == "rt" {
		trial.rt, _ = strconv.ParseFloat(value, 64)
	}
	if measuretype == "procedure" {
		trial.procedure = value
	}
	if measuretype == "accuracy" {
		trial.accuracy, _ = strconv.ParseFloat(value, 64)
	}
	if measuretype == "timeofclockcheck" {
		trial.clockcheck = 1
	}
	if measuretype == "pmresponseclock" {
		pmresponse, pmnumber := processPMResponse(value)
		if pmnumber == 3 && pmresponse == "correct" {
			trial.pm3 = 1
		}
		if pmnumber == 6 && pmresponse == "correct" {
			trial.pm6 = 1
		}
		if pmnumber == 9 && pmresponse == "correct" {
			trial.pm9 = 1
		}
		// Else if not one of the clock numbers
		if pmnumber == 0 && pmresponse == "falsealarm" {
			trial.falsealarm = 1
		}
	}
}

func processPMResponse(pmresponseclock string) (pmresponse string, pmnumber int) {
	clocktimes := []string{"00:03:00", "00:06:00", "00:09:00"}
	clockindex := []int{3, 6, 9}
	pmresponse = "falsealarm"
	pmnumber = 0
	for index, atime := range clocktimes {
		correctime, _ := time.Parse("15:04:05", atime)
		pmresponsetime, _ := time.Parse("15:04:05", pmresponseclock)
		responseoffset := math.Abs(correctime.Sub(pmresponsetime).Seconds())
		if responseoffset <= threshold {
			pmresponse = "correct"
			pmnumber = clockindex[index]
		}
	}
	return pmresponse, pmnumber
}
