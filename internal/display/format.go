package display

import (
	"fmt"
	"math"
)

// FormatSeconds converts seconds into the hh:mm:ss.ss without leading zeros
func FormatSeconds(seconds float64) string {
	h := math.Floor(seconds / 3600)
	m := math.Floor((seconds - (h * 3600)) / 60)
	s := seconds - (h * 3600) - (m * 60)
	switch {
	case h > 0:
		return fmt.Sprintf("%d:%02d:%05.2f", int(h), int(m), s)
	case m > 0:
		return fmt.Sprintf("%d:%05.2f", int(m), s)
	default:
		return fmt.Sprintf("%.2f", s)
	}
}
