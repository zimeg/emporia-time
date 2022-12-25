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
func TimeExec(args ...string) error {
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	cmd := exec.Command("time", args...)
	if errors.Is(cmd.Err, exec.ErrDot) {
		cmd.Err = nil
	}
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return err
	}

	fmt.Printf("%s", stdout.String())
	fmt.Fprintf(os.Stderr, "%s", stderr.String())
	return nil
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

	// perform and observe the command
	startTime := time.Now().UTC()

	prog := os.Args[1:]
	err := TimeExec(prog...)
	if err != nil {
		log.Fatalf("Error: Failed to execute command (%v)\n", err)
	}

	endTime := time.Now().UTC()

	// gather and display usage information
	watts, sureness := e.CollectEnergyUsage(startTime, endTime)
	outputUsage(watts, sureness)
}
