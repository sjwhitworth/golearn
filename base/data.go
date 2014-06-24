package base

// SortDirection specifies sorting direction...
type SortDirection int

const (
	// Descending says that Instances should be sorted high to low...
	Descending SortDirection = 1
	// Ascending states that Instances should be sorted low to high...
	Ascending SortDirection = 2
)

type DataGrid interface {
	GetAttribute(Attribute) (AttributeSpec, error)
	AllAttributes() []Attribute
	AddClassAttribute(Attribute) error
	RemoveClassAttribute(Attribute) error
	AllClassAttributes() []Attribute
}

type FixedDataGrid interface {
	DataGrid
	Get(AttributeSpec, int) []byte
	RowString(int) string
	Size() (int, int)
	MapOverRows([]AttributeSpec, func([][]byte, int) (bool, error)) error
}

type UpdatableDataGrid interface {
	FixedDataGrid
	Set(AttributeSpec, int, []byte)
	AddAttribute(Attribute) AttributeSpec
}
