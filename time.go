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
	e := &Emporia{
		token:  os.Getenv("EMPORIA_TOKEN"),
		device: os.Getenv("EMPORIA_DEVICE"),
	}
	if e.token == "" {
		log.Panicf("Error: EMPORIA_TOKEN environment variable not set\n")
	}
	if e.device == "" {
		log.Panicf("Error: EMPORIA_DEVICE environment variable not set\n")
	}

	// perform and observe the command
	startTime := time.Now().UTC()

	prog := os.Args[1:]
	err := TimeExec(prog...)
	if err != nil {
		log.Fatalf("Error: Failed to execute command (%v)\n", err)
	}
	endTime := time.Now().UTC()
	e.elapsedTime = endTime.Sub(startTime)

	// query emporia for usage stats
	time.Sleep(2 * time.Second) // delay to respect latency
	_, err = e.LookupEnergyUsage(startTime, endTime)
	if err != nil {
		log.Panicf("Error: Failed to gather energy usage data (%v)\n", err)
	}

	fmt.Printf("%12.2f watt %11.1f%% sure\n", e.usage, e.sureness*100)
}
