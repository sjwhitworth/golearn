package cross_validation

import (
	. "github.com/smartystreets/goconvey/convey"
	mat64 "github.com/gonum/matrix/mat64"
	"testing"
	"sync"
)

func TestShuffleMatrix(t *testing.T) {
	var vectorX, vectorY *mat64.Dense

	Convey("Given two equal vectors", t, func() {
			vectorX = mat64.NewDense(3, 1, []float64{1, 2, 3})
			vectorY = mat64.DenseCopyOf(vectorX)

			Convey("After shuffling", func() {
					wg := new(sync.WaitGroup)
					wg.Add(1)
					shuffleMatrix(vectorY, 0, wg)
					wg.Wait()
					result := vectorX.Equals(vectorY)

					Convey("The vectors should be different", func() {
							So(result, ShouldNotEqual, true)
						})
				})

		})
}
