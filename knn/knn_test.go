package knn

import (
	"github.com/sjwhitworth/golearn/base"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestKnnClassifier(t *testing.T) {
	Convey("Given labels, a classifier and data", t, func() {
		trainingData, err := base.ParseCSVToInstances("knn_train.csv", false)
		So(err, ShouldBeNil)

		testingData, err := base.ParseCSVToInstances("knn_test.csv", false)
		So(err, ShouldBeNil)

		cls := NewKnnClassifier("euclidean", 2)
		cls.Fit(trainingData)
		predictions := cls.Predict(testingData)
		So(predictions, ShouldNotEqual, nil)

		Convey("When predicting the label for our first vector", func() {
			result := base.GetClass(predictions, 0)
			Convey("The result should be 'blue", func() {
				So(result, ShouldEqual, "blue")
			})
		})

		Convey("When predicting the label for our second vector", func() {
			result2 := base.GetClass(predictions, 1)
			Convey("The result should be 'red", func() {
				So(result2, ShouldEqual, "red")
			})
		})
	})
}
