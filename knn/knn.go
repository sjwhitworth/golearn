// Package KNN implements a K Nearest Neighbors object, capable of both classification
// and regression. It accepts data in the form of a slice of float64s, which are then reshaped
// into a X by Y matrix.
package knn

import (
	"github.com/gonum/matrix/mat64"
	base "github.com/sjwhitworth/golearn/base"
	pairwiseMetrics "github.com/sjwhitworth/golearn/metrics/pairwise"
	util "github.com/sjwhitworth/golearn/utilities"
)

// A KNN Classifier. Consists of a data matrix, associated labels in the same order as the matrix, and a distance function.
// The accepted distance functions at this time are 'euclidean' and 'manhattan'.
type KNNClassifier struct {
	base.BaseEstimator
	Labels       []string
	DistanceFunc string
}

// Returns a new classifier
func NewKnnClassifier(distfunc string) *KNNClassifier {
	KNN := KNNClassifier{}
	KNN.DistanceFunc = distfunc
	return &KNN
}

func (KNN *KNNClassifier) Fit(labels []string, numbers []float64, rows int, cols int) {
	if rows != len(labels) {
		panic(mat64.ErrShape)
	}

	KNN.Data = mat64.NewDense(rows, cols, numbers)
	KNN.Labels = labels
}

// Returns a classification for the vector, based on a vector input, using the KNN algorithm.
// See http://en.wikipedia.org/wiki/K-nearest_neighbors_algorithm.
func (KNN *KNNClassifier) Predict(vector []float64, K int) string {

	convertedVector := util.FloatsToMatrix(vector)
	// Get the number of rows
	rows, _ := KNN.Data.Dims()
	rownumbers := make(map[int]float64)
	labels := make([]string, 0)
	maxmap := make(map[string]int)

	// Check what distance function we are using
	switch KNN.DistanceFunc {
	case "euclidean":
		{
			euclidean := pairwiseMetrics.NewEuclidean()
			for i := 0; i < rows; i++ {
				row := KNN.Data.RowView(i)
				rowMat := util.FloatsToMatrix(row)
				distance := euclidean.Distance(rowMat, convertedVector)
				rownumbers[i] = distance
			}
		}
	case "manhattan":
		{
			manhattan := pairwiseMetrics.NewEuclidean()
			for i := 0; i < rows; i++ {
				row := KNN.Data.RowView(i)
				rowMat := util.FloatsToMatrix(row)
				distance := manhattan.Distance(rowMat, convertedVector)
				rownumbers[i] = distance
			}
		}
	}

	sorted := util.SortIntMap(rownumbers)
	values := sorted[:K]

	for _, elem := range values {
		// It's when we access this map
		labels = append(labels, KNN.Labels[elem])

		if _, ok := maxmap[KNN.Labels[elem]]; ok {
			maxmap[KNN.Labels[elem]] += 1
		} else {
			maxmap[KNN.Labels[elem]] = 1
		}
	}

	sortedlabels := util.SortStringMap(maxmap)
	label := sortedlabels[0]

	return label
}

//A KNN Regressor. Consists of a data matrix, associated result variables in the same order as the matrix, and a name.
type KNNRegressor struct {
	base.BaseEstimator
	Values       []float64
	DistanceFunc string
}

// Mints a new classifier.
func NewKnnRegressor(distfunc string) *KNNRegressor {
	KNN := KNNRegressor{}
	KNN.DistanceFunc = distfunc
	return &KNN
}

func (KNN *KNNRegressor) Fit(values []float64, numbers []float64, rows int, cols int) {
	if rows != len(values) {
		panic(mat64.ErrShape)
	}

	KNN.Data = mat64.NewDense(rows, cols, numbers)
	KNN.Values = values
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
