package trees

import (
	"fmt"
	"math"
	"sort"
	"strings"

	"github.com/sjwhitworth/golearn/base"
)

// The "c" prefix to function names indicates that they were tailored for classification

// CNode is Node struct for Decision Tree Classifier
type CNode struct {
	Left       *CNode
	Right      *CNode
	Threshold  float64
	Feature    int64
	LeftLabel  int64
	RightLabel int64
	Use_not    bool
	maxDepth   int64
}

// CTree: Tree struct for Decision Tree Classifier
type CTree struct {
	RootNode    *CNode
	criterion   string
	maxDepth    int64
	labels      []int64
	triedSplits [][]float64
}

// Calculate Gini Impurity of Target Labels
func giniImpurity(y []int64, labels []int64) (float64, int64) {
	nInstances := len(y)
	gini := 0.0
	maxLabelCount := 0
	var maxLabel int64 = 0
	for label := range labels {
		numLabel := 0
		for target := range y {
			if y[target] == labels[label] {
				numLabel++
			}
		}
		p := float64(numLabel) / float64(nInstances)
		gini += p * (1 - p)
		if numLabel > maxLabelCount {
			maxLabel = labels[label]
			maxLabelCount = numLabel
		}
	}
	return gini, maxLabel
}

// Calculate Entropy loss of Target Labels
func entropy(y []int64, labels []int64) (float64, int64) {
	nInstances := len(y)
	entropy := 0.0
	maxLabelCount := 0
	var maxLabel int64 = 0
	for label := range labels {
		numLabel := 0
		for target := range y {
			if y[target] == labels[label] {
				numLabel++
			}
		}
		p := float64(numLabel) / float64(nInstances)

		logP := math.Log2(p)
		if p == 0 {
			logP = 0
		}
		entropy += -p * logP
		if numLabel > maxLabelCount {
			maxLabel = labels[label]
			maxLabelCount = numLabel
		}
	}
	return entropy, maxLabel
}

