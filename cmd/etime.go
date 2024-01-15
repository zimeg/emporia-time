package etime

import (
	"os/exec"

	"github.com/zimeg/emporia-time/internal/program"
	"github.com/zimeg/emporia-time/pkg/emporia"
	"github.com/zimeg/emporia-time/pkg/energy"
	"github.com/zimeg/emporia-time/pkg/times"
)

// CommandResult holds information from the run command
type CommandResult struct {
	energy.EnergyResult
	times.TimeMeasurement
	ExitCode int
}

// Setup prepares the command and client with provided inputs
func Setup(arguments []string) (command program.Command, client emporia.Emporia, err error) {
	command, err = program.ParseFlags(arguments)
	if err != nil || command.Flags.Help {
		return command, client, err
	}
	if config, err := emporia.SetupConfig(command.Flags); err != nil {
		return command, client, err
	} else {
		client.Config = config
	}
	return command, client, err
}

// Run executes the command and returns the usage statistics
func Run(command program.Command, client emporia.Emporia) (results CommandResult, err error) {
	if measurements, err := times.TimeExec(command); err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			results.ExitCode = exitError.ExitCode()
		} else {
			return results, err
		}
		results.TimeMeasurement = measurements
	} else {
		results.TimeMeasurement = measurements
	}
	if usage, err := client.CollectEnergyUsage(results.TimeMeasurement); err != nil {
		return results, err
	} else {
		results.EnergyResult = usage
	}
	return results, nil
}
