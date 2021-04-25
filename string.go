// Copyright (c) 2021, Peter Ohler, All rights reserved.

package ojg

import (
	"unicode/utf8"
)

const (
	hex = "0123456789abcdef"

	tokenStart = 'j'
	tokenOk    = 'u'
)

var (
	maxTokenLen = 64

	// Copied from sen/maps.go

	//   0123456789abcdef0123456789abcdef
	valueMap = "" +
		".........ab..a.................." + // 0x00
		"a.i.j.....jjafjcghhhhhhhhh..j.jj" + // 0x20
		"jjjjjjjjjjjjjjjjjjjjjjjjjjjk.mjj" + // 0x40
		".jjjjjjjjjjjjjjjjjjjjjjjjjjl.nj." + // 0x60
		"jjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjj" + // 0x80
		"jjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjj" + // 0xa0
		"jjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjj" + // 0xc0
		"jjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjv" //  0xe0
	//   0123456789abcdef0123456789abcdef
	tokenMap = "" +
		".........GJ..G.................." + // 0x00
		"G...u.....uuGuucuuuuuuuuuuI.u.uu" + // 0x20
		"uuuuuuuuuuuuuuuuuuuuuuuuuuuk.muu" + // 0x40
		".uuuuuuuuuuuuuuuuuuuuuuuuuul.nu." + // 0x60
		"uuuuuuuuuuuuuuuuuuuuuuuuuuuuuuuu" + // 0x80
		"uuuuuuuuuuuuuuuuuuuuuuuuuuuuuuuu" + // 0xa0
		"uuuuuuuuuuuuuuuuuuuuuuuuuuuuuuuu" + // 0xc0
		"uuuuuuuuuuuuuuuuuuuuuuuuuuuuuuuut" //  0xe0
)

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

func AppendSENString(buf []byte, s string) []byte {
	tokOk := false
	if 0 < len(s) {
		vm := valueMap
		tm := tokenMap
		if vm[s[0]] == tokenStart &&
			len(s) < maxTokenLen { // arbitrary length, longer strings look better in quotes
			tokOk = true
			for _, b := range []byte(s) {
				if tm[b] != tokenOk {
					tokOk = false
					break
				}
			}
		}
	}
	if !tokOk {
		buf = append(buf, '"')
	}
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
			buf = append(buf, []byte{'\n'}...)
		case '\r':
			buf = append(buf, []byte{'\\', 'r'}...)
		case '\t':
			buf = append(buf, []byte{'\t'}...)
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
	if !tokOk {
		buf = append(buf, '"')
	}
	return buf
}
