package jacobi

import (
	"github.com/mcanalesmayo/jacobi-go/model/matrix"
)

// RunJacobi runs the jacobi method to simulate the thermal transmission in a 2D space
func RunJacobi(initialValue float64, nDim int, maxIters int, tolerance float64, nThreads int, matrixType matrix.MatrixType) (matrix.Matrix, int, float64) {
	if nThreads == 1 {
		return runSinglethreadedJacobi(initialValue, nDim, maxIters, tolerance, matrixType)
	}
	return runMultithreadedJacobi(initialValue, nDim, maxIters, tolerance, nThreads, matrixType)
}
