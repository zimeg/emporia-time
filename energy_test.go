package main

import (
	"testing"
)

func TestConvertKWhToWZero(t *testing.T) {
	var kWh float64 = 0
	var expected float64 = 0
	actual := ScaleKWhToWs(kWh)
	if actual != expected {
		t.Fatalf("Incorrect conversion: %8fkWh != %8fW", kWh, expected)
	}
}

func TestConvertKWhToWsUnit(t *testing.T) {
	var kWh float64 = 1
	var expected float64 = 1 * 1000 * 3600
	actual := ScaleKWhToWs(kWh)
	if actual != expected {
		t.Fatalf("Incorrect conversion: %8fkWh == %8fW != %8fW", kWh, expected, actual)
	}
}

func TestConvertKWhToWKilo(t *testing.T) {
	var kWh float64 = 1000
	var expected float64 = 1000 * 1000 * 3600
	actual := ScaleKWhToWs(kWh)
	if actual != expected {
		t.Fatalf("Incorrect conversion: %8fkWh == %8fW != %8fW", kWh, expected, actual)
	}
}

func TestExtrapolateUsageEmpty(t *testing.T) {
	var measurements []float64 = []float64{}
	var durr float64 = 3

	// scaled measurements over the duration
	var expected float64 = 0.0
	var confidence float64 = 0.0

	estimated, sureness := ExtrapolateUsage(measurements, durr)
	if estimated != expected {
		t.Fatalf("Incorrect estimation when empty: %8f != %8f", estimated, expected)
	}

	if sureness != confidence {
		t.Fatalf("Incorrect sureness when empty: %8f != %8f", sureness, confidence)
	}
}

func TestExtrapolateUsageZero(t *testing.T) {
	// occasionally all measurements return as zero, so usage cannot be predicted
	var measurements []float64 = []float64{0, 0, 0}
	var durr float64 = 3

	// scaled measurements over the duration
	var expected float64 = 0.0
	var confidence float64 = 0.0

	estimated, sureness := ExtrapolateUsage(measurements, durr)
	if estimated != expected {
		t.Fatalf("Incorrect estimation when empty: %8f != %8f", estimated, expected)
	}

	if sureness != confidence {
		t.Fatalf("Incorrect sureness when empty: %8f != %8f", sureness, confidence)
	}
}

func TestExtrapolateUsageComplete(t *testing.T) {
	var measurements []float64 = []float64{3.64, 4.2, 2} // sum=9.84, avg=3.28
	var durr float64 = 3

	// scaled measurements over the duration
	var expected float64 = 9.84
	var confidence float64 = 1.0

	estimated, sureness := ExtrapolateUsage(measurements, durr)
	if estimated != expected {
		t.Fatalf("Incorrect estimation when empty: %8f != %8f", estimated, expected)
	}

	if sureness != confidence {
		t.Fatalf("Incorrect sureness when empty: %8f != %8f", sureness, confidence)
	}
}

func TestExtrapolateUsageUnderMeasured(t *testing.T) {
	var measurements []float64 = []float64{3, 4, 6, 3} // sum=16, avg=4
	var durr float64 = 5

	// scaled measurements over the duration
	var expected float64 = 20
	var confidence float64 = 0.8

	estimated, sureness := ExtrapolateUsage(measurements, durr)
	if estimated != expected {
		t.Fatalf("Incorrect estimation when empty: %8f != %8f", estimated, expected)
	}

	if sureness != confidence {
		t.Fatalf("Incorrect sureness when empty: %8f != %8f", sureness, confidence)
	}
}

func TestExtrapolateUsageHalfSecondUnderMeasured(t *testing.T) {
	var measurements []float64 = []float64{3, 4, 6, 3} // sum=16, avg=4
	var durr float64 = 4.5

	// scaled measurements over the duration
	var expected float64 = 18
	var confidence float64 = 4 / 4.5

	estimated, sureness := ExtrapolateUsage(measurements, durr)
	if estimated != expected {
		t.Fatalf("Incorrect estimation when empty: %8f != %8f", estimated, expected)
	}

	if sureness != confidence {
		t.Fatalf("Incorrect sureness when empty: %8f != %8f", sureness, confidence)
	}
}

func TestExtrapolateUsageHalfSecondOverMeasured(t *testing.T) {
	var measurements []float64 = []float64{3, 4, 6, 3, 4} // sum=20, avg=4
	var durr float64 = 4.5

	// scaled measurements over the duration
	var expected float64 = 18
	var confidence float64 = 1.0

	estimated, sureness := ExtrapolateUsage(measurements, durr)
	if estimated != expected {
		t.Fatalf("Incorrect estimation when empty: %8f != %8f", estimated, expected)
	}

	if sureness != confidence {
		t.Fatalf("Incorrect sureness when empty: %8f != %8f", sureness, confidence)
	}
}
