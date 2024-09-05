package config

import (
	"context"
	"encoding/json"
	"errors"
	"os"
	"time"

	"github.com/spf13/afero"
	"github.com/zimeg/emporia-time/pkg/api"
	"github.com/zimeg/emporia-time/pkg/cognito"
)

// Configure reveals what configurations to the changing settings
type Configure interface {
	API() api.Emporiac
}

// Config contains device configurations and user authentications
type Config struct {
	Device string   // Device is the specific machine to be measured
	Tokens TokenSet // Tokens contain the authentication information

	fs   afero.Fs     // fs wraps abstraction over a stable file system
	path string       // path is the location of the configuration file
	req  api.Emporiac // req is an HTTP client with some configurations
}

// TokenSet contains authentication information saved for cognito
type TokenSet struct {
	IdToken      string
	RefreshToken string
	ExpiresAt    time.Time
}

// Load collects configurations into a single structure
func Load(
	ctx context.Context,
	cog cognito.Cognitoir,
	fs afero.Fs,
	req api.Emporiac,
	flags Flags,
) (
	cfg Config,
	err error,
) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return Config{}, err
	}
	configDir := homeDir + "/.config/etime"
	if val, ok := os.LookupEnv("XDG_CONFIG_HOME"); ok {
		configDir = val + "/etime"
	}
	err = fs.MkdirAll(configDir, 0o755)
	if err != nil {
		return Config{}, err
	}
	path := configDir + "/settings.json"
	data, err := afero.ReadFile(fs, path)
	if err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			return Config{}, err
		}
	} else if len(data) > 0 {
		err := json.Unmarshal(data, &cfg)
		if err != nil {
			return Config{}, err
		}
	}
	cfg.fs = fs
	cfg.path = path
	cfg.req = req
	defer func() {
		if err != nil {
			return
		}
		err = cfg.save()
	}()
	tokens, err := cfg.GetTokens(ctx, cog, flags)
	if err != nil {
		return Config{}, err
	} else {
		cfg.SetTokens(tokens)
	}
	device, err := cfg.GetDevice(flags)
	if err != nil {
		return Config{}, err
	} else {
		cfg.SetDevice(device)
	}
	return cfg, nil
}

// API returns a client for making requests
func (cfg *Config) API() api.Emporiac {
	return cfg.req
}

// Save writes configuration data to the settings file
func (cfg *Config) save() error {
	data, err := json.MarshalIndent(cfg, "", "\t")
	if err != nil {
		return err
	}
	err = afero.WriteFile(cfg.fs, cfg.path, data, 0o660)
	if err != nil {
		return err
	}
	return nil
}
