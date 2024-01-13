package templates

import (
	"fmt"
	"os"
	"regexp"
	"strings"
)

// PrintHelpMessage outputs an informative message for this program
func PrintHelpMessage() {
	helpTemplate := `
Measure the time and energy used while executing a command

{{ Bold "USAGE" }}
  {{ CommandName }} [flags] <command> [args]

{{ Bold "DESCRIPTION" }}
  flags    optional flags to provide this program
  command  the program to execute and measure
  args     optional arguments for the command

{{ Bold "FLAGS" }}
  -h, --help           display this very informative message
  -p, --portable       output measurements on separate lines
  --device <string>    name or ID of the smart plug to measure
  --username <string>  account username for Emporia
  --password <string>  account password for Emporia

{{ Bold "OUTPUT" }}
  Command output is printed as specified by the command
  Time and energy usage information is output to stderr

  Time is counted with seconds and measured by the time command
  Usage is measured by the device and shown in joules and watts
  Sure is the ratio of recieved-to-expected measurements

{{ Bold "EXAMPLE" }}
  $ {{ CommandName }} sleep 12
         12.00 real         0.00 user         0.00 sys
        922.63 joules      76.87 watts      100.0%% sure

`
	if body, err := templateBuilder(helpTemplate, nil); err != nil {
		boldRegex := regexp.MustCompile(`{{ Bold "([^"]+)" }}`)
		body = boldRegex.ReplaceAllString(helpTemplate, "$1")
		commandNameRegex := regexp.MustCompile(`{{ CommandName }}`)
		body = commandNameRegex.ReplaceAllString(body, "etime")
		fmt.Fprintf(os.Stderr, strings.TrimLeft(body, "\n"))
	} else {
		fmt.Fprintf(os.Stderr, strings.TrimLeft(body, "\n"))
	}
}
