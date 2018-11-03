package jacobi

import (
	"math"
	"sync"
	"github.com/mcanalesmayo/jacobi-go/model/matrix"
)

type worker struct {
	// Define subproblem matrix
	matDef matrix.MatrixDef
	// For sharing values with threads working on adjacent submatrices
	toTop chan float64
	toBottom chan float64
	toRight chan float64
	toLeft chan float64
	fromTop chan float64
	fromBottom chan float64
	fromRight chan float64
	fromLeft chan float64
}

func (worker worker) mergeSubproblem(resMat *matrix.Matrix, subprobResMat matrix.Matrix) {
	coords := worker.matDef.Coords

	for i := coords.X0; i < coords.X1; i++ {
		for j := coords.Y0; j < coords.Y1; j++ {
			// Values are ordered by the sender
			(*resMat)[i][j] = subprobResMat[i][j]
		}
	}
}

func (worker worker) maxReduce(maxDiff float64) float64 {
	// TODO: Implement max reduce of all threads to compute actual maxDiff value
	return maxDiff
}

func (worker worker) sendBorderValues(mat matrix.Matrix) {
	coords := worker.matDef.Coords
	matLen := worker.matDef.Size

	// Since subproblem coordinates never change, this solution
	// isn't the best one in terms of performance, as these
	// checks are done for every jacobi iteration
	if coords.Y0 != 0 {
		for j := 0; j < matLen; j++ {
			worker.toTop <- mat[0][j]
		}
	}
	if coords.Y1 != 0 {
		for j := 0; j < matLen; j++ {
			worker.toBottom <- mat[matLen-1][j]
		}
	}
	if coords.X0 != 0 {
		for i := 0; i < matLen; i++ {
			worker.toRight <- mat[i][0]
		}
	}
	if coords.X1 != 0 {
		for i := 0; i < matLen; i++ {
			worker.toLeft <- mat[i][matLen-1]
		}
	}
}

func (worker worker) recvBorderValues(mat *matrix.Matrix) {
	coords := worker.matDef.Coords
	matLen := worker.matDef.Size

	// Since subproblem coordinates never change, this solution
	// isn't the best one in terms of performance, as these
	// checks are done for every jacobi iteration
	if coords.Y0 != 0 {
		for j := 0; j < matLen; j++ {
			(*mat)[0][j] = <- worker.fromTop
		}
	}
	if coords.Y1 != 0 {
		for j := 0; j < matLen; j++ {
			(*mat)[matLen-1][j] = <- worker.fromBottom
		}
	}
	if coords.X0 != 0 {
		for i := 0; i < matLen; i++ {
			(*mat)[i][0] = <- worker.fromRight
		}
	}
	if coords.X1 != 0 {
		for i := 0; i < matLen; i++ {
			(*mat)[i][matLen-1] = <- worker.fromLeft
		}
	}
}

func (worker worker) solveSubproblem(resMat *matrix.Matrix, initialValue float64, maxIters int, tolerance float64, wg *sync.WaitGroup) {
	defer wg.Done()

	var matA, matB matrix.Matrix
	initialValue, nDim, matDef := initialValue, worker.matDef.Size, worker.matDef
	// The algorithm requires computing each grid cell as a 3x3 filter with no corners
	// Therefore, we an aux matrix to keep the grid values in every iteration after computing new values
	matA = matrix.NewSubprobMatrix(initialValue, nDim+2, matDef)
	matB = matA.Clone()

	maxDiff, subprobSize := 1.0, worker.matDef.Size

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

		worker.recvBorderValues(&matA)
		// TODO: compute cells adjacent to outer cells
		// Actual max diff is maximum of all threads max diff
		maxDiff = worker.maxReduce(maxDiff)

		// Swap matrices
		matA, matB = matB, matA
	}

	worker.mergeSubproblem(resMat, matA)
}

func RunJacobiPar(initialValue float64, nDim int, maxIters int, tolerance float64, nThreads int) matrix.Matrix {
	// Resulting matrix
	resMat := matrix.NewMatrix(0.0, nDim)

	var wg sync.WaitGroup
	wg.Add(nThreads)

	for i := 0; i < nThreads; i++ {
		// TODO: Assign channels between threads
		go worker{

		}.solveSubproblem(&resMat, initialValue, maxIters, tolerance, &wg)
	}

	wg.Wait()

	// TODO: Return number of iterations and maximum diff
	return resMat
}