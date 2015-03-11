package pca

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
	"github.com/gonum/matrix/mat64"

)

func TestPCAWithZeroComponents(t *testing.T){
	Convey("Set to pca 0 components with first matrix", t, func(){
		X1 := mat64.NewDense(3,7, []float64{6,5,4,3,8,2,9,5,1,10,2,3,8,7,5,14,2,3,6,3,2})
		pca := NewPCA(0)
		rows, cols := pca.Transform(X1).Dims()
		So(rows, ShouldEqual, 3)
		So(cols, ShouldEqual, 3)
	})

	Convey("Set to pca 0 components with second matrix", t, func(){
		X1 := mat64.NewDense(10,5, []float64{
		 0.52984892,  0.1141001 ,  0.91599294,  0.9574267 ,  0.15361222,
         0.07057588,  0.46371013,  0.73091854,  0.84641034,  0.08122213,
         0.96221946,  0.60367214,  0.69851546,  0.91965564,  0.27040597,
         0.03152856,  0.97912403,  0.39487038,  0.12232594,  0.18474705,
         0.77061953,  0.35898551,  0.78684562,  0.11638404,  0.88908044,
         0.35828086,  0.47214831,  0.95781755,  0.74762736,  0.59850757,
         0.07806127,  0.96940955,  0.15751804,  0.00973325,  0.85041635,
         0.02663938,  0.49755131,  0.57984119,  0.12233871,  0.47967853,
         0.63903222,  0.88556565,  0.79797963,  0.13345186,  0.37415535,
         0.60605207,  0.52067165,  0.91217494,  0.57148943,  0.92210331})
		pca := NewPCA(0)
		rows, cols := pca.Transform(X1).Dims()
		So(rows, ShouldEqual, 10)
		So(cols, ShouldEqual, 5)
	})
}

func TestPCAWithNComponents(t *testing.T){
	Convey("Set to pca 3 components with 5x5 matrix", t, func(){
		X := mat64.NewDense(5,5, [] float64{
		 0.23030838,  0.05669317,  0.3187813 ,  0.34455114,  0.98062806,
         0.38995469,  0.2996771 ,  0.99043575,  0.04443827,  0.99527955,
         0.27266308,  0.14068906,  0.46999473,  0.03296131,  0.90855405,
         0.28360708,  0.8839966 ,  0.81107014,  0.52673877,  0.59432817,
         0.64107253,  0.56165215,  0.79811756,  0.48845398,  0.20506649})
		pca := NewPCA(3)
		rows, cols := pca.Transform(X).Dims()
		So(rows, ShouldEqual, 5)
		So(cols, ShouldEqual, 3)
	})

	Convey("Set to pca 2 components with 3x5 matrix",t, func(){
		X := mat64.NewDense(3,5, [] float64{
		 0.12294845,  0.55170713,  0.67572832,  0.60615516,  0.38184551,
         0.93486821,  0.15120374,  0.89760169,  0.74715672,  0.81373931,
         0.42821569,  0.47457753,  0.18960954,  0.42466159,  0.34166049})
		pca := NewPCA(2)
		rows, cols := pca.Transform(X).Dims()
		So(rows, ShouldEqual,3)
		So(cols, ShouldEqual,2)
	})
}