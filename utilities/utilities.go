// Package utilities implements a host of helpful miscellaneous functions to the library.
package utilities

import (
	"fmt"
	rand "math/rand"
	"sort"
	"strconv"

	mat64 "github.com/gonum/matrix/mat64"
)

type sortedIntMap struct {
	m map[int]float64
	s []int
}

func (sm *sortedIntMap) Len() int {
	return len(sm.m)
}

func (sm *sortedIntMap) Less(i, j int) bool {
	return sm.m[sm.s[i]] < sm.m[sm.s[j]]
}

func (sm *sortedIntMap) Swap(i, j int) {
	sm.s[i], sm.s[j] = sm.s[j], sm.s[i]
}

func SortIntMap(m map[int]float64) []int {
	sm := new(sortedIntMap)
	sm.m = m
	sm.s = make([]int, len(m))
	i := 0
	for key, _ := range m {
		sm.s[i] = key
		i++
	}
	sort.Sort(sm)
	return sm.s
}

type sortedStringMap struct {
	m map[string]int
	s []string
}

func (sm *sortedStringMap) Len() int {
	return len(sm.m)
}

func (sm *sortedStringMap) Less(i, j int) bool {
	return sm.m[sm.s[i]] < sm.m[sm.s[j]]
}

func (sm *sortedStringMap) Swap(i, j int) {
	sm.s[i], sm.s[j] = sm.s[j], sm.s[i]
}

func SortStringMap(m map[string]int) []string {
	sm := new(sortedStringMap)
	sm.m = m
	sm.s = make([]string, len(m))
	i := 0
	for key, _ := range m {
		sm.s[i] = key
		i++
	}
	sort.Sort(sm)
	return sm.s
}

func RandomArray(n int, k int) []float64 {
	ReturnedArray := make([]float64, n)
	for i := 0; i < n; i++ {
		ReturnedArray[i] = rand.Float64() * float64(rand.Intn(k))
	}
	return ReturnedArray
}

func ConvertLabelsToFloat(labels []string) []float64 {
	floats := make([]float64, 0)
	for _, elem := range labels {
		converted, err := strconv.ParseFloat(elem, 64)

		if err != nil {
			fmt.Println(err)
		}

		floats = append(floats, converted)
	}
	return floats
}

func FloatsToMatrix(floats []float64) *mat64.Dense {
	return mat64.NewDense(1, len(floats), floats)
}
