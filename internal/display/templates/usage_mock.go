package templates

type mockUsageStatistics struct {
	Real     float64
	User     float64
	Sys      float64
	Joules   float64
	Watts    float64
	Sureness float64
}

// GetReal returns the real time of a mocked command
func (results mockUsageStatistics) GetReal() float64 {
	return results.Real
}

// GetUser returns the user time of a mocked command
func (results mockUsageStatistics) GetUser() float64 {
	return results.User
}

// GetSys returns the sys time of a mocked command
func (results mockUsageStatistics) GetSys() float64 {
	return results.Sys
}

// GetJoules returns the joules in a mocked result
func (results mockUsageStatistics) GetJoules() float64 {
	return results.Joules
}

// GetWatts returns the watts in a mocked result
func (results mockUsageStatistics) GetWatts() float64 {
	return results.Watts
}

// GetSureness returns the sureness of a mocked result
func (results mockUsageStatistics) GetSureness() float64 {
	return results.Sureness
}
