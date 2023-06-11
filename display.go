package main

import (
	"fmt"
	"os"
	"strings"
	"text/template"
)

// templateBuilder creates a string using a template and variables
func templateBuilder(templateStr string, body interface{}) (string, error) {
	funcs := template.FuncMap{
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

// formatUsage arranges information about resource usage of a command
func formatUsage(results CommandResult) (string, error) {
	energyTemplate := strings.TrimSpace(`
{{12 | Time .TimeMeasurement.Command.Real}} real {{12 | Time .TimeMeasurement.Command.User}} user {{12 | Time .TimeMeasurement.Command.Sys}} sys
{{12 | Value .EnergyResult.Watts}} watt {{11 | Percent .EnergyResult.Sureness}}% sure`)

	body, err := templateBuilder(energyTemplate, results)
	if err != nil {
		return "", err
	}
	return body, nil
}

func outputHelp() {
	fmt.Printf("Measure the time and energy used while executing a command\n\n")

	fmt.Printf("%s\n", bold("USAGE"))
	fmt.Printf("  %s <command> [args]\n\n", os.Args[0])

	fmt.Printf("%s\n", bold("DESCRIPTION"))
	fmt.Printf("  <command>  the program to execute and measure\n")
	fmt.Printf("  [args]     optional arguments for the command\n\n")

	fmt.Printf("%s\n", bold("OUTPUT"))
	fmt.Printf("  Command output is printed as specified by the program\n")
	fmt.Printf("  Time and energy usage information is output to stderr\n\n")

	fmt.Printf("  Time is measured in seconds as defined by the time command\n")
	fmt.Printf("  Energy is measured in watts and collected from Emporia\n")
	fmt.Printf("  Sureness is the ratio of recieved-to-expected measurements\n\n")

	fmt.Printf("%s\n", bold("EXAMPLE"))
	fmt.Printf("  $ etime sleep 12\n")
	fmt.Printf("         12.00 real         0.00 user         0.00 sys\n")
	fmt.Printf("          9.53 watt        61.5%% sure\n\n")
}

func bold(str string) string {
	return "\x1b[1m" + str + "\x1b[0m"
}
