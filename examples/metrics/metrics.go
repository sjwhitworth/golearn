package main

import (
	"fmt"

	pariwiseMetrics "github.com/sjwhitworth/golearn/metrics/pairwise"
	mat "github.com/skelterjohn/go.matrix"
)

func main() {
	vectorX := mat.MakeDenseMatrix([]float64{1, 2, 3}, 3, 1)
	vectorY := mat.MakeDenseMatrix([]float64{3, 4, 5}, 3, 1)

	euclidean := pariwiseMetrics.NewEuclidean()
	polyKernel := pariwiseMetrics.NewPolyKernel(3)

	euclideanDistance, err := euclidean.Distance(vectorX, vectorY)
	polyKernelDistance, err := polyKernel.Distance(vectorX, vectorY)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("Vector X:")
	fmt.Println(vectorX.String())
	fmt.Println("Vector Y: ")
	fmt.Println(vectorY.String())
	fmt.Println("Euclidean           : ", euclideanDistance)
	fmt.Println("PolyKernel(degree 3): ", polyKernelDistance)
}
