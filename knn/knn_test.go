package knn

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestKnnClassifier(t *testing.T) {
	Convey("Given labels, a classifier and data", t, func() {
		labels := []string{"blue", "blue", "red", "red"}
		data := []float64{1, 1, 1, 1, 1, 1, 3, 3, 3, 6, 6, 6}
		cls := NewKnnClassifier(labels, data, 4, 3, "euclidean")

		Convey("When predicting the label for our first vector", func() {
			// The vector we're going to predict
			vector := []float64{1.2, 1.2, 1.5}
			result := cls.Predict(vector, 2)
			Convey("The result should be 'blue", func() {
				So(result, ShouldEqual, "blue")
			})
		})

		Convey("When predicting the label for our first vector", func() {
			// The vector we're going to predict
			vector2 := []float64{5, 5, 5}
			result2 := cls.Predict(vector2, 2)
			Convey("The result should be 'red", func() {
				So(result2, ShouldEqual, "red")
			})
		})
	})
}
