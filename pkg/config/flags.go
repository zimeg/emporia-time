package config

import (
	"bytes"
	"flag"

	"github.com/zimeg/emporia-time/internal/errors"
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

// ParseFlags prepares the command using provided arguments
func ParseFlags(args []string) (cmd []string, flags Flags, err error) {
	flagset := flag.NewFlagSet("etime", flag.ContinueOnError)
	flagset.SetOutput(&bytes.Buffer{})
	flagset.BoolVar(&flags.Help, "h", false, "display this very informative message")
	flagset.BoolVar(&flags.Help, "help", false, "display this very informative message")
	flagset.BoolVar(&flags.Portable, "p", false, "output measurements on separate lines")
	flagset.BoolVar(&flags.Portable, "portable", false, "output measurements on separate lines")
	flagset.BoolVar(&flags.Version, "version", false, "print the current version of this build")
	flagset.StringVar(&flags.Device, "device", "", "device to measure usage for")
	flagset.StringVar(&flags.Password, "password", "", "account password for Emporia")
	flagset.StringVar(&flags.Username, "username", "", "account username for Emporia")
	err = flagset.Parse(args[1:])
	if err != nil {
		return []string{}, Flags{}, errors.Wrap(errors.ErrConfigFlag, err)
	}
	cmd = flagset.Args()
	if len(cmd) <= 0 {
		flags.Help = true
	}
	return cmd, flags, nil
}
