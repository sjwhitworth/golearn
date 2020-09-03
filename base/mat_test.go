package base

import (
	. "github.com/smartystreets/goconvey/convey"
	"gonum.org/v1/gonum/mat"
	"testing"
)

func TestInlineMat64Creation(t *testing.T) {

	Convey("Given a literal array...", t, func() {
		X := mat.NewDense(4, 3, []float64{
			1, 0, 1,
			0, 1, 1,
			0, 0, 0,
			1, 1, 0,
		})
		inst := InstancesFromMat64(4, 3, X)
		attrs := inst.AllAttributes()
		Convey("Attributes should be well-defined...", func() {
			So(len(attrs), ShouldEqual, 3)
		})

		Convey("No class variables set by default...", func() {
			classAttrs := inst.AllClassAttributes()
			So(len(classAttrs), ShouldEqual, 0)
		})

		Convey("Getting values should work...", func() {
			as, err := inst.GetAttribute(attrs[0])
			So(err, ShouldBeNil)
			valBytes := inst.Get(as, 3)
			val := UnpackBytesToFloat(valBytes)
			So(val, ShouldAlmostEqual, 1.0)
		})

		Convey("Getting size should work...", func() {
			attrLen, rows := inst.Size()
			So(attrLen, ShouldEqual, 3)
			So(rows, ShouldEqual, 4)
		})

		Convey("Getting row string should work...", func() {
			So(inst.RowString(0), ShouldEqual, "0")
		})

		Convey("Getting attribute not in it should error...", func() {
			Y := mat.NewDense(1, 4, []float64{1, 2, 3, 4})
			ins := InstancesFromMat64(1, 4, Y)
			attr := ins.AllAttributes()
			_, err := inst.GetAttribute(attr[3])
			So(err.Error(), ShouldEqual, "Couldn't find a matching attribute")
		})

		Convey("Generate human-readable summary...", func() {
			output := inst.String()
			So(output, ShouldStartWith, "Instances with")
			So(output, ShouldContainSubstring, "Attributes:")
			So(output, ShouldContainSubstring, "Data:")
		})

	})

}

func TestStringWithExceedMaxRow(t *testing.T) {
	Convey("Given a long literal array...", t, func() {
		v := make([]float64, 35, 35)
		X := mat.NewDense(35, 1, v)
		inst := InstancesFromMat64(35, 1, X)
		output := inst.String()
		So(output, ShouldStartWith, "Instances with")
		So(output, ShouldContainSubstring, "Attributes:")
		So(output, ShouldContainSubstring, "Data:")
		So(output, ShouldContainSubstring, "undisplayed")

	})
}
