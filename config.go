package main

import (
	"encoding/json"
	"errors"
	"log"
	"os"
	"time"
)

type EmporiaConfig struct {
	EmporiaDevice  string
	EmporiaToken   string
	EmporiaRefresh string
	EmporiaExpires time.Time
}

// Init collects and sets the tokens needed when calling the Emporia API
func (e *Emporia) Init() {
	config := new(EmporiaConfig)
	config.LoadConfig()

	// generate new authentication tokens or refresh expired tokens
	if config.EmporiaToken == "" || config.EmporiaRefresh == "" {
		username, password := collectCredentials()
		tokens := GenerateTokens(username, password)
		config.SaveTokens(tokens)

	} else if time.Now().After(config.EmporiaExpires) {
		tokens := RefreshTokens(config.EmporiaRefresh)
		config.SaveTokens(tokens)
	}

	// select an available device
	if config.EmporiaDevice == "" {
		devices := getAvailableDevices(config.EmporiaToken)
		config.EmporiaDevice = selectDevice(devices)
	}

	config.SaveConfig()
	e.config = config
}

// LoadConfig unmarshals stored config data
func (conf *EmporiaConfig) LoadConfig() {
	configFilePath := findConfigFilePath()
	data, err := os.ReadFile(configFilePath)
	if errors.Is(err, os.ErrNotExist) {
		return
	}

	if err != nil {
		log.Panicf("Failed to read config file: %s\n", err)
	}

	if len(data) > 0 {
		err := json.Unmarshal(data, &conf)
		if err != nil {
			log.Panicf("Failed to parse config file: %s\n", err)
		}
	}
}

// SaveConfig saves config data to the config file
func (conf *EmporiaConfig) SaveConfig() {
	configFilePath := findConfigFilePath()
	data, err := json.MarshalIndent(conf, "", "\t")
	if err != nil {
		log.Panicf("Failed to encode config data: %s\n", err)
	}

	err = os.WriteFile(configFilePath, data, 0660)
	if err != nil {
		log.Fatal(err)
	}
}

// findConfigFilePath returns the path to stored local credentials
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
