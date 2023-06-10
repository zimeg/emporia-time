package main

import (
	"fmt"
	"log"
	"os"
)

// CommandResult holds information from the run command
type CommandResult struct {
	TimeMeasurement
	EnergyResult
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
	results.TimeMeasurement = TimeExec(prog...)

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
}
