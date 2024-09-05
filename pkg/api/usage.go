package api

import (
	"errors"
	"fmt"
	"log"
	"net/url"
	"time"

	"github.com/zimeg/emporia-time/pkg/energy"
	"github.com/zimeg/emporia-time/pkg/times"
)

// UsageResponse holds usage information from the response
type UsageResponse struct {
	Message           string
	FirstUsageInstant string
	UsageList         []float64
}

// GetChartUsage calls the Emporia API for usage information over and over until
// a certain confidence is reached
func (emp *Emporia) GetChartUsage(times times.TimeMeasurement) (energy.EnergyResult, error) {
	confidence := 0.80

	// Delay before lookup to respect latency
	time.Sleep(200 * time.Millisecond)
	chart, err := emp.LookupEnergyUsage(times)
	if err != nil {
		log.Printf("Error: Failed to gather energy usage data!\n")
		return energy.EnergyResult{}, err
	}

	results := energy.ExtrapolateUsage(energy.EnergyMeasurement{
		Chart:    chart,
		Duration: times.Elapsed,
	})

	// Repeat lookup for unsure results
	for results.Sureness < confidence {
		results, err = emp.GetChartUsage(times)
		if err != nil {
			return energy.EnergyResult{}, err
		}
	}
	return results, nil
}

// LookupEnergyUsage gathers device watt usage between the start and end times
func (emp *Emporia) LookupEnergyUsage(times times.TimeMeasurement) ([]float64, error) {
	response, err := emp.getEnergyUsage(times)
	if err != nil {
		return []float64{}, err
	} else if response.Message != "" {
		return []float64{}, errors.New(response.Message)
	}
	chart := response.UsageList
	for ii, kwh := range chart {
		chart[ii] = energy.ScaleKWhToWs(kwh)
	}
	return chart, nil
}

// getEnergyUsage performs a GET request to collect usage statistics
//
// https://github.com/magico13/PyEmVue/blob/master/api_docs.md#getchartusage---usage-over-a-range-of-time
func (emp *Emporia) getEnergyUsage(times times.TimeMeasurement) (UsageResponse, error) {
	response := UsageResponse{}
	params := url.Values{
		"apiMethod":  []string{"getChartUsage"},
		"deviceGid":  []string{emp.deviceID},
		"channel":    []string{"1,2,3"}, // ?
		"start":      []string{times.Start.Format(time.RFC3339)},
		"end":        []string{times.End.Format(time.RFC3339)},
		"scale":      []string{"1S"},
		"energyUnit": []string{"KilowattHours"},
	}
	url := fmt.Sprintf("%s/AppAPI?%s", RequestURL, params.Encode())
	err := emp.get(url, &response)
	if err != nil {
		return UsageResponse{}, err
	}
	return response, nil
}
