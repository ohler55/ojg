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
	if cap(o.buf) < o.InitSize {
		o.buf = make([]byte, 0, o.InitSize)
	} else {
		o.buf = o.buf[:0]
	}
	defer func() {
		if r := recover(); r != nil {
			o.buf = o.buf[:0]
		}
	}()
	o.buildSen(data, 0)

	return string(o.buf)
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
	defer func() {
		if r := recover(); r != nil {
			o.buf = o.buf[:0]
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
		o.buf = append(o.buf, Normal...)
	}
	if w != nil && 0 < len(o.buf) {
		_, err = o.w.Write(o.buf)
	}
	return
}

func (o *Options) buildSen(data interface{}, depth int) {
	switch td := data.(type) {
	case nil:
		o.buf = append(o.buf, []byte("null")...)

	case bool:
		if td {
			o.buf = append(o.buf, []byte("true")...)
		} else {
			o.buf = append(o.buf, []byte("false")...)
		}
	case gen.Bool:
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
	case gen.Int:
		o.buf = append(o.buf, []byte(strconv.FormatInt(int64(td), 10))...)

	case float32:
		o.buf = append(o.buf, []byte(strconv.FormatFloat(float64(td), 'g', -1, 32))...)
	case float64:
		o.buf = append(o.buf, []byte(strconv.FormatFloat(td, 'g', -1, 64))...)
	case gen.Float:
		o.buf = append(o.buf, []byte(strconv.FormatFloat(float64(td), 'g', -1, 64))...)

	case string:
		o.buildString(td)
	case gen.String:
		o.buildString(string(td))

	case time.Time:
		o.buildTime(td)
	case gen.Time:
		o.buildTime(time.Time(td))

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
			o.buildString(fmt.Sprintf("%v", td))
		}
	}
	if o.w != nil && o.WriteLimit < len(o.buf) {
		if _, err := o.w.Write(o.buf); err != nil {
			panic(err)
		}
		o.buf = o.buf[:0]
	}
}

func (o *Options) buildString(s string) {
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
		o.buf = append(o.buf, '"')
	}
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
			o.buf = append(o.buf, []byte{'\n'}...)
		case '\r':
			o.buf = append(o.buf, []byte{'\\', 'r'}...)
		case '\t':
			o.buf = append(o.buf, []byte{'\t'}...)
		case '\u2028':
			o.buf = append(o.buf, []byte(`\u2028`)...)
		case '\u2029':
			o.buf = append(o.buf, []byte(`\u2029`)...)
		default:
			if r < ' ' {
				o.buf = append(o.buf, []byte{'\\', 'u', '0', '0', hex[(r>>4)&0x0f], hex[r&0x0f]}...)
			} else if r < 0x80 {
				o.buf = append(o.buf, byte(r))
			} else {
				if len(o.utf) < utf8.UTFMax {
					o.utf = make([]byte, utf8.UTFMax)
				} else {
					o.utf = o.utf[:cap(o.utf)]
				}
				n := utf8.EncodeRune(o.utf, r)
				o.buf = append(o.buf, o.utf[:n]...)
			}
		}
	}
	if !tokOk {
		o.buf = append(o.buf, '"')
	}
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
			o.buf = append(o.buf, []byte(fmt.Sprintf("%d.%09d", secs, -(nano-(secs*int64(time.Second)))))...)
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

func (o *Options) buildArray(n gen.Array, depth int) {
	o.buf = append(o.buf, '[')
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
			o.buf = append(o.buf, []byte(cs)...)
			o.buildSen(m, d2)
		}
		o.buf = append(o.buf, []byte(is)...)
	} else {
		for j, m := range n {
			if 0 < j {
				o.buf = append(o.buf, ' ')
			}
			o.buildSen(m, depth)
		}
	}
	o.buf = append(o.buf, ']')
}

func (o *Options) buildSimpleArray(n []interface{}, depth int) {
	o.buf = append(o.buf, '[')
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
			o.buf = append(o.buf, []byte(cs)...)
			o.buildSen(m, d2)
		}
		o.buf = append(o.buf, []byte(is)...)
	} else {
		for j, m := range n {
			if 0 < j {
				o.buf = append(o.buf, ' ')
			}
			o.buildSen(m, depth)
		}
	}
	o.buf = append(o.buf, ']')
}

func (o *Options) buildObject(n gen.Object, depth int) {
	o.buf = append(o.buf, '{')
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
				o.buf = append(o.buf, []byte(cs)...)
				o.buildString(k)
				o.buf = append(o.buf, ':')
				o.buf = append(o.buf, ' ')
				o.buildSen(m, d2)
			}
		} else {
			for k, m := range n {
				if m == nil && o.OmitNil {
					continue
				}
				o.buf = append(o.buf, []byte(cs)...)
				o.buildString(k)
				o.buf = append(o.buf, ':')
				o.buf = append(o.buf, ' ')
				o.buildSen(m, d2)
			}
		}
		o.buf = append(o.buf, []byte(is)...)
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
					o.buf = append(o.buf, ' ')
				}
				o.buildString(k)
				o.buf = append(o.buf, ':')
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
					o.buf = append(o.buf, ' ')
				}
				o.buildString(k)
				o.buf = append(o.buf, ':')
				o.buildSen(m, 0)
			}
		}
	}
	o.buf = append(o.buf, '}')
}

func (o *Options) buildSimpleObject(n map[string]interface{}, depth int) {
	o.buf = append(o.buf, '{')
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
				o.buf = append(o.buf, []byte(cs)...)
				o.buildString(k)
				o.buf = append(o.buf, ':')
				o.buf = append(o.buf, ' ')
				o.buildSen(m, d2)
			}
		} else {
			for k, m := range n {
				if m == nil && o.OmitNil {
					continue
				}
				o.buf = append(o.buf, []byte(cs)...)
				o.buildString(k)
				o.buf = append(o.buf, ':')
				o.buf = append(o.buf, ' ')
				o.buildSen(m, d2)
			}
		}
		o.buf = append(o.buf, []byte(is)...)
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
					o.buf = append(o.buf, ' ')
				}
				o.buildString(k)
				o.buf = append(o.buf, ':')
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
					o.buf = append(o.buf, ' ')
				}
				o.buildString(k)
				o.buf = append(o.buf, ':')
				o.buildSen(m, 0)
			}
		}
	}
	o.buf = append(o.buf, '}')
}
