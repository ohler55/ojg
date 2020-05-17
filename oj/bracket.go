// Copyright (c) 2020, Peter Ohler, All rights reserved.

package oj

// Bracket is used as a flag to indicate the path should be displayed in a
// bracketed representation.
type Bracket byte

// Append a fragment string representation of the fragment to the buffer
// then returning the expanded buffer.
func (f Bracket) Append(buf []byte, bracket, first bool) []byte {
	return buf
}

func (f Bracket) get(top, data interface{}, rest Expr) (results []interface{}) {
	if 0 < len(rest) {
		results = rest[0].get(top, data, rest[1:])
	} else {
		results = append(results, data)
	}
	return
}
