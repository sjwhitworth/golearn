package filters

import (
	"fmt"
	"github.com/sjwhitworth/golearn/base"
	"math"
	"testing"
)

func TestChiMFreqTable(testEnv *testing.T) {

	inst, err := base.ParseCSVToInstances("../examples/datasets/chim.csv", true)
	if err != nil {
		panic(err)
	}

	freq := ChiMBuildFrequencyTable(inst.AllAttributes()[0], inst)

	if freq[0].Frequency["c1"] != 1 {
		testEnv.Error("Wrong frequency")
	}
	if freq[0].Frequency["c3"] != 4 {
		testEnv.Errorf("Wrong frequency %s", freq[1])
	}
	if freq[10].Frequency["c2"] != 1 {
		testEnv.Error("Wrong frequency")
	}
}

func TestChiClassCounter(testEnv *testing.T) {
	inst, err := base.ParseCSVToInstances("../examples/datasets/chim.csv", true)
	if err != nil {
		panic(err)
	}
	freq := ChiMBuildFrequencyTable(inst.AllAttributes()[0], inst)
	classes := chiCountClasses(freq)
	if classes["c1"] != 27 {
		testEnv.Error(classes)
	}
	if classes["c2"] != 12 {
		testEnv.Error(classes)
	}
	if classes["c3"] != 21 {
		testEnv.Error(classes)
	}
}

func TestStatisticValues(testEnv *testing.T) {
	inst, err := base.ParseCSVToInstances("../examples/datasets/chim.csv", true)
	if err != nil {
		panic(err)
	}
	freq := ChiMBuildFrequencyTable(inst.AllAttributes()[0], inst)
	chiVal := chiComputeStatistic(freq[5], freq[6])
	if math.Abs(chiVal-1.89) > 0.01 {
		testEnv.Error(chiVal)
	}

	chiVal = chiComputeStatistic(freq[1], freq[2])
	if math.Abs(chiVal-1.08) > 0.01 {
		testEnv.Error(chiVal)
	}
}

func TestChiSquareDistValues(testEnv *testing.T) {
	chiVal1 := chiSquaredPercentile(2, 4.61)
	chiVal2 := chiSquaredPercentile(3, 7.82)
	chiVal3 := chiSquaredPercentile(4, 13.28)
	if math.Abs(chiVal1-0.90) > 0.001 {
		testEnv.Error(chiVal1)
	}
	if math.Abs(chiVal2-0.95) > 0.001 {
		testEnv.Error(chiVal2)
	}
	if math.Abs(chiVal3-0.99) > 0.001 {
		testEnv.Error(chiVal3)
	}
}

func TestChiMerge1(testEnv *testing.T) {

	// Read the data
	inst, err := base.ParseCSVToInstances("../examples/datasets/chim.csv", true)
	if err != nil {
		panic(err)
	}
	_, rows := inst.Size()

	freq := chiMerge(inst, inst.AllAttributes()[0], 0.90, 0, rows)
	if len(freq) != 3 {
		testEnv.Error("Wrong length")
	}
	if freq[0].Value != 1.3 {
		testEnv.Error(freq[0])
	}
	if freq[1].Value != 56.2 {
		testEnv.Error(freq[1])
	}
	if freq[2].Value != 87.1 {
		testEnv.Error(freq[2])
	}
}

func TestChiMerge2(testEnv *testing.T) {
	//
	// See http://sci2s.ugr.es/keel/pdf/algorithm/congreso/1992-Kerber-ChimErge-AAAI92.pdf
	//   Randy Kerber, ChiMerge: Discretisation of Numeric Attributes, 1992
	inst, err := base.ParseCSVToInstances("../examples/datasets/iris_headers.csv", true)
	if err != nil {
		panic(err)
	}

	// Sort the instances
	allAttrs := inst.AllAttributes()
	sortAttrSpecs := base.ResolveAttributes(inst, allAttrs)[0:1]
	instSorted, err := base.Sort(inst, base.Ascending, sortAttrSpecs)
	if err != nil {
		panic(err)
	}

	// Perform Chi-Merge
	_, rows := inst.Size()
	freq := chiMerge(instSorted, allAttrs[0], 0.90, 0, rows)
	if len(freq) != 5 {
		testEnv.Errorf("Wrong length (%d)", len(freq))
		testEnv.Error(freq)
	}
	if freq[0].Value != 4.3 {
		testEnv.Error(freq[0])
	}
	if freq[1].Value != 5.5 {
		testEnv.Error(freq[1])
	}
	if freq[2].Value != 5.8 {
		testEnv.Error(freq[2])
	}
	if freq[3].Value != 6.3 {
		testEnv.Error(freq[3])
	}
	if freq[4].Value != 7.1 {
		testEnv.Error(freq[4])
	}
}

/*
func TestChiMerge3(testEnv *testing.T) {
	// See http://sci2s.ugr.es/keel/pdf/algorithm/congreso/1992-Kerber-ChimErge-AAAI92.pdf
	//   Randy Kerber, ChiMerge: Discretisation of Numeric Attributes, 1992
	inst, err := base.ParseCSVToInstances("../examples/datasets/iris_headers.csv", true)
	if err != nil {
		panic(err)
	}

	insts, err := base.LazySort(inst, base.Ascending, base.ResolveAllAttributes(inst, inst.AllAttributes()))
	if err != nil {
		testEnv.Error(err)
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

func TestChiMerge4(testEnv *testing.T) {
	// See http://sci2s.ugr.es/keel/pdf/algorithm/congreso/1992-Kerber-ChimErge-AAAI92.pdf
	//   Randy Kerber, ChiMerge: Discretisation of Numeric Attributes, 1992
	inst, err := base.ParseCSVToInstances("../examples/datasets/iris_headers.csv", true)
	if err != nil {
		panic(err)
	}

	filt := NewChiMergeFilter(inst, 0.90)
	filt.AddAttribute(inst.AllAttributes()[0])
	filt.AddAttribute(inst.AllAttributes()[1])
	filt.Train()
	instf := base.NewLazilyFilteredInstances(inst, filt)
	fmt.Println(instf)
	fmt.Println(instf.String())
	clsAttrs := instf.AllClassAttributes()
	if len(clsAttrs) != 1 {
		panic(fmt.Sprintf("%d != %d", len(clsAttrs), 1))
	}
	if clsAttrs[0].GetName() != "Species" {
		panic("Class Attribute wrong!")
	}
}
