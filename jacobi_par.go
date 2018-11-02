package jacobi

import (
	"math"
	"sync"
	"github.com/mcanalesmayo/jacobi-go/model/matrix"
)

type worker struct {
	subprob struct {
		coords struct {
			// Top-left corner
			x0 int
			y0 int
			// Bottom-right corner
			x1 int
			y1 int
		}
		// Precomputed size (x1 - x0) or (y1 - y0)
		size int
	}
	// For sharing values with threads working on adjacent submatrices
	toUpper chan float64
	toLower chan float64
	toRight chan float64
	toLeft chan float64
	fromUpper chan float64
	fromLower chan float64
	fromRight chan float64
	fromLeft chan float64
}

func (worker worker) mergeSubproblem(resMat *matrix.Matrix, subprobResMat matrix.Matrix) {
	coords := worker.subprob.coords

	for i := coords.x0; i < coords.x1; i++ {
		for j := coords.y0; j < coords.y1; j++ {
			// Values are ordered by the sender
			resMat[i][j] = subprobResMat[i][j]
		}
	}
}

func (worker worker) maxReduce(maxDiff float64) float64 {
	// TODO: Implement max reduce of all threads to compute actual maxDiff value
}

func (worker worker) sendBorderValues(mat matrix.Matrix) {
	coords := worker.subprob.coords
	matLen = worker.subprob.size

	// Since subproblem coordinates never change, this solution
	// isn't the best one in terms of performance, as these
	// checks are done for every jacobi iteration
	if coords.y0 != 0 {
		for j := 0; j < matLen; j++ {
			worker.toUpper <- mat[0][j]
		}
	}
	if coords.y1 != 0 {
		for j := 0; j < matLen; j++ {
			worker.toLower <- mat[matLen-1][j]
		}
	}
	if coords.x0 != 0 {
		for i := 0; i < matLen; i++ {
			worker.toLower <- mat[i][0]
		}
	}
	if coords.x1 != 0 {
		for i := 0; i < matLen; i++ {
			worker.toLower <- mat[i][matLen-1]
		}
	}
}

func (worker worker) recvBorderValues(mat *matrix.Matrix) {
	coords := worker.subprob.coords
	matLen = worker.subprob.size

	// Since subproblem coordinates never change, this solution
	// isn't the best one in terms of performance, as these
	// checks are done for every jacobi iteration
	if coords.y0 != 0 {
		for j := 0; j < matLen; j++ {
			mat[0][j] <- worker.fromUpper
		}
	}
	if coords.y1 != 0 {
		for j := 0; j < matLen; j++ {
			mat[matLen-1][j] <- worker.fromLower
		}
	}
	if coords.x0 != 0 {
		for i := 0; i < matLen; i++ {
			mat[i][0] <- worker.fromRight
		}
	}
	if coords.x1 != 0 {
		for i := 0; i < matLen; i++ {
			mat[i][matLen-1] <- worker.fromLeft
		}
	}
}

func (worker worker) solveSubproblem(resMat *matrix.Matrix) {
	var matA, matB matrix.Matrix
	// The algorithm requires computing each grid cell as a 3x3 filter with no corners
	// Therefore, we an aux matrix to keep the grid values in every iteration after computing new values
	matA = matrix.NewMatrix(initialValue, nDim+2)
	matB = matA.Clone()

	maxDiff, maxIters, subprobSize := 1.0, nDim + 1, worker.subprob.size

	for nIters := 0; maxDiff > tolerance && nIters < maxIters; nIters++ {
		maxDiff = 0.0

		worker.sendBorderValues(matA)

		for i := 1; i < subprobSize-1; i++ {
			for j := 1; j < subprobSize-1; j++ {
				// Compute new value with 3x3 filter with no corners
				matB[i][j] = 0.2*(matA[i][j] + matA[i-1][j] + matA[i+1][j] + matA[i][j-1] + matA[i][j+1])
				absDiff := math.Abs(matB[i][j] - matA[i][j])
				if (absDiff > maxDiff) {
					maxDiff = absDiff
				}
			}
		}

		worker.recvBorderValues(matA)
		// TODO: compute cells adjacent to outer cells
		// Actual max diff is maximum of all threads max diff
		maxDiff = worker.maxReduce(maxDiff)

		// Swap matrices
		matA, matB = matB, matA
	}

	worker.mergeSubproblem(resMat, matA)
}

func RunJacobi(initialValue float64, nDim int, maxIters int, tolerance float64, nThreads int) (matrix.Matrix, int, float64) {
	// Resulting matrix
	resMat := matrix.NewMatrix(0.0, nDim)

	var wg sync.WaitGroup
	wg.Add(nThreads)

	for i := range nThreads {
		// TODO: Assign channels between threads
		go worker{

		}.solveSubproblem(&resMat)
	}

	wg.Wait()

	return resMat, nIters, maxDiff
}