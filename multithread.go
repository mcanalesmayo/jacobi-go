package jacobi

import (
	"github.com/mcanalesmayo/jacobi-go/model/matrix"
	"math"
	"os"
	"sync"
)

const (
	invalidProblemParams = 1
)

type globalParams struct {
	nWorkers, size int
}

type adjacents struct {
	// For sharing values among adjacent workers
	toTopWorker, toBottomWorker, toRightWorker, toLeftWorker         chan float64
	fromTopWorker, fromBottomWorker, fromRightWorker, fromLeftWorker chan float64
	topValues, bottomValues, rightValues, leftValues                 []float64
}

type worker struct {
	// For identifying the worker
	id, rowNumber, columnNumber int
	// Global problem parameters
	globalParams globalParams
	// Subproblem matrix
	matDef matrix.MatrixDef
	// For communicating with adjacent workers
	adjacents adjacents
	// For reducing maxDiff
	maxDiffResToRoot, maxDiffResFromRoot []chan float64
}

// Creates the corresponding adjacents for each thread
func newAdjacents(nThreads, subprobSize int) []adjacents {
	res, nThreadsSqrt := make([]adjacents, nThreads), int(math.Sqrt(float64(nThreads)))

	for id := 0; id < nThreads; id++ {
		rowN, columnN := int(id/nThreadsSqrt), id%nThreadsSqrt

		if rowN == 0 {
			if columnN == 0 {
				// Worker for top-left corner matrix
				res[id] = adjacents{
					toTopWorker:      nil,
					fromTopWorker:    nil,
					toBottomWorker:   make(chan float64, subprobSize),
					fromBottomWorker: make(chan float64, subprobSize),
					toRightWorker:    make(chan float64, subprobSize),
					fromRightWorker:  make(chan float64, subprobSize),
					toLeftWorker:     nil,
					fromLeftWorker:   nil,
				}
			} else if columnN == nThreadsSqrt-1 {
				// Worker for top-right corner matrix
				res[id] = adjacents{
					toTopWorker:      nil,
					fromTopWorker:    nil,
					toBottomWorker:   make(chan float64, subprobSize),
					fromBottomWorker: make(chan float64, subprobSize),
					toRightWorker:    nil,
					fromRightWorker:  nil,
					toLeftWorker:     res[id-1].fromRightWorker,
					fromLeftWorker:   res[id-1].toRightWorker,
				}
			} else {
				// Worker for top matrix
				res[id] = adjacents{
					toTopWorker:      nil,
					fromTopWorker:    nil,
					toBottomWorker:   make(chan float64, subprobSize),
					fromBottomWorker: make(chan float64, subprobSize),
					toRightWorker:    make(chan float64, subprobSize),
					fromRightWorker:  make(chan float64, subprobSize),
					toLeftWorker:     res[id-1].fromRightWorker,
					fromLeftWorker:   res[id-1].toRightWorker,
				}
			}
		} else if rowN == nThreadsSqrt-1 {
			if columnN == 0 {
				// Worker for bottom-left corner matrix
				res[id] = adjacents{
					toTopWorker:      res[id-nThreadsSqrt].fromBottomWorker,
					fromTopWorker:    res[id-nThreadsSqrt].toBottomWorker,
					toBottomWorker:   nil,
					fromBottomWorker: nil,
					toRightWorker:    make(chan float64, subprobSize),
					fromRightWorker:  make(chan float64, subprobSize),
					toLeftWorker:     nil,
					fromLeftWorker:   nil,
				}
			} else if columnN == nThreadsSqrt-1 {
				// Worker for bottom-right corner matrix
				res[id] = adjacents{
					toTopWorker:      res[id-nThreadsSqrt].fromBottomWorker,
					fromTopWorker:    res[id-nThreadsSqrt].toBottomWorker,
					toBottomWorker:   nil,
					fromBottomWorker: nil,
					toRightWorker:    nil,
					fromRightWorker:  nil,
					toLeftWorker:     res[id-1].fromRightWorker,
					fromLeftWorker:   res[id-1].toRightWorker,
				}
			} else {
				// Worker for bottom matrix
				res[id] = adjacents{
					toTopWorker:      res[id-nThreadsSqrt].fromBottomWorker,
					fromTopWorker:    res[id-nThreadsSqrt].toBottomWorker,
					toBottomWorker:   nil,
					fromBottomWorker: nil,
					toRightWorker:    make(chan float64, subprobSize),
					fromRightWorker:  make(chan float64, subprobSize),
					toLeftWorker:     res[id-1].fromRightWorker,
					fromLeftWorker:   res[id-1].toRightWorker,
				}
			}
		} else {
			if columnN == 0 {
				// Worker for a left side matrix
				res[id] = adjacents{
					toTopWorker:      res[id-nThreadsSqrt].fromBottomWorker,
					fromTopWorker:    res[id-nThreadsSqrt].toBottomWorker,
					toBottomWorker:   make(chan float64, subprobSize),
					fromBottomWorker: make(chan float64, subprobSize),
					toRightWorker:    make(chan float64, subprobSize),
					fromRightWorker:  make(chan float64, subprobSize),
					toLeftWorker:     nil,
					fromLeftWorker:   nil,
				}
			} else if columnN == nThreadsSqrt-1 {
				// Worker for a right side matrix
				res[id] = adjacents{
					toTopWorker:      res[id-nThreadsSqrt].fromBottomWorker,
					fromTopWorker:    res[id-nThreadsSqrt].toBottomWorker,
					toBottomWorker:   make(chan float64, subprobSize),
					fromBottomWorker: make(chan float64, subprobSize),
					toRightWorker:    nil,
					fromRightWorker:  nil,
					toLeftWorker:     res[id-1].fromRightWorker,
					fromLeftWorker:   res[id-1].toRightWorker,
				}
			} else {
				// Worker for any of the rest of the submatrices
				res[id] = adjacents{
					toTopWorker:      res[id-nThreadsSqrt].fromBottomWorker,
					fromTopWorker:    res[id-nThreadsSqrt].toBottomWorker,
					toBottomWorker:   make(chan float64, subprobSize),
					fromBottomWorker: make(chan float64, subprobSize),
					toRightWorker:    make(chan float64, subprobSize),
					fromRightWorker:  make(chan float64, subprobSize),
					toLeftWorker:     res[id-1].fromRightWorker,
					fromLeftWorker:   res[id-1].toRightWorker,
				}
			}
		}

		res[id].topValues = make([]float64, subprobSize)
		res[id].bottomValues = make([]float64, subprobSize)
		res[id].leftValues = make([]float64, subprobSize)
		res[id].rightValues = make([]float64, subprobSize)
	}

	return res
}

