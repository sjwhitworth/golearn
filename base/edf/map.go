package edf

import (
	"fmt"
	mmap "github.com/Sentimentron/go-mmap"
	"os"
	"runtime"
	"runtime/debug"
)

type edfMode int

const (
	edfFileMode edfMode = iota
	edfAnonMode
	edfFreedMode
)

// EdfFile represents a mapped file on disk or
// and anonymous mapping for instance storage
type EdfFile struct {
	f           *os.File
	m           []mmap.Mmap
	segmentSize uint64
	pageSize    uint64
	mode        edfMode
}

// GetPageSize returns the pageSize of an EdfFile
func (e *EdfFile) GetPageSize() uint64 {
	return e.pageSize
}

// edfPosition represents a start and finish point
// within the mapping
type edfPosition struct {
	Segment uint64
	Byte    uint64
}

// edfRange represents a start and an end segment
// mapped in an EdfFile and also the byte offsets
// within that segment
type edfRange struct {
	Start       edfPosition
	End         edfPosition
	segmentSize uint64
}

// Size returns the size (in bytes) of a given edfRange
func (r *edfRange) Size() uint64 {
	ret := uint64(r.End.Segment-r.Start.Segment) * r.segmentSize
	ret += uint64(r.End.Byte - r.Start.Byte)
	return ret
}

// edfCallFree is a half-baked finalizer called on garbage
// collection to ensure that the mapping gets freed
func edfCallFree(e *EdfFile) {
	e.unmap(EDF_UNMAP_NOSYNC)
}

// EdfAnonMap maps the EdfFile structure into RAM
// IMPORTANT: everything's lost if unmapped
func EdfAnonMap() (*EdfFile, error) {

	var err error

	// Allocate return structure
	ret := new(EdfFile)

	// Create mapping references
	ret.m = make([]mmap.Mmap, 0)

	// Get the page size
	pageSize := int64(os.Getpagesize())

	// Segment size is the size of each mapped region
	ret.pageSize = uint64(pageSize)
	ret.segmentSize = uint64(EDF_LENGTH) * uint64(os.Getpagesize())

	// Set the mode
	ret.mode = edfAnonMode

	// Allocate 4 pages initially
	ret.truncate(4)

	// Generate the header
	ret.createHeader()
	err = ret.writeInitialData()

	// Make sure this gets unmapped on garbage collection
	runtime.SetFinalizer(ret, edfCallFree)

	return ret, err
}

// edfMap takes an os.File and returns an EdfMappedFile
// structure, which represents the mmap'd underlying file
//
// The `mode` parameter takes the following values
//      EDF_CREATE: edfMap will truncate the file to the right length and write the correct header information
//      EDF_READ_WRITE: edfMap will verify header information
//      EDF_READ_ONLY:  edfMap will verify header information
// IMPORTANT: EDF_LENGTH (edf.go) controls the size of the address
// space mapping. This means that the file can be truncated to the
// correct size without remapping. On 32-bit systems, this
// is set to 2GiB.
func edfMap(f *os.File, mode int) (*EdfFile, error) {
	var err error

	// Set up various things
	ret := new(EdfFile)
	ret.f = f
	ret.m = make([]mmap.Mmap, 0)

	// Figure out the flags
	protFlags := mmap.PROT_READ
	if mode == EDF_READ_WRITE || mode == EDF_CREATE {
		protFlags |= mmap.PROT_WRITE
	}
	mapFlags := mmap.MAP_FILE | mmap.MAP_SHARED
	// Get the page size
	pageSize := int64(os.Getpagesize())
	// Segment size is the size of each mapped region
	ret.pageSize = uint64(pageSize)
	ret.segmentSize = uint64(EDF_LENGTH) * uint64(os.Getpagesize())

	// Map the file
	for i := int64(0); i < EDF_SIZE; i += int64(EDF_LENGTH) * pageSize {
		thisMapping, err := mmap.Map(f, i*pageSize, int(int64(EDF_LENGTH)*pageSize), protFlags, mapFlags)
		if err != nil {
			// TODO: cleanup
			return nil, err
		}
		ret.m = append(ret.m, thisMapping)
	}

	// Verify or generate the header
	if mode == EDF_READ_WRITE || mode == EDF_READ_ONLY {
		err = ret.verifyHeader()
		if err != nil {
			return nil, err
		}
	} else if mode == EDF_CREATE {
		err = ret.truncate(4)
		if err != nil {
			return nil, err
		}
		ret.createHeader()
		err = ret.writeInitialData()
	} else {
		err = fmt.Errorf("Unrecognised flags")
	}

	// Make sure this gets unmapped on garbage collection
	runtime.SetFinalizer(ret, edfCallFree)

	// Set the mode
	ret.mode = edfFileMode

	return ret, err

}

