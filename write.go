// Copyright (c) 2020, Peter Ohler, All rights reserved.

package ojg

import (
	"fmt"
	"io"
	"sort"
	"strconv"
	"time"
	"unicode/utf8"

	"github.com/ohler55/ojg/gd"
)

const (
	spaces = "\n                                                                                                                                "

	hex = "0123456789abcdef"
)

// Options for writing data to JSON.
type Options struct {

	// Indent for the output.
	Indent int

	// Sort object members if true.
	Sort bool

	// SkipNil skips the writing of nil values in an object.
	SkipNil bool

	// InitSize is the initial buffer size.
	InitSize int

	// WriteLimit is the size of the buffer that will trigger a write when
	// using a writer.
	WriteLimit int

	// TimeFormat defines how time is encoded. Options are to use a time. layout
	// string format such as time.RFC3339Nano, "second" for a decimal
	// representation, "nano" for a an integer.
	TimeFormat string

	// TimeWrap if not empty encoded time as an object with a single member. For
	// example if set to "@" then and TimeFormat is RFC3339Nano then the encoded
	// time will look like '{"@":"2020-04-12T16:34:04.123456789Z"}'
	TimeWrap string

	buf []byte
	w   io.Writer
}

// String returns a JSON string for the data provided. The data can be a
// simple type of nil, bool, int, floats, time.Time, []interface{}, or
// map[string]interface{} or a gd.Node type, The args, if supplied can be an
// int as an indent or a *Options.
func String(data interface{}, args ...interface{}) string {
	var o Options

	if 0 < len(args) {
		switch ta := args[0].(type) {
		case int:
			o.Indent = ta
		case *Options:
			o = *ta
		}
	}
	if o.InitSize == 0 {
		o.InitSize = 256
	}
	o.buf = make([]byte, 0, o.InitSize)

	_ = o.buildJSON(data, 0)

	return string(o.buf)
}

// Write a JSON string for the data provided. The data can be a simple type of
// nil, bool, int, floats, time.Time, []interface{}, or map[string]interface{}
// or a gd.Node type, The args, if supplied can be an int as an indent or a
// *Options.
func Write(w io.Writer, data interface{}, args ...interface{}) (err error) {
	var o Options

	if 0 < len(args) {
		switch ta := args[0].(type) {
		case int:
			o.Indent = ta
		case *Options:
			o = *ta
		}
	}
	o.w = w
	if o.InitSize == 0 {
		o.InitSize = 256
	}
	if o.WriteLimit == 0 {
		o.WriteLimit = 1024
	}
	o.buf = make([]byte, 0, o.InitSize)
	if err = o.buildJSON(data, 0); err != nil {
		return
	}
	if w != nil && 0 < len(o.buf) {
		_, err = o.w.Write(o.buf)
	}
	return
}

