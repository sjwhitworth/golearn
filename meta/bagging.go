package meta

import (
	"fmt"
	base "github.com/sjwhitworth/golearn/base"
	"runtime"
	"strings"
)

// BaggedModels train Classifiers on subsets of the original
// Instances and combine the results through voting
type BaggedModel struct {
	base.BaseClassifier
	Models []base.Classifier
}

func (b *BaggedModel) generateTrainingInstances(from *base.Instances) *base.Instances {
	from = from.SampleWithReplacement(from.Rows)
	return from
}

func (b *BaggedModel) AddModel(m base.Classifier) {
	b.Models = append(b.Models, m)
}

// Train generates and trains each model on a randomised subset of
// Instances.
func (b *BaggedModel) Fit(from *base.Instances) {
	n := runtime.GOMAXPROCS(0)
	block := make(chan bool, n)
	for _, m := range b.Models {
		go func(c base.Classifier, f *base.Instances) {
			f = b.generateTrainingInstances(f)
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
	for _, m := range b.Models {
		go func(c base.Classifier, f *base.Instances) {
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
