package base

import (
	"bytes"
	"fmt"
	"github.com/sjwhitworth/golearn/base/edf"
	"math"
	"sync"
)

type DenseInstances struct {
	storage    *edf.EdfFile
	ponds      map[string]*Pond
	lock       sync.Mutex
	fixed      bool
	classAttrs map[AttributeSpec]bool
	maxRow     int
}

type AttributeSpec struct {
	pondName string
	position int
	attr     Attribute
}

func NewDenseInstances() *DenseInstances {
	storage, err := edf.EdfAnonMap()
	if err != nil {
		panic(err)
	}
	return &DenseInstances{
		storage,
		make(map[string]*Pond),
		sync.Mutex{},
		false,
		make(map[AttributeSpec]bool),
		0,
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
	inst.ponds[name] = pond
}

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

func (inst *DenseInstances) GetPond(name string) (*Pond, error) {

	inst.lock.Lock()
	defer inst.lock.Unlock()

	// Check if the pond exists
	if _, ok := inst.ponds[name]; !ok {
		return nil, fmt.Errorf("Pond '%s' doesn't exist", name)
	}

	// Return the pond
	return inst.ponds[name], nil
}

//
// Attribute creation and handling
//

// Adds an Attribute to this set of DenseInstances
// Creates a default Pond for it if a suitable one doesn't exist
// Returns an AttributeSpec for subsequent Set() calls
// IMPORTANT: will panic if storage has been allocated via SetSize
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
	if _, ok := inst.ponds[pond]; !ok {
		inst.createPond(pond, 8)
	}
	p := inst.ponds[pond]

	p.attributes = append(p.attributes, a)
	return AttributeSpec{pond, len(p.attributes) - 1, a}
}

func (inst *DenseInstances) AddAttributeToPond(newAttribute Attribute, pond string) (AttributeSpec, error) {
	inst.lock.Lock()
	defer inst.lock.Unlock()

	// Check if the pond exists
	if _, ok := inst.ponds[pond]; !ok {
		return AttributeSpec{"", 0, nil}, fmt.Errorf("Pond '%s' doesn't exist. Call CreatePond() first", pond)
	}

	p := inst.ponds[pond]
	for i, a := range p.attributes {
		if !a.Compatable(newAttribute) {
			return AttributeSpec{"", 0, nil}, fmt.Errorf("Attribute %s is not compatable with %s in pond '%s' (position %d)", newAttribute, a, i)
		}
	}

	p.attributes = append(p.attributes, newAttribute)
	return AttributeSpec{pond, len(p.attributes) - 1, newAttribute}, nil
}

func (inst *DenseInstances) GetAttribute(get Attribute) (AttributeSpec, error) {
	inst.lock.Lock()
	defer inst.lock.Unlock()

	for pondName := range inst.ponds {
		p := inst.ponds[pondName]
		for i, a := range p.attributes {
			if a.Equals(get) {
				return AttributeSpec{pondName, i, a}, nil
			}
		}
	}

	return AttributeSpec{"", 0, nil}, fmt.Errorf("Couldn't resolve %s", get)
}

func (inst *DenseInstances) AllAttributes() []Attribute {
	inst.lock.Lock()
	defer inst.lock.Unlock()

	ret := make([]Attribute, 0)
	for pondName := range inst.ponds {
		p := inst.ponds[pondName]
		for _, a := range p.attributes {
			ret = append(ret, a)
		}
	}

	return ret
}

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

func (inst *DenseInstances) AllClassAttributes() []Attribute {
	ret := make([]Attribute, 0)
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

// Extends this set of Instances to store rows additional rows.
// It's recommended to set rows to something quite large
// IMPORTANT: panics if the allocation fails
func (inst *DenseInstances) Extend(rows int) {

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
			panic(fmt.Sprintf("Allocation error: %s (rowSize %d, pageSize %d, rowsPerPage %.2f, tried to allocate %d page(s))", err, rowSize, pageSize, rowsPerPage, pagesNeeded))
		}
		// Resolve and assign those pages
		byteBlock := inst.storage.ResolveRange(r)
		for _, block := range byteBlock {
			p.alloc = append(p.alloc, block)
		}
	}
	inst.fixed = true
	inst.maxRow += rows
}

// Sets a particular Attribute (given as an AttributeSpec) on a particular
// row to a particular value.
// AttributeSpecs can be obtained using GetAttribute() or AddAttribute()
// IMPORTANT: Will panic() if the AttributeSpec isn't valid
// IMPORTANT: Will panic() if the row is too large
// IMPORTANT: Will panic() if the val is not the right length
func (inst *DenseInstances) Set(a AttributeSpec, row int, val []byte) {
	inst.ponds[a.pondName].set(a.position, row, val)
}

func (inst *DenseInstances) Get(a AttributeSpec, row int) []byte {
	return inst.ponds[a.pondName].get(a.position, row)
}

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
		p := inst.ponds[as.pondName]
		ret[i] = make([]byte, p.size)
	}
	return ret
}

// MapOverRows passes each row map into a function
// First argument is a list of AttributeSpec in the order
// they're needed in for the function
func (inst *DenseInstances) MapOverRows(asv []AttributeSpec, mapFunc func([][]byte, int) (bool, error)) error {
	rowBuf := make([][]byte, len(asv))
	for i := 0; i < inst.maxRow; i++ {
		for j, as := range asv {
			p := inst.ponds[as.pondName]
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

func (inst *DenseInstances) Size() (int, int) {

	return len(inst.AllAttributes()), inst.maxRow
}
