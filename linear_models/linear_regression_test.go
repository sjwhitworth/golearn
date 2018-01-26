package linear_models

import (
	"strconv"
	"testing"

	"github.com/amclay/golearn/base"
	. "github.com/smartystreets/goconvey/convey"
)

func TestLinearRegression(t *testing.T) {
	Convey("Doing a  linear regression", t, func() {
		lr := NewLinearRegression()

		Convey("With no training data", func() {
			Convey("Predicting", func() {
				testData, err := base.ParseCSVToInstances("../examples/datasets/exams.csv", true)
				So(err, ShouldBeNil)

				_, err = lr.Predict(testData)

				Convey("Should result in a NoTrainingDataError", func() {
					So(err, ShouldEqual, NoTrainingDataError)
				})

			})
		})

		Convey("With not enough training data", func() {
			trainingDatum, err := base.ParseCSVToInstances("../examples/datasets/exam.csv", true)
			So(err, ShouldBeNil)

			Convey("Fitting", func() {
				err = lr.Fit(trainingDatum)

				Convey("Should result in a NotEnoughDataError", func() {
					So(err, ShouldEqual, NotEnoughDataError)
				})
			})
		})

		Convey("With sufficient training data", func() {
			instances, err := base.ParseCSVToInstances("../examples/datasets/exams.csv", true)
			So(err, ShouldBeNil)
			trainData, testData := base.InstancesTrainTestSplit(instances, 0.1)

			Convey("Fitting and Predicting", func() {
				err := lr.Fit(trainData)
				So(err, ShouldBeNil)

				predictions, err := lr.Predict(testData)
				So(err, ShouldBeNil)

				Convey("It makes reasonable predictions", func() {
					_, rows := predictions.Size()

					for i := 0; i < rows; i++ {
						actualValue, _ := strconv.ParseFloat(base.GetClass(testData, i), 64)
						expectedValue, _ := strconv.ParseFloat(base.GetClass(predictions, i), 64)

						So(actualValue, ShouldAlmostEqual, expectedValue, actualValue*0.05)
					}
				})
			})
		})
	})
}

func BenchmarkLinearRegressionOneRow(b *testing.B) {
	// Omits error handling in favor of brevity
	trainData, _ := base.ParseCSVToInstances("../examples/datasets/exams.csv", true)
	testData, _ := base.ParseCSVToInstances("../examples/datasets/exam.csv", true)
	lr := NewLinearRegression()
	lr.Fit(trainData)

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		lr.Predict(testData)
	}
}
