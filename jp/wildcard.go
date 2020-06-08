// Copyright (c) 2020, Peter Ohler, All rights reserved.

package jp

// Wildcard is used as a flag to indicate the path should be displayed in a
// wildcarded representation.
type Wildcard byte

// Append a fragment string representation of the fragment to the buffer
// then returning the expanded buffer.
func (f Wildcard) Append(buf []byte, bracket, first bool) []byte {
	if bracket || f == '#' {
		buf = append(buf, "[*]"...)
	} else {
		if !first {
			buf = append(buf, '.')
		}
		buf = append(buf, '*')
	}
	return buf
}
