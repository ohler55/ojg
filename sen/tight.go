// Copyright (c) 2021, Peter Ohler, All rights reserved.

package sen

import (
	"fmt"
	"reflect"
	"sort"
	"strings"
	"unsafe"

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
			wr.tightMap(rv, nil)
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

func (wr *Writer) tightStruct(rv reflect.Value, si *sinfo) {
	if si == nil {
		si = getSinfo(rv.Interface())
	}
	fields := si.fields[wr.findex&ojg.MaskIndex]
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
			wr.buf = append(wr.buf, si.rt.PkgPath()...)
			wr.buf = append(wr.buf, '/')
			wr.buf = append(wr.buf, si.rt.Name()...)
			wr.buf = append(wr.buf, '"')
		} else {
			wr.buf = wr.appendString(wr.buf, si.rt.Name(), !wr.HTMLUnsafe)
		}
		wr.buf = append(wr.buf, ' ')
		comma = true
	}
	var addr uintptr
	if rv.CanAddr() {
		addr = rv.UnsafeAddr()
	}
	for _, fi := range fields {
		if 0 < addr {
			wr.buf, v, wrote, has = fi.Append(fi, wr.buf, rv, addr, !wr.HTMLUnsafe)
		} else {
			wr.buf, v, wrote, has = fi.iAppend(fi, wr.buf, rv, addr, !wr.HTMLUnsafe)
		}
		if wrote {
			wr.buf = append(wr.buf, ' ')
			comma = true
			continue
		}
		if !has {
			continue
		}
		var fv reflect.Value
		kind := fi.kind
		if kind == reflect.Ptr {
			if (*[2]uintptr)(unsafe.Pointer(&v))[1] != 0 { // Check for nil of any type
				fv = reflect.ValueOf(v).Elem()
				kind = fv.Kind()
				v = fv.Interface()
			} else if wr.OmitNil {
				wr.buf = wr.buf[:len(wr.buf)-fi.KeyLen()]
				continue
			}
		}
		switch kind {
		case reflect.Struct:
			if !fv.IsValid() {
				fv = reflect.ValueOf(v)
			}
			wr.tightStruct(fv, fi.elem)
		case reflect.Slice, reflect.Array:
			if !fv.IsValid() {
				fv = reflect.ValueOf(v)
			}
			wr.tightSlice(fv, fi.elem)
		case reflect.Map:
			if !fv.IsValid() {
				fv = reflect.ValueOf(v)
			}
			wr.tightMap(fv, fi.elem)
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

func (wr *Writer) tightSlice(rv reflect.Value, si *sinfo) {
	end := rv.Len()
	comma := false
	wr.buf = append(wr.buf, '[')
	for j := 0; j < end; j++ {
		rm := rv.Index(j)
		switch rm.Kind() {
		case reflect.Struct:
			wr.tightStruct(rm, si)
		case reflect.Slice, reflect.Array:
			wr.tightSlice(rm, si)
		case reflect.Map:
			wr.tightMap(rm, si)
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

func (wr *Writer) tightMap(rv reflect.Value, si *sinfo) {
	wr.buf = append(wr.buf, '{')
	keys := rv.MapKeys()
	if wr.Sort {
		sort.Slice(keys, func(i, j int) bool { return 0 > strings.Compare(keys[i].String(), keys[j].String()) })
	}
	comma := false
	for _, kv := range keys {
		rm := rv.MapIndex(kv)
		if rm.Kind() == reflect.Ptr {
			if rm.IsNil() {
				if wr.OmitNil {
					continue
				}
			} else {
				rm = rm.Elem()
			}
		}
		switch rm.Kind() {
		case reflect.Struct:
			wr.buf = ojg.AppendSENString(wr.buf, kv.String(), !wr.HTMLUnsafe)
			wr.buf = append(wr.buf, ':')
			wr.tightStruct(rm, si)
		case reflect.Slice, reflect.Array:
			if wr.OmitNil && rm.Len() == 0 {
				continue
			}
			wr.buf = ojg.AppendSENString(wr.buf, kv.String(), !wr.HTMLUnsafe)
			wr.buf = append(wr.buf, ':')
			wr.tightSlice(rm, si)
		case reflect.Map:
			if wr.OmitNil && rm.Len() == 0 {
				continue
			}
			wr.buf = ojg.AppendSENString(wr.buf, kv.String(), !wr.HTMLUnsafe)
			wr.buf = append(wr.buf, ':')
			wr.tightMap(rm, si)
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
