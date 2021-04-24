// Copyright (c) 2020, Peter Ohler, All rights reserved.

package alt

import (
	"strconv"
	"time"
	"unicode/utf8"

	"github.com/ohler55/ojg/gen"
)

var hex = "0123456789abcdef"

// String converts the value provided to a string. If conversion is not
// possible such as if the provided value is an array then the first option
// default value is returned or if not provided and empty string is
// returned. If the type is not a string or gen.String and there is a second
// optional default then that second default value is returned. This approach
// keeps the return as a single value and gives the caller the choice of how
// to indicate a bad value.
func String(v interface{}, defaults ...string) (s string) {
	switch ts := v.(type) {
	case string:
		s = ts
	case gen.String:
		s = string(ts)
	default:
		if 1 < len(defaults) {
			s = defaults[1]
		} else {
			switch tv := v.(type) {
			case nil:
				s = ""
			case bool:
				if tv {
					s = "true"
				} else {
					s = "false"
				}
			case int64:
				s = strconv.FormatInt(tv, 10)
			case int:
				s = strconv.FormatInt(int64(tv), 10)
			case int8:
				s = strconv.FormatInt(int64(tv), 10)
			case int16:
				s = strconv.FormatInt(int64(tv), 10)
			case int32:
				s = strconv.FormatInt(int64(tv), 10)
			case uint:
				s = strconv.FormatInt(int64(tv), 10)
			case uint8:
				s = strconv.FormatInt(int64(tv), 10)
			case uint16:
				s = strconv.FormatInt(int64(tv), 10)
			case uint32:
				s = strconv.FormatInt(int64(tv), 10)
			case uint64:
				s = strconv.FormatInt(int64(tv), 10)
			case float32:
				s = strconv.FormatFloat(float64(tv), 'g', -1, 32)
			case float64:
				s = strconv.FormatFloat(tv, 'g', -1, 64)
			case time.Time:
				s = tv.Format(time.RFC3339Nano)

			case gen.Bool:
				if tv {
					s = "true"
				} else {
					s = "false"
				}
			case gen.Int:
				s = strconv.FormatInt(int64(tv), 10)
			case gen.Float:
				s = strconv.FormatFloat(float64(tv), 'g', -1, 32)
			case gen.Time:
				s = time.Time(tv).Format(time.RFC3339Nano)
			case gen.Big:
				return string(tv)

			default:
				if 0 < len(defaults) {
					s = defaults[0]
				}
			}
		}
	}
	return
}

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
