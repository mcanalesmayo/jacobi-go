package jacobi

import (
	"math"
)

// Computes the new maxDiff value
func MaxDiff(prevMaxDiff, valA, valB float64) float64 {
	maxDiff, absDiff := prevMaxDiff, math.Abs(valB - valA)
	if (absDiff > maxDiff) {
		maxDiff = absDiff
	}

	return maxDiff
}