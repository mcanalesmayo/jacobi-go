package matrix

import (
	"fmt"
	"strings"
)

// OneDimMatrix represents a matrix in a 1D array
type OneDimMatrix struct {
	matrix []float64
	nDim   int
}

// NewOneDimMatrix creates and initializes a 2D array representing a matrix
func NewOneDimMatrix(initialValue float64, n int, topBoundary, bottomBoundary, leftBoundary, rightBoundary float64) interface{} {
	mat := OneDimMatrix{
		matrix: make([]float64, n*n),
		nDim:   n,
	}

	// Init inner cells value
	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			mat.SetCell(i, j, initialValue)
		}
	}

	// Init top, right and left boundaries
	for i := 0; i < n; i++ {
		mat.SetCell(0, i, topBoundary)
		mat.SetCell(i, 0, leftBoundary)
		mat.SetCell(i, n-1, rightBoundary)
	}

	// Init bottom boundary
	for j := 0; j < n; j++ {
		mat.SetCell(n-1, j, bottomBoundary)
	}

	return mat
}

func (mat OneDimMatrix) GetCell(i, j int) float64 {
	return mat.matrix[i*mat.nDim+j]
}

func (mat OneDimMatrix) SetCell(i, j int, value float64) {
	mat.matrix[i*mat.nDim+j] = value
}

func (mat OneDimMatrix) GetNDim() int {
	return mat.nDim
}

// Clone clones the portion of the matrix specified by a OneDimMatrixDef
func (mat OneDimMatrix) Clone(matDef MatrixDef) interface{} {
	x0, y0, x1, y1, length := matDef.Coords.X0, matDef.Coords.Y0, matDef.Coords.X1, matDef.Coords.Y1, matDef.Size

	clone := OneDimMatrix{
		nDim:   length,
		matrix: make([]float64, length*length),
	}
	for i := x0; i <= x1; i++ {
		for j := y0; j <= y1; j++ {
			clone.SetCell(i-x0, j-y0, mat.GetCell(i, j))
		}
	}

	return clone
}

// ToString returns the matrix in a human-readable format
func (mat OneDimMatrix) ToString() string {
	var resSb strings.Builder
	matStrBuf := make([]string, mat.nDim)
	rowStrBuf := make([]string, mat.nDim)

	for i := 0; i < mat.nDim; i++ {
		for j := 0; j < mat.nDim; j++ {
			rowStrBuf[j] = fmt.Sprintf("%.4f", mat.matrix[i*mat.nDim+j])
		}
		matStrBuf[i] = strings.Join(matStrBuf, " ")
	}
	resSb.WriteString(strings.Join(matStrBuf, "\n"))

	return resSb.String()
}
