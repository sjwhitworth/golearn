package filters

import (
	"testing"

	"github.com/amclay/golearn/base"
	. "github.com/smartystreets/goconvey/convey"
)

func TestBinning(t *testing.T) {
	Convey("Given some data and a reference", t, func() {
		inst1, err := base.ParseCSVToInstances("../examples/datasets/iris_headers.csv", true)
		So(err, ShouldBeNil)

		inst2, err := base.ParseCSVToInstances("../examples/datasets/iris_binned.csv", true)
		So(err, ShouldBeNil)

		//
		// Construct the binning filter
		binAttr := inst1.AllAttributes()[0]
		filt := NewBinningFilter(inst1, 10)
		filt.AddAttribute(binAttr)
		filt.Train()
		inst1f := base.NewLazilyFilteredInstances(inst1, filt)

		// Retrieve the categorical version of the original Attribute
		var cAttr base.Attribute
		for _, a := range inst1f.AllAttributes() {
			if a.GetName() == binAttr.GetName() {
				cAttr = a
			}
		}

		cAttrSpec, err := inst1f.GetAttribute(cAttr)
		So(err, ShouldEqual, nil)
		binAttrSpec, err := inst2.GetAttribute(binAttr)
		So(err, ShouldEqual, nil)

		//
		// Create the LazilyFilteredInstances
		// and check the values
		Convey("Discretized version should match reference", func() {
			_, rows := inst1.Size()
			for i := 0; i < rows; i++ {
				val1 := inst1f.Get(cAttrSpec, i)
				val2 := inst2.Get(binAttrSpec, i)
				val1s := cAttr.GetStringFromSysVal(val1)
				val2s := binAttr.GetStringFromSysVal(val2)
				So(val1s, ShouldEqual, val2s)
			}
		})
	})
}
