package clustering

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestClusterEquality(t *testing.T) {

	Convey("Should be able to determine if two cluster maps represent the same thing...", t, func() {

		Convey("When everything's exactly the same...", func() {

			m1 := ClusterMap(make(map[int][]int))
			m1[0] = []int{1, 2, 3}
			m1[1] = []int{4, 5}

			m2 := ClusterMap(make(map[int][]int))
			m2[0] = []int{1, 2, 3}
			m2[1] = []int{4, 5}

			ret, err := m1.Equals(m2)
			So(err, ShouldBeNil)
			So(ret, ShouldBeTrue)

		})

		Convey("With re-labelled clusters...", func() {
			m1 := ClusterMap(make(map[int][]int))
			m1[1] = []int{1, 2, 3}
			m1[0] = []int{4, 5}

			m2 := ClusterMap(make(map[int][]int))
			m2[1] = []int{1, 2, 3}
			m2[0] = []int{4, 5}

			ret, err := m1.Equals(m2)
			So(err, ShouldBeNil)
			So(ret, ShouldBeTrue)
		})

		Convey("With missing clusters...", func() {
			m1 := ClusterMap(make(map[int][]int))
			m1[1] = []int{1, 2, 3}

			m2 := ClusterMap(make(map[int][]int))
			m2[1] = []int{1, 2, 3}
			m2[0] = []int{4, 5}

			_, err := m1.Equals(m2)
			So(err, ShouldNotBeNil)
		})

		Convey("With missing points...", func() {
			m1 := ClusterMap(make(map[int][]int))
			m1[1] = []int{1, 3}
			m1[0] = []int{4, 5}

			m2 := ClusterMap(make(map[int][]int))
			m2[1] = []int{1, 2, 3}
			m2[0] = []int{4, 5}

			_, err := m1.Equals(m2)
			So(err, ShouldNotBeNil)
		})

		Convey("With invalid maps...", func() {
			m1 := ClusterMap(make(map[int][]int))
			m1[0] = []int{1, 2, 3}
			m1[1] = []int{4, 4, 5}

			m2 := ClusterMap(make(map[int][]int))
			m2[0] = []int{1, 2, 3}
			m2[1] = []int{4, 5}

			_, err := m1.Equals(m2)
			So(err, ShouldNotBeNil)
		})

	})

}
