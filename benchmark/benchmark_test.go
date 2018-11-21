package benchmark

import (
	"fmt"
	"github.com/mcanalesmayo/jacobi-go"
	"runtime"
	"testing"
)

func BenchmarkRunJacobi(b *testing.B) {
	experiments := []struct {
		initialValue float64
		nDim         int
		maxIters     int
		tolerance    float64
		nThreads     int
	}{
		{0.5, 16, 1000, 1.0e-4, 1},
		{0.5, 64, 1000, 1.0e-4, 1},
		{0.5, 256, 1000, 1.0e-4, 1},
		{0.5, 1024, 1000, 1.0e-4, 1},
		{0.5, 4098, 1000, 1.0e-4, 1},
		{0.5, 16, 1000, 1.0e-4, 4},
		{0.5, 64, 1000, 1.0e-4, 4},
		{0.5, 256, 1000, 1.0e-4, 4},
		{0.5, 1024, 1000, 1.0e-4, 4},
		{0.5, 4098, 1000, 1.0e-4, 4},
	}

	fmt.Printf("Running with GOMAXPROCS=%d\n", runtime.GOMAXPROCS(0))

	for _, params := range experiments {
		b.Run(fmt.Sprintf("%.4f,%d,%d,%.4f,%d", params.initialValue, params.nDim, params.maxIters, params.tolerance, params.nThreads), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				jacobi.RunJacobi(params.initialValue, params.nDim, params.maxIters, params.tolerance, params.nThreads)
			}
		})
	}
}
