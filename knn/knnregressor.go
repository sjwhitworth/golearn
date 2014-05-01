package knn

import (
	"fmt"
	"math"

	util "github.com/sjwhitworth/golearn/utilities"
	mat "github.com/skelterjohn/go.matrix"
)

//A KNN Regressor. Consists of a data matrix, associated result variables in the same order as the matrix, and a name.
type KNNRegressor struct {
	Data   *mat.DenseMatrix
	Name   string
	Labels []float64
}

//Mints a new classifier.
func (KNN *KNNRegressor) New(name string, labels []float64, numbers []float64, x int, y int) {

	KNN.Data = mat.MakeDenseMatrix(numbers, x, y)
	KNN.Name = name
	KNN.Labels = labels
}

//Computes the Euclidean distance between two vectors.
func (KNN *KNNRegressor) ComputeDistance(vector *mat.DenseMatrix, testrow *mat.DenseMatrix) float64 {
	var sum float64

	difference, err := testrow.MinusDense(vector)
	flat := difference.Array()

	if err != nil {
		fmt.Println(err)
	}

	for _, i := range flat {
		squared := math.Pow(i, 2)
		sum += squared
	}

	eucdistance := math.Sqrt(sum)
	return eucdistance
}

//Returns an average of the K nearest labels/variables, based on a vector input.
func (KNN *KNNRegressor) Predict(vector *mat.DenseMatrix, K int) (float64, []int) {

	rows := KNN.Data.Rows()
	rownumbers := make(map[int]float64)
	labels := make([]float64, 1)
	sum := 0.0

	for i := 0; i < rows; i++ {
		row := KNN.Data.GetRowVector(i)
		eucdistance := KNN.ComputeDistance(row, vector)
		rownumbers[i] = eucdistance
	}

	sorted := util.SortIntMap(rownumbers)
	values := sorted[:K]

	for _, elem := range values {
		value := KNN.Labels[elem]
		labels = append(labels, value)
		sum += value
	}

	average := sum / float64(K)
	return average, values
}
