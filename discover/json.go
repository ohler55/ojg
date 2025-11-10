// Copyright (c) 2025, Peter Ohler, All rights reserved.

package discover

import (
	"io"

	"github.com/ohler55/ojg/oj"
)

// JSON finds occurrence of JSON documents that are either maps or arrays. The
// callback function should return true to stop discovering.
func JSON(buf []byte, cb func(value any) (stop bool)) {
	Find(buf, func(found []byte) (bool, bool) {
		if value, err := oj.Parse(found); err == nil {
			return false, cb(value)
		}
		return true, false
	})
}

// ReadJSON finds occurrence of JSON documents that are either maps or arrays in
// a stream. The callback function should return true to stop discovering.
func ReadJSON(r io.Reader, cb func(value any) bool) {
	Read(r, func(found []byte) (bool, bool) {
		if value, err := oj.Parse(found); err == nil {
			return false, cb(value)
		}
		return true, false
	})
}
