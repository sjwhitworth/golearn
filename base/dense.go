package base

import (
	"bytes"
	"fmt"
	"github.com/sjwhitworth/golearn/base/edf"
	"math"
	"sync"
)

// DenseInstances stores each Attribute value explicitly
// in a large grid.
type DenseInstances struct {
	storage    *edf.EdfFile
	pondMap    map[string]int
	ponds      []*Pond
	lock       sync.Mutex
	fixed      bool
	classAttrs map[AttributeSpec]bool
	maxRow     int
	attributes []Attribute
}

// NewDenseInstances generates a new DenseInstances set
// with an anonymous EDF mapping and default settings.
func NewDenseInstances() *DenseInstances {
	storage, err := edf.EdfAnonMap()
	if err != nil {
		panic(err)
	}
	return &DenseInstances{
		storage,
		make(map[string]int),
		make([]*Pond, 0),
		sync.Mutex{},
		false,
		make(map[AttributeSpec]bool),
		0,
		make([]Attribute, 0),
	}
}

//
// Pond functions
//

// createPond adds a new Pond to this set of Instances
// IMPORTANT: do not call unless you've acquired the lock
func (inst *DenseInstances) createPond(name string, size int) {
	if inst.fixed {
		panic("Can't add additional Attributes")
	}

	// Resolve or create thread
	threads, err := inst.storage.GetThreads()
	if err != nil {
		panic(err)
	}

	ok := false
	for i := range threads {
		if threads[i] == name {
			ok = true
			break
		}
	}
	if ok {
		panic("Can't create pond: pond thread already exists")
	}

	// Write the pool's thread into the file
	thread := edf.NewThread(inst.storage, name)
	err = inst.storage.WriteThread(thread)
	if err != nil {
		panic(fmt.Sprintf("Can't write thread: %s", err))
	}

	// Create the pond information
	pond := new(Pond)
	pond.threadNo = thread.GetId()
	pond.parent = inst
	pond.attributes = make([]Attribute, 0)
	pond.size = size
	pond.alloc = make([][]byte, 0)
	// Store within instances
	inst.pondMap[name] = len(inst.ponds)
	inst.ponds = append(inst.ponds, pond)
}

// CreatePond adds a new Pond to this set of instances
// with a given name. If the size is 0, a bit-pond is added
// if the size of not 0, then the size of each pond attribute
// is set to that number of bytes.
func (inst *DenseInstances) CreatePond(name string, size int) (err error) {
	defer func() {
		if r := recover(); r != nil {
			var ok bool
			if err, ok = r.(error); !ok {
				err = fmt.Errorf("CreatePond: %v (not created)", r)
			}
		}
	}()

	inst.lock.Lock()
	defer inst.lock.Unlock()

	inst.createPond(name, size)
	return nil
}

// GetPond returns a reference to a Pond of a given name /
func (inst *DenseInstances) GetPond(name string) (*Pond, error) {
	inst.lock.Lock()
	defer inst.lock.Unlock()

	// Check if the pond exists
	if id, ok := inst.pondMap[name]; !ok {
		return nil, fmt.Errorf("Pond '%s' doesn't exist", name)
	} else {
		// Return the pond
		return inst.ponds[id], nil
	}
}

//
// Attribute creation and handling
//

// AddAttribute adds an Attribute to this set of DenseInstances
// Creates a default Pond for it if a suitable one doesn't exist.
// Returns an AttributeSpec for subsequent Set() calls.
//
// IMPORTANT: will panic if storage has been allocated via Extend.
func (inst *DenseInstances) AddAttribute(a Attribute) AttributeSpec {
	inst.lock.Lock()
	defer inst.lock.Unlock()

	if inst.fixed {
		panic("Can't add additional Attributes")
	}

	// Generate a default Pond name
	pond := "FLOAT"
	if _, ok := a.(*CategoricalAttribute); ok {
		pond = "CAT"
	} else if _, ok := a.(*FloatAttribute); ok {
		pond = "FLOAT"
	} else {
		panic("Unrecognised Attribute type")
	}

	// Create the pond if it doesn't exist
	if _, ok := inst.pondMap[pond]; !ok {
		inst.createPond(pond, 8)
	}
	id := inst.pondMap[pond]
	p := inst.ponds[id]
	p.attributes = append(p.attributes, a)
	inst.attributes = append(inst.attributes, a)
	return AttributeSpec{id, len(p.attributes) - 1, a}
}

