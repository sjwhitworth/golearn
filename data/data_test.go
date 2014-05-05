package data

import (
	"testing"
	. "github.com/smartystreets/goconvey/convey"
)

func TestParseCSV(t *testing.T) {

	Convey("Parse IRIS dataset", t, func() {
		dataFrame := ParseCSV("../examples/datasets/iris.csv", []int{0, 1, 3}, []int{4}, false)

		Convey("First row should be {5.1, 3.5, 0.2}", func() {
			So(dataFrame.Values.RowView(0), ShouldResemble, []float64{5.1, 3.5, 0.2})
		})

		Convey("First label should be Iris-setosa", func() {
			So(dataFrame.Labels[0], ShouldEqual, "Iris-setosa")
		})

		Convey("Headers should be empty", func() {
			So(dataFrame.Headers, ShouldResemble, []string{})
		})

		Convey("Number of features should be 3", func() {
			So(dataFrame.NFeature, ShouldEqual, 3)
		})

		Convey("Number of labels should be 1", func() {
			So(dataFrame.NLabel, ShouldEqual, 1)
		})

		Convey("Number of rows should be 150", func() {
			So(dataFrame.NRow, ShouldEqual, 150)
		})

	})
}
