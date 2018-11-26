package matrix

import (
	"fmt"
	"github.com/mcanalesmayo/jacobi-go/utils"
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

func (mat TwoDimMatrix) Get(i, j int) float64 {
	return mat[i][j]
}

func (mat TwoDimMatrix) Set(i, j int, value float64) {
	mat[i][j] = value
}

// CompareTo returns true if both rows contain equal values or both are nil,
// otherwise returns false
func (row row) CompareTo(anotherRow row) bool {
	if row == nil && anotherRow == nil {
		return true
	} else if len(row) != len(anotherRow) {
		return false
	} else {
		for i := range row {
			if !utils.CompareFloats(row[i], anotherRow[i], utils.Epsilon) {
				return false
			}
		}

		return true
	}
}

// CompareTo returns true if both matrices contain equal cells or both are nil,
// otherwise returns false
func (mat TwoDimMatrix) CompareTo(anotherMat TwoDimMatrix) bool {
	if mat == nil && anotherMat == nil {
		return true
	} else if len(mat) != len(anotherMat) {
		return false
	} else {
		for i := range mat {
			if !mat[i].CompareTo(anotherMat[i]) {
				return false
			}
		}

		return true
	}
}

// Clone clones the portion of the matrix specified by a TwoDimMatrixDef
func (mat TwoDimMatrix) Clone(matDef MatrixDef) Matrix {
	x0, y0, x1, y1, length := matDef.Coords.X0, matDef.Coords.Y0, matDef.Coords.X1, matDef.Coords.Y1, matDef.Size

	clone := make(TwoDimMatrix, length)
	for i := x0; i <= x1; i++ {
		clone[i-x0] = make(row, length)
		for j := y0; j <= y1; j++ {
			clone[i-x0][j-y0] = mat[i][j]
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
