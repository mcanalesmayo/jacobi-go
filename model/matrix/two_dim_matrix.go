package matrix

import (
	"fmt"
	"strings"
)

// TwoDimMatrix represents a matrix in a 2D array
type TwoDimMatrix []row

// row represents a 1D array belonging to a TwoDimMatrix
type row []float64

// NewTwoDimMatrix creates and initializes a 2D array representing a matrix
func NewTwoDimMatrix(initialValue float64, n int, topBoundary, bottomBoundary, leftBoundary, rightBoundary float64) TwoDimMatrix {
	mat := make(TwoDimMatrix, n, n)
	// Init inner cells value
	for i := range mat {
		// TODO: Look into how Go allocates the memory. Are rows contiguous? => Cache & Performance
		mat[i] = make([]float64, n, n)
		for j := range mat[i] {
			mat.SetCell(i, j, initialValue)
		}
	}

	// Init top, right and left boundaries
	for i := range mat {
		mat.SetCell(0, i, topBoundary)
		mat.SetCell(i, 0, leftBoundary)
		mat.SetCell(i, n-1, rightBoundary)
	}

	// Init bottom boundary
	for j := range mat {
		mat.SetCell(n-1, j, bottomBoundary)
	}

	return mat
}

// GetCell retrieves the value in the (i, j) position
func (mat TwoDimMatrix) GetCell(i, j int) float64 {
	return mat[i][j]
}

// SetCell updates the value in the (i, j) position
func (mat TwoDimMatrix) SetCell(i, j int, value float64) {
	mat[i][j] = value
}

// GetNDim retrieves the length of the matrix
func (mat TwoDimMatrix) GetNDim() int {
	return len(mat)
}

// Clone clones the portion of the matrix specified by a TwoDimMatrixDef
func (mat TwoDimMatrix) Clone(matDef MatrixDef) Matrix {
	x0, y0, x1, y1, length := matDef.Coords.X0, matDef.Coords.Y0, matDef.Coords.X1, matDef.Coords.Y1, matDef.Size

	clone := make(TwoDimMatrix, length)
	for i := x0; i <= x1; i++ {
		clone[i-x0] = make(row, length)
		for j := y0; j <= y1; j++ {
			clone.SetCell(i-x0, j-y0, mat.GetCell(i, j))
		}
	}

	return clone
}

// ToString returns the row in a human-readable format
func (row row) ToString() string {
	var resSb strings.Builder
	strBuf := make([]string, len(row))

	for k, el := range row {
		strBuf[k] = fmt.Sprintf("%.4f", el)
	}
	resSb.WriteString(strings.Join(strBuf, " "))

	return resSb.String()
}

// ToString returns the matrix in a human-readable format
func (mat TwoDimMatrix) ToString() string {
	var resSb strings.Builder
	strBuf := make([]string, len(mat))

	for k, row := range mat {
		strBuf[k] = row.ToString()
	}
	resSb.WriteString(strings.Join(strBuf, "\n"))

	return resSb.String()
}
