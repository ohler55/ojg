// Copyright (c) 2020, Peter Ohler, All rights reserved.

package sen

import (
	"fmt"
	"io"
	"sort"
	"strconv"
	"time"
	"unicode/utf8"

	"github.com/ohler55/ojg/alt"
	"github.com/ohler55/ojg/gen"
)

const (
	spaces      = "\n                                                                                                                                "
	tabs        = "\n\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t"
	hex         = "0123456789abcdef"
	maxTokenLen = 64
)

// String returns a SEN string for the data provided. The data can be a simple
// type of nil, bool, int, floats, time.Time, []interface{}, or
// map[string]interface{} or a Node type, The args, if supplied can be an int
// as an indent or a *Options.
func String(data interface{}, args ...interface{}) string {
	o := &DefaultOptions

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
	if cap(o.Buf) < o.InitSize {
		o.Buf = make([]byte, 0, o.InitSize)
	} else {
		o.Buf = o.Buf[:0]
	}
	defer func() {
		if r := recover(); r != nil {
			o.Buf = o.Buf[:0]
		}
	}()
	o.buildSen(data, 0)

	return string(o.Buf)
}

// Write a JSON string for the data provided. The data can be a simple type of
// nil, bool, int, floats, time.Time, []interface{}, or map[string]interface{}
// or a Node type, The args, if supplied can be an int as an indent or a
// *Options.
func Write(w io.Writer, data interface{}, args ...interface{}) (err error) {
	o := &DefaultOptions

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
	o.W = w
	if o.InitSize == 0 {
		o.InitSize = 256
	}
	if o.WriteLimit == 0 {
		o.WriteLimit = 1024
	}
	if cap(o.Buf) < o.InitSize {
		o.Buf = make([]byte, 0, o.InitSize)
	} else {
		o.Buf = o.Buf[:0]
	}
	defer func() {
		if r := recover(); r != nil {
			o.Buf = o.Buf[:0]
			if err, _ = r.(error); err == nil {
				err = fmt.Errorf("%v", r)
			}
		}
	}()
	if o.Color {
		o.cbuildJSON(data, 0)
	} else {
		o.buildSen(data, 0)
	}
	if o.Color {
		o.Buf = append(o.Buf, Normal...)
	}
	if w != nil && 0 < len(o.Buf) {
		_, err = o.W.Write(o.Buf)
	}
	return
}

func (o *Options) buildSen(data interface{}, depth int) {
	switch td := data.(type) {
	case nil:
		o.Buf = append(o.Buf, []byte("null")...)

	case bool:
		if td {
			o.Buf = append(o.Buf, []byte("true")...)
		} else {
			o.Buf = append(o.Buf, []byte("false")...)
		}
	case gen.Bool:
		if td {
			o.Buf = append(o.Buf, []byte("true")...)
		} else {
			o.Buf = append(o.Buf, []byte("false")...)
		}

	case int:
		o.Buf = append(o.Buf, []byte(strconv.FormatInt(int64(td), 10))...)
	case int8:
		o.Buf = append(o.Buf, []byte(strconv.FormatInt(int64(td), 10))...)
	case int16:
		o.Buf = append(o.Buf, []byte(strconv.FormatInt(int64(td), 10))...)
	case int32:
		o.Buf = append(o.Buf, []byte(strconv.FormatInt(int64(td), 10))...)
	case int64:
		o.Buf = append(o.Buf, []byte(strconv.FormatInt(td, 10))...)
	case uint:
		o.Buf = append(o.Buf, []byte(strconv.FormatInt(int64(td), 10))...)
	case uint8:
		o.Buf = append(o.Buf, []byte(strconv.FormatInt(int64(td), 10))...)
	case uint16:
		o.Buf = append(o.Buf, []byte(strconv.FormatInt(int64(td), 10))...)
	case uint32:
		o.Buf = append(o.Buf, []byte(strconv.FormatInt(int64(td), 10))...)
	case uint64:
		o.Buf = append(o.Buf, []byte(strconv.FormatInt(int64(td), 10))...)
	case gen.Int:
		o.Buf = append(o.Buf, []byte(strconv.FormatInt(int64(td), 10))...)

	case float32:
		o.Buf = append(o.Buf, []byte(strconv.FormatFloat(float64(td), 'g', -1, 32))...)
	case float64:
		o.Buf = append(o.Buf, []byte(strconv.FormatFloat(td, 'g', -1, 64))...)
	case gen.Float:
		o.Buf = append(o.Buf, []byte(strconv.FormatFloat(float64(td), 'g', -1, 64))...)

	case string:
		o.BuildString(td)
	case gen.String:
		o.BuildString(string(td))

	case time.Time:
		o.BuildTime(td)
	case gen.Time:
		o.BuildTime(time.Time(td))

	case []interface{}:
		o.buildSimpleArray(td, depth)
	case gen.Array:
		o.buildArray(td, depth)

	case map[string]interface{}:
		o.buildSimpleObject(td, depth)
	case gen.Object:
		o.buildObject(td, depth)

	default:
		if g, _ := data.(alt.Genericer); g != nil {
			o.buildSen(g.Generic(), depth)
			return
		}
		if simp, _ := data.(alt.Simplifier); simp != nil {
			data = simp.Simplify()
			o.buildSen(data, depth)
			return
		}
		if 0 < len(o.CreateKey) {
			ao := alt.Options{CreateKey: o.CreateKey, OmitNil: o.OmitNil, FullTypePath: o.FullTypePath}
			o.buildSen(alt.Decompose(data, &ao), depth)
			return
		} else {
			o.BuildString(fmt.Sprintf("%v", td))
		}
	}
	if o.W != nil && o.WriteLimit < len(o.Buf) {
		if _, err := o.W.Write(o.Buf); err != nil {
			panic(err)
		}
		o.Buf = o.Buf[:0]
	}
}

