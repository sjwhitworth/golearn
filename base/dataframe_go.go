package base

import (
	"fmt"
	"reflect"
	"strconv"

	"github.com/rocketlaunchr/dataframe-go"
)

// ConvertDataFrameToInstances converts a DataFrame-go dataframe object to Golearn Fixed Data Grid. Allows for compabitibility between dataframe and golearn's ML models.
// df is the dataframe Object. classAttrIndex is the index of the class Attribute in the data.i
func ConvertDataFrameToInstances(df *dataframe.DataFrame, classAttrIndex int) FixedDataGrid {

	// Creating Attributes based on Dataframe
	names := df.Names()
	attrs := make([]Attribute, len(names))

	newInst := NewDenseInstances()

	for i := range names {
		col := df.Series[i]
		if reflect.TypeOf(col.Value(0)).Kind() == reflect.String {
			attrs[i] = new(CategoricalAttribute)
			attrs[i].SetName(names[i])
		} else {
			attrs[i] = NewFloatAttribute(names[i])
		}
	}

	// Add the attributes
	newSpecs := make([]AttributeSpec, len(attrs))
	for i, a := range attrs {

		newSpecs[i] = newInst.AddAttribute(a)
	}
	// Adding the class attribute
	newInst.AddClassAttribute(attrs[classAttrIndex])

	// Allocate space
	nRows := df.NRows()
	newInst.Extend(df.NRows())

	// Write the data based on DataType
	for i := 0; i < nRows; i++ {
		for j := range names {
			col := df.Series[j]

			var val string
			switch v := col.Value(i).(type) {
			case string:
				val = v
			case int64:
				val = strconv.FormatInt(v, 10)
			case float64:
				val = fmt.Sprintf("%f", v)
			case float32:
				val = fmt.Sprintf("%f", v)
			}

			newInst.Set(newSpecs[j], i, newSpecs[j].GetAttribute().GetSysValFromString(val))
		}
	}

	return newInst
}
