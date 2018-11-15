package utils

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
