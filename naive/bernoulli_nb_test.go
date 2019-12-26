package naive

import (
	"github.com/sjwhitworth/golearn/base"
	"github.com/sjwhitworth/golearn/filters"
	. "github.com/smartystreets/goconvey/convey"
	"io/ioutil"
	"os"
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

func TestSerialize(t *testing.T) {
	Convey("Given simple training/test data", t, func() {
		trainingData, err := base.ParseCSVToInstances("test/simple_train.csv", false)
		So(err, ShouldBeNil)

		testData, err := base.ParseCSVToTemplatedInstances("test/simple_test.csv", false, trainingData)
		So(err, ShouldBeNil)

		nb := NewBernoulliNBClassifier()
		nb.Fit(convertToBinary(trainingData))
		oldPredictions, err := nb.Predict(convertToBinary(testData))

		Convey("Saving the classifer should work...", func() {
			f, err := ioutil.TempFile(os.TempDir(), "nb")
			So(err, ShouldBeNil)
			defer func() {
				f.Close()
			}()
			err = nb.Save(f.Name())
			So(err, ShouldBeNil)
			Convey("Loading the classifier should work...", func() {
				newNb := NewBernoulliNBClassifier()
				err := newNb.Load(f.Name())
				So(err, ShouldBeNil)
				Convey("Predictions should match...", func() {
					newPredictions, err := newNb.Predict(convertToBinary(testData))
					So(err, ShouldBeNil)
					So(base.InstancesAreEqual(oldPredictions, newPredictions), ShouldBeTrue)
				})
			})
		})
	})
}

func TestSimple(t *testing.T) {
	Convey("Given a simple training dataset", t, func() {
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
			testData, err := base.ParseCSVToTemplatedInstances("test/simple_test.csv", false, trainingData)
			So(err, ShouldBeNil)

			predictions, err := nb.Predict(convertToBinary(testData))
			So(err, ShouldBeNil)

			Convey("All simple predictions should be correct", func() {
				So(base.GetClass(predictions, 0), ShouldEqual, "blue")
				So(base.GetClass(predictions, 1), ShouldEqual, "red")
				So(base.GetClass(predictions, 2), ShouldEqual, "blue")
				So(base.GetClass(predictions, 3), ShouldEqual, "red")
			})
		})
	})
}
