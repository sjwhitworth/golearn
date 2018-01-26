package main

// This example program demonstrates Instances

import (
	"fmt"

	"github.com/amclay/golearn/base"
)

func main() {

	// Instances can be read using ParseCsvToInstances
	rawData, err := base.ParseCSVToInstances("../datasets/iris_headers.csv", true)
	if err != nil {
		panic(err)
	}

	// Instances can be printed, and you'll see a human-readable summary
	// if you do so. The first section is a line like
	//     Instances with 150 row(s) and 5 attribute(s)
	//
	// It next prints all the attributes
	//     FloatAttribute(Sepal length)
	//     FloatAttribute(Sepal width)
	//     FloatAttribute(Petal length)
	//     FloatAttribute(Petal width)
	//     CategoricalAttribute([Iris-setosa Iris-versicolor Iris-viriginica])
	// The final attribute has an asterisk (*) printed before it,
	// meaning that it is the class variable. It then prints out up to
	// 30 rows which correspond to those attributes.
	// 	5.10 3.50 1.40 0.20 Iris-setosa
	// 	4.90 3.00 1.40 0.20 Iris-setosa
	fmt.Println(rawData)

	// If two decimal places isn't enough, you can update the
	// Precision field on any FloatAttribute
	if attr, ok := rawData.AllAttributes()[0].(*base.FloatAttribute); !ok {
		panic("Invalid cast")
	} else {
		attr.Precision = 4
	}
	// Now the first column has more precision
	fmt.Println(rawData)

	// We can update the set of Instances, although the API
	// for doing so is not very sophisticated.

	// First, have to resolve Attribute Specifications
	as := base.ResolveAttributes(rawData, rawData.AllAttributes())

	// Attribute Specifications describe where a given column lives
	rawData.Set(as[0], 0, as[0].GetAttribute().GetSysValFromString("1.00"))

	// A SetClass method exists as a shortcut
	base.SetClass(rawData, 0, "Iris-unusual")
	fmt.Println(rawData)

	// There is a way of creating new Instances from scratch.
	// Inside an Instance, everything's stored as float64
	newData := make([]float64, 2)
	newData[0] = 1.0
	newData[1] = 0.0

	// Let's create some attributes
	attrs := make([]base.Attribute, 2)
	attrs[0] = base.NewFloatAttribute("Arbitrary Float Quantity")
	attrs[1] = new(base.CategoricalAttribute)
	attrs[1].SetName("Class")
	// Insert a standard class
	attrs[1].GetSysValFromString("A")

	// Now let's create the final instances set
	newInst := base.NewDenseInstances()

	// Add the attributes
	newSpecs := make([]base.AttributeSpec, len(attrs))
	for i, a := range attrs {
		newSpecs[i] = newInst.AddAttribute(a)
	}

	// Allocate space
	newInst.Extend(1)

	// Write the data
	newInst.Set(newSpecs[0], 0, newSpecs[0].GetAttribute().GetSysValFromString("1.0"))
	newInst.Set(newSpecs[1], 0, newSpecs[1].GetAttribute().GetSysValFromString("A"))

	fmt.Println(newInst)

}
