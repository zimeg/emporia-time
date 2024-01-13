package main

import (
	"fmt"
	"log"
	"os"

	etime "github.com/zimeg/emporia-time/cmd"
	"github.com/zimeg/emporia-time/internal/display/templates"
	"github.com/zimeg/emporia-time/pkg/emporia"
)

// main manages the lifecycle of this program
func main() {
	command, client, err := etime.Setup()
	if err != nil {
		log.Fatalf("Error: %s", err)
	} else if command.Flags.Help {
		os.Exit(0)
	}
	if available, err := emporia.EmporiaStatus(); err != nil {
		log.Fatalf("Error: %s", err)
	} else if !available {
		log.Fatalf("Error: Cannot measure energy during Emporia maintenance\n")
	}
	results, err := etime.Run(command, client)
	if err != nil {
		log.Fatalf("Error: %s", err)
	}
	if stats, err := templates.FormatUsage(results, command.Flags.Portable); err != nil {
		log.Fatalf("Error: %s", err)
	} else {
		fmt.Fprintf(os.Stderr, "%s\n", stats)
	}
	os.Exit(results.ExitCode)
}
