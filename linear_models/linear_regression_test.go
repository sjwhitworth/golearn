package linear_models

import (
	"fmt"
	"testing"

	"github.com/sjwhitworth/golearn/base"
)

func TestNoTrainingData(t *testing.T) {
	lr := NewLinearRegression()

	rawData, err := base.ParseCSVToInstances("../examples/datasets/exams.csv", true)
	if err != nil {
		t.Fatal(err)
	}

	_, err = lr.Predict(rawData)
	if err != NoTrainingDataError {
		t.Fatal("failed to error out even if no training data exists")
	}
}

func TestNotEnoughTrainingData(t *testing.T) {
	lr := NewLinearRegression()

	rawData, err := base.ParseCSVToInstances("../examples/datasets/exam.csv", true)
	if err != nil {
		t.Fatal(err)
	}

	err = lr.Fit(rawData)
	if err != NotEnoughDataError {
		t.Fatal("failed to error out even though there was not enough data")
	}
}

func TestLinearRegression(t *testing.T) {
	lr := NewLinearRegression()

	rawData, err := base.ParseCSVToInstances("../examples/datasets/exams.csv", true)
	if err != nil {
		t.Fatal(err)
	}

	trainData, testData := base.InstancesTrainTestSplit(rawData, 0.1)
	err = lr.Fit(trainData)
	if err != nil {
		t.Fatal(err)
	}

	predictions, err := lr.Predict(testData)
	if err != nil {
		t.Fatal(err)
	}

	for i := 0; i < predictions.Rows; i++ {
		fmt.Printf("Expected: %f || Predicted: %f\n", testData.Get(i, testData.ClassIndex), predictions.Get(i, predictions.ClassIndex))
	}
}
