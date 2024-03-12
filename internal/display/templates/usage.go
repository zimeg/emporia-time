package templates

import (
	"strings"
)

// UsageStatistics contains methods to retrieve information of command results
type UsageStatistics interface {
	GetReal() float64     // GetReal returns the real time of a command
	GetUser() float64     // GetUser returns the user time of a command
	GetSys() float64      // GetSys returns the sys time of a command
	GetJoules() float64   // GetJoules returns the joules in a result
	GetWatts() float64    // GetWatts returns the watts in a result
	GetSureness() float64 // GetSureness returns the sureness of a result
}

// FormatUsage arranges information about resource usage of a command
func FormatUsage(results UsageStatistics, isPortableFormat bool) (string, error) {
	var energyTemplate string
	switch isPortableFormat {
	case false:
		energyTemplate = strings.TrimSpace(`
{{12 | TimeF .GetReal}} real {{12 | TimeF .GetUser}} user {{12 | TimeF .GetSys}} sys
{{12 | Value .GetJoules}} joules {{10 | Value .GetWatts}} watts {{10 | Percent .GetSureness}}% sure`)
	case true:
		energyTemplate = strings.TrimSpace(`
real {{0 | Time .GetReal}}
user {{0 | Time .GetUser}}
sys {{0 | Time .GetSys}}
joules {{0 | Value .GetJoules}}
watts {{0 | Value .GetWatts}}
sure {{0 | Percent .GetSureness}}%`)
	}
	return templateBuilder(energyTemplate, results)
}
