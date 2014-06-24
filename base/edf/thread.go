package edf

import (
	"fmt"
)

// Threads are streams of data encapsulated within the file
type Thread struct {
	name string
	id   uint32
}

// NewThread returns a new thread
func NewThread(e *EdfFile, name string) *Thread {
	return &Thread{name, e.GetThreadCount() + 1}
}

// GetSpaceNeeded the number of bytes needed to serialize this
// Thread.
func (t *Thread) GetSpaceNeeded() int {
	return 8 + len(t.name)
}

// Serialize copies this thread to the output byte slice
// Returns the number of bytes used
func (t *Thread) Serialize(out []byte) int {
	// ret keeps track of written bytes
	ret := 0
	// Write the length of the name first
	nameLength := len(t.name)
	uint32ToBytes(uint32(nameLength), out)
	out = out[4:]
	ret += 4
	// Then write the string
	copy(out, t.name)
	out = out[nameLength:]
	ret += nameLength
	// Then the thread number
	uint32ToBytes(t.id, out)
	ret += 4
	return ret
}

// Deserialize copies the input byte slice into a thread
func (t *Thread) Deserialize(out []byte) int {
	ret := 0
	// Read the length of the thread's name
	nameLength := uint32FromBytes(out)
	ret += 4
	out = out[4:]
	// Copy out the string
	t.name = string(out[:nameLength])
	ret += int(nameLength)
	out = out[nameLength:]
	// Read the identifier
	t.id = uint32FromBytes(out)
	ret += 4
	return ret
}

// FindThread obtains the index of a thread in the EdfFile
func (e *EdfFile) FindThread(targetName string) (uint32, error) {

	var offset uint32
	var counter uint32
	// Resolve the initial thread block
	blockRange := e.GetPageRange(1, 1)
	if blockRange.Start.Segment != blockRange.End.Segment {
		return 0, fmt.Errorf("Thread block split across segments!")
	}
	bytes := e.m[blockRange.Start.Segment][blockRange.Start.Byte:blockRange.End.Byte]
	// Skip the first 8 bytes, since we don't support multiple thread blocks yet
	// TODO: fix that
	bytes = bytes[8:]
	counter = 1
	for {
		length := uint32FromBytes(bytes)
		if length == 0 {
			return 0, fmt.Errorf("No matching threads")
		}

		name := string(bytes[4 : 4+length])
		if name == targetName {
			offset = counter
			break
		}

		bytes = bytes[8+length:]
		counter++
	}

	return offset, nil
}

// WriteThread inserts a new thread into the EdfFile
func (e *EdfFile) WriteThread(t *Thread) error {

	offset, _ := e.FindThread(t.name)
	if offset != 0 {
		return fmt.Errorf("Writing a duplicate thread")
	}
	// Resolve the initial Thread block
	blockRange := e.GetPageRange(1, 1)
	if blockRange.Start.Segment != blockRange.End.Segment {
		return fmt.Errorf("Thread block split across segments!")
	}
	bytes := e.m[blockRange.Start.Segment][blockRange.Start.Byte:blockRange.End.Byte]
	// Skip the first 8 bytes, since we don't support multiple thread blocks yet
	// TODO: fix that
	bytes = bytes[8:]
	cur := 0
	for {
		length := uint32FromBytes(bytes)
		if length == 0 {
			break
		}
		cur += 8 + int(length)
		bytes = bytes[8+length:]
	}
	// cur should have now found an empty offset
	// Check that we have enough room left to insert
	roomLeft := len(bytes)
	roomNeeded := t.GetSpaceNeeded()
	if roomLeft < roomNeeded {
		return fmt.Errorf("Not enough space available")
	}
	// If everything's fine, serialise
	t.Serialize(bytes)
	// Increment thread count
	e.incrementThreadCount()
	return nil
}

// GetId returns this Thread's identifier
func (t *Thread) GetId() uint32 {
	return t.id
}
