package knn

import (
	base "github.com/sjwhitworth/golearn/base"
	pairwiseMetrics "github.com/sjwhitworth/golearn/metrics/pairwise"
	util "github.com/sjwhitworth/golearn/utilities"
	mat "github.com/skelterjohn/go.matrix"
)

//A KNN Classifier. Consists of a data matrix, associated labels in the same order as the matrix, and a name.
type KNNClassifier struct {
	base.BaseEstimator
	Labels       []string
	DistanceFunc string
}

//Mints a new classifier.
func (KNN *KNNClassifier) New(labels []string, numbers []float64, x int, y int, distfunc string) {

	KNN.Data = mat.MakeDenseMatrix(numbers, x, y)
	KNN.Labels = labels
	KNN.DistanceFunc = distfunc
}

// Returns a classification for the vector, based on a vector input, using the KNN algorithm.
// @todo: Lots of room to improve this. V messy.
func (KNN *KNNClassifier) Predict(vector *mat.DenseMatrix, K int) (string, []int) {

	rows := KNN.Data.Rows()
	rownumbers := make(map[int]float64)
	labels := make([]string, 0)
	maxmap := make(map[string]int)

	for i := 0; i < rows; i++ {
		row := KNN.Data.GetRowVector(i)

		//Will put code in to check errs later
		euclidean := pairwiseMetrics.NewEuclidean()
		eucdistance, _ := euclidean.Distance(row, vector)
		rownumbers[i] = eucdistance
	}

	sorted := util.SortIntMap(rownumbers)
	values := sorted[:K]

	for _, elem := range values {
		labels = append(labels, KNN.Labels[elem])

		if _, ok := maxmap[KNN.Labels[elem]]; ok {
			maxmap[KNN.Labels[elem]] += 1
		} else {
			maxmap[KNN.Labels[elem]] = 1
		}
	}

	sortedlabels := util.SortStringMap(maxmap)
	label := sortedlabels[0]

	return label, values
}