// AddAttributeToPond adds an Attribute to a given pond
func (inst *DenseInstances) AddAttributeToPond(newAttribute Attribute, pond string) (AttributeSpec, error) {
	inst.lock.Lock()
	defer inst.lock.Unlock()

	// Check if the pond exists
	if _, ok := inst.pondMap[pond]; !ok {
		return AttributeSpec{-1, 0, nil}, fmt.Errorf("Pond '%s' doesn't exist. Call CreatePond() first", pond)
	}

	id := inst.pondMap[pond]
	p := inst.ponds[id]
	for i, a := range p.attributes {
		if !a.Compatable(newAttribute) {
			return AttributeSpec{-1, 0, nil}, fmt.Errorf("Attribute %s is not compatable with %s in pond '%s' (position %d)", newAttribute, a, pond, i)
		}
	}

	p.attributes = append(p.attributes, newAttribute)
	inst.attributes = append(inst.attributes, newAttribute)
	return AttributeSpec{id, len(p.attributes) - 1, newAttribute}, nil
}

// GetAttribute returns an Attribute equal to the argument.
//
// TODO: Write a function to pre-compute this once we've allocated
// TODO: Write a utility function which retrieves all AttributeSpecs for
// a given instance set.
func (inst *DenseInstances) GetAttribute(get Attribute) (AttributeSpec, error) {
	inst.lock.Lock()
	defer inst.lock.Unlock()

	for i, p := range inst.ponds {
		for j, a := range p.attributes {
			if a.Equals(get) {
				return AttributeSpec{i, j, a}, nil
			}
		}
	}

	return AttributeSpec{-1, 0, nil}, fmt.Errorf("Couldn't resolve %s", get)
}

// AllAttributes returns a slice of all Attributes.
func (inst *DenseInstances) AllAttributes() []Attribute {
	inst.lock.Lock()
	defer inst.lock.Unlock()

	return inst.attributes
}

// AddClassAttribute sets an Attribute to be a class Attribute.
func (inst *DenseInstances) AddClassAttribute(a Attribute) error {

	as, err := inst.GetAttribute(a)
	if err != nil {
		return err
	}

	inst.lock.Lock()
	defer inst.lock.Unlock()

	inst.classAttrs[as] = true
	return nil
}

// RemoveClassAttribute removes an Attribute from the set of class Attributes.
func (inst *DenseInstances) RemoveClassAttribute(a Attribute) error {
	inst.lock.Lock()
	defer inst.lock.Unlock()

	as, err := inst.GetAttribute(a)
	if err != nil {
		return err
	}

	inst.lock.Lock()
	defer inst.lock.Unlock()

	inst.classAttrs[as] = false
	return nil
}

// AllClassAttributes returns a slice of Attributes which have
// been designated class Attributes.
func (inst *DenseInstances) AllClassAttributes() []Attribute {
	var ret []Attribute
	inst.lock.Lock()
	defer inst.lock.Unlock()
	for a := range inst.classAttrs {
		if inst.classAttrs[a] {
			ret = append(ret, a.attr)
		}
	}
	return ret
}

//
// Allocation functions
//

// Extend extends this set of Instances to store rows additional rows.
// It's recommended to set rows to something quite large.
//
// IMPORTANT: panics if the allocation fails
func (inst *DenseInstances) Extend(rows int) error {

	inst.lock.Lock()
	defer inst.lock.Unlock()

	// Get the size of each page
	pageSize := inst.storage.GetPageSize()

	for pondName := range inst.ponds {
		p := inst.ponds[pondName]

		// Compute pond row storage requirements
		rowSize := p.RowSize()

		// How many rows can we store per page?
		rowsPerPage := float64(pageSize) / float64(rowSize)

		// How many pages?
		pagesNeeded := uint32(math.Ceil(float64(rows) / rowsPerPage))

		// Allocate those pages
		r, err := inst.storage.AllocPages(pagesNeeded, p.threadNo)
		if err != nil {
			panic(fmt.Sprintf("Allocation error: %s (rowSize %d, pageSize %d, rowsPerPage %.2f, tried to allocate %d page(s) and extend by %d row(s))", err, rowSize, pageSize, rowsPerPage, pagesNeeded, rows))
		}
		// Resolve and assign those pages
		byteBlock := inst.storage.ResolveRange(r)
		for _, block := range byteBlock {
			p.alloc = append(p.alloc, block)
		}
	}
	inst.fixed = true
	inst.maxRow += rows
	return nil
}

// Set sets a particular Attribute (given as an AttributeSpec) on a particular
// row to a particular value.
//
// AttributeSpecs can be obtained using GetAttribute() or AddAttribute().
//
// IMPORTANT: Will panic() if the AttributeSpec isn't valid
//
// IMPORTANT: Will panic() if the row is too large
//
// IMPORTANT: Will panic() if the val is not the right length
func (inst *DenseInstances) Set(a AttributeSpec, row int, val []byte) {
	inst.ponds[a.pond].set(a.position, row, val)
}

