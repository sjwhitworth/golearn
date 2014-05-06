package linear_models 

import (
	"testing"
	. "github.com/smartystreets/goconvey/convey"
)

func TestLogisticRegression(t *testing.T) {
	Convey("Given labels, a classifier and data", t, func() {
		X := [][]float64{
			{0, 0, 0, 1},
			{0, 0, 1, 0},
			{0, 1, 0, 0},
			{1, 0, 0, 0},
		}
		y := []float64{-1, -1, 1, 1}
		lr := NewLogisticRegression("l2", 1.0, 1e-6)
		lr.Fit(X,y)

		Convey("When predicting the label of first vector", func() {
			pred_x := [][]float64{ {1,1,0,0} }
			pred_y := lr.Predict(pred_x)
			Convey("The result should be 1", func() {
				So(pred_y[0], ShouldEqual, 1.0)
			})
		})
		Convey("When predicting the label of second vector", func() {
			pred_x := [][]float64{ {0,0,1,1} }
			pred_y := lr.Predict(pred_x)
			Convey("The result should be -1", func() {
				So(pred_y[0], ShouldEqual, -1.0)
			})
		})
	})
}
