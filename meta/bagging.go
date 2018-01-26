package meta

import (
	"fmt"
	"math/rand"
	"runtime"
	"strings"
	"sync"

	"github.com/amclay/golearn/base"
)

// BaggedModel trains base.Classifiers on subsets of the original
// Instances and combine the results through voting
type BaggedModel struct {
	base.BaseClassifier
	Models             []base.Classifier
	RandomFeatures     int
	lock               sync.Mutex
	selectedAttributes map[int][]base.Attribute
	fitOn              base.FixedDataGrid
}

// generateTrainingAttrs selects RandomFeatures number of base.Attributes from
// the provided base.Instances.
func (b *BaggedModel) generateTrainingAttrs(model int, from base.FixedDataGrid) []base.Attribute {
	ret := make([]base.Attribute, 0)
	attrs := base.NonClassAttributes(from)
	if b.RandomFeatures == 0 {
		ret = attrs
	} else {
		for {
			if len(ret) >= b.RandomFeatures {
				break
			}
			attrIndex := rand.Intn(len(attrs))
			attr := attrs[attrIndex]
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
	for _, a := range from.AllClassAttributes() {
		ret = append(ret, a)
	}
	b.lock.Lock()
	b.selectedAttributes[model] = ret
	b.lock.Unlock()
	return ret
}

// generatePredictionInstances returns a modified version of the
// requested base.Instances with only the base.Attributes selected
// for training the model.
func (b *BaggedModel) generatePredictionInstances(model int, from base.FixedDataGrid) base.FixedDataGrid {
	selected := b.selectedAttributes[model]
	return base.NewInstancesViewFromAttrs(from, selected)
}

// generateTrainingInstances generates RandomFeatures number of
// attributes and returns a modified version of base.Instances
// for training the model
func (b *BaggedModel) generateTrainingInstances(model int, from base.FixedDataGrid) base.FixedDataGrid {
	_, rows := from.Size()
	insts := base.SampleWithReplacement(from, rows)
	selected := b.generateTrainingAttrs(model, from)
	return base.NewInstancesViewFromAttrs(insts, selected)
}

// AddModel adds a base.Classifier to the current model
func (b *BaggedModel) AddModel(m base.Classifier) {
	b.Models = append(b.Models, m)
}

// Fit generates and trains each model on a randomised subset of
// Instances.
func (b *BaggedModel) Fit(from base.FixedDataGrid) {
	var wait sync.WaitGroup
	b.selectedAttributes = make(map[int][]base.Attribute)
	for i, m := range b.Models {
		wait.Add(1)
		go func(c base.Classifier, f base.FixedDataGrid, model int) {
			l := b.generateTrainingInstances(model, f)
			c.Fit(l)
			wait.Done()
		}(m, from, i)
	}
	wait.Wait()
	b.fitOn = base.NewStructuralCopy(from)
}

// Predict gathers predictions from all the classifiers
// and outputs the most common (majority) class
//
// IMPORTANT: in the event of a tie, the first class which
// achieved the tie value is output.
func (b *BaggedModel) Predict(from base.FixedDataGrid) (base.FixedDataGrid, error) {
	n := runtime.NumCPU()
	// Channel to receive the results as they come in
	votes := make(chan base.DataGrid, n)
	// Count the votes for each class
	voting := make(map[int](map[string]int))

	// Create a goroutine to collect the votes
	var votingwait sync.WaitGroup
	votingwait.Add(1)
	go func() {
		for { // Need to resolve the voting problem
			incoming, ok := <-votes
			if ok {
				cSpecs := base.ResolveAttributes(incoming, incoming.AllClassAttributes())
				incoming.MapOverRows(cSpecs, func(row [][]byte, predRow int) (bool, error) {
					// Check if we've seen this class before...
					if _, ok := voting[predRow]; !ok {
						// If we haven't, create an entry
						voting[predRow] = make(map[string]int)
						// Continue on the current row
					}
					voting[predRow][base.GetClass(incoming, predRow)]++
					return true, nil
				})
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
					v, _ := c.Predict(l)
					votes <- v
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
	ret := base.GeneratePredictionVector(from)
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
		base.SetClass(ret, i, maxClass)
	}
	return ret, nil
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

// GetMetadata returns required serialization information for this classifier
func (b *BaggedModel) GetMetadata() base.ClassifierMetadataV1 {

	return base.ClassifierMetadataV1{
		FormatVersion:      1,
		ClassifierName:     "BaggedModel",
		ClassifierVersion:  "1.0",
		ClassifierMetadata: nil,
	}

}

func (b *BaggedModel) Save(filePath string) error {
	writer, err := base.CreateSerializedClassifierStub(filePath, b.GetMetadata())
	if err != nil {
		return err
	}
	err = b.SaveWithPrefix(writer, "")
	writer.Close()
	return err
}

func (b *BaggedModel) Load(filePath string) error {
	reader, err := base.ReadSerializedClassifierStub(filePath)
	if err != nil {
		return err
	}

	err = b.LoadWithPrefix(reader, "")
	reader.Close()
	return err
}

/* type BaggedModel struct {
	base.BaseClassifier
	Models             []base.Classifier
	RandomFeatures     int
	lock               sync.Mutex
	selectedAttributes map[int][]base.Attribute, always RandomFeatures in length
}*/
func (b *BaggedModel) SaveWithPrefix(writer *base.ClassifierSerializer, prefix string) error {
	b.lock.Lock()
	defer b.lock.Unlock()

	pI := func(n string, i int) string {
		return writer.Prefix(prefix, writer.Prefix(n, fmt.Sprintf("%d", i)))
	}

	// Export the number of random features
	randomFeaturesKey := writer.Prefix(prefix, "NUM_RANDOM_FEATURES")
	err := writer.WriteU64ForKey(randomFeaturesKey, uint64(b.RandomFeatures))
	if err != nil {
		return base.DescribeError("Can't write NUM_RANDOM_FEATURES", err)
	}

	// Write the number of classifiers
	classifiersKey := writer.Prefix(prefix, "NUM_CLASSIFIERS")
	err = writer.WriteU64ForKey(classifiersKey, uint64(len(b.Models)))
	if err != nil {
		return base.DescribeError("Can't write NUM_CLASSIFIERS", err)
	}

	// Save the classifiers
	for i, c := range b.Models {
		clsPrefix := fmt.Sprintf("%s/", pI("CLASSIFIERS", i))
		err = c.SaveWithPrefix(writer, clsPrefix)
		if err != nil {
			return base.FormatError(err, "Can't save classifier %d", i)
		}
	}

	// Save the instances template
	err = writer.WriteInstancesForKey(writer.Prefix(prefix, "REFERENCE_INSTANCES"), b.fitOn, false)
	if err != nil {
		return base.DescribeError("Can't write REFERENCE_INSTANCES", err)
	}

	// Save the selectedAttributes map
	selectedAttributesKey := writer.Prefix(prefix, "SELECTED_ATTRIBUTES")
	ser := make(map[int][][]byte)
	for key := range b.selectedAttributes {
		ser[key] = make([][]byte, 0)
		for _, a := range b.selectedAttributes[key] {
			bytes, err := base.SerializeAttribute(a)
			if err != nil {
				return base.DescribeError("Can't serialize Attribute", err)
			}
			ser[key] = append(ser[key], bytes)
		}
	}

	err = writer.WriteJSONForKey(selectedAttributesKey, ser)
	if err != nil {
		return base.DescribeError("Can't write selected attributes map", err)
	}

	return nil
}

// Remember: have to add the models before you use this.
func (b *BaggedModel) LoadWithPrefix(reader *base.ClassifierDeserializer, prefix string) error {
	b.lock.Lock()
	defer b.lock.Unlock()

	pI := func(n string, i int) string {
		return reader.Prefix(prefix, reader.Prefix(n, fmt.Sprintf("%d", i)))
	}

	// Read the essential info
	randomFeaturesKey := reader.Prefix(prefix, "NUM_RANDOM_FEATURES")
	randomFeatures, err := reader.GetU64ForKey(randomFeaturesKey)
	if err != nil {
		return base.DescribeError("Can't read NUM_RANDOM_FEATURES", err)
	}

	b.RandomFeatures = int(randomFeatures)

	// Reload the classifiers
	for i, m := range b.Models {
		clsPrefix := fmt.Sprintf("%s/", pI("CLASSIFIERS", i))
		err := m.LoadWithPrefix(reader, clsPrefix)
		if err != nil {
			return base.DescribeError("Can't read classifier", err)
		}
	}

	// Load the instances template
	tmp, err := reader.GetInstancesForKey(reader.Prefix(prefix, "REFERENCE_INSTANCES"))
	if err != nil {
		return base.DescribeError("Can't read REFERENCE_INSTACES", err)
	}
	b.fitOn = tmp

	// Reload the selected attributes
	selectedAttributesKey := reader.Prefix(prefix, "SELECTED_ATTRIBUTES")
	ser := make(map[int][][]byte)
	err = reader.GetJSONForKey(selectedAttributesKey, &ser)
	if err != nil {
		return base.DescribeError("Can't reload selected attributes", err)
	}

	b.selectedAttributes = make(map[int][]base.Attribute)

	for key := range ser {
		b.selectedAttributes[key] = make([]base.Attribute, len(ser[key]))
		for i, attrBytes := range ser[key] {
			attr, err := base.DeserializeAttribute(attrBytes)
			if err != nil {
				return base.DescribeError("Can't deserialize attribute", err)
			}
			attrNew, err := base.ReplaceDeserializedAttributeWithVersionFromInstances(attr, b.fitOn)
			if err != nil {
				return base.DescribeError("Can't replace attribute", err)
			}
			b.selectedAttributes[key][i] = attrNew
		}
	}
	return nil
}
