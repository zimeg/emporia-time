package program

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseFlags(t *testing.T) {
	tests := map[string]struct {
		arguments []string
		command   Command
		makesExit bool
	}{
		"plain arguments are treated as a command": {
			arguments: []string{"etime", "sleep", "12"},
			command: Command{
				Args: []string{"sleep", "12"},
			},
		},
		"flags before the command are parsed as flags": {
			arguments: []string{"etime", "-p", "make", "build"},
			command: Command{
				Args:  []string{"make", "build"},
				Flags: Flags{Portable: true},
			},
		},
		"flags after the command are parsed as command": {
			arguments: []string{"etime", "zip", "code.zip", "-r", "."},
			command: Command{
				Args: []string{"zip", "code.zip", "-r", "."},
			},
		},
		"overlapping command flags are for the command": {
			arguments: []string{"etime", "unzip", "-p", "code.zip"},
			command: Command{
				Args: []string{"unzip", "-p", "code.zip"},
			},
		},
		"duplicated flags that matched are set separate": {
			arguments: []string{"etime", "-p", "mkdir", "-p", "/tmp/words"},
			command: Command{
				Args:  []string{"mkdir", "-p", "/tmp/words"},
				Flags: Flags{Portable: true},
			},
		},
		"multiple energy flags are accepted at a time": {
			arguments: []string{"etime", "--username", "example", "--password", "123", "ls"},
			command: Command{
				Args:  []string{"ls"},
				Flags: Flags{Username: "example", Password: "123"},
			},
		},
		"help is noticed when help flags are provided": {
			arguments: []string{"etime", "-h"},
			command: Command{
				Args:  []string{},
				Flags: Flags{Help: true},
			},
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			command := ParseFlags(tt.arguments)
			assert.Equal(t, tt.command, command)
		})
	}
}
