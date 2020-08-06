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
	MAE string = "mae"
	MSE string = "mse"
)

// RNode - Node struct for Decision Tree Regressor
// It holds the information for each split
// Which feature to use, threshold, left prediction and right prediction
type regressorNode struct {
	Left         *regressorNode
	Right        *regressorNode
	Threshold    float64
	Feature      int64
	LeftPred     float64
	RightPred    float64
	isNodeNeeded bool
}

// CARTDecisionTreeRegressor - Tree struct for Decision Tree Regressor
// It contains the rootNode, as well as the hyperparameters chosen by user.
// Also keeps track of splits used at tree level.
type CARTDecisionTreeRegressor struct {
	RootNode    *regressorNode
	criterion   string
	maxDepth    int64
	triedSplits [][]float64
}

// Find average
func average(y []float64) float64 {
	mean := 0.0
	for _, value := range y {
		mean += value
	}
	mean /= float64(len(y))
	return mean
}

// Calculate Mean Absolute Error for a constant prediction
func meanAbsoluteError(y []float64, yBar float64) float64 {
	error := 0.0
	for _, target := range y {
		error += math.Abs(target - yBar)
	}
	error /= float64(len(y))
	return error
}

// Turn Mean Absolute Error into impurity function for decision trees.
func computeMaeImpurityAndAverage(y []float64) (float64, float64) {
	yHat := average(y)
	return meanAbsoluteError(y, yHat), yHat
}

// Calculate Mean Squared Error for constant prediction
func meanSquaredError(y []float64, yBar float64) float64 {
	error := 0.0
	for _, target := range y {
		itemError := target - yBar
		error += math.Pow(itemError, 2)
	}
	error /= float64(len(y))
	return error
}

// Convert mean squared error into impurity function for decision trees
func computeMseImpurityAndAverage(y []float64) (float64, float64) {
	yHat := average(y)
	return meanSquaredError(y, yHat), yHat
}

func calculateRegressionLoss(y []float64, criterion string) (float64, float64, error) {
	if criterion == MAE {
		loss, avg := computeMaeImpurityAndAverage(y)
		return loss, avg, nil
	} else if criterion == MSE {
		loss, avg := computeMseImpurityAndAverage(y)
		return loss, avg, nil
	} else {
		panic("Invalid impurity function, choose from MAE or MSE")
	}
}

