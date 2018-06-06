package neural

import (
	"github.com/sjwhitworth/golearn/base"
	. "github.com/smartystreets/goconvey/convey"
	"gonum.org/v1/gonum/mat"
	"testing"
)

func TestLayerStructureNoHidden(t *testing.T) {

	Convey("Creating a network...", t, func() {
		XORData, err := base.ParseCSVToInstances("xor.csv", false)
		So(err, ShouldEqual, nil)
		Convey("Create a MultiLayerNet with no layers...", func() {
			net := NewMultiLayerNet(make([]int, 0))
			net.MaxIterations = 0
			net.Fit(XORData)
			Convey("The network should be the right size...", func() {
				So(net.network.size, ShouldEqual, 3)
			})
			Convey("The right nodes should be connected in the network...", func() {
				So(net.network.GetWeight(1, 1), ShouldAlmostEqual, 1.000)
				So(net.network.GetWeight(2, 2), ShouldAlmostEqual, 1.000)
				So(net.network.GetWeight(3, 3), ShouldAlmostEqual, 0.000)
				So(net.network.GetWeight(1, 3), ShouldNotAlmostEqual, 0.000)
				So(net.network.GetWeight(2, 3), ShouldNotAlmostEqual, 0.000)
			})
		})
		Convey("Create a multilayer net with two hidden layers...", func() {
			net := NewMultiLayerNet([]int{3, 2})
			net.MaxIterations = 0
			net.Fit(XORData)
			Convey("The network should be the right size...", func() {
				So(net.network.size, ShouldEqual, 8)
			})
			Convey("The right nodes should be connected in the network...", func() {
				So(net.network.GetWeight(1, 1), ShouldAlmostEqual, 1.000)
				So(net.network.GetWeight(2, 2), ShouldAlmostEqual, 1.000)
				for i := 3; i <= 8; i++ {
					So(net.network.GetWeight(i, i), ShouldAlmostEqual, 0.000)
				}
				for i := 1; i <= 2; i++ {
					for j := 3; j <= 5; j++ {
						So(net.network.GetWeight(i, j), ShouldNotAlmostEqual, 0.000)
					}
				}
				for i := 3; i <= 5; i++ {
					for j := 6; j <= 7; j++ {
						So(net.network.GetWeight(i, j), ShouldNotAlmostEqual, 0.000)
					}
				}
				for i := 6; i <= 7; i++ {
					So(net.network.GetWeight(i, 8), ShouldNotAlmostEqual, 0.000)
				}
				for i := 8; i > 0; i-- {
					for j := i - 1; j > 0; j-- {
						So(net.network.GetWeight(i, j), ShouldAlmostEqual, 0.000)
					}
				}
			})
		})
		Convey("Create a MultiLayerNet with 1 hidden layer...", func() {
			net := NewMultiLayerNet([]int{3})
			net.LearningRate = 0.9
			net.MaxIterations = 0
			net.Fit(XORData)

			Convey("The network should be the right size...", func() {
				So(net.network.size, ShouldEqual, 6)
			})

			Convey("The right nodes should be connected in the network...", func() {
				So(net.network.GetWeight(1, 1), ShouldAlmostEqual, 1.000)
				So(net.network.GetWeight(2, 2), ShouldAlmostEqual, 1.000)
				So(net.network.GetWeight(1, 3), ShouldNotAlmostEqual, 0.000)
				So(net.network.GetWeight(1, 4), ShouldNotAlmostEqual, 0.000)
				So(net.network.GetWeight(1, 5), ShouldNotAlmostEqual, 0.000)

				So(net.network.GetWeight(2, 3), ShouldNotAlmostEqual, 0.000)
				So(net.network.GetWeight(2, 4), ShouldNotAlmostEqual, 0.000)
				So(net.network.GetWeight(2, 5), ShouldNotAlmostEqual, 0.000)

				So(net.network.GetWeight(3, 3), ShouldAlmostEqual, 0.000)
				So(net.network.GetWeight(3, 4), ShouldAlmostEqual, 0.000)
				So(net.network.GetWeight(3, 5), ShouldAlmostEqual, 0.000)

				So(net.network.GetWeight(4, 4), ShouldAlmostEqual, 0.000)
				So(net.network.GetWeight(4, 3), ShouldAlmostEqual, 0.000)
				So(net.network.GetWeight(4, 5), ShouldAlmostEqual, 0.000)

				So(net.network.GetWeight(5, 5), ShouldAlmostEqual, 0.000)
				So(net.network.GetWeight(5, 3), ShouldAlmostEqual, 0.000)
				So(net.network.GetWeight(5, 4), ShouldAlmostEqual, 0.000)

				So(net.network.GetWeight(3, 6), ShouldNotAlmostEqual, 0.000)
				So(net.network.GetWeight(4, 6), ShouldNotAlmostEqual, 0.000)
				So(net.network.GetWeight(5, 6), ShouldNotAlmostEqual, 0.000)

				for i := 1; i <= 6; i++ {
					So(net.network.GetWeight(6, i), ShouldAlmostEqual, 0.000)
				}

			})
		})
	})

}

