package main

import (
		mat "github.com/skelterjohn/go.matrix"
		"fmt"
		)

type KNNClassifier struct {
	Data mat.DenseMatrix
	Name string
}

func (Class *KNNClassifier) New(name string, numbers []float64, x int, y int){
	Class.Data = *mat.MakeDenseMatrix(numbers, x, y)
	Class.Name = name
}


func main(){
	knn := KNNClassifier{}
	knn.New("Testing", []float64{1,2,3,4,5,6},2,3)
	another := KNNClassifier{}
	another.New("Blah", []float64{2,4,5,3,4,6},2,3)
	fmt.Println(mat.Difference(&knn.Data, &another.Data))
}