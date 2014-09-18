package base

import (
	"bytes"
)

// AttributeGroups store related sequences of system values
// in memory for the DenseInstances structure.
type AttributeGroup interface {
	addStorage(a []byte)
	// Used for printing
	appendToRowBuf(row int, buffer *bytes.Buffer)
	// Adds a new Attribute
	AddAttribute(Attribute) error
	// Returns all Attributes
	Attributes() []Attribute
	// Gets the byte slice at a given column, row offset
	get(int, int) []byte
	// Stores the byte slice at a given column, row offset
	set(int, int, []byte)
	// Gets the size of each row in bytes (rounded up)
	RowSize() int
	// Gets references to underlying memory
	Storage() []AttributeGroupStorageRef
	// Returns a human-readable summary
	String() string
}

// AttributeGroupStorageRef is a reference to a particular set
// of allocated rows within a FixedAttributeGroup
type AttributeGroupStorageRef struct {
	Storage []byte
	Rows    int
}