// Merges the worker subproblem resulting matrix into the global resulting matrix
func (worker worker) mergeSubproblem(resMat, subprobResMat matrix.Matrix) {
	coords := worker.matDef.Coords
	x0, y0, x1, y1 := coords.X0, coords.Y0, coords.X1, coords.Y1

	for i := x0; i <= x1; i++ {
		for j := y0; j <= y1; j++ {
			// Values are ordered by the sender
			resMat.SetCell(i, j, subprobResMat.GetCell(i-x0, j-y0))
		}
	}
}

// Computes the new maxDiff taking into account subproblem matrix as well as other workers matrix (like a max-reduce on the global matrix)
func (worker worker) computeNewMaxDiff(matB, matA matrix.Matrix) float64 {
	matLen, maxDiff := worker.matDef.Size, 0.0

	// My subproblem maxDiff
	for i := 0; i < matLen; i++ {
		for j := 0; j < matLen; j++ {
			maxDiff = math.Max(maxDiff, math.Abs(matB.GetCell(i, j)-matA.GetCell(i, j)))
		}
	}

	return worker.maxReduce(maxDiff)
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
		for i := 0; i < worker.globalParams.nWorkers-1; i++ {
			maxMaxDiff = math.Max(maxMaxDiff, <-worker.maxDiffResToRoot[i])
		}

		// Fan out the result to the rest of the workers
		for i := 0; i < worker.globalParams.nWorkers-1; i++ {
			worker.maxDiffResFromRoot[i] <- maxMaxDiff
		}
	} else {
		// 'Non-root' workers send their results
		worker.maxDiffResToRoot[worker.id-1] <- maxDiff
		// Wait for result calculated by 'Root' worker
		maxMaxDiff = <-worker.maxDiffResFromRoot[worker.id-1]
	}

	return maxMaxDiff
}

