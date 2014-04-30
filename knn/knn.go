package knn

import (
	"fmt"
	"math"

	base "github.com/sjwhitworth/golearn/base"
	util "github.com/sjwhitworth/golearn/utilities"

	mat "github.com/skelterjohn/go.matrix"
)

//A KNN Classifier. Consists of a data matrix, associated labels in the same order as the matrix, and a name.
type KNNClassifier struct {
	base.BaseClassifier
}

//Mints a new classifier.
func (KNN *KNNClassifier) New(name string, labels []string, numbers []float64, x int, y int) {

	//Write in some error handling here
	// if x != len(KNN.Labels) {
	// 	return errors.New("KNN: There must be a label for each row")
	// }

	KNN.Data = *mat.MakeDenseMatrix(numbers, x, y)
	KNN.Name = name
	KNN.Labels = labels
}

//Computes the Euclidean distance between two vectors.
func (KNN *KNNClassifier) ComputeDistance(vector *mat.DenseMatrix, testrow *mat.DenseMatrix) float64 {
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

//Returns a classification for the vector, based on a vector input, using the KNN algorithm.
func (KNN *KNNClassifier) Predict(vector *mat.DenseMatrix, K int) (string, []int) {

	rows := KNN.Data.Rows()
	rownumbers := make(map[int]float64)
	labels := make([]string, 0)
	maxmap := make(map[string]int)

	for i := 0; i < rows; i++ {
		row := KNN.Data.GetRowVector(i)
		eucdistance := KNN.ComputeDistance(row, vector)
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
