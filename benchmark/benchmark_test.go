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

// BenchmarkMatrixTypes runs the simulation with different matrix implementations to see the difference in terms of performance.
// Current implementations:
// - OneDimMatrixType, which allocates the complete matrix at once and accesses are performed with one index ([i*len+j]).
// - TwoDimDividedMatrixType, which allocates the matrix row by row. This may lead the matrix to be divided in memory, causing performance degradation because of higher number of cache misses
// - TwoDimContiguousMatrixType, which allocates the matrix at once (as OneDimMatrixType), but accesses are performed with two indexes ([i][j]).
func BenchmarkMatrixTypes(b *testing.B) {
	// Interleaving of multiple threads may favor the TwoDimDividedMatrixType matrix to be divided in memory, as one thread's matrix allocation may be interleaved with
	// another thread's activity which requires memory allocation too
	experiments := []experiment{
		{matrix.TwoDimDividedMatrixType, 0.5, 2048, 1000, 1.0e-4, 4},
		{matrix.TwoDimContiguousMatrixType, 0.5, 2048, 1000, 1.0e-4, 4},
		{matrix.OneDimMatrixType, 0.5, 2048, 1000, 1.0e-4, 4},
	}

	for _, params := range experiments {
		b.Run(fmt.Sprintf("%.4f,%d,%d,%.4f,%d,%s", params.initialValue, params.nDim, params.maxIters, params.tolerance, params.nThreads, params.matrixType.ToString()), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				jacobi.RunJacobi(params.initialValue, params.nDim, params.maxIters, params.tolerance, params.nThreads, params.matrixType)
			}
		})
	}
}

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
				jacobi.RunJacobi(params.initialValue, params.nDim, params.maxIters, params.tolerance, params.nThreads, params.matrixType, )
			}
		})
	}
}
