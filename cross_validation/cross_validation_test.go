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

func TestTrainTestSplit(t *testing.T) {
	data := make([]float64, 100)
	labels := make([]float64, 20)

	for i := 0; i < 100; i++ {
		data[i] = float64(i)
		labels[i/5] = float64(i-4)
	}

	Convey("Given a data and a labels matrix", t, func() {
			dataMatrix := mat64.NewDense(20, 5, data)
			labelMatrix := mat64.NewDense(20, 1, labels)

			Convey("After splitting into 15 training and 5 testing instances without shuffling", func() {
					newData, _, _ := TrainTestSplit(5, nil, dataMatrix, labelMatrix)

					Convey("First 15 items of the original data matrix should equal new data matrix", func() {
							So(dataMatrix.RawMatrix().Data[:75], ShouldResemble, newData[0].RawMatrix().Data)
						})

				})
			Convey("After splitting into 15 training and 5 testing instances with shuffling", func() {
					newData, newLabels, _ := TrainTestSplit(5, 99, dataMatrix, labelMatrix)

					Convey("First 15 items of the original data matrix should not equal new data matrix", func() {
							So(dataMatrix.RawMatrix().Data[:75], ShouldNotResemble, newData[0].RawMatrix().Data)
						})

					Convey("First element of every row aligns with labels", func() {
							for i, v := range newLabels[0].RawMatrix().Data {
								So(newData[0].At(i, 0), ShouldEqual, v)
							}
						})

					Convey("After shuffling the same matrix with the same seed", func() {
							newestData, _, _ := TrainTestSplit(5, 99, dataMatrix, labelMatrix)

							Convey("First 15 items of the new data matrix should equal newest data matrix", func() {
									So(newData[0].RawMatrix().Data[:75], ShouldResemble, newestData[0].RawMatrix().Data)
								})
						})
				})
		})

}
