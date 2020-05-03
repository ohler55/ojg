// Copyright (c) 2020, Peter Ohler, All rights reserved.

package ojg

import (
	"github.com/ohler55/ojg/gd"
)

type nKey string

// String returns the key as a string.
func (k nKey) String() string {
	return string(k)
}

// Alter converts the node into it's native type. Note this will modify
// Objects and Arrays in place making them no longer usable as the
// original type. Use with care!
func (k nKey) Alter() interface{} {
	return string(k)
}

// Simplify makes a copy of the node but as simple types.
func (k nKey) Simplify() interface{} {
	return string(k)
}

// Dup returns a deep duplicate of the node.
func (k nKey) Dup() gd.Node {
	return k
}

// Empty returns true if the node is empty.
func (k nKey) Empty() bool {
	return false
}
