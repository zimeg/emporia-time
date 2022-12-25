package main

import (
	"bytes"
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"time"
)

// TimeExec performs the `args` command with timing, without interactivity
func TimeExec(args ...string) (time.Time, time.Time) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	cmd := exec.Command("time", args...)
	if errors.Is(cmd.Err, exec.ErrDot) {
		cmd.Err = nil
	}
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	startTime := time.Now().UTC()
	if err := cmd.Run(); err != nil {
		log.Fatalf("Error: Failed to execute command (%v)\n", err)
	}
	endTime := time.Now().UTC()

	fmt.Printf("%s", stdout.String())
	fmt.Fprintf(os.Stderr, "%s", stderr.String())

	return startTime, endTime
}

// main executes the command and displays energy stats
func main() {

	// share usage info on empty input
	if len(os.Args) == 1 {
		outputHelp()
		os.Exit(0)
	}

	// instantiate an Emporia client
	e := new(Emporia)
	e.Init()

	available, _ := EmporiaStatus()
	if !available {
		log.Panicf("Error: Cannot measure energy during Emporia maintenance\n")
	}

	// perform and measure the command
	prog := os.Args[1:]
	start, end := TimeExec(prog...)

	// gather and display usage information
	watts, sureness := e.CollectEnergyUsage(start, end)
	outputUsage(watts, sureness)
}
