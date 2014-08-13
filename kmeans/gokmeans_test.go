package gokmeans

import (
	"fmt"
	"testing"
)

var observations []Node = []Node{
	Node{20.0, 20.0, 20.0, 20.0},
	Node{21.0, 21.0, 21.0, 21.0},
	Node{100.5, 100.5, 100.5, 100.5},
	Node{50.1, 50.1, 50.1, 50.1},
	Node{64.2, 64.2, 64.2, 64.2},
}

var clusters []Node = []Node{
	Node{20.0, 20.0, 20.0, 20.0},
	Node{21.0, 21.0, 21.0, 21.0},
	Node{100.5, 100.5, 100.5, 100.5},
	Node{50.1, 50.1, 50.1, 50.1},
	Node{64.2, 64.2, 64.2, 64.2},
}

func TestDistance(t *testing.T) {
	if distance(observations[3], observations[3]) != 0 {
		t.Fail()
	}
}

func TestNearest(t *testing.T) {
	if Nearest(observations[3], clusters) != 3 {
		t.Fail()
	}
}

func TestMeanNode(t *testing.T) {
	values := []Node{
		Node{10, 10, 10, 10},
		Node{20, 20, 20, 20},
	}

	for _, value := range meanNode(values) {
		if value != 15 {
			fmt.Println(value)
			t.Fail()
		}
	}
}

func TestTrain(t *testing.T) {
	if worked, clusters := Train(observations, 2, 50); !worked || clusters == nil || len(clusters) != 2 {
		t.Log("Worked:", worked, "\nClusters:", clusters)
		t.Fail()
	}
}

func BenchmarkTrain(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Train(observations, 2, 50)
	}
}
