package clustering

import (
	"bufio"
	"github.com/sjwhitworth/golearn/base"
	"github.com/sjwhitworth/golearn/metrics/pairwise"
	. "github.com/smartystreets/goconvey/convey"
	"gonum.org/v1/gonum/mat"
	"math"
	"math/big"
	"os"
	"strconv"
	"testing"
)

func TestDBSCANDistanceQuery(t *testing.T) {

	Convey("Should be able to determine which points are in range...", t, func() {

		// Read in the synthetic test data
		inst, err := base.ParseCSVToInstances("synthetic.csv", false)
		So(err, ShouldBeNil)

		// Create a neighbours vector
		neighbours := big.NewInt(0)

		// Compute pairwise distances
		dist, err := computePairwiseDistances(inst, inst.AllAttributes(), pairwise.NewEuclidean())
		So(dist.At(0, 0), ShouldAlmostEqual, 0)
		So(dist.At(0, 1), ShouldAlmostEqual, 1)
		So(dist.At(1, 0), ShouldAlmostEqual, 1)
		So(dist.At(0, 2), ShouldAlmostEqual, math.Sqrt(5))
		So(dist.At(2, 0), ShouldAlmostEqual, math.Sqrt(5))
		So(err, ShouldBeNil)

		// Do the region query
		neighbours = regionQuery(0, neighbours, dist, 1)
		So(neighbours.Bit(0), ShouldEqual, 1)
		So(neighbours.Bit(1), ShouldEqual, 1)
		So(neighbours.Bit(2), ShouldEqual, 0)
		So(neighbours.Bit(3), ShouldEqual, 0)
		So(neighbours.Bit(4), ShouldEqual, 0)

	})

}

func TestDBSCANSynthetic(t *testing.T) {
	Convey("Synthetic DBSCAN test should work...", t, func() {

		inst, err := base.ParseCSVToInstances("synthetic.csv", false)
		So(err, ShouldBeNil)

		p := DBSCANParameters{
			ClusterParameters{
				inst.AllAttributes(),
				pairwise.NewEuclidean(),
			},
			1,
			1,
		}

		m, err := DBSCAN(inst, p)
		So(err, ShouldBeNil)

		So(len(m), ShouldEqual, 2)
		So(m[1], ShouldContain, 0)
		So(m[1], ShouldContain, 1)
		So(m[1], ShouldContain, 2)
		So(m[1], ShouldContain, 3)

	})
}

func TestDBSCANDistanceMetric(t *testing.T) {

	Convey("Check the distance function is sane...", t, func() {

		d1 := mat.NewDense(1, 2, nil)
		d2 := mat.NewDense(1, 2, nil)

		d1.Set(0, 0, 0.494260967249)
		d1.Set(0, 1, 1.45106696541)
		d2.Set(0, 0, -1.42808099324)
		d2.Set(0, 1, -0.83706376669)

		e := pairwise.NewEuclidean()
		So(e.Distance(d1, d2), ShouldAlmostEqual, 2.9882, 0.001)

	})

}

func TestDBSCAN(t *testing.T) {

	Convey("Loading some data and labels...", t, func() {

		inst, err := base.ParseCSVToInstances("dbscan.csv", false)
		So(err, ShouldBeNil)

		file, err := os.Open("dbscan_labels.csv")
		defer file.Close()
		So(err, ShouldBeNil)

		clusterMap := ClusterMap(make(map[int][]int))

		scanner := bufio.NewScanner(file)
		line := -1
		for scanner.Scan() {
			line = line + 1
			v, err := strconv.ParseInt(scanner.Text(), 10, 64)
			if err != nil {
				panic(err)
			}
			v = v + 1 // -1 are noise in scikit-learn's DBSCAN
			c := int(v)
			if c == 0 {
				continue
			}
			if _, ok := clusterMap[c]; !ok {
				clusterMap[c] = make([]int, 0)
			}
			clusterMap[c] = append(clusterMap[c], line)
		}

		Convey("Our DBSCAN implementation should match...", func() {
			p := DBSCANParameters{
				ClusterParameters{
					inst.AllAttributes(),
					pairwise.NewEuclidean(),
				},
				0.3,
				10,
			}
			m, err := DBSCAN(inst, p)
			Convey("There should be nothing in the result that's smaller than MinPts", func() {

				for id := range m {
					So(len(m[id]), ShouldBeGreaterThanOrEqualTo, 10)
				}

			})
			So(err, ShouldBeNil)
			eq, err := clusterMap.Equals(m)
			So(err, ShouldBeNil)
			So(eq, ShouldBeTrue)
		})
	})

}
