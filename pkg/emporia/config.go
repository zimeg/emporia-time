package emporia

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"os"
	"time"

	"github.com/zimeg/emporia-time/internal/program"
	"github.com/zimeg/emporia-time/pkg/cognito"
)

const (
	// EmporiaClientID is the AWS Cognito client ID used by Emporia
	EmporiaClientID string = "4qte47jbstod8apnfic0bunmrq"

	// EmporiaRegion is the AWS Cognito region used by Emporia
	EmporiaRegion string = "us-east-2"
)

// EmporiaConfig contains device configurations and user authentications
type EmporiaConfig struct {
	Device string
	Tokens EmporiaTokenSet
}

// EmporiaTokenSet contains authentication information needed by Emporia
type EmporiaTokenSet struct {
	IdToken      string
	RefreshToken string
	ExpiresAt    time.Time
}

// SetupConfig prepares the local configurations for a command
func SetupConfig(flags program.Flags) (EmporiaConfig, error) {
	ctx := context.Background()
	cog, err := cognito.NewClient(ctx, EmporiaClientID, EmporiaRegion)
	if err != nil {
		return EmporiaConfig{}, err
	}
	config, err := LoadConfigFile()
	if err != nil {
		return EmporiaConfig{}, err
	}
	tokens, err := config.gatherTokens(ctx, cog, flags)
	if err != nil {
		return EmporiaConfig{}, err
	} else {
		config.SetTokens(tokens)
	}
	err = config.gatherDevice(flags)
	if err != nil {
		return EmporiaConfig{}, err
	}
	config.SaveConfig()
	return config, nil
}

// SetDevice stores the active device in the config
func (config *EmporiaConfig) SetDevice(device string) {
	config.Device = device
}

// SetTokens stores newly gathered auth tokens in the config
func (config *EmporiaConfig) SetTokens(auth cognito.CognitoResponse) {
	config.Tokens.IdToken = *auth.IdToken
	if auth.RefreshToken != nil {
		config.Tokens.RefreshToken = *auth.RefreshToken
	}
	lifespan := time.Duration(auth.ExpiresIn)
	config.Tokens.ExpiresAt = time.Now().Add(time.Second * lifespan).UTC()
}

// LoadConfigFile unmarshals stored config data
func LoadConfigFile() (EmporiaConfig, error) {
	var config EmporiaConfig
	configFilePath := findConfigFilePath()
	if data, err := os.ReadFile(configFilePath); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return EmporiaConfig{}, nil
		}
		return EmporiaConfig{}, err
	} else if len(data) > 0 {
		if err := json.Unmarshal(data, &config); err != nil {
			return EmporiaConfig{}, err
		}
		return config, nil
	}
	return EmporiaConfig{}, nil
}

// SaveConfig saves config data to the config file
func (config *EmporiaConfig) SaveConfig() {
	configFilePath := findConfigFilePath()
	if data, err := json.MarshalIndent(config, "", "\t"); err != nil {
		log.Panicf("Failed to encode config data: %s\n", err)
	} else if err := os.WriteFile(configFilePath, data, 0o660); err != nil {
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
	if err := os.MkdirAll(configDir, 0o755); err != nil {
		log.Panicf("Failed to create config directory: %s\n", err)
	}

	configFile := configDir + "/settings.json"
	if file, err := os.OpenFile(configFile, os.O_CREATE, 0o600); err != nil {
		log.Panicf("Failed to open file: %s\n", err)
	} else {
		defer file.Close()
	}
	return configFile
}
