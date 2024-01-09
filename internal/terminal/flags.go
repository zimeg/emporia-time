package terminal

import (
	"flag"
)

// Flags holds command line flags specific to etime
type Flags struct {
	Device   string
	Help     bool
	Password string
	Portable bool
	Username string
}

// Command contains the command line configurations
type Command struct {
	Args  []string // Args contains arguments to use in the provided program
	Flags Flags    // Flags holds command specific flags and configurations
}

// ParseFlags prepares the command using provided arguments
func ParseFlags(arguments []string) Command {
	var flagset flag.FlagSet
	var flags Flags

	flagset.BoolVar(&flags.Help, "h", false, "display this very informative message")
	flagset.BoolVar(&flags.Help, "help", false, "display this very informative message")
	flagset.BoolVar(&flags.Portable, "p", false, "output measurements on separate lines")
	flagset.BoolVar(&flags.Portable, "portable", false, "output measurements on separate lines")

	flagset.StringVar(&flags.Device, "device", "", "device to measure usage for")
	flagset.StringVar(&flags.Password, "password", "", "account password for Emporia")
	flagset.StringVar(&flags.Username, "username", "", "account username for Emporia")

	flagset.Usage = PrintHelpMessage
	flagset.Parse(arguments[1:])

	return Command{
		Args:  flagset.Args(),
		Flags: flags,
	}
}
