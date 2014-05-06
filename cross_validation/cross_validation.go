package cross_validation

import (
	mat64 "github.com/gonum/matrix/mat64"
	"math/rand"
	"fmt"
	"sync"
)

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

// TrainTestSplit breaks the dataset data into training and testing groups.
// The inputs are the size of the test data (number of items of percentage),
//  an integer to be used as the random state or a nil if no shuffling
// of data is required, and the data and optional label matrices.
// The method returns two arrays. The first contains the training and
// testing subsets, and the second the corresponding labels, or nil for
// unlabeled data.
func TrainTestSplit(size interface{}, randomState interface{}, data ...*mat64.Dense) ([]*mat64.Dense, []*mat64.Dense, error) {
	datasetCount := len(data)

	if datasetCount > 2 {
		return nil, nil, fmt.Errorf("only one data and an optional label matrix are allowed")
	}

	// Get number of instances (rows).
	instanceCount, featureCount := data[0].Dims()

	// Create temporary matrices so the original matrices don't get shuffled.
	returnData := make([]*mat64.Dense, 2)
	tempDataMatrix := mat64.DenseCopyOf(data[0])
	var tempLabelMatrix *mat64.Dense
	var returnLabels []*mat64.Dense

	// Check for data compatibility.
	if datasetCount == 2 {
		labelCount, labelCols := data[1].Dims()
		if labelCount != instanceCount {
			return nil, nil, fmt.Errorf("incompatible data and label matrices")
		} else if labelCols != 1 {
			return nil, nil, fmt.Errorf("label matrix must have single column")
		} else {
			// Create label placeholders.
			tempLabelMatrix = mat64.DenseCopyOf(data[1])
			returnLabels = make([]*mat64.Dense, 2)
		}
	}

	// Determine train/test partitions.
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

	// Make sure partitions have proper sizes.
	if trainSize > instanceCount || trainSize < 0 {
		return nil, nil, fmt.Errorf("test size out of bounds")
	}

	// Create a deterministic shuffle.
	if seed, ok := randomState.(int); ok {
		wg := new(sync.WaitGroup)
		wg.Add(1)
		go shuffleMatrix(tempDataMatrix, int64(seed), wg)

		if datasetCount == 2 {
			wg.Add(1)
			go shuffleMatrix(tempLabelMatrix, int64(seed), wg)
		}
		wg.Wait()
	}

	returnData[0] = mat64.NewDense(trainSize, featureCount, tempDataMatrix.RawMatrix().Data[:trainSize*featureCount])
	returnData[1] = mat64.NewDense(testSize, featureCount, tempDataMatrix.RawMatrix().Data[trainSize*featureCount:])

	if datasetCount == 2 {
		returnLabels[0] = mat64.NewDense(trainSize, 1, tempLabelMatrix.RawMatrix().Data[:trainSize])
		returnLabels[1] = mat64.NewDense(testSize, 1, tempLabelMatrix.RawMatrix().Data[trainSize:])
	}

	return returnData, returnLabels, nil
}


