package cross_validation

import (
	mat64 "github.com/gonum/matrix/mat64"
	"math/rand"
	"fmt"
	"sync"
)

type dataset struct {
	data, labels *mat64.Dense
}

// shuffleMatrix can be used to shuffle any Matrix type (including DenseMatrix and Vec).
// It accepts a matrix and a seed. The latter is used to repeat shuffles.
// All changes to the matrix are made in place.
func shuffleMatrix(vector *mat64.Dense, seed int64, wg *sync.WaitGroup) {
	defer wg.Done()

	// Create repeatable random number generator.
	numGen := rand.New(rand.NewSource(seed))

	rows, cols := vector.Dims()
	temp := make([]float64, cols)

	// Fisher–Yates shuffle
	for i := 0; i < rows; i++ {
		j := numGen.Intn(i+1)
		if j != i {
			// Replace all values on that row.
			copy(temp, vector.RowView(i))
			vector.SetRow(i, vector.RowView(j))
			vector.SetRow(j, temp)
		}
	}
}

// NewDataset creates a new structure to hold both the data and the labels.
// For unlabled data, the second parameter should be nil.
// The function returns a pointer to a dataset type.
func NewDataset(data, labels *mat64.Dense) (*dataset, error) {
	if _, cols := labels.Dims(); cols != 1 {
		return nil, fmt.Errorf("label vector must have single column")
	}

	return &dataset{data, labels}, nil
}

// TrainTestSplit breaks the dataset data into training and testing groups.
// The inputs are the size of the test data (number of items of percentage),
// and an integer to be used as the random state or a nil if no shuffling
// of data is required.
// The method returns two arrays. The first contains the training and
// testing subsets, and the second the corresponding labels, or nil for
// unlabeled data.
func (ds *dataset) TrainTestSplit(size interface{}, randomState interface{}) ([]*mat64.Dense, []*mat64.Dense, error) {
	// Get number of instances (rows).
	instanceCount, featureCount := ds.data.Dims()

	var testSize int
	switch size := size.(type) {
	case int:
		// If size is an integer, treat it as the test data instance count.
		testSize = size
	case float64:
		// If size is a float, treat it as a percentage of the instances to be allocated to the test set.
		testSize = int(float64(instanceCount)*size + 0.5)
	default:
		return nil, nil, fmt.Errorf("expected a test size (int) or percentage (float64)")
	}
	trainSize := instanceCount - testSize

	if trainSize > instanceCount || trainSize < 0 {
		return nil, nil, fmt.Errorf("test size out of bounds")
	}

	returnData := make([]*mat64.Dense, 2)
	var returnLabels []*mat64.Dense

	returnData[0] = mat64.NewDense(trainSize, featureCount, ds.data.RawMatrix().Data[:trainSize*featureCount])
	returnData[1] = mat64.NewDense(testSize, featureCount, ds.data.RawMatrix().Data[trainSize*featureCount:])

	if ds.labels != nil {
		returnLabels = make([]*mat64.Dense, 2)
		returnLabels[0] = mat64.NewDense(trainSize, 1, ds.labels.RawMatrix().Data[:trainSize])
		returnLabels[1] = mat64.NewDense(testSize, 1, ds.labels.RawMatrix().Data[trainSize:])
	}

	// Create a deterministic shuffle, or a "random" one based on current time.
	wg := new(sync.WaitGroup)
	if seed, ok := randomState.(int); ok {
		wg.Add(2)
		go shuffleMatrix(returnData[0], int64(seed), wg)
		go shuffleMatrix(returnData[1], int64(seed), wg)

		if ds.labels != nil {
			wg.Add(2)
			go shuffleMatrix(returnLabels[0], int64(seed), wg)
			go shuffleMatrix(returnLabels[1], int64(seed), wg)
		}
	}
	wg.Wait()

	return returnData, returnLabels, nil
}


