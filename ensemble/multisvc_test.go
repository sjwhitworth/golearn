package ensemble

import (
	"fmt"
	"github.com/sjwhitworth/golearn/base"
	"github.com/sjwhitworth/golearn/evaluation"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestMultiSVM(t *testing.T) {
	Convey("Loading data...", t, func() {
		inst, err := base.ParseCSVToInstances("../examples/datasets/articles.csv", false)
		So(err, ShouldBeNil)
		X, Y := base.InstancesTrainTestSplit(inst, 0.4)

		m := NewMultiLinearSVC("l1", "l2", true, 1.0, 1e-4)
		m.Fit(X)

		Convey("Predictions should work...", func() {
			predictions, err := m.Predict(Y)
			cf, err := evaluation.GetConfusionMatrix(Y, predictions)
			So(err, ShouldEqual, nil)
			fmt.Println(evaluation.GetSummary(cf))
		})
	})
}
