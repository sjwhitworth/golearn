package naive

import (
    "github.com/sjwhitworth/golearn/base"
    "testing"
    . "github.com/smartystreets/goconvey/convey"
)

func TestNoFit(t *testing.T) {
    Convey("Given an empty BernoulliNaiveBayes", t, func() {
        nb := NewBernoulliNBClassifier()

        Convey("PredictOne should panic if Fit was not called", func() {
            testDoc := []float64{0.0, 1.0}
            So(func() { nb.PredictOne(testDoc) }, ShouldPanic)
        })
    })
}

func TestSimple(t *testing.T) {
    Convey("Given a simple training data", t, func() {
        trainingData, err1 := base.ParseCSVToInstances("test/simple_train.csv", false)
        if err1 != nil {
            t.Error(err1)
        }

        nb := NewBernoulliNBClassifier()
        nb.Fit(trainingData)

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
                testDoc := []float64{0.0, 2.0}
                So(func() { nb.PredictOne(testDoc) }, ShouldPanic)
            })

            Convey("Token 1 should be a good predictor of the blue class", func() {
                testDoc := []float64{0.0, 123.0, 0.0}
                So(nb.PredictOne(testDoc), ShouldEqual, "blue")

                testDoc = []float64{120.0, 123.0, 0.0}
                So(nb.PredictOne(testDoc), ShouldEqual, "blue")
            })

            Convey("Token 2 should be a good predictor of the red class", func() {
                testDoc := []float64{0.0, 0.0, 120.0}
                So(nb.PredictOne(testDoc), ShouldEqual, "red")

                testDoc = []float64{10.0, 0.0, 120.0}
                So(nb.PredictOne(testDoc), ShouldEqual, "red")
            })
        })

        Convey("Predict should work as expected", func() {
            testData, err := base.ParseCSVToInstances("test/simple_test.csv", false)
            if err != nil {
                t.Error(err)
            }
            predictions := nb.Predict(testData)

            Convey("All simple predicitions should be correct", func() {
                So(predictions.GetClass(0), ShouldEqual, "blue")
                So(predictions.GetClass(1), ShouldEqual, "red")
                So(predictions.GetClass(2), ShouldEqual, "blue")
                So(predictions.GetClass(3), ShouldEqual, "red")
            })
        })
    })
}
