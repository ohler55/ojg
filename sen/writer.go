// Copyright (c) 2021, Peter Ohler, All rights reserved.

package sen

import (
	"encoding"
	"encoding/json"
	"fmt"
	"io"
	"reflect"
	"sort"
	"strconv"
	"strings"
	"time"
	"unsafe"

	"github.com/ohler55/ojg"
	"github.com/ohler55/ojg/alt"
)

const (
	spaces = "\n                                                                                                                                "
	tabs   = "\n\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t"
)

// Writer is a SEN writer that includes a reused buffer for reduced
// allocations for repeated encoding calls.
type Writer struct {
	ojg.Options
	buf           []byte
	w             io.Writer
	appendArray   func(wr *Writer, data []interface{}, depth int)
	appendObject  func(wr *Writer, data map[string]interface{}, depth int)
	appendDefault func(wr *Writer, data interface{}, depth int)
	appendString  func(buf []byte, s string, htmlSafe bool) []byte
	findex        byte
	needSep       bool
}

// SEN writes data, SEN encoded. On error, an empty string is returned.
func (wr *Writer) SEN(data interface{}) string {
	defer func() {
		if r := recover(); r != nil {
			wr.buf = wr.buf[:0]
		}
	}()
	return string(wr.MustSEN(data))
}

// MustSEN writes data, SEN encoded as a []byte and not a string like the
// SEN() function. On error a panic is called with the error.
func (wr *Writer) MustSEN(data interface{}) []byte {
	wr.w = nil
	if wr.InitSize <= 0 {
		wr.InitSize = 256
	}
	if cap(wr.buf) < wr.InitSize {
		wr.buf = make([]byte, 0, wr.InitSize)
	} else {
		wr.buf = wr.buf[:0]
	}
	wr.calcFieldsIndex()
	if wr.Color {
		wr.colorSEN(data, 0)
	} else {
		wr.appendString = ojg.AppendSENString
		if wr.Tab || 0 < wr.Indent {
			wr.appendArray = appendArray
			if wr.Sort {
				wr.appendObject = appendSortObject
			} else {
				wr.appendObject = appendObject
			}
			wr.appendDefault = appendDefault
		} else {
			wr.appendArray = tightArray
			if wr.Sort {
				wr.appendObject = tightSortObject
			} else {
				wr.appendObject = tightObject
			}
			wr.appendDefault = tightDefault
		}
		wr.appendSEN(data, 0)
	}
	return wr.buf
}

// Write a SEN string for the data provided.
func (wr *Writer) Write(w io.Writer, data interface{}) (err error) {
	defer func() {
		if r := recover(); r != nil {
			wr.buf = wr.buf[:0]
			err = ojg.NewError(r)
		}
	}()
	wr.MustWrite(w, data)
	return
}

// MustWrite a SEN string for the data provided. If an error occurs panic is
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
	wr.calcFieldsIndex()
	if wr.Color {
		wr.colorSEN(data, 0)
	} else {
		wr.appendString = ojg.AppendSENString
		if wr.Tab || 0 < wr.Indent {
			wr.appendArray = appendArray
			if wr.Sort {
				wr.appendObject = appendSortObject
			} else {
				wr.appendObject = appendObject
			}
			wr.appendDefault = appendDefault
		} else {
			wr.appendArray = tightArray
			if wr.Sort {
				wr.appendObject = tightSortObject
			} else {
				wr.appendObject = tightObject
			}
			wr.appendDefault = tightDefault
		}
		wr.appendSEN(data, 0)
	}
	if 0 < len(wr.buf) {
		if _, err := wr.w.Write(wr.buf); err != nil {
			panic(err)
		}
	}
}

func (wr *Writer) calcFieldsIndex() {
	wr.findex = 0
	if wr.NestEmbed {
		wr.findex |= maskNested
	}
	if 0 < wr.Indent {
		wr.findex |= maskPretty
	}
	if wr.UseTags {
		wr.findex |= maskByTag
	} else if wr.KeyExact {
		wr.findex |= maskExact
	}
}

