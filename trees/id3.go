package trees

import (
	"bytes"
	"fmt"
	"github.com/sjwhitworth/golearn/base"
	eval "github.com/sjwhitworth/golearn/evaluation"
	"sort"
)

// NodeType determines whether a DecisionTreeNode is a leaf or not
type NodeType int

const (
	// LeafNode means there are no children
	LeafNode NodeType = 1
	// RuleNode means we should look at the next attribute value
	RuleNode NodeType = 2
)

// RuleGenerator implementations analyse instances and determine
// the best value to split on
type RuleGenerator interface {
	GenerateSplitAttribute(base.FixedDataGrid) base.Attribute
}

// DecisionTreeNode represents a given portion of a decision tree
type DecisionTreeNode struct {
	Type      NodeType
	Children  map[string]*DecisionTreeNode
	SplitAttr base.Attribute
	ClassDist map[string]int
	Class     string
	ClassAttr base.Attribute
}

func getClassAttr(from base.FixedDataGrid) base.Attribute {
	allClassAttrs := from.AllClassAttributes()
	return allClassAttrs[0]
}

// InferID3Tree builds a decision tree using a RuleGenerator
// from a set of Instances (implements the ID3 algorithm)
func InferID3Tree(from base.FixedDataGrid, with RuleGenerator) *DecisionTreeNode {
	// Count the number of classes at this node
	classes := base.GetClassDistribution(from)
	// If there's only one class, return a DecisionTreeLeaf with
	// the only class available
	if len(classes) == 1 {
		maxClass := ""
		for i := range classes {
			maxClass = i
		}
		ret := &DecisionTreeNode{
			LeafNode,
			nil,
			nil,
			classes,
			maxClass,
			getClassAttr(from),
		}
		return ret
	}

	// Only have the class attribute
	maxVal := 0
	maxClass := ""
	for i := range classes {
		if classes[i] > maxVal {
			maxClass = i
			maxVal = classes[i]
		}
	}

	// If there are no more Attributes left to split on,
	// return a DecisionTreeLeaf with the majority class
	cols, _ := from.Size()
	if cols == 2 {
		ret := &DecisionTreeNode{
			LeafNode,
			nil,
			nil,
			classes,
			maxClass,
			getClassAttr(from),
		}
		return ret
	}

	// Generate a return structure
	ret := &DecisionTreeNode{
		RuleNode,
		nil,
		nil,
		classes,
		maxClass,
		getClassAttr(from),
	}

	// Generate the splitting attribute
	splitOnAttribute := with.GenerateSplitAttribute(from)
	if splitOnAttribute == nil {
		// Can't determine, just return what we have
		return ret
	}
	// Split the attributes based on this attribute's value
	splitInstances := base.DecomposeOnAttributeValues(from, splitOnAttribute)
	// Create new children from these attributes
	ret.Children = make(map[string]*DecisionTreeNode)
	for k := range splitInstances {
		newInstances := splitInstances[k]
		ret.Children[k] = InferID3Tree(newInstances, with)
	}
	ret.SplitAttr = splitOnAttribute
	return ret
}

