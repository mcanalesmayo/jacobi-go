package utils

import (
	"math"
)

const (
	Epsilon = 1.0e-9
)

// MaxMaxDiff gets the new maxDiff value based on two maxDiff values
func MaxMaxDiff(maxDiffB, maxDiffA float64) float64 {
	var res float64

	if maxDiffB > maxDiffA {
		res = maxDiffB
	} else {
		res = maxDiffA
	}

	return res
}

// CompareFloats compares two fp numbers taking into account that fp representation
// results in approximations. Epsilon is taken as maximum tolerance for the difference
// between both fp numbers.
func CompareFloats(fpA, fpB float64, epsilon float64) bool {
	res := false

	if fpA == fpB {
		res = true
	} else if math.Abs(fpA - fpB) < epsilon {
		res = true
	}

	return res
}