package naive

import (
    "testing"
    "github.com/gonum/matrix/mat64"
)

// Test if panic is correctly called when matrices with different
// dimensions are used as arguments.
func TestFitPanic(t *testing.T) {
    defer func() {
        if recover() == nil {
            t.Fatalf("invalid matrix dim did not panic")
        }
    }()

    nb := NewBernoulliNBClassifier(2)

    X := mat64.NewDense(10, 20, nil)
    // simulating user mistake, one extra label
    y := make([]int, 11)

    nb.Fit(X, y)
}
