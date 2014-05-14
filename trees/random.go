package trees

import (
	"fmt"
	base "github.com/sjwhitworth/golearn/base"
	"math"
	"math/rand"
)

type RandomTreeRuleGenerator struct {
	Attributes int
}

func getSplitEntropy(s map[string]map[string]int) float64 {
	ret := 0.0
	count := 0
	for a := range s {
		total := 0.0
		for c := range s[a] {
			ret -= float64(s[a][c]) * math.Log(float64(s[a][c])) / math.Log(2)
			total += float64(s[a][c])
			count += s[a][c]
		}
		ret += total * math.Log(total) / math.Log(2)
	}
	return ret / float64(count)
}

func getBaseEntropy(s map[string]int) float64 {
	ret := 0.0
	count := 0
	for k := range s {
		count += s[k]
	}
	for k := range s {
		ret -= float64(s[k]) / float64(count) * math.Log(float64(s[k])/float64(count)) / math.Log(2)
	}
	return ret
}

// So WEKA returns a couple of possible attributes and evaluates
// the split criteria on each
func (r *RandomTreeRuleGenerator) GenerateSplitAttribute(f *base.Instances) base.Attribute {

	fmt.Println("GenerateSplitAttribute", r.Attributes)

	// First step is to generate the random attributes that we'll consider
	maximumAttribute := f.GetAttributeCount()
	consideredAttributes := make([]int, r.Attributes)
	attrCounter := 0
	for {
		if attrCounter >= r.Attributes {
			break
		}
		selectedAttribute := rand.Intn(maximumAttribute)
		if selectedAttribute != f.ClassIndex {
			consideredAttributes = append(consideredAttributes, selectedAttribute)
			attrCounter++
		}
	}

	// Next step is to compute the information gain at this node
	// for each randomly chosen attribute, and pick the one
	// which maximises it
	maxGain := math.Inf(-1)
	selectedAttribute := -1

	// Compute the base entropy
	classDist := f.GetClassDistribution()
	baseEntropy := getBaseEntropy(classDist)

	for _, s := range consideredAttributes {
		proposedClassDist := f.GetClassDistributionAfterSplit(f.GetAttr(s))
		localEntropy := getSplitEntropy(proposedClassDist)
		informationGain := baseEntropy - localEntropy
		if informationGain > maxGain {
			maxGain = localEntropy
			selectedAttribute = s
			fmt.Printf("Gain: %.4f, selectedAttribute: %s\n", informationGain, f.GetAttr(selectedAttribute))
		}
	}

	return f.GetAttr(selectedAttribute)
}

type RandomTree struct {
	base.BaseClassifier
	Root *DecisionTreeNode
	Rule RandomTreeRuleGenerator
}

func NewRandomTree(attrs int) *RandomTree {
	return &RandomTree{
		base.BaseClassifier{},
		nil,
		RandomTreeRuleGenerator{
			attrs,
		},
	}
}

// Train builds a RandomTree suitable for prediction
func (rt *RandomTree) Fit(from *base.Instances) {
	rt.Root = InferDecisionTree(from, &rt.Rule)
}

// Predict returns a set of Instances containing predictions
func (rt *RandomTree) Predict(from *base.Instances) *base.Instances {
	return rt.Root.Predict(from)
}

func (rt *RandomTree) String() string {
	return fmt.Sprintf("RandomTree(%s)", rt.Root)
}
