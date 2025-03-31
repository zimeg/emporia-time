package times

import (
	"bytes"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/zimeg/emporia-time/internal/errors"
	"github.com/zimeg/emporia-time/internal/logs"
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
	Code int
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
func TimeExec(args []string, logger logs.Logger) (TimeMeasurement, error) {
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
	results, warnings := parseTimeResults(stderr.buff.String())
	for _, warning := range warnings {
		logger.Warn(warning)
	}
	times.Elapsed = times.End.Sub(times.Start)
	times.Command = results
	if err != nil {
		var exits *exec.ExitError
		ok := errors.As(err, &exits)
		if !ok {
			return times, errors.Wrap(errors.ErrTimeCommand, err)
		}
		times.Command.Code = exits.ExitCode()
	}
	return times, nil
}

// parseTimeResults extracts the time information from output
func parseTimeResults(output string) (times CommandTime, errs []error) {
	lines := strings.TrimSpace(output)
	for line := range strings.SplitSeq(lines, "\n") {
		fields := strings.Fields(line)
		if len(fields) < 2 {
			continue
		}
		measurement, value := fields[0], fields[1]
		switch measurement {
		case "code":
			parsed, err := strconv.ParseInt(value, 10, 32)
			if err != nil {
				errs = append(errs, errors.Wrap(errors.ErrTimeParseCode, err))
			}
			times.Code = int(parsed)
		case "real":
			parsed, err := strconv.ParseFloat(value, 64)
			if err != nil {
				errs = append(errs, errors.Wrap(errors.ErrTimeParseReal, err))
			}
			times.Real = parsed
		case "user":
			parsed, err := strconv.ParseFloat(value, 64)
			if err != nil {
				errs = append(errs, errors.Wrap(errors.ErrTimeParseUser, err))
			}
			times.User = parsed
		case "sys":
			parsed, err := strconv.ParseFloat(value, 64)
			if err != nil {
				errs = append(errs, errors.Wrap(errors.ErrTimeParseSys, err))
			}
			times.Sys = parsed
		}
	}
	return times, errs
}
