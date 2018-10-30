package main

import (
	"fmt"
	"time"
	"math"
)

const (
	hot = 1.0
	cold = 0.0
	initialGrid = 0.5
	maxIterations = 1000
	tol = 1.0e-4
	nDim = 16
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

func initMatrix(n int, fillVal float64) [][]float64 {
	mat := make([][]float64, n, n)
	for i := range mat {
		// TODO: Look into how Go allocates the memory. Are rows contiguous? => Cache & Performance
		mat[i] = make([]float64, n, n)
		for j := range mat[i] {
			mat[i][j] = fillVal
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

func main() {
	var matA, matB Matrix
	matA = initMatrix(nDim+2, initialGrid)
	matB = cloneMatrix(matA)

	matrixIters := nDim + 1

	iterations, maxDiff := 0, 1.0

	fmt.Printf("Running simulation with tolerance=%f and max iterations=%d\n", tol, maxIterations)

	before := time.Now()

	for maxDiff > tol && iterations < maxIterations {
		maxDiff = 0.0

		for i := 1; i < matrixIters; i++ {
			for j := 1; j < matrixIters; j++ {
				matB[i][j] = 0.2*(matA[i][j] + matA[i-1][j] + matA[i+1][j] + matA[i][j-1] + matA[i][j+1])
				absDiff := math.Abs(matB[i][j] - matA[i][j])
				if (absDiff > maxDiff) {
					maxDiff = absDiff
				}
			}
		}

		// Swap matrices
		matA, matB = matB, matA
		iterations += 1
	}

	after := time.Now()

	fmt.Println("Final grid:")
	printMatrix(matA)
	fmt.Println("Results:")
	fmt.Printf("Iterations: %d\n", iterations)
	fmt.Printf("Tolerance: %.4f\n", maxDiff)
	fmt.Printf("Running time: %s\n", after.Sub(before))
}