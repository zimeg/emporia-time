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

// TimeMeasurement holds information of a command run
type TimeMeasurement struct {
	Start   time.Time
	End     time.Time
	Elapsed time.Duration
	Command CommandTime
}

// CommandTime contains the values from the time command
type CommandTime struct {
	Real string
	User string
	Sys  string
}

// TimeExec performs the command and prints outputs while measuring timing
func TimeExec(command ...string) (TimeMeasurement, error) {
	var times TimeMeasurement
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	timeFlags := []string{"-p"}
	timeArgs := append(timeFlags, command...)

	cmd := exec.Command("/usr/bin/time", timeArgs...)
	if errors.Is(cmd.Err, exec.ErrDot) {
		cmd.Err = nil
	}
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	times.Start = time.Now().UTC()
	err := cmd.Run()
	times.End = time.Now().UTC()
	times.Elapsed = times.End.Sub(times.Start)

	results, stderr, warning := parseTimeResults(stderr)
	if warning != nil {
		log.Printf("Warning: %s", warning)
	}
	times.Command = results

	fmt.Printf("%s", stdout.String())
	fmt.Fprintf(os.Stderr, "%s", stderr.String())

	return times, err
}

// parseTimeResults extracts the time information from output
func parseTimeResults(output bytes.Buffer) (CommandTime, bytes.Buffer, error) {
	times := CommandTime{}
	lines := strings.Split(output.String(), "\n")

	var cmd []string
	userTimeFound, sysTimeFound, realTimeFound := false, false, false
	for _, line := range lines {
		fields := strings.Fields(line)
		matched := false

		if len(fields) == 2 {
			name := fields[0]
			value := fields[1]
			switch name {
			case "user":
				times.User = trimTimeValue(value)
				matched = true
				userTimeFound = true
			case "sys":
				times.Sys = trimTimeValue(value)
				matched = true
				sysTimeFound = true
			case "real":
				times.Real = trimTimeValue(value)
				matched = true
				realTimeFound = true
			}
		}
		if !matched {
			cmd = append(cmd, line)
		}
	}

	buff := bytes.NewBufferString(strings.Join(cmd, "\n"))
	if !userTimeFound || !sysTimeFound || !realTimeFound {
		return times, *buff, errors.New("A time value is missing in the output!")
	}
	return times, *buff, nil
}

// trimTimeValue removes most leading zeros
func trimTimeValue(value string) string {
	re := regexp.MustCompile(`^0*:0?`)
	trim := re.ReplaceAllString(value, "")
	return trim
}
