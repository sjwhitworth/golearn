package trees

import (
	"errors"
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
	Left         *classifierNode
	Right        *classifierNode
	Threshold    float64
	Feature      int64
	LeftLabel    int64
	RightLabel   int64
	isNodeNeeded bool
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

// Convert a series of labels to frequency map for efficient impurity calculation
func convertToMap(y []int64, labels []int64) map[int64]int {
	labelCount := make(map[int64]int)
	for _, label := range labels {
		labelCount[label] = 0
	}
	for _, value := range y {
		labelCount[value]++
	}
	return labelCount
}

// Calculate Gini Impurity of Target Labels
func computeGiniImpurityAndModeLabel(y []int64, labels []int64) (float64, int64) {
	nInstances := len(y)
	gini := 0.0
	var maxLabel int64 = 0

	labelCount := convertToMap(y, labels)
	for _, label := range labels {
		if labelCount[label] > labelCount[maxLabel] {
			maxLabel = label
		}
		p := float64(labelCount[label]) / float64(nInstances)
		gini += p * (1 - p)
	}
	return gini, maxLabel
}

// Calculate Entropy loss of Target Labels
func computeEntropyAndModeLabel(y []int64, labels []int64) (float64, int64) {
	nInstances := len(y)
	entropy := 0.0
	var maxLabel int64 = 0

	labelCount := convertToMap(y, labels)
	for _, label := range labels {
		if labelCount[label] > labelCount[maxLabel] {
			maxLabel = label
		}
		p := float64(labelCount[label]) / float64(nInstances)
		logP := math.Log2(p)
		if p == 0 {
			logP = 0
		}
		entropy += (-p * logP)
	}
	return entropy, maxLabel
}

func calculateClassificationLoss(y []int64, labels []int64, criterion string) (float64, int64, error) {
	if len(y) == 0 {
		return 0, 0, errors.New("Need atleast 1 value to compute impurity")
	}
	if criterion == GINI {
		loss, modeLabel := computeGiniImpurityAndModeLabel(y, labels)
		return loss, modeLabel, nil
	} else if criterion == ENTROPY {
		loss, modeLabel := computeEntropyAndModeLabel(y, labels)
		return loss, modeLabel, nil
	} else {
		return 0, 0, errors.New("Invalid impurity function, choose from GINI or ENTROPY")
	}
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

// Function to Create New Decision Tree Classifier.
// It assigns all of the hyperparameters by user into the tree attributes.
func NewDecisionTreeClassifier(criterion string, maxDepth int64, labels []int64) *CARTDecisionTreeClassifier {
	var tree CARTDecisionTreeClassifier
	tree.criterion = strings.ToLower(criterion)
	tree.maxDepth = maxDepth
	tree.labels = labels

	return &tree
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
func classifierUpdateSplit(left [][]float64, leftY []int64, right [][]float64, rightY []int64, feature int64, threshold float64) ([][]float64, []int64, [][]float64, []int64) {

	for right[0][feature] < threshold {
		left = append(left, right[0])
		right = right[1:]
		leftY = append(leftY, rightY[0])
		rightY = rightY[1:]
	}

	return left, leftY, right, rightY
}

// Fit - Creates an Empty Root Node
// Trains the tree by calling recursive function classifierBestSplit
func (tree *CARTDecisionTreeClassifier) Fit(X base.FixedDataGrid) error {
	var emptyNode classifierNode
	var err error

	data := convertInstancesToProblemVec(X)
	y, err := classifierConvertInstancesToLabelVec(X)
	if err != nil {
		return err
	}

	emptyNode, err = classifierBestSplit(*tree, data, y, tree.labels, emptyNode, tree.criterion, tree.maxDepth, 0)

	if err != nil {
		return err
	}
	tree.RootNode = &emptyNode
	return nil
}

// Iterativly find and record the best split
// Stop If depth reaches maxDepth or nodes are pure
func classifierBestSplit(tree CARTDecisionTreeClassifier, data [][]float64, y []int64, labels []int64, upperNode classifierNode, criterion string, maxDepth int64, depth int64) (classifierNode, error) {

	// Ensure that we have not reached maxDepth. maxDepth =-1 means split until nodes are pure
	depth++

	if maxDepth != -1 && depth > maxDepth {
		return upperNode, nil
	}

	numFeatures := len(data[0])
	var bestGini, origGini float64
	var err error
	// Calculate loss based on Criterion Specified by user
	origGini, upperNode.LeftLabel, err = calculateClassificationLoss(y, labels, criterion)
	if err != nil {
		return upperNode, err
	}

	bestGini = origGini

	bestLeft, bestRight, bestLefty, bestRighty := data, data, y, y

	numData := len(data)

	bestLeftGini, bestRightGini := bestGini, bestGini

	upperNode.isNodeNeeded = true

	var leftN, rightN classifierNode

	// Iterate over all features
	for i := 0; i < numFeatures; i++ {

		featureVal := getFeature(data, int64(i))
		unique := findUnique(featureVal)
		sort.Float64s(unique)

		sortData, sortY := classifierReOrderData(featureVal, data, y)

		firstTime := true

		var left, right [][]float64
		var leftY, rightY []int64
		// Iterate over all possible thresholds for that feature
		for j := 0; j < len(unique)-1; j++ {

			threshold := (unique[j] + unique[j+1]) / 2
			// Ensure that same split has not been made before
			if validate(tree.triedSplits, int64(i), threshold) {
				// We need to split data from fresh when considering new feature for the first time.
				// Otherwise, we need to update the split by moving data points from left to right.
				if firstTime {
					left, right, leftY, rightY = classifierCreateSplit(sortData, int64(i), sortY, threshold)
					firstTime = false
				} else {
					left, leftY, right, rightY = classifierUpdateSplit(left, leftY, right, rightY, int64(i), threshold)
				}

				var leftGini, rightGini float64
				var leftLabels, rightLabels int64

				leftGini, leftLabels, _ = calculateClassificationLoss(leftY, labels, criterion)
				rightGini, rightLabels, _ = calculateClassificationLoss(rightY, labels, criterion)

				// Calculate weighted gini impurity of child nodes
				subGini := (leftGini * float64(len(left)) / float64(numData)) + (rightGini * float64(len(right)) / float64(numData))

				// If we find a split that reduces impurity
				if subGini < bestGini {
					bestGini = subGini

					bestLeft, bestRight = left, right

					bestLefty, bestRighty = leftY, rightY

					upperNode.Threshold, upperNode.Feature = threshold, int64(i)

					upperNode.LeftLabel, upperNode.RightLabel = leftLabels, rightLabels

					bestLeftGini, bestRightGini = leftGini, rightGini
				}
			}
		}
	}
	// If no split was found, we don't want to use this node, so we will flag it
	if bestGini == origGini {
		upperNode.isNodeNeeded = false
		return upperNode, nil
	}
	// Until nodes are not pure
	if bestGini > 0 {

		// If left node is pure, no need to split on left side again
		if bestLeftGini > 0 {
			tree.triedSplits = append(tree.triedSplits, []float64{float64(upperNode.Feature), upperNode.Threshold})
			// Recursive splitting logic
			leftN, err = classifierBestSplit(tree, bestLeft, bestLefty, labels, leftN, criterion, maxDepth, depth)
			if err != nil {
				return upperNode, err
			}
			if leftN.isNodeNeeded == true {
				upperNode.Left = &leftN
			}

		}
		// If right node is pure, no need to split on right side again
		if bestRightGini > 0 {
			tree.triedSplits = append(tree.triedSplits, []float64{float64(upperNode.Feature), upperNode.Threshold})
			// Recursive splitting logic
			rightN, err = classifierBestSplit(tree, bestRight, bestRighty, labels, rightN, criterion, maxDepth, depth)
			if err != nil {
				return upperNode, err
			}
			if rightN.isNodeNeeded == true {
				upperNode.Right = &rightN
			}

		}

	}
	// Return the node - contains all information regarding feature and threshold.
	return upperNode, nil
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
	test := convertInstancesToProblemVec(X_test)
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
func (tree *CARTDecisionTreeClassifier) Evaluate(test base.FixedDataGrid) (float64, error) {
	rootNode := *tree.RootNode
	xTest := convertInstancesToProblemVec(test)
	yTest, err := classifierConvertInstancesToLabelVec(test)
	if err != nil {
		return 0, err
	}
	return classifierEvaluateFromNode(rootNode, xTest, yTest), nil
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
func classifierConvertInstancesToLabelVec(X base.FixedDataGrid) ([]int64, error) {
	// Get the class Attributes
	classAttrs := X.AllClassAttributes()
	// Only support 1 class Attribute
	if len(classAttrs) != 1 {
		return []int64{0}, errors.New(fmt.Sprintf("%d ClassAttributes (1 expected)", len(classAttrs)))
	}
	// ClassAttribute must be numeric
	if _, ok := classAttrs[0].(*base.FloatAttribute); !ok {
		return []int64{0}, errors.New(fmt.Sprintf("%s: ClassAttribute must be a FloatAttribute", classAttrs[0]))
	}
	// Allocate return structure
	_, rows := X.Size()

	labelVec := make([]int64, rows)
	// Resolve class Attribute specification
	classAttrSpecs := base.ResolveAttributes(X, classAttrs)
	X.MapOverRows(classAttrSpecs, func(row [][]byte, rowNo int) (bool, error) {
		labelVec[rowNo] = int64(base.UnpackBytesToFloat(row[0]))
		return true, nil
	})
	return labelVec, nil
}
