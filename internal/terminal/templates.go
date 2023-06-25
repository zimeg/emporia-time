package terminal

import (
	"fmt"
	"os"
	"regexp"
	"strings"
	"text/template"
)

// TemplateBuilder creates a string using a template and variables
func TemplateBuilder(templateStr string, body interface{}) (string, error) {
	funcs := template.FuncMap{
		"Bold": func(f string) string {
			return fmt.Sprintf("\x1b[1m%s\x1b[0m", f)
		},
		"CommandName": func() string {
			return os.Args[0]
		},
		"Percent": func(f float64, spacing int) string {
			return fmt.Sprintf("%*.1f", spacing, f*100)
		},
		"Time": func(f string, spacing int) string {
			return fmt.Sprintf("%*s", spacing, f)
		},
		"Value": func(f float64, spacing int) string {
			return fmt.Sprintf("%*.2f", spacing, f)
		},
	}

	tmpl, err := template.New("outputs").Funcs(funcs).Parse(templateStr)
	if err != nil {
		return "", err
	}
	var result strings.Builder
	if tmpl.Execute(&result, body) != nil {
		return "", err
	}
	formattedString := result.String()
	return formattedString, nil
}

// printHelpMessage outputs an informative message for this program
func printHelpMessage() {
	helpTemplate := `
Measure the time and energy used while executing a command

{{ Bold "USAGE" }}
  {{ CommandName }} [flags] <command> [args]

{{ Bold "DESCRIPTION" }}
  flags    optional flags to provide this program
  command  the program to execute and measure
  args     optional arguments for the command

{{ Bold "FLAGS" }}
  -h, --help           this very informative message
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
	if body, err := TemplateBuilder(helpTemplate, nil); err != nil {
		boldRegex := regexp.MustCompile(`{{ Bold "([^"]+)" }}`)
		body = boldRegex.ReplaceAllString(helpTemplate, "$1")
		commandNameRegex := regexp.MustCompile(`{{ CommandName }}`)
		body = commandNameRegex.ReplaceAllString(body, "etime")
		fmt.Fprintf(os.Stderr, strings.TrimLeft(body, "\n"))
	} else {
		fmt.Fprintf(os.Stderr, strings.TrimLeft(body, "\n"))
	}
}
