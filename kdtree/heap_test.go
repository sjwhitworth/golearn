package kdtree

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestHeap(t *testing.T) {
	Convey("Given a heap", t, func() {
		h := newHeap()

		Convey("When heap is empty", func() {

			Convey("The size should be 0", func() {
				size := h.size()
				So(size, ShouldEqual, 0)
			})

			Convey("The maximum node should be empty", func() {
				max := h.maximum()
				So(max.value, ShouldBeEmpty)
				So(max.length, ShouldEqual, 0)
				So(max.srcRowNo, ShouldEqual, 0)
			})

			Convey("Extrace Max should be fine", func() {
				So(func() { h.extractMax() }, ShouldNotPanic)
			})
		})

		Convey("When insert 10 nodes", func() {
			for i := 0; i < 10; i++ {
				h.insert([]float64{}, float64(i), i)
			}
			max1 := h.maximum()
			h.extractMax()
			h.extractMax()
			h.extractMax()
			max2 := h.maximum()

			Convey("The max1.length should be 9", func() {
				So(max1.length, ShouldEqual, 9)
			})
			Convey("The max2.length should be 6", func() {
				So(max2.length, ShouldEqual, 6)
			})
		})

	})
}
