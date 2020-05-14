// Copyright (c) 2020, Peter Ohler, All rights reserved.

package oj

import (
	"fmt"
	"io"
	"sort"
	"strconv"
	"time"
	"unicode/utf8"
)

const (
	spaces = "\n                                                                                                                                "

	hex = "0123456789abcdef"
)

// JSON returns a JSON string for the data provided. The data can be a
// simple type of nil, bool, int, floats, time.Time, []interface{}, or
// map[string]interface{} or a Node type, The args, if supplied can be an
// int as an indent or a *Options.
func JSON(data interface{}, args ...interface{}) string {
	o := &defaultOptions

	if 0 < len(args) {
		switch ta := args[0].(type) {
		case int:
			oi := *o
			oi.Indent = ta
			o = &oi
		case *Options:
			o = ta
		}
	}
	if o.InitSize == 0 {
		o.InitSize = 256
	}
	if cap(o.buf) < o.InitSize {
		o.buf = make([]byte, 0, o.InitSize)
	} else {
		o.buf = o.buf[:0]
	}
	_ = o.buildJSON(data, 0)

	return string(o.buf)
}

// Write a JSON string for the data provided. The data can be a simple type of
// nil, bool, int, floats, time.Time, []interface{}, or map[string]interface{}
// or a Node type, The args, if supplied can be an int as an indent or a
// *Options.
func Write(w io.Writer, data interface{}, args ...interface{}) (err error) {
	o := &defaultOptions

	if 0 < len(args) {
		switch ta := args[0].(type) {
		case int:
			oi := *o
			oi.Indent = ta
			o = &oi
		case *Options:
			o = ta
		}
	}
	o.w = w
	if o.InitSize == 0 {
		o.InitSize = 256
	}
	if o.WriteLimit == 0 {
		o.WriteLimit = 1024
	}
	if cap(o.buf) < o.InitSize {
		o.buf = make([]byte, 0, o.InitSize)
	} else {
		o.buf = o.buf[:0]
	}
	if err = o.buildJSON(data, 0); err != nil {
		return
	}
	if w != nil && 0 < len(o.buf) {
		_, err = o.w.Write(o.buf)
	}
	return
}

func (o *Options) buildJSON(data interface{}, depth int) (err error) {
	switch td := data.(type) {
	case nil:
		o.buf = append(o.buf, []byte("null")...)

	case bool:
		if td {
			o.buf = append(o.buf, []byte("true")...)
		} else {
			o.buf = append(o.buf, []byte("false")...)
		}
	case Bool:
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
	case Int:
		o.buf = append(o.buf, []byte(strconv.FormatInt(int64(td), 10))...)

	case float32:
		o.buf = append(o.buf, []byte(strconv.FormatFloat(float64(td), 'g', -1, 32))...)
	case float64:
		o.buf = append(o.buf, []byte(strconv.FormatFloat(td, 'g', -1, 64))...)
	case Float:
		o.buf = append(o.buf, []byte(strconv.FormatFloat(float64(td), 'g', -1, 64))...)

	case string:
		o.buildString(td)
	case String:
		o.buildString(string(td))

	case time.Time:
		o.buildTime(td)
	case Time:
		o.buildTime(time.Time(td))

	case []interface{}:
		err = o.buildSimpleArray(td, depth)
	case Array:
		err = o.buildArray(td, depth)

	case map[string]interface{}:
		err = o.buildSimpleObject(td, depth)
	case Object:
		err = o.buildObject(td, depth)

	default:
		if g, _ := data.(Genericer); g != nil {
			return o.buildJSON(g.Generic(), depth)
		}
		if simp, _ := data.(Simplifier); simp != nil {
			data = simp.Simplify()
			return o.buildJSON(data, depth)
		}
		if 0 < len(o.CreateKey) {
			return o.buildJSON(Decompose(data), depth)
		} else {
			o.buildString(fmt.Sprintf("%v", td))
		}
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

func (o *Options) buildArray(n Array, depth int) (err error) {
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

func (o *Options) buildObject(n Object, depth int) (err error) {
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
				if m == nil && o.OmitNil {
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
				if m == nil && o.OmitNil {
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
				if m == nil && o.OmitNil {
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
				if m == nil && o.OmitNil {
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
				if m == nil && o.OmitNil {
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
				if m == nil && o.OmitNil {
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
				if m == nil && o.OmitNil {
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
				if m == nil && o.OmitNil {
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
