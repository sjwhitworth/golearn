package base

import (
	"bytes"
	"fmt"
)

// FixedAttributeGroups contain a particular number of rows of
// a particular number of Attributes, all of a given type.
type FixedAttributeGroup struct {
	parent     DataGrid
	attributes []Attribute
	size       int
	alloc      [][]byte
	maxRow     int
}

// String gets a human-readable summary
func (f *FixedAttributeGroup) String() string {
	return "FixedAttributeGroup"
}

// RowSize returns the size of each row in bytes
func (f *FixedAttributeGroup) RowSize() int {
	return len(f.attributes) * f.size
}

// Attributes returns a slice of Attributes in this FixedAttributeGroup
func (f *FixedAttributeGroup) Attributes() []Attribute {
	return f.attributes
}

// AddAttribute adds an attribute to this FixedAttributeGroup
func (f *FixedAttributeGroup) AddAttribute(a Attribute) error {
	f.attributes = append(f.attributes, a)
	return nil
}

// addStorage appends the given storage reference to this FixedAttributeGroup
func (f *FixedAttributeGroup) addStorage(a []byte) {
	f.alloc = append(f.alloc, a)
}

// Storage returns a slice of FixedAttributeGroupStorageRefs which can
// be used to access the memory in this pond.
func (f *FixedAttributeGroup) Storage() []AttributeGroupStorageRef {
	ret := make([]AttributeGroupStorageRef, len(f.alloc))
	rowSize := f.RowSize()
	for i, b := range f.alloc {
		ret[i] = AttributeGroupStorageRef{b, len(b) / rowSize}
	}
	return ret
}

func (f *FixedAttributeGroup) resolveBlock(col int, row int) (int, int) {

	if len(f.alloc) == 0 {
		panic("No blocks to resolve")
	}

	// Find where in the pond the byte is
	byteOffset := row*f.RowSize() + col*f.size
	return f.resolveBlockFromByteOffset(byteOffset, f.RowSize())
}

func (f *FixedAttributeGroup) resolveBlockFromByteOffset(byteOffset, rowSize int) (int, int) {
	curOffset := 0
	curBlock := 0
	blockOffset := 0
	for {
		if curBlock >= len(f.alloc) {
			panic("Don't have enough blocks to fulfill")
		}

		// Rows are not allowed to span blocks
		blockAdd := len(f.alloc[curBlock])
		blockAdd -= blockAdd % rowSize

		// Case 1: we need to skip this allocation
		if curOffset+blockAdd < byteOffset {
			curOffset += blockAdd
			curBlock++
		} else {
			blockOffset = byteOffset - curOffset
			break
		}
	}

	return curBlock, blockOffset
}

func (f *FixedAttributeGroup) set(col int, row int, val []byte) {

	// Double-check the length
	if len(val) != f.size {
		panic(fmt.Sprintf("Tried to call set() with %d bytes, should be %d", len(val), f.size))
	}

	// Find where in the pond the byte is
	curBlock, blockOffset := f.resolveBlock(col, row)

	// Copy the value in
	copied := copy(f.alloc[curBlock][blockOffset:], val)
	if copied != f.size {
		panic(fmt.Sprintf("set() terminated by only copying %d bytes into the current block (should be %d). Check EDF allocation", copied, f.size))
	}

	row++
	if row > f.maxRow {
		f.maxRow = row
	}
}

func (f *FixedAttributeGroup) get(col int, row int) []byte {
	curBlock, blockOffset := f.resolveBlock(col, row)
	return f.alloc[curBlock][blockOffset : blockOffset+f.size]
}

func (f *FixedAttributeGroup) appendToRowBuf(row int, buffer *bytes.Buffer) {
	for i, a := range f.attributes {
		postfix := " "
		if i == len(f.attributes)-1 {
			postfix = ""
		}
		buffer.WriteString(fmt.Sprintf("%s%s", a.GetStringFromSysVal(f.get(i, row)), postfix))
	}
}
