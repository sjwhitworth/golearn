package base

import (
	"bytes"
	"fmt"
)

// Ponds contain a particular number of rows of
// a particular number of Attributes, all of a given type.
type Pond struct {
	threadNo   uint32
	parent     DataGrid
	attributes []Attribute
	size       int
	alloc      [][]byte
	maxRow     int
}

func (p *Pond) String() string {
	if len(p.alloc) > 1 {
		return fmt.Sprintf("Pond(%d attributes\n thread: %d\n size: %d\n)", len(p.attributes), p.threadNo, p.size)
	}
	return fmt.Sprintf("Pond(%d attributes\n thread: %d\n size: %d\n %d \n)", len(p.attributes), p.threadNo, p.size, p.alloc[0][0:60])
}

// PondStorageRef is a reference to a particular set
// of allocated rows within a Pond
type PondStorageRef struct {
	Storage []byte
	Rows    int
}

// RowSize returns the size of each row in bytes
func (p *Pond) RowSize() int {
	return len(p.attributes) * p.size
}

// Attributes returns a slice of Attributes in this Pond
func (p *Pond) Attributes() []Attribute {
	return p.attributes
}

// Storage returns a slice of PondStorageRefs which can
// be used to access the memory in this pond.
func (p *Pond) Storage() []PondStorageRef {
	ret := make([]PondStorageRef, len(p.alloc))
	rowSize := p.RowSize()
	for i, b := range p.alloc {
		ret[i] = PondStorageRef{b, len(b)/rowSize}
	}
	return ret
}

func (p *Pond) resolveBlock(col int, row int) (int, int) {

	if len(p.alloc) == 0 {
		panic("No blocks to resolve")
	}

	// Find where in the pond the byte is
	byteOffset := row*p.RowSize() + col*p.size
	curOffset := 0
	curBlock := 0
	blockOffset := 0
	for {
		if curBlock >= len(p.alloc) {
			panic("Don't have enough blocks to fulfill")
		}

		// Rows are not allowed to span blocks
		blockAdd := len(p.alloc[curBlock])
		blockAdd -= blockAdd % p.RowSize()

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

func (p *Pond) set(col int, row int, val []byte) {

	// Double-check the length
	if len(val) != p.size {
		panic(fmt.Sprintf("Tried to call set() with %d bytes, should be %d", len(val), p.size))
	}

	// Find where in the pond the byte is
	curBlock, blockOffset := p.resolveBlock(col, row)

	// Copy the value in
	copied := copy(p.alloc[curBlock][blockOffset:], val)
	if copied != p.size {
		panic(fmt.Sprintf("set() terminated by only copying %d bytes into the current block (should be %d). Check EDF allocation", copied, p.size))
	}

	row++
	if row > p.maxRow {
		p.maxRow = row
	}
}

func (p *Pond) get(col int, row int) []byte {
	curBlock, blockOffset := p.resolveBlock(col, row)
	return p.alloc[curBlock][blockOffset : blockOffset+p.size]
}

func (p *Pond) appendToRowBuf(row int, buffer *bytes.Buffer) {
	for i, a := range p.attributes {
		postfix := " "
		if i == len(p.attributes)-1 {
			postfix = ""
		}
		buffer.WriteString(fmt.Sprintf("%s%s", a.GetStringFromSysVal(p.get(i, row)), postfix))
	}
}