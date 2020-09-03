package clustering

import (
	"github.com/sjwhitworth/golearn/base"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestExpectationMaximization(t *testing.T) {
	Convey("Doing EM-based clustering", t, func() {
		em, _ := NewExpectationMaximization(2)

		// Initialization tests
		// Trying to create NewExpectationMaximization with < 1 component
		Convey("With less than one component", func() {
			Convey("Creating a new instance", func() {
				_, err := NewExpectationMaximization(0)
				Convey("Should result in a InsufficientComponentsError", func() {
					So(err, ShouldEqual, InsufficientComponentsError)
				})
			})
		})

		// Data tests
		// Trying to Fit with fewer samples than components
		Convey("With insufficient training data", func() {
			Convey("Fitting", func() {
				testData, err := base.ParseCSVToInstances("./gaussian_mixture_single_obs.csv", false)
				So(err, ShouldBeNil)

				err = em.Fit(testData)

				Convey("Should result in a InsufficientDataError", func() {
					So(err, ShouldEqual, InsufficientDataError)
				})
			})
		})

		// Trying to Predict before having Fit
		Convey("With no training data", func() {
			Convey("Predicting", func() {
				testData, err := base.ParseCSVToInstances("./gaussian_mixture.csv", false)
				So(err, ShouldBeNil)

				_, err = em.Predict(testData)

				Convey("Should result in a NoTrainingDataError", func() {
					So(err, ShouldEqual, NoTrainingDataError)
				})
			})
		})

		// Computation tests
		// Test the predictions are resonable
		Convey("With sufficient training data", func() {
			instances, err := base.ParseCSVToInstances("./gaussian_mixture.csv", true)
			So(err, ShouldBeNil)

			Convey("Fitting", func() {
				err := em.Fit(instances)
				So(err, ShouldBeNil)

				first_mean := em.Params.Means.At(0, 0)

				Convey("It converges to reasonable a value", func() {
					So(first_mean, ShouldAlmostEqual, -5.973, .1)
				})
			})
		})

		Convey("Test more code", func() {
			trainData, _ := base.ParseCSVToInstances("./gaussian_mixture.csv", false)
			testData, _ := base.ParseCSVToInstances("./gaussian_mixture.csv", false)

			em, err := NewExpectationMaximization(1)
			if err != nil {
				panic(err)
			}
			em.Fit(trainData)
			em.Predict(testData)
		})
	})
}

func BenchmarkExpectationMaximizationOneRow(b *testing.B) {
	// Omits error handling in favor of brevity
	trainData, _ := base.ParseCSVToInstances("./gaussian_mixture.csv", false)
	testData, _ := base.ParseCSVToInstances("./gaussian_mixture.csv", false)

	em, err := NewExpectationMaximization(2)
	if err != nil {
		panic(err)
	}
	em.Fit(trainData)

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		em.Predict(testData)
	}
}
