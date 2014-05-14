package trees

import (
	"fmt"
	base "github.com/sjwhitworth/golearn/base"
	"strings"
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

// InferDecisionTree builds a decision tree using a RuleGenerator
// from a set of Instances
func InferDecisionTree(from *base.Instances, with RuleGenerator) *DecisionTreeNode {
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

	// If there are no more attribute left to split on,
	// return a DecisionTreeLeaf with the majority class
	if from.GetAttributeCount() == 1 {
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
		make(map[string]*DecisionTreeNode),
		nil,
		classes,
		maxClass,
		from.GetClassAttrPtr(),
	}

	// Generate a return structure
	// Generate the splitting attribute
	splitOnAttribute := with.GenerateSplitAttribute(from)
	// Split the attributes based on this attribute's value
	splitInstances := from.DecomposeOnAttributeValues(splitOnAttribute)
	// Create new children from these attributes
	for k := range splitInstances {
		newInstances := splitInstances[k]
		ret.Children[k] = InferDecisionTree(newInstances, with)
	}
	ret.SplitAttr = splitOnAttribute
	return ret
}

// String returns a human-readable representation of a given node
// and it's children
func (d *DecisionTreeNode) String() string {
	children := make([]string, 0)
	if d.Children != nil {
		for k := range d.Children {
			childStr := fmt.Sprintf("Rule(%s -> %s)", k, d.Children[k])
			children = append(children, childStr)
		}
		return fmt.Sprintf("(%s(%s))", d.SplitAttr, strings.Join(children, "\n\t"))
	}

	return fmt.Sprintf("Leaf(%s (%s))", d.Class, d.ClassDist)
}

func (d *DecisionTreeNode) Predict(what *base.Instances) *base.Instances {
	outputAttrs := make([]base.Attribute, 1)
	outputAttrs[0] = what.GetClassAttr()
	predictions := base.NewInstances(outputAttrs, what.Rows)
	for i := 0; i < what.Rows; i++ {
		cur := d
		for j := 0; j < what.Cols; j++ {
			at := what.GetAttr(j)
			classVar := at.GetStringFromSysVal(what.Get(i, j))
			if cur.Children == nil {
				predictions.SetAttrStr(i, 0, cur.Class)
			} else {
				if next, ok := cur.Children[classVar]; ok {
					cur = next
				} else {
					predictions.SetAttrStr(i, 0, cur.Class)
				}
			}
		}
	}
	return predictions
}
