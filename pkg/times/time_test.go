package times

import (
	"bytes"
	"strings"
	"testing"
)

func TestParseTimeResults(t *testing.T) {
	tests := []struct {
		Title     string
		Output    []string
		Times     CommandTime
		CmdOutput []string
		Error     error
	}{
		{
			"check a command like sleep without any outputs",
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
		{
			"capture and remove time values from the output",
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
		{
			"prefer the latest outputs for timing information",
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
		{
			"ensure the parser doesnt break on unexpected values",
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

	for _, tt := range tests {
		buff := bytes.NewBufferString(strings.Join(tt.Output, "\n"))
		times, cmd, err := parseTimeResults(*buff)
		stderr := bytes.NewBufferString(strings.Join(tt.CmdOutput, "\n"))

		if tt.Error != err {
			if tt.Error != nil && err != nil && tt.Error.Error() != err.Error() {
				t.Fatalf("An unexpected error was encountered!\nTEST: '%s'\nEXPECT: %+v\nACTUAL: %+v",
					tt.Title,
					tt.Error,
					err,
				)
			}
		}
		if tt.Times != times {
			t.Fatalf("Not all times were collected!\nTEST: '%s'\nEXPECT: %+v\nACTUAL: %+v",
				tt.Title,
				tt.Times,
				times,
			)
		}
		if !bytes.Equal(stderr.Bytes(), cmd.Bytes()) {
			t.Fatalf("The correct command output was not retained!\nTEST: '%s'\nEXPECT: %+v\nACTUAL: %+v",
				tt.Title,
				stderr.String(),
				cmd.String(),
			)
		}
	}
}

func TestParseTimeValue(t *testing.T) {
	tests := []struct {
		Title    string
		Value    []string
		Expected []float64
	}{
		{
			"retain the necessary leading zero",
			[]string{"0.00", "0.08"},
			[]float64{0.0, 0.08},
		},
		{
			"ignore padding on seconds as needed",
			[]string{"2.00", "4.00", "6.02", "12.80", "6.10"},
			[]float64{2.0, 4.0, 6.02, 12.8, 6.1},
		},
		{
			"parse times greater than sixty seconds",
			[]string{"573.66", "261.01", "850.70", "86405.12"},
			[]float64{573.66, 261.01, 850.7, 86405.12},
		},
	}

	for _, tt := range tests {
		for ii, val := range tt.Value {
			trimmed, err := parseTimeValue(val)
			if err != nil {
				t.Fatalf("An unexpected error was found!\nTEST: '%s'\nEXPECT: %+v\nACTUAL: %+v",
					tt.Title,
					nil,
					err,
				)
			}
			if trimmed != tt.Expected[ii] {
				t.Fatalf("The formatting of time seems off for '%s'!\nTEST: '%s'\nEXPECT: %+v\nACTUAL: %+v",
					val,
					tt.Title,
					tt.Expected[ii],
					trimmed,
				)
			}
		}
	}
}
