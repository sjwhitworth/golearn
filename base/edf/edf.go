package edf

// map.go:    handles mmaping, truncation, header creation, verification,
//            creation of initial thread contents block (todo)
//            creation of initial thread metadata block (todo)
// thread.go: handles extending thread contents block (todo)
//            extending thread metadata block (todo), adding threads (todo),
//            retrieving the segments and offsets relevant to a thread (todo)
//            resolution of threads by name (todo)
//            appending data to a thread (todo)
//            deleting threads (todo)

const (
	// EDF_VERSION is the file format version
	EDF_VERSION = 1
	// EDF_LENGTH is th number of OS pages in each slice
	EDF_LENGTH = 1024 * 1024
	// EDF_SIZE sets the maximum size of the mapping, represented with
	// EDF_LENGTH segments
	// Currently set arbitrarily to 8 MiB
	EDF_SIZE = 8 * (1024 * 1024)
)

const (
	// EDF_READ_ONLY means the file will only be read, modifications fail
	EDF_READ_ONLY = iota
	// EDF_READ_WRITE specifies that the file will be read and written
	EDF_READ_WRITE
	// EDF_CREATE means the file will be created and opened with EDF_READ_WRITE
	EDF_CREATE
)

const (
	// EDF_UNMAP_NOSYNC means the file won't be
	// Sync'd to disk before unmapping
	EDF_UNMAP_NOSYNC = iota
	// EDF_UNMAP_SYNC synchronises the EDF file to disk
	// during unmapping
	EDF_UNMAP_SYNC
)
