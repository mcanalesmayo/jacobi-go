package utils

import (
	"math"
)

const (
	// Epsilon is the default maximum tolerance that can be used for comparing
	// if two fp numbers are equal.
	Epsilon = 1.0e-9
)

// CompareFloats compares two fp numbers taking into account that fp representation
// results in approximations.
// epsilon is taken as maximum absoluteValue(fpA-fpB).
func CompareFloats(fpA, fpB float64, epsilon float64) bool {
	res := false

	if fpA == fpB {
		res = true
	} else if math.Abs(fpA-fpB) < epsilon {
		res = true
	}

	return res
}
