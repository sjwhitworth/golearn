package linear_models

import (
	"github.com/sjwhitworth/golearn/base"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestLinearSVC(t *testing.T) {
	Convey("Given labels, a classifier and data", t, func() {
		trainingData, err := base.ParseCSVToInstances("linearsvc_train.csv", false)
		So(err, ShouldBeNil)

		testingData, err := base.ParseCSVToInstances("linearsvc_test.csv", false)
		So(err, ShouldBeNil)

		cls, err := NewLinearSVC("l2", "l2", true, 1.0, 1e-4)
		So(err, ShouldBeNil)

		cls.Fit(trainingData)

		predictions, _ := cls.Predict(testingData)
		So(predictions, ShouldNotEqual, nil)

		Convey("When predicting the label for our first vector", func() {
			result := base.GetClass(predictions, 0)
			Convey("The result should be 1.0", func() {
				So(result, ShouldEqual, "1.0")
			})
		})

		Convey("When predicting the label for our second vector", func() {
			result := base.GetClass(predictions, 1)
			Convey("The result should be -1.0", func() {
				So(result, ShouldEqual, "-1.0")
			})
		})

	})
}
