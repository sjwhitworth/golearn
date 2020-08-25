package trees

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestIsolation(t *testing.T) {

	Convey("Doing an Isolation Forest Test", t, func() {

		var data [][]float64
		data = append(data, []float64{8, 9, 8, 3})
		data = append(data, []float64{4, 2, 5, 3})
		data = append(data, []float64{3, 2, 5, 9})
		data = append(data, []float64{2, 1, 5, 9})

		featureChosen := selectFeature(data)
		So(featureChosen, ShouldNotBeNil)

		min, max := minMax(0, data)
		So(min, ShouldEqual, 2)
		So(max, ShouldEqual, 8)

		min, max = minMax(featureChosen, data)

		val := selectValue(min, max)
		So(val, ShouldBeBetween, min, max)

		leftData, rightData := splitData(val, featureChosen, data)
		So(len(leftData), ShouldBeGreaterThan, 0)
		So(len(rightData), ShouldBeGreaterThan, 0)

		checked := checkData(data)
		So(checked, ShouldBeTrue)

		randomSubset := getRandomData(data, 2)
		So(len(randomSubset), ShouldEqual, 2)

	})
}
