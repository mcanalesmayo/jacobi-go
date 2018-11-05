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
func newAdjacentChns(nThreads, subprobSize int) []adjacentChns {
	res, nThreadsSqrt := make([]adjacentChns, nThreads), int(math.Sqrt(float64(nThreads)))

	for id := 0; id < nThreads; id++ {
		rowN, columnN := int(id / nThreadsSqrt), id % nThreadsSqrt

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
	for i := coords.X0; i <= coords.X1; i++ {
		for j := coords.Y0; j <= coords.Y1; j++ {
			// Values are ordered by the sender
			resMat[i+1][j+1] = subprobResMat[i+1][j+1]
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
		for i := 1; i < worker.globalParams.nWorkers; i++ {
			maxMaxDiff = MaxMaxDiff(maxMaxDiff, <- worker.maxDiffRes[i])
		}
		
		// Fan out the result to the rest of the workers
		for i := 1; i < worker.globalParams.nWorkers; i++ {
			worker.maxDiffRes[i] <- maxMaxDiff
		}
	} else {
		// 'Non-root' workers send their results
		worker.maxDiffRes[worker.id] <- maxDiff
		// Wait for result calculated by 'Root' worker
		maxMaxDiff = <- worker.maxDiffRes[worker.id]
	}

	return maxMaxDiff
}

// Sends the worker outer values to adjacent workers
func (worker worker) sendOuterCells(mat matrix.Matrix) {
	matLen, nThreadsSqrt := worker.matDef.Size, int(math.Sqrt(float64(worker.globalParams.nWorkers)))

	// Since subproblem coordinates never change, this solution
	// isn't the best one in terms of performance, as these
	// checks are done for every jacobi iteration
	rowN, columnN := int(worker.id / nThreadsSqrt), worker.id % nThreadsSqrt

	if rowN != 0 {
		for j := 1; j < matLen; j++ {
			worker.adjacent.toTop <- mat[1][j]
		}
	}
	if rowN != nThreadsSqrt-1 {
		for j := 1; j < matLen; j++ {
			worker.adjacent.toBottom <- mat[matLen-1][j]
		}
	}
	if columnN != 0 {
		for i := 1; i < matLen; i++ {
			worker.adjacent.toLeft <- mat[i][1]
		}
	}
	if columnN != nThreadsSqrt-1 {
		for i := 1; i < matLen; i++ {
			worker.adjacent.toRight <- mat[i][matLen-1]
		}
	}
}

// Gets the adjacent workers outer values
func (worker worker) recvAdjacentCells(mat matrix.Matrix) {
	matLen, nThreadsSqrt := worker.matDef.Size, int(math.Sqrt(float64(worker.globalParams.nWorkers)))
	rowN, columnN := int(worker.id / nThreadsSqrt), worker.id % nThreadsSqrt

	if rowN != 0 {
		for j := 1; j < matLen; j++ {
			mat[0][j] = <- worker.adjacent.fromTop
		}
	}
	if rowN != nThreadsSqrt-1 {
		for j := 1; j < matLen; j++ {
			mat[matLen][j] = <- worker.adjacent.fromBottom
		}
	}
	if columnN != 0 {
		for i := 1; i < matLen; i++ {
			mat[i][0] = <- worker.adjacent.fromLeft
		}
	}
	if columnN != nThreadsSqrt-1 {
		for i := 1; i < matLen; i++ {
			mat[i][matLen] = <- worker.adjacent.fromRight
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
		maxDiff = MaxMaxDiff(maxDiff, math.Abs(dst[1][k] - src[1][k]))
		// Bottom outer cells
		dst[matLen-2][k] = 0.2*(src[matLen-2][k] + src[matLen-2][k-1] + src[matLen-2][k+1] + src[matLen-3][k] + src[matLen-1][k])
		maxDiff = MaxMaxDiff(maxDiff, math.Abs(dst[matLen-2][k] - src[matLen-2][k]))
		// Left outer cells
		dst[k][1] = 0.2*(src[k][1] + src[k-1][1] + src[k+1][1] + src[k][0] + src[k][2])
		maxDiff = MaxMaxDiff(maxDiff, math.Abs(dst[k][1] - src[k][1]))
		// Right outer cells
		dst[k][matLen-2] = 0.2*(src[k][matLen-2] + src[k-1][matLen-2] + src[k+2][matLen-2] + src[k][matLen-3] + src[k][matLen-1])
		maxDiff = MaxMaxDiff(maxDiff, math.Abs(dst[k][matLen-2] - src[k][matLen-2]))
	}

	return maxDiff
}

// Runs the jacobi method for the worker subproblem to get its partial result
func (worker worker) solveSubproblem(resMat matrix.Matrix, initialValue float64, maxIters int, tolerance float64, wg *sync.WaitGroup) {
	defer wg.Done()

	var matA, matB matrix.Matrix
	maxDiff, matDef, matLen := 1.0, worker.matDef, worker.matDef.Size
	// Adjacent cells are needed to compute outer cells
	cloneMatDef := matrix.MatrixDef{
		matrix.Coords{matDef.Coords.X0, matDef.Coords.Y0, matDef.Coords.X1+1, matDef.Coords.Y1+1},
		matLen+2,
	}

	// The algorithm requires computing each grid cell as a 3x3 filter with no corners
	// Therefore, we need an aux matrix to keep the grid values in every iteration after computing new values
	matA, matB = resMat.Clone(cloneMatDef), resMat.Clone(cloneMatDef)

	for nIters := 0; maxDiff > tolerance && nIters < maxIters; nIters++ {
		maxDiff = 0.0

		worker.sendOuterCells(matA)

		// Outer cells are a special case which will be computed later on
		for i := 2; i < matLen-1; i++ {
			for j := 2; j < matLen-1; j++ {
				// Compute new value with 3x3 filter with no corners
				matB[i][j] = 0.2*(matA[i][j] + matA[i-1][j] + matA[i+1][j] + matA[i][j-1] + matA[i][j+1])
				maxDiff = MaxMaxDiff(maxDiff, math.Abs(matA[i][j] - matB[i][j]))
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
	resMat, maxDiffResChns := matrix.NewMatrix(initialValue, nDim+2, matrix.Hot, matrix.Cold, matrix.Hot, matrix.Hot), make([]chan float64, nThreads)
	for i := 0; i < nThreads; i++ {
		// These channels can also be unbuffered, as there's currently no computation between sending and receiving
		maxDiffResChns[i] = make(chan float64, 1)
	}
	subprobSize, nThreadsSqrt := int(math.Sqrt(float64((nDim*nDim)/nThreads))), int(math.Sqrt(float64(nThreads)))
	workerMatLen, adjacentChns := nDim/nThreadsSqrt, newAdjacentChns(nThreads, subprobSize)

	var wg sync.WaitGroup
	wg.Add(nThreads)
	for id := 0; id < nThreads; id++ {
		x0, y0 := (id/nThreadsSqrt)*workerMatLen, (id%nThreadsSqrt)*workerMatLen
		x1, y1 := x0 + workerMatLen - 1, y0 + workerMatLen - 1

		go worker{
			id: id,
			globalParams: globalParams{
				nWorkers: nThreads,
				size: nDim,
			},
			matDef: matrix.MatrixDef{
				Coords: matrix.Coords{x0, y0, x1, y1},
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