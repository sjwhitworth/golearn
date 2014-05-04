//@todo: A lot of code duplication here.

package knn

import (
	"github.com/gonum/matrix/mat64"
	"github.com/sjwhitworth/golearn/base"
	pairwiseMetrics "github.com/sjwhitworth/golearn/metrics/pairwise"
	util "github.com/sjwhitworth/golearn/utilities"
)

//A KNN Regressor. Consists of a data matrix, associated result variables in the same order as the matrix, and a name.
type KNNRegressor struct {
	base.BaseEstimator
	Values       []float64
	DistanceFunc string
}

// Mints a new classifier.
func NewKnnRegressor(values []float64, numbers []float64, x int, y int, distfunc string) *KNNRegressor {
	KNN := KNNRegressor{}
	KNN.Data = mat64.NewDense(x, y, numbers)
	KNN.Values = values
	KNN.DistanceFunc = distfunc
	return &KNN
}

//Returns an average of the K nearest labels/variables, based on a vector input.
func (KNN *KNNRegressor) Predict(vector *mat64.Dense, K int) float64 {

	// Get the number of rows
	rows, _ := KNN.Data.Dims()
	rownumbers := make(map[int]float64)
	labels := make([]float64, 0)

	// Check what distance function we are using
	switch KNN.DistanceFunc {
	case "euclidean":
		{
			euclidean := pairwiseMetrics.NewEuclidean()
			for i := 0; i < rows; i++ {
				row := KNN.Data.RowView(i)
				rowMat := util.FloatsToMatrix(row)
				distance := euclidean.Distance(rowMat, vector)
				rownumbers[i] = distance
			}
		}
	case "manhattan":
		{
			manhattan := pairwiseMetrics.NewEuclidean()
			for i := 0; i < rows; i++ {
				row := KNN.Data.RowView(i)
				rowMat := util.FloatsToMatrix(row)
				distance := manhattan.Distance(rowMat, vector)
				rownumbers[i] = distance
			}
		}
	}

	sorted := util.SortIntMap(rownumbers)
	values := sorted[:K]

	var sum float64
	for _, elem := range values {
		value := KNN.Values[elem]
		labels = append(labels, value)
		sum += value
	}

	average := sum / float64(K)
	return average
}
