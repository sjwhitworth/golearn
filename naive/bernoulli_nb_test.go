package naive

import (
	"github.com/sjwhitworth/golearn/base"
	"github.com/sjwhitworth/golearn/filters"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestNoFit(t *testing.T) {
	Convey("Given an empty BernoulliNaiveBayes", t, func() {
		nb := NewBernoulliNBClassifier()

		Convey("PredictOne should panic if Fit was not called", func() {
			testDoc := [][]byte{[]byte{0}, []byte{1}}
			So(func() { nb.PredictOne(testDoc) }, ShouldPanic)
		})
	})
}

func convertToBinary(src base.FixedDataGrid) base.FixedDataGrid {
	// Convert to binary
	b := filters.NewBinaryConvertFilter()
	attrs := base.NonClassAttributes(src)
	for _, a := range attrs {
		b.AddAttribute(a)
	}
	b.Train()
	ret := base.NewLazilyFilteredInstances(src, b)
	return ret
}

func TestSimple(t *testing.T) {
	Convey("Given a simple training data", t, func() {
		trainingData, err := base.ParseCSVToInstances("test/simple_train.csv", false)
		So(err, ShouldBeNil)

		nb := NewBernoulliNBClassifier()
		nb.Fit(convertToBinary(trainingData))

		Convey("Check if Fit is working as expected", func() {
			Convey("All data needed for prior should be correctly calculated", func() {
				So(nb.classInstances["blue"], ShouldEqual, 2)
				So(nb.classInstances["red"], ShouldEqual, 2)
				So(nb.trainingInstances, ShouldEqual, 4)
			})

			Convey("'red' conditional probabilities should be correct", func() {
				logCondProbTok0 := nb.condProb["red"][0]
				logCondProbTok1 := nb.condProb["red"][1]
				logCondProbTok2 := nb.condProb["red"][2]

				So(logCondProbTok0, ShouldAlmostEqual, 1.0)
				So(logCondProbTok1, ShouldAlmostEqual, 1.0/3.0)
				So(logCondProbTok2, ShouldAlmostEqual, 1.0)
			})

			Convey("'blue' conditional probabilities should be correct", func() {
				logCondProbTok0 := nb.condProb["blue"][0]
				logCondProbTok1 := nb.condProb["blue"][1]
				logCondProbTok2 := nb.condProb["blue"][2]

				So(logCondProbTok0, ShouldAlmostEqual, 1.0)
				So(logCondProbTok1, ShouldAlmostEqual, 1.0)
				So(logCondProbTok2, ShouldAlmostEqual, 1.0/3.0)
			})
		})

		Convey("PredictOne should work as expected", func() {
			Convey("Using a document with different number of cols should panic", func() {
				testDoc := [][]byte{[]byte{0}, []byte{2}}
				So(func() { nb.PredictOne(testDoc) }, ShouldPanic)
			})

			Convey("Token 1 should be a good predictor of the blue class", func() {
				testDoc := [][]byte{[]byte{0}, []byte{1}, []byte{0}}
				So(nb.PredictOne(testDoc), ShouldEqual, "blue")

				testDoc = [][]byte{[]byte{1}, []byte{1}, []byte{0}}
				So(nb.PredictOne(testDoc), ShouldEqual, "blue")
			})

			Convey("Token 2 should be a good predictor of the red class", func() {
				testDoc := [][]byte{[]byte{0}, []byte{0}, []byte{1}}
				So(nb.PredictOne(testDoc), ShouldEqual, "red")
				testDoc = [][]byte{[]byte{1}, []byte{0}, []byte{1}}
				So(nb.PredictOne(testDoc), ShouldEqual, "red")
			})
		})

		Convey("Predict should work as expected", func() {
			testData, err := base.ParseCSVToInstances("test/simple_test.csv", false)
			So(err, ShouldBeNil)

			predictions := nb.Predict(convertToBinary(testData))

			Convey("All simple predicitions should be correct", func() {
				So(base.GetClass(predictions, 0), ShouldEqual, "blue")
				So(base.GetClass(predictions, 1), ShouldEqual, "red")
				So(base.GetClass(predictions, 2), ShouldEqual, "blue")
				So(base.GetClass(predictions, 3), ShouldEqual, "red")
			})
		})
	})
}
