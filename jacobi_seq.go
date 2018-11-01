package jacobi

import (
	"math"
	"github.com/mcanalesmayo/jacobi-go/model/matrix"
)

func RunJacobi(initialValue float64, nDim int, maxIters int, tolerance float64) (matrix.Matrix, int, float64) {
	var matA, matB matrix.Matrix
	// The algorithm requires computing each grid cell as a 3x3 filter with no corners
	// Therefore, we an aux matrix to keep the grid values in every iteration after computing new values
	matA = matrix.NewMatrix(initialValue, nDim+2)
	matB = matA.Clone()

	matrixIters := nDim + 1

	nIters, maxDiff := 0, 1.0

	for maxDiff > tolerance && nIters < maxIters {
		maxDiff = 0.0

		for i := 1; i < matrixIters; i++ {
			for j := 1; j < matrixIters; j++ {
				// Compute new value with 3x3 filter with no corners
				matB[i][j] = 0.2*(matA[i][j] + matA[i-1][j] + matA[i+1][j] + matA[i][j-1] + matA[i][j+1])
				absDiff := math.Abs(matB[i][j] - matA[i][j])
				if (absDiff > maxDiff) {
					maxDiff = absDiff
				}
			}
		}

		// Swap matrices
		matA, matB = matB, matA
		nIters += 1
	}

	return matA, nIters, maxDiff
}