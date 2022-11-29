package main

import (
	"math"
	"net/url"
	"time"
)

type Emporia struct {
	device      string
	token       string
	chart       EmporiaUsageChart
	usage       float64
	elapsedTime time.Duration
	sureness    float64
}

type EmporiaUsageChart struct {
	Message           string
	FirstUsageInstant string
	UsageList         []float64
}

// LookupEnergyUsage gathers device usage stats between the start and end times
func (e *Emporia) LookupEnergyUsage(start time.Time, end time.Time) ([]float64, error) {

	// https://github.com/magico13/PyEmVue/blob/master/api_docs.md#getchartusage---usage-over-a-range-of-time
	params := url.Values{}
	params.Set("apiMethod", "getChartUsage")
	params.Set("deviceGid", e.device)
	params.Set("channel", "1,2,3") // ?
	params.Set("start", start.Format(time.RFC3339))
	params.Set("end", end.Format(time.RFC3339))
	params.Set("scale", "1S")
	params.Set("energyUnit", "KilowattHours")

	chart, err := e.getEnergyUsage(params)
	if err != nil {
		return []float64{}, err
	}

	_ = e.extrapolateUsage()
	return chart, nil
}

// extrapolateUsage scales the average measured energy rate over the elapsed
// time to account for missing measurements, returning estimated watts
func (e *Emporia) extrapolateUsage() float64 {
	var measuredUsage float64 = 0
	for _, uu := range e.chart.UsageList {
		measuredUsage += uu * 3600 * 1000 // convert kWh to W
	}

	// scale the summation across the entire duration
	var seconds float64 = e.elapsedTime.Seconds()
	var measurements = len(e.chart.UsageList)
	e.usage = measuredUsage * (seconds / float64(measurements))

	// share the observed-to-expected measurement ratio
	e.sureness = 0
	if e.usage > 0.0 {
		e.sureness = float64(measurements) / math.Ceil(seconds)
	}

	return e.usage
}
