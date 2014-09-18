package base

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"
)

// ParseCSVGetRows returns the number of rows in a given file.
func ParseCSVGetRows(filepath string) (int, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return 0, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	counter := 0
	for {
		_, err := reader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			return 0, err
		}
		counter++
	}
	return counter, nil
}

// ParseCSVGetAttributes returns an ordered slice of appropriate-ly typed
// and named Attributes.
func ParseCSVGetAttributes(filepath string, hasHeaders bool) []Attribute {
	attrs := ParseCSVSniffAttributeTypes(filepath, hasHeaders)
	names := ParseCSVSniffAttributeNames(filepath, hasHeaders)
	for i, attr := range attrs {
		attr.SetName(names[i])
	}
	return attrs
}

// ParseCSVSniffAttributeNames returns a slice containing the top row
// of a given CSV file, or placeholders if hasHeaders is false.
func ParseCSVSniffAttributeNames(filepath string, hasHeaders bool) []string {
	file, err := os.Open(filepath)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	headers, err := reader.Read()
	if err != nil {
		panic(err)
	}

	if hasHeaders {
		for i, h := range headers {
			headers[i] = strings.TrimSpace(h)
		}
		return headers
	}

	for i := range headers {
		headers[i] = fmt.Sprintf("%d", i)
	}
	return headers

}

// ParseCSVSniffAttributeTypes returns a slice of appropriately-typed Attributes.
//
// The type of a given attribute is determined by looking at the first data row
// of the CSV.
func ParseCSVSniffAttributeTypes(filepath string, hasHeaders bool) []Attribute {
	var attrs []Attribute
	// Open file
	file, err := os.Open(filepath)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	// Create the CSV reader
	reader := csv.NewReader(file)
	if hasHeaders {
		// Skip the headers
		_, err := reader.Read()
		if err != nil {
			panic(err)
		}
	}
	// Read the first line of the file
	columns, err := reader.Read()
	if err != nil {
		panic(err)
	}

	for _, entry := range columns {
		// Match the Attribute type with regular expressions
		entry = strings.Trim(entry, " ")
		matched, err := regexp.MatchString("^[-+]?[0-9]*\\.?[0-9]+([eE][-+]?[0-9]+)?$", entry)
		if err != nil {
			panic(err)
		}
		if matched {
			attrs = append(attrs, NewFloatAttribute(""))
		} else {
			attrs = append(attrs, new(CategoricalAttribute))
		}
	}

	return attrs
}

// ParseCSVBuildInstances updates an [[#UpdatableDataGrid]] from a filepath in place
func ParseCSVBuildInstances(filepath string, hasHeaders bool, u UpdatableDataGrid) {

	// Read the input
	file, err := os.Open(filepath)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	reader := csv.NewReader(file)

	rowCounter := 0

	specs := ResolveAttributes(u, u.AllAttributes())

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			panic(err)
		}
		if rowCounter == 0 {
			if hasHeaders {
				hasHeaders = false
				continue
			}
		}
		for i, v := range record {
			u.Set(specs[i], rowCounter, specs[i].attr.GetSysValFromString(v))
		}
		rowCounter++
	}

}

// ParseCSVToInstances reads the CSV file given by filepath and returns
// the read Instances.
func ParseCSVToInstances(filepath string, hasHeaders bool) (instances *DenseInstances, err error) {

	// Read the number of rows in the file
	rowCount, err := ParseCSVGetRows(filepath)
	if err != nil {
		return nil, err
	}

	if hasHeaders {
		rowCount--
	}

	// Read the row headers
	attrs := ParseCSVGetAttributes(filepath, hasHeaders)
	specs := make([]AttributeSpec, len(attrs))
	// Allocate the Instances to return
	instances = NewDenseInstances()
	for i, a := range attrs {
		spec := instances.AddAttribute(a)
		specs[i] = spec
	}
	instances.Extend(rowCount)

	// Read the input
	file, err := os.Open(filepath)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	reader := csv.NewReader(file)

	rowCounter := 0

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			return nil, err
		}
		if rowCounter == 0 {
			if hasHeaders {
				hasHeaders = false
				continue
			}
		}
		for i, v := range record {
			v = strings.Trim(v, " ")
			instances.Set(specs[i], rowCounter, attrs[i].GetSysValFromString(v))
		}
		rowCounter++
	}

	instances.AddClassAttribute(attrs[len(attrs)-1])

	return instances, nil
}

