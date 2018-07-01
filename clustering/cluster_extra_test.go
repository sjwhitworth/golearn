package clustering

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func Test(t *testing.T) {
	Convey("Only m[0]", t, func() {
		m1 := ClusterMap(make(map[int][]int))
		m1[0] = []int{1, 2}

		m2 := ClusterMap(make(map[int][]int))
		m2[0] = []int{1, 2}

		ret, err := m1.Equals(m2)
		So(err, ShouldBeNil)
		So(ret, ShouldBeTrue)

	})

	Convey("Nothing in m", t, func() {
		m1 := ClusterMap(make(map[int][]int))

		m2 := ClusterMap(make(map[int][]int))

		ret, err := m1.Equals(m2)
		So(err, ShouldBeNil)
		So(ret, ShouldBeTrue)

	})

	Convey("Many elements in m", t, func() {
		m1 := ClusterMap(make(map[int][]int))
		m1[0] = []int{1, 2, 3, 4, 5}
		m1[1] = []int{11, 12, 13, 14, 15}

		m2 := ClusterMap(make(map[int][]int))
		m2[0] = []int{1, 2, 3, 4, 5}
		m2[1] = []int{11, 12, 13, 14, 15}

		ret, err := m1.Equals(m2)
		So(err, ShouldBeNil)
		So(ret, ShouldBeTrue)

	})

	Convey("m[0] not the same", t, func() {
		m1 := ClusterMap(make(map[int][]int))
		m1[1] = []int{1, 2, 3}
		m1[0] = []int{4, 5}

		m2 := ClusterMap(make(map[int][]int))
		m2[1] = []int{1, 2, 3}
		m2[0] = []int{6, 5}

		_, err := m1.Equals(m2)
		So(err, ShouldNotBeNil)
	})

	Convey("m[0] size diff", t, func() {
		m1 := ClusterMap(make(map[int][]int))
		m1[1] = []int{1, 2, 3}
		m1[0] = []int{4, 5}

		m2 := ClusterMap(make(map[int][]int))
		m2[1] = []int{1, 2, 3}

		_, err := m1.Equals(m2)
		So(err, ShouldNotBeNil)
	})

	Convey("m[1] size diff", t, func() {
		m1 := ClusterMap(make(map[int][]int))
		m1[1] = []int{1, 3}
		m1[0] = []int{4, 5}

		m2 := ClusterMap(make(map[int][]int))
		m2[1] = []int{1, 2, 3}
		m1[0] = []int{4, 5}

		_, err := m1.Equals(m2)
		So(err, ShouldNotBeNil)
	})

	Convey("m[1] duplicate", t, func() {
		m1 := ClusterMap(make(map[int][]int))
		m1[1] = []int{1, 1}
		m1[0] = []int{4, 5}

		m2 := ClusterMap(make(map[int][]int))
		m1[1] = []int{1, 1}
		m1[0] = []int{4, 5}

		_, err := m1.Equals(m2)
		So(err, ShouldNotBeNil)
	})

	Convey("m[0] duplicate", t, func() {
		m1 := ClusterMap(make(map[int][]int))
		m1[1] = []int{1, 2}
		m1[0] = []int{4, 4}

		m2 := ClusterMap(make(map[int][]int))
		m1[1] = []int{1, 2}
		m1[0] = []int{4, 4}

		_, err := m1.Equals(m2)
		So(err, ShouldNotBeNil)
	})

}
