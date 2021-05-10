// Copyright (c) 2021, Peter Ohler, All rights reserved.

package sen

import (
	"fmt"
	"reflect"
	"sort"
	"strings"

	"github.com/ohler55/ojg"
	"github.com/ohler55/ojg/alt"
)

func tightDefault(wr *Writer, data interface{}, _ int) {
	if simp, _ := data.(alt.Simplifier); simp != nil {
		data = simp.Simplify()
		wr.appendSEN(data, 0)
		return
	}
	if g, _ := data.(alt.Genericer); g != nil {
		wr.appendSEN(g.Generic().Simplify(), 0)
		return
	}
	if !wr.NoReflect {
		rv := reflect.ValueOf(data)
		kind := rv.Kind()
		if kind == reflect.Ptr {
			rv = rv.Elem()
			kind = rv.Kind()
		}
		switch kind {
		case reflect.Struct:
			wr.tightStruct(rv, nil)
		case reflect.Slice, reflect.Array:
			wr.tightSlice(rv, nil)
		case reflect.Map:
			wr.tightSlice(rv, nil)
		default:
			// Not much should get here except Map, Complex and un-decomposable
			// values.
			dec := alt.Decompose(data, &wr.Options)
			wr.appendSEN(dec, 0)
			return
		}
	} else {
		wr.buf = ojg.AppendSENString(wr.buf, fmt.Sprintf("%v", data), !wr.HTMLUnsafe)
	}
}

func tightArray(wr *Writer, n []interface{}, _ int) {
	if 0 < len(n) {
		space := false
		wr.buf = append(wr.buf, '[')
		for _, m := range n {
			wr.appendSEN(m, 0)
			if wr.needSep {
				wr.buf = append(wr.buf, ' ')
				space = true
			} else {
				space = false
			}
		}
		if space {
			wr.buf[len(wr.buf)-1] = ']'
		} else {
			wr.buf = append(wr.buf, ']')
		}
	} else {
		wr.buf = append(wr.buf, "[]"...)
	}
}

func tightObject(wr *Writer, n map[string]interface{}, _ int) {
	comma := false
	wr.buf = append(wr.buf, '{')
	for k, m := range n {
		if m == nil && wr.OmitNil {
			continue
		}
		wr.buf = ojg.AppendSENString(wr.buf, k, !wr.HTMLUnsafe)
		wr.buf = append(wr.buf, ':')
		wr.appendSEN(m, 0)
		wr.buf = append(wr.buf, ' ')
		comma = true
	}
	if comma {
		wr.buf[len(wr.buf)-1] = '}'
	} else {
		wr.buf = append(wr.buf, '}')
	}
}

func tightSortObject(wr *Writer, n map[string]interface{}, _ int) {
	comma := false
	wr.buf = append(wr.buf, '{')
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
		wr.buf = ojg.AppendSENString(wr.buf, k, !wr.HTMLUnsafe)
		wr.buf = append(wr.buf, ':')
		wr.appendSEN(m, 0)
		wr.buf = append(wr.buf, ' ')
		comma = true
	}
	if comma {
		wr.buf[len(wr.buf)-1] = '}'
	} else {
		wr.buf = append(wr.buf, '}')
	}
}

