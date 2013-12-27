package main

import (
		mat "github.com/skelterjohn/go.matrix"
		rand "math/rand"
		"fmt"
		)

type KNNClassifier struct {
	Data mat.DenseMatrix
	Name string
	Labels []string
}

func RandomArray(n int) []float64 {
	ReturnedArray := make([]float64, n)
	for i := 0; i < n; i++ {
		ReturnedArray[i] = rand.Float64()
	}
	return ReturnedArray
}

//Mints a new classifier
func (KNN *KNNClassifier) New(name string, labels []string, numbers []float64, x int, y int){
	KNN.Data = *mat.MakeDenseMatrix(numbers, x, y)
	KNN.Name = name
	KNN.Labels = labels
}

//Computes a variety of distance metrics between two vectors
func (KNN *KNNClassifier) ComputeDistance(vector *mat.DenseMatrix) *mat.DenseMatrix {
	//Add switches for different distance metrics
	result, err := KNN.Data.TimesDense(vector)

	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(result)
	return result
}

//Returns a classification based on a vector input
func (KNN *KNNClassifier) Predict(vector mat.DenseMatrix) *mat.DenseMatrix {
	return KNN.ComputeDistance(&vector)
}

//Returns a label, given an index
func (KNN *KNNClassifier) GetLabel(index int) string {
	return KNN.Labels[index]
}

func main(){
	for {
		values := RandomArray(4)
		knn := KNNClassifier{}
		knn.New("Testing", []string{"this sucks", "hiya"}, values,2,2)
		knn.Predict(knn.Data)
	}
}