// Split the data into left node and right node based on feature and threshold - only needed for fresh nodes
func ctestSplit(data [][]float64, feature int64, y []int64, threshold float64) ([][]float64, [][]float64, []int64, []int64) {
	var left [][]float64
	var right [][]float64
	var lefty []int64
	var righty []int64

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

// Helper Function to check if data point is unique or not
func cstringInSlice(a float64, list []float64) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

// Isolate only unique values. Needed for splitting data.
func cfindUnique(data []float64) []float64 {
	var unique []float64
	for i := range data {
		if !cstringInSlice(data[i], unique) {
			unique = append(unique, data[i])
		}
	}
	return unique
}

// Isolate only the feature being considered for splitting
func cgetFeature(data [][]float64, feature int64) []float64 {
	var featureVals []float64
	for i := range data {
		featureVals = append(featureVals, data[i][feature])
	}
	return featureVals
}

// Function to Create New Decision Tree Classifier
func NewDecisionTreeClassifier(criterion string, maxDepth int64, labels []int64) *CTree {
	var tree CTree
	tree.criterion = strings.ToLower(criterion)
	tree.maxDepth = maxDepth
	tree.labels = labels

	return &tree
}

// Make sure that split being considered has not been done before
func cvalidate(triedSplits [][]float64, feature int64, threshold float64) bool {
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
type cSlice struct {
	sort.Float64Slice
	Idx []int
}

// Helper function for re-ordering data
func (s cSlice) cSwap(i, j int) {
	s.Float64Slice.Swap(i, j)
	s.Idx[i], s.Idx[j] = s.Idx[j], s.Idx[i]
}

// Final Helper Function for re-ordering data
func cNewSlice(n []float64) *cSlice {
	s := &cSlice{Float64Slice: sort.Float64Slice(n), Idx: make([]int, len(n))}

	for i := range s.Idx {
		s.Idx[i] = i
	}
	return s
}

// Reorder the data by feature being considered. Optimizes code by reducing the number of times we have to loop over data for splitting
func creOrderData(featureVal []float64, data [][]float64, y []int64) ([][]float64, []int64) {
	s := cNewSlice(featureVal)
	sort.Sort(s)

	indexes := s.Idx

	var dataSorted [][]float64
	var ySorted []int64

	for _, index := range indexes {
		dataSorted = append(dataSorted, data[index])
		ySorted = append(ySorted, y[index])
	}

	return dataSorted, ySorted
}

// Change data in Left Node and Right Node based on change in threshold
func cupdateSplit(left [][]float64, lefty []int64, right [][]float64, righty []int64, feature int64, threshold float64) ([][]float64, []int64, [][]float64, []int64) {

	for right[0][feature] < threshold {
		left = append(left, right[0])
		right = right[1:]
		lefty = append(lefty, righty[0])
		righty = righty[1:]
	}

	return left, lefty, right, righty
}

// Fit - Method visible to user to train tree
func (tree *CTree) Fit(X base.FixedDataGrid) {
	var emptyNode CNode

	data := classifierConvertInstancesToProblemVec(X)
	y := classifierConvertInstancesToLabelVec(X)
	emptyNode = cbestSplit(*tree, data, y, tree.labels, emptyNode, tree.criterion, tree.maxDepth, 0)

	tree.RootNode = &emptyNode
}

// Iterativly find and record the best split - recursive function
func cbestSplit(tree CTree, data [][]float64, y []int64, labels []int64, upperNode CNode, criterion string, maxDepth int64, depth int64) CNode {

	// Ensure that we have not reached maxDepth. maxDepth =-1 means split until nodes are pure
	depth++

	if maxDepth != -1 && depth > maxDepth {
		return upperNode
	}

	numFeatures := len(data[0])
	var bestGini float64
	var origGini float64

	// Calculate loss based on Criterion Specified by user
	if criterion == "gini" {
		origGini, upperNode.LeftLabel = giniImpurity(y, labels)
	} else if criterion == "entropy" {
		origGini, upperNode.LeftLabel = entropy(y, labels)
	} else {
		panic("Invalid impurity function, choose from GINI or ENTROPY")
	}

	bestGini = origGini

	bestLeft := data
	bestRight := data
	bestLefty := y
	bestRighty := y

	numData := len(data)

	bestLeftGini := bestGini
	bestRightGini := bestGini

	upperNode.Use_not = true

	var leftN CNode
	var rightN CNode
	// Iterate over all features
	for i := 0; i < numFeatures; i++ {
		featureVal := cgetFeature(data, int64(i))
		unique := cfindUnique(featureVal)
		sort.Float64s(unique)
		numUnique := len(unique)

		sortData, sortY := creOrderData(featureVal, data, y)

		firstTime := true

		var left, right [][]float64
		var lefty, righty []int64
		// Iterate over all possible thresholds for that feature
		for j := range unique {
			if j != (numUnique - 1) {
				threshold := (unique[j] + unique[j+1]) / 2
				// Ensure that same split has not been made before
				if cvalidate(tree.triedSplits, int64(i), threshold) {
					// We need to split data from fresh when considering new feature for the first time.
					// Otherwise, we need to update the split by moving data points from left to right.
					if firstTime {
						left, right, lefty, righty = ctestSplit(sortData, int64(i), sortY, threshold)
						firstTime = false
					} else {
						left, lefty, right, righty = cupdateSplit(left, lefty, right, righty, int64(i), threshold)
					}

					var leftGini float64
					var rightGini float64
					var leftLabels int64
					var rightLabels int64

					if criterion == "gini" {
						leftGini, leftLabels = giniImpurity(lefty, labels)
						rightGini, rightLabels = giniImpurity(righty, labels)
					} else if criterion == "entropy" {
						leftGini, leftLabels = entropy(lefty, labels)
						rightGini, rightLabels = entropy(righty, labels)
					}
					// Calculate weighted gini impurity of child nodes
					subGini := (leftGini * float64(len(left)) / float64(numData)) + (rightGini * float64(len(right)) / float64(numData))

					// If we find a split that reduces impurity
					if subGini < bestGini {
						bestGini = subGini
						bestLeft = left
						bestRight = right
						bestLefty = lefty
						bestRighty = righty
						upperNode.Threshold = threshold
						upperNode.Feature = int64(i)

						upperNode.LeftLabel = leftLabels
						upperNode.RightLabel = rightLabels

						bestLeftGini = leftGini
						bestRightGini = rightGini
					}
				}

			}
		}
	}
	// If no split was found, we don't want to use this node, so we will flag it
	if bestGini == origGini {
		upperNode.Use_not = false
		return upperNode
	}
	// Until nodes are not pure
	if bestGini > 0 {

		// If left node is pure, no need to split on left side again
		if bestLeftGini > 0 {
			tree.triedSplits = append(tree.triedSplits, []float64{float64(upperNode.Feature), upperNode.Threshold})
			// Recursive splitting logic
			leftN = cbestSplit(tree, bestLeft, bestLefty, labels, leftN, criterion, maxDepth, depth)
			if leftN.Use_not == true {
				upperNode.Left = &leftN
			}

		}
		// If right node is pure, no need to split on right side again
		if bestRightGini > 0 {
			tree.triedSplits = append(tree.triedSplits, []float64{float64(upperNode.Feature), upperNode.Threshold})
			// Recursive splitting logic
			rightN = cbestSplit(tree, bestRight, bestRighty, labels, rightN, criterion, maxDepth, depth)
			if rightN.Use_not == true {
				upperNode.Right = &rightN
			}

		}

	}
	// Return the node - contains all information regarding feature and threshold.
	return upperNode
}

// PrintTree : this function prints out entire tree for visualization - visible to user
func (tree *CTree) PrintTree() {
	rootNode := *tree.RootNode
	cprintTreeFromNode(rootNode, "")
}

// Tree struct has root node. That is used to print tree - invisible to user but called from PrintTree
func cprintTreeFromNode(tree CNode, spacing string) float64 {

	fmt.Print(spacing + "Feature ")
	fmt.Print(tree.Feature)
	fmt.Print(" < ")
	fmt.Println(tree.Threshold)

	if tree.Left == nil {
		fmt.Println(spacing + "---> True")
		fmt.Print("  " + spacing + "PREDICT    ")
		fmt.Println(tree.LeftLabel)
	}
	if tree.Right == nil {
		fmt.Println(spacing + "---> FALSE")
		fmt.Print("  " + spacing + "PREDICT    ")
		fmt.Println(tree.RightLabel)
	}

	if tree.Left != nil {
		fmt.Println(spacing + "---> True")
		cprintTreeFromNode(*tree.Left, spacing+"  ")
	}

	if tree.Right != nil {
		fmt.Println(spacing + "---> False")
		cprintTreeFromNode(*tree.Right, spacing+"  ")
	}

	return 0.0
}

// Predict a single data point by traversing the entire tree
func cpredictSingle(tree CNode, instance []float64) int64 {
	if instance[tree.Feature] < tree.Threshold {
		if tree.Left == nil {
			return tree.LeftLabel
		} else {
			return cpredictSingle(*tree.Left, instance)
		}
	} else {
		if tree.Right == nil {
			return tree.RightLabel
		} else {
			return cpredictSingle(*tree.Right, instance)
		}
	}
}

// Predict is visible to user. Given test data, they receive predictions for every datapoint.
func (tree *CTree) Predict(X_test base.FixedDataGrid) []int64 {
	root := *tree.RootNode
	test := classifierConvertInstancesToProblemVec(X_test)
	return cpredictFromNode(root, test)
}

// This function uses the rootnode from Predict. It is invisible to user, but called from predict method.
func cpredictFromNode(tree CNode, test [][]float64) []int64 {
	var preds []int64
	for i := range test {
		iPred := cpredictSingle(tree, test[i])
		preds = append(preds, iPred)
	}
	return preds
}

// Given Test data and label, return the accuracy of the classifier. Data has to be in float slice format before feeding.
func (tree *CTree) Evaluate(test base.FixedDataGrid) float64 {
	rootNode := *tree.RootNode
	xTest := classifierConvertInstancesToProblemVec(test)
	yTest := classifierConvertInstancesToLabelVec(test)
	return cevaluateFromNode(rootNode, xTest, yTest)
}

func cevaluateFromNode(tree CNode, xTest [][]float64, yTest []int64) float64 {
	preds := cpredictFromNode(tree, xTest)
	accuracy := 0.0
	for i := range preds {
		if preds[i] == yTest[i] {
			accuracy++
		}
	}
	accuracy /= float64(len(yTest))
	return accuracy
}

// Helper function to convert base.FixedDataGrid into required format. Called in Fit, Predict
func classifierConvertInstancesToProblemVec(X base.FixedDataGrid) [][]float64 {
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

// Helper function to convert base.FixedDataGrid into required format. Called in Fit, Predict
func classifierConvertInstancesToLabelVec(X base.FixedDataGrid) []int64 {
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
	labelVec := make([]int64, rows)
	// Resolve class Attribute specification
	classAttrSpecs := base.ResolveAttributes(X, classAttrs)
	X.MapOverRows(classAttrSpecs, func(row [][]byte, rowNo int) (bool, error) {
		labelVec[rowNo] = int64(base.UnpackBytesToFloat(row[0]))
		return true, nil
	})
	return labelVec
}
