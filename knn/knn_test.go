package knn

import (
	"github.com/sjwhitworth/golearn/base"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestKnnClassifier(t *testing.T) {
	Convey("Given labels, a classifier and data", t, func() {

		trainingData, err1 := base.ParseCSVToInstances("knn_train.csv", false)
		testingData, err2 := base.ParseCSVToInstances("knn_test.csv", false)

		if err1 != nil {
			t.Error(err1)
			return
		}
		if err2 != nil {
			t.Error(err2)
			return
		}

		cls := NewKnnClassifier("euclidean", 2)
		cls.Fit(trainingData)
		predictions := cls.Predict(testingData)

		Convey("When predicting the label for our first vector", func() {
			result := base.GetClass(predictions, 0)
			Convey("The result should be 'blue", func() {
				So(result, ShouldEqual, "blue")
			})
		})

		Convey("When predicting the label for our first vector", func() {
			result2 := base.GetClass(predictions, 1)
			Convey("The result should be 'red", func() {
				So(result2, ShouldEqual, "red")
			})
		})
	})
}
