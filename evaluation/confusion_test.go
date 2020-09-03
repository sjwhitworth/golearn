package evaluation

import (
	"github.com/sjwhitworth/golearn/base"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestMetrics(t *testing.T) {
	Convey("Quantities derived from a confusion matrix", t, func() {
		confusionMat := make(ConfusionMatrix)
		confusionMat["a"] = make(map[string]int)
		confusionMat["b"] = make(map[string]int)
		confusionMat["a"]["a"] = 75
		confusionMat["a"]["b"] = 5
		confusionMat["b"]["a"] = 10
		confusionMat["b"]["b"] = 10

		Convey("True Positives", func() {
			So(GetTruePositives("a", confusionMat), ShouldAlmostEqual, 75, 1)
			So(GetTruePositives("b", confusionMat), ShouldAlmostEqual, 10, 1)
		})

		Convey("True Negatives", func() {
			So(GetTrueNegatives("a", confusionMat), ShouldAlmostEqual, 10, 1)
			So(GetTrueNegatives("b", confusionMat), ShouldAlmostEqual, 75, 1)
		})

		Convey("False Positives", func() {
			So(GetFalsePositives("a", confusionMat), ShouldAlmostEqual, 10, 1)
			So(GetFalsePositives("b", confusionMat), ShouldAlmostEqual, 5, 1)
		})

		Convey("False Negatives", func() {
			So(GetFalseNegatives("a", confusionMat), ShouldAlmostEqual, 5, 1)
			So(GetFalseNegatives("b", confusionMat), ShouldAlmostEqual, 10, 1)
		})

		Convey("Precision", func() {
			So(GetPrecision("a", confusionMat), ShouldAlmostEqual, 0.88, 0.01)
			So(GetPrecision("b", confusionMat), ShouldAlmostEqual, 0.666, 0.01)
		})

		Convey("Recall", func() {
			So(GetRecall("a", confusionMat), ShouldAlmostEqual, 0.94, 0.01)
			So(GetRecall("b", confusionMat), ShouldAlmostEqual, 0.50, 0.01)
		})

		Convey("MicroPrecision", func() {
			So(GetMicroPrecision(confusionMat), ShouldAlmostEqual, 0.85, 0.01)
		})

		Convey("MicroRecall", func() {
			So(GetMicroRecall(confusionMat), ShouldAlmostEqual, 0.85, 0.01)
		})

		Convey("MacroPrecision", func() {
			So(GetMacroPrecision(confusionMat), ShouldAlmostEqual, 0.775, 0.01)
		})

		Convey("MacroRecall", func() {
			So(GetMacroRecall(confusionMat), ShouldAlmostEqual, 0.719, 0.01)
		})

		Convey("F1Score", func() {
			So(GetF1Score("a", confusionMat), ShouldAlmostEqual, 0.91, 0.1)
			So(GetF1Score("b", confusionMat), ShouldAlmostEqual, 0.571, 0.01)
		})

		Convey("Accuracy", func() {
			So(GetAccuracy(confusionMat), ShouldAlmostEqual, 0.85, 0.1)
		})

		Convey("Get Summary", func() {
			output := GetSummary(confusionMat)
			So(output, ShouldStartWith, "Reference Class")
			So(output, ShouldContainSubstring, "True Positives")
			So(output, ShouldContainSubstring, "False Positives")
			So(output, ShouldContainSubstring, "True Negatives")
			So(output, ShouldContainSubstring, "Precision")
			So(output, ShouldContainSubstring, "Recall")
			So(output, ShouldContainSubstring, "F1 Score")
			So(output, ShouldContainSubstring, "------")
			So(output, ShouldContainSubstring, "Overall accuracy:")
		})

		Convey("Show Confusion Matrix", func() {
			output := ShowConfusionMatrix(confusionMat)
			So(output, ShouldStartWith, "Reference Class")
			So(output, ShouldContainSubstring, "---------------")
		})

		Convey("Get Confusion Matrix", func() {
			X, _ := base.ParseCSVToInstances("../examples/datasets/iris_headers.csv", true)
			Y, _ := base.ParseCSVToInstances("../examples/datasets/exam.csv", true)
			Convey("Nomarl ref and gen matrices", func() {
				out, _ := GetConfusionMatrix(X, X)
				ret := make(map[string]map[string]int)
				So(out, ShouldHaveSameTypeAs, ret)
			})
			Convey("Row count mismatch", func() {
				_, err := GetConfusionMatrix(X, Y)
				So(err.Error(), ShouldStartWith, "Row count mismatch:")
			})
		})
	})
}
