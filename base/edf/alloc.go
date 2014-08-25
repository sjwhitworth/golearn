package edf

import (
	"fmt"
)

// contentEntry structs are stored in contentEntry blocks
// which always at block 2.
type contentEntry struct {
	// Which thread this entry is assigned to
	thread uint32
	// Which page this block starts at
	Start uint32
	// The page up to and including which the block ends
	End uint32
}

func (e *EdfFile) extend(additionalPages uint32) error {
	fileInfo, err := e.f.Stat()
	if err != nil {
		panic(err)
	}
	newSize := uint64(fileInfo.Size())/e.pageSize + uint64(additionalPages)
	return e.truncate(int64(newSize))
}

func (e *EdfFile) getFreeMapSize() uint64 {
	if e.f != nil {
		fileInfo, err := e.f.Stat()
		if err != nil {
			panic(err)
		}
		return uint64(fileInfo.Size()) / e.pageSize
	}
	return uint64(EDF_SIZE) / e.pageSize
}

func (e *EdfFile) getContiguousOffset(pagesRequested uint32) (uint32, error) {
	// Create the free bitmap
	bitmap := make([]bool, e.getFreeMapSize())
	for i := 0; i < 4; i++ {
		bitmap[i] = true
	}
	// Traverse the contents table and build a free bitmap
	block := uint64(2)
	for {
		// Get the range for this block
		r := e.getPageRange(block, block)
		if r.Start.Segment != r.End.Segment {
			return 0, fmt.Errorf("Contents block split across segments")
		}
		bytes := e.m[r.Start.Segment]
		bytes = bytes[r.Start.Byte : r.End.Byte+1]
		// Get the address of the next contents block
		block = uint64FromBytes(bytes)
		if block != 0 {
			// No point in checking this block for free space
			continue
		}
		bytes = bytes[8:]
		// Look for a blank entry in the table
		for i := 0; i < len(bytes); i += 12 {
			threadID := uint32FromBytes(bytes[i:])
			if threadID == 0 {
				continue
			}
			start := uint32FromBytes(bytes[i+4:])
			end := uint32FromBytes(bytes[i+8:])
			for j := start; j <= end; j++ {
				if int(j) >= len(bitmap) {
					break
				}
				bitmap[j] = true
			}
		}
		break
	}
	// Look through the freemap and find a good spot
	for i := 0; i < len(bitmap); i++ {
		if bitmap[i] {
			continue
		}
		for j := i; j < len(bitmap); j++ {
			if !bitmap[j] {
				diff := j - 1 - i
				if diff > int(pagesRequested) {
					return uint32(i), nil
				}
			}
		}
	}
	return 0, nil
}

// addNewContentsBlock adds a new contents block in the next available space
func (e *EdfFile) addNewContentsBlock() error {

	var toc contentEntry

	// Find the next available offset
	startBlock, err := e.getContiguousOffset(1)
	if startBlock == 0 && err == nil {
		// Increase the size of the file if necessary
		e.extend(uint32(e.pageSize))
	} else if err != nil {
		return err
	}

	// Traverse the contents blocks looking for one with a blank NEXT pointer
	block := uint64(2)
	for {
		// Get the range for this block
		r := e.getPageRange(block, block)
		if r.Start.Segment != r.End.Segment {
			return fmt.Errorf("Contents block split across segments")
		}
		bytes := e.m[r.Start.Segment]
		bytes = bytes[r.Start.Byte : r.End.Byte+1]
		// Get the address of the next contents block
		block = uint64FromBytes(bytes)
		if block == 0 {
			uint64ToBytes(uint64(startBlock), bytes)
			break
		}
	}
	// Add to the next available TOC space
	toc.Start = startBlock
	toc.End = startBlock + 1
	toc.thread = 1 // SYSTEM thread
	return e.addToTOC(&toc, false)
}

// addToTOC adds a ContentsEntry structure in the next available place
func (e *EdfFile) addToTOC(c *contentEntry, extend bool) error {
	// Traverse the contents table looking for a free spot
	block := uint64(2)
	for {
		// Get the range for this block
		r := e.getPageRange(block, block)
		if r.Start.Segment != r.End.Segment {
			return fmt.Errorf("Contents block split across segments")
		}
		bytes := e.m[r.Start.Segment]
		bytes = bytes[r.Start.Byte : r.End.Byte+1]
		// Get the address of the next contents block
		block = uint64FromBytes(bytes)
		if block != 0 {
			// No point in checking this block for free space
			continue
		}
		bytes = bytes[8:]
		// Look for a blank entry in the table
		cur := 0
		for {
			threadID := uint32FromBytes(bytes)
			if threadID == 0 {
				break
			}
			cur += 12
			bytes = bytes[12:]
			if len(bytes) < 12 {
				if extend {
					// Append a new contents block and try again
					e.addNewContentsBlock()
					return e.addToTOC(c, false)
				}
				return fmt.Errorf("Can't add to contents: no space available")
			}
		}
		// Write the contents information into this block
		uint32ToBytes(c.thread, bytes)
		bytes = bytes[4:]
		uint32ToBytes(c.Start, bytes)
		bytes = bytes[4:]
		uint32ToBytes(c.End, bytes)
		break
	}
	return nil
}

// AllocPages allocates a |pagesRequested| chunk of pages on the thread
// with the given identifier. Returns an edfRange describing the result.
func (e *EdfFile) AllocPages(pagesRequested uint32, thread uint32) (edfRange, error) {

	var ret edfRange
	var toc contentEntry

	// Parameter check
	if pagesRequested == 0 {
		return ret, fmt.Errorf("Must request some pages")
	}
	if thread == 0 {
		return ret, fmt.Errorf("Need a valid page identifier")
	}

	// Find the next available offset
	startBlock, err := e.getContiguousOffset(pagesRequested)
	if startBlock == 0 && err == nil {
		// Increase the size of the file if necessary
		e.extend(pagesRequested)
		return e.AllocPages(pagesRequested, thread)
	} else if err != nil {
		return ret, err
	}

	// Add to the table of contents
	toc.thread = thread
	toc.Start = startBlock
	toc.End = startBlock + pagesRequested
	err = e.addToTOC(&toc, true)

	// Compute the range
	ret = e.getPageRange(uint64(startBlock), uint64(startBlock+pagesRequested))

	return ret, err
}

// getThreadBlocks returns EdfRanges containing blocks assigned to a given thread.
func (e *EdfFile) getThreadBlocks(thread uint32) ([]edfRange, error) {
	var ret []edfRange
	// Traverse the contents table
	block := uint64(2)
	for {
		// Get the range for this block
		r := e.getPageRange(block, block)
		if r.Start.Segment != r.End.Segment {
			return nil, fmt.Errorf("Contents block split across segments")
		}
		bytes := e.m[r.Start.Segment]
		bytes = bytes[r.Start.Byte : r.End.Byte+1]
		// Get the address of the next contents block
		block = uint64FromBytes(bytes)
		bytes = bytes[8:]
		// Look for matching contents entries
		for {
			threadID := uint32FromBytes(bytes)
			if threadID == thread {
				blockStart := uint32FromBytes(bytes[4:])
				blockEnd := uint32FromBytes(bytes[8:])
				r = e.getPageRange(uint64(blockStart), uint64(blockEnd))
				ret = append(ret, r)
			}
			bytes = bytes[12:]
			if len(bytes) < 12 {
				break
			}
		}
		// Time to stop
		if block == 0 {
			break
		}
	}
	return ret, nil
}
