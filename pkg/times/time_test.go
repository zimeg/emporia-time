package times

import (
	"bytes"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseTimeResults(t *testing.T) {
	tests := map[string]struct {
		Output    []string
		Times     CommandTime
		CmdOutput []string
		Error     error
	}{
		"check a command like sleep without any outputs": {
			[]string{
				"real 1.00",
				"user 1.00",
				"sys 1.00",
			},
			CommandTime{
				Real: 1.0,
				User: 1.0,
				Sys:  1.0,
			},
			[]string{},
			nil,
		},
		"capture and remove time values from the output": {
			[]string{
				"example command output!",
				"real 3754.56",
				"user 1301.21",
				"sys 3.24",
			},
			CommandTime{
				Real: 3754.56,
				User: 1301.21,
				Sys:  3.24,
			},
			[]string{
				"example command output!",
			},
			nil,
		},
		"prefer the latest outputs for timing information": {
			[]string{
				"this command outputs something familiar",
				"real 3.00",
				"user 3.00",
				"sys 3.00",
				"but suppose that was hardcoded!",
				"the real output now follows",
				"real 240.00",
				"user 8.00",
				"sys 12.00",
			},
			CommandTime{
				Real: 240.0,
				User: 8.0,
				Sys:  12.0,
			},
			[]string{
				"this command outputs something familiar",
				"real 3.00",
				"user 3.00",
				"sys 3.00",
				"but suppose that was hardcoded!",
				"the real output now follows",
			},
			nil,
		},
		"ensure the parser doesnt break on unexpected values": {
			[]string{
				"gathering example account information",
				"user goatish_burr",
				"real true",
				"sys mountains",
				"real 240.00",
				"user 8.05",
				"sys 12.00",
			},
			CommandTime{
				Real: 240.0,
				User: 8.05,
				Sys:  12.0,
			},
			[]string{
				"gathering example account information",
				"user goatish_burr",
				"real true",
				"sys mountains",
			},
			nil,
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			buff := bytes.NewBufferString(strings.Join(tt.Output, "\n"))
			stdout := bytes.NewBufferString(strings.Join(tt.CmdOutput, "\n"))
			times, cmd, err := parseTimeResults(*buff)
			if tt.Error != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.Error, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.Times, times)
				if len(stdout.Bytes()) != 0 {
					assert.Equal(t, stdout.Bytes(), cmd.Bytes())
				} else {
					assert.Nil(t, cmd.Bytes())
				}
			}
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
