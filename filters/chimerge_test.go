package filters

import (
	"github.com/sjwhitworth/golearn/base"
	"math"
	"testing"
)

func TestChiMFreqTable(t *testing.T) {
	inst, err := base.ParseCSVToInstances("../examples/datasets/chim.csv", true)
	if err != nil {
		t.Fatal("Unable to parse CSV to instances: %s", err.Error())
	}

	freq := ChiMBuildFrequencyTable(inst.AllAttributes()[0], inst)

	if freq[0].Frequency["c1"] != 1 {
		t.Error("Wrong frequency")
	}
	if freq[0].Frequency["c3"] != 4 {
		t.Errorf("Wrong frequency %s", freq[1])
	}
	if freq[10].Frequency["c2"] != 1 {
		t.Error("Wrong frequency")
	}
}

func TestChiClassCounter(t *testing.T) {
	inst, err := base.ParseCSVToInstances("../examples/datasets/chim.csv", true)
	if err != nil {
		t.Fatal("Unable to parse CSV to instances: %s", err.Error())
	}
	freq := ChiMBuildFrequencyTable(inst.AllAttributes()[0], inst)
	classes := chiCountClasses(freq)
	if classes["c1"] != 27 {
		t.Error(classes)
	}
	if classes["c2"] != 12 {
		t.Error(classes)
	}
	if classes["c3"] != 21 {
		t.Error(classes)
	}
}

func TestStatisticValues(t *testing.T) {
	inst, err := base.ParseCSVToInstances("../examples/datasets/chim.csv", true)
	if err != nil {
		t.Fatal("Unable to parse CSV to instances: %s", err.Error())
	}
	freq := ChiMBuildFrequencyTable(inst.AllAttributes()[0], inst)
	chiVal := chiComputeStatistic(freq[5], freq[6])
	if math.Abs(chiVal-1.89) > 0.01 {
		t.Error(chiVal)
	}

	chiVal = chiComputeStatistic(freq[1], freq[2])
	if math.Abs(chiVal-1.08) > 0.01 {
		t.Error(chiVal)
	}
}

func TestChiSquareDistValues(t *testing.T) {
	chiVal1 := chiSquaredPercentile(2, 4.61)
	chiVal2 := chiSquaredPercentile(3, 7.82)
	chiVal3 := chiSquaredPercentile(4, 13.28)
	if math.Abs(chiVal1-0.90) > 0.001 {
		t.Error(chiVal1)
	}
	if math.Abs(chiVal2-0.95) > 0.001 {
		t.Error(chiVal2)
	}
	if math.Abs(chiVal3-0.99) > 0.001 {
		t.Error(chiVal3)
	}
}

func TestChiMerge1(t *testing.T) {
	inst, err := base.ParseCSVToInstances("../examples/datasets/chim.csv", true)
	if err != nil {
		t.Fatal("Unable to parse CSV to instances: %s", err.Error())
	}
	_, rows := inst.Size()

	freq := chiMerge(inst, inst.AllAttributes()[0], 0.90, 0, rows)
	if len(freq) != 3 {
		t.Error("Wrong length")
	}
	if freq[0].Value != 1.3 {
		t.Error(freq[0])
	}
	if freq[1].Value != 56.2 {
		t.Error(freq[1])
	}
	if freq[2].Value != 87.1 {
		t.Error(freq[2])
	}
}

func TestChiMerge2(t *testing.T) {
	//
	// See http://sci2s.ugr.es/keel/pdf/algorithm/congreso/1992-Kerber-ChimErge-AAAI92.pdf
	//   Randy Kerber, ChiMerge: Discretisation of Numeric Attributes, 1992
	inst, err := base.ParseCSVToInstances("../examples/datasets/iris_headers.csv", true)
	if err != nil {
		t.Fatal("Unable to parse CSV to instances: %s", err.Error())
	}

	// Sort the instances
	allAttrs := inst.AllAttributes()
	sortAttrSpecs := base.ResolveAttributes(inst, allAttrs)[0:1]
	instSorted, err := base.Sort(inst, base.Ascending, sortAttrSpecs)
	if err != nil {
		t.Fatalf("Sort failed: %s", err.Error())
	}

	// Perform Chi-Merge
	_, rows := inst.Size()
	freq := chiMerge(instSorted, allAttrs[0], 0.90, 0, rows)
	if len(freq) != 5 {
		t.Errorf("Wrong length (%d)", len(freq))
		t.Error(freq)
	}
	if freq[0].Value != 4.3 {
		t.Error(freq[0])
	}
	if freq[1].Value != 5.5 {
		t.Error(freq[1])
	}
	if freq[2].Value != 5.8 {
		t.Error(freq[2])
	}
	if freq[3].Value != 6.3 {
		t.Error(freq[3])
	}
	if freq[4].Value != 7.1 {
		t.Error(freq[4])
	}
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

func TestChiMerge4(t *testing.T) {
	// See http://sci2s.ugr.es/keel/pdf/algorithm/congreso/1992-Kerber-ChimErge-AAAI92.pdf
	//   Randy Kerber, ChiMerge: Discretisation of Numeric Attributes, 1992
	inst, err := base.ParseCSVToInstances("../examples/datasets/iris_headers.csv", true)
	if err != nil {
		t.Fatal("Unable to parse CSV to instances: %s", err.Error())
	}

	filt := NewChiMergeFilter(inst, 0.90)
	filt.AddAttribute(inst.AllAttributes()[0])
	filt.AddAttribute(inst.AllAttributes()[1])
	filt.Train()
	instf := base.NewLazilyFilteredInstances(inst, filt)
	clsAttrs := instf.AllClassAttributes()
	if len(clsAttrs) != 1 {
		t.Fatalf("%d != %d", len(clsAttrs), 1)
	}
	firstClassAttributeName := clsAttrs[0].GetName()
	expectedClassAttributeName := "Species"
	if firstClassAttributeName != expectedClassAttributeName {
		t.Fatalf("Expected class attribute '%s'; actual class attribute '%s'", expectedClassAttributeName, firstClassAttributeName)
	}
}
