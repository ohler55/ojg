// Copyright (c) 2021, Peter Ohler, All rights reserved.

package ojg

import (
	"unicode/utf8"
)

var hex = "0123456789abcdef"

// AppendJSONString appends a JSON encoding of a string to the provided byte
// slice.
func AppendJSONString(buf []byte, s string, htmlSafe bool) []byte {
	buf = append(buf, '"')
	// TBD keep track of start and reset on anything that requires an append
	//  append rest at end
	for _, r := range s {
		switch r {
		case '\\':
			buf = append(buf, []byte{'\\', '\\'}...)
		case '"':
			buf = append(buf, []byte{'\\', '"'}...)
		case '\b':
			buf = append(buf, []byte{'\\', 'b'}...)
		case '\f':
			buf = append(buf, []byte{'\\', 'f'}...)
		case '\n':
			buf = append(buf, []byte{'\\', 'n'}...)
		case '\r':
			buf = append(buf, []byte{'\\', 'r'}...)
		case '\t':
			buf = append(buf, []byte{'\\', 't'}...)
		case '&', '<', '>': // prefectly okay for JSON but commonly escaped
			if htmlSafe {
				buf = append(buf, []byte{'\\', 'u', '0', '0', hex[r>>4], hex[r&0x0f]}...)
			} else {
				buf = append(buf, byte(r))
			}
		case '\u2028':
			buf = append(buf, []byte(`\u2028`)...)
		case '\u2029':
			buf = append(buf, []byte(`\u2029`)...)
		default:
			if r < ' ' {
				buf = append(buf, []byte{'\\', 'u', '0', '0', hex[(r>>4)&0x0f], hex[r&0x0f]}...)
			} else if r < 0x80 {
				buf = append(buf, byte(r))
			} else {
				utf := make([]byte, utf8.UTFMax)
				n := utf8.EncodeRune(utf, r)
				buf = append(buf, utf[:n]...)
			}
		}
	}
	return append(buf, '"')
}
