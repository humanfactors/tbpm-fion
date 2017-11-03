package main

import (
	"github.com/montanaflynn/stats"
)

func (p Participant) AverageCorrectRT(procedure string, responsetype string) float64 {
	var rttype float64
	if responsetype == "correct" {
		rttype = 1
	} else {
		if responsetype == "incorrect" {
			rttype = 0
		}
	}
	rts := []float64{}
	for _, trial := range p.trials {
		if trial.accuracy == rttype && trial.procedure == procedure {
			rts = append(rts, trial.rt)
		}
	}
	meanrts, _ := stats.Mean(rts)
	return meanrts
}

func (p Participant) AverageAccuracy(procedure string) float64 {
	errs := []float64{}
	for _, trial := range p.trials {
		if trial.procedure == procedure {
			errs = append(errs, trial.accuracy)
		}
	}
	meanerrs, _ := stats.Mean(errs)
	return meanerrs
}

func (p Participant) FalseAlarms(procedure string) float64 {
	fas := []float64{}
	for _, trial := range p.trials {
		if trial.procedure == procedure {
			fas = append(fas, trial.falsealarm)
		}
	}
	meanfas, _ := stats.Sum(fas)
	return meanfas
}

func (p Participant) ClockChecks(procedure string) float64 {
	clockchecks := []float64{}
	for _, trial := range p.trials {
		if trial.procedure == procedure {
			clockchecks = append(clockchecks, trial.clockcheck)
		}
	}
	meanclockchecks, _ := stats.Sum(clockchecks)
	return meanclockchecks
}

func (p Participant) TrialPMScore(pmnumber int) (pmvalue float64) {
	pmvalue = 0
	for _, trial := range p.trials {
		if pmnumber == 3 && trial.pm3 == 1 {
			pmvalue = 1
		}
		if pmnumber == 7 && trial.pm7 == 1 {
			pmvalue = 1
		}
		if pmnumber == 9 && trial.pm9 == 1 {
			pmvalue = 1
		}
	}
	return pmvalue
}

func (p Participant) TotalPMScore() (totalpm float64) {
	totalpm, _ = stats.Sum([]float64{p.TrialPMScore(3), p.TrialPMScore(7), p.TrialPMScore(9)})
	totalpm = totalpm / 3
	return totalpm
}