// getNestedString returns the contents of node d
// prefixed by level number of tags (also prints children)
func (d *DecisionTreeNode) getNestedString(level int) string {
	buf := bytes.NewBuffer(nil)
	tmp := bytes.NewBuffer(nil)
	for i := 0; i < level; i++ {
		tmp.WriteString("\t")
	}
	buf.WriteString(tmp.String())
	if d.Children == nil {
		buf.WriteString(fmt.Sprintf("Leaf(%s)", d.Class))
	} else {
		buf.WriteString(fmt.Sprintf("Rule(%s)", d.SplitAttr.GetName()))
		keys := make([]string, 0)
		for k := range d.Children {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		for _, k := range keys {
			buf.WriteString("\n")
			buf.WriteString(tmp.String())
			buf.WriteString("\t")
			buf.WriteString(k)
			buf.WriteString("\n")
			buf.WriteString(d.Children[k].getNestedString(level + 1))
		}
	}
	return buf.String()
}

// String returns a human-readable representation of a given node
// and it's children
func (d *DecisionTreeNode) String() string {
	return d.getNestedString(0)
}

// computeAccuracy is a helper method for Prune()
func computeAccuracy(predictions base.FixedDataGrid, from base.FixedDataGrid) float64 {
	cf := eval.GetConfusionMatrix(from, predictions)
	return eval.GetAccuracy(cf)
}

// Prune eliminates branches which hurt accuracy
func (d *DecisionTreeNode) Prune(using base.FixedDataGrid) {
	// If you're a leaf, you're already pruned
	if d.Children == nil {
		return
	}
	if d.SplitAttr == nil {
		return
	}

	// Recursively prune children of this node
	sub := base.DecomposeOnAttributeValues(using, d.SplitAttr)
	for k := range d.Children {
		if sub[k] == nil {
			continue
		}
		subH, subV := sub[k].Size()
		if subH == 0 || subV == 0 {
			continue
		}
		d.Children[k].Prune(sub[k])
	}

	// Get a baseline accuracy
	baselineAccuracy := computeAccuracy(d.Predict(using), using)

	// Speculatively remove the children and re-evaluate
	tmpChildren := d.Children
	d.Children = nil
	newAccuracy := computeAccuracy(d.Predict(using), using)

	// Keep the children removed if better, else restore
	if newAccuracy < baselineAccuracy {
		d.Children = tmpChildren
	}
}

// Predict outputs a base.Instances containing predictions from this tree
func (d *DecisionTreeNode) Predict(what base.FixedDataGrid) base.FixedDataGrid {
	predictions := base.GeneratePredictionVector(what)
	classAttr := getClassAttr(predictions)
	classAttrSpec, err := predictions.GetAttribute(classAttr)
	if err != nil {
		panic(err)
	}
	predAttrs := base.AttributeDifferenceReferences(what.AllAttributes(), predictions.AllClassAttributes())
	predAttrSpecs := base.ResolveAttributes(what, predAttrs)
	what.MapOverRows(predAttrSpecs, func(row [][]byte, rowNo int) (bool, error) {
		cur := d
		for {
			if cur.Children == nil {
				predictions.Set(classAttrSpec, rowNo, classAttr.GetSysValFromString(cur.Class))
				break
			} else {
				at := cur.SplitAttr
				ats, err := what.GetAttribute(at)
				if err != nil {
					predictions.Set(classAttrSpec, rowNo, classAttr.GetSysValFromString(cur.Class))
					break
				}

				classVar := ats.GetAttribute().GetStringFromSysVal(what.Get(ats, rowNo))
				if next, ok := cur.Children[classVar]; ok {
					cur = next
				} else {
					var bestChild string
					for c := range cur.Children {
						bestChild = c
						if c > classVar {
							break
						}
					}
					cur = cur.Children[bestChild]
				}
			}
		}
		return true, nil
	})
	return predictions
}

//
// ID3 Tree type
//

// ID3DecisionTree represents an ID3-based decision tree
// using the Information Gain metric to select which attributes
// to split on at each node.
type ID3DecisionTree struct {
	base.BaseClassifier
	Root       *DecisionTreeNode
	PruneSplit float64
}

// NewID3DecisionTree returns a new ID3DecisionTree with the specified test-prune
// ratio. Of the ratio is less than 0.001, the tree isn't pruned
func NewID3DecisionTree(prune float64) *ID3DecisionTree {
	return &ID3DecisionTree{
		base.BaseClassifier{},
		nil,
		prune,
	}
}

// Fit builds the ID3 decision tree
func (t *ID3DecisionTree) Fit(on base.FixedDataGrid) {
	rule := new(InformationGainRuleGenerator)
	if t.PruneSplit > 0.001 {
		trainData, testData := base.InstancesTrainTestSplit(on, t.PruneSplit)
		t.Root = InferID3Tree(trainData, rule)
		t.Root.Prune(testData)
	} else {
		t.Root = InferID3Tree(on, rule)
	}
}

// Predict outputs predictions from the ID3 decision tree
func (t *ID3DecisionTree) Predict(what base.FixedDataGrid) base.FixedDataGrid {
	return t.Root.Predict(what)
}

// String returns a human-readable version of this ID3 tree
func (t *ID3DecisionTree) String() string {
	return fmt.Sprintf("ID3DecisionTree(%s\n)", t.Root)
}
