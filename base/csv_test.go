package base

import (
	"testing"
)

func TestParseCSVGetRows(testEnv *testing.T) {
	lineCount, err := ParseCSVGetRows("../examples/datasets/iris.csv")
	if err != nil {
		testEnv.Fatalf("Unable to parse CSV to get number of rows: %s", err.Error())
	}
	if lineCount != 150 {
		testEnv.Errorf("Should have %d lines, has %d", 150, lineCount)
	}

	lineCount, err = ParseCSVGetRows("../examples/datasets/iris_headers.csv")
	if err != nil {
		testEnv.Fatalf("Unable to parse CSV to get number of rows: %s", err.Error())
	}

	if lineCount != 151 {
		testEnv.Errorf("Should have %d lines, has %d", 151, lineCount)
	}

}

func TestParseCSVGetRowsWithMissingFile(testEnv *testing.T) {
	_, err := ParseCSVGetRows("../examples/datasets/non-existent.csv")
	if err == nil {
		testEnv.Fatal("Expected ParseCSVGetRows to return error when given path to non-existent file")
	}
}

func TestParseCCSVGetAttributes(testEnv *testing.T) {
	attrs := ParseCSVGetAttributes("../examples/datasets/iris_headers.csv", true)
	if attrs[0].GetType() != Float64Type {
		testEnv.Errorf("First attribute should be a float, %s", attrs[0])
	}
	if attrs[0].GetName() != "Sepal length" {
		testEnv.Errorf(attrs[0].GetName())
	}

	if attrs[4].GetType() != CategoricalType {
		testEnv.Errorf("Final attribute should be categorical, %s", attrs[4])
	}
	if attrs[4].GetName() != "Species" {
		testEnv.Error(attrs[4])
	}
}

func TestParseCsvSniffAttributeTypes(testEnv *testing.T) {
	attrs := ParseCSVSniffAttributeTypes("../examples/datasets/iris_headers.csv", true)
	if attrs[0].GetType() != Float64Type {
		testEnv.Errorf("First attribute should be a float, %s", attrs[0])
	}
	if attrs[1].GetType() != Float64Type {
		testEnv.Errorf("Second attribute should be a float, %s", attrs[1])
	}
	if attrs[2].GetType() != Float64Type {
		testEnv.Errorf("Third attribute should be a float, %s", attrs[2])
	}
	if attrs[3].GetType() != Float64Type {
		testEnv.Errorf("Fourth attribute should be a float, %s", attrs[3])
	}
	if attrs[4].GetType() != CategoricalType {
		testEnv.Errorf("Final attribute should be categorical, %s", attrs[4])
	}
}

func TestParseCSVSniffAttributeNamesWithHeaders(testEnv *testing.T) {
	attrs := ParseCSVSniffAttributeNames("../examples/datasets/iris_headers.csv", true)
	if attrs[0] != "Sepal length" {
		testEnv.Error(attrs[0])
	}
	if attrs[1] != "Sepal width" {
		testEnv.Error(attrs[1])
	}
	if attrs[2] != "Petal length" {
		testEnv.Error(attrs[2])
	}
	if attrs[3] != "Petal width" {
		testEnv.Error(attrs[3])
	}
	if attrs[4] != "Species" {
		testEnv.Error(attrs[4])
	}
}

func TestParseCSVToInstances(testEnv *testing.T) {
	inst, err := ParseCSVToInstances("../examples/datasets/iris_headers.csv", true)
	if err != nil {
		testEnv.Error(err)
		return
	}
	row1 := inst.RowString(0)
	row2 := inst.RowString(50)
	row3 := inst.RowString(100)

	if row1 != "5.10 3.50 1.40 0.20 Iris-setosa" {
		testEnv.Error(row1)
	}
	if row2 != "7.00 3.20 4.70 1.40 Iris-versicolor" {
		testEnv.Error(row2)
	}
	if row3 != "6.30 3.30 6.00 2.50 Iris-virginica" {
		testEnv.Error(row3)
	}
}

func TestParseCSVToInstancesWithMissingFile(testEnv *testing.T) {
	_, err := ParseCSVToInstances("../examples/datasets/non-existent.csv", true)
	if err == nil {
		testEnv.Fatal("Expected ParseCSVToInstances to return error when given path to non-existent file")
	}
}

func TestReadAwkwardInsatnces(testEnv *testing.T) {
	inst, err := ParseCSVToInstances("../examples/datasets/chim.csv", true)
	if err != nil {
		testEnv.Error(err)
		return
	}
	attrs := inst.AllAttributes()
	if attrs[0].GetType() != Float64Type {
		testEnv.Error("Should be float!")
	}
	if attrs[1].GetType() != CategoricalType {
		testEnv.Error("Should be discrete!")
	}
}
