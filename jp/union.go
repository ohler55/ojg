// Copyright (c) 2020, Peter Ohler, All rights reserved.

package jp

import "strconv"

// Union is a union operation for a JSON path expression which is a union of a
// Child and Nth fragment.
type Union []interface{}

// Append a fragment string representation of the fragment to the buffer
// then returning the expanded buffer.
func (f Union) Append(buf []byte, _, _ bool) []byte {
	buf = append(buf, '[')
	for i, x := range f {
		if 0 < i {
			buf = append(buf, ',')
		}
		switch tx := x.(type) {
		case string:
			buf = append(buf, '\'')
			buf = append(buf, tx...)
			buf = append(buf, '\'')
		case int64:
			buf = append(buf, strconv.FormatInt(tx, 10)...)
		}
	}
	buf = append(buf, ']')

	return buf
}

// NewUnion
func NewUnion(keys ...interface{}) (u Union) {
	for _, k := range keys {
		switch tk := k.(type) {
		case string:
			u = append(u, k)
		case int:
			u = append(u, int64(tk))
		case int64:
			u = append(u, tk)
		}
	}
	return
}