// getByteRange returns the segment offset and range of
// two positions in the file.
func (e *EdfFile) getByteRange(byteStart uint64, byteEnd uint64) edfRange {
	var ret edfRange
	ret.Start.Segment = byteStart / e.segmentSize
	ret.End.Segment = byteEnd / e.segmentSize
	ret.Start.Byte = byteStart % e.segmentSize
	ret.End.Byte = byteEnd % e.segmentSize
	ret.segmentSize = e.segmentSize
	return ret
}

// getPageRange returns the segment offset and range of
// two pages in the file.
func (e *EdfFile) getPageRange(pageStart uint64, pageEnd uint64) edfRange {
	return e.getByteRange(pageStart*e.pageSize, pageEnd*e.pageSize+e.pageSize-1)
}

// verifyHeader checks that this version of Golearn can
// read the file presented.
func (e *EdfFile) verifyHeader() error {
	// Check the magic bytes
	diff := (e.m[0][0] ^ byte('G')) | (e.m[0][1] ^ byte('O'))
	diff |= (e.m[0][2] ^ byte('L')) | (e.m[0][3] ^ byte('N'))
	if diff != 0 {
		return fmt.Errorf("Invalid magic bytes")
	}
	// Check the file version
	version := uint32FromBytes(e.m[0][4:8])
	if version != EDF_VERSION {
		return fmt.Errorf("Unsupported version: %d", version)
	}
	// Check the page size
	pageSize := uint32FromBytes(e.m[0][8:12])
	if pageSize != uint32(os.Getpagesize()) {
		return fmt.Errorf("Unsupported page size: (file: %d, system: %d", pageSize, os.Getpagesize())
	}
	return nil
}

// createHeader writes a valid header file into the file.
// Unexported since it can cause data loss.
func (e *EdfFile) createHeader() {
	e.m[0][0] = byte('G')
	e.m[0][1] = byte('O')
	e.m[0][2] = byte('L')
	e.m[0][3] = byte('N')
	uint32ToBytes(EDF_VERSION, e.m[0][4:8])
	uint32ToBytes(uint32(os.Getpagesize()), e.m[0][8:12])
	e.sync()
}

// writeInitialData writes system thread information
func (e *EdfFile) writeInitialData() error {
	var t thread
	t.name = "SYSTEM"
	t.id = 1
	err := e.WriteThread(&t)
	if err != nil {
		return err
	}
	t.name = "FIXED"
	t.id = 2
	err = e.WriteThread(&t)
	return err
}

// getThreadCount returns the number of threads in this file.
func (e *EdfFile) getThreadCount() uint32 {
	// The number of threads is stored in bytes 12-16 in the header
	return uint32FromBytes(e.m[0][12:])
}

// incrementThreadCount increments the record of the number
// of threads in this file
func (e *EdfFile) incrementThreadCount() uint32 {
	cur := e.getThreadCount()
	cur++
	uint32ToBytes(cur, e.m[0][12:])
	return cur
}

// GetThreads returns the thread identifier -> name map.
func (e *EdfFile) GetThreads() (map[uint32]string, error) {
	ret := make(map[uint32]string)
	count := e.getThreadCount()
	// The starting block
	block := uint64(1)
	for {
		// Decode the block offset
		r := e.getPageRange(block, block)
		if r.Start.Segment != r.End.Segment {
			return nil, fmt.Errorf("thread range split across segments")
		}
		bytes := e.m[r.Start.Segment]
		bytes = bytes[r.Start.Byte : r.End.Byte+1]
		// The first 8 bytes say where to go next
		block = uint64FromBytes(bytes)
		bytes = bytes[8:]

		for {
			length := uint32FromBytes(bytes)
			if length == 0 {
				break
			}
			t := &thread{}
			size := t.Deserialize(bytes)
			bytes = bytes[size:]
			ret[t.id] = t.name[0:len(t.name)]
		}
		// If next block offset is zero, no more threads to read
		if block == 0 {
			break
		}
	}
	// Hey? What's wrong with you!
	if len(ret) != int(count) {
		return ret, fmt.Errorf("thread mismatch: %d/%d, indicates possible corruption", len(ret), count)
	}
	return ret, nil
}

