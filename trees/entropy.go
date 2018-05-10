package trees

import (
	"github.com/sjwhitworth/golearn/base"
	"math"
	"sort"
)

//
// Information gain rule generator
//

// InformationGainRuleGenerator generates DecisionTreeRules which
// maximize information gain at each node.
type InformationGainRuleGenerator struct {
}

// GenerateSplitRule returns a DecisionTreeNode based on a non-class Attribute
// which maximises the information gain.
//
// IMPORTANT: passing a base.Instances with no Attributes other than the class
// variable will panic()
func (r *InformationGainRuleGenerator) GenerateSplitRule(f base.FixedDataGrid) *DecisionTreeRule {

	attrs := f.AllAttributes()
	classAttrs := f.AllClassAttributes()
	candidates := base.AttributeDifferenceReferences(attrs, classAttrs)

	return r.GetSplitRuleFromSelection(candidates, f)
}

// GetSplitRuleFromSelection returns a DecisionTreeRule which maximises
// the information gain amongst the considered Attributes.
//
// IMPORTANT: passing a zero-length consideredAttributes parameter will panic()
func (r *InformationGainRuleGenerator) GetSplitRuleFromSelection(consideredAttributes []base.Attribute, f base.FixedDataGrid) *DecisionTreeRule {

	var selectedAttribute base.Attribute

	// Parameter check
	if len(consideredAttributes) == 0 {
		panic("More Attributes should be considered")
	}

	// Next step is to compute the information gain at this node
	// for each randomly chosen attribute, and pick the one
	// which maximises it
	maxGain := math.Inf(-1)
	selectedVal := math.Inf(1)

	// Compute the base entropy
	classDist := base.GetClassDistribution(f)
	baseEntropy := getBaseEntropy(classDist)

	// Compute the information gain for each attribute
	for _, s := range consideredAttributes {
		var informationGain float64
		var splitVal float64
		if fAttr, ok := s.(*base.FloatAttribute); ok {
			var attributeEntropy float64
			attributeEntropy, splitVal = getNumericAttributeEntropy(f, fAttr)
			informationGain = baseEntropy - attributeEntropy
		} else {
			proposedClassDist := base.GetClassDistributionAfterSplit(f, s)
			localEntropy := getSplitEntropy(proposedClassDist)
			informationGain = baseEntropy - localEntropy
		}

		if informationGain > maxGain {
			maxGain = informationGain
			selectedAttribute = s
			selectedVal = splitVal
		}
	}

	// Pick the one which maximises IG
	return &DecisionTreeRule{selectedAttribute, selectedVal}
}

//
// Entropy functions
//

type numericSplitRef struct {
	val   float64
	class int
}

type splitVec []numericSplitRef

func (a splitVec) Len() int           { return len(a) }
func (a splitVec) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a splitVec) Less(i, j int) bool { return a[i].val < a[j].val }

func getNumericAttributeEntropy(f base.FixedDataGrid, attr *base.FloatAttribute) (float64, float64) {

	// Resolve Attribute
	attrSpec, err := f.GetAttribute(attr)
	if err != nil {
		panic(err)
	}

	// Build sortable vector
	_, rows := f.Size()
	refs := make([]numericSplitRef, rows)
	numClasses := 0
	class2Int := make(map[string]int)
	f.MapOverRows([]base.AttributeSpec{attrSpec}, func(val [][]byte, row int) (bool, error) {
		cls := base.GetClass(f, row)
		i, ok := class2Int[cls]
		if !ok {
			i = numClasses
			class2Int[cls] = i
			numClasses++
		}
		v := base.UnpackBytesToFloat(val[0])
		refs[row] = numericSplitRef{v, i}
		return true, nil
	})

	sort.Sort(splitVec(refs))

	minSplitEntropy := math.Inf(1)
	minSplitVal := math.Inf(1)
	prevVal := math.NaN()
	prevInd := 0

	splitDist := [2][]int{make([]int, numClasses), make([]int, numClasses)}
	// Before first split all refs are not smaller than val
	for _, x := range refs {
		splitDist[1][x.class]++
	}

	// Consider each possible function
	for i := 0; i < len(refs)-1; {
		val := refs[i].val + refs[i+1].val
		val /= 2
		if val == prevVal {
			i++
			continue
		}
		// refs is sorted, so we only need to update values that are
		// bigger than prevVal, but are lower than val
		for j := prevInd; j < len(refs) && refs[j].val < val; j++ {
			splitDist[0][refs[j].class]++
			splitDist[1][refs[j].class]--
			i++
			prevInd++
		}
		prevVal = val
		splitEntropy := getSplitEntropyFast(splitDist)
		if splitEntropy < minSplitEntropy {
			minSplitEntropy = splitEntropy
			minSplitVal = val
		}
	}

	return minSplitEntropy, minSplitVal
}

// getSplitEntropyFast determines the entropy of the target
// class distribution after splitting on an base.Attribute.
// It is similar to getSplitEntropy, but accepts array of slices,
// to avoid map access overhead.
func getSplitEntropyFast(s [2][]int) float64 {
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
			if s[a][c] != 0 {
				ret -= float64(s[a][c]) / float64(count) * math.Log(float64(s[a][c])/float64(count)) / math.Log(2)
			}
		}
		ret += total / float64(count) * math.Log(total/float64(count)) / math.Log(2)
	}
	return ret
}

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
