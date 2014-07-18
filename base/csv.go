package base

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"regexp"
	"strconv"
	"strings"
)

// ParseCSVGetRows returns the number of rows in a given file.
func ParseCSVGetRows(filepath string) int {
	file, err := os.Open(filepath)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	counter := 0
	for {
		_, err := reader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			panic(err)
		}
		counter++
	}
	return counter
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
	file, err := os.Open(filepath)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	reader := csv.NewReader(file)
	attrs := make([]Attribute, 0)
	if hasHeaders {
		_, err := reader.Read()
		if err != nil {
			panic(err)
		}
	}
	columns, err := reader.Read()
	if err != nil {
		panic(err)
	}

	for _, entry := range columns {
		entry = strings.Trim(entry, " ")
		matched, err := regexp.MatchString("^[-+]?[0-9]*\\.?[0-9]+([eE][-+]?[0-9]+)?$", entry)
		fmt.Println(entry, matched)
		if err != nil {
			panic(err)
		}
		if matched {
			attrs = append(attrs, NewFloatAttribute())
		} else {
			attrs = append(attrs, new(CategoricalAttribute))
		}
	}

	return attrs
}

// ParseCSVToInstances reads the CSV file given by filepath and returns
// the read Instances.
func ParseCSVToInstances(filepath string, hasHeaders bool) (instances *Instances, err error) {

	defer func() {
		if r := recover(); r != nil {
			var ok bool
			if err, ok = r.(error); !ok {
				err = fmt.Errorf("golearn: ParseCSVToInstances: %v", r)
			}
		}
	}()

	// Read the number of rows in the file
	rowCount := ParseCSVGetRows(filepath)
	if hasHeaders {
		rowCount--
	}

	// Read the row headers
	attrs := ParseCSVGetAttributes(filepath, hasHeaders)

	// Allocate the Instances to return
	instances = NewInstances(attrs, rowCount)

	// Read the input
	file, err := os.Open(filepath)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	reader := csv.NewReader(file)

	rowCounter := 0
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
		for i := range attrs {
			instances.SetAttrStr(rowCounter, i, strings.Trim(record[i], " "))
		}
		rowCounter++
	}

	return
}

//ParseCSV parses a CSV file and returns the number of columns and rows, the headers, the labels associated with
//classification, and the data that will be used for training.
func ParseCSV(filepath string, label int, columns []int) (int, int, []string, []string, []float64) {
	labels := make([]string, 0)
	data := make([]float64, 0)
	headers := make([]string, 0)
	rows := 0

	file, err := os.Open(filepath)
	if err != nil {
		fmt.Println("Error:", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)

	headerrow, _ := reader.Read()

	for _, col := range columns {
		entry := headerrow[col]
		headers = append(headers, entry)
	}

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			fmt.Println("Error:", err)
		}

		//
		labels = append(labels, record[label])

		//Iterate over our rows and append the values to a slice
		for _, col := range columns {
			entry := record[col]
			number, _ := strconv.ParseFloat(entry, 64)
			data = append(data, number)
		}
		rows++
	}
	cols := len(columns)
	return cols, rows, headers, labels, data
}
