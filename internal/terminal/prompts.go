package terminal

import (
	"flag"
	"os"

	"github.com/AlecAivazis/survey/v2"
	"github.com/AlecAivazis/survey/v2/terminal"
	"github.com/zimeg/emporia-time/internal/errors"
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
		err := survey.AskOne(&survey.Input{Message: prompt.Message}, &value)
		if err != nil {
			if errors.Is(err, terminal.InterruptErr) {
				os.Exit(TerminalInterruptCode)
			}
			return "", errors.Wrap(errors.ErrPromptInput, err)
		}
	case true:
		err := survey.AskOne(&survey.Password{Message: prompt.Message}, &value)
		if err != nil {
			if errors.Is(err, terminal.InterruptErr) {
				os.Exit(TerminalInterruptCode)
			}
			return "", errors.Wrap(errors.ErrPromptInput, err)
		}
	}
	return value, nil
}

// CollectSelect gathers the index of the selected value for a prompt
func CollectSelect(prompt Prompt) (int, error) {
	switch {
	case len(prompt.Options) == 0:
		return 0, errors.New(errors.ErrPromptSelectMissing)
	case prompt.Descriptions != nil && len(prompt.Options) != len(prompt.Descriptions):
		return 0, errors.New(errors.ErrPromptSelectDescription)
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
	err := survey.AskOne(&question, &selectedIndex)
	if err != nil {
		if errors.Is(err, terminal.InterruptErr) {
			os.Exit(TerminalInterruptCode)
		}
		return 0, errors.Wrap(errors.ErrPromptSelect, err)
	}
	return selectedIndex, nil
}
