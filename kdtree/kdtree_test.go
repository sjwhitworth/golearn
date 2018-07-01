package kdtree

import (
	"testing"

	"github.com/sjwhitworth/golearn/metrics/pairwise"
	. "github.com/smartystreets/goconvey/convey"
)

func TestKdtree(t *testing.T) {
	Convey("Test Build", t, func() {
		Convey("When no input data", func() {
			kd := New()
			data := [][]float64{}
			err := kd.Build(data)
			So(err.Error(), ShouldEqual, "no input data")
		})

		Convey("When amounts of features not the same", func() {
			kd := New()
			data := [][]float64{{3, 5}, {6, 7, 10}}
			err := kd.Build(data)
			So(err.Error(), ShouldEqual, "amounts of features are not the same")
		})

		Convey("When only one data", func() {
			kd := New()
			data := [][]float64{{3, 5}}
			err := kd.Build(data)
			So(err, ShouldBeNil)
		})

		Convey("When data all the same", func() {
			kd := New()
			data := [][]float64{{3, 5}, {3, 5}, {3, 5}}
			err := kd.Build(data)
			So(err, ShouldBeNil)
		})
	})

	Convey("Test Search", t, func() {
		Convey("Functionally test", func() {
			kd := New()
			data := [][]float64{{2, 3}, {5, 4}, {4, 7}, {8, 1}, {7, 2}, {9, 6}}
			kd.Build(data)
			euclidean := pairwise.NewEuclidean()

			Convey("When k is 3 with euclidean", func() {
				result, _, _ := kd.Search(3, euclidean, []float64{7, 3})

				Convey("The result[0] should be 4", func() {
					So(result[0], ShouldEqual, 4)
				})
				Convey("The result[1] should be 3", func() {
					So(result[1], ShouldEqual, 3)
				})
				Convey("The result[2] should be 1", func() {
					So(result[2], ShouldEqual, 1)
				})
			})

			Convey("When k is 2 with euclidean", func() {
				result, _, _ := kd.Search(2, euclidean, []float64{7, 3})

				Convey("The result[0] should be 4", func() {
					So(result[0], ShouldEqual, 4)
				})
				Convey("The result[1] should be 1", func() {
					So(result[1], ShouldEqual, 1)
				})
			})
		})

		Convey("When k is larger than amount of trainData", func() {
			kd := New()
			data := [][]float64{{3, 5}, {2, 1}}
			kd.Build(data)
			euclidean := pairwise.NewEuclidean()
			_, _, err := kd.Search(3, euclidean, []float64{7, 3})
			So(err.Error(), ShouldEqual, "k is largerer than amount of trainData")
		})

		Convey("When features of target is larger than trainData", func() {
			kd := New()
			data := [][]float64{{3, 5}, {2, 1}}
			kd.Build(data)
			euclidean := pairwise.NewEuclidean()
			_, _, err := kd.Search(1, euclidean, []float64{7, 3, 5})
			So(err.Error(), ShouldEqual, "amount of features is not equal")
		})

		Convey("When node.feature is -2", func() {
			kd := New()
			data := [][]float64{{3, 5}, {2, 1}}
			kd.Build(data)
			euclidean := pairwise.NewEuclidean()
			_, _, err := kd.Search(1, euclidean, []float64{7, 3})
			So(err, ShouldBeNil)
		})

		Convey("Search All Node (left)", func() {
			kd := New()
			data := [][]float64{{1, 2}, {5, 6}, {9, 10}}
			kd.Build(data)
			euclidean := pairwise.NewEuclidean()
			result, _, _ := kd.Search(1, euclidean, []float64{7, 3})
			So(result[0], ShouldEqual, 1)
		})

		Convey("Search when node length larger than heap max", func() {
			Convey("Search All Node (left)", func() {
				kd := New()
				data := [][]float64{{1, 2}, {5, 6}, {9, 10}}
				kd.Build(data)
				euclidean := pairwise.NewEuclidean()
				result, _, _ := kd.Search(1, euclidean, []float64{8, 7})
				So(result[0], ShouldEqual, 2)
			})

			Convey("Search All Node (right)", func() {
				kd := New()
				data := [][]float64{{1, 2}, {5, 4}, {9, 10}}
				kd.Build(data)
				euclidean := pairwise.NewEuclidean()
				result, _, _ := kd.Search(1, euclidean, []float64{3, 3})
				So(result[0], ShouldEqual, 0)
			})
		})
	})
}
