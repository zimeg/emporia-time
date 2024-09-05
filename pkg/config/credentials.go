package config

import (
	"context"
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/zimeg/emporia-time/internal/terminal"
	"github.com/zimeg/emporia-time/pkg/cognito"
)

// Credentials contains basic authentication information for gathering tokens
type Credentials struct {
	Username string
	Password string
}

// SetTokens stores newly gathered auth tokens in the config
func (config *Config) SetTokens(auth cognito.CognitoResponse) {
	if auth.IdToken != nil {
		token := *auth.IdToken
		config.Tokens.IdToken = token
		config.req.SetToken(token)
	} else {
		config.req.SetToken(config.Tokens.IdToken)
		return
	}
	if auth.RefreshToken != nil {
		config.Tokens.RefreshToken = *auth.RefreshToken
	}
	config.Tokens.ExpiresAt = time.Now().
		Add(time.Second * time.Duration(auth.ExpiresIn)).UTC()
}

// GetTokens gathers valid authentication tokens from provided configurations
func (cfg *Config) GetTokens(
	ctx context.Context,
	cog cognito.Cognitoir,
	flags Flags,
) (
	cognito.CognitoResponse,
	error,
) {
	if cfg.useCredentials(flags) {
		username, password, err := cfg.gatherCredentials(flags)
		if err != nil {
			return cognito.CognitoResponse{}, err
		}
		tokens, err := cog.GenerateTokens(ctx, username, password)
		if err != nil {
			return cognito.CognitoResponse{}, err
		}
		return tokens, nil
	}
	if time.Now().After(cfg.Tokens.ExpiresAt) {
		tokens, err := cog.RefreshTokens(ctx, cfg.Tokens.RefreshToken)
		if err != nil {
			return cognito.CognitoResponse{}, err
		}
		return tokens, nil
	}
	return cognito.CognitoResponse{}, nil
}

// useCredentials returns if new login credentials should be used
func (cfg *Config) useCredentials(flags Flags) bool {
	return (cfg.Tokens.IdToken == "" || cfg.Tokens.RefreshToken == "") ||
		(flags.Username != "" || os.Getenv("EMPORIA_USERNAME") != "") ||
		(flags.Password != "" || os.Getenv("EMPORIA_PASSWORD") != "")
}

// headlessLogin returns if all credentials are provided by flag or environment
func (cfg *Config) headlessLogin(flags Flags) bool {
	return (flags.Username != "" || os.Getenv("EMPORIA_USERNAME") != "") &&
		(flags.Password != "" || os.Getenv("EMPORIA_PASSWORD") != "")
}

// gatherCredentials prompts for an Emporia username and password
func (cfg *Config) gatherCredentials(
	flags Flags,
) (
	username string,
	password string,
	err error,
) {
	if !cfg.headlessLogin(flags) {
		fmt.Printf("Enter your Emporia credentials <https://web.emporiaenergy.com/>\n")
	}
	username, err = terminal.CollectInput(&terminal.Prompt{
		Message:     "Username",
		Flag:        flag.Lookup("username"),
		Environment: "EMPORIA_USERNAME",
	})
	if err != nil {
		return "", "", err
	}
	password, err = terminal.CollectInput(&terminal.Prompt{
		Message:     "Password",
		Flag:        flag.Lookup("password"),
		Environment: "EMPORIA_PASSWORD",
		Hidden:      true,
	})
	if err != nil {
		return "", "", err
	}
	return username, password, nil
}
