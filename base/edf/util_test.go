package edf

// Utility function tests

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestInt32Conversion(t *testing.T) {
	Convey("Given deadbeef", t, func() {
		buf := make([]byte, 4)
		original := uint32(0xDEAD)
		uint32ToBytes(original, buf)
		converted := uint32FromBytes(buf)
		Convey("Decoded value should be the original...", func() {
			So(converted, ShouldEqual, original)
		})
	})
}
