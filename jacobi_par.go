package jacobi

import (
	"math"
	"sync"
	"github.com/mcanalesmayo/jacobi-go/model/matrix"
)

type globalParams struct {
	nWorkers, size int
}

type adjacentChns struct {
	// For sharing values with threads working on adjacent submatrices
	toTop, toBottom, toRight, toLeft chan float64
	fromTop, fromBottom, fromRight, fromLeft chan float64
}

type worker struct {
	id int
	// Global problem parameters
	globalParams globalParams
	// Subproblem matrix
	matDef matrix.MatrixDef
	// For communicating with adjacent workers
	adjacent adjacentChns
	// For reducing maxDiff
	maxDiffRes []chan float64
}

// Creates the corresponding adjacentChns for each thread
func newAdjacentChns(nThreads, nDim, subprobSize int) []adjacentChns {
	res, nThreadsSqrt := make([]adjacentChns, nThreads), int(math.Sqrt(float64(nThreads)))

	for id := 0; id < nThreads; id++ {
		rowN := id % nThreadsSqrt
		columnN := int(id / nThreadsSqrt)

		if rowN == 0 {
			if columnN == 0 {
				// Worker for top-left corner matrix
				res[id] = adjacentChns{
					toTop: nil,
					fromTop: nil,
					toBottom: make(chan float64, subprobSize),
					fromBottom: make(chan float64, subprobSize),
					toRight: make(chan float64, subprobSize),
					fromRight: make(chan float64, subprobSize),
					toLeft: nil,
					fromLeft: nil,
				}
			} else if columnN == nThreadsSqrt - 1{
				// Worker for top-right corner matrix
				res[id] = adjacentChns{
					toTop: nil,
					fromTop: nil,
					toBottom: make(chan float64, subprobSize),
					fromBottom: make(chan float64, subprobSize),
					toRight: nil,
					fromRight: nil,
					toLeft: res[id-1].fromRight,
					fromLeft: res[id-1].toRight,
				}
			} else {
				// Worker for top matrix
				res[id] = adjacentChns{
					toTop: nil,
					fromTop: nil,
					toBottom: make(chan float64, subprobSize),
					fromBottom: make(chan float64, subprobSize),
					toRight: make(chan float64, subprobSize),
					fromRight: make(chan float64, subprobSize),
					toLeft: res[id-1].fromRight,
					fromLeft: res[id-1].toRight,
				}
			}
		} else if rowN == nThreadsSqrt - 1 {
			if columnN == 0 {
				// Worker for bottom-left corner matrix
				res[id] = adjacentChns{
					toTop: res[id-nThreadsSqrt].fromBottom,
					fromTop: res[id-nThreadsSqrt].toBottom,
					toBottom: nil,
					fromBottom: nil,
					toRight: make(chan float64, subprobSize),
					fromRight: make(chan float64, subprobSize),
					toLeft: nil,
					fromLeft: nil,
				}
			} else if columnN == nThreadsSqrt - 1{
				// Worker for bottom-right corner matrix
				res[id] = adjacentChns{
					toTop: res[id-nThreadsSqrt].fromBottom,
					fromTop: res[id-nThreadsSqrt].toBottom,
					toBottom: nil,
					fromBottom: nil,
					toRight: nil,
					fromRight: nil,
					toLeft: res[id-1].fromRight,
					fromLeft: res[id-1].toRight,
				}
			} else {
				// Worker for bottom matrix
				res[id] = adjacentChns{
					toTop: res[id-nThreadsSqrt].fromBottom,
					fromTop: res[id-nThreadsSqrt].toBottom,
					toBottom: nil,
					fromBottom: nil,
					toRight: make(chan float64, subprobSize),
					fromRight: make(chan float64, subprobSize),
					toLeft: res[id-1].fromRight,
					fromLeft: res[id-1].toRight,
				}
			}
		} else {
			if columnN == 0 {
				// Worker for a left side matrix
				res[id] = adjacentChns{
					toTop: res[id-nThreadsSqrt].fromBottom,
					fromTop: res[id-nThreadsSqrt].toBottom,
					toBottom: make(chan float64, subprobSize),
					fromBottom: make(chan float64, subprobSize),
					toRight: make(chan float64, subprobSize),
					fromRight: make(chan float64, subprobSize),
					toLeft: nil,
					fromLeft: nil,
				}
			} else if columnN == nThreadsSqrt - 1{
				// Worker for a right side matrix
				res[id] = adjacentChns{
					toTop: res[id-nThreadsSqrt].fromBottom,
					fromTop: res[id-nThreadsSqrt].toBottom,
					toBottom: make(chan float64, subprobSize),
					fromBottom: make(chan float64, subprobSize),
					toRight: nil,
					fromRight: nil,
					toLeft: res[id-1].fromRight,
					fromLeft: res[id-1].toRight,
				}
			} else {
				// Worker for any of the rest of the submatrices
				res[id] = adjacentChns{
					toTop: res[id-nThreadsSqrt].fromBottom,
					fromTop: res[id-nThreadsSqrt].toBottom,
					toBottom: make(chan float64, subprobSize),
					fromBottom: make(chan float64, subprobSize),
					toRight: make(chan float64, subprobSize),
					fromRight: make(chan float64, subprobSize),
					toLeft: res[id-1].fromRight,
					fromLeft: res[id-1].toRight,
				}
			}
		}
	}

	return res
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
		
		// Fan out the result to the rest of the workers
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
			worker.adjacent.toTop <- mat[0][j]
		}
	}
	if coords.Y1 != worker.globalParams.size {
		for j := 0; j < matLen; j++ {
			worker.adjacent.toBottom <- mat[matLen-1][j]
		}
	}
	if coords.X0 != 0 {
		for i := 0; i < matLen; i++ {
			worker.adjacent.toRight <- mat[i][0]
		}
	}
	if coords.X1 != worker.globalParams.size {
		for i := 0; i < matLen; i++ {
			worker.adjacent.toLeft <- mat[i][matLen-1]
		}
	}
}

