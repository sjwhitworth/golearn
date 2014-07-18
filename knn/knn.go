// Package knn implements a K Nearest Neighbors object, capable of both classification
// and regression. It accepts data in the form of a slice of float64s, which are then reshaped
// into a X by Y matrix.
package knn

import (
	"github.com/gonum/matrix/mat64"
	base "github.com/sjwhitworth/golearn/base"
	pairwiseMetrics "github.com/sjwhitworth/golearn/metrics/pairwise"
	util "github.com/sjwhitworth/golearn/utilities"
)

// A KNNClassifier consists of a data matrix, associated labels in the same order as the matrix, and a distance function.
// The accepted distance functions at this time are 'euclidean' and 'manhattan'.
type KNNClassifier struct {
	base.BaseEstimator
	TrainingData      *base.Instances
	DistanceFunc      string
	NearestNeighbours int
}

// NewKnnClassifier returns a new classifier
func NewKnnClassifier(distfunc string, neighbours int) *KNNClassifier {
	KNN := KNNClassifier{}
	KNN.DistanceFunc = distfunc
	KNN.NearestNeighbours = neighbours
	return &KNN
}

// Fit stores the training data for later
func (KNN *KNNClassifier) Fit(trainingData *base.Instances) {
	KNN.TrainingData = trainingData
}

// PredictOne returns a classification for the vector, based on a vector input, using the KNN algorithm.
// See http://en.wikipedia.org/wiki/K-nearest_neighbors_algorithm.
func (KNN *KNNClassifier) PredictOne(vector []float64) string {

	rows := KNN.TrainingData.Rows
	rownumbers := make(map[int]float64)
	labels := make([]string, 0)
	maxmap := make(map[string]int)

	convertedVector := util.FloatsToMatrix(vector)

	// Check what distance function we are using
	switch KNN.DistanceFunc {
	case "euclidean":
		{
			euclidean := pairwiseMetrics.NewEuclidean()
			for i := 0; i < rows; i++ {
				row := KNN.TrainingData.GetRowVectorWithoutClass(i)
				rowMat := util.FloatsToMatrix(row)
				distance := euclidean.Distance(rowMat, convertedVector)
				rownumbers[i] = distance
			}
		}
	case "manhattan":
		{
			manhattan := pairwiseMetrics.NewEuclidean()
			for i := 0; i < rows; i++ {
				row := KNN.TrainingData.GetRowVectorWithoutClass(i)
				rowMat := util.FloatsToMatrix(row)
				distance := manhattan.Distance(rowMat, convertedVector)
				rownumbers[i] = distance
			}
		}
	}

	sorted := util.SortIntMap(rownumbers)
	values := sorted[:KNN.NearestNeighbours]

	for _, elem := range values {
		label := KNN.TrainingData.GetClass(elem)
		labels = append(labels, label)

		if _, ok := maxmap[label]; ok {
			maxmap[label]++
		} else {
			maxmap[label] = 1
		}
	}

	sortedlabels := util.SortStringMap(maxmap)
	label := sortedlabels[0]

	return label
}

func (KNN *KNNClassifier) Predict(what *base.Instances) *base.Instances {
	ret := what.GeneratePredictionVector()
	for i := 0; i < what.Rows; i++ {
		ret.SetAttrStr(i, 0, KNN.PredictOne(what.GetRowVectorWithoutClass(i)))
	}
	return ret
}

// A KNNRegressor consists of a data matrix, associated result variables in the same order as the matrix, and a name.
type KNNRegressor struct {
	base.BaseEstimator
	Values       []float64
	DistanceFunc string
}

// NewKnnRegressor mints a new classifier.
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
