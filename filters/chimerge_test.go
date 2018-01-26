package filters

import (
	"fmt"
	"testing"

	"github.com/amclay/golearn/base"

	. "github.com/smartystreets/goconvey/convey"
)

func TestChiMergeFrequencyTable(t *testing.T) {
	Convey("Chi-Merge Frequency Table", t, func() {
		instances, err := base.ParseCSVToInstances("../examples/datasets/chim.csv", true)
		So(err, ShouldBeNil)

		frequencyTable := ChiMBuildFrequencyTable(instances.AllAttributes()[0], instances)

		Convey("Computes frequencies correctly", func() {
			So(frequencyTable[0].Frequency["c1"], ShouldEqual, 1)
			So(frequencyTable[0].Frequency["c3"], ShouldEqual, 4)
			So(frequencyTable[10].Frequency["c2"], ShouldEqual, 1)
		})

		Convey("Counts classes correctly", func() {
			classes := chiCountClasses(frequencyTable)

			So(classes["c1"], ShouldEqual, 27)
			So(classes["c2"], ShouldEqual, 12)
			So(classes["c3"], ShouldEqual, 21)
		})

		Convey("Computes statistics correctly", func() {
			So(chiComputeStatistic(frequencyTable[5], frequencyTable[6]), ShouldAlmostEqual, 1.89, 0.01)
			So(chiComputeStatistic(frequencyTable[1], frequencyTable[2]), ShouldAlmostEqual, 1.08, 0.01)
		})
	})
}

func TestChiSquaredDistribution(t *testing.T) {
	Convey("Chi-Squared Distribution percentiles are computed correctly", t, func() {
		So(chiSquaredPercentile(2, 4.61), ShouldAlmostEqual, 0.9, 0.001)
		So(chiSquaredPercentile(3, 7.82), ShouldAlmostEqual, 0.95, 0.001)
		So(chiSquaredPercentile(4, 13.28), ShouldAlmostEqual, 0.99, 0.001)
	})
}

func TestChiMergeDiscretization(t *testing.T) {
	Convey("Chi-Merge Discretization", t, func() {
		chimDatasetPath := "../examples/datasets/chim.csv"

		Convey(fmt.Sprintf("With the '%s' dataset", chimDatasetPath), func() {
			instances, err := base.ParseCSVToInstances(chimDatasetPath, true)
			So(err, ShouldBeNil)

			_, rows := instances.Size()

			frequencies := chiMerge(instances, instances.AllAttributes()[0], 0.9, 0, rows)
			values := []float64{}
			for _, entry := range frequencies {
				values = append(values, entry.Value)
			}

			Convey("Computes frequencies correctly", func() {
				So(values, ShouldResemble, []float64{1.3, 56.2, 87.1})
			})
		})

		irisHeadersDatasetpath := "../examples/datasets/iris_headers.csv"

		Convey(fmt.Sprintf("With the '%s' dataset", irisHeadersDatasetpath), func() {
			instances, err := base.ParseCSVToInstances(irisHeadersDatasetpath, true)
			So(err, ShouldBeNil)

			Convey("Sorting the instances first", func() {
				allAttributes := instances.AllAttributes()
				sortedAttributesSpecs := base.ResolveAttributes(instances, allAttributes)[0:1]
				sortedInstances, err := base.Sort(instances, base.Ascending, sortedAttributesSpecs)
				So(err, ShouldBeNil)

				_, rows := sortedInstances.Size()

				frequencies := chiMerge(sortedInstances, sortedInstances.AllAttributes()[0], 0.9, 0, rows)
				values := []float64{}
				for _, entry := range frequencies {
					values = append(values, entry.Value)
				}

				Convey("Computes frequencies correctly", func() {
					So(values, ShouldResemble, []float64{4.3, 5.5, 5.8, 6.3, 7.1})
				})
			})
		})
	})
}

func TestChiMergeFilter(t *testing.T) {
	Convey("Chi-Merge Filter", t, func() {
		// See http://sci2s.ugr.es/keel/pdf/algorithm/congreso/1992-Kerber-ChimErge-AAAI92.pdf
		//   Randy Kerber, ChiMerge: Discretisation of Numeric Attributes, 1992
		instances, err := base.ParseCSVToInstances("../examples/datasets/iris_headers.csv", true)
		So(err, ShouldBeNil)

		Convey("Create and train the filter", func() {
			filter := NewChiMergeFilter(instances, 0.90)
			filter.AddAttribute(instances.AllAttributes()[0])
			filter.AddAttribute(instances.AllAttributes()[1])
			filter.Train()

			Convey("Filter the dataset", func() {
				filteredInstances := base.NewLazilyFilteredInstances(instances, filter)

				classAttributes := filteredInstances.AllClassAttributes()

				Convey("There should only be one class attribute", func() {
					So(len(classAttributes), ShouldEqual, 1)
				})

				expectedClassAttribute := "Species"

				Convey(fmt.Sprintf("The class attribute should be %s", expectedClassAttribute), func() {
					So(classAttributes[0].GetName(), ShouldEqual, expectedClassAttribute)
				})
			})
		})
	})
}

/*
func TestChiMerge3(t *testing.T) {
	// See http://sci2s.ugr.es/keel/pdf/algorithm/congreso/1992-Kerber-ChimErge-AAAI92.pdf
	//   Randy Kerber, ChiMerge: Discretisation of Numeric Attributes, 1992
	inst, err := base.ParseCSVToInstances("../examples/datasets/iris_headers.csv", true)
	if err != nil {
		panic(err)
	}

	insts, err := base.LazySort(inst, base.Ascending, base.ResolveAllAttributes(inst, inst.AllAttributes()))
	if err != nil {
		t.Error(err)
	}
	filt := NewChiMergeFilter(inst, 0.90)
	filt.AddAttribute(inst.AllAttributes()[0])
	filt.Train()
	instf := base.NewLazilyFilteredInstances(insts, filt)
	fmt.Println(instf)
	fmt.Println(instf.String())
	rowStr := instf.RowString(0)
	ref := "4.300000 3.00 1.10 0.10 Iris-setosa"
	if rowStr != ref {
		panic(fmt.Sprintf("'%s' != '%s'", rowStr, ref))
	}
	clsAttrs := instf.AllClassAttributes()
	if len(clsAttrs) != 1 {
		panic(fmt.Sprintf("%d != %d", len(clsAttrs), 1))
	}
	if clsAttrs[0].GetName() != "Species" {
		panic("Class Attribute wrong!")
	}
}
*/
