package ensemble

import (
	base "github.com/sjwhitworth/golearn/base"
	meta "github.com/sjwhitworth/golearn/meta"
	trees "github.com/sjwhitworth/golearn/trees"
)

// RandomForest classifies instances using an ensemble
// of bagged random decision trees
type RandomForest struct {
	base.BaseClassifier
	ForestSize int
	Features   int
	Model      *meta.BaggedModel
}

// NewRandomForests generates and return a new random forests
// forestSize controls the number of trees that get built
// features controls the number of features used to build each tree
func NewRandomForest(forestSize int, features int) RandomForest {
	ret := RandomForest{
		base.BaseClassifier{},
		forestSize,
		features,
		nil,
	}
	return ret
}

// Train builds the RandomForest on the specified instances
func (f *RandomForest) Train(on *base.Instances) {
	f.Model = new(meta.BaggedModel)
	for i := 0; i < f.ForestSize; i++ {
		tree := new(trees.RandomTree)
		tree.Rules = new(trees.RandomTreeRule)
		tree.Attributes = f.Features
		f.Model.AddModel(tree)
	}
	f.Model.Train(on)
}

// Predict generates predictions from a trained RandomForest
func (f *RandomForest) Predict(with *base.Instances) *base.Instances {
	return f.Model.Predict(with)
}
