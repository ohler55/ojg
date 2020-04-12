// Copyright (c) 2020, Peter Ohler, All rights reserved.

package gd

import "fmt"

// Path represents a JSONPath.
type Path interface {
	fmt.Stringer

	// Get the next set of matching elements.
	Get(at Node) []Node

	// Set a child node value.
	Set(at, value Node) error

	// Remove removes nodes returns then in an array.
	Remove(at Node) []Node
}
