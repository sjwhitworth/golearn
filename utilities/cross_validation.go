package utilities

import (
	"fmt"
	mat "github.com/gonum/matrix/mat64"
	"math/rand"
	"time"
)

func shuffleMatrix(dataset *mat.Dense, numGen *rand.Rand) *mat.Dense {
	shuffledSet := mat.DenseCopyOf(dataset)
	rowCount, colCount := shuffledSet.Dims()
	temp := make([]float64, colCount)

	// Fisher–Yates shuffle
	for i := 0; i < rowCount; i++ {
		j := numGen.Intn(i + 1)
		if j != i {
			// Make a "hard" copy to avoid pointer craziness.
			copy(temp, shuffledSet.RowView(i))
			shuffledSet.SetRow(i, shuffledSet.RowView(j))
			shuffledSet.SetRow(j, temp)
		}
	}

	return shuffledSet
}

// TrainTestSplit splits input DenseMatrix into subsets for testing.
// The function expects a test size number (int) or percentage (float64), and a random state or nil to get "random" shuffle.
// It returns a list containing the train-test split and an error status.
func TrainTestSplit(size interface{}, randomState interface{}, datasets ...*mat.Dense) ([]*mat.Dense, error) {
	// Get number of instances (rows).
	instanceCount, _ := datasets[0].Dims()

	// Input should be one or two matrices.
	dataCount := len(datasets)
	if dataCount > 2 {
		return nil, fmt.Errorf("Expected 1 or 2 datasets, got %d\n", dataCount)
	}

	if dataCount == 2 {
		// Test for consistency.
		labelCount, labelFeatures := datasets[1].Dims()
		if labelCount != instanceCount {
			return nil, fmt.Errorf("Data and labels must have the same number of instances")
		} else if labelFeatures != 1 {
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
	var returnDatasets []*mat.Dense

	for _, dataset := range datasets {
		_, featureCount := dataset.Dims()

		tempMatrix := shuffleMatrix(dataset, numGen)

		// Features count is different on data and labels.
		returnDatasets = append(returnDatasets, mat.NewDense(trainSize, featureCount, tempMatrix.RawMatrix().Data[:trainSize*featureCount]))
		returnDatasets = append(returnDatasets, mat.NewDense(testSize, featureCount, tempMatrix.RawMatrix().Data[trainSize*featureCount:]))
	}

	return returnDatasets, nil
}
