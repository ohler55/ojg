// Copyright (c) 2020, Peter Ohler, All rights reserved.

package oj

// Wildcard is used as a flag to indicate the path should be displayed in a
// wildcarded representation.
type Wildcard byte

// Append a fragment string representation of the fragment to the buffer
// then returning the expanded buffer.
func (f Wildcard) Append(buf []byte, bracket, first bool) []byte {
	if bracket {
		buf = append(buf, "[*]"...)
	} else {
		if !first {
			buf = append(buf, '.')
		}
		buf = append(buf, '*')
	}
	return buf
}

func (f Wildcard) get(top, data interface{}, rest Expr) (results []interface{}) {
	switch td := data.(type) {
	case map[string]interface{}:
		for _, v := range td {
			if 0 < len(rest) {
				results = append(results, rest[0].get(top, v, rest[1:])...)
			} else {
				results = append(results, v)
			}
		}
	case Object:
		for _, v := range td {
			if 0 < len(rest) {
				results = append(results, rest[0].get(top, v, rest[1:])...)
			} else {
				results = append(results, v)
			}
		}
	case []interface{}:
		for _, v := range td {
			if 0 < len(rest) {
				results = append(results, rest[0].get(top, v, rest[1:])...)
			} else {
				results = append(results, v)
			}
		}
	case Array:
		for _, v := range td {
			if 0 < len(rest) {
				results = append(results, rest[0].get(top, v, rest[1:])...)
			} else {
				results = append(results, v)
			}
		}
	default:
		// TBD use reflections for map or struct
	}
	return
}

func (f Wildcard) first(top, data interface{}, rest Expr) (result interface{}, found bool) {
	switch td := data.(type) {
	case map[string]interface{}:
		for _, v := range td {
			if 0 < len(rest) {
				if result, found = rest[0].first(top, v, rest[1:]); found {
					break
				}
			} else {
				result = v
				found = true
				break
			}
		}
	case Object:
		for _, v := range td {
			if 0 < len(rest) {
				if result, found = rest[0].first(top, v, rest[1:]); found {
					break
				}
			} else {
				result = v
				found = true
				break
			}
		}
	case []interface{}:
		for _, v := range td {
			if 0 < len(rest) {
				if result, found = rest[0].first(top, v, rest[1:]); found {
					break
				}
			} else {
				result = v
				found = true
				break
			}
		}
	case Array:
		for _, v := range td {
			if 0 < len(rest) {
				if result, found = rest[0].first(top, v, rest[1:]); found {
					break
				}
			} else {
				result = v
				found = true
				break
			}
		}
	default:
		// TBD use reflections for map or struct
	}
	return
}
