package trees

import (
	"fmt"
	base "github.com/sjwhitworth/golearn/base"
	"math/rand"
)

// RandomTreeRuleGenerator is used to generate decision rules for Random Trees
type RandomTreeRuleGenerator struct {
	Attributes   int
	internalRule InformationGainRuleGenerator
}

// GenerateSplitAttribute returns the best attribute out of those randomly chosen
// which maximises Information Gain
func (r *RandomTreeRuleGenerator) GenerateSplitAttribute(f base.FixedDataGrid) base.Attribute {

	// First step is to generate the random attributes that we'll consider
	allAttributes := base.AttributeDifferenceReferences(f.AllAttributes(), f.AllClassAttributes())
	maximumAttribute := len(allAttributes)
	consideredAttributes := make([]base.Attribute, 0)

	attrCounter := 0
	for {
		if len(consideredAttributes) >= r.Attributes {
			break
		}
		selectedAttrIndex := rand.Intn(maximumAttribute)
		selectedAttribute := allAttributes[selectedAttrIndex]
		matched := false
		for _, a := range consideredAttributes {
			if a.Equals(selectedAttribute) {
				matched = true
				break
			}
		}
		if matched {
			continue
		}
		consideredAttributes = append(consideredAttributes, selectedAttribute)
		attrCounter++
	}

	return r.internalRule.GetSplitAttributeFromSelection(consideredAttributes, f)
}

// RandomTree builds a decision tree by considering a fixed number
// of randomly-chosen attributes at each node
type RandomTree struct {
	base.BaseClassifier
	Root *DecisionTreeNode
	Rule *RandomTreeRuleGenerator
}

// NewRandomTree returns a new RandomTree which considers attrs randomly
// chosen attributes at each node.
func NewRandomTree(attrs int) *RandomTree {
	return &RandomTree{
		base.BaseClassifier{},
		nil,
		&RandomTreeRuleGenerator{
			attrs,
			InformationGainRuleGenerator{},
		},
	}
}

// Fit builds a RandomTree suitable for prediction
func (rt *RandomTree) Fit(from base.FixedDataGrid) {
	rt.Root = InferID3Tree(from, rt.Rule)
}

// Predict returns a set of Instances containing predictions
func (rt *RandomTree) Predict(from base.FixedDataGrid) base.FixedDataGrid {
	return rt.Root.Predict(from)
}

// String returns a human-readable representation of this structure
func (rt *RandomTree) String() string {
	return fmt.Sprintf("RandomTree(%s)", rt.Root)
}

// Prune removes nodes from the tree which are detrimental
// to determining the accuracy of the test set (with)
func (rt *RandomTree) Prune(with base.FixedDataGrid) {
	rt.Root.Prune(with)
}
