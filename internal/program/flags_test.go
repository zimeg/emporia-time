package program

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseFlags(t *testing.T) {
	tests := map[string]struct {
		Arguments []string
		Command   Command
		Error     error
	}{
		"plain arguments are treated as a command": {
			Arguments: []string{"etime", "sleep", "12"},
			Command: Command{
				Args: []string{"sleep", "12"},
			},
		},
		"flags before the command are parsed as flags": {
			Arguments: []string{"etime", "-p", "make", "build"},
			Command: Command{
				Args:  []string{"make", "build"},
				Flags: Flags{Portable: true},
			},
		},
		"flags after the command are parsed as command": {
			Arguments: []string{"etime", "zip", "code.zip", "-r", "."},
			Command: Command{
				Args: []string{"zip", "code.zip", "-r", "."},
			},
		},
		"overlapping command flags are for the command": {
			Arguments: []string{"etime", "unzip", "-p", "code.zip"},
			Command: Command{
				Args: []string{"unzip", "-p", "code.zip"},
			},
		},
		"duplicated flags that matched are set separate": {
			Arguments: []string{"etime", "-p", "mkdir", "-p", "/tmp/words"},
			Command: Command{
				Args:  []string{"mkdir", "-p", "/tmp/words"},
				Flags: Flags{Portable: true},
			},
		},
		"multiple energy flags are accepted at a time": {
			Arguments: []string{"etime", "--username", "example", "--password", "123", "ls"},
			Command: Command{
				Args:  []string{"ls"},
				Flags: Flags{Username: "example", Password: "123"},
			},
		},
		"help is noticed when help flags are provided": {
			Arguments: []string{"etime", "-h"},
			Command: Command{
				Args:  []string{},
				Flags: Flags{Help: true},
			},
		},
		"help is noticed when no arguments are provided": {
			Arguments: []string{"etime"},
			Command: Command{
				Args:  []string{},
				Flags: Flags{Help: true},
			},
		},
		"help is noticed when no commands are provided": {
			Arguments: []string{"etime", "-p"},
			Command: Command{
				Args:  []string{},
				Flags: Flags{Help: true, Portable: true},
			},
		},
		"parsing errors are returned before the command": {
			Arguments: []string{"etime", "--help=2"},
			Command: Command{
				Args:  []string{},
				Flags: Flags{},
			},
			Error: fmt.Errorf("invalid boolean value"),
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			command, err := ParseFlags(tt.Arguments)
			if tt.Error != nil {
				assert.ErrorContains(t, err, tt.Error.Error())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.Command, command)
			}
		})
	}
}
