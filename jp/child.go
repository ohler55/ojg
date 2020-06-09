// Copyright (c) 2020, Peter Ohler, All rights reserved.

package jp

// Child is a child operation for a JSON path expression.
type Child string

// Append a fragment string representation of the fragment to the buffer
// then returning the expanded buffer.
func (f Child) Append(buf []byte, bracket, first bool) []byte {
	if bracket || !f.tokenOk() {
		buf = append(buf, "['"...)
		buf = append(buf, string(f)...)
		buf = append(buf, "']"...)
	} else {
		if !first {
			buf = append(buf, '.')
		}
		buf = append(buf, string(f)...)
	}
	return buf
}

func (f Child) tokenOk() bool {
	for _, b := range []byte(f) {
		if tokenMap[b] == '.' {
			return false
		}
	}
	return true
}
