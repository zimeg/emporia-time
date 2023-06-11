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
			"error and return command outputs if time value is missing",
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
			t.Logf("%+v", tt)
			t.Fatalf("The correct command output was not retained!\nTEST: '%s'\nEXPECT: %+v\nACTUAL: %+v",
				tt.Title,
				stderr.Bytes(),
				cmd.Bytes(),
			)
		}
	}
}
