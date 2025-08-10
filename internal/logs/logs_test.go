package logs

import (
	"bytes"
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zimeg/emporia-time/internal/errors"
)

func TestLogsWarn(t *testing.T) {
	tests := map[string]struct {
		err      error
		expected []string
	}{
		"outputs a warn log": {
			err: errors.Err{
				Code:    "err_unexpected_problem",
				Message: "something is broken",
			},
			expected: []string{
				"WARN something is broken code=err_unexpected_problem\n",
			},
		},
		"outputs warn logs": {
			err: errors.Err{
				Code:    "err_emporia_checkup",
				Message: "a reason or cause",
				Source: errors.Err{
					Code:    "err_emporia_status",
					Message: "a breaking thing",
					Source: errors.Err{
						Code:    "err_emporia_maintenance",
						Message: "the error in question",
						Source:  fmt.Errorf("some downtimes"),
					},
				},
			},
			expected: []string{
				"WARN the error in question code=err_emporia_maintenance source=\"some downtimes\"\n",
				"WARN a breaking thing code=err_emporia_status\n",
				"WARN a reason or cause code=err_emporia_checkup\n",
			},
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			buff := &bytes.Buffer{}
			logger := NewLogger(buff)
			logger.Warn(tt.err)
			assert.Equal(t, strings.Join(tt.expected, ""), buff.String())
		})
	}
}

func TestLogsError(t *testing.T) {
	tests := map[string]struct {
		err      error
		expected []string
	}{
		"outputs an error log": {
			err: errors.Err{
				Code:    "err_unexpected_problem",
				Message: "something is broken",
			},
			expected: []string{
				"ERRO something is broken code=err_unexpected_problem\n",
			},
		},
		"outputs error logs": {
			err: errors.Err{
				Code:    "err_emporia_checkup",
				Message: "a reason or cause",
				Source: errors.Err{
					Code:    "err_emporia_status",
					Message: "a breaking thing",
					Source: errors.Err{
						Code:    "err_emporia_maintenance",
						Message: "the error in question",
						Source:  fmt.Errorf("some downtimes"),
					},
				},
			},
			expected: []string{
				"ERRO the error in question code=err_emporia_maintenance source=\"some downtimes\"\n",
				"ERRO a breaking thing code=err_emporia_status\n",
				"ERRO a reason or cause code=err_emporia_checkup\n",
			},
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			buff := &bytes.Buffer{}
			logger := NewLogger(buff)
			logger.Error(tt.err)
			assert.Equal(t, strings.Join(tt.expected, ""), buff.String())
		})
	}
}
