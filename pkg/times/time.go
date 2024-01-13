package times

import (
	"bytes"
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/zimeg/emporia-time/internal/program"
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
	Real float64
	User float64
	Sys  float64
}

// GetReal returns the real time of a command
func (times TimeMeasurement) GetReal() float64 {
	return times.Command.Real
}

// GetUser returns the user time of a command
func (times TimeMeasurement) GetUser() float64 {
	return times.Command.User
}

// GetSys returns the sys time of a command
func (times TimeMeasurement) GetSys() float64 {
	return times.Command.Sys
}

// TimeExec performs the command and prints outputs while measuring timing
func TimeExec(command program.Command) (TimeMeasurement, error) {
	var times TimeMeasurement
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	timeFlags := []string{"-p"}
	timeArgs := append(timeFlags, command.Args...)

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
	if stderr.Len() > 0 {
		fmt.Fprintf(os.Stderr, "%s\n", stderr.String())
	}

	return times, err
}

// splitBuffer separates the command output from times
func splitBuffer(buff bytes.Buffer) (command, times []string) {
	output := buff.String()
	trimmed := strings.TrimRight(output, "\n")
	lines := strings.Split(trimmed, "\n")
	for i := len(lines) - 1; i >= 0; i-- {
		times = append([]string{lines[i]}, times...)
		if strings.Contains(lines[i], "real") {
			command = lines[:i]
			break
		}
	}
	return command, times
}

// parseTimeResults extracts the time information from output
func parseTimeResults(output bytes.Buffer) (times CommandTime, buff bytes.Buffer, err error) {
	commandLines, timeLines := splitBuffer(output)
	for _, line := range timeLines {
		fields := strings.Fields(line)
		measurement, value := fields[0], fields[1]
		switch measurement {
		case "real":
			if times.Real, err = parseTimeValue(value); err != nil {
				return times, buff, errors.New("Failed to parse the real time value!")
			}
		case "user":
			if times.User, err = parseTimeValue(value); err != nil {
				return times, buff, errors.New("Failed to parse the user time value!")
			}
		case "sys":
			if times.Sys, err = parseTimeValue(value); err != nil {
				return times, buff, errors.New("Failed to parse the sys time value!")
			}
		}
	}
	command := strings.Join(commandLines, "\n")
	buff.WriteString(command)
	return times, buff, err
}

// parseTimeValue converts a string to a float64
func parseTimeValue(value string) (float64, error) {
	return strconv.ParseFloat(value, 64)
}
