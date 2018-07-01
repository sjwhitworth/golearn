package linear_models

import (
	"github.com/sjwhitworth/golearn/base"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestLogistic(t *testing.T) {
	Convey("Doing a logistic test", t, func() {
		X, err := base.ParseCSVToInstances("train.csv", false)
		So(err, ShouldEqual, nil)
		Y, err := base.ParseCSVToInstances("test.csv", false)
		So(err, ShouldEqual, nil)
		_, err = NewLogisticRegression("l0", 1.0, 1e-6)
		So(err, ShouldNotBeNil)
		lr, err := NewLogisticRegression("l1", 1.0, 1e-6)
		So(err, ShouldBeNil)

		lr.Fit(X)

		Convey("When predicting the label of first vector", func() {
			Z, err := lr.Predict(Y)
			So(err, ShouldEqual, nil)
			Convey("The result should be 1", func() {
				So(Z.RowString(0), ShouldEqual, "-1.0")
			})
		})
		Convey("When predicting the label of second vector", func() {
			Z, err := lr.Predict(Y)
			So(err, ShouldEqual, nil)
			Convey("The result should be -1", func() {
				So(Z.RowString(1), ShouldEqual, "-1.0")
			})
		})
		So((*lr).String(), ShouldEqual, "LogisticRegression")
	})
}
