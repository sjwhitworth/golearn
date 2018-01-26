package linear_models

import (
	"fmt"

	"github.com/amclay/golearn/base"
)

func generateClassWeightVectorFromDist(X base.FixedDataGrid) []float64 {
	classDist := base.GetClassDistributionByBinaryFloatValue(X)
	ret := make([]float64, len(classDist))
	for i, c := range classDist {
		if c == 0 {
			ret[i] = 1.0
		} else {
			ret[i] = 1.0 / float64(c)
		}
	}
	return ret
}

func generateClassWeightVectorFromFixed(X base.FixedDataGrid) []float64 {
	classAttrs := X.AllClassAttributes()
	if len(classAttrs) != 1 {
		panic("Wrong number of class Attributes")
	}
	if _, ok := classAttrs[0].(*base.FloatAttribute); ok {
		ret := make([]float64, 2)
		for i := range ret {
			ret[i] = 1.0
		}
		return ret
	} else {
		panic("Must be a FloatAttribute")
	}
}

func convertInstancesToProblemVec(X base.FixedDataGrid) [][]float64 {
	// Allocate problem array
	_, rows := X.Size()
	problemVec := make([][]float64, rows)

	// Retrieve numeric non-class Attributes
	numericAttrs := base.NonClassFloatAttributes(X)
	numericAttrSpecs := base.ResolveAttributes(X, numericAttrs)

	// Convert each row
	X.MapOverRows(numericAttrSpecs, func(row [][]byte, rowNo int) (bool, error) {
		// Allocate a new row
		probRow := make([]float64, len(numericAttrSpecs))
		// Read out the row
		for i, _ := range numericAttrSpecs {
			probRow[i] = base.UnpackBytesToFloat(row[i])
		}
		// Add the row
		problemVec[rowNo] = probRow
		return true, nil
	})
	return problemVec
}

func convertInstancesToLabelVec(X base.FixedDataGrid) []float64 {
	// Get the class Attributes
	classAttrs := X.AllClassAttributes()
	// Only support 1 class Attribute
	if len(classAttrs) != 1 {
		panic(fmt.Sprintf("%d ClassAttributes (1 expected)", len(classAttrs)))
	}
	// ClassAttribute must be numeric
	if _, ok := classAttrs[0].(*base.FloatAttribute); !ok {
		panic(fmt.Sprintf("%s: ClassAttribute must be a FloatAttribute", classAttrs[0]))
	}
	// Allocate return structure
	_, rows := X.Size()
	labelVec := make([]float64, rows)
	// Resolve class Attribute specification
	classAttrSpecs := base.ResolveAttributes(X, classAttrs)
	X.MapOverRows(classAttrSpecs, func(row [][]byte, rowNo int) (bool, error) {
		labelVec[rowNo] = base.UnpackBytesToFloat(row[0])
		return true, nil
	})
	return labelVec
}
