package terminal

import (
	"flag"
	"os"
)

// Flags holds command line flags specific to etime
type Flags struct {
	Device   string
	Username string
	Password string
	Portable bool
}

// Command contains the command line configurations
type Command struct {
	Args  []string // Args contains arguments to use in the provided program
	Flags Flags    // Flags holds command specific flags and configurations
}

// ParseFlags prepares the command using provided arguments
func ParseFlags(arguments []string) Command {
	if len(arguments) <= 1 {
		printHelpMessage()
		os.Exit(0)
	}

	var flags Flags

	flag.BoolVar(&flags.Portable, "p", false, "display measurements on separate lines")
	flag.BoolVar(&flags.Portable, "portable", false, "display measurements on separate lines")

	flag.StringVar(&flags.Device, "device", "", "device to measure usage for")
	flag.StringVar(&flags.Username, "username", "", "account username for Emporia")
	flag.StringVar(&flags.Password, "password", "", "account password for Emporia")

	flag.Usage = printHelpMessage
	flag.Parse()

	return Command{
		Args:  flag.Args(),
		Flags: flags,
	}
}
