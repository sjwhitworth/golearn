package data

import (
	"testing"
	. "github.com/smartystreets/goconvey/convey"
)

func TestStringFrame(t *testing.T) {

	Convey("Add 2 rows", t, func() {
		strFrame := &StringFrame{}
		strFrame.Add([]string{"label11", "label12", "label13"})
		strFrame.Add([]string{"label21", "label22", "label23"})

		Convey("String at (1, 1) should be label22", func() {
			So(strFrame.At(1, 1), ShouldEqual, "label22")
		})

		Convey("First row should be {label11, label12, label13}", func() {
			So(strFrame.Row(0), ShouldResemble, []string{"label11", "label12", "label13"})
		})

		Convey("Second column should be {label12, label22}", func() {
			So(strFrame.Col(1), ShouldResemble, []string{"label12", "label22"})
		})

		Convey("Number of column should be 3", func() {
			So(strFrame.NCol(), ShouldEqual, 3)
		})

		Convey("Number of row should be 2", func() {
			So(strFrame.NRow(), ShouldEqual, 2)
		})

	})
}
