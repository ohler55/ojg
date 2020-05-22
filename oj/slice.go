// Copyright (c) 2020, Peter Ohler, All rights reserved.

package oj

import "strconv"

// Slice is a slice operation for a JSON path expression.
type Slice []int

// Append a fragment string representation of the fragment to the buffer
// then returning the expanded buffer.
func (f Slice) Append(buf []byte, _, _ bool) []byte {
	buf = append(buf, '[')
	if 0 < len(f) {
		for i, n := range f {
			if 0 < i {
				buf = append(buf, ':')
			}
			buf = append(buf, strconv.FormatInt(int64(n), 10)...)
			if 2 <= i {
				break
			}
		}
		if len(f) == 1 {
			buf = append(buf, ':')
		}
	} else {
		buf = append(buf, ':')
	}
	buf = append(buf, ']')

	return buf
}
