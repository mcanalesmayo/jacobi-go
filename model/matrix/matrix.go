package matrix

import (
	"github.com/mcanalesmayo/jacobi-go/utils"
)

const (
	// Hot is the value of a hot state of a cell
	Hot = 1.0
	// Cold is the value a cold state of a cell
	Cold = 0.0
	// TwoDimMatrixType is the code to represent a TwoDimMatrix
	TwoDimMatrixType = 0
	// OneDimMatrixType is the code to represent a OneDimMatrix
	OneDimMatrixType = 1
)

// MatrixType defines the underlying representation of a matrix
type MatrixType int

// ToString returns a string representation of a matrix
func (matrixType MatrixType) ToString() string {
	if matrixType == TwoDimMatrixType {
		return "Two dimensions matrix"
	}
	return "One dimension matrix"
}

// Matrix defines a matrix
type Matrix interface {
	utils.Stringable
	MatrixCloneable
	// GetCell retrieves the value in the (i, j) position
	GetCell(i, j int) float64
	// SetCell updates the value in the (i, j) position
	SetCell(i, j int, value float64)
	// GetNDim retrieves the length of the matrix
	GetNDim() int
}

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

// MatrixCloneable represents a matrix that can be cloned
type MatrixCloneable interface {
	// Clone returns a Matrix
	Clone(matDef MatrixDef) Matrix
}

// CompareMatrices returns true if both matrices contain equal cells or both are nil,
// otherwise returns false
func CompareMatrices(matA, matB Matrix) bool {
	matNDim := matA.GetNDim()

	if matNDim != matB.GetNDim() {
		return false
	} else {
		for i := 0; i < matNDim; i++ {
			for j := 0; j < matNDim; j++ {
				if !utils.CompareFloats(matA.GetCell(i, j), matB.GetCell(i, j), utils.Epsilon) {
					return false
				}
			}
		}

		return true
	}
}
