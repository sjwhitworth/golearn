package filters

import (
	base "github.com/sjwhitworth/golearn/base"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestBinning(testEnv *testing.T) {
	//
	// Read the data
	inst1, err := base.ParseCSVToInstances("../examples/datasets/iris_headers.csv", true)
	if err != nil {
		panic(err)
	}

	inst2, err := base.ParseCSVToInstances("../examples/datasets/iris_binned.csv", true)
	if err != nil {
		panic(err)
	}
	//
	// Construct the binning filter
	binAttr := inst1.AllAttributes()[0]
	filt := NewBinningFilter(inst1, 10)
	filt.AddAttribute(binAttr)
	filt.Train()
	inst1f := base.NewLazilyFilteredInstances(inst1, filt)

	// Retrieve the categorical version of the original Attribute

	//
	// Create the LazilyFilteredInstances
	// and check the values
	Convey("Discretized version should match reference", testEnv, func() {
		_, rows := inst1.Size()
		for i := 0; i < rows; i++ {
			So(inst1f.RowString(i), ShouldEqual, inst2.RowString(i))
		}
	})
}
