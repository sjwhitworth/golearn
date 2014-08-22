package cross_validation

import (
	//. "github.com/smartystreets/goconvey/convey"
	"github.com/gonum/matrix/mat64"
	"math/rand"
	"testing"
	"time"
)

var (
	flatValues, flatLabels []float64
	values, labels         *mat64.Dense
)

func init() {
	flatValues = make([]float64, 80)
	flatLabels = make([]float64, 20)

	for i := 0; i < 80; i++ {
		flatValues[i] = float64(i + 1)
		// Replaces labels four times per run but who cares.
		flatLabels[int(i/4)] = float64(rand.Intn(2))
	}

	values = mat64.NewDense(20, 4, flatValues)
	labels = mat64.NewDense(20, 1, flatLabels)
}

func TestTrainTrainTestSplit(t *testing.T) {
	nolab1, err := TrainTestSplit(4, nil, values)
	if err != nil {
		t.Error(err)
	}

	// Make sure the random generator gets a new seed (time).
	time.Sleep(time.Second)

	nolab2, _ := TrainTestSplit(4, nil, values)
	if nolab1[0].Equals(nolab2[0]) {
		t.Errorf("Shuffle with different seed returned same matrix")
	}

	nolab1, _ = TrainTestSplit(4, 1, values)
	nolab2, _ = TrainTestSplit(4, 1, values)
	// Comparing the determinants does not guarantee uniqueness, but it will do for now.
	if !nolab1[0].Equals(nolab2[0]) {
		t.Errorf("Shuffle with same seed returned different matrix")
	}

	// Same thing for data with labels.
	lab1, err := TrainTestSplit(0.1, 10, values, labels)
	if err != nil {
		t.Error(err)
	}

	lab2, _ := TrainTestSplit(0.1, 10, values, labels)
	if !lab1[0].Equals(lab2[0]) {
		t.Errorf("Shuffle with same seed returned different determinants")
	}
}
