// Package KNN implements a K Nearest Neighbors object, capable of both classification
// and regression. It accepts data in the form of a slice of float64s, which are then reshaped
// into a X by Y matrix.
package knn

import (
	base "github.com/sjwhitworth/golearn/base"
	"math"
	"sort"
)

type neighbour struct {
	rowNumber int
	distance  float64
}

type neighbourList []neighbour

func (n neighbourList) Len() int {
	return len(n)
}
func (n neighbourList) Swap(i, j int) {
	n[i], n[j] = n[j], n[i]
}
func (n neighbourList) Less(i, j int) bool {
	return n[i].distance < n[j].distance
}

// A KNN Classifier. Consists of a data matrix, associated labels in the same order as the matrix, and a distance function.
// The accepted distance functions at this time are 'euclidean' and 'manhattan'.
type KNNClassifier struct {
	base.BaseEstimator
	TrainingData      base.FixedDataGrid
	DistanceFunc      string
	NearestNeighbours int
}

// Returns a new classifier
func NewKnnClassifier(distfunc string, neighbours int) *KNNClassifier {
	KNN := KNNClassifier{}
	KNN.DistanceFunc = distfunc
	KNN.NearestNeighbours = neighbours
	return &KNN
}

// Train stores the training data for llater
func (KNN *KNNClassifier) Fit(trainingData base.FixedDataGrid) {
	KNN.TrainingData = trainingData
}

func (KNN *KNNClassifier) Predict(what base.FixedDataGrid) base.UpdatableDataGrid {

	var trainingPond *base.Pond
	var err error

	// TODO: Implement DataGrid.Compatable function

	// Generate the prediction vector
	ret := base.GeneratePredictionVector(what)

	// Process the attributes
	normalAttrs := base.NonClassFloatAttributes(what)

	// Map over the rows

	_, trainSize := KNN.TrainingData.Size()
	distances := make([]float64, trainSize)
	neighbours := neighbourList(make([]neighbour, trainSize))

	maxMap := make(map[string]int)
	normalAttrTestSpecs := base.ResolveAllAttributes(what, normalAttrs)

	if trainInst, ok := KNN.TrainingData.(*base.DenseInstances); ok {
		trainingPond, err = trainInst.GetPond("FLOAT")
		if err != nil {
			panic(err)
		}
		// TODO: Write more tests checking compatability,
		// inclusion of ClassAttributes in this group etc
	}

	rowFloats := make([]float64, len(normalAttrTestSpecs))

	what.MapOverRows(normalAttrTestSpecs, func(pred [][]byte, predRow int) (bool, error) {
		for i := 0; i < len(neighbours); i++ {
			neighbours[i].distance = math.Inf(1)
		}
		for a := range maxMap {
			maxMap[a] = 0
		}

		for i, a := range pred {
			rowFloats[i] = base.UnpackBytesToFloat(a)
		}

		cur := 0
		for _, ref := range trainingPond.Storage() {
			computeDistances(ref, rowFloats, distances)
			cur += ref.Rows
		}

		for i, d := range distances {
			neighbours[i].rowNumber = i
			neighbours[i].distance = d
		}
		sort.Sort(neighbours)

		for i := 0; i < KNN.NearestNeighbours; i++ {
			rowNumber := neighbours[i].rowNumber
			label := base.GetClass(KNN.TrainingData, rowNumber)
			maxMap[label]++
		}

		maxClass := ""
		maxVal := 0
		for i := range maxMap {
			if maxMap[i] > maxVal {
				maxClass = i
				maxVal = maxMap[i]
			}
		}

		base.SetClass(ret, predRow, maxClass)
		return true, nil
	})

	return ret
}

/*
//A KNN Regressor. Consists of a data matrix, associated result variables in the same order as the matrix, and a name.
type KNNRegressor struct {
	base.BaseEstimator
	Values       []float64
	DistanceFunc string
}

// Mints a new classifier.
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
	switch KNN.DistanceFunc {
	case "euclidean":
		{
			euclidean := pairwiseMetrics.NewEuclidean()
			for i := 0; i < rows; i++ {
				row := KNN.Data.RowView(i)
				rowMat := util.FloatsToMatrix(row)
				distance := euclidean.Distance(rowMat, vector)
				rownumbers[i] = distance
			}
		}
	case "manhattan":
		{
			manhattan := pairwiseMetrics.NewEuclidean()
			for i := 0; i < rows; i++ {
				row := KNN.Data.RowView(i)
				rowMat := util.FloatsToMatrix(row)
				distance := manhattan.Distance(rowMat, vector)
				rownumbers[i] = distance
			}
		}
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
}*/
