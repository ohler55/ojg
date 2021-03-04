// Copyright (c) 2020, Peter Ohler, All rights reserved.

package oj

import (
	"fmt"
	"sort"
	"strconv"
	"time"

	"github.com/ohler55/ojg/alt"
	"github.com/ohler55/ojg/gen"
)

func (o *Options) cbuildJSON(data interface{}, depth int) {
	switch td := data.(type) {
	case nil:
		o.buf = append(o.buf, o.NullColor...)
		o.buf = append(o.buf, []byte("null")...)

	case bool:
		o.buf = append(o.buf, o.BoolColor...)
		if td {
			o.buf = append(o.buf, []byte("true")...)
		} else {
			o.buf = append(o.buf, []byte("false")...)
		}
	case gen.Bool:
		o.buf = append(o.buf, o.BoolColor...)
		if td {
			o.buf = append(o.buf, []byte("true")...)
		} else {
			o.buf = append(o.buf, []byte("false")...)
		}

	case int:
		o.buf = append(o.buf, o.NumberColor...)
		o.buf = append(o.buf, []byte(strconv.FormatInt(int64(td), 10))...)
	case int8:
		o.buf = append(o.buf, o.NumberColor...)
		o.buf = append(o.buf, []byte(strconv.FormatInt(int64(td), 10))...)
	case int16:
		o.buf = append(o.buf, o.NumberColor...)
		o.buf = append(o.buf, []byte(strconv.FormatInt(int64(td), 10))...)
	case int32:
		o.buf = append(o.buf, o.NumberColor...)
		o.buf = append(o.buf, []byte(strconv.FormatInt(int64(td), 10))...)
	case int64:
		o.buf = append(o.buf, o.NumberColor...)
		o.buf = append(o.buf, []byte(strconv.FormatInt(td, 10))...)
	case uint:
		o.buf = append(o.buf, o.NumberColor...)
		o.buf = append(o.buf, []byte(strconv.FormatInt(int64(td), 10))...)
	case uint8:
		o.buf = append(o.buf, o.NumberColor...)
		o.buf = append(o.buf, []byte(strconv.FormatInt(int64(td), 10))...)
	case uint16:
		o.buf = append(o.buf, o.NumberColor...)
		o.buf = append(o.buf, []byte(strconv.FormatInt(int64(td), 10))...)
	case uint32:
		o.buf = append(o.buf, o.NumberColor...)
		o.buf = append(o.buf, []byte(strconv.FormatInt(int64(td), 10))...)
	case uint64:
		o.buf = append(o.buf, o.NumberColor...)
		o.buf = append(o.buf, []byte(strconv.FormatInt(int64(td), 10))...)
	case gen.Int:
		o.buf = append(o.buf, o.NumberColor...)
		o.buf = append(o.buf, []byte(strconv.FormatInt(int64(td), 10))...)

	case float32:
		o.buf = append(o.buf, o.NumberColor...)
		o.buf = append(o.buf, []byte(strconv.FormatFloat(float64(td), 'g', -1, 32))...)
	case float64:
		o.buf = append(o.buf, o.NumberColor...)
		o.buf = append(o.buf, []byte(strconv.FormatFloat(td, 'g', -1, 64))...)
	case gen.Float:
		o.buf = append(o.buf, o.NumberColor...)
		o.buf = append(o.buf, []byte(strconv.FormatFloat(float64(td), 'g', -1, 64))...)

	case string:
		o.buf = append(o.buf, o.StringColor...)
		o.buildString(td)
	case gen.String:
		o.buf = append(o.buf, o.StringColor...)
		o.buildString(string(td))

	case time.Time:
		o.buf = append(o.buf, o.TimeColor...)
		o.buildTime(td)
	case gen.Time:
		o.buf = append(o.buf, o.TimeColor...)
		o.buildTime(time.Time(td))

	case []interface{}:
		o.cbuildSimpleArray(td, depth)
	case gen.Array:
		o.cbuildArray(td, depth)

	case map[string]interface{}:
		o.cbuildSimpleObject(td, depth)
	case gen.Object:
		o.cbuildObject(td, depth)

	default:
		if g, _ := data.(alt.Genericer); g != nil {
			o.cbuildJSON(g.Generic(), depth)
			return
		}
		if simp, _ := data.(alt.Simplifier); simp != nil {
			data = simp.Simplify()
			o.cbuildJSON(data, depth)
			return
		}
		if 0 < len(o.CreateKey) {
			ao := alt.Options{
				CreateKey:    o.CreateKey,
				OmitNil:      o.OmitNil,
				FullTypePath: o.FullTypePath,
				UseTags:      o.UseTags,
			}
			o.cbuildJSON(alt.Decompose(data, &ao), depth)
			return
		}
		if !o.NoReflect {
			ao := alt.Options{
				CreateKey:    o.CreateKey,
				OmitNil:      o.OmitNil,
				FullTypePath: o.FullTypePath,
				UseTags:      o.UseTags,
			}
			if dec := alt.Decompose(data, &ao); dec != nil {
				o.cbuildJSON(dec, depth)
				return
			}
		}
		o.buildString(fmt.Sprintf("%v", td))
	}
	o.buf = append(o.buf, o.NoColor...)

	if o.w != nil && o.WriteLimit < len(o.buf) {
		if _, err := o.w.Write(o.buf); err != nil {
			panic(err)
		}
		o.buf = o.buf[:0]
	}
}

