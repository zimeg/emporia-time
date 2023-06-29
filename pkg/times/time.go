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

	"github.com/zimeg/emporia-time/internal/terminal"
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

// TimeExec performs the command and prints outputs while measuring timing
func TimeExec(command terminal.Command) (TimeMeasurement, error) {
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
	fmt.Fprintf(os.Stderr, "%s", stderr.String())

	return times, err
}

// parseTimeResults extracts the time information from output
func parseTimeResults(output bytes.Buffer) (CommandTime, bytes.Buffer, error) {
	lines := strings.Split(output.String(), "\n")

	var cmd []string
	var userTimeFound, sysTimeFound, realTimeFound bool
	var userTimeIndex, sysTimeIndex, realTimeIndex int
	var userTimeValue, sysTimeValue, realTimeValue string
	var userTimeFloat, sysTimeFloat, realTimeFloat float64
	for ii, line := range lines {
		fields := strings.Fields(line)
		matched := false

		if len(fields) == 2 {
			name := fields[0]
			value := fields[1]
			switch name {
			case "user":
				matched = true
				userTimeFound = true
				userTimeIndex = ii
				userTimeValue = value
			case "sys":
				matched = true
				sysTimeFound = true
				sysTimeIndex = ii
				sysTimeValue = value
			case "real":
				matched = true
				realTimeFound = true
				realTimeIndex = ii
				realTimeValue = value
			}
		}
		if !matched {
			cmd = append(cmd, line)
		}
	}

	var buff bytes.Buffer
	for ii, line := range lines {
		switch {
		case ii == userTimeIndex && userTimeFound:
		case ii == sysTimeIndex && sysTimeFound:
		case ii == realTimeIndex && realTimeFound:
		default:
			str := fmt.Sprintf("%s\n", line)
			buff.WriteString(str)
		}
	}
	if buff.Len() > 0 {
		buff.Truncate(buff.Len() - 1)
	}

	var err error
	if !userTimeFound || !sysTimeFound || !realTimeFound {
		return CommandTime{}, buff, errors.New("A time value is missing in the output!")
	}
	if userTimeFloat, err = parseTimeValue(userTimeValue); err != nil {
		return CommandTime{}, buff, errors.New("Failed to parse the user time value!")
	}
	if sysTimeFloat, err = parseTimeValue(sysTimeValue); err != nil {
		return CommandTime{}, buff, errors.New("Failed to parse the sys time value!")
	}
	if realTimeFloat, err = parseTimeValue(realTimeValue); err != nil {
		return CommandTime{}, buff, errors.New("Failed to parse the real time value!")
	}

	times := CommandTime{
		User: userTimeFloat,
		Sys:  sysTimeFloat,
		Real: realTimeFloat,
	}
	return times, buff, nil
}

// parseTimeValue converts a string to a float64
func parseTimeValue(value string) (float64, error) {
	return strconv.ParseFloat(value, 64)
}
