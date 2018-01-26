package perceptron

import (
	"fmt"
	"path/filepath"
	"testing"

	"github.com/amclay/golearn/base"
	"github.com/amclay/golearn/evaluation"
)

func TestProcessData(t *testing.T) {
	absPath, _ := filepath.Abs("../examples/datasets/house-votes-84.csv")
	rawData, err := base.ParseCSVToInstances(absPath, true)
	trainData, _ := base.InstancesTrainTestSplit(rawData, 0.5)

	if err != nil {
		t.Fatal("Could not test processData. Could not load CSV")
	}

	if rawData == nil {
		t.Fatal("Could not test processData. Could not load CSV")
	}

	result := processData(trainData)
	_, size := trainData.Size()

	if len(result) != size {
		t.Errorf("Expected %d, Got %d", size, len(result))
	}
}

func TestFit(t *testing.T) {

	a := NewAveragePerceptron(10, 1.2, 0.5, 0.3)

	if a == nil {

		t.Errorf("Unable to create average perceptron")
	}

	absPath, _ := filepath.Abs("../examples/datasets/house-votes-84.csv")
	rawData, err := base.ParseCSVToInstances(absPath, true)
	if err != nil {
		t.Fail()
	}

	trainData, _ := base.InstancesTrainTestSplit(rawData, 0.7)
	a.Fit(trainData)

	if a.trained == false {
		t.Errorf("Perceptron was not trained")
	}

}

func TestPredict(t *testing.T) {

	a := NewAveragePerceptron(10, 1.2, 0.5, 0.3)

	if a == nil {

		t.Errorf("Unable to create average perceptron")
	}

	absPath, _ := filepath.Abs("../examples/datasets/house-votes-84.csv")
	rawData, err := base.ParseCSVToInstances(absPath, true)
	if err != nil {
		t.Fail()
	}

	trainData, testData := base.InstancesTrainTestSplit(rawData, 0.5)
	a.Fit(trainData)

	if a.trained == false {
		t.Errorf("Perceptron was not trained")
	}

	predictions := a.Predict(testData)
	cf, err := evaluation.GetConfusionMatrix(testData, predictions)
	if err != nil {
		t.Errorf("Couldn't get confusion matrix: %s", err)
		t.Fail()
	}
	fmt.Println(evaluation.GetSummary(cf))
	fmt.Println(trainData)
	fmt.Println(testData)
	if evaluation.GetAccuracy(cf) < 0.65 {
		t.Errorf("Perceptron not trained correctly")
	}
}

func TestCreateAveragePerceptron(t *testing.T) {

	a := NewAveragePerceptron(10, 1.2, 0.5, 0.3)

	if a == nil {

		t.Errorf("Unable to create average perceptron")
	}
}

func BenchmarkFit(b *testing.B) {

	a := NewAveragePerceptron(10, 1.2, 0.5, 0.3)
	absPath, _ := filepath.Abs("../examples/datasets/house-votes-84.csv")
	rawData, _ := base.ParseCSVToInstances(absPath, true)
	trainData, _ := base.InstancesTrainTestSplit(rawData, 0.5)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		a.Fit(trainData)
	}
}
