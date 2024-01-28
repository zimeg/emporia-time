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
		Error  error
	}{
		"parse the portable output of the time command": {
			[]string{
				"real 6.00",
				"user 2.20",
				"sys 3.80",
			},
			CommandTime{
				Real: 6.0,
				User: 2.2,
				Sys:  3.8,
			},
			nil,
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			output := strings.Join(tt.Output, "\n")
			times, err := parseTimeResults(output)
			assert.Equal(t, tt.Error, err)
			assert.Equal(t, tt.Times, times)
		})
	}
}

func TestParseTimeValue(t *testing.T) {
	tests := map[string]struct {
		Value    []string
		Expected []float64
	}{
		"retain the necessary leading zero": {
			[]string{"0.00", "0.08"},
			[]float64{0.0, 0.08},
		},
		"ignore padding on seconds as needed": {
			[]string{"2.00", "4.00", "6.02", "12.80", "6.10"},
			[]float64{2.0, 4.0, 6.02, 12.8, 6.1},
		},
		"parse times greater than sixty seconds": {
			[]string{"573.66", "261.01", "850.70", "86405.12"},
			[]float64{573.66, 261.01, 850.7, 86405.12},
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			for ii, val := range tt.Value {
				trimmed, err := parseTimeValue(val)
				assert.NoError(t, err)
				assert.Equal(t, tt.Expected[ii], trimmed)
			}
		})
	}
}
