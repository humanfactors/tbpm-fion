package main

import (
	"encoding/csv"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/zpatrick/go-config"
)

// Error Checker Utility Function
func checkError(message string, err error) {
	if err != nil {
		log.Fatal(message, err)
	}
}

// Determine Participant ABI or Control
var pStatus = regexp.MustCompile(".*?(control|abi).*?")

func getParticipantStatus(filename string) (status string) {
	status = pStatus.FindStringSubmatch(filename)[1]
	if status != "control" && status != "abi" {
		status = "filename did not contain control or abi"
	}
	return status
}

// Configuration Settings
var threshold float64
var outfilename string

func main() {

	// Load settings
	iniFile := config.NewINIFile("settings.ini")
	settings := config.NewConfig([]config.Provider{iniFile})
	threshold, _ = settings.Float("global.pmthreshold")
	outfilename, _ = settings.String("global.outfilename")

	// Open outfile & Write Headers
	file, err := os.Create(outfilename)
	checkError("Cannot create file", err)
	defer file.Close()
	writer := csv.NewWriter(file)
	defer writer.Flush()
	writer.Write([]string{"subj", "status", "testdate",
		"Accuacy_PracTrialProc", "CorrectRT_PracTrialProc", "IncorrectRT_PracTrialProc",
		"Accuacy_LDTrialProc", "CorrectRT_LDTrialProc", "IncorrectRT_LDTrialProc",
		"Accuacy_TBTrialProc", "CorrectRT_TBTrialProc", "IncorrectRT_TBTrialProc", "ClockChecks_TOTAL",
		"FalseAlarm_TOTAL", "PM_3min", "PM_6min", "PM_9min", "PM_TOTAL"})

	// For all arguments, process file and then write out results
	for _, file := range os.Args[1:] {
		thissubject, err := processFile(file)
		checkError("Error reading participant file", err)

		thissubject.status = getParticipantStatus(strings.ToLower(file))
		outvars := []string{
			thissubject.id,
			thissubject.status,
			thissubject.sessiondate,
			strconv.FormatFloat(thissubject.AverageAccuracy("PracTrialProc"), 'f', 10, 64),
			strconv.FormatFloat(thissubject.AverageCorrectRT("PracTrialProc", "correct"), 'f', 10, 64),
			strconv.FormatFloat(thissubject.AverageCorrectRT("PracTrialProc", "incorrect"), 'f', 10, 64),
			strconv.FormatFloat(thissubject.AverageAccuracy("LDTrialProc"), 'f', 10, 64),
			strconv.FormatFloat(thissubject.AverageCorrectRT("LDTrialProc", "correct"), 'f', 10, 64),
			strconv.FormatFloat(thissubject.AverageCorrectRT("LDTrialProc", "incorrect"), 'f', 10, 64),
			strconv.FormatFloat(thissubject.AverageAccuracy("TBTrialProc"), 'f', 10, 64),
			strconv.FormatFloat(thissubject.AverageCorrectRT("TBTrialProc", "correct"), 'f', 10, 64),
			strconv.FormatFloat(thissubject.AverageCorrectRT("TBTrialProc", "incorrect"), 'f', 10, 64),
			strconv.FormatFloat(thissubject.ClockChecks("TBTrialProc"), 'f', 1, 64),
			strconv.FormatFloat(thissubject.FalseAlarms("TBTrialProc"), 'f', 1, 64),
			strconv.FormatFloat(thissubject.TrialPMScore(3), 'f', 1, 64),
			strconv.FormatFloat(thissubject.TrialPMScore(6), 'f', 1, 64),
			strconv.FormatFloat(thissubject.TrialPMScore(9), 'f', 1, 64),
			strconv.FormatFloat(thissubject.TotalPMScore(), 'f', 3, 64),
		}
		writer.Write(outvars)
	}
}
