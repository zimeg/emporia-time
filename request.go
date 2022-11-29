package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

var EmporiaBaseURL = "https://api.emporiaenergy.com"

// getEnergyUsage performs a GET request to `/AppAPI` with configured params
func (e *Emporia) getEnergyUsage(params url.Values) ([]float64, error) {
	EmporiaURL := fmt.Sprintf("%s/AppAPI?%s", EmporiaBaseURL, params.Encode())

	client := &http.Client{}
	req, err := http.NewRequest("GET", EmporiaURL, nil)
	req.Header.Add("authToken", e.token)

	resp, err := client.Do(req)
	if err != nil {
		return []float64{}, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)

	err = json.Unmarshal(body, &e.chart)
	if err != nil {
		return []float64{}, err
	}
	if e.chart.Message != "" {
		return []float64{}, errors.New(e.chart.Message)
	}

	return e.chart.UsageList, nil
}
