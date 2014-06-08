package filters

import (
	"fmt"
	base "github.com/sjwhitworth/golearn/base"
	"math"
)

// ChiMergeFilter implements supervised discretisation
// by merging successive numeric intervals if the difference
// in their class distribution is not statistically signficant.
// See Bramer, "Principles of Data Mining", 2nd Edition
//  pp 105--115
type ChiMergeFilter struct {
	Attributes   []int
	Instances    *base.Instances
	Tables       map[int][]*FrequencyTableEntry
	Significance float64
	MinRows      int
	MaxRows      int
	_Trained     bool
}

// Create a ChiMergeFilter with some helpful intialisations.
func NewChiMergeFilter(inst *base.Instances, significance float64) ChiMergeFilter {
	return ChiMergeFilter{
		make([]int, 0),
		inst,
		make(map[int][]*FrequencyTableEntry),
		significance,
		0,
		0,
		false,
	}
}

// Build trains a ChiMergeFilter on the ChiMergeFilter.Instances given
func (c *ChiMergeFilter) Build() {
	for _, attr := range c.Attributes {
		tab := chiMerge(c.Instances, attr, c.Significance, c.MinRows, c.MaxRows)
		c.Tables[attr] = tab
		c._Trained = true
	}
}

// AddAllNumericAttributes adds every suitable attribute
// to the ChiMergeFilter for discretisation
func (b *ChiMergeFilter) AddAllNumericAttributes() {
	for i := 0; i < b.Instances.Cols; i++ {
		if i == b.Instances.ClassIndex {
			continue
		}
		attr := b.Instances.GetAttr(i)
		if attr.GetType() != base.Float64Type {
			continue
		}
		b.Attributes = append(b.Attributes, i)
	}
}

// Run discretises the set of Instances `on'
//
// IMPORTANT: ChiMergeFilter discretises in place.
func (c *ChiMergeFilter) Run(on *base.Instances) {
	if !c._Trained {
		panic("Call Build() beforehand")
	}
	for attr := range c.Tables {
		table := c.Tables[attr]
		for i := 0; i < on.Rows; i++ {
			val := on.Get(i, attr)
			dis := 0
			for j, k := range table {
				if k.Value < val {
					dis = j
					continue
				}
				break
			}
			on.Set(i, attr, float64(dis))
		}
		newAttribute := new(base.CategoricalAttribute)
		newAttribute.SetName(on.GetAttr(attr).GetName())
		for _, k := range table {
			newAttribute.GetSysValFromString(fmt.Sprintf("%f", k.Value))
		}
		on.ReplaceAttr(attr, newAttribute)
	}
}

// AddAttribute add a given numeric Attribute `attr' to the
// filter.
//
// IMPORTANT: This function panic()s if it can't locate the
// attribute in the Instances set.
func (c *ChiMergeFilter) AddAttribute(attr base.Attribute) {
	if attr.GetType() != base.Float64Type {
		panic("ChiMerge only works on Float64Attributes")
	}
	attrIndex := c.Instances.GetAttrIndex(attr)
	if attrIndex == -1 {
		panic("Invalid attribute!")
	}
	c.Attributes = append(c.Attributes, attrIndex)
}

type FrequencyTableEntry struct {
	Value     float64
	Frequency map[string]int
}

func (t *FrequencyTableEntry) String() string {
	return fmt.Sprintf("%.2f %s", t.Value, t.Frequency)
}

func ChiMBuildFrequencyTable(attr int, inst *base.Instances) []*FrequencyTableEntry {
	ret := make([]*FrequencyTableEntry, 0)
	var attribute *base.FloatAttribute
	attribute, ok := inst.GetAttr(attr).(*base.FloatAttribute)
	if !ok {
		panic("only use Chi-M on numeric stuff")
	}
	for i := 0; i < inst.Rows; i++ {
		value := inst.Get(i, attr)
		valueConv := attribute.GetUsrVal(value)
		class := inst.GetClass(i)
		// Search the frequency table for the value
		found := false
		for _, entry := range ret {
			if entry.Value == valueConv {
				found = true
				entry.Frequency[class] += 1
			}
		}
		if !found {
			newEntry := &FrequencyTableEntry{
				valueConv,
				make(map[string]int),
			}
			newEntry.Frequency[class] = 1
			ret = append(ret, newEntry)
		}
	}

	return ret
}