func (wr *Writer) appendSEN(data interface{}, depth int) {
	wr.needSep = true
	switch td := data.(type) {
	case nil:
		wr.buf = append(wr.buf, "null"...)

	case bool:
		if td {
			wr.buf = append(wr.buf, "true"...)
		} else {
			wr.buf = append(wr.buf, "false"...)
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

	case float32:
		wr.buf = strconv.AppendFloat(wr.buf, float64(td), 'g', -1, 32)
	case float64:
		wr.buf = strconv.AppendFloat(wr.buf, float64(td), 'g', -1, 64)

	case string:
		wr.buf = wr.appendString(wr.buf, td, !wr.HTMLUnsafe)

	case time.Time:
		wr.buf = wr.AppendTime(wr.buf, td, true)

	case []interface{}:
		wr.appendArray(wr, td, depth)
		wr.needSep = false

	case map[string]interface{}:
		wr.appendObject(wr, td, depth)
		wr.needSep = false

	case alt.Simplifier:
		wr.appendSEN(td.Simplify(), depth)
	case alt.Genericer:
		wr.appendSEN(td.Generic().Simplify(), depth)
	case json.Marshaler:
		out, err := td.MarshalJSON()
		if err != nil {
			panic(err)
		}
		wr.buf = append(wr.buf, out...)
	case encoding.TextMarshaler:
		out, err := td.MarshalText()
		if err != nil {
			panic(err)
		}
		wr.buf = wr.appendString(wr.buf, string(out), !wr.HTMLUnsafe)

	default:
		wr.appendDefault(wr, data, depth)
		if 0 < len(wr.buf) {
			switch wr.buf[len(wr.buf)-1] {
			case '}', ']':
				wr.needSep = false
			default:
			}
		}
	}
	if wr.w != nil && wr.WriteLimit < len(wr.buf) {
		if _, err := wr.w.Write(wr.buf); err != nil {
			panic(err)
		}
		wr.buf = wr.buf[:0]
	}
}

func appendDefault(wr *Writer, data interface{}, depth int) {
	if !wr.NoReflect {
		rv := reflect.ValueOf(data)
		kind := rv.Kind()
		if kind == reflect.Ptr {
			rv = rv.Elem()
			kind = rv.Kind()
		}
		switch kind {
		case reflect.Struct:
			wr.appendStruct(rv, depth, nil)
		case reflect.Slice, reflect.Array:
			wr.appendSlice(rv, depth, nil)
		case reflect.Map:
			wr.appendMap(rv, depth, nil)
		default:
			// Not much should get here except Complex and non-decomposable
			// values.
			dec := alt.Decompose(data, &wr.Options)
			wr.appendSEN(dec, depth)
			return
		}
	} else {
		wr.buf = wr.appendString(wr.buf, fmt.Sprintf("%v", data), !wr.HTMLUnsafe)
	}
}

func appendArray(wr *Writer, n []interface{}, depth int) {
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
	if 0 < len(n) {
		wr.buf = append(wr.buf, '[')
		for _, m := range n {
			wr.buf = append(wr.buf, cs...)
			wr.appendSEN(m, d2)
		}
		wr.buf = append(wr.buf, is...)
		wr.buf = append(wr.buf, ']')
	} else {
		wr.buf = append(wr.buf, "[]"...)
	}
}

func appendObject(wr *Writer, n map[string]interface{}, depth int) {
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
	wr.buf = append(wr.buf, '{')
	for k, m := range n {
		if m == nil && wr.OmitNil {
			continue
		}
		wr.buf = append(wr.buf, cs...)
		wr.buf = wr.appendString(wr.buf, k, !wr.HTMLUnsafe)
		wr.buf = append(wr.buf, ": "...)
		wr.appendSEN(m, d2)
	}
	wr.buf = append(wr.buf, is...)
	wr.buf = append(wr.buf, '}')
}

func appendSortObject(wr *Writer, n map[string]interface{}, depth int) {
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
	keys := make([]string, 0, len(n))
	for k := range n {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	wr.buf = append(wr.buf, '{')
	for _, k := range keys {
		m := n[k]
		if m == nil && wr.OmitNil {
			continue
		}
		wr.buf = append(wr.buf, cs...)
		wr.buf = wr.appendString(wr.buf, k, !wr.HTMLUnsafe)
		wr.buf = append(wr.buf, ": "...)
		wr.appendSEN(m, d2)
	}
	wr.buf = append(wr.buf, is...)
	wr.buf = append(wr.buf, '}')
}

func (wr *Writer) appendStruct(rv reflect.Value, depth int, si *sinfo) {
	if si == nil {
		si = getSinfo(rv.Interface())
	}
	d2 := depth + 1
	fields := si.fields[wr.findex]
	wr.buf = append(wr.buf, '{')
	empty := true
	var v interface{}
	var has bool
	var wrote bool
	indented := false
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
	if 0 < len(wr.CreateKey) {
		wr.buf = append(wr.buf, cs...)
		wr.buf = wr.appendString(wr.buf, wr.CreateKey, !wr.HTMLUnsafe)
		wr.buf = append(wr.buf, `: "`...)
		if wr.FullTypePath {
			wr.buf = append(wr.buf, si.rt.PkgPath()...)
			wr.buf = append(wr.buf, '/')
			wr.buf = append(wr.buf, si.rt.Name()...)
		} else {
			wr.buf = append(wr.buf, si.rt.Name()...)
		}
		wr.buf = append(wr.buf, '"')
		empty = false
	}
	var addr uintptr
	if rv.CanAddr() {
		addr = rv.UnsafeAddr()
	}
	for _, fi := range fields {
		if !indented {
			wr.buf = append(wr.buf, cs...)
			indented = true
		}
		if 0 < addr {
			wr.buf, v, wrote, has = fi.Append(fi, wr.buf, rv, addr, !wr.HTMLUnsafe)
		} else {
			wr.buf, v, wrote, has = fi.iAppend(fi, wr.buf, rv, addr, !wr.HTMLUnsafe)
		}
		if wrote {
			indented = false
			empty = false
			continue
		}
		if !has {
			continue
		}
		indented = false
		var fv reflect.Value
		kind := fi.kind
	Retry:
		switch kind {
		case reflect.Ptr:
			if (*[2]uintptr)(unsafe.Pointer(&v))[1] != 0 { // Check for nil of any type
				fv = reflect.ValueOf(v).Elem()
				kind = fv.Kind()
				v = fv.Interface()
				goto Retry
			}
			if wr.OmitNil {
				wr.buf = wr.buf[:len(wr.buf)-fi.keyLen()]
				indented = true
				continue
			}
			wr.buf = append(wr.buf, "null"...)
		case reflect.Interface:
			if wr.OmitNil && (*[2]uintptr)(unsafe.Pointer(&v))[1] == 0 {
				wr.buf = wr.buf[:len(wr.buf)-fi.keyLen()]
				indented = true
				continue
			}
			wr.appendSEN(v, 0)
		case reflect.Struct:
			if !fv.IsValid() {
				fv = reflect.ValueOf(v)
			}
			wr.appendStruct(fv, d2, fi.elem)
		case reflect.Slice, reflect.Array:
			if !fv.IsValid() {
				fv = reflect.ValueOf(v)
			}
			wr.appendSlice(fv, d2, fi.elem)
		case reflect.Map:
			if !fv.IsValid() {
				fv = reflect.ValueOf(v)
			}
			wr.appendMap(fv, d2, fi.elem)
		default:
			wr.appendSEN(v, d2)
		}
		empty = false
	}
	if indented {
		wr.buf = wr.buf[:len(wr.buf)-len(cs)]
	}
	if !empty {
		wr.buf = append(wr.buf, is...)
	}
	wr.buf = append(wr.buf, '}')
}

func (wr *Writer) appendSlice(rv reflect.Value, depth int, si *sinfo) {
	end := rv.Len()
	if end == 0 {
		wr.buf = append(wr.buf, "[]"...)
		return
	}
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
	wr.buf = append(wr.buf, '[')
	for j := 0; j < end; j++ {
		wr.buf = append(wr.buf, cs...)
		rm := rv.Index(j)
		switch rm.Kind() {
		case reflect.Struct:
			wr.appendStruct(rm, d2, si)
		case reflect.Slice, reflect.Array:
			wr.appendSlice(rm, d2, si)
		case reflect.Map:
			wr.appendMap(rm, d2, si)
		default:
			wr.appendSEN(rm.Interface(), d2)
		}
	}
	wr.buf = append(wr.buf, is...)
	wr.buf = append(wr.buf, ']')
}

func (wr *Writer) appendMap(rv reflect.Value, depth int, si *sinfo) {
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
	keys := rv.MapKeys()
	if wr.Sort {
		sort.Slice(keys, func(i, j int) bool { return 0 > strings.Compare(keys[i].String(), keys[j].String()) })
	}
	empty := true
	wr.buf = append(wr.buf, '{')
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
			wr.buf = append(wr.buf, cs...)
			wr.buf = wr.appendString(wr.buf, kv.String(), !wr.HTMLUnsafe)
			wr.buf = append(wr.buf, ": "...)
			wr.appendStruct(rm, d2, si)
		case reflect.Slice, reflect.Array:
			if wr.OmitNil && rm.Len() == 0 {
				continue
			}
			wr.buf = append(wr.buf, cs...)
			wr.buf = wr.appendString(wr.buf, kv.String(), !wr.HTMLUnsafe)
			wr.buf = append(wr.buf, ": "...)
			wr.appendSlice(rm, d2, si)
		case reflect.Map:
			if wr.OmitNil && rm.Len() == 0 {
				continue
			}
			wr.buf = append(wr.buf, cs...)
			wr.buf = wr.appendString(wr.buf, kv.String(), !wr.HTMLUnsafe)
			wr.buf = append(wr.buf, ": "...)
			wr.appendMap(rm, d2, si)
		default:
			wr.buf = append(wr.buf, cs...)
			wr.buf = wr.appendString(wr.buf, kv.String(), !wr.HTMLUnsafe)
			wr.buf = append(wr.buf, ": "...)
			wr.appendSEN(rm.Interface(), d2)
		}
		empty = false
	}
	if !empty {
		wr.buf = append(wr.buf, is...)
	}
	wr.buf = append(wr.buf, '}')
}
