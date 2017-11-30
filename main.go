// TBPM-FION is an eprime extraction utilitity for a specific experimental design
// Run by the neuropsychology department at UWA.

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
var pStatus = regexp.MustCompile(".*?([0-9]{4}).*?")

func getParticipantStatus(filename string) (status string) {
	statusnum := pStatus.FindStringSubmatch(filename)[1][0]
	switch statusnum {
	case "1"[0]:
		status = "ABI"
	case "2"[0]:
		status = "Control"
	default:
		status = "UnableToGenerate"
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
	writer.Write([]string{"id", "status", "testdate",
		"Accuacy_PracTrialProc", "CorrectRT_PracTrialProc", "IncorrectRT_PracTrialProc", "TotalNonResponses_PracTrialProc",
		"Accuacy_LDTrialProc", "CorrectRT_LDTrialProc", "IncorrectRT_LDTrialProc", "TotalNonResponses_LDTrialProc",
		"Accuacy_TBTrialProc", "CorrectRT_TBTrialProc", "IncorrectRT_TBTrialProc", "TotalNonResponses_TBTrialProc", "ClockChecks_TOTAL", "FalseAlarm_TOTAL", "PM_3min", "PM_7min", "PM_9min", "PM_TOTAL"})

	// For all arguments, process file and then write out results
	for _, file := range os.Args[1:] {
		thissubject, err := processFile(file)
		checkError("Error reading participant file", err)

		thissubject.status = getParticipantStatus(strings.ToLower(file))
		outvars := []string{
			thissubject.id,
			thissubject.status,
			thissubject.sessiondate,
			strconv.FormatFloat(thissubject.AverageAccuracy("PracTrialProc"), 'f', 4, 64),
			strconv.FormatFloat(thissubject.AverageRT("PracTrialProc", "correct"), 'f', 4, 64),
			strconv.FormatFloat(thissubject.AverageRT("PracTrialProc", "incorrect"), 'f', 4, 64),
			strconv.FormatFloat(thissubject.TotalNonResponseTrials("PracTrialProc"), 'f', 1, 64),
			strconv.FormatFloat(thissubject.AverageAccuracy("LDTrialProc"), 'f', 4, 64),
			strconv.FormatFloat(thissubject.AverageRT("LDTrialProc", "correct"), 'f', 4, 64),
			strconv.FormatFloat(thissubject.AverageRT("LDTrialProc", "incorrect"), 'f', 4, 64),
			strconv.FormatFloat(thissubject.TotalNonResponseTrials("LDTrialProc"), 'f', 1, 64),
			strconv.FormatFloat(thissubject.AverageAccuracy("TBTrialProc"), 'f', 4, 64),
			strconv.FormatFloat(thissubject.AverageRT("TBTrialProc", "correct"), 'f', 4, 64),
			strconv.FormatFloat(thissubject.AverageRT("TBTrialProc", "incorrect"), 'f', 4, 64),
			strconv.FormatFloat(thissubject.TotalNonResponseTrials("TBTrialProc"), 'f', 1, 64),
			strconv.FormatFloat(thissubject.ClockChecks("TBTrialProc"), 'f', 1, 64),
			strconv.FormatFloat(thissubject.FalseAlarms("TBTrialProc"), 'f', 1, 64),
			strconv.FormatFloat(thissubject.TrialPMScore(3), 'f', 1, 64),
			strconv.FormatFloat(thissubject.TrialPMScore(7), 'f', 1, 64),
			strconv.FormatFloat(thissubject.TrialPMScore(9), 'f', 1, 64),
			strconv.FormatFloat(thissubject.TotalPMScore(), 'f', 3, 64),
		}
		writer.Write(outvars)
	}
}
