// Copyright (c) 2020, Peter Ohler, All rights reserved.

package oj

import (
	"fmt"
	"io"
	"reflect"
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

// Writer is a JSON writer that includes a reused buffer for reduced
// allocations for repeated encoding calls.
type Writer struct {
	ojg.Options
	buf    []byte
	w      io.Writer
	strict bool
}

func (wr *Writer) JSON(data interface{}) string {
	defer func() {
		if r := recover(); r != nil {
			wr.buf = wr.buf[:0]
		}
	}()
	return wr.MustJSON(data)
}

func (wr *Writer) MustJSON(data interface{}) string {
	if wr.InitSize <= 0 {
		wr.InitSize = 256
	}
	if cap(wr.buf) < wr.InitSize {
		wr.buf = make([]byte, 0, wr.InitSize)
	} else {
		wr.buf = wr.buf[:0]
	}
	wr.buildJSON(data, 0, false)

	return string(wr.buf)
}

// Write a JSON string for the data provided.
func (wr *Writer) Write(w io.Writer, data interface{}) (err error) {
	defer func() {
		if r := recover(); r != nil {
			wr.buf = wr.buf[:0]
			if err, _ = r.(error); err == nil {
				err = fmt.Errorf("%v", r)
			}
		}
	}()
	wr.MustWrite(w, data)
	return
}

// Writep a JSON string for the data provided. If an error occurs panic is
// called with the error.
func (wr *Writer) MustWrite(w io.Writer, data interface{}) {
	wr.w = w
	if wr.InitSize <= 0 {
		wr.InitSize = 256
	}
	if wr.WriteLimit <= 0 {
		wr.WriteLimit = 1024
	}
	if cap(wr.buf) < wr.InitSize {
		wr.buf = make([]byte, 0, wr.InitSize)
	} else {
		wr.buf = wr.buf[:0]
	}
	if wr.Color {
		wr.cbuildJSON(data, 0) // TBD embedded
	} else {
		wr.buildJSON(data, 0, false)
	}
	if 0 < len(wr.buf) {
		if _, err := wr.w.Write(wr.buf); err != nil {
			panic(err)
		}
	}
}

func (wr *Writer) buildJSON(data interface{}, depth int, embedded bool) {
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
		wr.buf = strconv.AppendInt(wr.buf, int64(td), 10)
	case int8:
		wr.buf = strconv.AppendInt(wr.buf, int64(td), 10)
	case int16:
		wr.buf = strconv.AppendInt(wr.buf, int64(td), 10)
	case int32:
		wr.buf = strconv.AppendInt(wr.buf, int64(td), 10)
	case int64:
		wr.buf = strconv.AppendInt(wr.buf, td, 10)
	case uint:
		wr.buf = strconv.AppendUint(wr.buf, uint64(td), 10)
	case uint8:
		wr.buf = strconv.AppendUint(wr.buf, uint64(td), 10)
	case uint16:
		wr.buf = strconv.AppendUint(wr.buf, uint64(td), 10)
	case uint32:
		wr.buf = strconv.AppendUint(wr.buf, uint64(td), 10)
	case uint64:
		wr.buf = strconv.AppendUint(wr.buf, td, 10)
	case gen.Int:
		wr.buf = strconv.AppendInt(wr.buf, int64(td), 10)

	case float32:
		wr.buf = strconv.AppendFloat(wr.buf, float64(td), 'g', -1, 32)
	case float64:
		wr.buf = strconv.AppendFloat(wr.buf, float64(td), 'g', -1, 64)
	case gen.Float:
		wr.buf = strconv.AppendFloat(wr.buf, float64(td), 'g', -1, 64)

	case string:
		wr.buf = ojg.AppendJSONString(wr.buf, td, !wr.HTMLUnsafe)
	case gen.String:
		wr.buf = ojg.AppendJSONString(wr.buf, string(td), !wr.HTMLUnsafe)

	case time.Time:
		wr.buildTime(td)
	case gen.Time:
		wr.buildTime(time.Time(td))

	case []interface{}:
		wr.buildSimpleArray(td, depth)
	case gen.Array:
		wr.buildArray(td, depth)
	case []gen.Node:
		wr.buildArray(gen.Array(td), depth)

	case map[string]interface{}:
		wr.buildSimpleObject(td, depth)
	case gen.Object:
		wr.buildObject(td, depth)

	default:
		if g, _ := data.(alt.Genericer); g != nil {
			wr.buildJSON(g.Generic(), depth, false)
			return
		}
		if simp, _ := data.(alt.Simplifier); simp != nil {
			data = simp.Simplify()
			wr.buildJSON(data, depth, false)
			return
		}
		if 0 < len(wr.CreateKey) {
			// TBD reuse writer options
			ao := alt.Options{
				CreateKey:    wr.CreateKey,
				OmitNil:      wr.OmitNil,
				FullTypePath: wr.FullTypePath,
				UseTags:      wr.UseTags,
				KeyExact:     wr.KeyExact,
				NestEmbed:    wr.NestEmbed,
				BytesAs:      wr.BytesAs,
			}
			wr.buildJSON(alt.Decompose(data, &ao), depth, embedded)
			return
		}
		if !wr.NoReflect {
			rv := reflect.ValueOf(data)
			kind := rv.Kind()
			if kind == reflect.Ptr {
				rv = rv.Elem()
				kind = rv.Kind()
			}
			// TBD reuse writer options
			ao := alt.Options{
				CreateKey:    wr.CreateKey,
				OmitNil:      wr.OmitNil,
				FullTypePath: wr.FullTypePath,
				UseTags:      wr.UseTags,
				KeyExact:     wr.KeyExact,
				NestEmbed:    wr.NestEmbed,
				BytesAs:      wr.BytesAs,
			}
			switch kind {
			case reflect.Struct:
				wr.buildStruct(rv, &ao, depth, embedded)
			case reflect.Slice, reflect.Array:
				wr.buildSlice(rv, &ao, depth)
			default:
				// Not much should get here except Map, Complex and un-decomposable
				// values.
				dec := alt.Decompose(data, &ao)
				wr.buildJSON(dec, depth, false)
				return
			}
		} else if wr.strict {
			panic(fmt.Errorf("%T can not be encoded as a JSON element", data))
		} else {
			wr.buf = ojg.AppendJSONString(wr.buf, fmt.Sprintf("%v", td), !wr.HTMLUnsafe)
		}
	}
	if wr.w != nil && wr.WriteLimit < len(wr.buf) {
		if _, err := wr.w.Write(wr.buf); err != nil {
			panic(err)
		}
		wr.buf = wr.buf[:0]
	}
}

