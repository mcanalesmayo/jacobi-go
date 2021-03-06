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
  Input: initialValue, sideLength, maxIters, tolerance
  Output: matrix

  nIters <- 0
  maxDiff <- MAX_NUM
  prevIterMatrix <- initMatrix(initialValue, sideLength)
  matrix <- initEmptyMatrix(sideLength)

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

In the current multithreaded implementation, each worker is assigned a square submatrix of the problem of `sideLength/sqrt(nRoutines)` rows x `sideLength/sqrt(nRoutines)` columns. Hence, there's a precondition by which the number of threads needs to be a perfect square root and the size of one side of the matrix must be divisible by the number of threads. Each worker has to share its outer cells values with its adjacent workers, like shown in the animation below, corresponding to an example of a 16x16 matrix solved by 16 workers:

![Example of a 16x16 matrix solved by 16 workers](doc/img/examples/workers_submatrices.gif)

This is the best solution when it comes to minimizing the communication (I/O) overhead, as the number of values which are needed to be shared among all workers for each iteration is `4*sideLength*(sqrt(nRoutines)-1)`, hence an upper boundary of `O(sideLength*sqrt(nRoutines))`. Another possible implementation would be slicing the matrix in small submatrices of `sideLength/nRoutines` rows x `sideLength` columns, but that one would require `sideLength*(nRoutines-1)` values to be shared. The upper boundary in this case is `O(sideLength*nRoutines)`, which is worse than the current implementation.

On the other hand, the current implementation has a clear drawback, which is the precondition. To solve the problem with a 4-cores CPU the current implementation would be the best approach, but what if it's needed to be solved with a 8-cores CPU? Only 4 of them could be used because of the preconditions. That's why, in some cases (depending the underlying hardware resources), that other possible implementation may fit better, even though it's slightly more I/O bound.

## Run and analyze benchmarks
By using the built-in tools we can easily run the benchmark and take a look at some hardware metrics to analyze the performance of the application. As prerequisite for visualizing the metrics, GraphViz must be installed.

To run the benchmark:
```
cd jacobi-go
# For CPU and memory profiles
go test -v -cpuprofile=cpuprof.out -memprofile=memprof.out -bench=. benchmark/benchmark_multithreading.go
# For traces
go test -v -trace=trace.out -bench=. benchmark/benchmark_multithreading.go
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

## Results
The following results were gathered by running the benchmark in a MacBook Pro with an i7-7820HQ (4 cores):

| Routines | Initial Value | Each side length | Max iters | Max diff | ns/op        | Speedup over single routine version |
|----------|---------------|------------------|-----------|----------|--------------|-------------------------------------|
| 1        | 0.5           | 16               | 1000      | 1.0e-3   | 3801827      | N/A                                 |
| 1        | 0.5           | 64               | 1000      | 1.0e-3   | 223633618    | N/A                                 |
| 1        | 0.5           | 256              | 1000      | 1.0e-3   | 3410725471   | N/A                                 |
| 1        | 0.5           | 1024             | 1000      | 1.0e-3   | 55885176608  | N/A                                 |
| 1        | 0.5           | 4096             | 1000      | 1.0e-3   | 897391537991 | N/A                                 |
| 4        | 0.5           | 16               | 1000      | 1.0e-3   | 7995262      | 0,4755                              |
| 4        | 0.5           | 64               | 1000      | 1.0e-3   | 102926714    | 2,1728                              |
| 4        | 0.5           | 256              | 1000      | 1.0e-3   | 1085940009   | 3,1408                              |
| 4        | 0.5           | 1024             | 1000      | 1.0e-3   | 16899117229  | 3,3070                              |
| 4        | 0.5           | 4096             | 1000      | 1.0e-3   | 255977729080 | 3,5057                              |

The below routines are a small sample after the benchmark converged into solutions with similar resulting time for a 16x16 matrix solved by 4 routines:

![Overview for 16x16 matrix solved by 4 Go routines](doc/img/traces/16_4_goroutine.png)

Same for 64x64:

![Overview for 64x64 matrix solved by 4 Go routines](doc/img/traces/64_4_goroutine.png)

For 256x256:

![Overview for 256x256 matrix solved by 4 Go routines](doc/img/traces/256_4_goroutine.png)

For 1024x1024:

![Overview for 1024x1024 matrix solved by 4 Go routines](doc/img/traces/1024_4_goroutine.png)

For 4096x4096:

![Overview for 4096x4096 matrix solved by 4 Go routines](doc/img/traces/4096_4_goroutine.png)

Taking a deeper look into the traces, we can see that for the 16x16 problem it's not worth to run several routines, as the execution time is almost hidden by the send<->receive synchronization. This is because for each iteration all workers need to share the outer cells of their subproblem and also reduce the max value of the diff. Keep in mind that the scale of the following image is in µs and there are a lot of holes, which can be considered as synchronization points:

![Traces for 16x16 matrix solved by 4 Go routines](doc/img/traces/16_4_traces.png)

Same for 4096x4096, notice the scale of the traces is bigger (ms) and there are way less synchronization points. The execution time in this case hides the synchronization time:

![Traces for 4096x4096 matrix solved by 4 Go routines](doc/img/traces/4096_4_traces.png)

It makes sense that for small matrix sizes it's not worth to use multithreading capabilities as there's a significant communication/synchronization overhead. As shown in the table, the speedup of the 4 routines version vs the single routine version for a 16x16 matrix was `0,4755`, which means it's around 2 times slower. With the 64x64 matrix we got a better result, 2 times faster. Progressively increasing the matrix size of the experiments will hide more and more the synchronization time and will end up reaching values close to 4, which is the number of cores of the CPU used in these experiments. This is because the workload that can be run in parallel (independent workload) will be close to 100% of the program workload.

The main workloads involved in the problem are:
* Computing the inner cells values of the worker matrix (its assigned subproblem). The time of this workload depends on `(subSideLength-1)*(subSideLength-1)`, which is `O(subSideLength^2)`, being `subSideLength` the length of each side of a routine submatrix.
* Reducing the `maxDiff` values, which depends on the number of routines (not the size of the matrix). We can consider it to be `O(nRoutines)`, as the reduce is done by a single routine which collects the values from all other routines and computes the max value. Hence, this time will get hidden by the time to compute the inner cells values.
* Sharing the outer cells values to adjacent workers. This will also get hidden, as the number of outer cells values is `4*sideLength` in the worst case, which makes this workload grow by `O(sideLength)`.

In terms of space, the total memory used by the program grows by `2*(sideLength^2)`, because the main routine uses a `sideLength` x `sideLength` matrix and the sum of all other routines submatrices size is also `sideLength` x `sideLength`. Therefore, the space usage grows by `O(sideLength^2)`, which is the same as the single-threaded version one.