func TestLayeredXOR(t *testing.T) {

	Convey("Given an XOR dataset...", t, func() {

		XORData, err := base.ParseCSVToInstances("xor.csv", false)
		So(err, ShouldEqual, nil)

		net := NewMultiLayerNet([]int{3})
		net.MaxIterations = 20000
		net.Fit(XORData)

		Convey("After running for 20000 iterations, should have some predictive power...", func() {

			Convey("The right nodes should be connected in the network...", func() {
				So(net.network.GetWeight(1, 1), ShouldAlmostEqual, 1.000)
				So(net.network.GetWeight(2, 2), ShouldAlmostEqual, 1.000)

				for i := 1; i <= 6; i++ {
					So(net.network.GetWeight(6, i), ShouldAlmostEqual, 0.000)
				}

			})
			out := mat.NewDense(6, 1, []float64{1.0, 0.0, 0.0, 0.0, 0.0, 0.0})
			net.network.Activate(out, 2)
			So(out.At(5, 0), ShouldAlmostEqual, 1.0, 0.1)

			Convey("And Predict() should do OK too...", func() {

				pred := net.Predict(XORData)

				for _, a := range pred.AllAttributes() {
					af, ok := a.(*base.FloatAttribute)
					So(ok, ShouldBeTrue)

					af.Precision = 1
				}

				So(base.GetClass(pred, 0), ShouldEqual, "0.0")
				So(base.GetClass(pred, 1), ShouldEqual, "1.0")
				So(base.GetClass(pred, 2), ShouldEqual, "1.0")
				So(base.GetClass(pred, 3), ShouldEqual, "0.0")

			})
		})

	})

}

func TestLayeredXORInline(t *testing.T) {

	Convey("Given an inline XOR dataset...", t, func() {

		data := mat.NewDense(4, 3, []float64{
			1, 0, 1,
			0, 1, 1,
			0, 0, 0,
			1, 1, 0,
		})

		XORData := base.InstancesFromMat64(4, 3, data)
		classAttr := base.GetAttributeByName(XORData, "2")
		XORData.AddClassAttribute(classAttr)

		net := NewMultiLayerNet([]int{3})
		net.MaxIterations = 20000
		net.Fit(XORData)

		Convey("After running for 20000 iterations, should have some predictive power...", func() {

			Convey("The right nodes should be connected in the network...", func() {
				So(net.network.GetWeight(1, 1), ShouldAlmostEqual, 1.000)
				So(net.network.GetWeight(2, 2), ShouldAlmostEqual, 1.000)

				for i := 1; i <= 6; i++ {
					So(net.network.GetWeight(6, i), ShouldAlmostEqual, 0.000)
				}

			})
			out := mat.NewDense(6, 1, []float64{1.0, 0.0, 0.0, 0.0, 0.0, 0.0})
			net.network.Activate(out, 2)
			So(out.At(5, 0), ShouldAlmostEqual, 1.0, 0.1)

			Convey("And Predict() should do OK too...", func() {

				pred := net.Predict(XORData)

				for _, a := range pred.AllAttributes() {
					af, ok := a.(*base.FloatAttribute)
					So(ok, ShouldBeTrue)

					af.Precision = 1
				}

				So(base.GetClass(pred, 0), ShouldEqual, "1.0")
				So(base.GetClass(pred, 1), ShouldEqual, "1.0")
				So(base.GetClass(pred, 2), ShouldEqual, "0.0")
				So(base.GetClass(pred, 3), ShouldEqual, "0.0")

			})
		})

	})

}