func (o *Options) BuildString(s string) {
	tokOk := false
	if 0 < len(s) &&
		valueMap[s[0]] == tokenStart &&
		len(s) < maxTokenLen { // arbitrary length, longer strings look better in quotes
		tokOk = true
		for _, b := range []byte(s) {
			if tokenMap[b] != tokenOk {
				tokOk = false
				break
			}
		}
	}
	if !tokOk {
		o.Buf = append(o.Buf, '"')
	}
	for _, r := range s {
		switch r {
		case '\\':
			o.Buf = append(o.Buf, []byte{'\\', '\\'}...)
		case '"':
			o.Buf = append(o.Buf, []byte{'\\', '"'}...)
		case '\b':
			o.Buf = append(o.Buf, []byte{'\\', 'b'}...)
		case '\f':
			o.Buf = append(o.Buf, []byte{'\\', 'f'}...)
		case '\n':
			o.Buf = append(o.Buf, []byte{'\n'}...)
		case '\r':
			o.Buf = append(o.Buf, []byte{'\\', 'r'}...)
		case '\t':
			o.Buf = append(o.Buf, []byte{'\t'}...)
		case '\u2028':
			o.Buf = append(o.Buf, []byte(`\u2028`)...)
		case '\u2029':
			o.Buf = append(o.Buf, []byte(`\u2029`)...)
		default:
			if r < ' ' {
				o.Buf = append(o.Buf, []byte{'\\', 'u', '0', '0', hex[(r>>4)&0x0f], hex[r&0x0f]}...)
			} else if r < 0x80 {
				o.Buf = append(o.Buf, byte(r))
			} else {
				if len(o.Utf) < utf8.UTFMax {
					o.Utf = make([]byte, utf8.UTFMax)
				}
				n := utf8.EncodeRune(o.Utf, r)
				o.Buf = append(o.Buf, o.Utf[:n]...)
			}
		}
	}
	if !tokOk {
		o.Buf = append(o.Buf, '"')
	}
}

