package matrix

import (
	"fmt"
)

const (
	hot = 1.0
	cold = 0.0
)

type Matrix [][]float64

func CloneMatrix(mat Matrix) [][]float64 {
	length := len(mat)
	clone := make([][]float64, length, length)
	for i := range clone {
		clone[i] = make([]float64, length, length)
		copy(clone[i], mat[i])
	}
	
	return clone
}

func InitMatrix(n int, initialValue float64) Matrix {
	mat := make(Matrix, n, n)
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

func PrintMatrix(mat Matrix) {
	for _, row := range mat {
		for _, el := range row {
			fmt.Printf("%.4f ", el)
		}
		fmt.Println()
	}
}