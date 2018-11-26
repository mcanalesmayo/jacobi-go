package jacobi

import (
	"github.com/mcanalesmayo/jacobi-go/model/matrix"
	"math"
)

// runSinglethreadedJacobi runs a single-threaded version of the jacobi method
func runSinglethreadedJacobi(matrixType matrix.MatrixType, initialValue float64, nDim int, maxIters int, tolerance float64) (matrix.Matrix, int, float64) {
	// The algorithm requires computing each grid cell as a 3x3 filter with no corners
	// Therefore, we need an aux matrix to keep the grid values in every iteration after computing new values
	var matA matrix.Matrix
	if matrixType == matrix.TwoDimMatrixType {
		matA = matrix.NewTwoDimMatrix(initialValue, nDim+2, matrix.Hot, matrix.Cold, matrix.Hot, matrix.Hot)
	} else {
		matA = matrix.NewOneDimMatrix(initialValue, nDim+2, matrix.Hot, matrix.Cold, matrix.Hot, matrix.Hot)
	}
	matB := matA.Clone(matrix.MatrixDef{
		matrix.Coords{0, 0, nDim + 1, nDim + 1},
		nDim + 2,
	})

	matrixIters, nIters, maxDiff := nDim+1, 0, math.MaxFloat64

	for maxDiff > tolerance && nIters < maxIters {
		maxDiff = 0.0

		for i := 1; i < matrixIters; i++ {
			for j := 1; j < matrixIters; j++ {
				// Compute new value with 3x3 filter with no corners
				matB.Set(i, j, 0.2 * (matA.Get(i, j) + matA.Get(i-1, j) + matA.Get(i+1, j) + matA.Get(i, j-1) + matA.Get(i, j+1)))
				maxDiff = math.Max(maxDiff, math.Abs(matA.Get(i, j)-matB.Get(i, j)))
			}
		}

		// Swap matrices
		matA, matB = matB, matA
		nIters++
	}

	return matA, nIters, maxDiff
}
