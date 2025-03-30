package api

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/zimeg/emporia-time/internal/errors"
	"github.com/zimeg/emporia-time/pkg/energy"
	"github.com/zimeg/emporia-time/pkg/times"
)

// Emporiac communicates with the Emporia API
type Emporiac interface {
	GetChartUsage(times times.TimeMeasurement) (results energy.EnergyResult, err error)
	GetCustomerDevices() (devices []Device, err error)
	Status() (available bool, err error)

	SetDevice(deviceID string)
	SetToken(token string)
}

// Emporia holds information for and from the Emporia API
type Emporia struct {
	client interface {
		// Do does the HTTP request and is often implmented using net/http
		Do(req *http.Request) (*http.Response, error)
	}
	deviceID string
	token    string
}

// RequestURL is the base URL of the API
const RequestURL string = "https://api.emporiaenergy.com"

// New creates a new client to interact with Emporia HTTP APIs
func New() *Emporia {
	return &Emporia{
		client: &http.Client{},
	}
}

// SetToken sets the token for the client
func (emp *Emporia) SetToken(token string) {
	emp.token = token
}

// SetDevice sets the device ID for the client
func (emp *Emporia) SetDevice(deviceID string) {
	emp.deviceID = deviceID
}

// get makes an authenticated GET request to URL and saves the response to data
func (emp *Emporia) get(url string, data any) error {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return errors.Wrap(errors.ErrEmporiaRequest, err)
	}
	if emp.token != "" {
		req.Header.Add("authToken", emp.token)
	}
	resp, err := emp.client.Do(req)
	if err != nil {
		return errors.Wrap(errors.ErrEmporiaResponse, err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return errors.Wrap(errors.ErrEmporiaResult, err)
	}
	err = json.Unmarshal(body, &data)
	if err != nil {
		return errors.Wrap(errors.ErrEmporiaFormat, err)
	}
	return nil
}