func (wr *Writer) buildTime(t time.Time) {
	if wr.TimeMap {
		wr.buf = append(wr.buf, []byte(`{"`)...)
		wr.buf = append(wr.buf, wr.CreateKey...)
		wr.buf = append(wr.buf, []byte(`":"`)...)
		if wr.FullTypePath {
			wr.buf = append(wr.buf, []byte("time/Time")...)
		} else {
			wr.buf = append(wr.buf, []byte("Time")...)
		}
		wr.buf = append(wr.buf, []byte(`","value":`)...)
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
	if wr.Tab || 0 < wr.Indent {
		var is string
		var cs string
		d2 := depth + 1
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
		for j, m := range n {
			if 0 < j {
				wr.buf = append(wr.buf, ',')
			}
			wr.buf = append(wr.buf, []byte(cs)...)
			wr.buildJSON(m, d2, false)
		}
		wr.buf = append(wr.buf, []byte(is)...)
	} else {
		for j, m := range n {
			if 0 < j {
				wr.buf = append(wr.buf, ',')
			}
			wr.buildJSON(m, depth, false)
		}
	}
	wr.buf = append(wr.buf, ']')
}

func (wr *Writer) buildSimpleArray(n []interface{}, depth int) {
	wr.buf = append(wr.buf, '[')
	if wr.Tab || 0 < wr.Indent {
		var is string
		var cs string
		d2 := depth + 1
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
		for j, m := range n {
			if 0 < j {
				wr.buf = append(wr.buf, ',')
			}
			wr.buf = append(wr.buf, []byte(cs)...)
			wr.buildJSON(m, d2, false)
		}
		wr.buf = append(wr.buf, []byte(is)...)
	} else {
		for j, m := range n {
			if 0 < j {
				wr.buf = append(wr.buf, ',')
			}
			wr.buildJSON(m, depth, false)
		}
	}
	wr.buf = append(wr.buf, ']')
}

func (wr *Writer) buildObject(n gen.Object, depth int) {
	wr.buf = append(wr.buf, '{')
	first := true
	d2 := depth + 1
	if wr.Tab || 0 < wr.Indent {
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
				if first {
					first = false
				} else {
					wr.buf = append(wr.buf, ',')
				}
				wr.buf = append(wr.buf, []byte(cs)...)
				wr.buf = ojg.AppendJSONString(wr.buf, k, !wr.HTMLUnsafe)
				wr.buf = append(wr.buf, ':')
				wr.buf = append(wr.buf, ' ')
				wr.buildJSON(m, d2, false)
			}
		} else {
			for k, m := range n {
				if m == nil && wr.OmitNil {
					continue
				}
				if first {
					first = false
				} else {
					wr.buf = append(wr.buf, ',')
				}
				wr.buf = append(wr.buf, []byte(cs)...)
				wr.buf = ojg.AppendJSONString(wr.buf, k, !wr.HTMLUnsafe)
				wr.buf = append(wr.buf, ':')
				wr.buf = append(wr.buf, ' ')
				wr.buildJSON(m, d2, false)
			}
		}
		wr.buf = append(wr.buf, []byte(is)...)
	} else {
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
					wr.buf = append(wr.buf, ',')
				}
				wr.buf = ojg.AppendJSONString(wr.buf, k, !wr.HTMLUnsafe)
				wr.buf = append(wr.buf, ':')
				wr.buildJSON(m, d2, false)
			}
		} else {
			for k, m := range n {
				if m == nil && wr.OmitNil {
					continue
				}
				if first {
					first = false
				} else {
					wr.buf = append(wr.buf, ',')
				}
				wr.buf = ojg.AppendJSONString(wr.buf, k, !wr.HTMLUnsafe)
				wr.buf = append(wr.buf, ':')
				wr.buildJSON(m, d2, false)
			}
		}
	}
	wr.buf = append(wr.buf, '}')
}

