package times

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zimeg/emporia-time/internal/errors"
)

func TestParseTimeResults(t *testing.T) {
	tests := map[string]struct {
		Output []string
		Times  CommandTime
		Errors []error
	}{
		"parse the portable output of the time command": {
			[]string{
				"code 0",
				"real 6.00",
				"user 121.20",
				"sys 0.80",
			},
			CommandTime{
				Code: 0,
				Real: 6.0,
				User: 121.2,
				Sys:  0.8,
			},
			nil,
		},
		"continues parsing even if one output errors": {
			[]string{
				"code 12",
				"real 1.23",
				"user ****",
				"sys 90000",
			},
			CommandTime{
				Code: 12,
				Real: 1.23,
				Sys:  90000,
			},
			[]error{
				errors.Err{
					Code: errors.ErrTimeParseUser,
				},
			},
		},
		"continues parsing even if all outputs error": {
			[]string{
				"code x",
				"real ~6.00",
				"user ~121.20",
				"sys ~0.80",
			},
			CommandTime{},
			[]error{
				errors.Err{
					Code: errors.ErrTimeParseCode,
				},
				errors.Err{
					Code: errors.ErrTimeParseReal,
				},
				errors.Err{
					Code: errors.ErrTimeParseUser,
				},
				errors.Err{
					Code: errors.ErrTimeParseSys,
				},
			},
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			output := strings.Join(tt.Output, "\n")
			times, errs := parseTimeResults(output)
			assert.Equal(t, tt.Times, times)
			require.Equal(t, len(tt.Errors), len(errs))
			for ii, err := range tt.Errors {
				assert.True(t, errors.Is(err, errs[ii]))
			}
		})
	}
}
