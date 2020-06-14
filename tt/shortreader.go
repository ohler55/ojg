// Copyright (c) 2020, Peter Ohler, All rights reserved.

package tt

import "fmt"

// ShortReader readons only the designated amount and then returns an
// error.
type ShortReader struct {
	Max     int
	Content []byte
	pos     int
}

// Read the next batch of bytes.
func (r *ShortReader) Read(p []byte) (n int, err error) {
	start := r.pos
	r.pos += len(p)
	if r.Max < r.pos {
		return 0, fmt.Errorf("fail now")
	}
	return copy(p, r.Content[start:]), nil
}
