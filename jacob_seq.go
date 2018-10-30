package main

import (
	"flag"
	"fmt"
	"time"
	"math"
)

const (
	hot = 1.0
	cold = 0.0
)

type Matrix [][]float64

func cloneMatrix(mat Matrix) [][]float64 {
	length := len(mat)
	clone := make([][]float64, length, length)
	for i := range clone {
		clone[i] = make([]float64, length, length)
		copy(clone[i], mat[i])
	}
	
	return clone
}

func initMatrix(n int, initialValue float64) [][]float64 {
	mat := make([][]float64, n, n)
	// Init inner cells value
	for i := range mat {
		// TODO: Look into how Go allocates the memory. Are rows contiguous? => Cache & Performance
		mat[i] = make([]float64, n, n)
		for j := range mat[i] {
			mat[i][j] = initialValue
		}
	}

	// Init hot boundary
	for i := range mat {
		mat[i][0] = hot
		mat[i][n-1] = hot
		mat[0][i] = hot
	}

	// Init cold boundary
	for j := range mat {
		mat[n-1][j] = cold
	}

	return mat
}

func printMatrix(mat Matrix) {
	for _, row := range mat {
		for _, el := range row {
			fmt.Printf("%.4f ", el)
		}
		fmt.Println()
	}
}

func runJacobi(initialValue float64, nDim int, maxIters int, tolerance float64) (Matrix, int, float64) {
	var matA, matB Matrix
	// The algorithm requires computing each grid cell as a 3x3 filter with no corners
	// Therefore, we an aux matrix to keep the grid values in every iteration after computing new values
	matA = initMatrix(nDim+2, initialValue)
	matB = cloneMatrix(matA)

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

func main() {
	initialValuePtr := flag.Float64("init-value", 0.5, "initial value of the inner grid cells")
	nDimPtr := flag.Int("size", 16, "size of each side of the grid")
	maxIterationsPtr := flag.Int("max-iters", 1000, "maximum number of iterations")
	tolerancePtr := flag.Float64("tol", 1.0e-4, "difference tolerance")

	fmt.Printf("Running simulation with initial value=%.4f, num dims=%d, max iterations=%d and tolerance=%.4f\n",
		*initialValuePtr, *nDimPtr, *maxIterationsPtr, *tolerancePtr)

	before := time.Now()

	resMatrix, nIters, maxDiff := runJacobi(*initialValuePtr, *nDimPtr, *maxIterationsPtr, *tolerancePtr)

	after := time.Now()

	fmt.Println("Results:")
	fmt.Println("Final grid:")
	printMatrix(resMatrix)
	fmt.Printf("Number of iterations: %d\n", nIters)
	fmt.Printf("Latest diff: %.4f\n", maxDiff)
	fmt.Printf("Running time: %s\n", after.Sub(before))
}