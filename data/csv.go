/* Data - consists of helper functions for parsing different data formats */

package data

import (
	"encoding/csv"
	"io"
	"os"
	"strconv"

	"github.com/gonum/matrix/mat64"
)

//Parses a CSV file, returning the number of columns and rows, the headers, the labels associated with
//classification, and the data that will be used for training.
func ParseCSV(filepath string, featureCols []int, labelCols []int, header bool) *DataFrame {
	headers := make([]string, 0)
	data := make([]float64, 0)
	labels := make([]string, 0)
	rows := 0

	file, err := os.Open(filepath)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	reader := csv.NewReader(file)

	if header {
		record, err := reader.Read()
		if err != nil {
			panic(err)
		}

		for _, col := range append(featureCols, labelCols...) {
			headers = append(headers, record[col])
		}
		rows += 1
	}

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			panic(err)
		}

		//Iterate over our rows and append the values to a slice
		for _, col := range featureCols {
			entry := record[col]
			number, _ := strconv.ParseFloat(entry, 64)
			data = append(data, number)
		}

		for _, col := range labelCols {
			labels = append(labels, record[col])
		}

		rows += 1
	}

	return &DataFrame{
		Headers:  headers,
		Labels:   labels,
		Values:   mat64.NewDense(rows, len(featureCols), data),
		NRow:     rows,
		NFeature: len(featureCols),
		NLabel:   len(labelCols),
	}
}
