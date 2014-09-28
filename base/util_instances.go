package base

import (
	"fmt"
	"math/rand"
)

// This file contains utility functions relating to efficiently
// generating predictions and instantiating DataGrid implementations.

// GeneratePredictionVector selects the class Attributes from a given
// FixedDataGrid and returns something which can hold the predictions.
func GeneratePredictionVector(from FixedDataGrid) UpdatableDataGrid {
	classAttrs := from.AllClassAttributes()
	_, rowCount := from.Size()
	ret := NewDenseInstances()
	for _, a := range classAttrs {
		ret.AddAttribute(a)
		ret.AddClassAttribute(a)
	}
	ret.Extend(rowCount)
	return ret
}

// GetClass is a shortcut for returning the string value of the current
// class on a given row.
//
// IMPORTANT: GetClass will panic if the number of ClassAttributes is
// set to anything other than one.
func GetClass(from DataGrid, row int) string {

	// Get the Attribute
	classAttrs := from.AllClassAttributes()
	if len(classAttrs) > 1 {
		panic("More than one class defined")
	} else if len(classAttrs) == 0 {
		panic("No class defined!")
	}
	classAttr := classAttrs[0]

	// Fetch and convert the class value
	classAttrSpec, err := from.GetAttribute(classAttr)
	if err != nil {
		panic(fmt.Errorf("Can't resolve class Attribute %s", err))
	}

	classVal := from.Get(classAttrSpec, row)
	if classVal == nil {
		panic("Class values shouldn't be missing")
	}

	return classAttr.GetStringFromSysVal(classVal)
}

// SetClass is a shortcut for updating the given class of a row.
//
// IMPORTANT: SetClass will panic if the number of class Attributes
// is anything other than one.
func SetClass(at UpdatableDataGrid, row int, class string) {

	// Get the Attribute
	classAttrs := at.AllClassAttributes()
	if len(classAttrs) > 1 {
		panic("More than one class defined")
	} else if len(classAttrs) == 0 {
		panic("No class Attributes are defined")
	}

	classAttr := classAttrs[0]

	// Fetch and convert the class value
	classAttrSpec, err := at.GetAttribute(classAttr)
	if err != nil {
		panic(fmt.Errorf("Can't resolve class Attribute %s", err))
	}

	classBytes := classAttr.GetSysValFromString(class)
	at.Set(classAttrSpec, row, classBytes)
}

// GetClassDistribution returns a map containing the count of each
// class type (indexed by the class' string representation).
func GetClassDistribution(inst FixedDataGrid) map[string]int {
	ret := make(map[string]int)
	_, rows := inst.Size()
	for i := 0; i < rows; i++ {
		cls := GetClass(inst, i)
		ret[cls]++
	}
	return ret
}

// GetClassDistributionAfterSplit returns the class distribution
// after a speculative split on a given Attribute.
func GetClassDistributionAfterSplit(inst FixedDataGrid, at Attribute) map[string]map[string]int {

	ret := make(map[string]map[string]int)

	// Find the attribute we're decomposing on
	attrSpec, err := inst.GetAttribute(at)
	if err != nil {
		panic(fmt.Sprintf("Invalid attribute %s (%s)", at, err))
	}

	_, rows := inst.Size()

	for i := 0; i < rows; i++ {
		splitVar := at.GetStringFromSysVal(inst.Get(attrSpec, i))
		classVar := GetClass(inst, i)
		if _, ok := ret[splitVar]; !ok {
			ret[splitVar] = make(map[string]int)
			i--
			continue
		}
		ret[splitVar][classVar]++
	}

	return ret
}