// Split the data into left and right based on trehsold and feature.
func regressorCreateSplit(data [][]float64, feature int64, y []float64, threshold float64) ([][]float64, [][]float64, []float64, []float64) {
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

// Interface for creating new Decision Tree Regressor
func NewDecisionTreeRegressor(criterion string, maxDepth int64) *CARTDecisionTreeRegressor {
	var tree CARTDecisionTreeRegressor
	tree.maxDepth = maxDepth
	tree.criterion = strings.ToLower(criterion)
	return &tree
}

// Re order data based on a feature for optimizing code
// Helps in updating splits without reiterating entire dataset
func regressorReOrderData(featureVal []float64, data [][]float64, y []float64) ([][]float64, []float64) {
	s := NewSlice(featureVal)
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

// Update the left and right data based on change in threshold
func regressorUpdateSplit(left [][]float64, leftY []float64, right [][]float64, rightY []float64, feature int64, threshold float64) ([][]float64, []float64, [][]float64, []float64) {

	for right[0][feature] < threshold {
		left = append(left, right[0])
		right = right[1:]
		leftY = append(leftY, rightY[0])
		rightY = rightY[1:]
	}

	return left, leftY, right, rightY
}

// Fit - Build the tree using the data
// Creates empty root node and builds tree by calling regressorBestSplit
func (tree *CARTDecisionTreeRegressor) Fit(X base.FixedDataGrid) error {
	var emptyNode regressorNode
	var err error

	data := regressorConvertInstancesToProblemVec(X)
	y, err := regressorConvertInstancesToLabelVec(X)
	if err != nil {
		return err
	}

	emptyNode, err = regressorBestSplit(*tree, data, y, emptyNode, tree.criterion, tree.maxDepth, 0)
	if err != nil {
		return err
	}
	tree.RootNode = &emptyNode
	return nil
}

// Builds the tree by iteratively finding the best split.
// Recursive function - stops if maxDepth is reached or nodes are pure
func regressorBestSplit(tree CARTDecisionTreeRegressor, data [][]float64, y []float64, upperNode regressorNode, criterion string, maxDepth int64, depth int64) (regressorNode, error) {

	// Ensure that we have not reached maxDepth. maxDepth =-1 means split until nodes are pure
	depth++

	if depth > maxDepth && maxDepth != -1 {
		return upperNode, nil
	}

	numFeatures := len(data[0])
	var bestLoss, origLoss float64
	var err error
	origLoss, upperNode.LeftPred, err = calculateRegressionLoss(y, criterion)
	if err != nil {
		return upperNode, err
	}

	bestLoss = origLoss

	bestLeft, bestRight, bestLefty, bestRighty := data, data, y, y

	numData := len(data)

	bestLeftLoss, bestRightLoss := bestLoss, bestLoss

	upperNode.isNodeNeeded = true

	var leftN, rightN regressorNode
	// Iterate over all features
	for i := 0; i < numFeatures; i++ {

		featureVal := getFeature(data, int64(i))
		unique := findUnique(featureVal)
		sort.Float64s(unique)

		sortData, sortY := regressorReOrderData(featureVal, data, y)

		firstTime := true

		var left, right [][]float64
		var leftY, rightY []float64

		for j := 0; j < len(unique)-1; j++ {
			threshold := (unique[j] + unique[j+1]) / 2
			if validate(tree.triedSplits, int64(i), threshold) {
				if firstTime {
					left, right, leftY, rightY = regressorCreateSplit(sortData, int64(i), sortY, threshold)
					firstTime = false
				} else {
					left, leftY, right, rightY = regressorUpdateSplit(left, leftY, right, rightY, int64(i), threshold)
				}

				var leftLoss, rightLoss float64
				var leftPred, rightPred float64

				leftLoss, leftPred, _ = calculateRegressionLoss(leftY, criterion)
				rightLoss, rightPred, _ = calculateRegressionLoss(rightY, criterion)

				subLoss := (leftLoss * float64(len(left)) / float64(numData)) + (rightLoss * float64(len(right)) / float64(numData))

				if subLoss < bestLoss {
					bestLoss = subLoss

					bestLeft, bestRight = left, right
					bestLefty, bestRighty = leftY, rightY

					upperNode.Threshold, upperNode.Feature = threshold, int64(i)

					upperNode.LeftPred, upperNode.RightPred = leftPred, rightPred

					bestLeftLoss, bestRightLoss = leftLoss, rightLoss
				}
			}
		}
	}

	if bestLoss == origLoss {
		upperNode.isNodeNeeded = false
		return upperNode, nil
	}

	if bestLoss > 0 {

		if bestLeftLoss > 0 {
			tree.triedSplits = append(tree.triedSplits, []float64{float64(upperNode.Feature), upperNode.Threshold})
			leftN, err = regressorBestSplit(tree, bestLeft, bestLefty, leftN, criterion, maxDepth, depth)
			if err != nil {
				return upperNode, err
			}
			if leftN.isNodeNeeded == true {
				upperNode.Left = &leftN
			}
		}

		if bestRightLoss > 0 {
			tree.triedSplits = append(tree.triedSplits, []float64{float64(upperNode.Feature), upperNode.Threshold})
			rightN, err = regressorBestSplit(tree, bestRight, bestRighty, rightN, criterion, maxDepth, depth)
			if err != nil {
				return upperNode, err
			}
			if rightN.isNodeNeeded == true {
				upperNode.Right = &rightN
			}
		}
	}
	return upperNode, nil
}

// Print Tree for Visualtion - calls regressorPrintTreeFromNode()
func (tree *CARTDecisionTreeRegressor) String() string {
	rootNode := *tree.RootNode
	return regressorPrintTreeFromNode(rootNode, "")
}

// Recursively explore the entire tree and print out all details such as threshold, feature, prediction
func regressorPrintTreeFromNode(tree regressorNode, spacing string) string {
	returnString := ""
	returnString += spacing + "Feature "
	returnString += strconv.FormatInt(tree.Feature, 10)
	returnString += " < "
	returnString += fmt.Sprintf("%.3f", tree.Threshold)
	returnString += "\n"

	if tree.Left == nil {
		returnString += spacing + "---> True" + "\n"
		returnString += "  " + spacing + "PREDICT    "
		returnString += fmt.Sprintf("%.3f", tree.LeftPred) + "\n"
	}
	if tree.Right == nil {
		returnString += spacing + "---> False" + "\n"
		returnString += "  " + spacing + "PREDICT    "
		returnString += fmt.Sprintf("%.3f", tree.RightPred) + "\n"
	}

	if tree.Left != nil {
		returnString += spacing + "---> True" + "\n"
		returnString += regressorPrintTreeFromNode(*tree.Left, spacing+"  ")
	}

	if tree.Right != nil {
		returnString += spacing + "---> False" + "\n"
		returnString += regressorPrintTreeFromNode(*tree.Right, spacing+"  ")
	}

	return returnString
}

// Predict a single data point by navigating to rootNodes.
// Uses a recursive logic
func regressorPredictSingle(tree regressorNode, instance []float64) float64 {
	if instance[tree.Feature] < tree.Threshold {
		if tree.Left == nil {
			return tree.LeftPred
		} else {
			return regressorPredictSingle(*tree.Left, instance)
		}
	} else {
		if tree.Right == nil {
			return tree.RightPred
		} else {
			return regressorPredictSingle(*tree.Right, instance)
		}
	}
}

// Predict method for multiple data points.
// First converts input data into usable format, and then calls regressorPredictFromNode
func (tree *CARTDecisionTreeRegressor) Predict(X_test base.FixedDataGrid) []float64 {
	root := *tree.RootNode
	test := regressorConvertInstancesToProblemVec(X_test)
	return regressorPredictFromNode(root, test)
}

// Use tree's root node to print out entire tree.
// Iterates over all data points and calls regressorPredictSingle to predict individual datapoints.
func regressorPredictFromNode(tree regressorNode, test [][]float64) []float64 {
	var preds []float64
	for i := range test {
		i_pred := regressorPredictSingle(tree, test[i])
		preds = append(preds, i_pred)
	}
	return preds
}

// Helper function to convert base.FixedDataGrid into required format. Called in Fit, Predict
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

// Helper function to convert base.FixedDataGrid into required format. Called in Fit, Predict
func regressorConvertInstancesToLabelVec(X base.FixedDataGrid) ([]float64, error) {
	// Get the class Attributes
	classAttrs := X.AllClassAttributes()
	// Only support 1 class Attribute
	if len(classAttrs) != 1 {
		return []float64{0}, errors.New(fmt.Sprintf("%d ClassAttributes (1 expected)", len(classAttrs)))
	}
	// ClassAttribute must be numeric
	if _, ok := classAttrs[0].(*base.FloatAttribute); !ok {
		return []float64{0}, errors.New(fmt.Sprintf("%s: ClassAttribute must be a FloatAttribute", classAttrs[0]))
	}
	// Allocate return structure
	_, rows := X.Size()

	labelVec := make([]float64, rows)
	// Resolve class Attribute specification
	classAttrSpecs := base.ResolveAttributes(X, classAttrs)
	X.MapOverRows(classAttrSpecs, func(row [][]byte, rowNo int) (bool, error) {
		labelVec[rowNo] = base.UnpackBytesToFloat(row[0])
		return true, nil
	})
	return labelVec, nil
}
