package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/spf13/afero"
	"github.com/zimeg/emporia-time/cmd/etime"
	"github.com/zimeg/emporia-time/internal/display/templates"
	"github.com/zimeg/emporia-time/pkg/api"
	"github.com/zimeg/emporia-time/pkg/cognito"
	"github.com/zimeg/emporia-time/pkg/config"
)

// Root facilitates the setup and execution of the command
func Root(
	ctx context.Context,
	cog cognito.Cognitoir,
	fs afero.Fs,
	req api.Emporiac,
	args []string,
	version string,
) (
	etime.CommandResult,
	error,
) {
	cmd, flags, err := config.ParseFlags(args)
	if err != nil {
		return etime.CommandResult{}, err
	} else if flags.Version {
		fmt.Printf("%s\n", version)
		return etime.CommandResult{}, nil
	} else if flags.Help {
		templates.PrintHelpMessage(os.Stderr)
		return etime.CommandResult{}, nil
	}
	cfg, err := config.Load(ctx, cog, fs, req, flags)
	if err != nil {
		return etime.CommandResult{}, err
	}
	results, err := etime.Run(cmd, &cfg)
	if err != nil {
		return etime.CommandResult{}, err
	}
	stats, err := templates.FormatUsage(results, flags.Portable)
	if err != nil {
		return results, err
	} else {
		fmt.Fprintf(os.Stderr, "%s\n", stats)
	}
	return results, nil
}
