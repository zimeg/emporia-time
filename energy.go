package main

import (
	"math"
	"time"
)

type Emporia struct {
	device      string
	token       string
	resp        EmporiaUsageResp
	chart       []float64
	usage       float64
	elapsedTime time.Duration
	sureness    float64
}

type EmporiaUsageResp struct {
	Message           string
	FirstUsageInstant string
	UsageList         []float64
}

var HourToSeconds float64 = 3600
var KiloToUnit float64 = 1000

func ScaleKWhToWs(kwh float64) float64 {
	return kwh * KiloToUnit * HourToSeconds
}

// LookupEnergyUsage gathers device usage stats between the start and end times
func (e *Emporia) LookupEnergyUsage(start time.Time, end time.Time) ([]float64, error) {
	params := formatUsageParams(e.device, start, end)
	chart, err := e.getEnergyUsage(params)
	if err != nil {
		return []float64{}, err
	}

	for ii, kwh := range chart {
		chart[ii] = ScaleKWhToWs(kwh)
	}

	e.usage, e.sureness = ExtrapolateUsage(chart, e.elapsedTime.Seconds())
	return chart, nil
}

// extrapolateUsage scales the average measured energy rate over the elapsed
// time to account for missing measurements, returning estimated watts
func ExtrapolateUsage(measured []float64, durr float64) (float64, float64) {
	var sum float64 = 0
	for _, mm := range measured {
		sum += mm
	}

	// scale the summation across the entire duration
	measurements := len(measured)
	estimated := sum * (durr / float64(measurements))

	// calculate the observed-to-expected measurement ratio
	var sureness float64 = 0
	if estimated > 0.0 {
		sureness = float64(measurements) / math.Ceil(durr)
	}

	return estimated, sureness
}
