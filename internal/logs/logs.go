package logs

import (
	"io"
	"time"

	"github.com/charmbracelet/log"
	"github.com/zimeg/emporia-time/internal/errors"
)

// Logger outputs information
type Logger struct {
	charms *log.Logger
}

// NewLogger creates a new logger
func NewLogger(out io.Writer) *Logger {
	logger := log.NewWithOptions(
		out,
		log.Options{
			TimeFormat: time.RFC3339Nano,
		},
	)
	return &Logger{
		charms: logger,
	}
}

// Error logs the error
func (l *Logger) Error(err error) {
	cast := errors.Unwrap(err)
	if cast.Source != nil {
		if errors.Unwrap(cast.Source).Code != errors.ErrUnexpectedProblem {
			l.Error(cast.Source)
		} else {
			l.charms.Error(cast.Message, "code", cast.Code, "source", cast.Source)
			return
		}
	}
	l.charms.Error(cast.Message, "code", cast.Code)
}

// Fatal logs the error and exits
func (l *Logger) Fatal(err error) {
	cast := errors.Unwrap(err)
	if cast.Source != nil {
		if errors.Unwrap(cast.Source).Code != errors.ErrUnexpectedProblem {
			l.Error(cast.Source)
		} else {
			l.charms.Fatal(cast.Message, "code", cast.Code, "source", cast.Source)
			return
		}
	}
	l.charms.Fatal(cast.Message, "code", cast.Code)
}
