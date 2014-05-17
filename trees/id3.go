package trees

import (
	"bytes"
	"fmt"
	base "github.com/sjwhitworth/golearn/base"
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
	GenerateSplitAttribute(*base.Instances) base.Attribute
}

// DecisionTreeNode represents a given portion of a decision tree
type DecisionTreeNode struct {
	Type      NodeType
	Children  map[string]*DecisionTreeNode
	SplitAttr base.Attribute
	ClassDist map[string]int
	Class     string
	ClassAttr *base.Attribute
}

// InferID3Tree builds a decision tree using a RuleGenerator
// from a set of Instances (implements the ID3 algorithm)
func InferID3Tree(from *base.Instances, with RuleGenerator) *DecisionTreeNode {
	// Count the number of classes at this node
	classes := from.CountClassValues()
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
			from.GetClassAttrPtr(),
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
	if from.GetAttributeCount() == 2 {
		ret := &DecisionTreeNode{
			LeafNode,
			nil,
			nil,
			classes,
			maxClass,
			from.GetClassAttrPtr(),
		}
		return ret
	}

	ret := &DecisionTreeNode{
		RuleNode,
		nil,
		nil,
		classes,
		maxClass,
		from.GetClassAttrPtr(),
	}

	// Generate a return structure
	// Generate the splitting attribute
	splitOnAttribute := with.GenerateSplitAttribute(from)
	if splitOnAttribute == nil {
		// Can't determine, just return what we have
		return ret
	}
	// Split the attributes based on this attribute's value
	splitInstances := from.DecomposeOnAttributeValues(splitOnAttribute)
	// Create new children from these attributes
	ret.Children = make(map[string]*DecisionTreeNode)
	for k := range splitInstances {
		newInstances := splitInstances[k]
		ret.Children[k] = InferID3Tree(newInstances, with)
	}
	ret.SplitAttr = splitOnAttribute
	return ret
}

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

func computeAccuracy(predictions *base.Instances, from *base.Instances) float64 {
	cf := eval.GetConfusionMatrix(from, predictions)
	return eval.GetAccuracy(cf)
}

// Prune eliminates branches which hurt accuracy
func (d *DecisionTreeNode) Prune(using *base.Instances) {
	// If you're a leaf, you're already pruned
	if d.Children == nil {
		return
	} else {
		// Recursively prune children of this node
		sub := using.DecomposeOnAttributeValues(d.SplitAttr)
		for k := range d.Children {
			d.Children[k].Prune(sub[k])
		}
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
func (d *DecisionTreeNode) Predict(what *base.Instances) *base.Instances {
	outputAttrs := make([]base.Attribute, 1)
	outputAttrs[0] = what.GetClassAttr()
	predictions := base.NewInstances(outputAttrs, what.Rows)
	for i := 0; i < what.Rows; i++ {
		cur := d
		for {
			if cur.Children == nil {
				predictions.SetAttrStr(i, 0, cur.Class)
				break
			} else {
				at := cur.SplitAttr
				j := what.GetAttrIndex(at)
				classVar := at.GetStringFromSysVal(what.Get(i, j))
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
	}
	return predictions
}
