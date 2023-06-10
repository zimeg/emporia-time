package main

import (
	"bytes"
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"time"
)

// TimeMeasurement
type TimeMeasurement struct {
	Start   time.Time
	End     time.Time
	Elapsed time.Duration
	Command CommandTime
}

type CommandTime struct {
	Real string
	User string
	Sys  string
}

// TimeExec performs the `args` command with timing, without interactivity
func TimeExec(args ...string) TimeMeasurement {
	var times TimeMeasurement
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	timeFlags := []string{"--format", "real %E\nuser %U\nsys %S"}
	timeArgs := append(timeFlags, args...)

	cmd := exec.Command("/usr/bin/time", timeArgs...)
	if errors.Is(cmd.Err, exec.ErrDot) {
		cmd.Err = nil
	}
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	times.Start = time.Now().UTC()
	if err := cmd.Run(); err != nil {
		log.Fatalf("Error: Failed to execute command (%v)\n", err)
	}
	times.End = time.Now().UTC()

	times.Elapsed = times.End.Sub(times.Start)
	times.Command, stderr = parseTimeResults(stderr)

	fmt.Printf("%s", stdout.String())
	fmt.Fprintf(os.Stderr, "%s", stderr.String())

	return times
}

// parseTimeResults extracts the time information from output
func parseTimeResults(output bytes.Buffer) (CommandTime, bytes.Buffer) {
	times := CommandTime{}
	lines := strings.Split(output.String(), "\n")
	var cmd []string
	for _, line := range lines {
		fields := strings.Fields(line)
		if len(fields) == 2 {
			name := fields[0]
			value := fields[1]

			switch name {
			case "user":
				times.User = trimTimeValue(value)
			case "sys":
				times.Sys = trimTimeValue(value)
			case "real":
				times.Real = trimTimeValue(value)
			default:
				cmd = append(cmd, line)
			}
		}
	}

	buff := bytes.NewBufferString(strings.Join(cmd, "\n"))
	return times, *buff
}

// trimTimeValue removes most leading zeros
func trimTimeValue(value string) string {
	re := regexp.MustCompile(`^0*:0?`)
	trim := re.ReplaceAllString(value, "")
	return trim
}
