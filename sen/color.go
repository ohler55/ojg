// Copyright (c) 2020, Peter Ohler, All rights reserved.

package sen

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
		o.Buf = append(o.Buf, o.NullColor...)
		o.Buf = append(o.Buf, []byte("null")...)

	case bool:
		o.Buf = append(o.Buf, o.BoolColor...)
		if td {
			o.Buf = append(o.Buf, []byte("true")...)
		} else {
			o.Buf = append(o.Buf, []byte("false")...)
		}
	case gen.Bool:
		o.Buf = append(o.Buf, o.BoolColor...)
		if td {
			o.Buf = append(o.Buf, []byte("true")...)
		} else {
			o.Buf = append(o.Buf, []byte("false")...)
		}

	case int:
		o.Buf = append(o.Buf, o.NumberColor...)
		o.Buf = append(o.Buf, []byte(strconv.FormatInt(int64(td), 10))...)
	case int8:
		o.Buf = append(o.Buf, o.NumberColor...)
		o.Buf = append(o.Buf, []byte(strconv.FormatInt(int64(td), 10))...)
	case int16:
		o.Buf = append(o.Buf, o.NumberColor...)
		o.Buf = append(o.Buf, []byte(strconv.FormatInt(int64(td), 10))...)
	case int32:
		o.Buf = append(o.Buf, o.NumberColor...)
		o.Buf = append(o.Buf, []byte(strconv.FormatInt(int64(td), 10))...)
	case int64:
		o.Buf = append(o.Buf, o.NumberColor...)
		o.Buf = append(o.Buf, []byte(strconv.FormatInt(td, 10))...)
	case uint:
		o.Buf = append(o.Buf, o.NumberColor...)
		o.Buf = append(o.Buf, []byte(strconv.FormatInt(int64(td), 10))...)
	case uint8:
		o.Buf = append(o.Buf, o.NumberColor...)
		o.Buf = append(o.Buf, []byte(strconv.FormatInt(int64(td), 10))...)
	case uint16:
		o.Buf = append(o.Buf, o.NumberColor...)
		o.Buf = append(o.Buf, []byte(strconv.FormatInt(int64(td), 10))...)
	case uint32:
		o.Buf = append(o.Buf, o.NumberColor...)
		o.Buf = append(o.Buf, []byte(strconv.FormatInt(int64(td), 10))...)
	case uint64:
		o.Buf = append(o.Buf, o.NumberColor...)
		o.Buf = append(o.Buf, []byte(strconv.FormatInt(int64(td), 10))...)
	case gen.Int:
		o.Buf = append(o.Buf, o.NumberColor...)
		o.Buf = append(o.Buf, []byte(strconv.FormatInt(int64(td), 10))...)

	case float32:
		o.Buf = append(o.Buf, o.NumberColor...)
		o.Buf = append(o.Buf, []byte(strconv.FormatFloat(float64(td), 'g', -1, 32))...)
	case float64:
		o.Buf = append(o.Buf, o.NumberColor...)
		o.Buf = append(o.Buf, []byte(strconv.FormatFloat(td, 'g', -1, 64))...)
	case gen.Float:
		o.Buf = append(o.Buf, o.NumberColor...)
		o.Buf = append(o.Buf, []byte(strconv.FormatFloat(float64(td), 'g', -1, 64))...)

	case string:
		o.Buf = append(o.Buf, o.StringColor...)
		o.BuildString(td)
	case gen.String:
		o.Buf = append(o.Buf, o.StringColor...)
		o.BuildString(string(td))

	case time.Time:
		o.Buf = append(o.Buf, o.StringColor...)
		o.BuildTime(td)
	case gen.Time:
		o.Buf = append(o.Buf, o.StringColor...)
		o.BuildTime(time.Time(td))

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
			ao := alt.Options{CreateKey: o.CreateKey, OmitNil: o.OmitNil, FullTypePath: o.FullTypePath}
			o.cbuildJSON(alt.Decompose(data, &ao), depth)
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

func (o *Options) cbuildArray(n gen.Array, depth int) {
	o.Buf = append(o.Buf, o.SyntaxColor...)
	o.Buf = append(o.Buf, '[')

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
		if 0 < j && len(cs) == 0 {
			o.Buf = append(o.Buf, o.SyntaxColor...)
			o.Buf = append(o.Buf, ' ')
		}
		o.Buf = append(o.Buf, []byte(cs)...)
		o.cbuildJSON(m, d2)
	}
	o.Buf = append(o.Buf, []byte(is)...)
	o.Buf = append(o.Buf, o.SyntaxColor...)
	o.Buf = append(o.Buf, ']')
}

