package energy

import (
	"testing"
	"time"
)

func TestScaleKWhToWs(t *testing.T) {
	tests := []struct {
		Title      string
		KWh        float64
		ExpectedWs float64
	}{
		{
			"ensure the zero value is zero",
			0,
			0,
		},
		{
			"convert a single KWh to Ws",
			1,
			1 * 1000 * 3600,
		},
		{
			"convert 1000 KWh to Ws",
			1000,
			1000 * 1000 * 3600,
		},
	}

	for _, tt := range tests {
		actual := ScaleKWhToWs(tt.KWh)
		if tt.ExpectedWs != actual {
			t.Fatalf("An unexpected energy conversion was found!\nTEST: '%s'\nEXPECT: %8f\nACTUAL: %8f",
				tt.Title,
				tt.ExpectedWs,
				actual,
			)
		}
	}
}

func TestExtrapolateUsage(t *testing.T) {
	tests := []struct {
		Title          string
		Measurements   []float64
		Duration       time.Duration
		ExpectedResult EnergyResult
	}{
		{
			"handle the measurements of instant commands",
			[]float64{0},
			time.Duration(0 * float64(time.Second)),
			EnergyResult{Watts: 0, Sureness: 1},
		},
		{
			"return unsure results if no measurements are returned",
			[]float64{},
			time.Duration(3 * float64(time.Second)),
			EnergyResult{Watts: 0, Sureness: 0},
		},
		{
			"return unsure results if all measurements are zero",
			[]float64{0, 0, 0},
			time.Duration(3 * float64(time.Second)),
			EnergyResult{Watts: 0, Sureness: 0},
		},
		{
			"confidently compute results for complete measurements",
			[]float64{3.64, 4.2, 2}, // sum=9.84, avg=3.28
			time.Duration(3 * float64(time.Second)),
			EnergyResult{Watts: 9.84, Sureness: 1},
		},
		{
			"extrapolate a missing second of measured results",
			[]float64{3, 4, 6, 3}, // sum=16, avg=4
			time.Duration(5 * float64(time.Second)),
			EnergyResult{Watts: 20, Sureness: 0.8},
		},
		{
			"extrapolate a half second of undermeasured results",
			[]float64{3, 4, 6, 3}, // sum=16, avg=4
			time.Duration(4.5 * float64(time.Second)),
			EnergyResult{Watts: 18, Sureness: 4 / 4.5},
		},
		{
			"extrapolate a half second of overmeasured results",
			[]float64{3, 4, 6, 3, 4}, // sum=20, avg=4
			time.Duration(4.5 * float64(time.Second)),
			EnergyResult{Watts: 18, Sureness: 1},
		},
	}

	for _, tt := range tests {
		actual := ExtrapolateUsage(tt.Measurements, tt.Duration)
		if tt.ExpectedResult.Watts != actual.Watts {
			t.Fatalf("An unexpected watt estimation was found!\nFAILED: '%s'\nEXPECT: %8f\nACTUAL: %8f",
				tt.Title,
				tt.ExpectedResult.Watts,
				actual.Watts,
			)
		}
		if tt.ExpectedResult.Sureness != actual.Sureness {
			t.Fatalf("An unexpected sureness score was found!\nFAILED: '%s'\nEXPECT: %8f\nACTUAL: %8f",
				tt.Title,
				tt.ExpectedResult.Sureness,
				actual.Sureness,
			)
		}
	}
}
