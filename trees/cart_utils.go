package trees

import (
	"github.com/sjwhitworth/golearn/base"
)

// Helper Function to check if data point is unique or not.
// We will use this to isolate unique values of a feature
func stringInSlice(a float64, list []float64) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

// Isolate only unique values. This way, we can try only unique splits and not redundant ones.
func findUnique(data []float64) []float64 {
	var unique []float64
	for i := range data {
		if !stringInSlice(data[i], unique) {
			unique = append(unique, data[i])
		}
	}
	return unique
}

// Isolate only the feature being considered for splitting. Reduces the complexity in managing splits.
func getFeature(data [][]float64, feature int64) []float64 {
	var featureVals []float64
	for i := range data {
		featureVals = append(featureVals, data[i][feature])
	}
	return featureVals
}

// Make sure that split being considered has not been done before.
// Else we will unnecessarily try splits that won't improve Impurity.
func validate(triedSplits [][]float64, feature int64, threshold float64) bool {
	for i := range triedSplits {
		split := triedSplits[i]
		featureTried, thresholdTried := split[0], split[1]
		if int64(featureTried) == feature && thresholdTried == threshold {
			return false
		}
	}
	return true
}

// Helper function to convert base.FixedDataGrid into required format. Called in Fit, Predict
func convertInstancesToProblemVec(X base.FixedDataGrid) [][]float64 {
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
