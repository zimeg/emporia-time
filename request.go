package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/AlecAivazis/survey/v2"
)

var EmporiaBaseURL = "https://api.emporiaenergy.com"

type Emporia struct {
	resp   EmporiaUsageResp
	config *EmporiaConfig
}

type EmporiaUsageResp struct {
	Message           string
	FirstUsageInstant string
	UsageList         []float64
}

type EmporiaDeviceResp struct {
	Devices []EmporiaDevice
}

type EmporiaDevice struct {
	DeviceGid          int
	LocationProperties struct {
		DeviceName string
	}
}

// LookupEnergyUsage gathers device watt usage between the start and end times
func (e *Emporia) LookupEnergyUsage(start time.Time, end time.Time) ([]float64, error) {
	params := formatUsageParams(e.config.EmporiaDevice, start, end)
	chart, err := e.getEnergyUsage(params)
	if err != nil {
		return []float64{}, err
	}

	for ii, kwh := range chart {
		chart[ii] = ScaleKWhToWs(kwh)
	}

	return chart, nil
}

// formatUsageParams returns URL values for the API
func formatUsageParams(device string, start time.Time, end time.Time) url.Values {
	params := url.Values{}

	// https://github.com/magico13/PyEmVue/blob/master/api_docs.md#getchartusage---usage-over-a-range-of-time
	params.Set("apiMethod", "getChartUsage")
	params.Set("deviceGid", device)
	params.Set("channel", "1,2,3") // ?
	params.Set("start", start.Format(time.RFC3339))
	params.Set("end", end.Format(time.RFC3339))
	params.Set("scale", "1S")
	params.Set("energyUnit", "KilowattHours")

	return params
}

// getEnergyUsage performs a GET request to `/AppAPI` with configured params
func (e *Emporia) getEnergyUsage(params url.Values) ([]float64, error) {
	EmporiaURL := fmt.Sprintf("%s/AppAPI?%s", EmporiaBaseURL, params.Encode())

	client := &http.Client{}
	req, err := http.NewRequest("GET", EmporiaURL, nil)
	req.Header.Add("authToken", e.config.EmporiaToken)

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

// TODO refactor, split, remove conf
func (conf *EmporiaConfig) selectAvailableDevices() string {

	// gather device info
	EmporiaDeviceURL := EmporiaBaseURL + "/customers/devices"

	client := &http.Client{}
	req, err := http.NewRequest("GET", EmporiaDeviceURL, nil)
	req.Header.Add("authToken", conf.EmporiaToken)

	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("Failed to gather device information: %s\n", err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)

	var devs EmporiaDeviceResp
	err = json.Unmarshal(body, &devs)
	if err != nil {
		log.Fatalf("Failed to parse device information: %s\n", err)
	}

	// select a device
	if len(devs.Devices) == 0 {
		log.Fatalf("No devices found")
	}

	var devices []string
	var gids []int

	for _, val := range devs.Devices {
		devices = append(devices, val.LocationProperties.DeviceName)
		gids = append(gids, val.DeviceGid)
	}

	var selected string
	prompt := &survey.Select{
		Message: "Select a device:",
		Options: devices,
		Description: func(value string, index int) string {
			return fmt.Sprintf("#%d", gids[index])
		},
	}
	survey.AskOne(prompt, &selected)
	fmt.Printf("\n")

	var gid int
	for index, val := range devices {
		if val == selected {
			gid = gids[index]
		}
	}

	return fmt.Sprintf("%d", gid)
}