// DecomposeOnAttributeValues divides the instance set depending on the
// value of a given Attribute, constructs child instances, and returns
// them in a map keyed on the string value of that Attribute.
//
// IMPORTANT: calls panic() if the AttributeSpec of at cannot be determined.
func DecomposeOnAttributeValues(inst FixedDataGrid, at Attribute) map[string]FixedDataGrid {
	// Find the Attribute we're decomposing on
	attrSpec, err := inst.GetAttribute(at)
	if err != nil {
		panic(fmt.Sprintf("Invalid Attribute index %s", at))
	}
	// Construct the new Attribute set
	newAttrs := make([]Attribute, 0)
	for _, a := range inst.AllAttributes() {
		if a.Equals(at) {
			continue
		}
		newAttrs = append(newAttrs, a)
	}
	// Create the return map
	ret := make(map[string]FixedDataGrid)

	// Create the return row mapping
	rowMaps := make(map[string][]int)

	// Build full Attribute set
	fullAttrSpec := ResolveAttributes(inst, newAttrs)
	fullAttrSpec = append(fullAttrSpec, attrSpec)

	// Decompose
	inst.MapOverRows(fullAttrSpec, func(row [][]byte, rowNo int) (bool, error) {
		// Find the output instance set
		targetBytes := row[len(row)-1]
		targetAttr := fullAttrSpec[len(fullAttrSpec)-1].attr
		targetSet := targetAttr.GetStringFromSysVal(targetBytes)
		if _, ok := rowMaps[targetSet]; !ok {
			rowMaps[targetSet] = make([]int, 0)
		}
		rowMap := rowMaps[targetSet]
		rowMaps[targetSet] = append(rowMap, rowNo)
		return true, nil
	})

	for a := range rowMaps {
		ret[a] = NewInstancesViewFromVisible(inst, rowMaps[a], newAttrs)
	}

	return ret
}

// InstancesTrainTestSplit takes a given Instances (src) and a train-test fraction
// (prop) and returns an array of two new Instances, one containing approximately
// that fraction and the other containing what's left.
//
// IMPORTANT: this function is only meaningful when prop is between 0.0 and 1.0.
// Using any other values may result in odd behaviour.
func InstancesTrainTestSplit(src FixedDataGrid, prop float64) (FixedDataGrid, FixedDataGrid) {
	trainingRows := make([]int, 0)
	testingRows := make([]int, 0)
	src = Shuffle(src)

	// Create the return structure
	_, rows := src.Size()
	for i := 0; i < rows; i++ {
		trainOrTest := rand.Intn(101)
		if trainOrTest > int(100*prop) {
			trainingRows = append(trainingRows, i)
		} else {
			testingRows = append(testingRows, i)
		}
	}

	allAttrs := src.AllAttributes()

	return NewInstancesViewFromVisible(src, trainingRows, allAttrs), NewInstancesViewFromVisible(src, testingRows, allAttrs)

}

// LazyShuffle randomizes the row order without re-ordering the rows
// via an InstancesView.
func LazyShuffle(from FixedDataGrid) FixedDataGrid {
	_, rows := from.Size()
	rowMap := make(map[int]int)
	for i := 0; i < rows; i++ {
		j := rand.Intn(i + 1)
		rowMap[i] = j
		rowMap[j] = i
	}
	return NewInstancesViewFromRows(from, rowMap)
}

// Shuffle randomizes the row order either in place (if DenseInstances)
// or using LazyShuffle.
func Shuffle(from FixedDataGrid) FixedDataGrid {
	_, rows := from.Size()
	if inst, ok := from.(*DenseInstances); ok {
		for i := 0; i < rows; i++ {
			j := rand.Intn(i + 1)
			inst.swapRows(i, j)
		}
		return inst
	} else {
		return LazyShuffle(from)
	}
}

// SampleWithReplacement returns a new FixedDataGrid containing
// an equal number of random rows drawn from the original FixedDataGrid
//
// IMPORTANT: There's a high chance of seeing duplicate rows
// whenever size is close to the row count.
func SampleWithReplacement(from FixedDataGrid, size int) FixedDataGrid {
	rowMap := make(map[int]int)
	_, rows := from.Size()
	for i := 0; i < size; i++ {
		srcRow := rand.Intn(rows)
		rowMap[i] = srcRow
	}
	return NewInstancesViewFromRows(from, rowMap)
}

// CheckCompatible checks whether two DataGrids have the same Attributes
// and if they do, it returns them.
func CheckCompatible(s1 FixedDataGrid, s2 FixedDataGrid) []Attribute {
	s1Attrs := s1.AllAttributes()
	s2Attrs := s2.AllAttributes()
	interAttrs := AttributeIntersect(s1Attrs, s2Attrs)
	if len(interAttrs) != len(s1Attrs) {
		return nil
	} else if len(interAttrs) != len(s2Attrs) {
		return nil
	}
	return interAttrs
}
