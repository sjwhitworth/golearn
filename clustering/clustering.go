/* This package implements clustering algorithms */
package clustering

import (
	"fmt"

	"github.com/sjwhitworth/golearn/base"
	"github.com/sjwhitworth/golearn/metrics/pairwise"
)

// ClusterParameters takes a number of variables common to all clustering
// algorithms.

type ClusterParameters struct {
	// Attributes represents the set of Attributes which
	// can be used for clustering
	Attributes []base.Attribute

	// Metric is used to compute pairwise distance
	Metric pairwise.PairwiseDistanceFunc
}

// ClusterMap contains the cluster identifier as a key, followed by a vector of point
// indices that cluster contains.
type ClusterMap map[int][]int

// Invert returns an alternative form of cluster map where the key represents the point
// index and the value represents the cluster index it's assigned to
func (ref ClusterMap) Invert() (map[int]int, error) {
	ret := make(map[int]int)
	for c := range ref {
		for _, p := range ref[c] {
			if _, ok := ret[p]; ok {
				return nil, fmt.Errorf("Not a valid cluster map (points appear in more than one cluster)")
			} else {
				ret[p] = c
			}
		}
	}
	return ret, nil
}

// Equals checks whether a bijection exists between two ClusterMaps (i.e. the clusters in one can
// be re-labelled to become the clusters of another)
func (ref ClusterMap) Equals(other ClusterMap) (bool, error) {
	if len(ref) != len(other) {
		return false, fmt.Errorf("ref and other do not contain the same number of clusters (%d and %d)", len(ref), len(other))
	}

	refInv, err := ref.Invert()
	if err != nil {
		return false, fmt.Errorf("ref: %s", err)
	}

	otherInv, err := other.Invert()
	if err != nil {
		return false, fmt.Errorf("other: %s", err)
	}

	clusterIdMap := make(map[int]int)

	// Range through each point index
	for p := range refInv {
		c1 := refInv[p]                // Get the cluster index of this point
		if c2, ok := otherInv[p]; ok { // Check if the other map has this point
			// if so, c2 is the point's cluster in the other map
			if c3, ok := clusterIdMap[c2]; ok { // what's our correspondance with c2?
				if c1 != c3 {
					// if c1 is not what we've currently got, error out
					return false, fmt.Errorf("ref point %d (cluster %d) is assigned to a different cluster (%d) in ref %+v", p, c2, c1, clusterIdMap)
				}
			} else {
				clusterIdMap[c2] = c1
			}
		} else {
			return false, fmt.Errorf("failed to find reference point %d in src", p)
		}
	}

	// Check that after transformation, key contains the same points
	arraysEqual := func(a1, a2 []int) bool {

		cnt := make(map[int]bool)
		for _, a := range a1 {
			cnt[a] = true
		}

		for _, a := range a2 {
			if _, ok := cnt[a]; !ok {
				return false
			}
		}

		return true

	}
	newMap := ClusterMap(make(map[int][]int))
	for cOld := range other {
		cNew := clusterIdMap[cOld]
		if !arraysEqual(ref[cNew], other[cOld]) {
			return false, fmt.Errorf("Re-labelled cluster %d => %d doesn't contain the same points (%d, %d)", cOld, cNew, ref[cNew], other[cOld])
		}
		newMap[cNew] = other[cOld]
	}

	return true, nil

}
