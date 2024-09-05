package api

import (
	"fmt"
)

// Device represents a device that can be measured
type Device struct {
	DeviceGid          int
	LocationProperties struct {
		DeviceName string
	}
}

// DeviceResponse contains a slice of available devices
type DeviceResponse struct {
	Devices []Device
	Message string
}

// GetCustomerDevices returns customer devices for the Emporia account
func (emp *Emporia) GetCustomerDevices() ([]Device, error) {
	response := DeviceResponse{}
	url := fmt.Sprintf("%s/customers/devices", RequestURL)
	err := emp.get(url, &response)
	if err != nil {
		return []Device{}, err
	}
	if response.Message != "" {
		return []Device{}, fmt.Errorf("%s", response.Message)
	}
	return response.Devices, nil
}
