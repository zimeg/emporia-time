package main

import (
	"context"
	"log"
	"os"

	"github.com/spf13/afero"
	"github.com/zimeg/emporia-time/cmd"
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
	req := api.New()
	cog, err := cognito.NewClient(ctx, clientID, region)
	if err != nil {
		log.Fatalf("Error: %s", err)
	}
	result, err := cmd.Root(ctx, cog, fs, req, os.Args, version)
	if err != nil {
		log.Fatalf("Error: %s", err)
	}
	os.Exit(result.ExitCode)
}
