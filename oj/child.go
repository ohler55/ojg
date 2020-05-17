// Copyright (c) 2020, Peter Ohler, All rights reserved.

package oj

// Child is a child operation for a JSON path expression.
type Child string

// Append a fragment string representation of the fragment to the buffer
// then returning the expanded buffer.
func (f Child) Append(buf []byte, bracket, first bool) []byte {
	if bracket {
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

func (f Child) get(top, data interface{}, rest Expr) (results []interface{}) {
	switch td := data.(type) {
	case map[string]interface{}:
		if v, has := td[string(f)]; has {
			if 0 < len(rest) {
				results = rest[0].get(top, v, rest[1:])
			} else {
				results = append(results, v)
			}
		}
	case Object:
		if v, has := td[string(f)]; has {
			if 0 < len(rest) {
				results = rest[0].get(top, v, rest[1:])
			} else {
				results = append(results, v)
			}
		}
	default:
		// TBD use reflections for map or struct
	}
	return
}
