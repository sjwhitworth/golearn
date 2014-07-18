package meta

import (
	"fmt"
	base "github.com/sjwhitworth/golearn/base"
	"math/rand"
	"runtime"
	"strings"
	"sync"
)

// BaggedModel trains base.Classifiers on subsets of the original
// Instances and combine the results through voting
type BaggedModel struct {
	base.BaseClassifier
	Models             []base.Classifier
	RandomFeatures     int
	lock               sync.Mutex
	selectedAttributes map[int][]base.Attribute
}

// generateTrainingAttrs selects RandomFeatures number of base.Attributes from
// the provided base.Instances.
func (b *BaggedModel) generateTrainingAttrs(model int, from *base.Instances) []base.Attribute {
	ret := make([]base.Attribute, 0)
	if b.RandomFeatures == 0 {
		for j := 0; j < from.Cols; j++ {
			attr := from.GetAttr(j)
			ret = append(ret, attr)
		}
	} else {
		for {
			if len(ret) >= b.RandomFeatures {
				break
			}
			attrIndex := rand.Intn(from.Cols)
			if attrIndex == from.ClassIndex {
				continue
			}
			attr := from.GetAttr(attrIndex)
			matched := false
			for _, a := range ret {
				if a.Equals(attr) {
					matched = true
					break
				}
			}
			if !matched {
				ret = append(ret, attr)
			}
		}
	}
	ret = append(ret, from.GetClassAttr())
	b.lock.Lock()
	b.selectedAttributes[model] = ret
	b.lock.Unlock()
	return ret
}

// generatePredictionInstances returns a modified version of the
// requested base.Instances with only the base.Attributes selected
// for training the model.
func (b *BaggedModel) generatePredictionInstances(model int, from *base.Instances) *base.Instances {
	selected := b.selectedAttributes[model]
	return from.SelectAttributes(selected)
}

// generateTrainingInstances generates RandomFeatures number of
// attributes and returns a modified version of base.Instances
// for training the model
func (b *BaggedModel) generateTrainingInstances(model int, from *base.Instances) *base.Instances {
	insts := from.SampleWithReplacement(from.Rows)
	selected := b.generateTrainingAttrs(model, from)
	return insts.SelectAttributes(selected)
}

// AddModel adds a base.Classifier to the current model
func (b *BaggedModel) AddModel(m base.Classifier) {
	b.Models = append(b.Models, m)
}

// Train generates and trains each model on a randomised subset of
// Instances.
func (b *BaggedModel) Fit(from *base.Instances) {
	var wait sync.WaitGroup
	b.selectedAttributes = make(map[int][]base.Attribute)
	for i, m := range b.Models {
		wait.Add(1)
		go func(c base.Classifier, f *base.Instances, model int) {
			l := b.generateTrainingInstances(model, f)
			c.Fit(l)
			wait.Done()
		}(m, from, i)
	}
	wait.Wait()
}

// Predict gathers predictions from all the classifiers
// and outputs the most common (majority) class
//
// IMPORTANT: in the event of a tie, the first class which
// achieved the tie value is output.
func (b *BaggedModel) Predict(from *base.Instances) *base.Instances {
	n := runtime.NumCPU()
	// Channel to receive the results as they come in
	votes := make(chan *base.Instances, n)
	// Count the votes for each class
	voting := make(map[int](map[string]int))

	// Create a goroutine to collect the votes
	var votingwait sync.WaitGroup
	votingwait.Add(1)
	go func() {
		for {
			incoming, ok := <-votes
			if ok {
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
			} else {
				votingwait.Done()
				break
			}
		}
	}()

	// Create workers to process the predictions
	processpipe := make(chan int, n)
	var processwait sync.WaitGroup
	for i := 0; i < n; i++ {
		processwait.Add(1)
		go func() {
			for {
				if i, ok := <-processpipe; ok {
					c := b.Models[i]
					l := b.generatePredictionInstances(i, from)
					votes <- c.Predict(l)
				} else {
					processwait.Done()
					break
				}
			}
		}()
	}

	// Send all the models to the workers for prediction
	for i := range b.Models {
		processpipe <- i
	}
	close(processpipe) // Finished sending models to be predicted
	processwait.Wait() // Predictors all finished processing
	close(votes)       // Close the vote channel and allow it to drain
	votingwait.Wait()  // All the votes are in

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
