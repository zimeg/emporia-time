package emporia

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/zimeg/emporia-time/pkg/energy"
	"github.com/zimeg/emporia-time/pkg/times"
)

const EmporiaBaseURL = "https://api.emporiaenergy.com"

type Emporia struct {
	Resp   EmporiaUsageResp
	Config EmporiaConfig
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

// CollectEnergyUsage repeatedly calls the Emporia API for usage information
// until a certain confidence is reached
func (emp *Emporia) CollectEnergyUsage(times times.TimeMeasurement) (energy.EnergyResult, error) {
	confidence := 0.80

	// Delay before lookup to respect latency
	time.Sleep(200 * time.Millisecond)
	chart, err := emp.LookupEnergyUsage(times)
	if err != nil {
		log.Printf("Error: Failed to gather energy usage data!\n")
		return energy.EnergyResult{}, err
	}

	var results energy.EnergyResult
	results = energy.ExtrapolateUsage(chart, times.Elapsed.Seconds())

	// Repeat lookup for unsure results
	for results.Sureness < confidence {
		results, err = emp.CollectEnergyUsage(times)
		if err != nil {
			return energy.EnergyResult{}, err
		}
	}
	return results, nil
}

// LookupEnergyUsage gathers device watt usage between the start and end times
func (emp *Emporia) LookupEnergyUsage(times times.TimeMeasurement) ([]float64, error) {
	params := emp.formatUsageParams(times)
	chart, err := emp.getEnergyUsage(params)
	if err != nil {
		return []float64{}, err
	}
	for ii, kwh := range chart {
		chart[ii] = energy.ScaleKWhToWs(kwh)
	}
	return chart, nil
}

// formatUsageParams returns URL values for the API
func (emp *Emporia) formatUsageParams(times times.TimeMeasurement) url.Values {
	params := url.Values{}

	// https://github.com/magico13/PyEmVue/blob/master/api_docs.md#getchartusage---usage-over-a-range-of-time
	params.Set("apiMethod", "getChartUsage")
	params.Set("deviceGid", emp.Config.Device)
	params.Set("channel", "1,2,3") // ?
	params.Set("start", times.Start.Format(time.RFC3339))
	params.Set("end", times.End.Format(time.RFC3339))
	params.Set("scale", "1S")
	params.Set("energyUnit", "KilowattHours")

	return params
}

// getEnergyUsage performs a GET request to `/AppAPI` with configured params
func (emp *Emporia) getEnergyUsage(params url.Values) ([]float64, error) {
	EmporiaURL := fmt.Sprintf("%s/AppAPI?%s", EmporiaBaseURL, params.Encode())

	client := &http.Client{}
	req, err := http.NewRequest("GET", EmporiaURL, nil)
	req.Header.Add("authToken", emp.Config.Tokens.IdToken)

	resp, err := client.Do(req)
	if err != nil {
		return []float64{}, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)

	if err := json.Unmarshal(body, &emp.Resp); err != nil {
		return []float64{}, err
	} else if emp.Resp.Message != "" {
		return []float64{}, errors.New(emp.Resp.Message)
	}
	return emp.Resp.UsageList, nil
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

// getAvailableDevices returns customer devices for the Emporia account
func getAvailableDevices(token string) []EmporiaDevice {
	EmporiaDeviceURL := EmporiaBaseURL + "/customers/devices"

	client := &http.Client{}
	req, err := http.NewRequest("GET", EmporiaDeviceURL, nil)
	req.Header.Add("authToken", token)

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

	return devs.Devices
}
