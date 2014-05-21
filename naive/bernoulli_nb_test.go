package naive

import (
    "math"
    "github.com/sjwhitworth/golearn/base"
    "testing"
    . "github.com/smartystreets/goconvey/convey"
)

func TestFit(t *testing.T) {
    Convey("Given a simple training data", t, func() {
        trainingData, err1 := base.ParseCSVToInstances("test/simple_train.csv", false)
        if err1 != nil {
            t.Error(err1)
        }

        nb := NewBernoulliNBClassifier()
        nb.Fit(trainingData)

        Convey("All log(prior) should be correctly calculated", func() {
            logPriorBlue := nb.logClassPrior["blue"]
            logPriorRed := nb.logClassPrior["red"]

            So(logPriorBlue, ShouldAlmostEqual, math.Log(0.5))
            So(logPriorRed, ShouldAlmostEqual, math.Log(0.5))
        })

        Convey("'red' conditional probabilities should be correct", func() {
            logCondProbTok0 := nb.logCondProb["red"][0]
            logCondProbTok1 := nb.logCondProb["red"][1]
            logCondProbTok2 := nb.logCondProb["red"][2]

            So(logCondProbTok0, ShouldAlmostEqual, math.Log(1.0))
            So(logCondProbTok1, ShouldAlmostEqual, math.Log(1.0/3.0))
            So(logCondProbTok2, ShouldAlmostEqual, math.Log(1.0))
        })

        Convey("'blue' conditional probabilities should be correct", func() {
            logCondProbTok0 := nb.logCondProb["blue"][0]
            logCondProbTok1 := nb.logCondProb["blue"][1]
            logCondProbTok2 := nb.logCondProb["blue"][2]

            So(logCondProbTok0, ShouldAlmostEqual, math.Log(1.0))
            So(logCondProbTok1, ShouldAlmostEqual, math.Log(1.0))
            So(logCondProbTok2, ShouldAlmostEqual, math.Log(1.0/3.0))
        })
    })
}
