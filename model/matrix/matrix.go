package matrix

import (
	"fmt"
)

const (
	// Predefined values representing hot and cold state of a cell
	Hot = 1.0
	Cold = 0.0
)

type Matrix [][]float64
type Row []float64
type Coords struct {
	// Top-left corner and bottom-right corner
	X0, Y0, X1, Y1 int
}
type MatrixDef struct {
	Coords Coords
	// Precomputed matrix size: len(matrix)
	Size int
}

// Returns true if both rows contain equal values or both are nil
// Returns false otherwise
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

// Returns true if both matrices contain equal cells or both are nil
// Returns false otherwise
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

// Clones the portion of the matrix specified by the matDef argument
func (mat Matrix) Clone(matDef MatrixDef) Matrix {
	x0, y0, x1, y1, length := matDef.Coords.X0, matDef.Coords.Y0, matDef.Coords.X1, matDef.Coords.Y1, matDef.Size
	clone := make(Matrix, length)

	for i := x0; i <= x1; i++ {
		clone[i] = make([]float64, length)
		for j := y0; j <= y1; j++ {
			clone[i-x0][j-y0] = mat[i][j]
		}
	}
	
	return clone
}

// Initializes a new matrix with the specified values
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

// Prints the matrix in a human-readable format
func (mat Matrix) Print() {
	for _, row := range mat {
		for _, el := range row {
			fmt.Printf("%.4f ", el)
		}
		fmt.Println()
	}
}