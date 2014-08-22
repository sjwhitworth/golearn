package base

import (
	"testing"
)

func TestParseCSVGetRows(t *testing.T) {
	lineCount, err := ParseCSVGetRows("../examples/datasets/iris.csv")
	if err != nil {
		t.Fatalf("Unable to parse CSV to get number of rows: %s", err.Error())
	}
	if lineCount != 150 {
		t.Errorf("Should have %d lines, has %d", 150, lineCount)
	}

	lineCount, err = ParseCSVGetRows("../examples/datasets/iris_headers.csv")
	if err != nil {
		t.Fatalf("Unable to parse CSV to get number of rows: %s", err.Error())
	}

	if lineCount != 151 {
		t.Errorf("Should have %d lines, has %d", 151, lineCount)
	}

}

func TestParseCSVGetRowsWithMissingFile(t *testing.T) {
	_, err := ParseCSVGetRows("../examples/datasets/non-existent.csv")
	if err == nil {
		t.Fatal("Expected ParseCSVGetRows to return error when given path to non-existent file")
	}
}

func TestParseCCSVGetAttributes(t *testing.T) {
	attrs := ParseCSVGetAttributes("../examples/datasets/iris_headers.csv", true)
	if attrs[0].GetType() != Float64Type {
		t.Errorf("First attribute should be a float, %s", attrs[0])
	}
	if attrs[0].GetName() != "Sepal length" {
		t.Errorf(attrs[0].GetName())
	}

	if attrs[4].GetType() != CategoricalType {
		t.Errorf("Final attribute should be categorical, %s", attrs[4])
	}
	if attrs[4].GetName() != "Species" {
		t.Error(attrs[4])
	}
}

func TestParseCsvSniffAttributeTypes(t *testing.T) {
	attrs := ParseCSVSniffAttributeTypes("../examples/datasets/iris_headers.csv", true)
	if attrs[0].GetType() != Float64Type {
		t.Errorf("First attribute should be a float, %s", attrs[0])
	}
	if attrs[1].GetType() != Float64Type {
		t.Errorf("Second attribute should be a float, %s", attrs[1])
	}
	if attrs[2].GetType() != Float64Type {
		t.Errorf("Third attribute should be a float, %s", attrs[2])
	}
	if attrs[3].GetType() != Float64Type {
		t.Errorf("Fourth attribute should be a float, %s", attrs[3])
	}
	if attrs[4].GetType() != CategoricalType {
		t.Errorf("Final attribute should be categorical, %s", attrs[4])
	}
}

func TestParseCSVSniffAttributeNamesWithHeaders(t *testing.T) {
	attrs := ParseCSVSniffAttributeNames("../examples/datasets/iris_headers.csv", true)
	if attrs[0] != "Sepal length" {
		t.Error(attrs[0])
	}
	if attrs[1] != "Sepal width" {
		t.Error(attrs[1])
	}
	if attrs[2] != "Petal length" {
		t.Error(attrs[2])
	}
	if attrs[3] != "Petal width" {
		t.Error(attrs[3])
	}
	if attrs[4] != "Species" {
		t.Error(attrs[4])
	}
}

func TestParseCSVToInstances(t *testing.T) {
	inst, err := ParseCSVToInstances("../examples/datasets/iris_headers.csv", true)
	if err != nil {
		t.Error(err)
		return
	}
	row1 := inst.RowString(0)
	row2 := inst.RowString(50)
	row3 := inst.RowString(100)

	if row1 != "5.10 3.50 1.40 0.20 Iris-setosa" {
		t.Error(row1)
	}
	if row2 != "7.00 3.20 4.70 1.40 Iris-versicolor" {
		t.Error(row2)
	}
	if row3 != "6.30 3.30 6.00 2.50 Iris-virginica" {
		t.Error(row3)
	}
}

func TestParseCSVToInstancesWithMissingFile(t *testing.T) {
	_, err := ParseCSVToInstances("../examples/datasets/non-existent.csv", true)
	if err == nil {
		t.Fatal("Expected ParseCSVToInstances to return error when given path to non-existent file")
	}
}

func TestReadAwkwardInsatnces(t *testing.T) {
	inst, err := ParseCSVToInstances("../examples/datasets/chim.csv", true)
	if err != nil {
		t.Error(err)
		return
	}
	attrs := inst.AllAttributes()
	if attrs[0].GetType() != Float64Type {
		t.Error("Should be float!")
	}
	if attrs[1].GetType() != CategoricalType {
		t.Error("Should be discrete!")
	}
}
