package base

import (
	"fmt"
	"reflect"
	"strconv"

	dataframe "github.com/rocketlaunchr/dataframe-go"

	"github.com/sjwhitworth/golearn/base"
)

// ConvertDataFrameToInstances converts a DataFrame-go dataframe object to Golearn Fixed Data Grid. Allows for compabitibility between dataframe and golearn's ML models.
// df is the dataframe Object. classAttrIndex is the index of the class Attribute in the data.
func ConvertDataFrameToInstances(df *dataframe.DataFrame, classAttrIndex int) base.FixedDataGrid {

	// Creating Attributes based on Datafraem
	names := df.Names()
	attrs := make([]base.Attribute, len(names))

	newInst := base.NewDenseInstances()

	for i := range names {
		col := df.Series[i]
		if reflect.TypeOf(col.Value(0)).Kind() == reflect.String {
			attrs[i] = new(base.CategoricalAttribute)
			attrs[i].SetName(names[i])
		} else {
			attrs[i] = base.NewFloatAttribute(names[i])
		}
	}

	// Add the attributes
	newSpecs := make([]base.AttributeSpec, len(attrs))
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
			str, ok := col.Value(i).(string)
			val = str
			if ok != true {
				int_64, ok := col.Value(i).(int64)
				val = strconv.FormatInt(int_64, 10)
				if ok != true {
					float_64, ok := col.Value(i).(float64)
					val = fmt.Sprintf("%f", float_64)
					if ok != true {
						float_32, ok := col.Value(i).(float32)
						if ok == true {
							val = fmt.Sprintf("%f", float_32)
						}
					}
				}
			}
			newInst.Set(newSpecs[j], i, newSpecs[j].GetAttribute().GetSysValFromString(val))
		}

	}

	return newInst
}
