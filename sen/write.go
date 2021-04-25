// Copyright (c) 2020, Peter Ohler, All rights reserved.

package sen

import (
	"fmt"
	"io"
	"sort"
	"strconv"
	"time"

	"github.com/ohler55/ojg"
	"github.com/ohler55/ojg/alt"
	"github.com/ohler55/ojg/gen"
)

const (
	spaces = "\n                                                                                                                                "
	tabs   = "\n\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t"
)

type Writer struct {
	ojg.Options
	buf []byte
	w   io.Writer
}

func (wr *Writer) buildSen(data interface{}, depth int) {
	switch td := data.(type) {
	case nil:
		wr.buf = append(wr.buf, []byte("null")...)

	case bool:
		if td {
			wr.buf = append(wr.buf, []byte("true")...)
		} else {
			wr.buf = append(wr.buf, []byte("false")...)
		}
	case gen.Bool:
		if td {
			wr.buf = append(wr.buf, []byte("true")...)
		} else {
			wr.buf = append(wr.buf, []byte("false")...)
		}

	case int:
		wr.buf = append(wr.buf, []byte(strconv.FormatInt(int64(td), 10))...)
	case int8:
		wr.buf = append(wr.buf, []byte(strconv.FormatInt(int64(td), 10))...)
	case int16:
		wr.buf = append(wr.buf, []byte(strconv.FormatInt(int64(td), 10))...)
	case int32:
		wr.buf = append(wr.buf, []byte(strconv.FormatInt(int64(td), 10))...)
	case int64:
		wr.buf = append(wr.buf, []byte(strconv.FormatInt(td, 10))...)
	case uint:
		wr.buf = append(wr.buf, []byte(strconv.FormatInt(int64(td), 10))...)
	case uint8:
		wr.buf = append(wr.buf, []byte(strconv.FormatInt(int64(td), 10))...)
	case uint16:
		wr.buf = append(wr.buf, []byte(strconv.FormatInt(int64(td), 10))...)
	case uint32:
		wr.buf = append(wr.buf, []byte(strconv.FormatInt(int64(td), 10))...)
	case uint64:
		wr.buf = append(wr.buf, []byte(strconv.FormatInt(int64(td), 10))...)
	case gen.Int:
		wr.buf = append(wr.buf, []byte(strconv.FormatInt(int64(td), 10))...)

	case float32:
		wr.buf = append(wr.buf, []byte(strconv.FormatFloat(float64(td), 'g', -1, 32))...)
	case float64:
		wr.buf = append(wr.buf, []byte(strconv.FormatFloat(td, 'g', -1, 64))...)
	case gen.Float:
		wr.buf = append(wr.buf, []byte(strconv.FormatFloat(float64(td), 'g', -1, 64))...)

	case string:
		wr.buf = ojg.AppendSENString(wr.buf, td, !wr.HTMLUnsafe)
	case gen.String:
		wr.buf = ojg.AppendSENString(wr.buf, string(td), !wr.HTMLUnsafe)

	case time.Time:
		wr.BuildTime(td)
	case gen.Time:
		wr.BuildTime(time.Time(td))

	case []interface{}:
		wr.buildSimpleArray(td, depth)
	case gen.Array:
		wr.buildArray(td, depth)

	case map[string]interface{}:
		wr.buildSimpleObject(td, depth)
	case gen.Object:
		wr.buildObject(td, depth)

	default:
		if g, _ := data.(alt.Genericer); g != nil {
			wr.buildSen(g.Generic(), depth)
			return
		}
		if simp, _ := data.(alt.Simplifier); simp != nil {
			data = simp.Simplify()
			wr.buildSen(data, depth)
			return
		}
		if 0 < len(wr.CreateKey) {
			ao := alt.Options{CreateKey: wr.CreateKey, OmitNil: wr.OmitNil, FullTypePath: wr.FullTypePath}
			wr.buildSen(alt.Decompose(data, &ao), depth)
			return
		} else {
			wr.buildSen(alt.Decompose(data, &alt.Options{OmitNil: wr.OmitNil}), depth)
			return
		}
	}
	if wr.w != nil && wr.WriteLimit < len(wr.buf) {
		if _, err := wr.w.Write(wr.buf); err != nil {
			panic(err)
		}
		wr.buf = wr.buf[:0]
	}
}

