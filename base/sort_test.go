package base

import "testing"

func isSortedAsc(inst FixedDataGrid, attr AttributeSpec) bool {
	valPrev := 0.0
	_, rows := inst.Size()
	for i := 0; i < rows; i++ {
		cur := UnpackBytesToFloat(inst.Get(attr, i))
		if i > 0 {
			if valPrev > cur {
				return false
			}
		}
		valPrev = cur
	}
	return true
}

func isSortedDesc(inst FixedDataGrid, attr AttributeSpec) bool {
	valPrev := 0.0
	_, rows := inst.Size()
	for i := 0; i < rows; i++ {
		cur := UnpackBytesToFloat(inst.Get(attr, i))
		if i > 0 {
			if valPrev < cur {
				return false
			}
		}
		valPrev = cur
	}
	return true
}

func TestSortDesc(testEnv *testing.T) {
	inst1, err := ParseCSVToInstances("../examples/datasets/iris_headers.csv", true)
	if err != nil {
		testEnv.Error(err)
		return
	}
	inst2, err := ParseCSVToInstances("../examples/datasets/iris_sorted_desc.csv", true)
	if err != nil {
		testEnv.Error(err)
		return
	}

	as1 := ResolveAllAttributes(inst1)
	as2 := ResolveAllAttributes(inst2)

	if isSortedDesc(inst1, as1[0]) {
		testEnv.Error("Can't test descending sort order")
	}
	if !isSortedDesc(inst2, as2[0]) {
		testEnv.Error("Reference data not sorted in descending order!")
	}

	Sort(inst1, Descending, as1[0:len(as1)-1])
	if err != nil {
		testEnv.Error(err)
	}
	if !isSortedDesc(inst1, as1[0]) {
		testEnv.Error("Instances are not sorted in descending order")
		testEnv.Error(inst1)
	}
	if !inst2.Equal(inst1) {
		testEnv.Error("Instances don't match")
		testEnv.Error(inst1)
		testEnv.Error(inst2)
	}
}

func TestSortAsc(testEnv *testing.T) {
	inst, err := ParseCSVToInstances("../examples/datasets/iris_headers.csv", true)
	as1 := ResolveAllAttributes(inst)
	if isSortedAsc(inst, as1[0]) {
		testEnv.Error("Can't test ascending sort on something ascending already")
	}
	if err != nil {
		testEnv.Error(err)
		return
	}
	Sort(inst, Ascending, as1[0:1])
	if !isSortedAsc(inst, as1[0]) {
		testEnv.Error("Instances are not sorted in ascending order")
		testEnv.Error(inst)
	}

	inst2, err := ParseCSVToInstances("../examples/datasets/iris_sorted_asc.csv", true)
	if err != nil {
		testEnv.Error(err)
		return
	}
	as2 := ResolveAllAttributes(inst2)
	if !isSortedAsc(inst2, as2[0]) {
		testEnv.Error("This file should be sorted in ascending order")
	}

	if !inst2.Equal(inst) {
		testEnv.Error("Instances don't match")
		testEnv.Error(inst)
		testEnv.Error(inst2)
	}

}
