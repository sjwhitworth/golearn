package linear_models

import (
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

	_, _ = predictions.Size()
}

func BenchmarkLinearRegressionOneRow(b *testing.B) {
	// Omits error handling in favor of brevity
	trainData, _ := base.ParseCSVToInstances("../examples/datasets/exams.csv", true)
	testData, _ := base.ParseCSVToInstances("../examples/datasets/exam.csv", true)
	lr := NewLinearRegression()
	lr.Fit(trainData)

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		lr.Predict(testData)
	}
}
