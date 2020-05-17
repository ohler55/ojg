// Copyright (c) 2020, Peter Ohler, All rights reserved.

package oj

// Descent is used as a flag to indicate the path should be displayed in a
// recursive descent representation.
type Descent byte

// Append a fragment string representation of the fragment to the buffer
// then returning the expanded buffer.
func (f Descent) Append(buf []byte, bracket, first bool) []byte {
	if bracket {
		buf = append(buf, "[..]"...)
	} else {
		buf = append(buf, '.')
	}
	return buf
}

func (f Descent) get(top, data interface{}, rest Expr) (results []interface{}) {
	if 0 < len(rest) {
		stack := make([]interface{}, 0, 32)
		stack = append(stack, data)
		r := rest[0]
		r2 := rest[1:]
		var v interface{}
		for 0 < len(stack) {
			v = stack[len(stack)-1]
			stack = stack[:len(stack)-1]
			switch tv := v.(type) {
			case map[string]interface{}:
				if ra := r.get(top, v, r2); 0 < len(ra) {
					results = append(results, ra...)
				}
				for _, v2 := range tv {
					stack = append(stack, v2)
				}
			case Object:
				if ra := r.get(top, v, r2); 0 < len(ra) {
					results = append(results, ra...)
				}
				for _, v2 := range tv {
					stack = append(stack, v2)
				}
			case []interface{}:
				if ra := r.get(top, v, r2); 0 < len(ra) {
					results = append(results, ra...)
				}
				for _, v2 := range tv {
					stack = append(stack, v2)
				}
			case Array:
				if ra := r.get(top, v, r2); 0 < len(ra) {
					results = append(results, ra...)
				}
				for _, v2 := range tv {
					stack = append(stack, v2)
				}
			default:
				// TBD use reflections for map or struct
			}
		}
		// Free up anything still on the stack.
		stack = stack[0:cap(stack)]
		for i := len(stack) - 1; 0 <= i; i-- {
			stack[i] = nil
		}
	} else {
		results = append(results, data)
		switch td := data.(type) {
		case map[string]interface{}:
			for _, v := range td {
				results = append(results, v)
			}
		case Object:
			for _, v := range td {
				results = append(results, v)
			}
		case []interface{}:
			for _, v := range td {
				results = append(results, v)
			}
		case Array:
			for _, v := range td {
				results = append(results, v)
			}
		default:
			// TBD use reflections for map or struct
		}
	}
	return
}
