package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/spf13/afero"
	"github.com/zimeg/emporia-time/cmd/etime"
	"github.com/zimeg/emporia-time/internal/display/templates"
	"github.com/zimeg/emporia-time/internal/errors"
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
		return etime.CommandResult{}, errors.Wrap(errors.ErrConfigParse, err)
	} else if flags.Version {
		fmt.Printf("%s\n", version)
		return etime.CommandResult{}, nil
	} else if flags.Help {
		err := templates.PrintHelpMessage(os.Stderr)
		if err != nil {
			return etime.CommandResult{}, errors.Wrap(errors.ErrConfigHome, err)
		}
		return etime.CommandResult{}, nil
	}
	cfg, err := config.Setup(ctx, cog, fs, req, flags)
	if err != nil {
		return etime.CommandResult{}, errors.Wrap(errors.ErrConfigSetup, err)
	}
	results, err := etime.Run(cmd, &cfg)
	if err != nil {
		return etime.CommandResult{}, errors.Wrap(errors.ErrTimeCommand, err)
	}
	stats, err := templates.FormatUsage(results, flags.Portable)
	if err != nil {
		return results, errors.Wrap(errors.ErrTemplateFormat, err)
	} else {
		fmt.Fprintf(os.Stderr, "%s\n", stats)
	}
	return results, nil
}
