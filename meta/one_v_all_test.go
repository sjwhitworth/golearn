package meta

import (
	"fmt"
	"github.com/sjwhitworth/golearn/base"
	"github.com/sjwhitworth/golearn/evaluation"
	"github.com/sjwhitworth/golearn/linear_models"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestOneVsAllModel(t *testing.T) {

	classifierFunc := func() base.Classifier {
		m, err := linear_models.NewLinearSVC("l1", "l2", true, 1.0, 1e-4)
		if err != nil {
			panic(err)
		}
		return m
	}

	Convey("Given data", t, func() {
		inst, err := base.ParseCSVToInstances("../examples/datasets/iris_headers.csv", true)
		So(err, ShouldBeNil)

		X, Y := base.InstancesTrainTestSplit(inst, 0.4)

		m := NewOneVsAllModel(classifierFunc)
		m.Fit(X)

		Convey("The maximum class index should be 2", func() {
			So(m.maxClassVal, ShouldEqual, 2)
		})

		Convey("There should be three of everything...", func() {
			So(len(m.filters), ShouldEqual, 3)
			So(len(m.classifiers), ShouldEqual, 3)
		})

		Convey("Predictions should work...", func() {
			predictions, err := m.Predict(Y)
			So(err, ShouldEqual, nil)
			cf, err := evaluation.GetConfusionMatrix(Y, predictions)
			So(err, ShouldEqual, nil)
			fmt.Println(evaluation.GetAccuracy(cf))
			fmt.Println(evaluation.GetSummary(cf))
		})
	})
}
