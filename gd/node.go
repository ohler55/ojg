// Copyright (c) 2020, Peter Ohler, All rights reserved.

package gd

import (
	"fmt"
	"strings"
)

type Node interface {
	fmt.Stringer

	// Alter converts the node into it's native type. Note this will modify
	// Objects and Arrays in place making them no longer usable as the
	// original type. Use with care!
	Alter() interface{}

	// Native makes a copy of the node but as native types.
	Native() interface{}

	// Dup returns a deep duplicate of the node.
	Dup() Node

	// Empty returns true if the node is empty.
	Empty() bool

	// AsBool returns the Bool value of the node if possible. The ok return is
	// true if successful.
	AsBool() (v Bool, ok bool)

	// AsInt returns the Int value of the node if possible. The ok return is
	// true if successful.
	AsInt() (v Int, ok bool)

	// AsFloat returns the Float value of the node if possible. The ok return
	// is true if successful.
	AsFloat() (v Float, ok bool)

	// AsTime() (v Time, ok bool)
	// AsString() (v String, ok bool)
	// AsArray() (v Array, ok bool)
	// AsObject() (v Object, ok bool)

	// JSON returns a JSON encoded version of the instance.
	JSON(indent ...int) string

	// BuildJSON appends a JSON encoded version of the instance to the
	// provided builder.
	BuildJSON(b *strings.Builder, indent, depth int)
}
