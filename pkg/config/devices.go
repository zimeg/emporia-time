package config

import (
	"errors"
	"fmt"
	"os"
	"strconv"

	"github.com/zimeg/emporia-time/internal/terminal"
)

// Device is the machine being measured
type Device struct {
	DeviceID string // DeviceID is the unique numeric identifier
}

// SetDevice stores the active device throughout configurations
func (cfg *Config) SetDevice(device Device) {
	deviceID := device.DeviceID
	cfg.Device = deviceID
	cfg.req.SetDevice(deviceID)
}

// gatherDevice prompts and stores the choice of an Emporia device
func (cfg *Config) GetDevice(flags Flags) (Device, error) {
	names, gids, gidLabels := []string{}, []string{}, []string{}
	device := ""
	devices, err := cfg.API().GetCustomerDevices()
	if err != nil {
		return Device{}, err
	}
	if len(devices) == 0 {
		return Device{}, errors.New("No available devices found!")
	}
	switch {
	case flags.Device != "":
		device = flags.Device
	case os.Getenv("EMPORIA_DEVICE") != "":
		device = os.Getenv("EMPORIA_DEVICE")
	case cfg.Device != "":
		device = cfg.Device
	}
	for _, val := range devices {
		deviceGid := strconv.Itoa(val.DeviceGid)
		if deviceGid == device || val.LocationProperties.DeviceName == device {
			response := Device{
				DeviceID: strconv.Itoa(val.DeviceGid),
			}
			return response, nil
		}
		names = append(names, val.LocationProperties.DeviceName)
		gids = append(gids, deviceGid)
		gidLabels = append(gidLabels, fmt.Sprintf("#%d", val.DeviceGid))
	}
	if device != "" {
		return Device{}, errors.New("No matching device found!")
	}
	selection, err := terminal.CollectSelect(terminal.Prompt{
		Message:      "Select a device:",
		Options:      names,
		Descriptions: gidLabels,
	})
	if err != nil {
		return Device{}, err
	}
	response := Device{
		DeviceID: gids[selection],
	}
	return response, nil
}
