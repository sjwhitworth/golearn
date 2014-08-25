package ensemble

import (
	"errors"
	"fmt"
	"github.com/sjwhitworth/golearn/base"
	"github.com/sjwhitworth/golearn/meta"
	"github.com/sjwhitworth/golearn/trees"
)

// RandomForest classifies instances using an ensemble
// of bagged random decision trees.
type RandomForest struct {
	base.BaseClassifier
	ForestSize int
	Features   int
	Model      *meta.BaggedModel
}

// NewRandomForest generates and return a new random forests
// forestSize controls the number of trees that get built
// features controls the number of features used to build each tree.
func NewRandomForest(forestSize int, features int) *RandomForest {
	ret := &RandomForest{
		base.BaseClassifier{},
		forestSize,
		features,
		nil,
	}
	return ret
}

// Fit builds the RandomForest on the specified instances
func (f *RandomForest) Fit(on base.FixedDataGrid) error {
	numNonClassAttributes := len(base.NonClassAttributes(on))
	if numNonClassAttributes < f.Features {
		return errors.New(fmt.Sprintf(
			"Random forest with %d features cannot fit data grid with %d non-class attributes",
			f.Features,
			numNonClassAttributes,
		))
	}

	f.Model = new(meta.BaggedModel)
	f.Model.RandomFeatures = f.Features
	for i := 0; i < f.ForestSize; i++ {
		tree := trees.NewID3DecisionTree(0.00)
		f.Model.AddModel(tree)
	}
	f.Model.Fit(on)
	return nil
}

// Predict generates predictions from a trained RandomForest.
func (f *RandomForest) Predict(with base.FixedDataGrid) (base.FixedDataGrid, error) {
	return f.Model.Predict(with), nil
}

// String returns a human-readable representation of this tree.
func (f *RandomForest) String() string {
	return fmt.Sprintf("RandomForest(ForestSize: %d, Features:%d, %s\n)", f.ForestSize, f.Features, f.Model)
}