func (o *Options) cbuildSimpleArray(n []interface{}, depth int) {
	o.Buf = append(o.Buf, o.SyntaxColor...)
	o.Buf = append(o.Buf, '[')

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
		if 0 < j && len(cs) == 0 {
			o.Buf = append(o.Buf, o.SyntaxColor...)
			o.Buf = append(o.Buf, ' ')
		}
		o.Buf = append(o.Buf, []byte(cs)...)
		o.cbuildJSON(m, d2)
	}
	o.Buf = append(o.Buf, []byte(is)...)
	o.Buf = append(o.Buf, o.SyntaxColor...)
	o.Buf = append(o.Buf, ']')
}

func (o *Options) cbuildObject(n gen.Object, depth int) {
	o.Buf = append(o.Buf, o.SyntaxColor...)
	o.Buf = append(o.Buf, '{')

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
			} else if len(cs) == 0 {
				o.Buf = append(o.Buf, o.SyntaxColor...)
				o.Buf = append(o.Buf, ' ')
			}
			o.Buf = append(o.Buf, []byte(cs)...)
			o.Buf = append(o.Buf, o.KeyColor...)
			o.BuildString(k)
			o.Buf = append(o.Buf, o.SyntaxColor...)
			o.Buf = append(o.Buf, ':')
			if 0 < o.Indent {
				o.Buf = append(o.Buf, ' ')
			}
			o.cbuildJSON(m, d2)
		}
	} else {
		for k, m := range n {
			if m == nil && o.OmitNil {
				continue
			}
			if first {
				first = false
			} else if len(cs) == 0 {
				o.Buf = append(o.Buf, o.SyntaxColor...)
				o.Buf = append(o.Buf, ' ')
			}
			o.Buf = append(o.Buf, []byte(cs)...)
			o.Buf = append(o.Buf, o.KeyColor...)
			o.BuildString(k)
			o.Buf = append(o.Buf, o.SyntaxColor...)
			o.Buf = append(o.Buf, ':')
			if 0 < o.Indent {
				o.Buf = append(o.Buf, ' ')
			}
			o.cbuildJSON(m, d2)
		}
	}
	o.Buf = append(o.Buf, []byte(is)...)
	o.Buf = append(o.Buf, o.SyntaxColor...)
	o.Buf = append(o.Buf, '}')
}

func (o *Options) cbuildSimpleObject(n map[string]interface{}, depth int) {
	o.Buf = append(o.Buf, o.SyntaxColor...)
	o.Buf = append(o.Buf, '{')

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
			} else if len(cs) == 0 {
				o.Buf = append(o.Buf, o.SyntaxColor...)
				o.Buf = append(o.Buf, ' ')
			}
			o.Buf = append(o.Buf, []byte(cs)...)
			o.Buf = append(o.Buf, o.KeyColor...)
			o.BuildString(k)
			o.Buf = append(o.Buf, o.SyntaxColor...)
			o.Buf = append(o.Buf, ':')
			if 0 < o.Indent {
				o.Buf = append(o.Buf, ' ')
			}
			o.cbuildJSON(m, d2)
		}
	} else {
		for k, m := range n {
			if m == nil && o.OmitNil {
				continue
			}
			if first {
				first = false
			} else if len(cs) == 0 {
				o.Buf = append(o.Buf, o.SyntaxColor...)
				o.Buf = append(o.Buf, ' ')
			}
			o.Buf = append(o.Buf, []byte(cs)...)
			o.Buf = append(o.Buf, o.KeyColor...)
			o.BuildString(k)
			o.Buf = append(o.Buf, o.SyntaxColor...)
			o.Buf = append(o.Buf, ':')
			if 0 < o.Indent {
				o.Buf = append(o.Buf, ' ')
			}
			o.cbuildJSON(m, d2)
		}
	}
	o.Buf = append(o.Buf, []byte(is)...)
	o.Buf = append(o.Buf, o.SyntaxColor...)
	o.Buf = append(o.Buf, '}')
}
