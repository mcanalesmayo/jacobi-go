package main

import (
	"flag"
	"fmt"
	"time"
	"math"
	"github.com/mcanalesmayo/jacobi-go/model/matrix"
)

func RunJacobi(initialValue float64, nDim int, maxIters int, tolerance float64) (matrix.Matrix, int, float64) {
	var matA, matB matrix.Matrix
	// The algorithm requires computing each grid cell as a 3x3 filter with no corners
	// Therefore, we an aux matrix to keep the grid values in every iteration after computing new values
	matA = matrix.InitMatrix(nDim+2, initialValue)
	matB = matrix.CloneMatrix(matA)

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

	resMatrix, nIters, maxDiff := RunJacobi(*initialValuePtr, *nDimPtr, *maxIterationsPtr, *tolerancePtr)

	after := time.Now()

	fmt.Println("Results:")
	fmt.Println("Final grid:")
	matrix.PrintMatrix(resMatrix)
	fmt.Printf("Number of iterations: %d\n", nIters)
	fmt.Printf("Latest diff: %.4f\n", maxDiff)
	fmt.Printf("Running time: %s\n", after.Sub(before))
}