// Gets the adjacent workers outer values
func (worker worker) recvAdjacentCells(mat matrix.Matrix) {
	coords := worker.matDef.Coords
	matLen := worker.matDef.Size

	if coords.Y0 != 0 {
		for j := 0; j < matLen; j++ {
			mat[0][j] = <- worker.adjacent.fromTop
		}
	}
	if coords.Y1 != worker.globalParams.size {
		for j := 0; j < matLen; j++ {
			mat[matLen-1][j] = <- worker.adjacent.fromBottom
		}
	}
	if coords.X0 != 0 {
		for i := 0; i < matLen; i++ {
			mat[i][0] = <- worker.adjacent.fromRight
		}
	}
	if coords.X1 != worker.globalParams.size {
		for i := 0; i < matLen; i++ {
			mat[i][matLen-1] = <- worker.adjacent.fromLeft
		}
	}
}

// Computes the outer cells of this worker submatrix, which are adjacent to other workers submatrices
// Returns the updated maxDiff value
func (worker worker) computeOuterCells(dst, src matrix.Matrix, prevMaxDiff float64) float64 {
	maxDiff, matLen := prevMaxDiff, worker.matDef.Size
	// TODO: This is probably not the best way to compute the outer cells in terms of performance
	for k := 1; k < matLen-1; k++ {
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
	maxDiff, matLen := 1.0, worker.matDef.Size
	// TODO: Setup boundaries for matrix.NewMatrix, depending on worker submatrix location (worker id)


	// The algorithm requires computing each grid cell as a 3x3 filter with no corners
	// Therefore, we need an aux matrix to keep the grid values in every iteration after computing new values
	matA = matrix.NewMatrix(initialValue, matLen, matrix.Hot, matrix.Cold, matrix.Hot, matrix.Hot)
	matB = matA.Clone()

	for nIters := 0; maxDiff > tolerance && nIters < maxIters; nIters++ {
		maxDiff = 0.0

		worker.sendOuterCells(matA)

		// Outer cells are a special case which will be computed later on
		for i := 1; i < matLen-1; i++ {
			for j := 1; j < matLen-1; j++ {
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
	// TODO: Check preconditions


	// Resulting matrix, init value doesn't matter at this point as workers will overwrite all cells
	resMat := matrix.NewMatrix(0.0, nDim+2, matrix.Hot, matrix.Cold, matrix.Hot, matrix.Hot)

	var wg sync.WaitGroup

	subprobSize := int(math.Sqrt(float64((nDim*nDim)/nThreads)))
	maxDiffResChns := make([]chan float64, nThreads)
	adjacentChns := newAdjacentChns(nThreads, nDim, subprobSize)
	for i := 0; i < nThreads; i++ {
		// These channels can also be unbuffered, as there's currently no computation between sending and receiving
		maxDiffResChns[i] = make(chan float64, 1)
	}

	wg.Add(nThreads)
	for id := 0; id < nThreads; id++ {
		// TODO: Initialize workers params
		go worker{
			id: id,
			globalParams: globalParams{
				nWorkers: nThreads,
				size: nDim,
			},
			matDef: matrix.MatrixDef{
				Size: subprobSize,
			},
			adjacent: adjacentChns[id],
			maxDiffRes: maxDiffResChns,
		}.solveSubproblem(resMat, initialValue, maxIters, tolerance, &wg)
	}
	wg.Wait()

	// TODO: Return number of iterations and maximum diff
	return resMat
}