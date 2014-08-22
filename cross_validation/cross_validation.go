package cross_validation

import (
	"fmt"
	"github.com/gonum/matrix/mat64"
	"math/rand"
	"sync"
	"time"
)

func shuffleMatrix(returnDatasets []*mat64.Dense, dataset mat64.Matrix, testSize int, seed int64, wg *sync.WaitGroup) {
	numGen := rand.New(rand.NewSource(seed))

	// We don't want to alter the original dataset.
	shuffledSet := mat64.DenseCopyOf(dataset)
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
	trainSize := rowCount - testSize
	returnDatasets[0] = mat64.NewDense(trainSize, colCount, shuffledSet.RawMatrix().Data[:trainSize*colCount])
	returnDatasets[1] = mat64.NewDense(testSize, colCount, shuffledSet.RawMatrix().Data[trainSize*colCount:])

	wg.Done()
}

// TrainTestSplit splits input DenseMatrix into subsets for testing.
// The function expects a test size number (int) or percentage (float64), and a random state or nil to get "random" shuffle.
// It returns a list containing the train-test split and an error status.
func TrainTestSplit(size interface{}, randomState interface{}, datasets ...*mat64.Dense) ([]*mat64.Dense, error) {
	// Get number of instances (rows).
	instanceCount, _ := datasets[0].Dims()

	// Input should be one or two matrices.
	dataCount := len(datasets)
	if dataCount > 2 {
		return nil, fmt.Errorf("expected 1 or 2 datasets, got %d\n", dataCount)
	}

	if dataCount == 2 {
		// Test for consistency.
		labelCount, labelFeatures := datasets[1].Dims()
		if labelCount != instanceCount {
			return nil, fmt.Errorf("data and labels must have the same number of instances")
		} else if labelFeatures != 1 {
			return nil, fmt.Errorf("label matrix must have single feature")
		}
	}

	var testSize int
	switch size := size.(type) {
	// If size is an integer, treat it as the test data instance count.
	case int:
		testSize = size
	case float64:
		// If size is a float, treat it as a percentage of the instances to be allocated to the test set.
		testSize = int(float64(instanceCount)*size + 0.5)
	default:
		return nil, fmt.Errorf("expected a test instance count (int) or percentage (float64)")
	}

	var randSeed int64
	// Create a deterministic shuffle, or a "random" one based on current time.
	if seed, ok := randomState.(int); ok {
		randSeed = int64(seed)
	} else {
		// Use seconds since epoch as seed
		randSeed = time.Now().Unix()
	}

	// Wait group for goroutine syncronization.
	wg := new(sync.WaitGroup)
	wg.Add(dataCount)

	// Return slice will hold training and test data and optional labels matrix.
	returnDatasets := make([]*mat64.Dense, 2*dataCount)

	for i, dataset := range datasets {
		// Send proper returnDataset slice.
		// This is needed so goroutine doesn't mess up the expected return order.
		// Perhaps returning a map is a better solution...
		go shuffleMatrix(returnDatasets[i:i+2], dataset, testSize, randSeed, wg)
	}
	wg.Wait()

	return returnDatasets, nil
}
