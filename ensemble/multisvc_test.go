package ensemble

import (
	"github.com/sjwhitworth/golearn/base"
	"github.com/sjwhitworth/golearn/evaluation"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestMultiSVMUnweighted(t *testing.T) {
	Convey("Loading data...", t, func() {
		inst, err := base.ParseCSVToInstances("../examples/datasets/articles.csv", false)
		So(err, ShouldBeNil)
		X, Y := base.InstancesTrainTestSplit(inst, 0.4)

		m := NewMultiLinearSVC("l1", "l2", true, 1.0, 1e-4, nil)
		m.Fit(X)

		Convey("Predictions should work...", func() {
			predictions, err := m.Predict(Y)
			cf, err := evaluation.GetConfusionMatrix(Y, predictions)
			So(err, ShouldEqual, nil)
			So(evaluation.GetAccuracy(cf), ShouldBeGreaterThan, 0.70)
		})
	})
}

func TestMultiSVMWeighted(t *testing.T) {
	Convey("Loading data...", t, func() {
		weights := make(map[string]float64)
		weights["Finance"] = 0.1739
		weights["Tech"] = 0.0750
		weights["Politics"] = 0.4928

		inst, err := base.ParseCSVToInstances("../examples/datasets/articles.csv", false)
		So(err, ShouldBeNil)
		X, Y := base.InstancesTrainTestSplit(inst, 0.4)

		m := NewMultiLinearSVC("l1", "l2", true, 0.62, 1e-4, weights)
		m.Fit(X)

		Convey("Predictions should work...", func() {
			predictions, err := m.Predict(Y)
			cf, err := evaluation.GetConfusionMatrix(Y, predictions)
			So(err, ShouldEqual, nil)
			So(evaluation.GetAccuracy(cf), ShouldBeGreaterThan, 0.70)
		})
	})
}
