package matrix

import (
	"github.com/mcanalesmayo/jacobi-go/utils"
)

const (
	// Hot is the value of a hot state of a cell
	Hot = 1.0
	// Cold is the value a cold state of a cell
	Cold = 0.0
	// TwoDimMatrix is the code to represent a TwoDimMatrix
	TwoDimMatrixType = 0
	// OneDimMatrix is the code to represent a OneDimMatrix
	OneDimMatrixType = 1
)

type MatrixType int

func (matrixType MatrixType) ToString() string {
	if matrixType == TwoDimMatrixType {
		return "Two dimensions matrix"
	}
	return "One dimension matrix"
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

type Matrix interface {
	utils.Stringable
	MatrixCloneable
	// Get retrieves the value in the (i, j) position of the matrix
	GetCell(i, j int) float64
	// Set updates the value in the (i, j) position
	SetCell(i, j int, value float64)
	GetNDim() int
}

type MatrixCloneable interface {
	// Clone returns a Matrix
	Clone(matDef MatrixDef) interface{}
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