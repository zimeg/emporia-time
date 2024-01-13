package main

import (
	"fmt"
	"strings"
	"testing"

	"github.com/zimeg/emporia-time/internal/display"
	"github.com/zimeg/emporia-time/pkg/energy"
	"github.com/zimeg/emporia-time/pkg/times"
)

func TestFormatUsage_Formatted(t *testing.T) {
	tests := []struct {
		Title  string
		Result CommandResult
	}{
		{
			"the values of a zero command are output",
			CommandResult{
				TimeMeasurement: times.TimeMeasurement{
					Command: times.CommandTime{
						Real: 0.0,
						User: 0.0,
						Sys:  0.0,
					},
				},
				EnergyResult: energy.EnergyResult{
					Joules:   0.0,
					Watts:    0.0,
					Sureness: 0.0,
				},
			},
		},
		{
			"the values of a regular command are output",
			CommandResult{
				TimeMeasurement: times.TimeMeasurement{
					Command: times.CommandTime{
						Real: 4.0,
						User: 2.1,
						Sys:  1.4,
					},
				},
				EnergyResult: energy.EnergyResult{
					Joules:   12.00,
					Watts:    3.00,
					Sureness: 0.9620,
				},
			},
		},
		{
			"extended time values are shown as hh:mm:ss.ss",
			CommandResult{
				TimeMeasurement: times.TimeMeasurement{
					Command: times.CommandTime{
						Real: 75284924792499.01,
						User: 43200.56,
						Sys:  14400.12,
					},
				},
				EnergyResult: energy.EnergyResult{
					Joules:   18821231198124.75,
					Watts:    4.00,
					Sureness: 0.9620,
				},
			},
		},
	}

	for _, tt := range tests {
		output, err := formatUsage(tt.Result, false)
		if err != nil {
			t.Error("An unexpected error was encountered while formatting!")
		}

		if !strings.Contains(output, fmt.Sprintf("%s real", display.FormatSeconds(tt.Result.TimeMeasurement.Command.Real))) {
			t.Error("The `real` measurement is missing in the output!")
		}
		if !strings.Contains(output, fmt.Sprintf(" %s user", display.FormatSeconds(tt.Result.TimeMeasurement.Command.User))) {
			t.Error("The `user` measurement is missing in the output!")
		}
		if !strings.Contains(output, fmt.Sprintf(" %s sys", display.FormatSeconds(tt.Result.TimeMeasurement.Command.Sys))) {
			t.Error("The `sys` measurement is missing in the output!")
		}
		if !strings.Contains(output, fmt.Sprintf("%.2f joules", tt.Result.EnergyResult.Joules)) {
			t.Error("The `joules` measurement is missing in the output!")
		}
		if !strings.Contains(output, fmt.Sprintf(" %.2f watts", tt.Result.EnergyResult.Watts)) {
			t.Error("The `watts` measurement is missing in the output!")
		}
		if !strings.Contains(output, fmt.Sprintf(" %.1f%% sure", tt.Result.EnergyResult.Sureness*100)) {
			t.Error("The `sure` measurement is missing in the output!")
		}
	}
}

func TestFormatUsage_Portable(t *testing.T) {
	tests := []struct {
		Title  string
		Result CommandResult
	}{
		{
			"the values of a zero command are output on separate lines",
			CommandResult{
				TimeMeasurement: times.TimeMeasurement{
					Command: times.CommandTime{
						Real: 0.0,
						User: 0.0,
						Sys:  0.0,
					},
				},
				EnergyResult: energy.EnergyResult{
					Joules:   0.0,
					Watts:    0.0,
					Sureness: 0.0,
				},
			},
		},
		{
			"the values of a regular command are output on separate lines",
			CommandResult{
				TimeMeasurement: times.TimeMeasurement{
					Command: times.CommandTime{
						Real: 4.0,
						User: 2.1,
						Sys:  1.4,
					},
				},
				EnergyResult: energy.EnergyResult{
					Joules:   12.00,
					Watts:    3.00,
					Sureness: 0.9620,
				},
			},
		},
	}

	for _, tt := range tests {
		output, err := formatUsage(tt.Result, true)
		if err != nil {
			t.Error("An unexpected error was encountered while formatting!")
		}

		var (
			realTimeCount int
			userTimeCount int
			sysTimeCount  int
			joulesCount   int
			wattsCount    int
			surenessCount int
		)
		for _, line := range strings.Split(output, "\n") {
			switch line {
			case fmt.Sprintf("real %.2f", tt.Result.TimeMeasurement.Command.Real):
				realTimeCount += 1
			case fmt.Sprintf("user %.2f", tt.Result.TimeMeasurement.Command.User):
				userTimeCount += 1
			case fmt.Sprintf("sys %.2f", tt.Result.TimeMeasurement.Command.Sys):
				sysTimeCount += 1
			case fmt.Sprintf("joules %.2f", tt.Result.EnergyResult.Joules):
				joulesCount += 1
			case fmt.Sprintf("watts %.2f", tt.Result.EnergyResult.Watts):
				wattsCount += 1
			case fmt.Sprintf("sure %.1f%%", tt.Result.EnergyResult.Sureness*100):
				surenessCount += 1
			default:
				t.Error("An unexpected value appeared in the ouput:", line)
			}
		}
		if realTimeCount != 1 || userTimeCount != 1 || sysTimeCount != 1 ||
			joulesCount != 1 || wattsCount != 1 || surenessCount != 1 {
			t.Error("A measurement appeared an unexpected amount of times!")
		}
	}
}
