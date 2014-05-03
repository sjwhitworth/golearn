/* Data - consists of helper functions for parsing different data formats */

package data

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strconv"
)

//Parses a CSV file, returning the number of columns and rows, the headers, the labels associated with
//classification, and the data that will be used for training.
func ParseCsv(filepath string, label int, columns []int) (int, int, []string, []string, []float64) {
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
		rows += 1
	}
	cols := len(columns)
	return cols, rows, headers, labels, data
}
