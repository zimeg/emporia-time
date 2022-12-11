package main

import (
	"encoding/json"
	"log"
	"os"
)

type EmporiaConfig struct {
	EmporiaDevice string
	EmporiaToken  string
}

func (e *Emporia) Init() {
	config := new(EmporiaConfig)
	config.ParseConfig()

	// TODO update expired EmporiaToken

	config.SaveConfig()
	e.config = config
}

// ParseConfig unmarshals stored config data
func (c *EmporiaConfig) ParseConfig() {
	configFilePath := findConfigFilePath()
	data, err := os.ReadFile(configFilePath)
	if err != nil {
		log.Panicf("Failed to read config file: %s\n", err)
	}

	if len(data) > 0 {
		err := json.Unmarshal(data, &c)
		if err != nil {
			log.Panicf("Failed to parse config file: %s\n", err)
		}
	}
}

// SaveConfig saves config data to the config file
func (c *EmporiaConfig) SaveConfig() {
	configFilePath := findConfigFilePath()
	data, err := json.MarshalIndent(c, "", "\t")
	if err != nil {
		log.Panicf("Failed to encode config data: %s\n", err)
	}

	err = os.WriteFile(configFilePath, data, 0660)
	if err != nil {
		log.Fatal(err)
	}
}

// findConfigFilePath returns a path to stored local credentials
func findConfigFilePath() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Panicf("Failed to find home directory: %s\n", err)
	}

	configDir := homeDir + "/.config/etime"
	if val, ok := os.LookupEnv("XDG_CONFIG_HOME"); ok {
		configDir = val + "/etime"
	}

	err = os.MkdirAll(configDir, 0755)
	if err != nil {
		log.Panicf("Failed to create config directory: %s\n", err)
	}

	configFile := configDir + "/settings.json"
	file, err := os.OpenFile(configFile, os.O_CREATE, 0600)
	if err != nil {
		log.Panicf("Failed to open file: %s\n", err)
	}
	defer file.Close()
	return configFile
}
