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
	// File format version
	EDF_VERSION = 1
	// Number of OS-pages in each slice
	EDF_LENGTH = 1024 * 1024
	// Sets the maximum size of the mapping, represented with
	// EDF_LENGTH segments
	// Currently set arbitrarily to 4 GiB
	EDF_SIZE = 4 * (1024 * 1024 * 1024)
)

const (
	// File will only be read, modifications fail
	EDF_READ_ONLY = iota
	// File will be read and written
	EDF_READ_WRITE
	// File will be created and opened with RW
	EDF_CREATE
)

const (
	// Don't call Sync before unmapping
	EDF_UNMAP_NOSYNC = iota
	EDF_UNMAP_SYNC
)
