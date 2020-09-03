package neural

import (
	. "github.com/smartystreets/goconvey/convey"
	"gonum.org/v1/gonum/mat"
	"testing"
)

func TestNetworkWith1Layer(t *testing.T) {

	Convey("Given the Network from Han and Kamber (p 334)...", t, func() {
		// Contains 6 nodes, 3 input nodes and uses Sigmoid
		n := NewNetwork(6, 3, Sigmoid)

		// Set the weights
		n.SetWeight(1, 4, 0.2)
		n.SetWeight(1, 5, -0.3)
		n.SetWeight(2, 4, 0.4)
		n.SetWeight(2, 5, 0.1)
		n.SetWeight(3, 4, -0.5)
		n.SetWeight(3, 5, 0.2)
		n.SetWeight(4, 6, -0.3)
		n.SetWeight(5, 6, -0.2)

		// Set the biases
		n.SetBias(4, -0.4)
		n.SetBias(5, 0.2)
		n.SetBias(6, 0.1)

		// Create the Activation vector
		// NewDense is rows then columns
		a := mat.NewDense(6, 1, make([]float64, 6))
		// Set is rows then columns
		a.Set(0, 0, 1)
		a.Set(2, 0, 1)

		// ROBOTS ACTIVATE
		n.Activate(a, 2)
		Convey("The feed-forward results should be right...", func() {
			So(a.At(5, 0), ShouldAlmostEqual, 0.474, 0.01)
			So(a.At(4, 0), ShouldAlmostEqual, 0.525, 0.01)
			So(a.At(3, 0), ShouldAlmostEqual, 0.332, 0.01)

			// Set the observed error on the output node
			e := mat.NewDense(6, 1, make([]float64, 6))
			e.Set(5, 0, 1.0-a.At(5, 0))

			// Run back-propagated error
			b := n.Error(a, e, 2)
			Convey("The back-prop results should be right...", func() {
				So(b.At(5, 0), ShouldAlmostEqual, 0.1311, 0.001)
				So(b.At(4, 0), ShouldAlmostEqual, -0.0065, 0.001)
				So(b.At(3, 0), ShouldAlmostEqual, -0.0087, 0.001)
				So(b.At(2, 0), ShouldAlmostEqual, 0.000)
				So(b.At(1, 0), ShouldAlmostEqual, 0.000)
				So(b.At(0, 0), ShouldAlmostEqual, 0.000)

				Convey("The weight update results should be right...", func() {
					n.UpdateWeights(a, b, 0.9)
					for i := 1; i <= 6; i++ {
						for j := 1; j <= 6; j++ {
							v := n.GetWeight(i, j)
							switch i {
							case 1:
								switch j {
								case 1:
									So(v, ShouldAlmostEqual, 1.000)
								case 4:
									So(v, ShouldAlmostEqual, 0.192, 0.001)
								case 5:
									So(v, ShouldAlmostEqual, -0.306, 0.001)
								default:
									So(v, ShouldAlmostEqual, 0.000)
								}
							case 2:
								switch j {
								case 2:
									So(v, ShouldAlmostEqual, 1.000)
								case 4:
									So(v, ShouldAlmostEqual, 0.400, 0.001)
								case 5:
									So(v, ShouldAlmostEqual, 0.100, 0.001)
								default:
									So(v, ShouldAlmostEqual, 0.000)
								}
							case 3:
								switch j {
								case 3:
									So(v, ShouldAlmostEqual, 1.000)
								case 4:
									So(v, ShouldAlmostEqual, -0.508, 0.001)
								case 5:
									So(v, ShouldAlmostEqual, 0.194, 0.001)
								default:
									So(v, ShouldAlmostEqual, 0.000)
								}
							case 4:
								switch j {
								case 6:
									So(v, ShouldAlmostEqual, -0.261, 0.001)
								default:
									So(v, ShouldAlmostEqual, 0.000)
								}
							case 5:
								switch j {
								case 6:
									So(v, ShouldAlmostEqual, -0.138, 0.001)
								default:
									So(v, ShouldAlmostEqual, 0.000)
								}

							default:
								So(v, ShouldAlmostEqual, 0.000)
							}
						}
					}
				})

				Convey("The bias update results should be right...", func() {
					n.UpdateBias(b, 0.9)
					So(n.GetBias(6), ShouldAlmostEqual, 0.218, 0.001)
					So(n.GetBias(5), ShouldAlmostEqual, 0.194, 0.001)
					So(n.GetBias(4), ShouldAlmostEqual, -0.408, 0.001)
				})

			})
		})

	})

}
