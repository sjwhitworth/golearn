package trees

import (
	"fmt"
	base "github.com/sjwhitworth/golearn/base"
	"math/rand"
)

type RandomTreeRuleGenerator struct {
	Attributes   int
	internalRule InformationGainRuleGenerator
}

// So WEKA returns a couple of possible attributes and evaluates
// the split criteria on each
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

type RandomTree struct {
	base.BaseClassifier
	Root *DecisionTreeNode
	Rule *RandomTreeRuleGenerator
}

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

// Train builds a RandomTree suitable for prediction
func (rt *RandomTree) Fit(from *base.Instances) {
	rt.Root = InferID3Tree(from, rt.Rule)
}

// Predict returns a set of Instances containing predictions
func (rt *RandomTree) Predict(from *base.Instances) *base.Instances {
	return rt.Root.Predict(from)
}

func (rt *RandomTree) String() string {
	return fmt.Sprintf("RandomTree(%s)", rt.Root)
}

func (rt *RandomTree) Prune(with *base.Instances) {
	rt.Root.Prune(with)
}
