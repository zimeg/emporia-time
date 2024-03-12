package templates

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zimeg/emporia-time/internal/display"
)

func TestFormatUsage_Formatted(t *testing.T) {
	tests := map[string]mockUsageStatistics{
		"values of zero commands are output": {
			Real:     0.0,
			User:     0.0,
			Sys:      0.0,
			Joules:   0.0,
			Watts:    0.0,
			Sureness: 0.0,
		},
		"values of regular commands are output": {
			Real:     4.0,
			User:     2.1,
			Sys:      1.4,
			Joules:   12.00,
			Watts:    3.00,
			Sureness: 0.9620,
		},
		"extended time values are shown as hh:mm:ss.ss": {
			Real:     75284924792499.01,
			User:     43200.56,
			Sys:      14400.12,
			Joules:   18821231198124.75,
			Watts:    4.00,
			Sureness: 0.9620,
		},
	}
	for name, mockResults := range tests {
		t.Run(name, func(t *testing.T) {
			output, err := FormatUsage(mockResults, false)
			assert.NoError(t, err)
			assert.Contains(t, output, fmt.Sprintf("%s real", display.FormatSeconds(mockResults.Real)))
			assert.Contains(t, output, fmt.Sprintf(" %s user", display.FormatSeconds(mockResults.User)))
			assert.Contains(t, output, fmt.Sprintf(" %s sys", display.FormatSeconds(mockResults.Sys)))
			assert.Contains(t, output, fmt.Sprintf("%.2f joules", mockResults.Joules))
			assert.Contains(t, output, fmt.Sprintf("%.2f watts", mockResults.Watts))
			assert.Contains(t, output, fmt.Sprintf(" %.1f%% sure", mockResults.Sureness*100))
		})
	}
}

func TestFormatUsage_Portable(t *testing.T) {
	tests := map[string]mockUsageStatistics{
		"values of zero commands output on separate lines": {
			Real:     0.0,
			User:     0.0,
			Sys:      0.0,
			Joules:   0.0,
			Watts:    0.0,
			Sureness: 0.0,
		},
		"values of regular commands output on separate lines": {
			Real:     4.0,
			User:     2.1,
			Sys:      1.4,
			Joules:   12.00,
			Watts:    3.00,
			Sureness: 0.9620,
		},
	}
	for name, mockResults := range tests {
		t.Run(name, func(t *testing.T) {
			var (
				realTimeCount int
				userTimeCount int
				sysTimeCount  int
				joulesCount   int
				wattsCount    int
				surenessCount int
			)
			output, err := FormatUsage(mockResults, true)
			assert.NoError(t, err)
			for _, line := range strings.Split(output, "\n") {
				switch line {
				case fmt.Sprintf("real %.2f", mockResults.Real):
					realTimeCount += 1
				case fmt.Sprintf("user %.2f", mockResults.User):
					userTimeCount += 1
				case fmt.Sprintf("sys %.2f", mockResults.Sys):
					sysTimeCount += 1
				case fmt.Sprintf("joules %.2f", mockResults.Joules):
					joulesCount += 1
				case fmt.Sprintf("watts %.2f", mockResults.Watts):
					wattsCount += 1
				case fmt.Sprintf("sure %.1f%%", mockResults.Sureness*100):
					surenessCount += 1
				default:
					t.Error("An unexpected value appeared in the output:", line)
				}
			}
			assert.Equal(t, 1, realTimeCount)
			assert.Equal(t, 1, userTimeCount)
			assert.Equal(t, 1, sysTimeCount)
			assert.Equal(t, 1, joulesCount)
			assert.Equal(t, 1, wattsCount)
			assert.Equal(t, 1, surenessCount)
		})
	}
}
