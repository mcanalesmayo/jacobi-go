package jacobi

// Gets the new maxDiff value
func MaxMaxDiff(maxDiffB, maxDiffA float64) float64 {
	if maxDiffB > maxDiffA {
		return maxDiffB
	} else {
		return maxDiffA
	}
}
