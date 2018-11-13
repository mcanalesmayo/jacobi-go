package benchmark

import (
	"fmt"
	"github.com/mcanalesmayo/jacobi-go"
	"testing"
)

func BenchmarkRunSinglethreadedJacobi(b *testing.B) {
	experiments := []struct {
		initialValue float64
		nDim         int
		maxIters     int
		tolerance    float64
	}{
		{0.5, 16, 1000, 1.0e-4},
		{0.5, 128, 1000, 1.0e-4},
		{0.5, 1024, 1000, 1.0e-4},
	}

	for _, params := range experiments {
		b.Run(fmt.Sprintf("sequential,%.4f,%d,%d,%.4f", params.initialValue, params.nDim, params.maxIters, params.tolerance), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				jacobi.RunSinglethreadedJacobi(params.initialValue, params.nDim, params.maxIters, params.tolerance)
			}
		})
	}
}