// Sends the worker outer values to adjacent workers
func (worker worker) sendOuterCells(mat matrix.Matrix) {
	matLen, nThreadsSqrt := worker.matDef.Size, int(math.Sqrt(float64(worker.globalParams.nWorkers)))

	// Since subproblem coordinates never change, this solution
	// isn't the best one in terms of performance, as these
	// checks are done for every jacobi iteration
	if worker.rowNumber != 0 {
		for j := 0; j < matLen; j++ {
			worker.adjacents.toTopWorker <- mat.GetCell(0, j)
		}
	}
	if worker.rowNumber != nThreadsSqrt-1 {
		for j := 0; j < matLen; j++ {
			worker.adjacents.toBottomWorker <- mat.GetCell(matLen-1, j)
		}
	}
	if worker.columnNumber != 0 {
		for i := 0; i < matLen; i++ {
			worker.adjacents.toLeftWorker <- mat.GetCell(i, 0)
		}
	}
	if worker.columnNumber != nThreadsSqrt-1 {
		for i := 0; i < matLen; i++ {
			worker.adjacents.toRightWorker <- mat.GetCell(i, matLen-1)
		}
	}
}

// Gets the adjacent workers outer values
func (worker worker) recvAdjacentCells(mat matrix.Matrix) {
	matLen, nThreadsSqrt := worker.matDef.Size, int(math.Sqrt(float64(worker.globalParams.nWorkers)))

	if worker.rowNumber != 0 {
		for j := 0; j < matLen; j++ {
			worker.adjacents.topValues[j] = <-worker.adjacents.fromTopWorker
		}
	}
	if worker.rowNumber != nThreadsSqrt-1 {
		for j := 0; j < matLen; j++ {
			worker.adjacents.bottomValues[j] = <-worker.adjacents.fromBottomWorker
		}
	}
	if worker.columnNumber != 0 {
		for i := 0; i < matLen; i++ {
			worker.adjacents.leftValues[i] = <-worker.adjacents.fromLeftWorker
		}
	}
	if worker.columnNumber != nThreadsSqrt-1 {
		for i := 0; i < matLen; i++ {
			worker.adjacents.rightValues[i] = <-worker.adjacents.fromRightWorker
		}
	}
}

