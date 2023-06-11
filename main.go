package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
)

// CommandResult holds information from the run command
type CommandResult struct {
	EnergyResult
	ExitCode int
	TimeMeasurement
}

// main executes the command and displays energy stats
func main() {
	if len(os.Args) == 1 {
		outputHelp()
		os.Exit(0)
	}

	// Instantiate an Emporia client
	e := new(Emporia)
	e.Init()

	if available, err := EmporiaStatus(); err != nil {
		log.Panicf("Error: %s", err)
	} else if !available {
		log.Panicf("Error: Cannot measure energy during Emporia maintenance\n")
	}

	// Perform and measure the command
	results := CommandResult{}
	prog := os.Args[1:]
	if measurements, err := TimeExec(prog...); err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			results.ExitCode = exitError.ExitCode()
		} else {
			log.Printf("Error: %s", err)
		}
		results.TimeMeasurement = measurements
	} else {
		results.TimeMeasurement = measurements
	}

	if usage, err := e.CollectEnergyUsage(results.TimeMeasurement); err != nil {
		log.Fatalf("Error: %s", err)
	} else {
		results.EnergyResult = usage
	}

	// Output the resulting measurements
	if stats, err := formatUsage(results); err != nil {
		log.Fatalf("Error: %s", err)
	} else {
		fmt.Fprintf(os.Stderr, "%s\n", stats)
	}
	os.Exit(results.ExitCode)
}