func (wr *Writer) buildSimpleObject(n map[string]interface{}, depth int) {
	wr.buf = append(wr.buf, '{')
	first := true
	d2 := depth + 1
	if wr.Tab || 0 < wr.Indent {
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
				if first {
					first = false
				} else {
					wr.buf = append(wr.buf, ',')
				}
				wr.buf = append(wr.buf, []byte(cs)...)
				wr.buf = ojg.AppendJSONString(wr.buf, k, !wr.HTMLUnsafe)
				wr.buf = append(wr.buf, ':')
				wr.buf = append(wr.buf, ' ')
				wr.buildJSON(m, d2, false)
			}
		} else {
			for k, m := range n {
				if m == nil && wr.OmitNil {
					continue
				}
				if first {
					first = false
				} else {
					wr.buf = append(wr.buf, ',')
				}
				wr.buf = append(wr.buf, []byte(cs)...)
				wr.buf = ojg.AppendJSONString(wr.buf, k, !wr.HTMLUnsafe)
				wr.buf = append(wr.buf, ':')
				wr.buf = append(wr.buf, ' ')
				wr.buildJSON(m, d2, false)
			}
		}
		wr.buf = append(wr.buf, []byte(is)...)
	} else {
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
					wr.buf = append(wr.buf, ',')
				}
				wr.buf = ojg.AppendJSONString(wr.buf, k, !wr.HTMLUnsafe)
				wr.buf = append(wr.buf, ':')
				wr.buildJSON(m, d2, false)
			}
		} else {
			for k, m := range n {
				if m == nil && wr.OmitNil {
					continue
				}
				if first {
					first = false
				} else {
					wr.buf = append(wr.buf, ',')
				}
				wr.buf = ojg.AppendJSONString(wr.buf, k, !wr.HTMLUnsafe)
				wr.buf = append(wr.buf, ':')
				wr.buildJSON(m, d2, false)
			}
		}
	}
	wr.buf = append(wr.buf, '}')
}

