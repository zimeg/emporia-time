package main

import (
	"fmt"
	"strings"
	"testing"

	"github.com/zimeg/emporia-time/pkg/energy"
	"github.com/zimeg/emporia-time/pkg/times"
)

func TestFormatUsage_Regular(t *testing.T) {
	tests := []struct {
		Title  string
		Result CommandResult
	}{
		{
			"the values of a zero command are output",
			CommandResult{
				TimeMeasurement: times.TimeMeasurement{
					Command: times.CommandTime{
						Real: "0.00",
						User: "0.00",
						Sys:  "0.00",
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
						Real: "4.00",
						User: "2.10",
						Sys:  "1.40",
					},
				},
				EnergyResult: energy.EnergyResult{
					Joules:   12.00,
					Watts:    3.00,
					Sureness: .9620,
				},
			},
		},
	}

	for _, tt := range tests {
		output, err := formatUsage(tt.Result, false)
		if err != nil {
			t.Error("An unexpected error was encountered while formatting!")
		}

		if !strings.Contains(output, fmt.Sprintf(" %s real", tt.Result.TimeMeasurement.Command.Real)) {
			t.Error("The `real` measurement is missing in the output!")
		}
		if !strings.Contains(output, fmt.Sprintf(" %s user", tt.Result.TimeMeasurement.Command.User)) {
			t.Error("The `user` measurement is missing in the output!")
		}
		if !strings.Contains(output, fmt.Sprintf(" %s sys", tt.Result.TimeMeasurement.Command.Sys)) {
			t.Error("The `sys` measurement is missing in the output!")
		}
		if !strings.Contains(output, fmt.Sprintf(" %.2f joules", tt.Result.EnergyResult.Joules)) {
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
						Real: "0.00",
						User: "0.00",
						Sys:  "0.00",
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
						Real: "4.00",
						User: "2.10",
						Sys:  "1.40",
					},
				},
				EnergyResult: energy.EnergyResult{
					Joules:   12.00,
					Watts:    3.00,
					Sureness: .9620,
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
			realTimeFound int
			userTimeFound int
			sysTimeFound  int
			joulesFound   int
			wattsFound    int
			surenessFound int
		)
		for _, line := range strings.Split(output, "\n") {
			switch line {
			case fmt.Sprintf("%s real", tt.Result.TimeMeasurement.Command.Real):
				realTimeFound += 1
			case fmt.Sprintf("%s user", tt.Result.TimeMeasurement.Command.User):
				userTimeFound += 1
			case fmt.Sprintf("%s sys", tt.Result.TimeMeasurement.Command.Sys):
				sysTimeFound += 1
			case fmt.Sprintf("%.2f joules", tt.Result.EnergyResult.Joules):
				joulesFound += 1
			case fmt.Sprintf("%.2f watts", tt.Result.EnergyResult.Watts):
				wattsFound += 1
			case fmt.Sprintf("%.1f%% sure", tt.Result.EnergyResult.Sureness*100):
				surenessFound += 1
			default:
				t.Error("An unexpected value appeared in the ouput:", line)
			}
		}
		if realTimeFound != 1 || userTimeFound != 1 || sysTimeFound != 1 ||
			joulesFound != 1 || wattsFound != 1 || surenessFound != 1 {
			t.Error("A measurement appeared an unexpected amount of times!")
		}
	}
}
