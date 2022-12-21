package main

import (
	"fmt"
	"log"

	"github.com/AlecAivazis/survey/v2"
)

// collectCredentials prompts for an Emporia username and password
func collectCredentials() (string, string) {
	var username string
	var password string

	fmt.Printf("Enter your Emporia credentials <https://web.emporiaenergy.com/>\n")
	survey.AskOne(&survey.Input{Message: "Username"}, &username)
	survey.AskOne(&survey.Password{Message: "Password"}, &password)

	return username, password
}

// selectDevice prompts for a choice of an owned Emporia device
func selectDevice(devices []EmporiaDevice) string {
	var names []string
	var gids []int

	if len(devices) == 0 {
		log.Fatalf("No devices found")
	}

	for _, val := range devices {
		names = append(names, val.LocationProperties.DeviceName)
		gids = append(gids, val.DeviceGid)
	}

	var selected string
	prompt := &survey.Select{
		Message: "Select a device:",
		Options: names,
		Description: func(value string, index int) string {
			return fmt.Sprintf("#%d", gids[index])
		},
	}
	survey.AskOne(prompt, &selected)

	var gid int
	for index, val := range names {
		if val == selected {
			gid = gids[index]
		}
	}

	return fmt.Sprintf("%d", gid)
}
