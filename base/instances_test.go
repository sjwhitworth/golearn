package base

import (
	"fmt"
	"testing"
)

func TestIntAttribute(t *testing.T) {

	attrs := make([]Attribute, 0)
	intAttr1 := NewIntAttribute()
	intAttr2 := NewIntAttribute()
	intAttr1.SetName("Attribute 1")
	intAttr2.SetName("Attribute 2")
	attrs = append(attrs, intAttr1)
	attrs = append(attrs, intAttr2)

	instances := NewInstances(attrs, 5)
	for i := 0; i < 5; i++ {
		instances.Set(i, 0, intAttr1.GetSysValFromString(fmt.Sprintf("%d", i)))
		instances.Set(i, 1, intAttr2.GetSysValFromString(fmt.Sprintf("%d", i%2)))
	}

	fmt.Println(instances)

	row1 := instances.RowStr(0)
	row2 := instances.RowStr(1)
	row3 := instances.RowStr(2)
	if row1 != "0 0" {
		t.Error(row1)
	}
	if row2 != "1 1" {
		t.Error(row2)
	}
	if row3 != "2 0" {
		t.Error(row3)
	}
}
