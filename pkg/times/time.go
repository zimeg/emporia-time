package times

import (
	"bytes"
	"errors"
	"log"
	"os"
	"strconv"
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
func TimeExec(args []string) (TimeMeasurement, error) {
	times := TimeMeasurement{}
	stderr := bufferWriter{
		buff:   &bytes.Buffer{},
		std:    os.Stderr,
		bounds: makeBounds(),
	}

	cmd := timerCommand(args, stderr)
	times.Start = time.Now().UTC()
	err := cmd.Run()
	times.End = time.Now().UTC()

	results, warning := parseTimeResults(stderr.buff.String())
	if warning != nil {
		log.Printf("Warning: %s", warning)
	}
	times.Elapsed = times.End.Sub(times.Start)
	times.Command = results

	return times, err
}

// parseTimeResults extracts the time information from output
func parseTimeResults(output string) (times CommandTime, err error) {
	lines := strings.TrimSpace(output)
	for _, line := range strings.Split(lines, "\n") {
		fields := strings.Fields(line)
		if len(fields) < 2 {
			continue
		}
		measurement, value := fields[0], fields[1]
		switch measurement {
		case "real":
			if times.Real, err = parseTimeValue(value); err != nil {
				return times, errors.New("Failed to parse the real time value!")
			}
		case "user":
			if times.User, err = parseTimeValue(value); err != nil {
				return times, errors.New("Failed to parse the user time value!")
			}
		case "sys":
			if times.Sys, err = parseTimeValue(value); err != nil {
				return times, errors.New("Failed to parse the sys time value!")
			}
		}
	}
	return times, err
}

// parseTimeValue converts a string to a float64
func parseTimeValue(value string) (float64, error) {
	return strconv.ParseFloat(value, 64)
}