func chiSquaredPdf(k float64, x float64) float64 {
	if x < 0 {
		return 0
	}
	top := math.Pow(x, (k/2)-1) * math.Exp(-x/2)
	bottom := math.Pow(2, k/2) * math.Gamma(k/2)
	return top / bottom
}

func chiSquaredPercentile(k int, x float64) float64 {
	// Implements Yahya et al.'s "A Numerical Procedure
	//  for Computing Chi-Square Percentage Points"
	// InterStat Journal 01/2007; April 25:page:1-8.
	steps := 32
	intervals := 4 * steps
	w := x / (4.0 * float64(steps))
	values := make([]float64, intervals+1)
	for i := 0; i < intervals+1; i++ {
		c := w * float64(i)
		v := chiSquaredPdf(float64(k), c)
		values[i] = v
	}

	ret1 := values[0] + values[len(values)-1]
	ret2 := 0.0
	ret3 := 0.0
	ret4 := 0.0

	for i := 2; i < intervals-1; i += 4 {
		ret2 += values[i]
	}

	for i := 4; i < intervals-3; i += 4 {
		ret3 += values[i]
	}

	for i := 1; i < intervals; i += 2 {
		ret4 += values[i]
	}

	return (2.0 * w / 45) * (7*ret1 + 12*ret2 + 14*ret3 + 32*ret4)
}

func chiCountClasses(entries []*FrequencyTableEntry) map[string]int {
	classCounter := make(map[string]int)
	for _, e := range entries {
		for k := range e.Frequency {
			classCounter[k] += e.Frequency[k]
		}
	}
	return classCounter
}

func chiComputeStatistic(entry1 *FrequencyTableEntry, entry2 *FrequencyTableEntry) float64 {

	// Sum the number of things observed per class
	classCounter := make(map[string]int)
	for k := range entry1.Frequency {
		classCounter[k] += entry1.Frequency[k]
	}
	for k := range entry2.Frequency {
		classCounter[k] += entry2.Frequency[k]
	}

	// Sum the number of things observed per value
	entryObservations1 := 0
	entryObservations2 := 0
	for k := range entry1.Frequency {
		entryObservations1 += entry1.Frequency[k]
	}
	for k := range entry2.Frequency {
		entryObservations2 += entry2.Frequency[k]
	}

	totalObservations := entryObservations1 + entryObservations2
	// Compute the expected values per class
	expectedClassValues1 := make(map[string]float64)
	expectedClassValues2 := make(map[string]float64)
	for k := range classCounter {
		expectedClassValues1[k] = float64(classCounter[k])
		expectedClassValues1[k] *= float64(entryObservations1)
		expectedClassValues1[k] /= float64(totalObservations)
	}
	for k := range classCounter {
		expectedClassValues2[k] = float64(classCounter[k])
		expectedClassValues2[k] *= float64(entryObservations2)
		expectedClassValues2[k] /= float64(totalObservations)
	}

	// Compute chi-squared value
	chiSum := 0.0
	for k := range expectedClassValues1 {
		numerator := float64(entry1.Frequency[k])
		numerator -= expectedClassValues1[k]
		numerator = math.Pow(numerator, 2)
		denominator := float64(expectedClassValues1[k])
		if denominator < 0.5 {
			denominator = 0.5
		}
		chiSum += numerator / denominator
	}
	for k := range expectedClassValues2 {
		numerator := float64(entry2.Frequency[k])
		numerator -= expectedClassValues2[k]
		numerator = math.Pow(numerator, 2)
		denominator := float64(expectedClassValues2[k])
		if denominator < 0.5 {
			denominator = 0.5
		}
		chiSum += numerator / denominator
	}

	return chiSum
}

