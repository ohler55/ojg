// Copyright (c) 2020, Peter Ohler, All rights reserved.

package oj

// Expr is a JSON path expression composed of fragments.
type Expr []Frag

// String returns a string representation of the expression.
func (x Expr) String() string {
	return string(x.Append(nil))
}

// Append a string representation of the expression to a byte slice and return
// the expanded buffer.
func (x Expr) Append(buf []byte) []byte {
	bracket := false
	for i, frag := range x {
		if _, ok := frag.(Bracket); ok {
			bracket = true
			continue
		}
		buf = frag.Append(buf, bracket, i == 0)
	}
	return buf
}
