package times

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseTimeResults(t *testing.T) {
	tests := map[string]struct {
		Output []string
		Times  CommandTime
		Errors []error
	}{
		"parse the portable output of the time command": {
			[]string{
				"real 6.00",
				"user 121.20",
				"sys 0.80",
			},
			CommandTime{
				Real: 6.0,
				User: 121.2,
				Sys:  0.8,
			},
			nil,
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			output := strings.Join(tt.Output, "\n")
			times, errs := parseTimeResults(output)
			assert.Equal(t, tt.Errors, errs)
			assert.Equal(t, tt.Times, times)
		})
	}
}
