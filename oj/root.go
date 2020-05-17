// Copyright (c) 2020, Peter Ohler, All rights reserved.

package oj

// Root is used as a flag to indicate the path should be displayed in a
// recursive root representation.
type Root byte

// Append a fragment string representation of the fragment to the buffer
// then returning the expanded buffer.
func (f Root) Append(buf []byte, bracket, first bool) []byte {
	buf = append(buf, '$')
	return buf
}

func (f Root) get(top, _ interface{}, rest Expr) (results []interface{}) {
	if 0 < len(rest) {
		results = rest[0].get(top, top, rest[1:])
	} else {
		results = append(results, top)
	}
	return
}

func (f Root) first(top, _ interface{}, rest Expr) (result interface{}, found bool) {
	if 0 < len(rest) {
		result, found = rest[0].first(top, top, rest[1:])
	} else {
		result = top
		found = true
	}
	return
}
