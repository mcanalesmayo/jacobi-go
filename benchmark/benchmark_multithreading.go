package benchmark

import (
	"fmt"
	"github.com/mcanalesmayo/jacobi-go"
	"github.com/mcanalesmayo/jacobi-go/model/matrix"
	"runtime"
	"testing"
)

type experiment struct {
	matrixType   matrix.MatrixType
	initialValue float64
	nDim         int
	maxIters     int
	tolerance    float64
	nThreads     int
}

// BenchmarkSingleVsMultithreading runs the single-threaded and multi-threaded versions for different matrix sizes.
func BenchmarkSingleVsMultithreading(b *testing.B) {
	experiments := []experiment{
		{matrix.TwoDimContiguousMatrixType, 0.5, 16, 1000, 1.0e-4, 1},
		{matrix.TwoDimContiguousMatrixType, 0.5, 64, 1000, 1.0e-4, 1},
		{matrix.TwoDimContiguousMatrixType, 0.5, 256, 1000, 1.0e-4, 1},
		{matrix.TwoDimContiguousMatrixType, 0.5, 1024, 1000, 1.0e-4, 1},
		{matrix.TwoDimContiguousMatrixType, 0.5, 4096, 1000, 1.0e-4, 1},
		{matrix.TwoDimContiguousMatrixType, 0.5, 16, 1000, 1.0e-4, 4},
		{matrix.TwoDimContiguousMatrixType, 0.5, 64, 1000, 1.0e-4, 4},
		{matrix.TwoDimContiguousMatrixType, 0.5, 256, 1000, 1.0e-4, 4},
		{matrix.TwoDimContiguousMatrixType, 0.5, 1024, 1000, 1.0e-4, 4},
		{matrix.TwoDimContiguousMatrixType, 0.5, 4096, 1000, 1.0e-4, 4},
	}

	fmt.Printf("Running with GOMAXPROCS=%d\n", runtime.GOMAXPROCS(0))

	for _, params := range experiments {
		b.Run(fmt.Sprintf("%.4f,%d,%d,%.4f,%d,%s", params.initialValue, params.nDim, params.maxIters, params.tolerance, params.nThreads, params.matrixType.ToString()), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				jacobi.RunJacobi(params.initialValue, params.nDim, params.maxIters, params.tolerance, params.nThreads, params.matrixType)
			}
		})
	}
}
