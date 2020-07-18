package trees

import (
	"fmt"
	"math"
	"sort"
	"strings"

	"github.com/sjwhitworth/golearn/base"
)

// The "r" prefix to all function names indicates that they were tailored to support regression.

// See cart_classifier for details on functions.
type RNode struct {
	Left      *RNode
	Right     *RNode
	Threshold float64
	Feature   int64
	LeftPred  float64
	RightPred float64
	Use_not   bool
}

type RTree struct {
	RootNode    *RNode
	criterion   string
	maxDepth    int64
	triedSplits [][]float64
}

func meanAbsoluteError(y []float64, yBar float64) float64 {
	error := 0.0
	for _, target := range y {
		error += math.Abs(target - yBar)
	}
	error /= float64(len(y))
	return error
}

func average(y []float64) float64 {
	mean := 0.0
	for _, value := range y {
		mean += value
	}
	mean /= float64(len(y))
	return mean
}

func maeImpurity(y []float64) (float64, float64) {
	yHat := average(y)
	return meanAbsoluteError(y, yHat), yHat
}

func meanSquaredError(y []float64, yBar float64) float64 {
	error := 0.0
	for _, target := range y {
		item_error := target - yBar
		error += math.Pow(item_error, 2)
	}
	error /= float64(len(y))
	return error
}

func mseImpurity(y []float64) (float64, float64) {
	yHat := average(y)
	return meanSquaredError(y, yHat), yHat
}

func rtestSplit(data [][]float64, feature int64, y []float64, threshold float64) ([][]float64, [][]float64, []float64, []float64) {
	var left [][]float64
	var lefty []float64
	var right [][]float64
	var righty []float64

	for i := range data {
		example := data[i]
		if example[feature] < threshold {
			left = append(left, example)
			lefty = append(lefty, y[i])
		} else {
			right = append(right, example)
			righty = append(righty, y[i])
		}
	}

	return left, right, lefty, righty
}

func rstringInSlice(a float64, list []float64) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

func rfindUnique(data []float64) []float64 {
	var unique []float64
	for i := range data {
		if !rstringInSlice(data[i], unique) {
			unique = append(unique, data[i])
		}
	}
	return unique
}

func rgetFeature(data [][]float64, feature int64) []float64 {
	var featureVals []float64
	for i := range data {
		featureVals = append(featureVals, data[i][feature])
	}
	return featureVals
}

func NewDecisionTreeRegressor(criterion string, maxDepth int64) *RTree {
	var tree RTree
	tree.maxDepth = maxDepth
	tree.criterion = strings.ToLower(criterion)
	return &tree
}

func rvalidate(triedSplits [][]float64, feature int64, threshold float64) bool {
	for i := range triedSplits {
		split := triedSplits[i]
		featureTried, thresholdTried := split[0], split[1]
		if int64(featureTried) == feature && thresholdTried == threshold {
			return false
		}
	}
	return true
}

// Helper struct for re-rdering data
type rSlice struct {
	sort.Float64Slice
	Idx []int
}

// Helper function for re-ordering data
func (s rSlice) rSwap(i, j int) {
	s.Float64Slice.Swap(i, j)
	s.Idx[i], s.Idx[j] = s.Idx[j], s.Idx[i]
}

// Final Helper Function for re-ordering data
func rNewSlice(n []float64) *rSlice {
	s := &rSlice{Float64Slice: sort.Float64Slice(n), Idx: make([]int, len(n))}

	for i := range s.Idx {
		s.Idx[i] = i
	}
	return s
}

func rreOrderData(featureVal []float64, data [][]float64, y []float64) ([][]float64, []float64) {
	s := rNewSlice(featureVal)
	sort.Sort(s)

	indexes := s.Idx

	var dataSorted [][]float64
	var ySorted []float64

	for _, index := range indexes {
		dataSorted = append(dataSorted, data[index])
		ySorted = append(ySorted, y[index])
	}

	return dataSorted, ySorted

}

