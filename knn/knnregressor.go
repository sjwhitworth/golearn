//@todo: A lot of code duplication here.

package knn

import (
	util "github.com/sjwhitworth/golearn/utilities"
	mat "github.com/skelterjohn/go.matrix"
)

//A KNN Regressor. Consists of a data matrix, associated result variables in the same order as the matrix, and a name.
type KNNRegressor struct {
	Data         *mat.DenseMatrix
	Labels       []float64
	DistanceFunc string
}

//Mints a new classifier.
func NewKnnRegressor(labels []float64, numbers []float64, x int, y int, distfunc string) *KNNRegressor {
	KNN := KNNRegressor{}
	KNN.Data = mat.MakeDenseMatrix(numbers, x, y)
	KNN.Labels = labels
	return &KNN
}

//Returns an average of the K nearest labels/variables, based on a vector input.
func (KNN *KNNRegressor) Predict(vector *mat.DenseMatrix, K int) (float64, []int) {

	rows := KNN.Data.Rows()
	rownumbers := make(map[int]float64)
	labels := make([]float64, 1)
	sum := 0.0

	for i := 0; i < rows; i++ {
		row := KNN.Data.GetRowVector(i)
		eucdistance, _ := util.ComputeDistance(KNN.DistanceFunc, row, vector)
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
