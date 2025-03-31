package main

import (
	"context"
	"os"

	"github.com/spf13/afero"
	"github.com/zimeg/emporia-time/cmd"
	"github.com/zimeg/emporia-time/internal/logs"
	"github.com/zimeg/emporia-time/pkg/api"
	"github.com/zimeg/emporia-time/pkg/cognito"
)

// version is the title of this current build
var version = "development"

const (
	clientID string = "4qte47jbstod8apnfic0bunmrq" // Emporia AWS Cognito client ID
	region   string = "us-east-2"                  // Emporia AWS region
)

// main manages the lifecycle of this program
func main() {
	ctx := context.Background()
	fs := afero.NewOsFs()
	logger := logs.NewLogger(os.Stderr)
	req := api.New()
	cog, err := cognito.NewClient(ctx, clientID, region)
	if err != nil {
		logger.Fatal(err)
	}
	result, err := cmd.Root(ctx, cog, fs, logger, req, os.Args, version)
	if err != nil {
		logger.Fatal(err)
	}
	os.Exit(result.Command.Code)
}
