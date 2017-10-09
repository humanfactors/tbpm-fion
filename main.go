package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/Sirupsen/logrus"
	"golang.org/x/text/encoding/unicode"
	"golang.org/x/text/transform"
)

var log = logrus.New()

func checkError(message string, err error) {
	if err != nil {
		log.Fatal(message, err)
	}
}

type Trial struct {
	subject       string
	procedure     string
	trialnumber   string
	stimulus      string
	accuracy      string
	correctanswer string
	response      string
	rt            string
	stimuli       string
}

var allkeyvars = []string{"TBTrialList:", "LDTrialList:", "PracTrialList:", "Procedure:", "Stimulus2.ACC:", "Stimulus.ACC:", "Stimulus.RESP:", "Stimulus.CRESP:", "Stimulus2.RESP:", "Stimulus2.CRESP:", "Stimuli:", "Stimulus.RT:", "Stimulus2.RT:"}
var searchmap = map[string]string{
	"Procedure:":       "procedure",
	"TBTrialList:":     "trialnumber",
	"LDTrialList:":     "trialnumber",
	"PracTrialList:":   "trialnumber",
	"Stimulus.RT:":     "rt",
	"Stimulus2.RT:":    "rt",
	"Stimulus2.ACC:":   "accuracy",
	"Stimulus.ACC:":    "accuracy",
	"Stimulus.RESP:":   "response",
	"Stimulus2.RESP:":  "response",
	"Stimulus.CRESP:":  "correctanswer",
	"Stimulus2.CRESP:": "correctanswer",
	"Stimuli:":         "stimuli",
}

var subject string
var logframe = "LogFrame End"

func MatchVars(line string, varlist []string) (match string, found bool) {
	/*	if match in line, return match and true else "" and false	*/
	for _, hvar := range varlist {
		if strings.Contains(line, hvar) {
			return hvar, true
		}
	}
	return "", false
}

func main() {

	// Setup logging output file
	logfile, err := os.OpenFile("debugging.log", os.O_CREATE|os.O_WRONLY, 0666)
	if err == nil {
		log.Out = logfile
	} else {
		log.Info("Failed to log to file, using default stderr")
	}

	// Log file processed and generate output name
	filename := fmt.Sprintf("processed-%s.csv", time.Now().Format("2006-01-02"))
	log.Info("[Processing : %s]", filename)
	
	// Create output file writer
	file, err := os.Create(filename)
	checkError("Cannot create file", err)
	defer file.Close()
	writer := csv.NewWriter(file)
	defer writer.Flush()
	writer.Write([]string{"subj", "procedure", "trial", "stimulus", "correctanswer", "selectedresponse", "correct", "rt"})

	// For each file dragged onto the executable
	for _, file := range os.Args {
		fmt.Println("Processing subject %s", file)

		data, err := NewScannerUTF16(file)
		checkError("Data file not in UTF16", err)

		scanner := bufio.NewScanner(data)
		if err := scanner.Err(); err != nil {
			fmt.Fprintln(os.Stderr, "reading inputfile:", err)
		}

		thistrial := &Trial{}

		for scanner.Scan() {
			line := scanner.Text()
			if strings.Contains(line, "Subject:") {
				subject = strings.Fields(line)[len(strings.Fields(line))-1]
				thistrial.subject = subject
			}

			if strings.Contains(line, ":") {

				split := strings.Fields(line)
				value := split[len(split)-1]
				if len(split) < 2 {
					value = "NA"
				}

				linevar, linebool := MatchVars(line, allkeyvars)

				if linebool { // if key variable found

					thistrial.subject = subject

					if strings.Contains(searchmap[linevar], "procedure") {
						thistrial.procedure = value
					} 
					if strings.Contains(searchmap[linevar], "stimulus") {
						thistrial.stimulus = value
					}
					if strings.Contains(searchmap[linevar], "correctanswer") {
						thistrial.correctanswer = value
					}
					if strings.Contains(searchmap[linevar], "stimuli") {
						thistrial.stimuli = value
					}
					if strings.Contains(searchmap[linevar], "response") {
						thistrial.response = value
					}
					if strings.Contains(searchmap[linevar], "accuracy") {
						thistrial.accuracy = value
					}
					if strings.Contains(searchmap[linevar], "trialnumber") {
						thistrial.trialnumber = value
					}
					if strings.Contains(searchmap[linevar], "rt") {
						thistrial.rt = value
					}

				}

			}
			if strings.Contains(line, logframe) {
				if thistrial.trialnumber != "" {
					if thistrial.rt == "" {
						thistrial.rt = "NA"
					}
					out := []string{
						thistrial.subject,
						thistrial.procedure,
						thistrial.trialnumber,
						thistrial.stimuli,
						thistrial.correctanswer,
						thistrial.response,
						thistrial.accuracy,
						thistrial.rt,
					}
					writer.Write(out)
				}
				thistrial = &Trial{}
			}
		}
	}

}

type utfScanner interface {
	Read(p []byte) (n int, err error)
}

func NewScannerUTF16(filename string) (utfScanner, error) {

	// Read the file into a []byte:
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}

	// Make an tranformer that converts MS-Win default to UTF8:
	win16be := unicode.UTF16(unicode.BigEndian, unicode.IgnoreBOM)
	// Make a transformer that is like win16be, but abides by BOM:
	utf16bom := unicode.BOMOverride(win16be.NewDecoder())

	// Make a Reader that uses utf16bom:
	unicodeReader := transform.NewReader(file, utf16bom)
	return unicodeReader, nil
}
