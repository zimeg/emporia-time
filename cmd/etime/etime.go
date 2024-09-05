package etime

import (
	"fmt"
	"os/exec"

	"github.com/zimeg/emporia-time/pkg/config"
	"github.com/zimeg/emporia-time/pkg/energy"
	"github.com/zimeg/emporia-time/pkg/times"
)

// CommandResult holds information from the run command
type CommandResult struct {
	energy.EnergyResult
	times.TimeMeasurement
	ExitCode int
}

// Run executes the command and returns the usage statistics
func Run(cmd []string, cfg config.Configure) (results CommandResult, err error) {
	available, err := cfg.API().Status()
	if err != nil {
		return CommandResult{}, err
	} else if !available {
		return CommandResult{}, fmt.Errorf("Error: Cannot measure energy during Emporia maintenance")
	}
	measurements, err := times.TimeExec(cmd)
	if err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			results.ExitCode = exitError.ExitCode()
		} else {
			return CommandResult{}, err
		}
		results.TimeMeasurement = measurements
	} else {
		results.TimeMeasurement = measurements
	}
	usage, err := cfg.API().GetChartUsage(results.TimeMeasurement)
	if err != nil {
		return results, err
	} else {
		results.EnergyResult = usage
	}
	return results, nil
}
