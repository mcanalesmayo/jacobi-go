package jacobi

import (
	"sync"
	"github.com/mcanalesmayo/jacobi-go/model/matrix"
)

type worker struct {
	id int
	// Define subproblem matrix
	matDef matrix.MatrixDef
	// Number of global workers
	globalParams struct {
		nWorkers, size int
	}
	// For sharing values with threads working on adjacent submatrices
	toTop, toBottom, toRight, toLeft chan float64
	fromTop, fromBottom, fromRight, fromLeft chan float64
	// For reducing maxDiff
	maxDiffRes []chan float64
}

// Merges the worker subproblem resulting matrix into the global resulting matrix
func (worker worker) mergeSubproblem(resMat, subprobResMat matrix.Matrix) {
	coords := worker.matDef.Coords

	for i := coords.X0; i < coords.X1; i++ {
		for j := coords.Y0; j < coords.Y1; j++ {
			// Values are ordered by the sender
			resMat[i][j] = subprobResMat[i][j]
		}
	}
}

// For the sake of simplicity, reduction is centralized on the 'root' worker, which will fan out the resulting value
// TODO: Look into a better way to do a parallel reduce
func (worker worker) maxReduce(maxDiff float64) float64 {
	isRoot := worker.id == 0

	// maximum maxDiff found at this point
	var maxMaxDiff float64
	if isRoot {
		// Reduction centralized in the 'root' worker
		// Collect and reduce maxDiff values from all workers
		maxMaxDiff = maxDiff
		for i := 0; i < worker.globalParams.nWorkers - 1; i++ {
			otherMaxDiff := <- worker.maxDiffRes[i]
			if otherMaxDiff > maxMaxDiff {
				maxMaxDiff = otherMaxDiff
			}
		}
		
		// Fan out the result
		for i := 0; i < worker.globalParams.nWorkers - 1; i++ {
			worker.maxDiffRes[i] <- maxMaxDiff
		}
	} else {
		// 'Non-root' workers send their results
		worker.maxDiffRes[worker.id] <- maxDiff
		// Wait for result
		maxMaxDiff = <- worker.maxDiffRes[worker.id]
	}

	return maxMaxDiff
}

// Sends the worker outer values to adjacent workers
func (worker worker) sendOuterCells(mat matrix.Matrix) {
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
	if coords.Y1 != worker.globalParams.size {
		for j := 0; j < matLen; j++ {
			worker.toBottom <- mat[matLen-1][j]
		}
	}
	if coords.X0 != 0 {
		for i := 0; i < matLen; i++ {
			worker.toRight <- mat[i][0]
		}
	}
	if coords.X1 != worker.globalParams.size {
		for i := 0; i < matLen; i++ {
			worker.toLeft <- mat[i][matLen-1]
		}
	}
}

// Gets the adjacent workers outer values
func (worker worker) recvAdjacentCells(mat matrix.Matrix) {
	coords := worker.matDef.Coords
	matLen := worker.matDef.Size

	if coords.Y0 != 0 {
		for j := 0; j < matLen; j++ {
			mat[0][j] = <- worker.fromTop
		}
	}
	if coords.Y1 != worker.globalParams.size {
		for j := 0; j < matLen; j++ {
			mat[matLen-1][j] = <- worker.fromBottom
		}
	}
	if coords.X0 != 0 {
		for i := 0; i < matLen; i++ {
			mat[i][0] = <- worker.fromRight
		}
	}
	if coords.X1 != worker.globalParams.size {
		for i := 0; i < matLen; i++ {
			mat[i][matLen-1] = <- worker.fromLeft
		}
	}
}

// Computes the outer cells of this worker submatrix, which are adjacent to other workers submatrices
// Returns the updated maxDiff value
func (worker worker) computeOuterCells(dst, src matrix.Matrix, prevMaxDiff float64) float64 {
	maxDiff, matLen := prevMaxDiff, worker.matDef.Size
	// TODO: This is probably not the best way to compute the outer cells in terms of performance
	for k := 1; k < matLen - 1; k++ {
		// Top outer cells
		dst[1][k] = 0.2*(src[1][k] + src[1][k-1] + src[1][k+1] + src[0][k] + src[2][k])
		maxDiff = MaxDiff(maxDiff, dst[1][k], src[1][k])
		// Bottom outer cells
		dst[matLen-2][k] = 0.2*(src[matLen-2][k] + src[matLen-2][k-1] + src[matLen-2][k+1] + src[matLen-3][k] + src[matLen-1][k])
		maxDiff = MaxDiff(maxDiff, dst[matLen-2][k], src[matLen-2][k])
		// Left outer cells
		dst[k][1] = 0.2*(src[k][1] + src[k-1][1] + src[k+1][1] + src[k][0] + src[k][2])
		maxDiff = MaxDiff(maxDiff, dst[k][1], src[k][1])
		// Right outer cells
		dst[k][matLen-2] = 0.2*(src[k][matLen-2] + src[k-1][matLen-2] + src[k+2][matLen-2] + src[k][matLen-3] + src[k][matLen-1])
		maxDiff = MaxDiff(maxDiff, dst[k][matLen-2], src[k][matLen-2])
	}

	return maxDiff
}

// Runs the jacobi method for the worker subproblem to get its partial result
func (worker worker) solveSubproblem(resMat matrix.Matrix, initialValue float64, maxIters int, tolerance float64, wg *sync.WaitGroup) {
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

		worker.sendOuterCells(matA)

		for i := 1; i < subprobSize-1; i++ {
			for j := 1; j < subprobSize-1; j++ {
				// Compute new value with 3x3 filter with no corners
				matB[i][j] = 0.2*(matA[i][j] + matA[i-1][j] + matA[i+1][j] + matA[i][j-1] + matA[i][j+1])
				maxDiff = MaxDiff(maxDiff, matA[i][j], matB[i][j])
			}
		}

		worker.recvAdjacentCells(matA)
		maxDiff = worker.computeOuterCells(matB, matA, maxDiff)
		// Actual max diff is maximum of all threads maxDiff
		maxDiff = worker.maxReduce(maxDiff)

		// Swap matrices
		matA, matB = matB, matA
	}

	worker.mergeSubproblem(resMat, matA)
}

// Parallel version of the jacobi method using Go routines
func RunJacobiPar(initialValue float64, nDim int, maxIters int, tolerance float64, nThreads int) matrix.Matrix {
	// Resulting matrix
	resMat := matrix.NewMatrix(0.0, nDim)

	var wg sync.WaitGroup
	wg.Add(nThreads)

	for i := 0; i < nThreads; i++ {
		// TODO: Initialize workers params
		go worker{

		}.solveSubproblem(resMat, initialValue, maxIters, tolerance, &wg)
	}

	wg.Wait()

	// TODO: Return number of iterations and maximum diff
	return resMat
}