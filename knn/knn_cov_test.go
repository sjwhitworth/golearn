package knn

import (
	"testing"

	"github.com/sjwhitworth/golearn/base"
	. "github.com/smartystreets/goconvey/convey"
)

func TestKnnClassifierCov(t *testing.T) {
	Convey("Test predict", t, func() {
		Convey("distance function", func() {
			trainingData, err := base.ParseCSVToInstances("knn_train_1.csv", false)
			So(err, ShouldBeNil)

			testingData, err := base.ParseCSVToInstances("knn_test_1.csv", false)
			So(err, ShouldBeNil)

			Convey("use euclidean", func() {
				cls := NewKnnClassifier("euclidean", "kdtree", 2)
				cls.AllowOptimisations = false
				cls.Fit(trainingData)
				predictions, err := cls.Predict(testingData)
				So(err, ShouldBeNil)
				So(predictions, ShouldNotEqual, nil)
				result := base.GetClass(predictions, 0)
				So(result, ShouldEqual, "blue")
			})

			Convey("use manhattan", func() {
				cls := NewKnnClassifier("manhattan", "kdtree", 2)
				cls.AllowOptimisations = false
				cls.Fit(trainingData)
				predictions, err := cls.Predict(testingData)
				So(err, ShouldBeNil)
				So(predictions, ShouldNotEqual, nil)
				result := base.GetClass(predictions, 0)
				So(result, ShouldEqual, "blue")
			})

			Convey("use cosine", func() {
				cls := NewKnnClassifier("cosine", "kdtree", 2)
				cls.AllowOptimisations = false
				cls.Fit(trainingData)
				predictions, err := cls.Predict(testingData)
				So(err, ShouldBeNil)
				So(predictions, ShouldNotEqual, nil)
				result := base.GetClass(predictions, 0)
				So(result, ShouldEqual, "blue")
			})

			Convey("use undefined distance function", func() {
				cls := NewKnnClassifier("abcd", "kdtree", 2)
				cls.AllowOptimisations = false
				cls.Fit(trainingData)
				predictions, err := cls.Predict(testingData)
				So(predictions, ShouldBeNil)
				So(err.Error(), ShouldEqual, "unsupported distance function")
			})
		})

		Convey("searching algorithm", func() {
			trainingData, err := base.ParseCSVToInstances("knn_train_1.csv", false)
			So(err, ShouldBeNil)

			testingData, err := base.ParseCSVToInstances("knn_test_1.csv", false)
			So(err, ShouldBeNil)

			Convey("use undefined searching algorithm", func() {
				cls := NewKnnClassifier("cosine", "abcd", 2)
				cls.AllowOptimisations = false
				cls.Fit(trainingData)
				predictions, err := cls.Predict(testingData)
				So(predictions, ShouldBeNil)
				So(err.Error(), ShouldEqual, "unsupported searching algorithm")
			})
		})

		Convey("check features", func() {
			Convey("use different dataset", func() {
				trainingData, err := base.ParseCSVToInstances("knn_train_1.csv", false)
				So(err, ShouldBeNil)
				testingData, err := base.ParseCSVToInstances("knn_test_2.csv", false)
				So(err, ShouldBeNil)
				cls := NewKnnClassifier("cosine", "linear", 2)
				cls.AllowOptimisations = false
				cls.Fit(trainingData)
				predictions, err := cls.Predict(testingData)
				So(predictions, ShouldBeNil)
				So(err.Error(), ShouldEqual, "attributes not compatible")
			})
		})
	})
}
