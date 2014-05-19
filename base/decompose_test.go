package base

import "testing"

func TestDecomp(testEnv *testing.T) {
	inst, err := ParseCSVToInstances("../examples/datasets/iris_binned.csv", true)
	if err != nil {
		testEnv.Error(err)
		return
	}
	decomp := inst.DecomposeOnAttributeValues(inst.GetAttr(0))

	row0 := decomp["0.00"].RowStr(0)
	row1 := decomp["1.00"].RowStr(0)
	/*	row2 := decomp["2.00"].RowStr(0)
		row3 := decomp["3.00"].RowStr(0)
		row4 := decomp["4.00"].RowStr(0)
		row5 := decomp["5.00"].RowStr(0)
		row6 := decomp["6.00"].RowStr(0)
		row7 := decomp["7.00"].RowStr(0)*/
	row8 := decomp["8.00"].RowStr(0)
	//	row9 := decomp["9.00"].RowStr(0)

	if row0 != "3.10 1.50 0.20 Iris-setosa" {
		testEnv.Error(row0)
	}
	if row1 != "3.00 1.40 0.20 Iris-setosa" {
		testEnv.Error(row1)
	}
	if row8 != "2.90 6.30 1.80 Iris-virginica" {
		testEnv.Error(row8)
	}
}