func rupdateSplit(left [][]float64, lefty []float64, right [][]float64, righty []float64, feature int64, threshold float64) ([][]float64, []float64, [][]float64, []float64) {

	for right[0][feature] < threshold {
		left = append(left, right[0])
		right = right[1:]
		lefty = append(lefty, righty[0])
		righty = righty[1:]
	}

	return left, lefty, right, righty
}

func sum(y []int64) int64 {
	var sum_ int64 = 0
	for i := range y {
		sum_ += y[i]
	}
	return sum_
}

// Extra Method for creating simple to use interface. Many params are either redundant for user but are needed only for recursive logic.
func (tree *RTree) Fit(X base.FixedDataGrid) {
	var emptyNode RNode
	data := regressorConvertInstancesToProblemVec(X)
	y := regressorConvertInstancesToLabelVec(X)

	emptyNode = rbestSplit(*tree, data, y, emptyNode, tree.criterion, tree.maxDepth, 0)

	tree.RootNode = &emptyNode
}

// Essentially the Fit Method
func rbestSplit(tree RTree, data [][]float64, y []float64, upperNode RNode, criterion string, maxDepth int64, depth int64) RNode {

	depth++

	if depth > maxDepth && maxDepth != -1 {
		return upperNode
	}

	numFeatures := len(data[0])
	var bestLoss float64
	var origLoss float64

	if criterion == "mae" {
		origLoss, upperNode.LeftPred = maeImpurity(y)
	} else {
		origLoss, upperNode.LeftPred = mseImpurity(y)
	}

	bestLoss = origLoss

	bestLeft := data
	bestRight := data
	bestLefty := y
	bestRighty := y

	numData := len(data)

	bestLeftLoss := bestLoss
	bestRightLoss := bestLoss

	upperNode.Use_not = true

	var leftN RNode
	var rightN RNode
	// Iterate over all features
	for i := 0; i < numFeatures; i++ {
		featureVal := rgetFeature(data, int64(i))
		unique := rfindUnique(featureVal)
		sort.Float64s(unique)
		numUnique := len(unique)

		sortData, sortY := rreOrderData(featureVal, data, y)

		firstTime := true

		var left, right [][]float64
		var lefty, righty []float64

		for j := range unique {
			if j != (numUnique - 1) {
				threshold := (unique[j] + unique[j+1]) / 2
				if rvalidate(tree.triedSplits, int64(i), threshold) {
					if firstTime {
						left, right, lefty, righty = rtestSplit(sortData, int64(i), sortY, threshold)
						firstTime = false
					} else {
						left, lefty, right, righty = rupdateSplit(left, lefty, right, righty, int64(i), threshold)
					}

					var leftLoss float64
					var rightLoss float64
					var leftPred float64
					var rightPred float64

					if criterion == "mae" {
						leftLoss, leftPred = maeImpurity(lefty)
						rightLoss, rightPred = maeImpurity(righty)
					} else {
						leftLoss, leftPred = mseImpurity(lefty)
						rightLoss, rightPred = mseImpurity(righty)
					}

					subLoss := (leftLoss * float64(len(left)) / float64(numData)) + (rightLoss * float64(len(right)) / float64(numData))

					if subLoss < bestLoss {
						bestLoss = subLoss
						bestLeft = left
						bestRight = right
						bestLefty = lefty
						bestRighty = righty
						upperNode.Threshold = threshold
						upperNode.Feature = int64(i)

						upperNode.LeftPred = leftPred
						upperNode.RightPred = rightPred

						bestLeftLoss = leftLoss
						bestRightLoss = rightLoss
					}
				}

			}
		}
	}

	if bestLoss == origLoss {
		upperNode.Use_not = false
		return upperNode
	}

	if bestLoss > 0 {

		if bestLeftLoss > 0 {
			tree.triedSplits = append(tree.triedSplits, []float64{float64(upperNode.Feature), upperNode.Threshold})
			leftN = rbestSplit(tree, bestLeft, bestLefty, leftN, criterion, maxDepth, depth)
			if leftN.Use_not == true {
				upperNode.Left = &leftN
			}

		}
		if bestRightLoss > 0 {
			tree.triedSplits = append(tree.triedSplits, []float64{float64(upperNode.Feature), upperNode.Threshold})
			rightN = rbestSplit(tree, bestRight, bestRighty, rightN, criterion, maxDepth, depth)
			if rightN.Use_not == true {
				upperNode.Right = &rightN
			}

		}

	}

	return upperNode
}

