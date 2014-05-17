package trees

import (
	base "github.com/sjwhitworth/golearn/base"
	"math"
)

//
// Information gain rule generator
//

type InformationGainRuleGenerator struct {
}

// GetSplitAttribute returns the non-class Attribute which maximises the
// information gain.
//
// IMPORTANT: passing a base.Instances with no Attributes other than the class
// variable will panic()
func (r *InformationGainRuleGenerator) GenerateSplitAttribute(f *base.Instances) base.Attribute {
	allAttributes := make([]int, 0)
	for i := 0; i < f.Cols; i++ {
		if i != f.ClassIndex {
			allAttributes = append(allAttributes, i)
		}
	}
	return r.GetSplitAttributeFromSelection(allAttributes, f)
}

// GetSplitAttribute from selection returns the class Attribute which maximises
// the information gain amongst consideredAttributes
//
// IMPORTANT: passing a zero-length consideredAttributes parameter will panic()
func (r *InformationGainRuleGenerator) GetSplitAttributeFromSelection(consideredAttributes []int, f *base.Instances) base.Attribute {

	// Next step is to compute the information gain at this node
	// for each randomly chosen attribute, and pick the one
	// which maximises it
	maxGain := math.Inf(-1)
	selectedAttribute := -1

	// Compute the base entropy
	classDist := f.GetClassDistribution()
	baseEntropy := getBaseEntropy(classDist)

	// Compute the information gain for each attribute
	for _, s := range consideredAttributes {
		proposedClassDist := f.GetClassDistributionAfterSplit(f.GetAttr(s))
		localEntropy := getSplitEntropy(proposedClassDist)
		informationGain := baseEntropy - localEntropy
		if informationGain > maxGain {
			maxGain = informationGain
			selectedAttribute = s
		}
	}

	// Pick the one which maximises IG

	return f.GetAttr(selectedAttribute)
}

//
// Entropy functions
//

// getSplitEntropy determines the entropy of the target
// class distribution after splitting on an base.Attribute
func getSplitEntropy(s map[string]map[string]int) float64 {
	ret := 0.0
	count := 0
	for a := range s {
		for c := range s[a] {
			count += s[a][c]
		}
	}
	for a := range s {
		total := 0.0
		for c := range s[a] {
			total += float64(s[a][c])
		}
		for c := range s[a] {
			ret -= float64(s[a][c]) / float64(count) * math.Log(float64(s[a][c])/float64(count)) / math.Log(2)
		}
		ret += total / float64(count) * math.Log(total/float64(count)) / math.Log(2)
	}
	return ret
}

// getBaseEntropy determines the entropy of the target
// class distribution before splitting on an base.Attribute
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
