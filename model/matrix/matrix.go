package matrix

import (
	"fmt"
)

const (
	hot = 1.0
	cold = 0.0
)

type Matrix [][]float64
type Row []float64

func (rowA Row) CompareTo(rowB Row) bool {
	if rowA == nil && rowB == nil {
		return true
	} else if len(rowA) != len(rowB) {
		return false
	} else {
		for i := range rowA {
			if (rowA[i] != rowB[i]) {
				return false
			}
		}

		return true
	}
}

func (matA Matrix) CompareTo(matB Matrix) bool {
	if matA == nil && matB == nil {
		return true
	} else if len(matA) != len(matB) {
		return false
	} else {
		for i := range matA {
			// Need to assign to vars so that Row methods can be used
			var rowA, rowB Row = matA[i], matB[i]
			if (!rowA.CompareTo(rowB)) {
				return false
			}
		}

		return true
	}
}

func (mat Matrix) Clone() Matrix {
	length := len(mat)
	clone := make(Matrix, length, length)
	for i := range clone {
		clone[i] = make([]float64, length, length)
		copy(clone[i], mat[i])
	}
	
	return clone
}

func NewMatrix(initialValue float64, n int) Matrix {
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

func (mat Matrix) Print() {
	for _, row := range mat {
		for _, el := range row {
			fmt.Printf("%.4f ", el)
		}
		fmt.Println()
	}
}