// BuildTime appends a time string to the buffer.
func (wr *Writer) BuildTime(t time.Time) {
	if wr.TimeMap {
		wr.buf = append(wr.buf, []byte(`{"`)...)
		wr.buf = append(wr.buf, wr.CreateKey...)
		wr.buf = append(wr.buf, []byte(`":`)...)
		if wr.FullTypePath {
			wr.buf = append(wr.buf, []byte(`"time/Time"`)...)
		} else {
			wr.buf = append(wr.buf, []byte("Time")...)
		}
		wr.buf = append(wr.buf, []byte(` value:`)...)
	} else if 0 < len(wr.TimeWrap) {
		wr.buf = append(wr.buf, []byte(`{"`)...)
		wr.buf = append(wr.buf, []byte(wr.TimeWrap)...)
		wr.buf = append(wr.buf, []byte(`":`)...)
	}
	switch wr.TimeFormat {
	case "", "nano":
		wr.buf = append(wr.buf, []byte(strconv.FormatInt(t.UnixNano(), 10))...)
	case "second":
		// Decimal format but float is not accurate enough so build the output
		// in two parts.
		nano := t.UnixNano()
		secs := nano / int64(time.Second)
		if 0 < nano {
			wr.buf = append(wr.buf, []byte(fmt.Sprintf("%d.%09d", secs, nano-(secs*int64(time.Second))))...)
		} else {
			wr.buf = append(wr.buf, []byte(fmt.Sprintf("%d.%09d", secs, -(nano-(secs*int64(time.Second)))))...)
		}
	default:
		wr.buf = append(wr.buf, '"')
		wr.buf = append(wr.buf, []byte(t.Format(wr.TimeFormat))...)
		wr.buf = append(wr.buf, '"')
	}
	if 0 < len(wr.TimeWrap) || wr.TimeMap {
		wr.buf = append(wr.buf, '}')
	}
}

func (wr *Writer) buildArray(n gen.Array, depth int) {
	wr.buf = append(wr.buf, '[')
	d2 := depth + 1
	var is string
	var cs string

	if wr.Tab || 0 < wr.Indent {
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
		} else {
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
		for _, m := range n {
			wr.buf = append(wr.buf, []byte(cs)...)
			wr.buildSen(m, d2)
		}
		wr.buf = append(wr.buf, []byte(is)...)
	} else {
		var prev interface{}
		for j, m := range n {
			if 0 < j {
				switch prev.(type) {
				case []interface{}, map[string]interface{}:
				default:
					wr.buf = append(wr.buf, ' ')
				}
			}
			prev = m
			wr.buildSen(m, depth)
		}
	}
	wr.buf = append(wr.buf, ']')
}

func (wr *Writer) buildSimpleArray(n []interface{}, depth int) {
	wr.buf = append(wr.buf, '[')
	d2 := depth + 1
	var is string
	var cs string

	if wr.Tab || 0 < wr.Indent {
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
		} else {
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
		for _, m := range n {
			wr.buf = append(wr.buf, []byte(cs)...)
			wr.buildSen(m, d2)
		}
		wr.buf = append(wr.buf, []byte(is)...)
	} else {
		var prev interface{}
		for j, m := range n {
			if 0 < j {
				switch prev.(type) {
				case []interface{}, map[string]interface{}:
				default:
					wr.buf = append(wr.buf, ' ')
				}
			}
			prev = m
			wr.buildSen(m, depth)
		}
	}
	wr.buf = append(wr.buf, ']')
}

