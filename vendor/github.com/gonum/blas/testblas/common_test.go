package testblas

import (
	"math"
	"testing"

	"github.com/gonum/blas"
	"github.com/gonum/floats"
)

func TestFlattenBanded(t *testing.T) {
	for i, test := range []struct {
		dense     [][]float64
		ku        int
		kl        int
		condensed [][]float64
	}{
		{
			dense:     [][]float64{{3}},
			ku:        0,
			kl:        0,
			condensed: [][]float64{{3}},
		},
		{
			dense: [][]float64{
				{3, 4, 0},
			},
			ku: 1,
			kl: 0,
			condensed: [][]float64{
				{3, 4},
			},
		},
		{
			dense: [][]float64{
				{3, 4, 0, 0, 0},
			},
			ku: 1,
			kl: 0,
			condensed: [][]float64{
				{3, 4},
			},
		},
		{
			dense: [][]float64{
				{3, 4, 0},
				{0, 5, 8},
				{0, 0, 2},
				{0, 0, 0},
				{0, 0, 0},
			},
			ku: 1,
			kl: 0,
			condensed: [][]float64{
				{3, 4},
				{5, 8},
				{2, math.NaN()},
				{math.NaN(), math.NaN()},
				{math.NaN(), math.NaN()},
			},
		},
		{
			dense: [][]float64{
				{3, 4, 6},
				{0, 5, 8},
				{0, 0, 2},
				{0, 0, 0},
				{0, 0, 0},
			},
			ku: 2,
			kl: 0,
			condensed: [][]float64{
				{3, 4, 6},
				{5, 8, math.NaN()},
				{2, math.NaN(), math.NaN()},
				{math.NaN(), math.NaN(), math.NaN()},
				{math.NaN(), math.NaN(), math.NaN()},
			},
		},
		{
			dense: [][]float64{
				{3, 4, 6},
				{1, 5, 8},
				{0, 6, 2},
				{0, 0, 7},
				{0, 0, 0},
			},
			ku: 2,
			kl: 1,
			condensed: [][]float64{
				{math.NaN(), 3, 4, 6},
				{1, 5, 8, math.NaN()},
				{6, 2, math.NaN(), math.NaN()},
				{7, math.NaN(), math.NaN(), math.NaN()},
				{math.NaN(), math.NaN(), math.NaN(), math.NaN()},
			},
		},
		{
			dense: [][]float64{
				{1, 2, 0},
				{3, 4, 5},
				{6, 7, 8},
				{0, 9, 10},
				{0, 0, 11},
			},
			ku: 1,
			kl: 2,
			condensed: [][]float64{
				{math.NaN(), math.NaN(), 1, 2},
				{math.NaN(), 3, 4, 5},
				{6, 7, 8, math.NaN()},
				{9, 10, math.NaN(), math.NaN()},
				{11, math.NaN(), math.NaN(), math.NaN()},
			},
		},
		{
			dense: [][]float64{
				{1, 0, 0},
				{3, 4, 0},
				{6, 7, 8},
				{0, 9, 10},
				{0, 0, 11},
			},
			ku: 0,
			kl: 2,
			condensed: [][]float64{
				{math.NaN(), math.NaN(), 1},
				{math.NaN(), 3, 4},
				{6, 7, 8},
				{9, 10, math.NaN()},
				{11, math.NaN(), math.NaN()},
			},
		},
		{
			dense: [][]float64{
				{1, 0, 0, 0, 0},
				{3, 4, 0, 0, 0},
				{1, 3, 5, 0, 0},
			},
			ku: 0,
			kl: 2,
			condensed: [][]float64{
				{math.NaN(), math.NaN(), 1},
				{math.NaN(), 3, 4},
				{1, 3, 5},
			},
		},
	} {
		condensed := flattenBanded(test.dense, test.ku, test.kl)
		correct := flatten(test.condensed)
		if !floats.Same(condensed, correct) {
			t.Errorf("Case %v mismatch. Want %v, got %v.", i, correct, condensed)
		}
	}
}

func TestFlattenTriangular(t *testing.T) {
	for i, test := range []struct {
		a   [][]float64
		ans []float64
		ul  blas.Uplo
	}{
		{
			a: [][]float64{
				{1, 2, 3},
				{0, 4, 5},
				{0, 0, 6},
			},
			ul:  blas.Upper,
			ans: []float64{1, 2, 3, 4, 5, 6},
		},
		{
			a: [][]float64{
				{1, 0, 0},
				{2, 3, 0},
				{4, 5, 6},
			},
			ul:  blas.Lower,
			ans: []float64{1, 2, 3, 4, 5, 6},
		},
	} {
		a := flattenTriangular(test.a, test.ul)
		if !floats.Equal(a, test.ans) {
			t.Errorf("Case %v. Want %v, got %v.", i, test.ans, a)
		}
	}
}
