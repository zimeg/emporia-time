package main

import (
	"bytes"
	"errors"
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
				Real: "1.00",
				User: "1.00",
				Sys:  "1.00",
			},
			[]string{},
			nil,
		},
		{
			"capture and remove time values from the output",
			[]string{
				"example command output!",
				"real 1:02:34.56",
				"user 21:41.21",
				"sys 3.24",
			},
			CommandTime{
				Real: "1:02:34.56",
				User: "21:41.21",
				Sys:  "3.24",
			},
			[]string{
				"example command output!",
			},
			nil,
		},
		{
			"error and return command outputs if a time value is missing",
			[]string{
				"something strange happened here...",
				"sys 10:00.00",
			},
			CommandTime{
				Sys: "10:00.00",
			},
			[]string{
				"something strange happened here...",
			},
			errors.New("A time value is missing in the output!"),
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
				"real 4:00.00",
				"user 8.00",
				"sys 12.00",
			},
			CommandTime{
				Real: "4:00.00",
				User: "8.00",
				Sys:  "12.00",
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
	}

	for _, tt := range tests {
		buff := bytes.NewBufferString(strings.Join(tt.Output, "\n"))
		times, cmd, err := parseTimeResults(*buff)
		stderr := bytes.NewBufferString(strings.Join(tt.CmdOutput, "\n"))

		if tt.Error != err {
			if tt.Error != nil && err != nil && tt.Error.Error() == err.Error() {
				continue
			}
			t.Fatalf("An unexpected error was encountered!\nTEST: '%s'\nEXPECT: %+v\nACTUAL: %+v",
				tt.Title,
				tt.Error,
				err,
			)
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

func TestTrimTimeValue(t *testing.T) {
	tests := []struct {
		Title    string
		Value    []string
		Expected []string
	}{
		{
			"retain the necessary leading zero",
			[]string{"0.00", "0.08"},
			[]string{"0.00", "0.08"},
		},
		{
			"ignore padding on seconds as needed",
			[]string{"2.00", "04.00", "0:06.02", "0:12.80", "0:00:06.10"},
			[]string{"2.00", "4.00", "6.02", "12.80", "6.10"},
		},
		{
			"ignore padding on minutes when possible",
			[]string{"9:33.66", "04:21.01", "00:14:10.70"},
			[]string{"9:33.66", "4:21.01", "14:10.70"},
		},
	}

	for _, tt := range tests {
		for ii, val := range tt.Value {
			trimmed := trimTimeValue(val)
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
