package trees

import (
	"sort"
)

type Slice struct {
	sort.Float64Slice
	Idx []int
}

func (s Slice) Swap(i, j int) {
	s.Float64Slice.Swap(i, j)
	s.Idx[i], s.Idx[j] = s.Idx[j], s.Idx[i]
}

func NewSlice(n []float64) *Slice {
	s := &Slice{Float64Slice: sort.Float64Slice(n), Idx: make([]int, len(n))}

	for i := range s.Idx {
		s.Idx[i] = i
	}
	return s
}