// Computes the outer cells of this worker submatrix, which are adjacent to other workers submatrices
// Returns the updated maxDiff value
func (worker worker) computeOuterCells(dst, src matrix.Matrix) {
	matLen := worker.matDef.Size

	// Outer cells in the corners are a special case
	// Top-left corner
	dst.SetCell(0, 0, 0.2 * (src.GetCell(0, 0) + worker.adjacents.leftValues[0] + src.GetCell(0, 1) + worker.adjacents.topValues[0] + src.GetCell(1, 0)))
	// Top-right corner
	dst.SetCell(0, matLen-1, 0.2 * (src.GetCell(0, matLen-1) + src.GetCell(0, matLen-2) + worker.adjacents.rightValues[0] + worker.adjacents.topValues[matLen-1] + src.GetCell(1, matLen-1)))
	// Bottom-left corner
	dst.SetCell(matLen-1, 0, 0.2 * (src.GetCell(matLen-1, 0) + worker.adjacents.leftValues[matLen-1] + src.GetCell(matLen-1, 1) + src.GetCell(matLen-2, 0) + worker.adjacents.bottomValues[0]))
	// Bottom-right corner
	dst.SetCell(matLen-1, matLen-1, 0.2 * (src.GetCell(matLen-1, matLen-1) + src.GetCell(matLen-1, matLen-2) + worker.adjacents.rightValues[matLen-1] + src.GetCell(matLen-2, matLen-1) + worker.adjacents.bottomValues[matLen-1]))

	// Rest of outer cells
	// TODO: This is probably not the best way to compute the outer cells in terms of performance
	for k := 1; k < matLen-1; k++ {
		// Top outer cells
		dst.SetCell(0, k, 0.2 * (src.GetCell(0, k) + src.GetCell(0, k-1) + src.GetCell(0, k+1) + worker.adjacents.topValues[k] + src.GetCell(1, k)))
		// Bottom outer cells
		dst.SetCell(matLen-1, k, 0.2 * (src.GetCell(matLen-1, k) + src.GetCell(matLen-1, k-1) + src.GetCell(matLen-1, k+1) + src.GetCell(matLen-2, k) + worker.adjacents.bottomValues[k]))
		// Left outer cells
		dst.SetCell(k, 0, 0.2 * (src.GetCell(k, 0) + worker.adjacents.leftValues[k] + src.GetCell(k, 1) + src.GetCell(k-1, 0) + src.GetCell(k+1, 0)))
		// Right outer cells
		dst.SetCell(k, matLen-1, 0.2 * (src.GetCell(k, matLen-1) + src.GetCell(k, matLen-2) + worker.adjacents.rightValues[k] + src.GetCell(k-1, matLen-1) + src.GetCell(k+1, matLen-1)))
	}
}

func (worker worker) setupBoundaries(initialValue, topBoundary, bottomBoundary, leftBoundary, rightBoundary float64) {
	matLen, nThreadsSqrt := worker.matDef.Size, int(math.Sqrt(float64(worker.globalParams.nWorkers)))

	// By default adjacent cell will have the initial value
	for k := 0; k < matLen; k++ {
		worker.adjacents.topValues[k] = initialValue
		worker.adjacents.bottomValues[k] = initialValue
		worker.adjacents.leftValues[k] = initialValue
		worker.adjacents.rightValues[k] = initialValue
	}

	// Overwrite adjacent cells in special cases
	if worker.rowNumber == 0 {
		for j := 0; j < matLen; j++ {
			worker.adjacents.topValues[j] = topBoundary
		}
	}
	if worker.rowNumber == nThreadsSqrt-1 {
		for j := 0; j < matLen; j++ {
			worker.adjacents.bottomValues[j] = bottomBoundary
		}
	}
	if worker.columnNumber == 0 {
		for i := 0; i < matLen; i++ {
			worker.adjacents.leftValues[i] = leftBoundary
		}
	}
	if worker.columnNumber == nThreadsSqrt-1 {
		for i := 0; i < matLen; i++ {
			worker.adjacents.rightValues[i] = rightBoundary
		}
	}
}