func (o *Options) buildJSON(data interface{}, depth int) (err error) {
Top:
	switch td := data.(type) {
	case nil:
		o.buf = append(o.buf, []byte("null")...)

	case bool:
		if td {
			o.buf = append(o.buf, []byte("true")...)
		} else {
			o.buf = append(o.buf, []byte("false")...)
		}
	case gd.Bool:
		if td {
			o.buf = append(o.buf, []byte("true")...)
		} else {
			o.buf = append(o.buf, []byte("false")...)
		}

	case int:
		o.buf = append(o.buf, []byte(strconv.FormatInt(int64(td), 10))...)
	case int8:
		o.buf = append(o.buf, []byte(strconv.FormatInt(int64(td), 10))...)
	case int16:
		o.buf = append(o.buf, []byte(strconv.FormatInt(int64(td), 10))...)
	case int32:
		o.buf = append(o.buf, []byte(strconv.FormatInt(int64(td), 10))...)
	case int64:
		o.buf = append(o.buf, []byte(strconv.FormatInt(td, 10))...)
	case uint:
		o.buf = append(o.buf, []byte(strconv.FormatInt(int64(td), 10))...)
	case uint8:
		o.buf = append(o.buf, []byte(strconv.FormatInt(int64(td), 10))...)
	case uint16:
		o.buf = append(o.buf, []byte(strconv.FormatInt(int64(td), 10))...)
	case uint32:
		o.buf = append(o.buf, []byte(strconv.FormatInt(int64(td), 10))...)
	case uint64:
		o.buf = append(o.buf, []byte(strconv.FormatInt(int64(td), 10))...)
	case gd.Int:
		o.buf = append(o.buf, []byte(strconv.FormatInt(int64(td), 10))...)

	case float32:
		o.buf = append(o.buf, []byte(strconv.FormatFloat(float64(td), 'g', -1, 64))...)
	case float64:
		o.buf = append(o.buf, []byte(strconv.FormatFloat(td, 'g', -1, 64))...)
	case gd.Float:
		o.buf = append(o.buf, []byte(strconv.FormatFloat(float64(td), 'g', -1, 64))...)

	case string:
		o.buildString(td)
	case gd.String:
		o.buildString(string(td))

	case time.Time:
		o.buildTime(td)
	case gd.Time:
		o.buildTime(time.Time(td))

	case []interface{}:
		err = o.buildSimpleArray(td, depth)
	case gd.Array:
		err = o.buildArray(td, depth)

	case map[string]interface{}:
		err = o.buildSimpleObject(td, depth)
	case gd.Object:
		err = o.buildObject(td, depth)

	default:
		if simp, _ := data.(gd.Simplifier); simp != nil {
			data = simp.Simplify()
			goto Top
		}
		o.buildString(fmt.Sprintf("%v", td))
	}
	if o.w != nil && o.WriteLimit < len(o.buf) {
		_, err = o.w.Write(o.buf)
		o.buf = o.buf[:0]
	}
	return
}

func (o *Options) buildString(s string) {
	o.buf = append(o.buf, '"')
	for _, r := range s {
		switch r {
		case '\\':
			o.buf = append(o.buf, []byte{'\\', '\\'}...)
		case '"':
			o.buf = append(o.buf, []byte{'\\', '"'}...)
		case '\b':
			o.buf = append(o.buf, []byte{'\\', 'b'}...)
		case '\f':
			o.buf = append(o.buf, []byte{'\\', 'f'}...)
		case '\n':
			o.buf = append(o.buf, []byte{'\\', 'n'}...)
		case '\r':
			o.buf = append(o.buf, []byte{'\\', 'r'}...)
		case '\t':
			o.buf = append(o.buf, []byte{'\\', 't'}...)
		case '&', '<', '>': // prefectly okay for JSON but commonly escaped
			o.buf = append(o.buf, []byte{'\\', 'u', '0', '0', hex[r>>4], hex[r&0x0f]}...)
		case '\u2028':
			o.buf = append(o.buf, []byte(`\u2028`)...)
		case '\u2029':
			o.buf = append(o.buf, []byte(`\u2029`)...)
		default:
			if r < ' ' {
				o.buf = append(o.buf, []byte{'\\', 'u', hex[r>>12], hex[(r>>8)&0x0f], hex[(r>>4)&0x0f], hex[r&0x0f]}...)
			} else if r < 0x80 {
				o.buf = append(o.buf, byte(r))
			} else {
				n := len(o.buf)
				need := n + utf8.UTFMax
				if cap(o.buf) < need {
					buf := make([]byte, n, n+need)
					copy(buf, o.buf)
					o.buf = buf
				}
				utf8.EncodeRune(o.buf[n:need], r)
				o.buf = o.buf[:need]
			}
		}
	}
	o.buf = append(o.buf, '"')
}

func (o *Options) buildTime(t time.Time) {
	if 0 < len(o.TimeWrap) {
		o.buf = append(o.buf, []byte(`{"`)...)
		o.buf = append(o.buf, []byte(o.TimeWrap)...)
		o.buf = append(o.buf, []byte(`":`)...)
	}
	switch o.TimeFormat {
	case "", "nano":
		o.buf = append(o.buf, []byte(strconv.FormatInt(t.UnixNano(), 10))...)
	case "second":
		// Decimal format but float is not accurate enough so build the output
		// in two parts.
		nano := t.UnixNano()
		secs := nano / int64(time.Second)
		if 0 < nano {
			o.buf = append(o.buf, []byte(fmt.Sprintf("%d.%09d", secs, nano-(secs*int64(time.Second))))...)
		} else {
			o.buf = append(o.buf, []byte(fmt.Sprintf("%d.%09d", secs, -nano-(secs*int64(time.Second))))...)
		}
	default:
		o.buf = append(o.buf, '"')
		o.buf = append(o.buf, []byte(t.Format(o.TimeFormat))...)
		o.buf = append(o.buf, '"')
	}
	if 0 < len(o.TimeWrap) {
		o.buf = append(o.buf, '}')
	}
}