// BuildTime appends a time string to the buffer.
func (o *Options) BuildTime(t time.Time) {
	if 0 < len(o.TimeWrap) {
		o.Buf = append(o.Buf, []byte(`{"`)...)
		o.Buf = append(o.Buf, []byte(o.TimeWrap)...)
		o.Buf = append(o.Buf, []byte(`":`)...)
	}
	switch o.TimeFormat {
	case "", "nano":
		o.Buf = append(o.Buf, []byte(strconv.FormatInt(t.UnixNano(), 10))...)
	case "second":
		// Decimal format but float is not accurate enough so build the output
		// in two parts.
		nano := t.UnixNano()
		secs := nano / int64(time.Second)
		if 0 < nano {
			o.Buf = append(o.Buf, []byte(fmt.Sprintf("%d.%09d", secs, nano-(secs*int64(time.Second))))...)
		} else {
			o.Buf = append(o.Buf, []byte(fmt.Sprintf("%d.%09d", secs, -(nano-(secs*int64(time.Second)))))...)
		}
	default:
		o.Buf = append(o.Buf, '"')
		o.Buf = append(o.Buf, []byte(t.Format(o.TimeFormat))...)
		o.Buf = append(o.Buf, '"')
	}
	if 0 < len(o.TimeWrap) {
		o.Buf = append(o.Buf, '}')
	}
}

func (o *Options) buildArray(n gen.Array, depth int) {
	o.Buf = append(o.Buf, '[')
	d2 := depth + 1
	var is string
	var cs string

	if o.Tab || 0 < o.Indent {
		if o.Tab {
			x := depth + 1
			if len(tabs) < x {
				x = len(tabs)
			}
			is = tabs[0:x]
			x = d2 + 1
			if len(tabs) < x {
				x = len(tabs)
			}
			cs = tabs[0:x]
		} else {
			x := depth*o.Indent + 1
			if len(spaces) < x {
				x = len(spaces)
			}
			is = spaces[0:x]
			x = d2*o.Indent + 1
			if len(spaces) < x {
				x = len(spaces)
			}
			cs = spaces[0:x]
		}
		for _, m := range n {
			o.Buf = append(o.Buf, []byte(cs)...)
			o.buildSen(m, d2)
		}
		o.Buf = append(o.Buf, []byte(is)...)
	} else {
		for j, m := range n {
			if 0 < j {
				o.Buf = append(o.Buf, ' ')
			}
			o.buildSen(m, depth)
		}
	}
	o.Buf = append(o.Buf, ']')
}

func (o *Options) buildSimpleArray(n []interface{}, depth int) {
	o.Buf = append(o.Buf, '[')
	d2 := depth + 1
	var is string
	var cs string

	if o.Tab || 0 < o.Indent {
		if o.Tab {
			x := depth + 1
			if len(tabs) < x {
				x = len(tabs)
			}
			is = tabs[0:x]
			x = d2 + 1
			if len(tabs) < x {
				x = len(tabs)
			}
			cs = tabs[0:x]
		} else {
			x := depth*o.Indent + 1
			if len(spaces) < x {
				x = len(spaces)
			}
			is = spaces[0:x]
			x = d2*o.Indent + 1
			if len(spaces) < x {
				x = len(spaces)
			}
			cs = spaces[0:x]
		}
		for _, m := range n {
			o.Buf = append(o.Buf, []byte(cs)...)
			o.buildSen(m, d2)
		}
		o.Buf = append(o.Buf, []byte(is)...)
	} else {
		for j, m := range n {
			if 0 < j {
				o.Buf = append(o.Buf, ' ')
			}
			o.buildSen(m, depth)
		}
	}
	o.Buf = append(o.Buf, ']')
}

