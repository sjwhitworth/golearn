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
func (r *RandomTreeRuleGenerator) GenerateSplitAttribute(f *base.Instances) base.Attribute {

	// First step is to generate the random attributes that we'll consider
	maximumAttribute := f.GetAttributeCount()
	consideredAttributes := make([]int, r.Attributes)
	attrCounter := 0
	for {
		if len(consideredAttributes) >= r.Attributes {
			break
		}
		selectedAttribute := rand.Intn(maximumAttribute)
		fmt.Println(selectedAttribute, attrCounter, consideredAttributes, len(consideredAttributes))
		if selectedAttribute != f.ClassIndex {
			matched := false
			for _, a := range consideredAttributes {
				if a == selectedAttribute {
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
func (rt *RandomTree) Fit(from *base.Instances) {
	rt.Root = InferID3Tree(from, rt.Rule)
}

// Predict returns a set of Instances containing predictions
func (rt *RandomTree) Predict(from *base.Instances) *base.Instances {
	return rt.Root.Predict(from)
}

// String returns a human-readable representation of this structure
func (rt *RandomTree) String() string {
	return fmt.Sprintf("RandomTree(%s)", rt.Root)
}

// Prune removes nodes from the tree which are detrimental
// to determining the accuracy of the test set (with)
func (rt *RandomTree) Prune(with *base.Instances) {
	rt.Root.Prune(with)
}
