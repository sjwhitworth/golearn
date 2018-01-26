package meta

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"github.com/amclay/golearn/base"
	"github.com/amclay/golearn/evaluation"
	"github.com/amclay/golearn/linear_models"
	. "github.com/smartystreets/goconvey/convey"
)

func TestOneVsAllModel(t *testing.T) {

	classifierFunc := func(c string) base.Classifier {
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

		Convey("Saving and reloading should work...", func() {
			predictions, err := m.Predict(Y)
			So(err, ShouldEqual, nil)
			f, err := ioutil.TempFile(os.TempDir(), "tmpCls")
			defer func() {
				f.Close()
			}()
			err = m.Save(f.Name())
			So(err, ShouldBeNil)
			Convey("Reloaded classifier should output the same predictions", func() {
				m := NewOneVsAllModel(classifierFunc)
				err := m.Load(f.Name())
				So(err, ShouldBeNil)
				newPredictions, err := m.Predict(Y)
				So(err, ShouldBeNil)
				So(base.InstancesAreEqual(predictions, newPredictions), ShouldBeTrue)
			})
		})
	})
}
