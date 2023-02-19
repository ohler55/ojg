// Copyright (c) 2020, Peter Ohler, All rights reserved.

package sen

import (
	"sort"
	"strconv"
	"time"

	"github.com/ohler55/ojg"
	"github.com/ohler55/ojg/alt"
)

func (wr *Writer) colorSEN(data any, depth int) {
	switch td := data.(type) {
	case nil:
		wr.buf = append(wr.buf, wr.NullColor...)
		wr.buf = append(wr.buf, []byte("null")...)

	case bool:
		wr.buf = append(wr.buf, wr.BoolColor...)
		if td {
			wr.buf = append(wr.buf, []byte("true")...)
		} else {
			wr.buf = append(wr.buf, []byte("false")...)
		}

	case int:
		wr.buf = append(wr.buf, wr.NumberColor...)
		wr.buf = strconv.AppendInt(wr.buf, int64(td), 10)
	case int8:
		wr.buf = append(wr.buf, wr.NumberColor...)
		wr.buf = strconv.AppendInt(wr.buf, int64(td), 10)
	case int16:
		wr.buf = append(wr.buf, wr.NumberColor...)
		wr.buf = strconv.AppendInt(wr.buf, int64(td), 10)
	case int32:
		wr.buf = append(wr.buf, wr.NumberColor...)
		wr.buf = strconv.AppendInt(wr.buf, int64(td), 10)
	case int64:
		wr.buf = append(wr.buf, wr.NumberColor...)
		wr.buf = strconv.AppendInt(wr.buf, td, 10)
	case uint:
		wr.buf = append(wr.buf, wr.NumberColor...)
		wr.buf = strconv.AppendUint(wr.buf, uint64(td), 10)
	case uint8:
		wr.buf = append(wr.buf, wr.NumberColor...)
		wr.buf = strconv.AppendUint(wr.buf, uint64(td), 10)
	case uint16:
		wr.buf = append(wr.buf, wr.NumberColor...)
		wr.buf = strconv.AppendUint(wr.buf, uint64(td), 10)
	case uint32:
		wr.buf = append(wr.buf, wr.NumberColor...)
		wr.buf = strconv.AppendUint(wr.buf, uint64(td), 10)
	case uint64:
		wr.buf = append(wr.buf, wr.NumberColor...)
		wr.buf = strconv.AppendUint(wr.buf, td, 10)

	case float32:
		wr.buf = append(wr.buf, wr.NumberColor...)
		wr.buf = strconv.AppendFloat(wr.buf, float64(td), 'g', -1, 32)
	case float64:
		wr.buf = append(wr.buf, wr.NumberColor...)
		wr.buf = strconv.AppendFloat(wr.buf, td, 'g', -1, 64)

	case string:
		wr.buf = append(wr.buf, wr.StringColor...)
		wr.buf = ojg.AppendSENString(wr.buf, td, !wr.HTMLUnsafe)

	case time.Time:
		wr.buf = append(wr.buf, wr.TimeColor...)
		wr.buf = wr.AppendTime(wr.buf, td, true)

	case []any:
		wr.colorArray(td, depth)

	case map[string]any:
		wr.colorObject(td, depth)

	default:
		if simp, _ := data.(alt.Simplifier); simp != nil {
			data = simp.Simplify()
			wr.colorSEN(data, depth)
			return
		}
		if g, _ := data.(alt.Genericer); g != nil {
			wr.colorSEN(g.Generic().Simplify(), depth)
			return
		}
		if 0 < len(wr.CreateKey) {
			ao := alt.Options{CreateKey: wr.CreateKey, OmitNil: wr.OmitNil, FullTypePath: wr.FullTypePath}
			wr.colorSEN(alt.Decompose(data, &ao), depth)
			return
		}
		wr.colorSEN(alt.Decompose(data, &alt.Options{OmitNil: wr.OmitNil}), depth)
	}
	wr.buf = append(wr.buf, wr.NoColor...)

	if wr.w != nil && wr.WriteLimit < len(wr.buf) {
		if _, err := wr.w.Write(wr.buf); err != nil {
			panic(err)
		}
		wr.buf = wr.buf[:0]
	}
}

