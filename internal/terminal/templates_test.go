package terminal

import (
	"testing"
)

func TestFormatSeconds(t *testing.T) {
	tests := []struct {
		Title    string
		Seconds  float64
		Expected string
	}{
		{
			Title:    "the zero time is preserved",
			Seconds:  0.0,
			Expected: "0.00",
		},
		{
			Title:    "single digit seconds are matched",
			Seconds:  8.05,
			Expected: "8.05",
		},
		{
			Title:    "many seconds passed without formatting",
			Seconds:  12.8,
			Expected: "12.80",
		},
		{
			Title:    "a minute and a few seconds passed",
			Seconds:  64.46,
			Expected: "1:04.46",
		},
		{
			Title:    "an amount of time greater than a minute",
			Seconds:  222.22,
			Expected: "3:42.22",
		},
		{
			Title:    "many minutes passed and are formatted",
			Seconds:  1342.22,
			Expected: "22:22.22",
		},
		{
			Title:    "slightly past one whole hour",
			Seconds:  3663.36,
			Expected: "1:01:03.36",
		},
		{
			Title:    "around a third of an hour from seconds",
			Seconds:  5025.67,
			Expected: "1:23:45.67",
		},
		{
			Title:    "multiple hours were measured",
			Seconds:  52331.98,
			Expected: "14:32:11.98",
		},
		{
			Title:    "multiple days were counted and shown in hours",
			Seconds:  314159.27,
			Expected: "87:15:59.27",
		},
	}

	for _, tt := range tests {
		actual := FormatSeconds(tt.Seconds)
		if actual != tt.Expected {
			t.Errorf("A time is not formatted correctly!\nTEST: '%s'\nINPUT: %.2f\nEXPECT: %+v\nACTUAL: %+v",
				tt.Title,
				tt.Seconds,
				tt.Expected,
				actual,
			)
		}
	}
}
