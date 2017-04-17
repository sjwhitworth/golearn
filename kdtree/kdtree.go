package kdtree

import (
	"errors"
	"github.com/gonum/matrix/mat64"
	"github.com/sjwhitworth/golearn/metrics/pairwise"
	"sort"
)

type node struct {
	feature  int
	value    []float64
	srcRowNo int
	left     *node
	right    *node
}

// Tree is a kdtree.
type Tree struct {
	firstDiv *node
	data     [][]float64
}

type SortData struct {
	RowData [][]float64
	Data    []int
	Feature int
}

func (d SortData) Len() int { return len(d.Data) }
func (d SortData) Less(i, j int) bool {
	return d.RowData[d.Data[i]][d.Feature] < d.RowData[d.Data[j]][d.Feature]
}
func (d SortData) Swap(i, j int) { d.Data[i], d.Data[j] = d.Data[j], d.Data[i] }

// New return a Tree pointer.
func New() *Tree {
	return &Tree{}
}

// Build builds the kdtree with specific data.
func (t *Tree) Build(data [][]float64) error {
	if len(data) == 0 {
		return errors.New("no input data")
	}
	size := len(data[0])
	for _, v := range data {
		if len(v) != size {
			return errors.New("amounts of features are not the same")
		}
	}

	t.data = data

	newData := make([]int, len(data))
	for k, _ := range newData {
		newData[k] = k
	}

	if len(data) == 1 {
		t.firstDiv = &node{feature: -1, srcRowNo: 0}
		t.firstDiv.value = make([]float64, len(data[0]))
		copy(t.firstDiv.value, data[0])
	} else {
		t.firstDiv = t.buildHandle(newData, 0)
	}

	return nil
}

// buildHandle builds the kdtree recursively.
func (t *Tree) buildHandle(data []int, featureIndex int) *node {
	n := &node{feature: featureIndex}

	tmp := SortData{RowData: t.data, Data: data, Feature: featureIndex}
	sort.Sort(tmp)
	middle := len(data) / 2

	n.srcRowNo = data[middle]
	n.value = make([]float64, len(t.data[data[middle]]))
	copy(n.value, t.data[data[middle]])

	divPoint := middle
	for i := middle + 1; i < len(data); i++ {
		if t.data[data[i]][featureIndex] == t.data[data[middle]][featureIndex] {
			divPoint = i
		} else {
			break
		}
	}

	if divPoint == 1 {
		n.left = &node{feature: -1}
		n.left.value = make([]float64, len(t.data[data[0]]))
		copy(n.left.value, t.data[data[0]])
		n.left.srcRowNo = data[0]
	} else {
		n.left = t.buildHandle(data[:divPoint], (featureIndex+1)%len(t.data[data[0]]))
	}

	if divPoint == (len(data) - 2) {
		n.right = &node{feature: -1}
		n.right.value = make([]float64, len(t.data[data[divPoint+1]]))
		copy(n.right.value, t.data[data[divPoint+1]])
		n.left.srcRowNo = data[divPoint+1]
	} else if divPoint != (len(data) - 1) {
		n.right = t.buildHandle(data[divPoint+1:], (featureIndex+1)%len(t.data[data[0]]))
	}

	return n
}

// Search return []int contained k nearest neighbor from
// specific distance function.
func (t *Tree) Search(k int, disType pairwise.PairwiseDistanceFunc, target []float64) ([]int, error) {
	if k > len(t.data) {
		return []int{}, errors.New("k is largerer than amount of trainData")
	}

	if len(target) != len(t.data[0]) {
		return []int{}, errors.New("amount of features is not equal")
	}

	h := newHeap()
	t.searchHandle(k, disType, target, h, t.firstDiv)

	out := make([]int, k)
	i := k - 1
	for h.size() != 0 {
		out[i] = h.maximum().srcRowNo
		i--
		h.extractMax()
	}

	return out, nil
}

func (t *Tree) searchHandle(k int, disType pairwise.PairwiseDistanceFunc, target []float64, h *heap, n *node) {
	if n.feature == -1 {
		vectorX := mat64.NewDense(len(target), 1, target)
		vectorY := mat64.NewDense(len(target), 1, n.value)
		length := disType.Distance(vectorX, vectorY)
		h.insert(n.value, length, n.srcRowNo)
		return
	}

	dir := true
	if target[n.feature] <= n.value[n.feature] {
		t.searchHandle(k, disType, target, h, n.left)
	} else {
		dir = false
		t.searchHandle(k, disType, target, h, n.right)
	}

	vectorX := mat64.NewDense(len(target), 1, target)
	vectorY := mat64.NewDense(len(target), 1, n.value)
	length := disType.Distance(vectorX, vectorY)

	if k > h.size() {
		h.insert(n.value, length, n.srcRowNo)
		if dir {
			t.searchAllNodes(k, disType, target, h, n.right)
		} else {
			t.searchAllNodes(k, disType, target, h, n.left)
		}
	} else if h.maximum().length > length {
		h.extractMax()
		h.insert(n.value, length, n.srcRowNo)
		if dir {
			t.searchAllNodes(k, disType, target, h, n.right)
		} else {
			t.searchAllNodes(k, disType, target, h, n.left)
		}
	}
}

func (t *Tree) searchAllNodes(k int, disType pairwise.PairwiseDistanceFunc, target []float64, h *heap, n *node) {
	vectorX := mat64.NewDense(len(target), 1, target)
	vectorY := mat64.NewDense(len(target), 1, n.value)
	length := disType.Distance(vectorX, vectorY)

	if k > h.size() {
		h.insert(n.value, length, n.srcRowNo)
	} else if h.maximum().length > length {
		h.extractMax()
		h.insert(n.value, length, n.srcRowNo)
	}

	if n.left != nil {
		t.searchAllNodes(k, disType, target, h, n.left)
	}
	if n.right != nil {
		t.searchAllNodes(k, disType, target, h, n.right)
	}
}
