# jacobi-go
[![CircleCI](https://circleci.com/gh/mcanalesmayo/jacobi-go.svg?style=svg)](https://circleci.com/gh/mcanalesmayo/jacobi-go)
[![Go Report Card](https://goreportcard.com/badge/github.com/mcanalesmayo/jacobi-go)](https://goreportcard.com/report/github.com/mcanalesmayo/jacobi-go)
[![Godoc](https://img.shields.io/badge/go-documentation-blue.svg)](https://godoc.org/github.com/mcanalesmayo/jacobi-go)
[![MIT Licensed](https://img.shields.io/badge/license-MIT-blue.svg)](https://raw.githubusercontent.com/mcanalesmayo/jacobi-go/master/LICENSE)

## Description
Go implementation of a simulation of thermal transmission in a 2D space.

The purpose of this project is to compare the performance of a single-threaded implementation with a multithreaded one. Additionally, they can be compared with a single-threaded, multithreaded and distributed C implementation available in [this repo](https://github.com/mcanalesmayo/jacobi-mpi).

The simulation algorithm is really simple:
```
Algorithm thermalTransmission is:
  Input: initialValue, numDimensions, maxIters, tolerance
  Output: matrix

  nIters <- 0
  maxDiff <- MAX_NUM
  prevIterMatrix <- initMatrix(initialValue, numDimensions)
  matrix <- initEmptyMatrix(numDimensions)

  while maxDiff > tolerance AND nIters < maxIters do
    for each (i,j) in prevIterMatrix do
      matrix[i,j] <- arithmeticMean(prevIterMatrix[i,j],
        prevIterMatrix[i-1,j], prevIterMatrix[i+1,j]
        prevIterMatrix[i,j-1], prevIterMatrix[i,j+1])
    end

    maxDiff <- maxReduce(absoluteValue(prevIterMatrix-matrix))
    nIters++
    prevIterMatrix <- matrix
  end
```

## Run and analyze benchmarks
By using the built-in tools we can easily run the benchmark and take a look at some hardware metrics to analyze the performance of the application. As prerequisite for visualizing the metrics, GraphViz must be installed.

To run the benchmark:
```
cd jacobi-go
# For CPU and memory profiles
go test -v -cpuprofile=cpuprof.out -memprofile=memprof.out -bench=. benchmark/benchmark_test.go
# For traces
go test -v -trace=trace.out -bench=. benchmark/benchmark_test.go
```

To visualize the cpu metrics (same thing works for memory metrics) in PNG format or via web browser:
```
go tool pprof -png cpuprof.out
go tool pprof -http=localhost:8080 cpuprof.out
```

For traces another built-in tool has to be used, which allows to visualize the metrics via web browser:
```
go tool trace -http=localhost:8080 trace.out
```