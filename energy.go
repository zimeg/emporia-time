package main

var HourToSeconds float64 = 3600
var KiloToUnit float64 = 1000

func ScaleKWhToWs(kwh float64) float64 {
	return kwh * KiloToUnit * HourToSeconds
}

// ExtrapolateUsage scales the average measured energy rate over the elapsed
// time to account for missing measurements, returning est. watts and sureness
func ExtrapolateUsage(measurements []float64, durr float64) (float64, float64) {
	var sum float64 = 0
	for _, mm := range measurements {
		sum += mm
	}

	// cannot estimate an empty measurement
	measured := float64(len(measurements))
	if measured == 0 || sum == 0 {
		return 0.0, 0.0
	}

	// scale the summation across the entire duration
	estimated := sum * (durr / measured)

	// calculate the observed-to-expected measurement ratio
	sureness := measured / durr
	if measured > durr {
		sureness = 1.0
	}

	return estimated, sureness
}
