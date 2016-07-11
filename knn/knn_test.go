package knn

import (
	"github.com/sjwhitworth/golearn/base"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestKnnClassifierWithoutOptimisations(t *testing.T) {
	Convey("Given labels, a classifier and data", t, func() {
		trainingData, err := base.ParseCSVToInstances("knn_train_1.csv", false)
		So(err, ShouldBeNil)

		testingData, err := base.ParseCSVToInstances("knn_test_1.csv", false)
		So(err, ShouldBeNil)

		cls := NewKnnClassifier("euclidean", 2)
		cls.AllowOptimisations = false
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

func TestKnnClassifierWithOptimisations(t *testing.T) {
	Convey("Given labels, a classifier and data", t, func() {
		trainingData, err := base.ParseCSVToInstances("knn_train_1.csv", false)
		So(err, ShouldBeNil)

		testingData, err := base.ParseCSVToInstances("knn_test_1.csv", false)
		So(err, ShouldBeNil)

		cls := NewKnnClassifier("euclidean", 2)
		cls.AllowOptimisations = true
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

func TestKnnClassifierWithTemplatedInstances1(t *testing.T) {
	Convey("Given two basically identical files...", t, func() {
		trainingData, err := base.ParseCSVToInstances("knn_train_2.csv", true)
		So(err, ShouldBeNil)
		testingData, err := base.ParseCSVToTemplatedInstances("knn_test_2.csv", true, trainingData)
		So(err, ShouldBeNil)

		cls := NewKnnClassifier("euclidean", 2)
		cls.Fit(trainingData)
		predictions := cls.Predict(testingData)
		So(predictions, ShouldNotBeNil)
	})
}

func TestKnnClassifierWithTemplatedInstances1Subset(t *testing.T) {
	Convey("Given two basically identical files...", t, func() {
		trainingData, err := base.ParseCSVToInstances("knn_train_2.csv", true)
		So(err, ShouldBeNil)
		testingData, err := base.ParseCSVToTemplatedInstances("knn_test_2_subset.csv", true, trainingData)
		So(err, ShouldBeNil)

		cls := NewKnnClassifier("euclidean", 2)
		cls.Fit(trainingData)
		predictions := cls.Predict(testingData)
		So(predictions, ShouldNotBeNil)
	})
}
