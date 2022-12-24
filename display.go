package main

import (
	"fmt"
	"os"
)

func outputUsage(watts float64, sureness float64) {
	fmt.Printf("%12.2f watt %11.1f%% sure\n", watts, sureness*100)
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
