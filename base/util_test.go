package base

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestClassDistributionAfterSplit(t *testing.T) {
	Convey("Given the PlayTennis dataset", t, func() {
		inst, err := ParseCSVToInstances("../examples/datasets/tennis.csv", true)
		So(err, ShouldEqual, nil)

		Convey("Splitting on Sunny should give the right result...", func() {
			result := GetClassDistributionAfterSplit(inst, inst.AllAttributes()[0])
			So(result["sunny"]["no"], ShouldEqual, 3)
			So(result["sunny"]["yes"], ShouldEqual, 2)
			So(result["overcast"]["yes"], ShouldEqual, 4)
			So(result["rainy"]["yes"], ShouldEqual, 3)
			So(result["rainy"]["no"], ShouldEqual, 2)
		})

	})
}

func TestPackAndUnpack(t *testing.T) {
	Convey("Given some uint64", t, func() {
		x := uint64(0xDEADBEEF)
		Convey("When the integer is packed", func() {
			packed := PackU64ToBytes(x)
			Convey("And then unpacked", func() {
				unpacked := UnpackBytesToU64(packed)
				Convey("The unpacked version should be the same", func() {
					So(x, ShouldEqual, unpacked)
				})
			})
		})
	})

	Convey("Given another uint64", t, func() {
		x := uint64(1)
		Convey("When the integer is packed", func() {
			packed := PackU64ToBytes(x)
			Convey("And then unpacked", func() {
				unpacked := UnpackBytesToU64(packed)
				Convey("The unpacked version should be the same", func() {
					So(x, ShouldEqual, unpacked)
				})
			})
		})
	})
}

func TestPackAndUnpackFloat(t *testing.T) {
	Convey("Given some float", t, func() {
		x := 1.2011
		Convey("When the float gets packed", func() {
			packed := PackFloatToBytes(x)
			Convey("And then unpacked", func() {
				unpacked := UnpackBytesToFloat(packed)
				Convey("The unpacked version should be the same", func() {
					So(unpacked, ShouldEqual, x)
				})
			})
		})
	})
}
