package base

import (
	"fmt"
)

// NonClassFloatAttributes returns all FloatAttributes which
// aren't designated as a class Attribute
func NonClassFloatAttributes(d DataGrid) []Attribute {
	classAttrs := d.AllClassAttributes()
	allAttrs := d.AllAttributes()
	ret := make([]Attribute, 0)
	for _, a := range allAttrs {
		matched := false
		if _, ok := a.(*FloatAttribute); !ok {
			continue
		}
		for _, b := range classAttrs {
			if a.Equals(b) {
				matched = true
				break
			}
		}
		if !matched {
			ret = append(ret, a)
		}
	}
	return ret
}

// ResolveAllAttributes returns AttributeSpecs describing
// all of the Attributes
func ResolveAllAttributes(d DataGrid, attrs []Attribute) []AttributeSpec {
	ret := make([]AttributeSpec, len(attrs))
	for i, a := range attrs {
		spec, err := d.GetAttribute(a)
		if err != nil {
			panic(fmt.Errorf("Error resolving Attribute %s: %s", a, err))
		}
		ret[i] = spec
	}
	return ret
}

// GeneratePredictionVector selects the class Attributes from a given
// FixedDataGrid and returns something which can hold the predictions
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

func GetClass(from FixedDataGrid, row int) string {

	// Get the Attribute
	classAttrs := from.AllClassAttributes()
	if len(classAttrs) > 1 {
		panic("More than one class defined")
	}
	classAttr := classAttrs[0]

	// Fetch and convert the class value
	classAttrSpec, err := from.GetAttribute(classAttr)
	if err != nil {
		panic(fmt.Errorf("Can't resolve class Attribute %s", err))
	}

	return classAttr.GetStringFromSysVal(from.Get(classAttrSpec, row))
}

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
