// Package knn implements a K Nearest Neighbors object, capable of both classification
// and regression. It accepts data in the form of a slice of float64s, which are then reshaped
// into a X by Y matrix.
package knn

import (
	"github.com/gonum/matrix/mat64"
	base "github.com/sjwhitworth/golearn/base"
	pairwiseMetrics "github.com/sjwhitworth/golearn/metrics/pairwise"
	util "github.com/sjwhitworth/golearn/utilities"
)

// A KNNClassifier consists of a data matrix, associated labels in the same order as the matrix, and a distance function.
// The accepted distance functions at this time are 'euclidean' and 'manhattan'.
type KNNClassifier struct {
	base.BaseEstimator
	TrainingData      base.FixedDataGrid
	DistanceFunc      string
	NearestNeighbours int
}

// NewKnnClassifier returns a new classifier
func NewKnnClassifier(distfunc string, neighbours int) *KNNClassifier {
	KNN := KNNClassifier{}
	KNN.DistanceFunc = distfunc
	KNN.NearestNeighbours = neighbours
	return &KNN
}

// Fit stores the training data for later
func (KNN *KNNClassifier) Fit(trainingData base.FixedDataGrid) {
	KNN.TrainingData = trainingData
}

// Predict returns a classification for the vector, based on a vector input, using the KNN algorithm.
func (KNN *KNNClassifier) Predict(what base.FixedDataGrid) base.FixedDataGrid {

	// Check what distance function we are using
	var distanceFunc pairwiseMetrics.PairwiseDistanceFunc
	switch KNN.DistanceFunc {
	case "euclidean":
		distanceFunc = pairwiseMetrics.NewEuclidean()
	case "manhattan":
		distanceFunc = pairwiseMetrics.NewManhattan()
	case "cosine":
		distanceFunc = pairwiseMetrics.NewCosine()
	default:
		panic("unsupported distance function")

	}
	// Check compatability
	allAttrs := base.CheckCompatable(what, KNN.TrainingData)
	if allAttrs == nil {
		// Don't have the same Attributes
		return nil
	}

	// Remove the Attributes which aren't numeric
	allNumericAttrs := make([]base.Attribute, 0)
	for _, a := range allAttrs {
		if fAttr, ok := a.(*base.FloatAttribute); ok {
			allNumericAttrs = append(allNumericAttrs, fAttr)
		}
	}

	// Generate return vector
	ret := base.GeneratePredictionVector(what)

	// Resolve Attribute specifications for both
	whatAttrSpecs := base.ResolveAttributes(what, allNumericAttrs)
	trainAttrSpecs := base.ResolveAttributes(KNN.TrainingData, allNumericAttrs)

	// Reserve storage for most the most similar items
	distances := make(map[int]float64)

	// Reserve storage for voting map
	maxmap := make(map[string]int)

	// Reserve storage for row computations
	trainRowBuf := make([]float64, len(allNumericAttrs))
	predRowBuf := make([]float64, len(allNumericAttrs))

	// Iterate over all outer rows
	what.MapOverRows(whatAttrSpecs, func(predRow [][]byte, predRowNo int) (bool, error) {
		// Read the float values out
		for i, _ := range allNumericAttrs {
			predRowBuf[i] = base.UnpackBytesToFloat(predRow[i])
		}

		predMat := util.FloatsToMatrix(predRowBuf)

		// Find the closest match in the training data
		KNN.TrainingData.MapOverRows(trainAttrSpecs, func(trainRow [][]byte, srcRowNo int) (bool, error) {

			// Read the float values out
			for i, _ := range allNumericAttrs {
				trainRowBuf[i] = base.UnpackBytesToFloat(trainRow[i])
			}

			// Compute the distance
			trainMat := util.FloatsToMatrix(trainRowBuf)
			distances[srcRowNo] = distanceFunc.Distance(predMat, trainMat)
			return true, nil
		})

		sorted := util.SortIntMap(distances)
		values := sorted[:KNN.NearestNeighbours]

		// Reset maxMap
		for a := range maxmap {
			maxmap[a] = 0
		}

		// Refresh maxMap
		for _, elem := range values {
			label := base.GetClass(KNN.TrainingData, elem)
			if _, ok := maxmap[label]; ok {
				maxmap[label]++
			} else {
				maxmap[label] = 1
			}
		}

		// Sort the maxMap
		var maxClass string
		maxVal := -1
		for a := range maxmap {
			if maxmap[a] > maxVal {
				maxVal = maxmap[a]
				maxClass = a
			}
		}

		base.SetClass(ret, predRowNo, maxClass)
		return true, nil

	})

	return ret
}

// A KNNRegressor consists of a data matrix, associated result variables in the same order as the matrix, and a name.
type KNNRegressor struct {
	base.BaseEstimator
	Values       []float64
	DistanceFunc string
}

// NewKnnRegressor mints a new classifier.
func NewKnnRegressor(distfunc string) *KNNRegressor {
	KNN := KNNRegressor{}
	KNN.DistanceFunc = distfunc
	return &KNN
}

func (KNN *KNNRegressor) Fit(values []float64, numbers []float64, rows int, cols int) {
	if rows != len(values) {
		panic(mat64.ErrShape)
	}

	KNN.Data = mat64.NewDense(rows, cols, numbers)
	KNN.Values = values
}

func (KNN *KNNRegressor) Predict(vector *mat64.Dense, K int) float64 {
	// Get the number of rows
	rows, _ := KNN.Data.Dims()
	rownumbers := make(map[int]float64)
	labels := make([]float64, 0)

	// Check what distance function we are using
	var distanceFunc pairwiseMetrics.PairwiseDistanceFunc
	switch KNN.DistanceFunc {
	case "euclidean":
		distanceFunc = pairwiseMetrics.NewEuclidean()
	case "manhattan":
		distanceFunc = pairwiseMetrics.NewManhattan()
	case "cosine":
		distanceFunc = pairwiseMetrics.NewCosine()
	default:
		panic("unsupported distance function")
	}

	for i := 0; i < rows; i++ {
		row := KNN.Data.RowView(i)
		rowMat := util.FloatsToMatrix(row)
		distance := distanceFunc.Distance(rowMat, vector)
		rownumbers[i] = distance
	}

	sorted := util.SortIntMap(rownumbers)
	values := sorted[:K]

	var sum float64
	for _, elem := range values {
		value := KNN.Values[elem]
		labels = append(labels, value)
		sum += value
	}

	average := sum / float64(K)
	return average
}
