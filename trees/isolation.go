package trees

import (
	"math"
	"math/rand"
	"time"

	"github.com/sjwhitworth/golearn/base"
)

type IsolationForest struct {
	nTrees   int
	maxDepth int
	subSpace int
	trees    []regressorNode
}

func selectFeature(data [][]float64) int64 {
	rand.Seed(time.Now().UnixNano())
	return int64(rand.Intn(len(data[0])))
}

func minMax(feature int64, data [][]float64) (float64, float64) {

	var min, max float64

	min = math.Inf(1)
	max = math.Inf(-1)
	for _, instance := range data {
		if instance[feature] > max {
			max = instance[feature]
		}
		if instance[feature] < min {
			min = instance[feature]
		}
	}

	return min, max
}

func selectValue(min, max float64) float64 {
	rand.Seed(time.Now().UnixNano())
	val := min + (rand.Float64() * (max - min))
	if val == min {
		val += 0.000001
	} else if val == max {
		val -= 0.000001
	}
	return val
}

func splitData(val float64, feature int64, data [][]float64) ([][]float64, [][]float64) {
	var leftData, rightData [][]float64
	for _, instance := range data {
		if instance[feature] <= val {
			leftData = append(leftData, instance)
		} else {
			rightData = append(rightData, instance)
		}
	}
	return leftData, rightData
}

func checkData(data [][]float64) bool {
	for _, instance := range data {
		for i, val := range instance {
			if val != data[0][i] {
				return true
			}
		}
	}
	return false
}

func buildTree(data [][]float64, upperNode regressorNode, depth int, maxDepth int) regressorNode {
	depth++

	upperNode.isNodeNeeded = true
	if (depth > maxDepth) || (len(data) <= 1) || !checkData(data) {
		upperNode.isNodeNeeded = false
		return upperNode
	}

	var featureToSplit int64
	var splitVal float64
	var min, max float64
	min, max = 0.0, 0.0

	for min == max {
		featureToSplit = selectFeature(data)
		min, max = minMax(featureToSplit, data)
		splitVal = selectValue(min, max)
	}

	leftData, rightData := splitData(splitVal, featureToSplit, data)

	upperNode.Feature = featureToSplit
	upperNode.Threshold = splitVal
	upperNode.LeftPred = float64(len(leftData))
	upperNode.RightPred = float64(len(rightData))

	var leftNode, rightNode regressorNode
	leftNode = buildTree(leftData, leftNode, depth, maxDepth)
	rightNode = buildTree(rightData, rightNode, depth, maxDepth)

	if leftNode.isNodeNeeded {
		upperNode.Left = &leftNode
	}
	if rightNode.isNodeNeeded {
		upperNode.Right = &rightNode
	}

	return upperNode
}

func getRandomData(data [][]float64, subSpace int) [][]float64 {
	var randomData [][]float64
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < subSpace; i++ {
		randomData = append(randomData, data[rand.Intn(len(data))])
	}
	return randomData
}

func NewIsolationForest(nTrees int, maxDepth int, subSpace int) IsolationForest {
	var iForest IsolationForest
	iForest.nTrees = nTrees
	iForest.maxDepth = maxDepth
	iForest.subSpace = subSpace
	return iForest
}

func (iForest *IsolationForest) Fit(X base.FixedDataGrid) {
	data := preprocessData(X)
	nTrees := iForest.nTrees
	subSpace := iForest.subSpace
	maxDepth := iForest.maxDepth

	var forest []regressorNode
	for i := 0; i < nTrees; i++ {
		subData := getRandomData(data, subSpace)
		var tree regressorNode

		tree = buildTree(subData, tree, 0, maxDepth)
		forest = append(forest, tree)
	}
	iForest.trees = forest
}

func pathLength(tree regressorNode, instance []float64, path float64) float64 {
	path++

	if instance[tree.Feature] <= tree.Threshold {
		if tree.Left == nil {
			if tree.LeftPred <= 1 {
				return path
			} else {
				return path + cFactor(int(tree.LeftPred))
			}
		}
		path = pathLength(*tree.Left, instance, path)
	} else {
		if tree.Right == nil {
			if tree.RightPred <= 1 {
				return path
			} else {
				return path + cFactor(int(tree.RightPred))
			}
		}
		path = pathLength(*tree.Right, instance, path)
	}
	return path
}

func evaluateInstance(instance []float64, forest []regressorNode) []float64 {
	var paths []float64
	for _, tree := range forest {
		paths = append(paths, pathLength(tree, instance, 0))
	}
	return paths
}

func cFactor(n int) float64 {
	return 2.0*(math.Log(float64(n-1))+0.5772156649) - (float64(2.0*(n-1)) / float64(n))
}

func anomalyScore(instance []float64, forest []regressorNode, n int) float64 {
	paths := evaluateInstance(instance, forest)
	E := 0.0
	for _, path := range paths {
		E += path
	}
	E /= float64(len(paths))
	c := cFactor(n)
	return math.Pow(2, (-1 * E / c))
}
func (iForest *IsolationForest) Predict(X base.FixedDataGrid) []float64 {
	data := preprocessData(X)

	var preds []float64
	for _, instance := range data {
		score := anomalyScore(instance, iForest.trees, iForest.subSpace)
		preds = append(preds, score)
	}
	return preds
}

func preprocessData(X base.FixedDataGrid) [][]float64 {
	data := convertInstancesToProblemVec(X)
	class, err := regressorConvertInstancesToLabelVec(X)
	if err != nil {
		panic(err)
	}
	for i, point := range class {
		data[i] = append(data[i], point)
	}
	return data
}
