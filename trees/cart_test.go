package trees

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestRegressor(t *testing.T) {

	Convey("Doing a CART Test", t, func() {
		// For Classification Trees:

		// Is Gini being calculated correctly
		gini, giniMaxLabel := computeGiniImpurityAndModeLabel([]int64{1, 0, 0, 1}, []int64{0, 1})
		So(gini, ShouldEqual, 0.5)
		So(giniMaxLabel, ShouldNotBeNil)

		// Is Entropy being calculated correctly
		entropy, entropyMaxLabel := computeEntropyAndModeLabel([]int64{1, 0, 0, 1}, []int64{0, 1})
		So(entropy, ShouldEqual, 1.0)
		So(entropyMaxLabel, ShouldNotBeNil)

		// Is Data being split into left and right properly
		classifierData := [][]float64{[]float64{1, 3, 6},
			[]float64{1, 2, 3},
			[]float64{1, 9, 6},
			[]float64{1, 11, 1}}

		classifiery := []int64{0, 1, 0, 0}

		leftdata, rightdata, lefty, righty := classifierCreateSplit(classifierData, 1, classifiery, 5.0)

		So(len(leftdata), ShouldEqual, 2)
		So(len(lefty), ShouldEqual, 2)
		So(len(rightdata), ShouldEqual, 2)
		So(len(righty), ShouldEqual, 2)

		// Is isolating unique values working properly
		So(len(findUnique([]float64{10, 1, 1})), ShouldEqual, 2)

		// is data reordered correctly
		orderedData, orderedY := classifierReOrderData(getFeature(classifierData, 1), classifierData, classifiery)

		So(orderedData[1][1], ShouldEqual, 3.0)
		So(orderedY[0], ShouldEqual, 1)

		// Is split being updated properly based on threshold
		leftdata, lefty, rightdata, righty = classifierUpdateSplit(leftdata, lefty, rightdata, righty, 1, 9.5)
		So(len(leftdata), ShouldEqual, 3)
		So(len(rightdata), ShouldEqual, 1)

		// Is the root Node null when tree is not trained?
		tree := NewDecisionTreeClassifier("gini", -1, []int64{0, 1})
		So(tree.RootNode, ShouldBeNil)
		So(tree.triedSplits, ShouldBeEmpty)

		// ------------------------------------------
		// For Regression Trees

		// Is MAE being calculated correctly
		mae, maeMaxLabel := computeMaeImpurityAndAverage([]float64{1, 3, 5})
		So(mae, ShouldEqual, (4.0 / 3.0))
		So(maeMaxLabel, ShouldNotBeNil)

		// Is Entropy being calculated correctly
		mse, mseMaxLabel := computeMseImpurityAndAverage([]float64{1, 3, 5})
		So(mse, ShouldEqual, (8.0 / 3.0))
		So(mseMaxLabel, ShouldNotBeNil)

		// Is Data being split into left and right properly
		data := [][]float64{[]float64{1, 3, 6},
			[]float64{1, 2, 3},
			[]float64{1, 9, 6},
			[]float64{1, 11, 1}}

		y := []float64{1, 2, 3, 4}

		leftData, rightData, leftY, rightY := regressorCreateSplit(data, 1, y, 5.0)

		So(len(leftData), ShouldEqual, 2)
		So(len(leftY), ShouldEqual, 2)
		So(len(rightData), ShouldEqual, 2)
		So(len(rightY), ShouldEqual, 2)

		// is data reordered correctly
		regressorOrderedData, regressorOrderedY := regressorReOrderData(getFeature(data, 1), data, y)

		So(regressorOrderedData[1][1], ShouldEqual, 3.0)
		So(regressorOrderedY[0], ShouldEqual, 2)

		// Is split being updated properly based on threshold
		leftData, leftY, rightData, rightY = regressorUpdateSplit(leftData, leftY, rightData, rightY, 1, 9.5)
		So(len(leftData), ShouldEqual, 3)
		So(len(rightData), ShouldEqual, 1)

		// Is the root Node null when tree is not trained?
		regressorTreetree := NewDecisionTreeRegressor("mae", -1)
		So(regressorTreetree.RootNode, ShouldBeNil)
		So(regressorTreetree.triedSplits, ShouldBeEmpty)

	})

}