func (o *Options) buildArray(n gd.Array, depth int) (err error) {
	o.buf = append(o.buf, '[')
	if 0 < o.Indent {
		x := depth*o.Indent + 1
		if len(spaces) < x {
			x = depth*o.Indent + 1
		}
		is := spaces[0:x]
		d2 := depth + 1
		x = d2*o.Indent + 1
		if len(spaces) < x {
			x = depth*o.Indent + 1
		}
		cs := spaces[0:x]

		for j, m := range n {
			if 0 < j {
				o.buf = append(o.buf, ',')
			}
			o.buf = append(o.buf, []byte(cs)...)
			if m == nil {
				o.buf = append(o.buf, []byte("null")...)
			} else if err = o.buildJSON(m, d2); err != nil {
				return
			}
		}
		o.buf = append(o.buf, []byte(is)...)
	} else {
		for j, m := range n {
			if 0 < j {
				o.buf = append(o.buf, ',')
			}
			if m == nil {
				o.buf = append(o.buf, []byte("null")...)
			} else if err = o.buildJSON(m, depth); err != nil {
				return
			}
		}
	}
	o.buf = append(o.buf, ']')

	return
}

func (o *Options) buildSimpleArray(n []interface{}, depth int) (err error) {
	o.buf = append(o.buf, '[')
	if 0 < o.Indent {
		x := depth*o.Indent + 1
		if len(spaces) < x {
			x = depth*o.Indent + 1
		}
		is := spaces[0:x]
		d2 := depth + 1
		x = d2*o.Indent + 1
		if len(spaces) < x {
			x = depth*o.Indent + 1
		}
		cs := spaces[0:x]

		for j, m := range n {
			if 0 < j {
				o.buf = append(o.buf, ',')
			}
			o.buf = append(o.buf, []byte(cs)...)
			if m == nil {
				o.buf = append(o.buf, []byte("null")...)
			} else if err = o.buildJSON(m, d2); err != nil {
				return
			}
		}
		o.buf = append(o.buf, []byte(is)...)
	} else {
		for j, m := range n {
			if 0 < j {
				o.buf = append(o.buf, ',')
			}
			if m == nil {
				o.buf = append(o.buf, []byte("null")...)
			} else if err = o.buildJSON(m, depth); err != nil {
				return
			}
		}
	}
	o.buf = append(o.buf, ']')
	return
}

