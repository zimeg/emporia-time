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
	fmt.Printf("%s", stderr.String())
	return nil
}

// main executes the command and displays energy stats
func main() {
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
	elapsedTime := endTime.Sub(startTime)

	// query emporia for usage stats
	time.Sleep(2 * time.Second) // delay to respect latency
	chart, err := e.LookupEnergyUsage(startTime, endTime)
	if err != nil {
		log.Panicf("Error: Failed to gather energy usage data (%v)\n", err)
	}

	// display the estimated usage stats
	usage, sureness := ExtrapolateUsage(chart, elapsedTime.Seconds())
	fmt.Printf("%12.2f watt %11.1f%% sure\n", usage, sureness*100)
}