// Runs the jacobi method for the worker subproblem to get its partial result
func (worker worker) solveSubproblem(resMat matrix.Matrix, initialValue float64, maxIters int, tolerance float64, wg *sync.WaitGroup) {
	defer wg.Done()

	maxDiff, matDef, matLen := math.MaxFloat64, worker.matDef, worker.matDef.Size

	// The algorithm requires computing each grid cell as a 3x3 filter with no corners
	// Therefore, we need an aux matrix to keep the grid values in every iteration after computing new values
	matA, matB := resMat.Clone(matDef).(matrix.Matrix), resMat.Clone(matDef).(matrix.Matrix)

	worker.setupBoundaries(initialValue, matrix.Hot, matrix.Cold, matrix.Hot, matrix.Hot)

	for nIters := 0; maxDiff > tolerance && nIters < maxIters; nIters++ {
		worker.sendOuterCells(matA)

		// Outer cells are a special case which will be computed later on
		for i := 1; i < matLen-1; i++ {
			for j := 1; j < matLen-1; j++ {
				// Compute new value with 3x3 filter with no corners
				matB.SetCell(i, j, 0.2 * (matA.GetCell(i, j) + matA.GetCell(i-1, j) + matA.GetCell(i+1, j) + matA.GetCell(i, j-1) + matA.GetCell(i, j+1)))
			}
		}

		worker.recvAdjacentCells(matA)
		worker.computeOuterCells(matB, matA)
		// Actual max diff is maximum of all threads maxDiff
		maxDiff = worker.computeNewMaxDiff(matB, matA)

		// Swap matrices
		matA, matB = matB, matA
	}

	worker.mergeSubproblem(resMat, matA)
}

func validatePreconditions(nDim, nThreads int) bool {
	if nThreadsSqrt := int(math.Sqrt(float64(nThreads))); nThreadsSqrt*nThreadsSqrt == nThreads && nDim%nThreads == 0 {
		return true
	}
	return false
}

// runMultithreadedJacobi runs a multi-threaded version of the jacobi method using Go routines
func runMultithreadedJacobi(matrixType matrix.MatrixType, initialValue float64, nDim int, maxIters int, tolerance float64, nThreads int) (matrix.Matrix, int, float64) {
	if !validatePreconditions(nDim, nThreads) {
		os.Exit(invalidProblemParams)
	}

	var resMat matrix.Matrix
	if matrixType == matrix.TwoDimMatrixType {
		resMat = matrix.NewTwoDimMatrix(initialValue, nDim+2, matrix.Hot, matrix.Cold, matrix.Hot, matrix.Hot).(matrix.Matrix)
	} else {
		resMat = matrix.NewOneDimMatrix(initialValue, nDim+2, matrix.Hot, matrix.Cold, matrix.Hot, matrix.Hot).(matrix.Matrix)
	}

	maxDiffResToRoot, maxDiffResFromRoot := make([]chan float64, nThreads), make([]chan float64, nThreads)
	for i := 0; i < nThreads-1; i++ {
		// These channels can also be unbuffered, as there's currently no computation between sending and receiving
		maxDiffResToRoot[i] = make(chan float64, 1)
		maxDiffResFromRoot[i] = make(chan float64, 1)
	}
	subprobSize, nThreadsSqrt := int(math.Sqrt(float64(nDim*nDim/nThreads))), int(math.Sqrt(float64(nThreads)))
	workerMatLen, adjacents := nDim/nThreadsSqrt, newAdjacents(nThreads, subprobSize)

	var wg sync.WaitGroup
	wg.Add(nThreads)
	for id := 0; id < nThreads; id++ {
		x0, y0 := id/nThreadsSqrt*workerMatLen+1, id%nThreadsSqrt*workerMatLen+1
		x1, y1 := x0+workerMatLen-1, y0+workerMatLen-1

		go worker{
			id:           id,
			rowNumber:    int(id / nThreadsSqrt),
			columnNumber: id % nThreadsSqrt,
			globalParams: globalParams{
				nWorkers: nThreads,
				size:     nDim,
			},
			matDef: matrix.MatrixDef{
				Coords: matrix.Coords{x0, y0, x1, y1},
				Size:   subprobSize,
			},
			adjacents:          adjacents[id],
			maxDiffResToRoot:   maxDiffResToRoot,
			maxDiffResFromRoot: maxDiffResFromRoot,
		}.solveSubproblem(resMat, initialValue, maxIters, tolerance, &wg)
	}
	wg.Wait()

	// TODO: Return number of iterations and maximum diff
	return resMat, 0, 0.0
}
