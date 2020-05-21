// Copyright (c) 2020, Peter Ohler, All rights reserved.

package oj

import "strconv"

// Slice is a slice operation for a JSON path expression.
type Slice struct {
	Start int
	End   int
	Step  int
}

// Append a fragment string representation of the fragment to the buffer
// then returning the expanded buffer.
func (f *Slice) Append(buf []byte, _, _ bool) []byte {
	buf = append(buf, '[')
	if 0 != f.Start {
		buf = append(buf, strconv.FormatInt(int64(f.Start), 10)...)
	}
	buf = append(buf, ':')
	if 0 != f.End {
		buf = append(buf, strconv.FormatInt(int64(f.End), 10)...)
	}
	if 0 < f.Step {
		buf = append(buf, ':')
		buf = append(buf, strconv.FormatInt(int64(f.Step), 10)...)
	}
	buf = append(buf, ']')

	return buf
}
