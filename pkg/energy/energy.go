package energy

import (
	"time"
)

// HourToSeconds multiplies units of hours to seconds
const HourToSeconds float64 = 3600

// KiloToUnit multiplies thousands to ones
const KiloToUnit float64 = 1000

// EnergyMeasurement holds response values of energy usage
type EnergyMeasurement struct {
	Chart    []float64     // Chart contains energy measurements as watt seconds
	Duration time.Duration // Duration is the amount of time used running a command
}

// EnergyResult contains the calculated energy measurements
type EnergyResult struct {
	Joules   float64       // Joules is the total energy used during the duration
	Watts    float64       // Watts is the average power output over the duration
	Sureness float64       // Sureness is a confidence score for resulting energy
	Duration time.Duration // Duration is the amount of time used running a command
}

// GetJoules returns the joules in a result
func (result EnergyResult) GetJoules() float64 {
	return result.Joules
}

// GetWatts returns the watts in a result
func (result EnergyResult) GetWatts() float64 {
	return result.Watts
}

// GetSureness returns the sureness of a result
func (result EnergyResult) GetSureness() float64 {
	return result.Sureness
}

// ScaleKWhToWs converts kilowatt-hours to watt-seconds
func ScaleKWhToWs(kwh float64) float64 {
	return kwh * KiloToUnit * HourToSeconds
}

// ExtrapolateUsage scales the measured energy response over the actual duration
//
// The measured joule summation is scaled across the actual timed duration using
// a ratio of observed-to-expected measurements
func ExtrapolateUsage(measurements EnergyMeasurement) EnergyResult {
	actualSeconds := measurements.Duration.Seconds()
	measuredSeconds := float64(len(measurements.Chart))
	measuredJoules := 0.0
	for _, measuredWattSecond := range measurements.Chart {
		measuredJoules += measuredWattSecond
	}
	if actualSeconds == 0.0 {
		return EnergyResult{Joules: 0, Watts: 0, Sureness: 1}
	}
	if measuredSeconds == 0 || measuredJoules == 0 {
		return EnergyResult{Joules: 0, Watts: 0, Sureness: 0}
	}
	estimatedJoules := measuredJoules * (actualSeconds / measuredSeconds)
	estimatedWatts := estimatedJoules / actualSeconds
	sureness := measuredSeconds / actualSeconds
	if measuredSeconds > actualSeconds {
		sureness = 1.0
	}
	return EnergyResult{
		Joules:   estimatedJoules,
		Watts:    estimatedWatts,
		Sureness: sureness,
	}
}
