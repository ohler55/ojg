// Copyright (c) 2020, Peter Ohler, All rights reserved.

package jp

import (
	"unsafe"
)

// Expr is a JSON path expression composed of fragments. An Expr implements
// JSONPath as described by https://goessner.net/articles/JsonPath. Where the
// definition is unclear Oj has implemented the description based on the best
// judgement of the author.
type Expr []Frag

// String returns a string representation of the expression.
func (x Expr) String() string {
	return string(x.Append(nil))
}

// BracketString returns a string representation of the expression using the
// bracket notation.
func (x Expr) BracketString() string {
	return string(x.Append(nil, true))
}

// Append a string representation of the expression to a byte slice and return
// the expanded buffer.
func (x Expr) Append(buf []byte, brackets ...bool) []byte {
	bracket := 0 < len(brackets) && brackets[0]
	for i, frag := range x {
		if _, ok := frag.(Bracket); ok {
			bracket = true
			continue
		}
		start := len(buf)
		buf = frag.Append(buf, bracket, i == 0)
		// A descent only writes one '.', so a fragment in brackets right
		// after it needs the second one to read back as a descent.
		if !bracket && 0 < i && start < len(buf) && buf[start] == '[' {
			if _, ok := x[i-1].(Descent); ok {
				buf = append(buf, 0)
				copy(buf[start+1:], buf[start:])
				buf[start] = '.'
			}
		}
	}
	if 0 < len(x) {
		if _, ok := x[len(x)-1].(Descent); ok {
			buf = append(buf, '.')
		}
	}
	return buf
}

// Normal returns true if the only fragments in the expression are root, at,
// child, and nth.
func (x Expr) Normal() bool {
	for _, f := range x {
		switch f.(type) {
		case Child, Nth, Root, At, Bracket:
			// normal
		default:
			return false
		}
	}
	return true
}

func isNil(v any) bool {
	return (*[2]uintptr)(unsafe.Pointer(&v))[1] == 0
}