func (tree *RTree) PrintTree() {
	rootNode := *tree.RootNode
	printTreeFromNode(rootNode, "")
}

func printTreeFromNode(tree RNode, spacing string) float64 {

	fmt.Print(spacing + "Feature ")
	fmt.Print(tree.Feature)
	fmt.Print(" < ")
	fmt.Println(tree.Threshold)

	if tree.Left == nil {
		fmt.Println(spacing + "---> True")
		fmt.Print("  " + spacing + "PREDICT    ")
		fmt.Println(tree.LeftPred)
	}
	if tree.Right == nil {
		fmt.Println(spacing + "---> FALSE")
		fmt.Print("  " + spacing + "PREDICT    ")
		fmt.Println(tree.RightPred)
	}

	if tree.Left != nil {
		fmt.Println(spacing + "---> True")
		printTreeFromNode(*tree.Left, spacing+"  ")
	}

	if tree.Right != nil {
		fmt.Println(spacing + "---> False")
		printTreeFromNode(*tree.Right, spacing+"  ")
	}

	return 0.0
}

func predictSingle(tree RNode, instance []float64) float64 {
	if instance[tree.Feature] < tree.Threshold {
		if tree.Left == nil {
			return tree.LeftPred
		} else {
			return predictSingle(*tree.Left, instance)
		}
	} else {
		if tree.Right == nil {
			return tree.RightPred
		} else {
			return predictSingle(*tree.Right, instance)
		}
	}
}

func (tree *RTree) Predict(X_test base.FixedDataGrid) []float64 {
	root := *tree.RootNode
	test := regressorConvertInstancesToProblemVec(X_test)
	return predictFromNode(root, test)
}

func predictFromNode(tree RNode, test [][]float64) []float64 {
	var preds []float64
	for i := range test {
		i_pred := predictSingle(tree, test[i])
		preds = append(preds, i_pred)
	}
	return preds
}

// Helper function to convert base.FixedDataGrid into required format. Called in Fit
func regressorConvertInstancesToProblemVec(X base.FixedDataGrid) [][]float64 {
	// Allocate problem array
	_, rows := X.Size()
	problemVec := make([][]float64, rows)

	// Retrieve numeric non-class Attributes
	numericAttrs := base.NonClassFloatAttributes(X)
	numericAttrSpecs := base.ResolveAttributes(X, numericAttrs)

	// Convert each row
	X.MapOverRows(numericAttrSpecs, func(row [][]byte, rowNo int) (bool, error) {
		// Allocate a new row
		probRow := make([]float64, len(numericAttrSpecs))
		// Read out the row
		for i, _ := range numericAttrSpecs {
			probRow[i] = base.UnpackBytesToFloat(row[i])
		}
		// Add the row
		problemVec[rowNo] = probRow
		return true, nil
	})
	return problemVec
}

// Helper function to convert base.FixedDataGrid into required format. Called in Fit
func regressorConvertInstancesToLabelVec(X base.FixedDataGrid) []float64 {
	// Get the class Attributes
	classAttrs := X.AllClassAttributes()
	// Only support 1 class Attribute
	if len(classAttrs) != 1 {
		panic(fmt.Sprintf("%d ClassAttributes (1 expected)", len(classAttrs)))
	}
	// ClassAttribute must be numeric
	if _, ok := classAttrs[0].(*base.FloatAttribute); !ok {
		panic(fmt.Sprintf("%s: ClassAttribute must be a FloatAttribute", classAttrs[0]))
	}
	// Allocate return structure
	_, rows := X.Size()
	// labelVec := make([]float64, rows)
	labelVec := make([]float64, rows)
	// Resolve class Attribute specification
	classAttrSpecs := base.ResolveAttributes(X, classAttrs)
	X.MapOverRows(classAttrSpecs, func(row [][]byte, rowNo int) (bool, error) {
		labelVec[rowNo] = base.UnpackBytesToFloat(row[0])
		return true, nil
	})
	return labelVec
}
