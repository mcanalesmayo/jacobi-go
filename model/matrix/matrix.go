package matrix

import (
	"fmt"
)

const (
	// Hot is the value of a hot state of a cell
	Hot = 1.0
	// Cold is the value a cold state of a cell
	Cold = 0.0
)

// Matrix represents a 2D array
type Matrix [][]float64

// Row represents a 1D array belonging to a Matrix
type Row []float64

// Coords defines a 2D square
type Coords struct {
	// Top-left corner and bottom-right corner
	X0, Y0, X1, Y1 int
}

// MatrixDef defines a submatrix inside a Matrix
type MatrixDef struct {
	Coords Coords
	// Precomputed matrix size: len(matrix)
	Size int
}

// CompareTo returns true if both rows contain equal values or both are nil,
// otherwise returns false
func (row Row) CompareTo(anotherRow Row) bool {
	if row == nil && anotherRow == nil {
		return true
	} else if len(row) != len(anotherRow) {
		return false
	} else {
		for i := range row {
			if row[i] != anotherRow[i] {
				return false
			}
		}

		return true
	}
}

// CompareTo returns true if both matrices contain equal cells or both are nil,
// otherwise returns false
func (mat Matrix) CompareTo(anotherMat Matrix) bool {
	if mat == nil && anotherMat == nil {
		return true
	} else if len(mat) != len(anotherMat) {
		return false
	} else {
		for i := range mat {
			// Need to assign to vars so that Row methods can be used
			var rowA, rowB Row = mat[i], anotherMat[i]
			if !rowA.CompareTo(rowB) {
				return false
			}
		}

		return true
	}
}

// Clone clones the portion of the matrix specified by a MatrixDef
func (mat Matrix) Clone(matDef MatrixDef) Matrix {
	x0, y0, x1, y1, length := matDef.Coords.X0, matDef.Coords.Y0, matDef.Coords.X1, matDef.Coords.Y1, matDef.Size

	clone := make(Matrix, length)
	for i := x0; i <= x1; i++ {
		clone[i-x0] = make([]float64, length)
		for j := y0; j <= y1; j++ {
			//fmt.Printf("(%d,%d)-(%d,%d); clone[%d][%d] = mat[%d][%d]\n", x0, y0, x1, y1, i-x0, j-y0, i, j)
			clone[i-x0][j-y0] = mat[i][j]
		}
	}

	return clone
}

// NewMatrix creates a new matrix with the specified initializing values
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

// Print prints the matrix in a human-readable format
func (mat Matrix) Print() {
	for _, row := range mat {
		for _, el := range row {
			fmt.Printf("%.4f ", el)
		}
		fmt.Println()
	}
}
