// Copyright (c) 2020, Peter Ohler, All rights reserved.

package jp

// Frag represents a JSONPath fragment. A JSONPath expression is composed of
// fragments (Frag) linked together to form a full path expression.
type Frag interface {

	// Append a fragment string representation of the fragment to the buffer
	// then returning the expanded buffer.
	Append(buf []byte, bracket, first bool) []byte

	locate(pp Expr, data any, rest Expr, max int) (locs []Expr)

	// Walk the matching elements in the data and call cb on the matches or
	// follow on to the matching if not the last fragment in an
	// expression. The rest argument is the rest of the expression after this
	// fragment. The path is the normalized path up to this point. Data is the
	// data element to act on.
	Walk(rest, path Expr, data any, cb func(path Expr, data, parent any))
}
