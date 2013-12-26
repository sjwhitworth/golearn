package main

import (
		mat "github.com/skelterjohn/go.matrix"
		"fmt"
		)

type KNNClassifier struct {
	Data mat.DenseMatrix
	Name string
	Labels []string
}

//Initialises a new classifier
func (KNN *KNNClassifier) New(name string, labels []string, numbers []float64, x int, y int){
	KNN.Data = *mat.MakeDenseMatrix(numbers, x, y)
	KNN.Name = name
	KNN.Labels = labels
}

func (KNN *KNNClassifier) ComputeDistance(vector mat.DenseMatrix) mat.DenseMatrix {
	//Add switches for different distance metrics
	result, _ := KNN.Data.TimesDense(&vector)
	return *result
}

func (KNN *KNNClassifier) Predict(vector mat.DenseMatrix) mat.DenseMatrix {
	blah := KNN.ComputeDistance(vector)
	//return *mat.Difference(&KNN.Data, &vector)
	return blah
}

func (KNN *KNNClassifier) GetLabel(index int) string {
	return KNN.Labels[index]
}

func main(){
	knn := KNNClassifier{}
	dense := *mat.MakeDenseMatrix([]float64{4,5,1,3,4,2},2,3)
	knn.New("Testing", []string{"this sucks", "hiya"}, []float64{1,2,3,4,5,6},2,3)
	//hey := knn.ComputeDistance(dense)
	blof := knn.Predict(dense)
	fmt.Println(blof)
}