// ParseCSVToInstancesTemplated reads the CSV file given by filepath and returns
// the read Instances, using another already read DenseInstances as a template.
func ParseCSVToTemplatedInstances(filepath string, hasHeaders bool, template *DenseInstances) (instances *DenseInstances, err error) {

	// Read the number of rows in the file
	rowCount, err := ParseCSVGetRows(filepath)
	if err != nil {
		return nil, err
	}

	if hasHeaders {
		rowCount--
	}

	// Read the row headers
	attrs := ParseCSVGetAttributes(filepath, hasHeaders)
	templateAttrs := template.AllAttributes()
	for i, a := range attrs {
		for _, b := range templateAttrs {
			if a.Equals(b) {
				attrs[i] = b
			} else if a.GetName() == b.GetName() {
				attrs[i] = b
			}
		}
	}

	specs := make([]AttributeSpec, len(attrs))
	// Allocate the Instances to return
	instances = NewDenseInstances()

	templateAgs := template.AllAttributeGroups()
	for ag := range templateAgs {
		agTemplate := templateAgs[ag]
		if _, ok := agTemplate.(*BinaryAttributeGroup); ok {
			instances.CreateAttributeGroup(ag, 0)
		} else {
			instances.CreateAttributeGroup(ag, 8)
		}
	}

	for i, a := range templateAttrs {
		s, err := template.GetAttribute(a)
		if err != nil {
			panic(err)
		}
		if ag, ok := template.agRevMap[s.pond]; !ok {
			panic(ag)
		} else {
			spec, err := instances.AddAttributeToAttributeGroup(a, ag)
			if err != nil {
				panic(err)
			}
			specs[i] = spec
		}
	}

	instances.Extend(rowCount)

	// Read the input
	file, err := os.Open(filepath)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	reader := csv.NewReader(file)

	rowCounter := 0

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			return nil, err
		}
		if rowCounter == 0 {
			if hasHeaders {
				hasHeaders = false
				continue
			}
		}
		for i, v := range record {
			v = strings.Trim(v, " ")
			instances.Set(specs[i], rowCounter, attrs[i].GetSysValFromString(v))
		}
		rowCounter++
	}

	for _, a := range template.AllClassAttributes() {
		instances.AddClassAttribute(a)
	}

	return instances, nil
}

// ParseCSVToInstancesWithAttributeGroups reads the CSV file given by filepath,
// and returns the read DenseInstances, but also makes sure to group any Attributes
// specified in the first argument and also any class Attributes specified in the second
func ParseCSVToInstancesWithAttributeGroups(filepath string, attrGroups, classAttrGroups map[string]string, attrOverrides map[int]Attribute, hasHeaders bool) (instances *DenseInstances, err error) {

	// Read row count
	rowCount, err := ParseCSVGetRows(filepath)
	if err != nil {
		return nil, err
	}

	// Read the row headers
	attrs := ParseCSVGetAttributes(filepath, hasHeaders)
	for i := range attrs {
		if a, ok := attrOverrides[i]; ok {
			attrs[i] = a
		}
	}

	specs := make([]AttributeSpec, len(attrs))
	// Allocate the Instances to return
	instances = NewDenseInstances()

	//
	// Create all AttributeGroups
	agsToCreate := make(map[string]int)
	combinedAgs := make(map[string]string)
	for a := range attrGroups {
		agsToCreate[attrGroups[a]] = 0
		combinedAgs[a] = attrGroups[a]
	}
	for a := range classAttrGroups {
		agsToCreate[classAttrGroups[a]] = 8
		combinedAgs[a] = classAttrGroups[a]
	}

	// Decide the sizes
	for _, a := range attrs {
		if ag, ok := combinedAgs[a.GetName()]; ok {
			if _, ok := a.(*BinaryAttribute); ok {
				agsToCreate[ag] = 0
			} else {
				agsToCreate[ag] = 8
			}
		}
	}

	// Create them
	for i := range agsToCreate {
		size := agsToCreate[i]
		err = instances.CreateAttributeGroup(i, size)
		if err != nil {
			panic(err)
		}
	}

	// Add the Attributes to them
	for i, a := range attrs {
		var spec AttributeSpec
		if ag, ok := combinedAgs[a.GetName()]; ok {
			spec, err = instances.AddAttributeToAttributeGroup(a, ag)
			if err != nil {
				panic(err)
			}
			specs[i] = spec
		} else {
			spec = instances.AddAttribute(a)
		}
		specs[i] = spec
		if _, ok := classAttrGroups[a.GetName()]; ok {
			err = instances.AddClassAttribute(a)
			if err != nil {
				panic(err)
			}
		}
	}
	// Allocate
	instances.Extend(rowCount)

	// Read the input
	file, err := os.Open(filepath)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	reader := csv.NewReader(file)

	rowCounter := 0

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			return nil, err
		}
		if rowCounter == 0 {
			// Skip header row
			rowCounter++
			continue
		}
		for i, v := range record {
			v = strings.Trim(v, " ")
			instances.Set(specs[i], rowCounter, attrs[i].GetSysValFromString(v))
		}
		rowCounter++
	}

	// Add class Attributes
	for _, a := range instances.AllAttributes() {
		name := a.GetName() // classAttrGroups
		if _, ok := classAttrGroups[name]; ok {
			err = instances.AddClassAttribute(a)
			if err != nil {
				panic(err)
			}
		}
	}
	return instances, nil

}
