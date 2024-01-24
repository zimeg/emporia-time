package program

import (
	"bytes"
	"flag"
)

// Flags holds command line flags specific to etime
type Flags struct {
	Device   string
	Help     bool
	Password string
	Portable bool
	Username string
	Version  bool
}

// Command contains the command line configurations
type Command struct {
	Args  []string // Args contains arguments to use in the provided program
	Flags Flags    // Flags holds command specific flags and configurations
}

// ParseFlags prepares the command using provided arguments
func ParseFlags(arguments []string) (Command, error) {
	var flagset = flag.NewFlagSet("etime", flag.ContinueOnError)
	var flags Flags

	flagset.BoolVar(&flags.Help, "h", false, "display this very informative message")
	flagset.BoolVar(&flags.Help, "help", false, "display this very informative message")
	flagset.BoolVar(&flags.Portable, "p", false, "output measurements on separate lines")
	flagset.BoolVar(&flags.Portable, "portable", false, "output measurements on separate lines")
	flagset.BoolVar(&flags.Version, "version", false, "print the current version of this build")

	flagset.StringVar(&flags.Device, "device", "", "device to measure usage for")
	flagset.StringVar(&flags.Password, "password", "", "account password for Emporia")
	flagset.StringVar(&flags.Username, "username", "", "account username for Emporia")

	flagset.SetOutput(&bytes.Buffer{})
	err := flagset.Parse(arguments[1:])
	if err != nil {
		return Command{}, err
	}
	commandArgs := flagset.Args()
	if len(commandArgs) <= 0 {
		flags.Help = true
	}
	return Command{Args: commandArgs, Flags: flags}, nil
}
