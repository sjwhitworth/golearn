package trees

import (
	"github.com/sjwhitworth/golearn/base"
	"math"
)

//
// Information gain rule generator
//

type InformationGainRuleGenerator struct {
}

// GenerateSplitAttribute returns the non-class Attribute which maximises the
// information gain.
//
// IMPORTANT: passing a base.Instances with no Attributes other than the class
// variable will panic()
func (r *InformationGainRuleGenerator) GenerateSplitAttribute(f base.FixedDataGrid) base.Attribute {

	attrs := f.AllAttributes()
	classAttrs := f.AllClassAttributes()
	candidates := base.AttributeDifferenceReferences(attrs, classAttrs)

	return r.GetSplitAttributeFromSelection(candidates, f)
}

// GetSplitAttributeFromSelection returns the class Attribute which maximises
// the information gain amongst consideredAttributes
//
// IMPORTANT: passing a zero-length consideredAttributes parameter will panic()
func (r *InformationGainRuleGenerator) GetSplitAttributeFromSelection(consideredAttributes []base.Attribute, f base.FixedDataGrid) base.Attribute {

	var selectedAttribute base.Attribute

	// Parameter check
	if len(consideredAttributes) == 0 {
		panic("More Attributes should be considered")
	}

	// Next step is to compute the information gain at this node
	// for each randomly chosen attribute, and pick the one
	// which maximises it
	maxGain := math.Inf(-1)

	// Compute the base entropy
	classDist := base.GetClassDistribution(f)
	baseEntropy := getBaseEntropy(classDist)

	// Compute the information gain for each attribute
	for _, s := range consideredAttributes {
		proposedClassDist := base.GetClassDistributionAfterSplit(f, s)
		localEntropy := getSplitEntropy(proposedClassDist)
		informationGain := baseEntropy - localEntropy
		if informationGain > maxGain {
			maxGain = informationGain
			selectedAttribute = s
		}
	}

	// Pick the one which maximises IG
	return selectedAttribute
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
