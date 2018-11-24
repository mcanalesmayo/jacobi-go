package utils

import (
	"math"
)

const (
	// Epsilon is the default maximum tolerance that can be used for the difference
	// between two fp numbers.
	Epsilon = 1.0e-9
)

// CompareFloats compares two fp numbers taking into account that fp representation
// results in approximations. Epsilon is taken as maximum tolerance for the difference
// between both fp numbers.
func CompareFloats(fpA, fpB float64, epsilon float64) bool {
	res := false

	if fpA == fpB {
		res = true
	} else if math.Abs(fpA-fpB) < epsilon {
		res = true
	}

	return res
}