func (wr *Writer) colorArray(n []any, depth int) {
	wr.buf = append(wr.buf, wr.SyntaxColor...)
	wr.buf = append(wr.buf, '[')
	wr.buf = append(wr.buf, wr.NoColor...)

	d2 := depth + 1
	var is string
	var cs string
	if wr.Tab {
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
	} else if 0 < wr.Indent {
		x := depth*wr.Indent + 1
		if len(spaces) < x {
			x = len(spaces)
		}
		is = spaces[0:x]
		x = d2*wr.Indent + 1
		if len(spaces) < x {
			x = len(spaces)
		}
		cs = spaces[0:x]
	}
	for j, m := range n {
		if 0 < j && len(cs) == 0 {
			wr.buf = append(wr.buf, ' ')
		}
		wr.buf = append(wr.buf, []byte(cs)...)
		wr.colorSEN(m, d2)
	}
	wr.buf = append(wr.buf, []byte(is)...)
	wr.buf = append(wr.buf, wr.SyntaxColor...)
	wr.buf = append(wr.buf, ']')
}

func (wr *Writer) colorObject(n map[string]any, depth int) {
	wr.buf = append(wr.buf, wr.SyntaxColor...)
	wr.buf = append(wr.buf, '{')
	wr.buf = append(wr.buf, wr.NoColor...)

	d2 := depth + 1
	var is string
	var cs string
	first := true
	if wr.Tab {
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
	} else if 0 < wr.Indent {
		x := depth*wr.Indent + 1
		if len(spaces) < x {
			x = len(spaces)
		}
		is = spaces[0:x]
		x = d2*wr.Indent + 1
		if len(spaces) < x {
			x = len(spaces)
		}
		cs = spaces[0:x]
	}
	if wr.Sort {
		keys := make([]string, 0, len(n))
		for k := range n {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		for _, k := range keys {
			m := n[k]
			switch tm := m.(type) {
			case nil:
				if wr.OmitNil {
					continue
				}
			case string:
				if wr.OmitEmpty && len(tm) == 0 {
					continue
				}
			case map[string]any:
				if wr.OmitEmpty && len(tm) == 0 {
					continue
				}
			case []any:
				if wr.OmitEmpty && len(tm) == 0 {
					continue
				}
			}
			if first {
				first = false
			} else if len(cs) == 0 {
				wr.buf = append(wr.buf, ' ')
			}
			wr.buf = append(wr.buf, []byte(cs)...)
			wr.buf = append(wr.buf, wr.KeyColor...)
			wr.buf = ojg.AppendSENString(wr.buf, k, !wr.HTMLUnsafe)
			wr.buf = append(wr.buf, wr.NoColor...)
			wr.buf = append(wr.buf, wr.SyntaxColor...)
			wr.buf = append(wr.buf, ':')
			wr.buf = append(wr.buf, wr.NoColor...)
			if 0 < wr.Indent {
				wr.buf = append(wr.buf, ' ')
			}
			wr.colorSEN(m, d2)
		}
	} else {
		for k, m := range n {
			switch tm := m.(type) {
			case nil:
				if wr.OmitNil {
					continue
				}
			case string:
				if wr.OmitEmpty && len(tm) == 0 {
					continue
				}
			case map[string]any:
				if wr.OmitEmpty && len(tm) == 0 {
					continue
				}
			case []any:
				if wr.OmitEmpty && len(tm) == 0 {
					continue
				}
			}
			if first {
				first = false
			} else if len(cs) == 0 {
				wr.buf = append(wr.buf, ' ')
			}
			wr.buf = append(wr.buf, []byte(cs)...)
			wr.buf = append(wr.buf, wr.KeyColor...)
			wr.buf = ojg.AppendSENString(wr.buf, k, !wr.HTMLUnsafe)
			wr.buf = append(wr.buf, wr.NoColor...)
			wr.buf = append(wr.buf, wr.SyntaxColor...)
			wr.buf = append(wr.buf, ':')
			wr.buf = append(wr.buf, wr.NoColor...)
			if 0 < wr.Indent {
				wr.buf = append(wr.buf, ' ')
			}
			wr.colorSEN(m, d2)
		}
	}
	wr.buf = append(wr.buf, []byte(is)...)
	wr.buf = append(wr.buf, wr.SyntaxColor...)
	wr.buf = append(wr.buf, '}')
}
