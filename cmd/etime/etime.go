package etime

import (
	"github.com/zimeg/emporia-time/internal/errors"
	"github.com/zimeg/emporia-time/internal/logs"
	"github.com/zimeg/emporia-time/pkg/config"
	"github.com/zimeg/emporia-time/pkg/energy"
	"github.com/zimeg/emporia-time/pkg/times"
)

// CommandResult holds information from the run command
type CommandResult struct {
	energy.EnergyResult
	times.TimeMeasurement
}

// Run executes the command and returns the usage statistics
func Run(
	cmd []string,
	cfg config.Configure,
	logger logs.Logger,
) (
	results CommandResult,
	err error,
) {
	available, err := cfg.API().Status()
	if err != nil {
		return CommandResult{}, errors.Wrap(errors.ErrEmporiaCheckup, err)
	} else if !available {
		return CommandResult{}, errors.New(errors.ErrEmporiaMaintenance)
	}
	measurements, err := times.TimeExec(cmd, logger)
	if err != nil {
		return CommandResult{}, errors.Wrap(errors.ErrTimeExecution, err)
	} else {
		results.TimeMeasurement = measurements
	}
	usage, err := cfg.API().GetChartUsage(results.TimeMeasurement)
	if err != nil {
		return results, errors.Wrap(errors.ErrEmporiaChart, err)
	} else {
		results.EnergyResult = usage
	}
	return results, nil
}