// Get gets a particular Attribute (given as an AttributeSpec) on a particular
// row.
// AttributeSpecs can be obtained using GetAttribute() or AddAttribute()
func (inst *DenseInstances) Get(a AttributeSpec, row int) []byte {
	return inst.ponds[a.pond].get(a.position, row)
}

// RowString returns a string representation of a given row.
func (inst *DenseInstances) RowString(row int) string {
	var buffer bytes.Buffer
	first := true
	for name := range inst.ponds {
		if first {
			first = false
		} else {
			buffer.WriteString(" ")
		}
		p := inst.ponds[name]
		p.appendToRowBuf(row, &buffer)
	}
	return buffer.String()
}

//
// Row handling functions
//

func (inst *DenseInstances) allocateRowVector(asv []AttributeSpec) [][]byte {
	ret := make([][]byte, len(asv))
	for i, as := range asv {
		p := inst.ponds[as.pond]
		ret[i] = make([]byte, p.size)
	}
	return ret
}

// MapOverRows passes each row map into a function.
// First argument is a list of AttributeSpec in the order
// they're needed in for the function. The second is the function
// to call on each row.
func (inst *DenseInstances) MapOverRows(asv []AttributeSpec, mapFunc func([][]byte, int) (bool, error)) error {
	rowBuf := make([][]byte, len(asv))
	for i := 0; i < inst.maxRow; i++ {
		for j, as := range asv {
			p := inst.ponds[as.pond]
			rowBuf[j] = p.get(as.position, i)
		}
		ok, err := mapFunc(rowBuf, i)
		if err != nil {
			return err
		}
		if !ok {
			break
		}
	}
	return nil
}

// Size returns the number of Attributes as the first return value
// and the maximum allocated row as the second value.
func (inst *DenseInstances) Size() (int, int) {
	return len(inst.AllAttributes()), inst.maxRow
}

// swapRows swaps over rows i and j
func (inst *DenseInstances) swapRows(i, j int) {
	as := ResolveAllAttributes(inst)
	for _, a := range as {
		v1 := inst.Get(a, i)
		v2 := inst.Get(a, j)
		v3 := make([]byte, len(v2))
		copy(v3, v2)
		inst.Set(a, j, v1)
		inst.Set(a, i, v3)
	}
}

// Equal checks whether a given Instance set is exactly the same
// as another: same size and same values (as determined by the Attributes)
//
// IMPORTANT: does not explicitly check if the Attributes are considered equal.
func (inst *DenseInstances) Equal(other DataGrid) bool {

	_, rows := inst.Size()

	for _, a := range inst.AllAttributes() {
		as1, err := inst.GetAttribute(a)
		if err != nil {
			panic(err) // That indicates some kind of error
		}
		as2, err := inst.GetAttribute(a)
		if err != nil {
			return false // Obviously has different Attributes
		}
		for i := 0; i < rows; i++ {
			b1 := inst.Get(as1, i)
			b2 := inst.Get(as2, i)
			if !byteSeqEqual(b1, b2) {
				return false
			}
		}
	}

	return true
}

// String returns a human-readable summary of this dataset.
func (inst *DenseInstances) String() string {
	var buffer bytes.Buffer

	// Get all Attribute information
	as := ResolveAllAttributes(inst)

	// Print header
	cols, rows := inst.Size()
	buffer.WriteString("Instances with ")
	buffer.WriteString(fmt.Sprintf("%d row(s) ", rows))
	buffer.WriteString(fmt.Sprintf("%d attribute(s)\n", cols))
	buffer.WriteString(fmt.Sprintf("Attributes: \n"))

	for _, a := range as {
		prefix := "\t"
		if inst.classAttrs[a] {
			prefix = "*\t"
		}
		buffer.WriteString(fmt.Sprintf("%s%s\n", prefix, a.attr))
	}

	buffer.WriteString("\nData:\n")
	maxRows := 30
	if rows < maxRows {
		maxRows = rows
	}

	for i := 0; i < maxRows; i++ {
		buffer.WriteString("\t")
		for _, a := range as {
			val := inst.Get(a, i)
			buffer.WriteString(fmt.Sprintf("%s ", a.attr.GetStringFromSysVal(val)))
		}
		buffer.WriteString("\n")
	}

	missingRows := rows - maxRows
	if missingRows != 0 {
		buffer.WriteString(fmt.Sprintf("\t...\n%d row(s) undisplayed", missingRows))
	} else {
		buffer.WriteString("All rows displayed")
	}

	return buffer.String()
}
