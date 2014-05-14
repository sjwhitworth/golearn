package meta

import (
	"fmt"
	base "github.com/sjwhitworth/golearn/base"
	"math/rand"
	"runtime"
	"strings"
	"time"
)

// BaggedModels train Classifiers on subsets of the original
// Instances and combine the results through voting
type BaggedModel struct {
	base.BaseClassifier
	Models           []base.Classifier
	SelectedFeatures map[int][]base.Attribute
	// If this is greater than 0, select up to d features
	// for feeding into each classifier
	RandomFeatures int
}

func (b *BaggedModel) generateRandomAttributes(from *base.Instances) []base.Attribute {
	if b.RandomFeatures > from.GetAttributeCount()-1 {
		panic("Can't have more random features")
	}
	ret := make([]base.Attribute, 0)
	for {
		if len(ret) > b.RandomFeatures {
			break
		}
		attrIndex := rand.Intn(from.GetAttributeCount())
		if attrIndex == from.ClassIndex {
			continue
		}
		matched := false
		newAttr := from.GetAttr(attrIndex)
		for _, a := range ret {
			if a.Equals(newAttr) {
				matched = true
				break
			}
		}
		if !matched {
			ret = append(ret, newAttr)
		}
	}
	ret = append(ret, from.GetClassAttr())
	return ret
}

func (b *BaggedModel) generateTrainingInstances(from *base.Instances) ([]base.Attribute, *base.Instances) {

	var attrs []base.Attribute
	from = from.SampleWithReplacement(from.Rows)

	if b.RandomFeatures > 0 {
		attrs = b.generateRandomAttributes(from)
		from = from.SelectAttributes(attrs)
	} else {
		attrs = make([]base.Attribute, 0)
	}

	return attrs, from
}

func (b *BaggedModel) generateTestingInstances(from *base.Instances, model int) *base.Instances {
	attrs := b.SelectedFeatures[model]
	return from.SelectAttributes(attrs)
}

func (b *BaggedModel) AddModel(m base.Classifier) {
	b.Models = append(b.Models, m)
}

// Train generates and trains each model on a randomised subset of
// Instances.
func (b *BaggedModel) Fit(from *base.Instances) {
	n := runtime.GOMAXPROCS(0)
	block := make(chan bool, n)
	for i, m := range b.Models {
		go func(c base.Classifier, f *base.Instances) {
			a, f := b.generateTrainingInstances(f)
			b.SelectedFeatures[i] = a
			rand.Seed(time.Now().UnixNano())
			c.Fit(f)
			block <- true
		}(m, from)
	}
	for i := 0; i < len(b.Models); i++ {
		<-block
	}
}

// Predict gathers predictions from all the classifiers
// and outputs the most common (majority) class
//
// IMPORTANT: in the event of a tie, the first class which
// achieved the tie value is output.
func (b *BaggedModel) Predict(from *base.Instances) *base.Instances {
	n := runtime.GOMAXPROCS(0)
	// Channel to receive the results as they come in
	votes := make(chan *base.Instances, n)
	// Dispatch prediction generation
	for i, m := range b.Models {
		go func(c base.Classifier, f *base.Instances) {
			f = b.generateTestingInstances(f, i)
			p := c.Predict(f)
			votes <- p
		}(m, from)
	}
	// Count the votes for each class
	voting := make(map[int](map[string]int))
	for _ = range b.Models { // Have to do this - deadlocks otherwise
		incoming := <-votes
		// Step through each prediction
		for j := 0; j < incoming.Rows; j++ {
			// Check if we've seen this class before...
			if _, ok := voting[j]; !ok {
				// If we haven't, create an entry
				voting[j] = make(map[string]int)
				// Continue on the current row
				j--
				continue
			}
			voting[j][incoming.GetClass(j)]++
		}
	}

	// Generate the overall consensus
	ret := from.GeneratePredictionVector()
	for i := range voting {
		maxClass := ""
		maxCount := 0
		// Find the most popular class
		for c := range voting[i] {
			votes := voting[i][c]
			if votes > maxCount {
				maxClass = c
				maxCount = votes
			}
		}
		ret.SetAttrStr(i, 0, maxClass)
	}
	return ret
}

// String returns a human-readable representation of the
// BaggedModel and everything it contains
func (b *BaggedModel) String() string {
	children := make([]string, 0)
	for i, m := range b.Models {
		children = append(children, fmt.Sprintf("%d: %s", i, m))
	}
	return fmt.Sprintf("BaggedModel(\n%s)", strings.Join(children, "\n\t"))
}
