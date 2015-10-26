// Package knn implements a K Nearest Neighbors object, capable of both classification
// and regression. It accepts data in the form of a slice of float64s, which are then reshaped
// into a X by Y matrix.
package knn

import (
	"fmt"
	"github.com/gonum/matrix"
	"github.com/gonum/matrix/mat64"
	"github.com/sjwhitworth/golearn/base"
	"github.com/sjwhitworth/golearn/metrics/pairwise"
	"github.com/sjwhitworth/golearn/utilities"
)

// A KNNClassifier consists of a data matrix, associated labels in the same order as the matrix, and a distance function.
// The accepted distance functions at this time are 'euclidean' and 'manhattan'.
// Optimisations only occur when things are identically group into identical
// AttributeGroups, which don't include the class variable, in the same order.
type KNNClassifier struct {
	base.BaseEstimator
	TrainingData       base.FixedDataGrid
	DistanceFunc       string
	NearestNeighbours  int
	AllowOptimisations bool
}

// NewKnnClassifier returns a new classifier
func NewKnnClassifier(distfunc string, neighbours int) *KNNClassifier {
	KNN := KNNClassifier{}
	KNN.DistanceFunc = distfunc
	KNN.NearestNeighbours = neighbours
	KNN.AllowOptimisations = true
	return &KNN
}

// Fit stores the training data for later
func (KNN *KNNClassifier) Fit(trainingData base.FixedDataGrid) {
	KNN.TrainingData = trainingData
}

func (KNN *KNNClassifier) canUseOptimisations(what base.FixedDataGrid) bool {
	// Check that the two have exactly the same layout
	if !base.CheckStrictlyCompatible(what, KNN.TrainingData) {
		return false
	}
	// Check that the two are DenseInstances
	whatd, ok1 := what.(*base.DenseInstances)
	_, ok2 := KNN.TrainingData.(*base.DenseInstances)
	if !ok1 || !ok2 {
		return false
	}
	// Check that no Class Attributes are mixed in with the data
	classAttrs := whatd.AllClassAttributes()
	normalAttrs := base.NonClassAttributes(whatd)
	// Retrieve all the AGs
	ags := whatd.AllAttributeGroups()
	classAttrGroups := make([]base.AttributeGroup, 0)
	for agName := range ags {
		ag := ags[agName]
		attrs := ag.Attributes()
		matched := false
		for _, a := range attrs {
			for _, c := range classAttrs {
				if a.Equals(c) {
					matched = true
				}
			}
		}
		if matched {
			classAttrGroups = append(classAttrGroups, ag)
		}
	}
	for _, cag := range classAttrGroups {
		attrs := cag.Attributes()
		common := base.AttributeIntersect(normalAttrs, attrs)
		if len(common) != 0 {
			return false
		}
	}

	// Check that all of the Attributes are numeric
	for _, a := range normalAttrs {
		if _, ok := a.(*base.FloatAttribute); !ok {
			return false
		}
	}
	// If that's fine, return true
	return true
}

// Predict returns a classification for the vector, based on a vector input, using the KNN algorithm.
func (KNN *KNNClassifier) Predict(what base.FixedDataGrid) base.FixedDataGrid {
	// Check what distance function we are using
	var distanceFunc pairwise.PairwiseDistanceFunc
	switch KNN.DistanceFunc {
	case "euclidean":
		distanceFunc = pairwise.NewEuclidean()
	case "manhattan":
		distanceFunc = pairwise.NewManhattan()
	default:
		panic("unsupported distance function")
	}
	// Check Compatibility
	allAttrs := base.CheckCompatible(what, KNN.TrainingData)
	if allAttrs == nil {
		// Don't have the same Attributes
		return nil
	}

	// Use optimised version if permitted
	if KNN.AllowOptimisations {
		if KNN.DistanceFunc == "euclidean" {
			if KNN.canUseOptimisations(what) {
				return KNN.optimisedEuclideanPredict(what.(*base.DenseInstances))
			}
		}
	}
	fmt.Println("Optimisations are switched off")

	// Remove the Attributes which aren't numeric
	allNumericAttrs := make([]base.Attribute, 0)
	for _, a := range allAttrs {
		if fAttr, ok := a.(*base.FloatAttribute); ok {
			allNumericAttrs = append(allNumericAttrs, fAttr)
		}
	}

	// If every Attribute is a FloatAttribute, then we remove the last one
	// because that is the Attribute we are trying to predict.
	if len(allNumericAttrs) == len(allAttrs) {
		allNumericAttrs = allNumericAttrs[:len(allNumericAttrs)-1]
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

	_, maxRow := what.Size()
	curRow := 0

	// Iterate over all outer rows
	what.MapOverRows(whatAttrSpecs, func(predRow [][]byte, predRowNo int) (bool, error) {

		if (curRow%1) == 0 && curRow > 0 {
			fmt.Printf("KNN: %.2f %% done\n", float64(curRow)*100.0/float64(maxRow))
		}
		curRow++

		// Read the float values out
		for i, _ := range allNumericAttrs {
			predRowBuf[i] = base.UnpackBytesToFloat(predRow[i])
		}

		predMat := utilities.FloatsToMatrix(predRowBuf)

		// Find the closest match in the training data
		KNN.TrainingData.MapOverRows(trainAttrSpecs, func(trainRow [][]byte, srcRowNo int) (bool, error) {
			// Read the float values out
			for i, _ := range allNumericAttrs {
				trainRowBuf[i] = base.UnpackBytesToFloat(trainRow[i])
			}

			// Compute the distance
			trainMat := utilities.FloatsToMatrix(trainRowBuf)
			distances[srcRowNo] = distanceFunc.Distance(predMat, trainMat)
			return true, nil
		})

		sorted := utilities.SortIntMap(distances)
		values := sorted[:KNN.NearestNeighbours]

		maxClass := KNN.vote(maxmap, values)

		base.SetClass(ret, predRowNo, maxClass)
		return true, nil

	})

	return ret
}

func (KNN *KNNClassifier) vote(maxmap map[string]int, values []int) string {
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
	return maxClass
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
		panic(matrix.ErrShape)
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
	var distanceFunc pairwise.PairwiseDistanceFunc
	switch KNN.DistanceFunc {
	case "euclidean":
		distanceFunc = pairwise.NewEuclidean()
	case "manhattan":
		distanceFunc = pairwise.NewManhattan()
	default:
		panic("unsupported distance function")
	}

	for i := 0; i < rows; i++ {
		row := KNN.Data.RowView(i)
		distance := distanceFunc.Distance(utilities.VectorToMatrix(row), vector)
		rownumbers[i] = distance
	}

	sorted := utilities.SortIntMap(rownumbers)
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
