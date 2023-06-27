package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/zimeg/emporia-time/internal/terminal"
	"github.com/zimeg/emporia-time/pkg/emporia"
	"github.com/zimeg/emporia-time/pkg/energy"
	"github.com/zimeg/emporia-time/pkg/times"
)

// CommandResult holds information from the run command
type CommandResult struct {
	EnergyResult    energy.EnergyResult
	ExitCode        int
	TimeMeasurement times.TimeMeasurement
}

// main executes the command and displays energy stats
func main() {
	command := terminal.ParseFlags(os.Args)
	client := new(emporia.Emporia)

	if config, err := emporia.SetupConfig(command.Flags); err != nil {
		log.Fatalf("Error: %s", err)
	} else {
		client.Config = config
	}

	if available, err := emporia.EmporiaStatus(); err != nil {
		log.Fatalf("Error: %s", err)
	} else if !available {
		log.Fatalf("Error: Cannot measure energy during Emporia maintenance\n")
	}

	// Perform and measure the command
	results := CommandResult{}
	if measurements, err := times.TimeExec(command); err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			results.ExitCode = exitError.ExitCode()
		} else {
			log.Printf("Error: %s", err)
		}
		results.TimeMeasurement = measurements
	} else {
		results.TimeMeasurement = measurements
	}

	if usage, err := client.CollectEnergyUsage(results.TimeMeasurement); err != nil {
		log.Fatalf("Error: %s", err)
	} else {
		results.EnergyResult = usage
	}

	// Output the resulting measurements
	if stats, err := formatUsage(results, command.Flags.Portable); err != nil {
		log.Fatalf("Error: %s", err)
	} else {
		fmt.Fprintf(os.Stderr, "%s\n", stats)
	}
	os.Exit(results.ExitCode)
}

// formatUsage arranges information about resource usage of a command
func formatUsage(results CommandResult, isPortableFormat bool) (string, error) {
	var energyTemplate string
	switch isPortableFormat {
	case false:
		energyTemplate = strings.TrimSpace(`
{{12 | Time .TimeMeasurement.Command.Real}} real {{12 | Time .TimeMeasurement.Command.User}} user {{12 | Time .TimeMeasurement.Command.Sys}} sys
{{12 | Value .EnergyResult.Joules}} joules {{10 | Value .EnergyResult.Watts}} watts {{10 | Percent .EnergyResult.Sureness}}% sure`)
	case true:
		energyTemplate = strings.TrimSpace(`
{{0 | Time .TimeMeasurement.Command.Real}} real
{{0 | Time .TimeMeasurement.Command.User}} user
{{0 | Time .TimeMeasurement.Command.Sys}} sys
{{0 | Value .EnergyResult.Joules}} joules
{{0 | Value .EnergyResult.Watts}} watts
{{0 | Percent .EnergyResult.Sureness}}% sure`)
	}

	body, err := terminal.TemplateBuilder(energyTemplate, results)
	if err != nil {
		return "", err
	}
	return body, nil
}
