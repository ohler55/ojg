// Copyright (c) 2020, Peter Ohler, All rights reserved.

package oj

import "strconv"

// Union is a union operation for a JSON path expression which is a union of a
// Child and Nth fragment.
type Union struct {
	Indexes []int
	Keys    []string
}

// Append a fragment string representation of the fragment to the buffer
// then returning the expanded buffer.
func (f *Union) Append(buf []byte, _, _ bool) []byte {
	buf = append(buf, '[')
	comma := false
	for _, s := range f.Keys {
		if comma {
			buf = append(buf, ',')
		} else {
			comma = true
		}
		buf = append(buf, '\'')
		buf = append(buf, s...)
		buf = append(buf, '\'')
	}
	for _, i := range f.Indexes {
		if comma {
			buf = append(buf, ',')
		} else {
			comma = true
		}
		buf = append(buf, strconv.FormatInt(int64(i), 10)...)
	}
	buf = append(buf, ']')

	return buf
}