// sync writes information to physical storage.
func (e *EdfFile) sync() error {
	// Do nothing if we're mapped anonymously
	if e.mode == edfAnonMode {
		return nil
	}
	for _, m := range e.m {
		err := m.Sync(mmap.MS_SYNC)
		if err != nil {
			return err
		}
	}
	return nil
}

// truncate changes the size of the underlying file
// The size of the address space doesn't change.
func (e *EdfFile) truncateFile(size int64) error {
	pageSize := int64(os.Getpagesize())
	newSize := pageSize * size

	// Synchronise
	// e.sync()

	// Double-check that we're not reducing file size
	fileInfo, err := e.f.Stat()
	if err != nil {
		return err
	}
	if fileInfo.Size() > newSize {
		return fmt.Errorf("Can't reduce file size!")
	}

	// Truncate the file
	err = e.f.Truncate(newSize)
	if err != nil {
		return err
	}

	// Verify that the file is larger now than it was
	fileInfo, err = e.f.Stat()
	if err != nil {
		return err
	}
	if fileInfo.Size() != newSize {
		return fmt.Errorf("Truncation failed: %d, %d", fileInfo.Size(), newSize)
	}
	return err
}

func (e *EdfFile) truncateMem(size int64) error {
	pageSize := int64(os.Getpagesize())
	newSize := pageSize * size

	currentSize := 0
	for _, m := range e.m {
		currentSize += len(m)
	}

	if int64(currentSize) > newSize {
		return fmt.Errorf("Can't reduce size")
	}

	// Allocate some more memory
	for i := uint64(currentSize); i < uint64(newSize); i += e.segmentSize {
		newMap := make([]byte, e.segmentSize)
		e.m = append(e.m, newMap)
	}

	return nil
}

func (e *EdfFile) truncate(size int64) error {
	if e.mode == edfAnonMode {
		return e.truncateMem(size)
	} else if e.mode == edfFileMode {
		return e.truncateFile(size)
	}
	panic("Unsupported")
}

// unmap unlinks the EdfFile from the address space.
// EDF_UNMAP_NOSYNC skips calling Sync() on the underlying
// file before this happens.
// IMPORTANT: attempts to use this mapping after unmap() is
// called will result in crashes.
func (e *EdfFile) unmap(flags int) error {

	var ret error

	// Check if the file has already been freed
	if e.mode == edfFreedMode {
		fmt.Fprintln(os.Stderr, "Potential double-free")
		debug.PrintStack()
		return nil
	} else if e.mode == edfAnonMode {
		// If it's anonymous, don't do anything
		e.m = nil
		e.mode = edfFreedMode
		return nil
	}

	// Sync the file
	if flags != EDF_UNMAP_NOSYNC {
		e.sync()
	}

	e.mode = edfFreedMode
	// Unmap the file
	for _, m := range e.m {
		err := m.Unmap()
		if err != nil {
			ret = err
		}
	}
	return ret
}

// ResolveRange returns a slice of byte slices representing
// the underlying memory referenced by edfRange.
//
// WARNING: slow.
func (e *EdfFile) ResolveRange(r edfRange) [][]byte {
	var ret [][]byte
	segCounter := 0
	for segment := r.Start.Segment; segment <= r.End.Segment; segment++ {
		if segment == r.Start.Segment {
			ret = append(ret, e.m[segment][r.Start.Byte:])
		} else if segment == r.End.Segment {
			ret = append(ret, e.m[segment][:r.End.Byte+1])
		} else {
			ret = append(ret, e.m[segment])
		}
		segCounter++
	}
	return ret
}
