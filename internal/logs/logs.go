package logs

import (
	"io"
	"time"

	"github.com/charmbracelet/log"
	"github.com/zimeg/emporia-time/internal/errors"
)

// Logger outputs information
type Logger interface {
	Warn(err error)
	Error(err error)
	Fatal(err error)
}

// Logger outputs information
type Logs struct {
	charms *log.Logger
}

// NewLogger creates a new logger
func NewLogger(out io.Writer) *Logs {
	logger := log.NewWithOptions(
		out,
		log.Options{
			TimeFormat: time.RFC3339Nano,
		},
	)
	return &Logs{
		charms: logger,
	}
}

// Warn logs the warning
func (l *Logs) Warn(err error) {
	cast := errors.Unwrap(err)
	if cast.Source != nil {
		if errors.Unwrap(cast.Source).Code != errors.ErrUnexpectedProblem {
			l.Warn(cast.Source)
		} else {
			l.charms.Warn(cast.Message, "code", cast.Code, "source", cast.Source)
			return
		}
	}
	l.charms.Warn(cast.Message, "code", cast.Code)
}

// Error logs the error
func (l *Logs) Error(err error) {
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
func (l *Logs) Fatal(err error) {
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
