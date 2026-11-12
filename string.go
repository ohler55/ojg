// Copyright (c) 2021, Peter Ohler, All rights reserved.

package ojg

import (
	"fmt"
	"reflect"
	"unicode/utf8"
)

const hex = "0123456789abcdef"

var (
	maxTokenLen = 64

	// Copied from sen/maps.go

	//   0123456789abcdef0123456789abcdef
	jMap = "" +
		`........btn.fr..................` + // 0x00
		`oo"ooohooooooooooooooooooooohoho` + // 0x20
		`oooooooooooooooooooooooooooo\ooo` + // 0x40
		"ooooooooooooooooooooooooooooooo." + // 0x60
		`88888888888888888888888888888888` + // 0x80
		`88888888888888888888888888888888` + // 0xa0
		`88888888888888888888888888888888` + // 0xc0
		`88888888888888888888888888888888` //  0xe0

	//   0123456789abcdef0123456789abcdef
	senMap = "" +
		`........bxx.fr..................` + // 0x00
		`xx"xoxhxxxooxoox0000000000xxhxho` + // 0x20
		`ooooooooooooooooooooooooooox\xoo` + // 0x40
		"xooooooooooooooooooooooooooxoxo." + // 0x60
		`88888888888888888888888888888888` + // 0x80
		`88888888888888888888888888888888` + // 0xa0
		`88888888888888888888888888888888` + // 0xc0
		`88888888888888888888888888888888` //  0xe0
)

// AppendJSONString appends a JSON encoding of a string to the provided byte
// slice.
func AppendJSONString(buf []byte, s string, htmlSafe bool) []byte {
	buf = append(buf, '"')
	start := 0
	skip := 0
	for i, b := range []byte(s) {
		if i < skip {
			continue
		}
		c := jMap[b]
		switch c {
		case 'o':
			continue
		case '.':
			if start < i {
				buf = append(buf, s[start:i]...)
			}
			buf = append(buf, `\u00`...)
			buf = append(buf, hex[(b>>4)&0x0f])
			buf = append(buf, hex[b&0x0f])
			start = i + 1
		case 'h':
			if htmlSafe {
				if start < i {
					buf = append(buf, s[start:i]...)
				}
				buf = append(buf, `\u00`...)
				buf = append(buf, hex[(b>>4)&0x0f])
				buf = append(buf, hex[b&0x0f])
				start = i + 1
			}
		case '8':
			r, cnt := utf8.DecodeRuneInString(s[i:])
			switch r {
			case '\u2028':
				if start < i {
					buf = append(buf, s[start:i]...)
				}
				buf = append(buf, `\u2028`...)
				start = i + cnt
				skip = start
			case '\u2029':
				if start < i {
					buf = append(buf, s[start:i]...)
				}
				buf = append(buf, `\u2029`...)
				start = i + cnt
				skip = start
			case utf8.RuneError:
				if start < i {
					buf = append(buf, s[start:i]...)
				}
				buf = append(buf, `\ufffd`...)
				start = i + cnt
				skip = start
			default:
				skip = i + cnt
			}
		default:
			if start < i {
				buf = append(buf, s[start:i]...)
			}
			buf = append(buf, '\\')
			buf = append(buf, c)
			start = i + 1
		}
	}
	if start < len(s) {
		buf = append(buf, s[start:]...)
	}
	return append(buf, '"')
}

// AppendSENString appends a SEN encoding of a string to the provided byte
// slice.
func AppendSENString(buf []byte, s string, htmlSafe bool) []byte {
	return appendSENString(buf, s, htmlSafe, false)
}

// AppendSENKey appends a SEN encoding of a map key to the provided byte
// slice. A key is read back as a string even when it is written bare so the
// literals are left bare.
func AppendSENKey(buf []byte, s string, htmlSafe bool) []byte {
	return appendSENString(buf, s, htmlSafe, true)
}

func appendSENString(buf []byte, s string, htmlSafe bool, key bool) []byte {
	if len(s) == 0 {
		return append(buf, `""`...)
	}
	b0 := len(buf)
	m := senMap[s[0]]
	quote := maxTokenLen < len(s) || (m != 'o' && m != '8' && !(!htmlSafe && m == 'h'))
	if 1 < len(s) && (s[0] == '+' || s[0] == '-') {
		// A lone sign stays bare as an operator. More than that reads as a
		// number or a concatenation.
		quote = true
	}
	if !key && (s == "null" || s == "true" || s == "false") {
		// These are the literals and are read back as the literal and not as
		// the string. A key is read back as a string even when it is bare.
		quote = true
	}
	if !key && s[0] == 0xEF {
		// A leading byte order mark is taken for a mark by the tokenizer and
		// dropped, and any other character with the same lead byte is taken
		// for a broken mark. Only the first element of a document can be at
		// the start of it so a key is never affected.
		quote = true
	}
	buf = append(buf, '"')
	start := 0
	skip := 0
	for i, b := range []byte(s) {
		if i < skip {
			continue
		}
		c := senMap[b]
		switch c {
		case 'o', '0':
			continue
		case 'x':
			quote = true
		case '.':
			quote = true
			if start < i {
				buf = append(buf, s[start:i]...)
			}
			buf = append(buf, `\u00`...)
			buf = append(buf, hex[(b>>4)&0x0f])
			buf = append(buf, hex[b&0x0f])
			start = i + 1
		case 'h':
			if htmlSafe {
				quote = true
				if start < i {
					buf = append(buf, s[start:i]...)
				}
				buf = append(buf, `\u00`...)
				buf = append(buf, hex[(b>>4)&0x0f])
				buf = append(buf, hex[b&0x0f])
				start = i + 1
			}
		case '8':
			r, cnt := utf8.DecodeRuneInString(s[i:])
			switch r {
			case '\u2028':
				quote = true
				if start < i {
					buf = append(buf, s[start:i]...)
				}
				buf = append(buf, `\u2028`...)
				start = i + cnt
				skip = start
			case '\u2029':
				quote = true
				if start < i {
					buf = append(buf, s[start:i]...)
				}
				buf = append(buf, `\u2029`...)
				start = i + cnt
				skip = start
			case utf8.RuneError:
				quote = true
				if start < i {
					buf = append(buf, s[start:i]...)
				}
				buf = append(buf, `\ufffd`...)
				start = i + cnt
				skip = start
			default:
				skip = i + cnt
			}
		default:
			if start < i {
				buf = append(buf, s[start:i]...)
			}
			buf = append(buf, '\\')
			buf = append(buf, c)
			start = i + 1
			quote = true
		}
	}
	if start < len(s) {
		buf = append(buf, s[start:]...)
	}
	if quote {
		return append(buf, '"')
	}
	copy(buf[b0:], buf[b0+1:])

	return buf[:len(buf)-1]
}

// KeyString returns the string form of a reflect map key. Keys that are not
// strings are formatted with the fmt package the same way alt.reflectMap
// formats them.
func KeyString(kv reflect.Value) string {
	if kv.Kind() == reflect.String {
		return kv.String()
	}
	return fmt.Sprint(kv.Interface())
}
