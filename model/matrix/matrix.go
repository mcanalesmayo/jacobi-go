package matrix

import (
	"fmt"
)

const (
	Hot = 1.0
	Cold = 0.0
)

type Matrix [][]float64
type Row []float64
type MatrixDef struct {
	Coords struct {
		// Top-left corner and bottom-right corner
		X0, Y0, X1, Y1 int
	}
	// Precomputed matrix size: len(matrix)
	Size int
}

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

func NewMatrix(initialValue float64, n int, topBoundary, bottomBoundary, leftBoundary, rightBoundary float64) Matrix {
	mat := make(Matrix, n, n)
	// Init inner cells value
	for i := range mat {
		// TODO: Look into how Go allocates the memory. Are rows contiguous? => Cache & Performance
		mat[i] = make([]float64, n, n)
		for j := range mat[i] {
			mat[i][j] = initialValue
		}
	}

	// Init top, right and left boundaries
	for i := range mat {
		mat[0][i] = topBoundary
		mat[i][0] = leftBoundary
		mat[i][n-1] = rightBoundary
	}

	// Init bottom boundary
	for j := range mat {
		mat[n-1][j] = bottomBoundary
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