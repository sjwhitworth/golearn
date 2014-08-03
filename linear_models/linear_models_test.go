package linear_models

import (
	"github.com/sjwhitworth/golearn/base"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestLogisticRegression(t *testing.T) {
	Convey("Given labels, a classifier and data", t, func() {
		// Load data
		X, err := base.ParseCSVToInstances("train.csv", false)
		So(err, ShouldEqual, nil)
		Y, err := base.ParseCSVToInstances("test.csv", false)
		So(err, ShouldEqual, nil)

		// Setup the problem
		lr := NewLogisticRegression("l2", 1.0, 1e-6)
		lr.Fit(X)

		Convey("When predicting the label of first vector", func() {
			Z := lr.Predict(Y)
			Convey("The result should be 1", func() {
				So(Z.RowString(0), ShouldEqual, "1.00")
			})
		})
		Convey("When predicting the label of second vector", func() {
			Z := lr.Predict(Y)
			Convey("The result should be -1", func() {
				So(Z.RowString(1), ShouldEqual, "-1.00")
			})
		})
	})
}