func (o *Options) cbuildArray(n gen.Array, depth int) {
	o.buf = append(o.buf, o.SyntaxColor...)
	o.buf = append(o.buf, '[')
	o.buf = append(o.buf, o.NoColor...)

	d2 := depth + 1
	var is string
	var cs string
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
	} else if 0 < o.Indent {
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
	for j, m := range n {
		if 0 < j {
			o.buf = append(o.buf, o.SyntaxColor...)
			o.buf = append(o.buf, ',')
			o.buf = append(o.buf, o.NoColor...)
		}
		o.buf = append(o.buf, []byte(cs)...)
		o.cbuildJSON(m, d2)
	}
	o.buf = append(o.buf, []byte(is)...)
	o.buf = append(o.buf, o.SyntaxColor...)
	o.buf = append(o.buf, ']')
}

func (o *Options) cbuildSimpleArray(n []interface{}, depth int) {
	o.buf = append(o.buf, o.SyntaxColor...)
	o.buf = append(o.buf, '[')
	o.buf = append(o.buf, o.NoColor...)

	d2 := depth + 1
	var is string
	var cs string
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
	} else if 0 < o.Indent {
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
	for j, m := range n {
		if 0 < j {
			o.buf = append(o.buf, o.SyntaxColor...)
			o.buf = append(o.buf, ',')
			o.buf = append(o.buf, o.NoColor...)
		}
		o.buf = append(o.buf, []byte(cs)...)
		o.cbuildJSON(m, d2)
	}
	o.buf = append(o.buf, []byte(is)...)
	o.buf = append(o.buf, o.SyntaxColor...)
	o.buf = append(o.buf, ']')
}

func (o *Options) cbuildObject(n gen.Object, depth int) {
	o.buf = append(o.buf, o.SyntaxColor...)
	o.buf = append(o.buf, '{')
	o.buf = append(o.buf, o.NoColor...)

	d2 := depth + 1
	var is string
	var cs string
	first := true
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
	} else if 0 < o.Indent {
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
			if first {
				first = false
			} else {
				o.buf = append(o.buf, o.SyntaxColor...)
				o.buf = append(o.buf, ',')
				o.buf = append(o.buf, o.NoColor...)
			}
			o.buf = append(o.buf, []byte(cs)...)
			o.buf = append(o.buf, o.KeyColor...)
			o.buildString(k)
			o.buf = append(o.buf, o.NoColor...)
			o.buf = append(o.buf, o.SyntaxColor...)
			o.buf = append(o.buf, ':')
			o.buf = append(o.buf, o.NoColor...)
			if 0 < o.Indent {
				o.buf = append(o.buf, ' ')
			}
			o.cbuildJSON(n[k], d2)
		}
	} else {
		for k, m := range n {
			if m == nil && o.OmitNil {
				continue
			}
			if first {
				first = false
			} else {
				o.buf = append(o.buf, o.SyntaxColor...)
				o.buf = append(o.buf, ',')
				o.buf = append(o.buf, o.NoColor...)
			}
			o.buf = append(o.buf, []byte(cs)...)
			o.buf = append(o.buf, o.KeyColor...)
			o.buildString(k)
			o.buf = append(o.buf, o.NoColor...)
			o.buf = append(o.buf, o.SyntaxColor...)
			o.buf = append(o.buf, ':')
			o.buf = append(o.buf, o.NoColor...)
			if 0 < o.Indent {
				o.buf = append(o.buf, ' ')
			}
			o.cbuildJSON(m, d2)
		}
	}
	o.buf = append(o.buf, []byte(is)...)
	o.buf = append(o.buf, o.SyntaxColor...)
	o.buf = append(o.buf, '}')
}

func (o *Options) cbuildSimpleObject(n map[string]interface{}, depth int) {
	o.buf = append(o.buf, o.SyntaxColor...)
	o.buf = append(o.buf, '{')
	o.buf = append(o.buf, o.NoColor...)

	d2 := depth + 1
	var is string
	var cs string
	first := true
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
	} else if 0 < o.Indent {
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
			if first {
				first = false
			} else {
				o.buf = append(o.buf, o.SyntaxColor...)
				o.buf = append(o.buf, ',')
				o.buf = append(o.buf, o.NoColor...)
			}
			o.buf = append(o.buf, []byte(cs)...)
			o.buf = append(o.buf, o.KeyColor...)
			o.buildString(k)
			o.buf = append(o.buf, o.NoColor...)
			o.buf = append(o.buf, o.SyntaxColor...)
			o.buf = append(o.buf, ':')
			o.buf = append(o.buf, o.NoColor...)
			if 0 < o.Indent {
				o.buf = append(o.buf, ' ')
			}
			o.cbuildJSON(n[k], d2)
		}
	} else {
		for k, m := range n {
			if m == nil && o.OmitNil {
				continue
			}
			if first {
				first = false
			} else {
				o.buf = append(o.buf, o.SyntaxColor...)
				o.buf = append(o.buf, ',')
				o.buf = append(o.buf, o.NoColor...)
			}
			o.buf = append(o.buf, []byte(cs)...)
			o.buf = append(o.buf, o.KeyColor...)
			o.buildString(k)
			o.buf = append(o.buf, o.NoColor...)
			o.buf = append(o.buf, o.SyntaxColor...)
			o.buf = append(o.buf, ':')
			o.buf = append(o.buf, o.NoColor...)
			if 0 < o.Indent {
				o.buf = append(o.buf, ' ')
			}
			o.cbuildJSON(m, d2)
		}
	}
	o.buf = append(o.buf, []byte(is)...)
	o.buf = append(o.buf, o.SyntaxColor...)
	o.buf = append(o.buf, '}')
}
