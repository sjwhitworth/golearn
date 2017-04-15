package kdtree

import (
	"errors"
)

type heapNode struct {
	value  []float64
	length float64
}

type heap struct {
	tree []heapNode
}

// newHeap return a pointer of heap.
func newHeap() *heap {
	h := &heap{}
	h.tree = make([]heapNode, 0)
	return &heap{}
}

// maximum return the max heapNode in the heap.
func (h *heap) maximum() (heapNode, error) {
	if len(h.tree) == 0 {
		return heapNode{}, h.errEmpty()
	}

	return h.tree[0], nil
}

// extractMax remove the Max heapNode in the heap.
func (h *heap) extractMax() {
	if len(h.tree) == 0 {
		return
	}

	h.tree[0] = h.tree[len(h.tree)-1]
	h.tree = h.tree[:len(h.tree)-1]

	target := 1
	for true {
		largest := target
		if target*2-1 >= len(h.tree) {
			break
		}
		if h.tree[target*2-1].length > h.tree[target].length {
			largest = target * 2
		}

		if target*2 >= len(h.tree) {
			break
		}
		if h.tree[target*2].length > h.tree[largest-1].length {
			largest = target*2 + 1
		}

		if largest == target {
			break
		}
		h.tree[largest-1], h.tree[target-1] = h.tree[target-1], h.tree[largest-1]
		target = largest
	}
}

// insert put a new heapNode into heap.
func (h *heap) insert(value []float64, length float64) {
	node := heapNode{}
	node.length = length
	node.value = make([]float64, len(value))
	copy(node.value, value)
	h.tree = append(h.tree, node)

	target := len(h.tree)
	for target != 1 {
		if h.tree[(target/2)-1].length >= h.tree[target-1].length {
			break
		}
		h.tree[target-1], h.tree[(target/2)-1] = h.tree[(target/2)-1], h.tree[target-1]
		target /= 2
	}
}

// errEmpty is return an error which is returned
// when heap is empty.
func (h *heap) errEmpty() error {
	return errors.New("empty heap")
}
