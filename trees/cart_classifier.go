package trees

import (
	"fmt"
	"math"
	"sort"
	"strconv"
	"strings"

	"github.com/sjwhitworth/golearn/base"
)

const (
	GINI    string = "gini"
	ENTROPY string = "entropy"
)

// CNode is Node struct for Decision Tree Classifier.
// It holds the information for each split (which feature to use, what threshold, and which label to assign for each side of the split)
type classifierNode struct {
	Left       *classifierNode
	Right      *classifierNode
	Threshold  float64
	Feature    int64
	LeftLabel  int64
	RightLabel int64
	Use_not    bool
	maxDepth   int64
}

// CARTDecisionTreeClassifier: Tree struct for Decision Tree Classifier
// It contains the rootNode, as well as all of the hyperparameters chosen by the user.
// It also keeps track of all splits done at the tree level.
type CARTDecisionTreeClassifier struct {
	RootNode    *classifierNode
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

// Split the data into left node and right node based on feature and threshold
func classifierCreateSplit(data [][]float64, feature int64, y []int64, threshold float64) ([][]float64, [][]float64, []int64, []int64) {
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

// Helper Function to check if data point is unique or not.
// We will use this to isolate unique values of a feature
func classifierStringInSlice(a float64, list []float64) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

// Isolate only unique values. This way, we can try only unique splits and not redundant ones.
func classifierFindUnique(data []float64) []float64 {
	var unique []float64
	for i := range data {
		if !classifierStringInSlice(data[i], unique) {
			unique = append(unique, data[i])
		}
	}
	return unique
}

// Isolate only the feature being considered for splitting. Reduces the complexity in managing splits.
func classifierGetFeature(data [][]float64, feature int64) []float64 {
	var featureVals []float64
	for i := range data {
		featureVals = append(featureVals, data[i][feature])
	}
	return featureVals
}

// Function to Create New Decision Tree Classifier.
// It assigns all of the hyperparameters by user into the tree attributes.
func NewDecisionTreeClassifier(criterion string, maxDepth int64, labels []int64) *CARTDecisionTreeClassifier {
	var tree CARTDecisionTreeClassifier
	tree.criterion = strings.ToLower(criterion)
	tree.maxDepth = maxDepth
	tree.labels = labels

	return &tree
}

// Make sure that split being considered has not been done before.
// Else we will unnecessarily try splits that won't improve Impurity.
func classifierValidate(triedSplits [][]float64, feature int64, threshold float64) bool {
	for i := range triedSplits {
		split := triedSplits[i]
		featureTried, thresholdTried := split[0], split[1]
		if int64(featureTried) == feature && thresholdTried == threshold {
			return false
		}
	}
	return true
}

// Reorder the data by feature being considered. Optimizes code by reducing the number of times we have to loop over data for splitting
func classifierReOrderData(featureVal []float64, data [][]float64, y []int64) ([][]float64, []int64) {
	s := NewSlice(featureVal)
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

// Update the left and right side of the split based on the threshold.
func classifierUpdateSplit(left [][]float64, lefty []int64, right [][]float64, righty []int64, feature int64, threshold float64) ([][]float64, []int64, [][]float64, []int64) {

	for right[0][feature] < threshold {
		left = append(left, right[0])
		right = right[1:]
		lefty = append(lefty, righty[0])
		righty = righty[1:]
	}

	return left, lefty, right, righty
}

// Fit - Creates an Emppty Root Node
// Trains the tree by calling recursive function classifierBestSplit
func (tree *CARTDecisionTreeClassifier) Fit(X base.FixedDataGrid) {
	var emptyNode classifierNode

	data := classifierConvertInstancesToProblemVec(X)
	y := classifierConvertInstancesToLabelVec(X)
	emptyNode = classifierBestSplit(*tree, data, y, tree.labels, emptyNode, tree.criterion, tree.maxDepth, 0)

	tree.RootNode = &emptyNode
}

// Iterativly find and record the best split
// Stop If depth reaches maxDepth or nodes are pure
func classifierBestSplit(tree CARTDecisionTreeClassifier, data [][]float64, y []int64, labels []int64, upperNode classifierNode, criterion string, maxDepth int64, depth int64) classifierNode {

	// Ensure that we have not reached maxDepth. maxDepth =-1 means split until nodes are pure
	depth++

	if maxDepth != -1 && depth > maxDepth {
		return upperNode
	}

	numFeatures := len(data[0])
	var bestGini float64
	var origGini float64

	// Calculate loss based on Criterion Specified by user
	if criterion == GINI {
		origGini, upperNode.LeftLabel = giniImpurity(y, labels)
	} else if criterion == ENTROPY {
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

	var leftN classifierNode
	var rightN classifierNode
	// Iterate over all features
	for i := 0; i < numFeatures; i++ {
		featureVal := classifierGetFeature(data, int64(i))
		unique := classifierFindUnique(featureVal)
		sort.Float64s(unique)
		numUnique := len(unique)

		sortData, sortY := classifierReOrderData(featureVal, data, y)

		firstTime := true

		var left, right [][]float64
		var lefty, righty []int64
		// Iterate over all possible thresholds for that feature
		for j := range unique {
			if j != (numUnique - 1) {
				threshold := (unique[j] + unique[j+1]) / 2
				// Ensure that same split has not been made before
				if classifierValidate(tree.triedSplits, int64(i), threshold) {
					// We need to split data from fresh when considering new feature for the first time.
					// Otherwise, we need to update the split by moving data points from left to right.
					if firstTime {
						left, right, lefty, righty = classifierCreateSplit(sortData, int64(i), sortY, threshold)
						firstTime = false
					} else {
						left, lefty, right, righty = classifierUpdateSplit(left, lefty, right, righty, int64(i), threshold)
					}

					var leftGini float64
					var rightGini float64
					var leftLabels int64
					var rightLabels int64

					if criterion == GINI {
						leftGini, leftLabels = giniImpurity(lefty, labels)
						rightGini, rightLabels = giniImpurity(righty, labels)
					} else if criterion == ENTROPY {
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
			leftN = classifierBestSplit(tree, bestLeft, bestLefty, labels, leftN, criterion, maxDepth, depth)
			if leftN.Use_not == true {
				upperNode.Left = &leftN
			}

		}
		// If right node is pure, no need to split on right side again
		if bestRightGini > 0 {
			tree.triedSplits = append(tree.triedSplits, []float64{float64(upperNode.Feature), upperNode.Threshold})
			// Recursive splitting logic
			rightN = classifierBestSplit(tree, bestRight, bestRighty, labels, rightN, criterion, maxDepth, depth)
			if rightN.Use_not == true {
				upperNode.Right = &rightN
			}

		}

	}
	// Return the node - contains all information regarding feature and threshold.
	return upperNode
}

// String : this function prints out entire tree for visualization.
// Calls a recursive function to print the tree - classifierPrintTreeFromNode
func (tree *CARTDecisionTreeClassifier) String() string {
	rootNode := *tree.RootNode
	return classifierPrintTreeFromNode(rootNode, "")
}

func classifierPrintTreeFromNode(tree classifierNode, spacing string) string {
	returnString := ""
	returnString += spacing + "Feature "
	returnString += strconv.FormatInt(tree.Feature, 10)
	returnString += " < "
	returnString += fmt.Sprintf("%.3f", tree.Threshold)
	returnString += "\n"

	if tree.Left == nil {
		returnString += spacing + "---> True" + "\n"
		returnString += "  " + spacing + "PREDICT    "
		returnString += strconv.FormatInt(tree.LeftLabel, 10) + "\n"

	}
	if tree.Right == nil {

		returnString += spacing + "---> False" + "\n"
		returnString += "  " + spacing + "PREDICT    "
		returnString += strconv.FormatInt(tree.RightLabel, 10) + "\n"
	}

	if tree.Left != nil {
		returnString += spacing + "---> True" + "\n"
		returnString += classifierPrintTreeFromNode(*tree.Left, spacing+"  ")
	}

	if tree.Right != nil {
		returnString += spacing + "---> False" + "\n"
		returnString += classifierPrintTreeFromNode(*tree.Right, spacing+"  ")
	}

	return returnString
}

// Predict a single data point by traversing the entire tree
// Uses recursive logic to navigate the tree.
func classifierPredictSingle(tree classifierNode, instance []float64) int64 {
	if instance[tree.Feature] < tree.Threshold {
		if tree.Left == nil {
			return tree.LeftLabel
		} else {
			return classifierPredictSingle(*tree.Left, instance)
		}
	} else {
		if tree.Right == nil {
			return tree.RightLabel
		} else {
			return classifierPredictSingle(*tree.Right, instance)
		}
	}
}

// Given test data, return predictions for every datapoint. calls classifierPredictFromNode
func (tree *CARTDecisionTreeClassifier) Predict(X_test base.FixedDataGrid) []int64 {
	root := *tree.RootNode
	test := classifierConvertInstancesToProblemVec(X_test)
	return classifierPredictFromNode(root, test)
}

// This function uses the rootnode from Predict.
// It iterates through every data point and calls the recursive function to give predictions and then summarizes them.
func classifierPredictFromNode(tree classifierNode, test [][]float64) []int64 {
	var preds []int64
	for i := range test {
		iPred := classifierPredictSingle(tree, test[i])
		preds = append(preds, iPred)
	}
	return preds
}

// Given Test data and label, return the accuracy of the classifier.
// First it retreives predictions from the data, then compares for accuracy.
// Calls classifierEvaluateFromNode
func (tree *CARTDecisionTreeClassifier) Evaluate(test base.FixedDataGrid) float64 {
	rootNode := *tree.RootNode
	xTest := classifierConvertInstancesToProblemVec(test)
	yTest := classifierConvertInstancesToLabelVec(test)
	return classifierEvaluateFromNode(rootNode, xTest, yTest)
}

// Retrieve predictions and then calculate accuracy.
func classifierEvaluateFromNode(tree classifierNode, xTest [][]float64, yTest []int64) float64 {
	preds := classifierPredictFromNode(tree, xTest)
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