func (wr *Writer) buildStruct(rv reflect.Value, opt *alt.Options, depth int, embedded bool) {
	st := ojg.GetStruct(rv.Interface())
	var fields []*ojg.Field
	d2 := depth + 1
	if opt.NestEmbed {
		if opt.UseTags {
			fields = st.OutTag
		} else if opt.KeyExact {
			fields = st.OutName
		} else {
			fields = st.OutLow
		}
	} else {
		if opt.UseTags {
			fields = st.ByTag
		} else if opt.KeyExact {
			fields = st.ByName
		} else {
			fields = st.ByLow
		}
	}
	wr.buf = append(wr.buf, '{')
	first := true
	if wr.Tab || 0 < wr.Indent {
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
		for _, fi := range fields {
			v, omit := fi.Value(rv, opt.OmitNil, embedded)
			if omit || (opt.OmitNil && v == nil) {
				continue
			}
			if first {
				first = false
			} else {
				wr.buf = append(wr.buf, ',')
			}
			wr.buf = append(wr.buf, []byte(cs)...)
			wr.buf = ojg.AppendJSONString(wr.buf, fi.Key, !wr.HTMLUnsafe)
			wr.buf = append(wr.buf, []byte{':', ' '}...)
			wr.buildJSON(v, d2, true)
		}
		wr.buf = append(wr.buf, []byte(is)...)
	} else {
		var v interface{}
		var has bool
		var wrote bool
		empty := true
		for _, fi := range fields {
			wr.buf, v, wrote, has = fi.Append(wr.buf, rv, opt.OmitNil, embedded)
			if !has {
				if wrote {
					empty = false
				}
				continue
			}
			wr.buildJSON(v, d2, true)
			wr.buf = append(wr.buf, ',')
			empty = false
		}
		if !empty {
			wr.buf = wr.buf[:len(wr.buf)-1]
		}
	}
	wr.buf = append(wr.buf, '}')
}

func (wr *Writer) buildSlice(rv reflect.Value, opt *alt.Options, depth int) {
	d2 := depth + 1
	end := rv.Len()
	wr.buf = append(wr.buf, '[')
	if wr.Tab || 0 < wr.Indent {
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
		for j := 0; j < end; j++ {
			if 0 < j {
				wr.buf = append(wr.buf, ',')
			}
			wr.buf = append(wr.buf, []byte(cs)...)
			rm := rv.Index(j)
			switch rm.Kind() {
			case reflect.Struct:
				wr.buildStruct(rm, opt, d2, false)
			case reflect.Slice, reflect.Array:
				wr.buildSlice(rm, opt, d2)
			default:
				wr.buildJSON(rm.Interface(), d2, false)
			}
		}
		wr.buf = append(wr.buf, []byte(is)...)
	} else {
		for j := 0; j < end; j++ {
			if 0 < j {
				wr.buf = append(wr.buf, ',')
			}
			rm := rv.Index(j)
			switch rm.Kind() {
			case reflect.Struct:
				wr.buildStruct(rm, opt, d2, false)
			case reflect.Slice, reflect.Array:
				wr.buildSlice(rm, opt, d2)
			default:
				wr.buildJSON(rm.Interface(), d2, false)
			}
		}
	}
	wr.buf = append(wr.buf, ']')
}
