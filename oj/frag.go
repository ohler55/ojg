// Copyright (c) 2020, Peter Ohler, All rights reserved.

package oj

// Frag represents a JSONPath fragment. A JSONPath expression is composed of
// fragments (Frag) linked together to form a full path expression.
type Frag interface {

	// Append a fragment string representation of the fragment to the buffer
	// then returning the expanded buffer.
	Append(buf []byte, bracket, first bool) []byte

	get(top, data interface{}, rest Expr) []interface{}
	first(top, data interface{}, rest Expr) (result interface{}, found bool)
}
