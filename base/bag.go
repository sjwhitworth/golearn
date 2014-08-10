package base

import (
	"fmt"
)

// BinaryAttributeGroups contain only BinaryAttributes
// Compact each Attribute to a bit for better storage
type BinaryAttributeGroup struct {
	FixedAttributeGroup
}

func (b *BinaryAttributeGroup) RowSize() int {
	return (len(b.attributes) + 7) / 8
}

// String gets a human-readable view of this group
func (b *BinaryAttributeGroup) String() string {
	if len(b.alloc) > 1 {
		return fmt.Sprintf("BinaryAttributeGroup(%d attributes\n thread: %d\n size: %d\n)", len(b.attributes), b.threadNo, b.size)
	}
	return fmt.Sprintf("BinaryAttributeGroup(%d attributes\n thread: %d\n size: %d\n)", len(b.attributes), b.threadNo, b.size)
}

func (b *BinaryAttributeGroup) getByteOffset(col, row int) int {
	return row*b.RowSize() + col/8
}

func (b *BinaryAttributeGroup) set(col, row int, val []byte) {
	// Resolve the block
	curBlock, blockOffset := b.resolveBlock(col, row)

	// If the value is 1, OR it
	if val[0] > 0 {
		b.alloc[curBlock][blockOffset] |= (1 << (uint(col) % 8))
	} else {
		// Otherwise, AND its complement
		b.alloc[curBlock][blockOffset] &= ^(1 << (uint(col) % 8))
	}

	row++
	if row > b.maxRow {
		b.maxRow = row
	}
}

func (b *BinaryAttributeGroup) resolveBlock(col, row int) (int, int) {

	byteOffset := row*b.RowSize() + (col / 3)
	rowSize := b.RowSize()
	return b.FixedAttributeGroup.resolveBlockFromByteOffset(byteOffset, rowSize)

}

func (b *BinaryAttributeGroup) get(col, row int) []byte {
	curBlock, blockOffset := b.resolveBlock(col, row)
	if b.alloc[curBlock][blockOffset]&(1<<(uint(col%8))) > 0 {
		return []byte{1}
	} else {
		return []byte{0}
	}
}
