// Copyright (c) 2020, Peter Ohler, All rights reserved.

package gen

// Genericer is the interface for the Generic() function that converts types
// to generic types.
type Genericer interface {

	// Generic should return a Node that represents the object. Generally this
	// includes the use of a creation key consistent with call to the
	// reflection based Generic() function.
	Generic() Node
}
