package emporia

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/zimeg/emporia-time/internal/program"
	"github.com/zimeg/emporia-time/internal/terminal"
)

// EmporiaCredentials contains basic authentication information
type EmporiaCredentials struct {
	Username string
	Password string
}

// headlessLogin returns if all credentials are provided by flag or environment
func (config *EmporiaConfig) headlessLogin(flags program.Flags) bool {
	return (flags.Username != "" || os.Getenv("EMPORIA_USERNAME") != "") &&
		(flags.Password != "" || os.Getenv("EMPORIA_PASSWORD") != "")
}

// useCredentials returns if new login credentials should be used
func (config *EmporiaConfig) useCredentials(flags program.Flags) bool {
	return (config.Tokens.IdToken == "" || config.Tokens.RefreshToken == "") ||
		(flags.Username != "" || os.Getenv("EMPORIA_USERNAME") != "") ||
		(flags.Password != "" || os.Getenv("EMPORIA_PASSWORD") != "")
}

// gatherTokens collects and sets the tokens needed for calling the Emporia API
func (config *EmporiaConfig) gatherTokens(flags program.Flags) error {
	if config.useCredentials(flags) {
		if credentials, err := config.gatherCredentials(flags); err != nil {
			return err
		} else if resp, err := GenerateTokens(credentials); err != nil {
			return err
		} else {
			config.SetTokens(resp)
		}
	} else if time.Now().After(config.Tokens.ExpiresAt) {
		if resp, err := RefreshTokens(config.Tokens.RefreshToken); err != nil {
			return err
		} else {
			config.SetTokens(resp)
		}
	}
	return nil
}

// gatherCredentials prompts for an Emporia username and password
func (config *EmporiaConfig) gatherCredentials(flags program.Flags) (EmporiaCredentials, error) {
	credentials := EmporiaCredentials{}
	if !config.headlessLogin(flags) {
		fmt.Printf("Enter your Emporia credentials <https://web.emporiaenergy.com/>\n")
	}
	if username, err := terminal.CollectInput(&terminal.Prompt{
		Message:     "Username",
		Flag:        flag.Lookup("username"),
		Environment: "EMPORIA_USERNAME",
	}); err != nil {
		return EmporiaCredentials{}, err
	} else {
		credentials.Username = username
	}
	if password, err := terminal.CollectInput(&terminal.Prompt{
		Message:     "Password",
		Flag:        flag.Lookup("password"),
		Environment: "EMPORIA_PASSWORD",
		Hidden:      true,
	}); err != nil {
		return EmporiaCredentials{}, err
	} else {
		credentials.Password = password
	}
	return credentials, nil
}

// gatherDevice prompts and stores the choice of an Emporia device
func (config *EmporiaConfig) gatherDevice(flags program.Flags) error {
	var names, gids, gidLabels = []string{}, []string{}, []string{}
	var device string

	devices := getAvailableDevices(config.Tokens.IdToken)
	if len(devices) == 0 {
		return errors.New("No available devices found!")
	}

	switch {
	case flags.Device != "":
		device = flags.Device
	case os.Getenv("EMPORIA_DEVICE") != "":
		device = os.Getenv("EMPORIA_DEVICE")
	case config.Device != "":
		device = config.Device
	}

	for _, val := range devices {
		deviceGid := strconv.Itoa(val.DeviceGid)
		if deviceGid == device || val.LocationProperties.DeviceName == device {
			config.SetDevice(deviceGid)
			return nil
		}
		names = append(names, val.LocationProperties.DeviceName)
		gids = append(gids, deviceGid)
		gidLabels = append(gidLabels, fmt.Sprintf("#%d", val.DeviceGid))
	}
	if device != "" {
		return errors.New("No matching device found!")
	}

	if selection, err := terminal.CollectSelect(terminal.Prompt{
		Message:      "Select a device:",
		Options:      names,
		Descriptions: gidLabels,
	}); err != nil {
		return err
	} else {
		device = gids[selection]
		config.SetDevice(device)
	}
	return nil
}
