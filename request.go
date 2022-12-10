package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
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

	err = json.Unmarshal(body, &e.resp)
	if err != nil {
		return []float64{}, err
	}
	if e.resp.Message != "" {
		return []float64{}, errors.New(e.resp.Message)
	}

	return e.resp.UsageList, nil
}

func formatUsageParams(device string, start time.Time, end time.Time) url.Values {

	// https://github.com/magico13/PyEmVue/blob/master/api_docs.md#getchartusage---usage-over-a-range-of-time
	params := url.Values{}
	params.Set("apiMethod", "getChartUsage")
	params.Set("deviceGid", device)
	params.Set("channel", "1,2,3") // ?
	params.Set("start", start.Format(time.RFC3339))
	params.Set("end", end.Format(time.RFC3339))
	params.Set("scale", "1S")
	params.Set("energyUnit", "KilowattHours")

	return params
}

// EmporiaStatus returns if the Emporia API is available
func EmporiaStatus() (bool, error) {

	// https://github.com/magico13/PyEmVue/blob/master/api_docs.md#detection-of-maintenance
	EmporiaStatusURL := "https://s3.amazonaws.com/com.emporiaenergy.manual.ota/maintenance/maintenance.json"

	resp, err := http.Get(EmporiaStatusURL)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	status := resp.StatusCode == 403
	return status, nil
}
