package meta

import (
	"math/rand"
	"testing"
	"time"

	"github.com/sjwhitworth/golearn/base"
	"github.com/sjwhitworth/golearn/evaluation"
	"github.com/sjwhitworth/golearn/filters"
	"github.com/sjwhitworth/golearn/trees"
	. "github.com/smartystreets/goconvey/convey"
)

func BenchmarkBaggingRandomForestFit(t *testing.B) {
	inst, err := base.ParseCSVToInstances("../examples/datasets/iris_headers.csv", true)
	if err != nil {
		t.Fatalf("Unable to parse CSV to instances: %s", err.Error())
	}

	rand.Seed(time.Now().UnixNano())
	filt := filters.NewChiMergeFilter(inst, 0.90)
	for _, a := range base.NonClassFloatAttributes(inst) {
		filt.AddAttribute(a)
	}
	filt.Train()
	instf := base.NewLazilyFilteredInstances(inst, filt)

	rf := new(BaggedModel)
	for i := 0; i < 10; i++ {
		rf.AddModel(trees.NewRandomTree(2))
	}

	t.ResetTimer()
	for i := 0; i < 20; i++ {
		rf.Fit(instf)
	}
}

func BenchmarkBaggingRandomForestPredict(t *testing.B) {
	inst, err := base.ParseCSVToInstances("../examples/datasets/iris_headers.csv", true)
	if err != nil {
		t.Fatalf("Unable to parse CSV to instances: %s", err.Error())
	}

	rand.Seed(time.Now().UnixNano())
	filt := filters.NewChiMergeFilter(inst, 0.90)
	for _, a := range base.NonClassFloatAttributes(inst) {
		filt.AddAttribute(a)
	}
	filt.Train()
	instf := base.NewLazilyFilteredInstances(inst, filt)

	rf := new(BaggedModel)
	for i := 0; i < 10; i++ {
		rf.AddModel(trees.NewRandomTree(2))
	}

	rf.Fit(instf)
	t.ResetTimer()
	for i := 0; i < 20; i++ {
		rf.Predict(instf)
	}
}

func TestBaggedModelRandomForest(t *testing.T) {
	Convey("Given data", t, func() {
		inst, err := base.ParseCSVToInstances("../examples/datasets/iris_headers.csv", true)
		So(err, ShouldBeNil)

		Convey("Splitting the data into training and test data", func() {
			trainData, testData := base.InstancesTrainTestSplit(inst, 0.6)

			Convey("Filtering the split datasets", func() {
				rand.Seed(time.Now().UnixNano())
				filt := filters.NewChiMergeFilter(inst, 0.90)
				for _, a := range base.NonClassFloatAttributes(inst) {
					filt.AddAttribute(a)
				}
				filt.Train()
				trainDataf := base.NewLazilyFilteredInstances(trainData, filt)
				testDataf := base.NewLazilyFilteredInstances(testData, filt)

				Convey("Fitting and Predicting with a Bagged Model of 10 Random Trees", func() {
					rf := new(BaggedModel)
					for i := 0; i < 10; i++ {
						rf.AddModel(trees.NewRandomTree(2))
					}

					rf.Fit(trainDataf)
					predictions := rf.Predict(testDataf)

					confusionMat, err := evaluation.GetConfusionMatrix(testDataf, predictions)
					So(err, ShouldBeNil)

					Convey("Predictions are somewhat accurate", func() {
						So(evaluation.GetAccuracy(confusionMat), ShouldBeGreaterThan, 0.5)
					})
				})
			})
		})
	})
}
