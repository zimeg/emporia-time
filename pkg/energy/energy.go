package energy

import (
	"time"
)

const HourToSeconds float64 = 3600
const KiloToUnit float64 = 1000

// EnergyResult contains the calculated energy measurements
type EnergyResult struct {
	Watts    float64
	Sureness float64
}

// ScaleKWhToWs converts kilowatt-hours to watt-seconds
func ScaleKWhToWs(kwh float64) float64 {
	return kwh * KiloToUnit * HourToSeconds
}

// ExtrapolateUsage scales the average measured energy rate over the elapsed
// time to account for missing measurements, returning est. watts and sureness
func ExtrapolateUsage(measurements []float64, duration time.Duration) EnergyResult {
	seconds := time.Duration.Seconds(duration)
	count := float64(len(measurements))
	sum := 0.0
	for _, mm := range measurements {
		sum += mm
	}

	if duration == 0 {
		return EnergyResult{Watts: 0, Sureness: 1}
	}
	if count == 0 || sum == 0 {
		return EnergyResult{Watts: 0, Sureness: 0}
	}

	// scale the summation across the entire duration
	estimated := sum * (seconds / count)

	// calculate the observed-to-expected measurement ratio
	sureness := count / seconds
	if count > seconds {
		sureness = 1.0
	}

	return EnergyResult{Watts: estimated, Sureness: sureness}
}
