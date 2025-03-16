package templates

import (
	"bytes"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPrintHelpMessage(t *testing.T) {
	tests := map[string]struct {
		expected []string
	}{
		"outputs the expected help message": {
			expected: []string{
				"Measure the time and energy used while executing a command",
				"",
				"\x1b[1mUSAGE\x1b[0m",
				fmt.Sprintf("  %s [flags] <command> [args]", os.Args[0]),
				"",
				"\x1b[1mDESCRIPTION\x1b[0m",
				"  flags    optional flags to provide this program",
				"  command  the program to execute and measure",
				"  args     optional arguments for the command",
				"",
				"\x1b[1mFLAGS\x1b[0m",
				"  -h, --help           display this very informative message",
				"  -p, --portable       output measurements on separate lines",
				"  --device <string>    name or ID of the smart plug to measure",
				"  --username <string>  account username for Emporia",
				"  --password <string>  account password for Emporia",
				"  --version            print the current version of this build",
				"",
				"\x1b[1mOUTPUT\x1b[0m",
				"  Command output is printed as specified by the command",
				"  Time and energy usage information is output to stderr",
				"",
				"  Time is counted with seconds and measured by the time command",
				"  Usage is measured by the device and shown in joules and watts",
				"  Sure is the ratio of received-to-expected measurements",
				"",
				"\x1b[1mEXAMPLE\x1b[0m",
				fmt.Sprintf("  $ %s sleep 12", os.Args[0]),
				"         12.00 real         0.00 user         0.00 sys",
				"        922.63 joules      76.87 watts      100.0% sure",
				"",
				"",
			},
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			buff := bytes.Buffer{}
			PrintHelpMessage(&buff)
			expected := strings.Join(tt.expected, "\n")
			assert.Equal(t, expected, buff.String())
		})
	}
}