func (o *Options) buildObject(n gen.Object, depth int) {
	o.Buf = append(o.Buf, '{')
	d2 := depth + 1
	var is string
	var cs string

	if o.Tab || 0 < o.Indent {
		if o.Tab {
			x := depth + 1
			if len(tabs) < x {
				x = len(tabs)
			}
			is = tabs[0:x]
			x = d2 + 1
			if len(tabs) < x {
				x = len(tabs)
			}
			cs = tabs[0:x]
		} else {
			x := depth*o.Indent + 1
			if len(spaces) < x {
				x = len(spaces)
			}
			is = spaces[0:x]
			x = d2*o.Indent + 1
			if len(spaces) < x {
				x = len(spaces)
			}
			cs = spaces[0:x]
		}
		if o.Sort {
			keys := make([]string, 0, len(n))
			for k := range n {
				keys = append(keys, k)
			}
			sort.Strings(keys)
			for _, k := range keys {
				m := n[k]
				if m == nil && o.OmitNil {
					continue
				}
				o.Buf = append(o.Buf, []byte(cs)...)
				o.BuildString(k)
				o.Buf = append(o.Buf, ':')
				o.Buf = append(o.Buf, ' ')
				o.buildSen(m, d2)
			}
		} else {
			for k, m := range n {
				if m == nil && o.OmitNil {
					continue
				}
				o.Buf = append(o.Buf, []byte(cs)...)
				o.BuildString(k)
				o.Buf = append(o.Buf, ':')
				o.Buf = append(o.Buf, ' ')
				o.buildSen(m, d2)
			}
		}
		o.Buf = append(o.Buf, []byte(is)...)
	} else {
		first := true
		if o.Sort {
			keys := make([]string, 0, len(n))
			for k := range n {
				keys = append(keys, k)
			}
			sort.Strings(keys)
			for _, k := range keys {
				m := n[k]
				if m == nil && o.OmitNil {
					continue
				}
				if first {
					first = false
				} else {
					o.Buf = append(o.Buf, ' ')
				}
				o.BuildString(k)
				o.Buf = append(o.Buf, ':')
				o.buildSen(m, 0)
			}
		} else {
			for k, m := range n {
				if m == nil && o.OmitNil {
					continue
				}
				if first {
					first = false
				} else {
					o.Buf = append(o.Buf, ' ')
				}
				o.BuildString(k)
				o.Buf = append(o.Buf, ':')
				o.buildSen(m, 0)
			}
		}
	}
	o.Buf = append(o.Buf, '}')
}

func (o *Options) buildSimpleObject(n map[string]interface{}, depth int) {
	o.Buf = append(o.Buf, '{')
	d2 := depth + 1
	var is string
	var cs string

	if o.Tab || 0 < o.Indent {
		if o.Tab {
			x := depth + 1
			if len(tabs) < x {
				x = len(tabs)
			}
			is = tabs[0:x]
			x = d2 + 1
			if len(tabs) < x {
				x = len(tabs)
			}
			cs = tabs[0:x]
		} else {
			x := depth*o.Indent + 1
			if len(spaces) < x {
				x = len(spaces)
			}
			is = spaces[0:x]
			x = d2*o.Indent + 1
			if len(spaces) < x {
				x = len(spaces)
			}
			cs = spaces[0:x]
		}
		if o.Sort {
			keys := make([]string, 0, len(n))
			for k := range n {
				keys = append(keys, k)
			}
			sort.Strings(keys)
			for _, k := range keys {
				m := n[k]
				if m == nil && o.OmitNil {
					continue
				}
				o.Buf = append(o.Buf, []byte(cs)...)
				o.BuildString(k)
				o.Buf = append(o.Buf, ':')
				o.Buf = append(o.Buf, ' ')
				o.buildSen(m, d2)
			}
		} else {
			for k, m := range n {
				if m == nil && o.OmitNil {
					continue
				}
				o.Buf = append(o.Buf, []byte(cs)...)
				o.BuildString(k)
				o.Buf = append(o.Buf, ':')
				o.Buf = append(o.Buf, ' ')
				o.buildSen(m, d2)
			}
		}
		o.Buf = append(o.Buf, []byte(is)...)
	} else {
		first := true
		if o.Sort {
			keys := make([]string, 0, len(n))
			for k := range n {
				keys = append(keys, k)
			}
			sort.Strings(keys)
			for _, k := range keys {
				m := n[k]
				if m == nil && o.OmitNil {
					continue
				}
				if first {
					first = false
				} else {
					o.Buf = append(o.Buf, ' ')
				}
				o.BuildString(k)
				o.Buf = append(o.Buf, ':')
				o.buildSen(m, 0)
			}
		} else {
			for k, m := range n {
				if m == nil && o.OmitNil {
					continue
				}
				if first {
					first = false
				} else {
					o.Buf = append(o.Buf, ' ')
				}
				o.BuildString(k)
				o.Buf = append(o.Buf, ':')
				o.buildSen(m, 0)
			}
		}
	}
	o.Buf = append(o.Buf, '}')
}
