package kdtree

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestHeap(t *testing.T) {
	h := newHeap()

	Convey("Given a heap", t, func() {

		Convey("When heap is empty", func() {
			_, err := h.maximum()

			Convey("The err should be errEmpty", func() {
				So(err, ShouldEqual, h.errEmpty())
			})
		})

		Convey("When insert 5 nodes", func() {
			for i := 0; i < 5; i++ {
				h.insert([]float64{}, float64(i))
			}
			max1, _ := h.maximum()
			h.extractMax()
			max2, _ := h.maximum()

			Convey("The max1.value should be 4", func() {
				So(max1.value, ShouldEqual, 4)
			})
			Convey("The max2.value should be 3", func() {
				So(max2.value, ShouldEqual, 3)
			})

		})

	})
}
