package terminal

import (
	"errors"
	"flag"
	"os"

	"github.com/AlecAivazis/survey/v2"
	"github.com/AlecAivazis/survey/v2/terminal"
)

// TerminalInterruptCode is the exit status for interrupted prompts
const TerminalInterruptCode = 130

// Prompt holds the information needed for gathering inputs
type Prompt struct {
	Message     string     // Promptful message for the input
	Flag        *flag.Flag // Substituting flag, preferred if set
	Environment string     // Related environment variable name

	Options      []string // Available selection options
	Descriptions []string // Optional option descriptions
	Hidden       bool     // Mask any text inputs
}

// CollectInput gathers a value for the prompt from flag, environment, or input
func CollectInput(prompt *Prompt) (string, error) {
	var value string
	if flag := prompt.Flag; flag != nil && flag.Value.String() != "" {
		return flag.Value.String(), nil
	} else if value = os.Getenv(prompt.Environment); value != "" {
		return value, nil
	}
	switch prompt.Hidden {
	case false:
		if err := survey.AskOne(&survey.Input{Message: prompt.Message}, &value); err != nil {
			if err == terminal.InterruptErr {
				os.Exit(TerminalInterruptCode)
			}
			return "", err
		}
	case true:
		if err := survey.AskOne(&survey.Password{Message: prompt.Message}, &value); err != nil {
			if err == terminal.InterruptErr {
				os.Exit(TerminalInterruptCode)
			}
			return "", err
		}
	}
	return value, nil
}

// CollectInput gathers a value for the prompt from flag, environment, or select
func CollectSelect(prompt Prompt) (int, error) {
	switch {
	case len(prompt.Options) == 0:
		return 0, errors.New("No options to select from!")
	case prompt.Descriptions != nil && len(prompt.Options) != len(prompt.Descriptions):
		return 0, errors.New("Mismatched option and description count for select")
	}

	question := survey.Select{
		Message: prompt.Message,
		Options: prompt.Options,
	}
	if prompt.Descriptions != nil {
		question.Description = func(value string, index int) string {
			return prompt.Descriptions[index]
		}
	}

	var selectedIndex int
	if err := survey.AskOne(&question, &selectedIndex); err != nil {
		if err == terminal.InterruptErr {
			os.Exit(TerminalInterruptCode)
		}
		return 0, err
	}
	return selectedIndex, nil
}