func (wr *Writer) buildObject(n gen.Object, depth int) {
	wr.buf = append(wr.buf, '{')
	d2 := depth + 1
	var is string
	var cs string

	if wr.Tab || 0 < wr.Indent {
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
		} else {
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
				if m == nil && wr.OmitNil {
					continue
				}
				wr.buf = append(wr.buf, []byte(cs)...)
				wr.buf = ojg.AppendSENString(wr.buf, k, !wr.HTMLUnsafe)
				wr.buf = append(wr.buf, ':')
				wr.buf = append(wr.buf, ' ')
				wr.buildSen(m, d2)
			}
		} else {
			for k, m := range n {
				if m == nil && wr.OmitNil {
					continue
				}
				wr.buf = append(wr.buf, []byte(cs)...)
				wr.buf = ojg.AppendSENString(wr.buf, k, !wr.HTMLUnsafe)
				wr.buf = append(wr.buf, ':')
				wr.buf = append(wr.buf, ' ')
				wr.buildSen(m, d2)
			}
		}
		wr.buf = append(wr.buf, []byte(is)...)
	} else {
		first := true
		if wr.Sort {
			keys := make([]string, 0, len(n))
			for k := range n {
				keys = append(keys, k)
			}
			sort.Strings(keys)
			for _, k := range keys {
				m := n[k]
				if m == nil && wr.OmitNil {
					continue
				}
				if first {
					first = false
				} else {
					wr.buf = append(wr.buf, ' ')
				}
				wr.buf = ojg.AppendSENString(wr.buf, k, !wr.HTMLUnsafe)
				wr.buf = append(wr.buf, ':')
				wr.buildSen(m, 0)
			}
		} else {
			for k, m := range n {
				if m == nil && wr.OmitNil {
					continue
				}
				if first {
					first = false
				} else {
					wr.buf = append(wr.buf, ' ')
				}
				wr.buf = ojg.AppendSENString(wr.buf, k, !wr.HTMLUnsafe)
				wr.buf = append(wr.buf, ':')
				wr.buildSen(m, 0)
			}
		}
	}
	wr.buf = append(wr.buf, '}')
}

func (wr *Writer) buildSimpleObject(n map[string]interface{}, depth int) {
	wr.buf = append(wr.buf, '{')
	d2 := depth + 1
	var is string
	var cs string

	if wr.Tab || 0 < wr.Indent {
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
		} else {
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
				if m == nil && wr.OmitNil {
					continue
				}
				wr.buf = append(wr.buf, []byte(cs)...)
				wr.buf = ojg.AppendSENString(wr.buf, k, !wr.HTMLUnsafe)
				wr.buf = append(wr.buf, ':')
				wr.buf = append(wr.buf, ' ')
				wr.buildSen(m, d2)
			}
		} else {
			for k, m := range n {
				if m == nil && wr.OmitNil {
					continue
				}
				wr.buf = append(wr.buf, []byte(cs)...)
				wr.buf = ojg.AppendSENString(wr.buf, k, !wr.HTMLUnsafe)
				wr.buf = append(wr.buf, ':')
				wr.buf = append(wr.buf, ' ')
				wr.buildSen(m, d2)
			}
		}
		wr.buf = append(wr.buf, []byte(is)...)
	} else {
		first := true
		if wr.Sort {
			keys := make([]string, 0, len(n))
			for k := range n {
				keys = append(keys, k)
			}
			sort.Strings(keys)
			for _, k := range keys {
				m := n[k]
				if m == nil && wr.OmitNil {
					continue
				}
				if first {
					first = false
				} else {
					wr.buf = append(wr.buf, ' ')
				}
				wr.buf = ojg.AppendSENString(wr.buf, k, !wr.HTMLUnsafe)
				wr.buf = append(wr.buf, ':')
				wr.buildSen(m, 0)
			}
		} else {
			for k, m := range n {
				if m == nil && wr.OmitNil {
					continue
				}
				if first {
					first = false
				} else {
					wr.buf = append(wr.buf, ' ')
				}
				wr.buf = ojg.AppendSENString(wr.buf, k, !wr.HTMLUnsafe)
				wr.buf = append(wr.buf, ':')
				wr.buildSen(m, 0)
			}
		}
	}
	wr.buf = append(wr.buf, '}')
}
