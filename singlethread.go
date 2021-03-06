package jacobi

import (
	"github.com/mcanalesmayo/jacobi-go/model/matrix"
	"math"
)

// runSinglethreadedJacobi runs a single-threaded version of the jacobi method
func runSinglethreadedJacobi(initialValue float64, nDim int, maxIters int, tolerance float64, matrixType matrix.MatrixType) (matrix.Matrix, int, float64) {
	// The algorithm requires computing each grid cell as a 3x3 filter with no corners
	// Therefore, we need an aux matrix to keep the grid values in every iteration after computing new values
	var matA matrix.Matrix
	if matrixType == matrix.OneDimMatrixType {
		matA = matrix.NewOneDimMatrix(initialValue, nDim+2, matrix.Hot, matrix.Cold, matrix.Hot, matrix.Hot)
	} else {
		matA = matrix.NewTwoDimMatrix(initialValue, nDim+2, matrix.Hot, matrix.Cold, matrix.Hot, matrix.Hot, matrixType)
	}
	matB := matA.Clone(matrix.MatrixDef{
		matrix.Coords{0, 0, nDim + 1, nDim + 1},
		nDim + 2,
	}).(matrix.Matrix)

	matrixIters, nIters, maxDiff := nDim+1, 0, math.MaxFloat64

	for maxDiff > tolerance && nIters < maxIters {
		maxDiff = 0.0

		for i := 1; i < matrixIters; i++ {
			for j := 1; j < matrixIters; j++ {
				// Compute new value with 3x3 filter with no corners
				matB.SetCell(i, j, 0.2*(matA.GetCell(i, j)+matA.GetCell(i-1, j)+matA.GetCell(i+1, j)+matA.GetCell(i, j-1)+matA.GetCell(i, j+1)))
				maxDiff = math.Max(maxDiff, math.Abs(matA.GetCell(i, j)-matB.GetCell(i, j)))
			}
		}

		// Swap matrices
		matA, matB = matB, matA
		nIters++
	}

	return matA, nIters, maxDiff
}
