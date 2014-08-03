package base

import (
    "fmt"
)

// AttributeSpec is a pointer to a particular Attribute
// within a particular Instance structure and encodes position
// and storage information associated with that Attribute.
type AttributeSpec struct {
    pondName string
    position int
    attr     Attribute
}

// GetAttribute returns an AttributeSpec which matches a given
// Attribute.
func (a *AttributeSpec) GetAttribute() Attribute {
    return a.attr
}

// String returns a human-readable description of this AttributeSpec.
func (a *AttributeSpec) String() string {
    return fmt.Sprintf("AttributeSpec(Attribute: '%s', Pond: '%s'/%d)", a.attr, a.pondName, a.position)
}
