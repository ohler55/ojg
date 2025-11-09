// Copyright (c) 2025, Peter Ohler, All rights reserved.

package discover

import (
	"io"

	"github.com/ohler55/ojg/oj"
)

// JSONbytes finds potential occurrence of JSON documents that are either maps
// or arrays. This is a best effort search to find potential JSON documents. It
// is possible that document will not parse without errors. The callback
// function should return a true back return value to back up to the next open
// character after the current start. If back is false scanning continues
// after the end of the found section. If stop is true then no further
// scanning is attempted and the function returns.
func JSONbytes(buf []byte, cb func(found []byte) (back, stop bool)) {
	jsonBytes(buf, cb, nil)
}

func jsonBytes(
	buf []byte,
	cb func(found []byte) (back, stop bool),
	more func(buf []byte, start, i int) ([]byte, int, int, bool)) {

	// TBD
}

// JSON finds occurrence of JSON documents that are either maps or arrays. The
// callback function should return true to stop discovering.
func JSON(buf []byte, cb func(value any) (stop bool)) {
	JSONbytes(buf, func(found []byte) (bool, bool) {
		if value, err := oj.Parse(found); err == nil {
			return false, cb(value)
		}
		return true, false
	})
}

// ReadJSONbytes finds potential occurrence of JSON documents that are either
// maps or arrays in a stream. This is a best effort search to find potential
// JSON documents. It is possible that document will not parse without
// errors. The callback function should return a true back return value to
// back up to the next open character after the current start. If back is
// false scanning continues after the end of the found section. If stop is
// true then no further scanning is attempted and the function returns.
func ReadJSONbytes(r io.Reader, cb func(b []byte) (back, stop bool)) {
	jsonBytes(nil, cb, func(buf []byte, start, i int) ([]byte, int, int, bool) {
		return readMore(r, buf, start, i)
	})
}

// ReadJSON finds occurrence of JSON documents that are either maps or arrays in
// a stream. The callback function should return true to stop discovering.
func ReadJSON(r io.Reader, cb func(value any) bool) {
	ReadJSONbytes(r, func(found []byte) (bool, bool) {
		if value, err := oj.Parse(found); err == nil {
			return false, cb(value)
		}
		return true, false
	})
}