func (wr *Writer) tightStruct(rv reflect.Value, st *ojg.Struct) {
	if st == nil {
		st = ojg.GetStruct(rv.Interface())
	}
	fields := st.Fields[wr.findex&ojg.MaskIndex]
	wr.buf = append(wr.buf, '{')
	var v interface{}
	var has bool
	var wrote bool
	comma := false
	if 0 < len(wr.CreateKey) {
		wr.buf = wr.appendString(wr.buf, wr.CreateKey, !wr.HTMLUnsafe)
		wr.buf = append(wr.buf, ':')
		if wr.FullTypePath {
			wr.buf = append(wr.buf, '"')
			wr.buf = append(wr.buf, st.Type.PkgPath()...)
			wr.buf = append(wr.buf, '/')
			wr.buf = append(wr.buf, st.Type.Name()...)
			wr.buf = append(wr.buf, '"')
		} else {
			wr.buf = wr.appendString(wr.buf, st.Type.Name(), !wr.HTMLUnsafe)
		}
		wr.buf = append(wr.buf, ' ')
		comma = true
	}
	for _, fi := range fields {
		wr.buf, v, wrote, has = fi.Append(fi, wr.buf, rv, !wr.HTMLUnsafe)
		if wrote {
			wr.buf = append(wr.buf, ' ')
			comma = true
			continue
		}
		if !has {
			continue
		}
		var fv reflect.Value
		kind := fi.Kind
		if kind == reflect.Ptr {
			fv = reflect.ValueOf(v).Elem()
			if !fv.IsValid() {
				continue
			}
			kind = fv.Kind()
			v = fv.Interface()
		}
		switch kind {
		case reflect.Struct:
			if !fv.IsValid() {
				fv = reflect.ValueOf(v)
			}
			wr.tightStruct(fv, fi.Elem)
		case reflect.Slice, reflect.Array:
			if !fv.IsValid() {
				fv = reflect.ValueOf(v)
			}
			wr.tightSlice(fv, fi.Elem)
		case reflect.Map:
			if !fv.IsValid() {
				fv = reflect.ValueOf(v)
			}
			wr.tightMap(fv, fi.Elem)
		default:
			wr.appendSEN(v, 0)
		}
		wr.buf = append(wr.buf, ' ')
		comma = true
	}
	if comma {
		wr.buf[len(wr.buf)-1] = '}'
	} else {
		wr.buf = append(wr.buf, '}')
	}
}

func (wr *Writer) tightSlice(rv reflect.Value, st *ojg.Struct) {
	end := rv.Len()
	comma := false
	wr.buf = append(wr.buf, '[')
	for j := 0; j < end; j++ {
		rm := rv.Index(j)
		if rm.Kind() == reflect.Ptr {
			rm = rm.Elem()
		}
		switch rm.Kind() {
		case reflect.Struct:
			wr.tightStruct(rm, st)
		case reflect.Slice, reflect.Array:
			wr.tightSlice(rm, st)
		case reflect.Map:
			wr.tightMap(rm, st)
		default:
			wr.appendSEN(rm.Interface(), 0)
		}
		wr.buf = append(wr.buf, ' ')
		comma = true
	}
	if comma {
		wr.buf[len(wr.buf)-1] = ']'
	} else {
		wr.buf = append(wr.buf, ']')
	}
}

func (wr *Writer) tightMap(rv reflect.Value, st *ojg.Struct) {
	wr.buf = append(wr.buf, '{')
	keys := rv.MapKeys()
	if wr.Sort {
		sort.Slice(keys, func(i, j int) bool { return 0 < strings.Compare(keys[i].String(), keys[j].String()) })
	}
	comma := false
	for _, kv := range keys {
		rm := rv.MapIndex(kv)
		if rm.Kind() == reflect.Ptr {
			if wr.OmitNil && rm.IsNil() {
				continue
			}
			rm = rm.Elem()
		}
		switch rm.Kind() {
		case reflect.Struct:
			wr.buf = ojg.AppendSENString(wr.buf, kv.String(), !wr.HTMLUnsafe)
			wr.buf = append(wr.buf, ':')
			wr.tightStruct(rm, st)
		case reflect.Slice, reflect.Array:
			if wr.OmitNil && rm.IsNil() {
				continue
			}
			wr.buf = ojg.AppendSENString(wr.buf, kv.String(), !wr.HTMLUnsafe)
			wr.buf = append(wr.buf, ':')
			wr.tightSlice(rm, st)
		case reflect.Map:
			if wr.OmitNil && rm.IsNil() {
				continue
			}
			wr.buf = ojg.AppendSENString(wr.buf, kv.String(), !wr.HTMLUnsafe)
			wr.buf = append(wr.buf, ':')
			wr.tightMap(rm, st)
		default:
			wr.buf = ojg.AppendSENString(wr.buf, kv.String(), !wr.HTMLUnsafe)
			wr.buf = append(wr.buf, ':')
			wr.appendSEN(rm.Interface(), 0)
		}
		wr.buf = append(wr.buf, ' ')
		comma = true
	}
	if comma {
		wr.buf[len(wr.buf)-1] = '}'
	} else {
		wr.buf = append(wr.buf, '}')
	}
}
