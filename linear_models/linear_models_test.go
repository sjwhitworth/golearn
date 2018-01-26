package linear_models

import (
	"testing"

	"github.com/amclay/golearn/base"
	. "github.com/smartystreets/goconvey/convey"
)

func TestLogisticRegression(t *testing.T) {
	Convey("Given labels, a classifier and data", t, func() {
		// Load data
		X, err := base.ParseCSVToInstances("train.csv", false)
		So(err, ShouldEqual, nil)
		Y, err := base.ParseCSVToInstances("test.csv", false)
		So(err, ShouldEqual, nil)

		// Setup the problem
		lr, err := NewLogisticRegression("l2", 1.0, 1e-6)
		So(err, ShouldBeNil)

		lr.Fit(X)

		Convey("When predicting the label of first vector", func() {
			Z, err := lr.Predict(Y)
			So(err, ShouldEqual, nil)
			Convey("The result should be 1", func() {
				So(Z.RowString(0), ShouldEqual, "1.0")
			})
		})
		Convey("When predicting the label of second vector", func() {
			Z, err := lr.Predict(Y)
			So(err, ShouldEqual, nil)
			Convey("The result should be -1", func() {
				So(Z.RowString(1), ShouldEqual, "-1.0")
			})
		})
	})
}
