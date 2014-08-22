package filters

import (
	"github.com/sjwhitworth/golearn/base"
	"math"
)

func ChiMBuildFrequencyTable(attr base.Attribute, inst base.FixedDataGrid) []*FrequencyTableEntry {
	ret := make([]*FrequencyTableEntry, 0)
	attribute := attr.(*base.FloatAttribute)

	attrSpec, err := inst.GetAttribute(attr)
	if err != nil {
		panic(err)
	}
	attrSpecs := []base.AttributeSpec{attrSpec}

	err = inst.MapOverRows(attrSpecs, func(row [][]byte, rowNo int) (bool, error) {
		value := row[0]
		valueConv := attribute.GetFloatFromSysVal(value)
		class := base.GetClass(inst, rowNo)
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
		return true, nil
	})

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