func (o *Options) buildObject(n gd.Object, depth int) (err error) {
	o.buf = append(o.buf, '{')
	if 0 < o.Indent {
		x := depth*o.Indent + 1
		if len(spaces) < x {
			x = depth*o.Indent + 1
		}
		is := spaces[0:x]
		d2 := depth + 1
		x = d2*o.Indent + 1
		if len(spaces) < x {
			x = depth*o.Indent + 1
		}
		cs := spaces[0:x]
		if o.Sort {
			keys := make([]string, 0, len(n))
			for k := range n {
				keys = append(keys, k)
			}
			sort.Strings(keys)
			for i, k := range keys {
				m := n[k]
				if m == nil && o.SkipNil {
					continue
				}
				if 0 < i {
					o.buf = append(o.buf, ',')
				}
				o.buf = append(o.buf, []byte(cs)...)
				o.buildString(k)
				o.buf = append(o.buf, ':')
				if m := n[k]; m == nil {
					o.buf = append(o.buf, []byte("null")...)
				} else if err = o.buildJSON(m, d2); err != nil {
					return
				}
			}
		} else {
			first := true
			for k, m := range n {
				if m == nil && o.SkipNil {
					continue
				}
				if first {
					first = false
				} else {
					o.buf = append(o.buf, ',')
				}
				o.buf = append(o.buf, []byte(cs)...)
				o.buildString(k)
				o.buf = append(o.buf, ':')
				if m == nil {
					o.buf = append(o.buf, []byte("null")...)
				} else if err = o.buildJSON(m, d2); err != nil {
					return
				}
			}
		}
		o.buf = append(o.buf, []byte(is)...)
	} else {
		if o.Sort {
			keys := make([]string, 0, len(n))
			for k := range n {
				keys = append(keys, k)
			}
			sort.Strings(keys)
			for i, k := range keys {
				m := n[k]
				if m == nil && o.SkipNil {
					continue
				}
				if 0 < i {
					o.buf = append(o.buf, ',')
				}
				o.buildString(k)
				o.buf = append(o.buf, ':')
				if m == nil {
					o.buf = append(o.buf, []byte("null")...)
				} else if err = o.buildJSON(m, 0); err != nil {
					return
				}
			}
		} else {
			first := true
			for k, m := range n {
				if m == nil && o.SkipNil {
					continue
				}
				if first {
					first = false
				} else {
					o.buf = append(o.buf, ',')
				}
				o.buildString(k)
				o.buf = append(o.buf, ':')
				if m == nil {
					o.buf = append(o.buf, []byte("null")...)
				} else if err = o.buildJSON(m, 0); err != nil {
					return
				}
			}
		}
	}
	o.buf = append(o.buf, '}')

	return
}

func (o *Options) buildSimpleObject(n map[string]interface{}, depth int) (err error) {
	o.buf = append(o.buf, '{')
	if 0 < o.Indent {
		x := depth*o.Indent + 1
		if len(spaces) < x {
			x = depth*o.Indent + 1
		}
		is := spaces[0:x]
		d2 := depth + 1
		x = d2*o.Indent + 1
		if len(spaces) < x {
			x = depth*o.Indent + 1
		}
		cs := spaces[0:x]
		if o.Sort {
			keys := make([]string, 0, len(n))
			for k := range n {
				keys = append(keys, k)
			}
			sort.Strings(keys)
			for i, k := range keys {
				m := n[k]
				if m == nil && o.SkipNil {
					continue
				}
				if 0 < i {
					o.buf = append(o.buf, ',')
				}
				o.buf = append(o.buf, []byte(cs)...)
				o.buildString(k)
				o.buf = append(o.buf, ':')
				if m := n[k]; m == nil {
					o.buf = append(o.buf, []byte("null")...)
				} else if err = o.buildJSON(m, d2); err != nil {
					return
				}
			}
		} else {
			first := true
			for k, m := range n {
				if m == nil && o.SkipNil {
					continue
				}
				if first {
					first = false
				} else {
					o.buf = append(o.buf, ',')
				}
				o.buf = append(o.buf, []byte(cs)...)
				o.buildString(k)
				o.buf = append(o.buf, ':')
				if m == nil {
					o.buf = append(o.buf, []byte("null")...)
				} else if err = o.buildJSON(m, d2); err != nil {
					return
				}
			}
		}
		o.buf = append(o.buf, []byte(is)...)
	} else {
		if o.Sort {
			keys := make([]string, 0, len(n))
			for k := range n {
				keys = append(keys, k)
			}
			sort.Strings(keys)
			for i, k := range keys {
				m := n[k]
				if m == nil && o.SkipNil {
					continue
				}
				if 0 < i {
					o.buf = append(o.buf, ',')
				}
				o.buildString(k)
				o.buf = append(o.buf, ':')
				if m == nil {
					o.buf = append(o.buf, []byte("null")...)
				} else if err = o.buildJSON(m, 0); err != nil {
					return
				}
			}
		} else {
			first := true
			for k, m := range n {
				if m == nil && o.SkipNil {
					continue
				}
				if first {
					first = false
				} else {
					o.buf = append(o.buf, ',')
				}
				o.buildString(k)
				o.buf = append(o.buf, ':')
				if m == nil {
					o.buf = append(o.buf, []byte("null")...)
				} else if err = o.buildJSON(m, 0); err != nil {
					return
				}
			}
		}
	}
	o.buf = append(o.buf, '}')

	return
}
