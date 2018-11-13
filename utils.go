package jacobi

// MaxMaxDiff gets the new maxDiff value based on two maxDiff values
func MaxMaxDiff(maxDiffB, maxDiffA float64) float64 {
	if maxDiffB > maxDiffA {
		return maxDiffB
	} else {
		return maxDiffA
	}
}
