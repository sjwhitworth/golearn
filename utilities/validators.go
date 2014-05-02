package utilities

import (
	"fmt"
	mat "github.com/skelterjohn/go.matrix"
	"math/rand"
	"time"
)

func shuffleMatrix(dataset *mat.DenseMatrix, numGen *rand.Rand) *mat.DenseMatrix {
	shuffledSet := dataset.Copy()

	for i := 0; i < shuffledSet.Rows(); i++ {
		j := numGen.Intn(i + 1)
		shuffledSet.SwapRows(i, j)
	}

	return shuffledSet
}

// TrainTestSplit splits input DenseMatrix into subsets for testing.
// The function expects a test size number (int) or percentage (float64), and a random state or nil to get "random" shuffle.
// It returns a list containing the train-test split and an error status.
func TrainTestSplit(size interface{}, randomState interface{}, datasets ...*mat.DenseMatrix) ([]*mat.DenseMatrix, error) {
	// Get number of instances (rows).
	instanceCount := datasets[0].Rows()

	// Input should be one or two matrices.
	dataCount := len(datasets)
	if dataCount > 2 {
		return nil, fmt.Errorf("Expected 1 or 2 datasets, got %d\n", dataCount)
	}

	if dataCount == 2 {
		// Test for consistency.
		if datasets[1].Rows() != instanceCount {
			return nil, fmt.Errorf("Data and labels must have the same number of instances")
		} else if datasets[1].Cols() != 1 {
			return nil, fmt.Errorf("Label matrix must have single feature")
		}
	}

	var trainSize, testSize int
	switch size := size.(type) {
	// If size is an integer, treat it as the test data instance count.
	case int:
		trainSize = instanceCount - size
		testSize = size
	case float64:
		// If size is a float, treat it as a percentage of the instances to be allocated to the test set.
		trainSize = int(float64(instanceCount)*(1-size) + 0.5)
		testSize = int(float64(instanceCount)*size + 0.5)
	default:
		return nil, fmt.Errorf("Expected a test instance count (int) or percentage (float64)")
	}

	// Create a deterministic shuffle, or a "random" one based on current time.
	var randSource rand.Source
	if seed, ok := randomState.(int); ok {
		randSource = rand.NewSource(int64(seed))
	} else {
		randSource = rand.NewSource(time.Now().Unix())
	}
	numGen := rand.New(randSource)

	// Return slice will hold training and test data and optional labels matrix.
	var returnDatasets []*mat.DenseMatrix

	for _, dataset := range datasets {
		tempMatrix := shuffleMatrix(dataset, numGen)
		// Features count is different on data and labels.
		featureCount := tempMatrix.Cols()
		returnDatasets = append(returnDatasets, tempMatrix.GetMatrix(0, 0, trainSize, featureCount))
		returnDatasets = append(returnDatasets, tempMatrix.GetMatrix(trainSize, 0, testSize, featureCount))
	}

	return returnDatasets, nil
}