func chiMergeMergeZipAdjacent(freq []*FrequencyTableEntry, minIndex int) []*FrequencyTableEntry {
	mergeEntry1 := freq[minIndex]
	mergeEntry2 := freq[minIndex+1]
	classCounter := make(map[string]int)
	for k := range mergeEntry1.Frequency {
		classCounter[k] += mergeEntry1.Frequency[k]
	}
	for k := range mergeEntry2.Frequency {
		classCounter[k] += mergeEntry2.Frequency[k]
	}
	newVal := freq[minIndex].Value
	newEntry := &FrequencyTableEntry{
		newVal,
		classCounter,
	}
	lowerSlice := freq
	upperSlice := freq
	if minIndex > 0 {
		lowerSlice = freq[0:minIndex]
		upperSlice = freq[minIndex+1:]
	} else {
		lowerSlice = make([]*FrequencyTableEntry, 0)
		upperSlice = freq[1:]
	}
	upperSlice[0] = newEntry
	freq = append(lowerSlice, upperSlice...)
	return freq
}

func chiMergePrintTable(freq []*FrequencyTableEntry) {
	classes := chiCountClasses(freq)
	fmt.Printf("Attribute value\t")
	for k := range classes {
		fmt.Printf("\t%s", k)
	}
	fmt.Printf("\tTotal\n")
	for _, f := range freq {
		fmt.Printf("%.2f\t", f.Value)
		total := 0
		for k := range classes {
			fmt.Printf("\t%d", f.Frequency[k])
			total += f.Frequency[k]
		}
		fmt.Printf("\t%d\n", total)
	}
}

// Produces a value mapping table
//   inst: The base.Instances which need discretising
//   sig:  The significance level (e.g. 0.95)
//   minrows: The minimum number of rows required in the frequency table
//   maxrows: The maximum number of rows allowed in the frequency table
//            If the number of rows is above this, statistically signficant
//            adjacent rows will be merged
//   precision: internal number of decimal places to round E value to
//              (useful for verification)
func chiMerge(inst *base.Instances, attr int, sig float64, minrows int, maxrows int) []*FrequencyTableEntry {

	// Parameter sanity checking
	if !(2 <= minrows) {
		minrows = 2
	}
	if !(minrows < maxrows) {
		maxrows = minrows + 1
	}
	if sig == 0 {
		sig = 10
	}

	// Build a frequency table
	freq := ChiMBuildFrequencyTable(attr, inst)
	// Count the number of classes
	classes := chiCountClasses(freq)
	for {
		// chiMergePrintTable(freq) DEBUG
		if len(freq) <= minrows {
			break
		}
		minChiVal := math.Inf(1)
		// There may be more than one index to merge
		minChiIndexes := make([]int, 0)
		for i := 0; i < len(freq)-1; i++ {
			chiVal := chiComputeStatistic(freq[i], freq[i+1])
			if chiVal < minChiVal {
				minChiVal = chiVal
				minChiIndexes = make([]int, 0)
			}
			if chiVal == minChiVal {
				minChiIndexes = append(minChiIndexes, i)
			}
		}
		// Only merge if:
		//  We're above the maximum number of rows
		//  OR the chiVal is significant
		//   AS LONG AS we're above the minimum row count
		merge := false
		if len(freq) > maxrows {
			merge = true
		}
		// Compute the degress of freedom |classes - 1| * |rows - 1|
		degsOfFree := len(classes) - 1
		sigVal := chiSquaredPercentile(degsOfFree, minChiVal)
		if sigVal < sig {
			merge = true
		}
		// If we don't need to merge, then break
		if !merge {
			break
		}
		// Otherwise merge the rows i, i+1 by taking
		//  The higher of the two things as the value
		//  Combining the class frequencies
		for i, v := range minChiIndexes {
			freq = chiMergeMergeZipAdjacent(freq, v-i)
		}
	}
	return freq
}
