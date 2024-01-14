package display

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFormatSeconds(t *testing.T) {
	tests := map[string]struct {
		Seconds  float64
		Expected string
	}{
		"the zero time is preserved": {
			Seconds:  0.0,
			Expected: "0.00",
		},
		"single digit seconds are matched": {
			Seconds:  8.05,
			Expected: "8.05",
		},
		"many seconds passed without formatting": {
			Seconds:  12.8,
			Expected: "12.80",
		},
		"a minute and a few seconds passed": {
			Seconds:  64.46,
			Expected: "1:04.46",
		},
		"an amount of time greater than a minute": {
			Seconds:  222.22,
			Expected: "3:42.22",
		},
		"many minutes passed and are formatted": {
			Seconds:  1342.22,
			Expected: "22:22.22",
		},
		"slightly past one whole hour": {
			Seconds:  3663.36,
			Expected: "1:01:03.36",
		},
		"around a third of an hour from seconds": {
			Seconds:  5025.67,
			Expected: "1:23:45.67",
		},
		"multiple hours were measured": {
			Seconds:  52331.98,
			Expected: "14:32:11.98",
		},
		"multiple days were counted and shown in hours": {
			Seconds:  314159.27,
			Expected: "87:15:59.27",
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			actual := FormatSeconds(tt.Seconds)
			assert.Equalf(t, tt.Expected, actual, "Failed to format %.2f seconds", tt.Seconds)
		})
	}
}
