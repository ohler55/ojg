// Copyright (c) 2020, Peter Ohler, All rights reserved.

package oj

import (
	"fmt"
	"reflect"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/ohler55/ojg"
	"github.com/ohler55/ojg/alt"
)

func (wr *Writer) tightJSON(data interface{}) {

	// TBD if marshal and nil (as apposed to empty) the null
	//  use wr.strict field as indicator of marshal called?

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
		wr.buf = ojg.AppendJSONString(wr.buf, td, !wr.HTMLUnsafe)

	case time.Time:
		wr.buildTime(td)

	case []interface{}:
		wr.tightArray(td)

	case map[string]interface{}:
		wr.tightObject(td)

	default:

		// TBD make this a separate function and point to it from wr?

		if g, _ := data.(alt.Genericer); g != nil {
			wr.tightJSON(g.Generic())
			return
		}
		if simp, _ := data.(alt.Simplifier); simp != nil {
			data = simp.Simplify()
			wr.tightJSON(data)
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
			default:
				// Not much should get here except Map, Complex and un-decomposable
				// values.
				dec := alt.Decompose(data, &wr.Options)
				wr.tightJSON(dec)
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

func (wr *Writer) tightArray(n []interface{}) {
	if 0 < len(n) {
		wr.buf = append(wr.buf, '[')
		for _, m := range n {
			wr.tightJSON(m)
			wr.buf = append(wr.buf, ',')
		}
		wr.buf[len(wr.buf)-1] = ']'
	} else {
		wr.buf = append(wr.buf, "[]"...)
	}
}

func (wr *Writer) tightObject(n map[string]interface{}) {
	comma := false
	wr.buf = append(wr.buf, '{')
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
			wr.buf = ojg.AppendJSONString(wr.buf, k, !wr.HTMLUnsafe)
			wr.buf = append(wr.buf, ':')
			wr.tightJSON(m)
			wr.buf = append(wr.buf, ',')
			comma = true
		}
	} else {
		for k, m := range n {
			if m == nil && wr.OmitNil {
				continue
			}
			wr.buf = ojg.AppendJSONString(wr.buf, k, !wr.HTMLUnsafe)
			wr.buf = append(wr.buf, ':')
			wr.tightJSON(m)
			wr.buf = append(wr.buf, ',')
			comma = true
		}
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
		wr.buf = append(wr.buf, '"')
		wr.buf = append(wr.buf, wr.CreateKey...)
		wr.buf = append(wr.buf, `":"`...)
		if wr.FullTypePath {
			wr.buf = append(wr.buf, (st.Type.PkgPath() + "/" + st.Type.Name())...)
		} else {
			wr.buf = append(wr.buf, st.Type.Name()...)
		}
		wr.buf = append(wr.buf, `",`...)
		comma = true
	}
	for _, fi := range fields {
		wr.buf, v, wrote, has = fi.Append(fi, wr.buf, rv, !wr.HTMLUnsafe)
		if wrote {
			wr.buf = append(wr.buf, ',')
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
			wr.tightJSON(v)
		}
		wr.buf = append(wr.buf, ',')
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
			wr.tightJSON(rm.Interface())
		}
		wr.buf = append(wr.buf, ',')
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
			wr.buf = ojg.AppendJSONString(wr.buf, kv.String(), !wr.HTMLUnsafe)
			wr.buf = append(wr.buf, ':')
			wr.tightStruct(rm, st)
		case reflect.Slice, reflect.Array:
			if wr.OmitNil && rm.IsNil() {
				continue
			}
			wr.buf = ojg.AppendJSONString(wr.buf, kv.String(), !wr.HTMLUnsafe)
			wr.buf = append(wr.buf, ':')
			wr.tightSlice(rm, st)
		case reflect.Map:
			if wr.OmitNil && rm.IsNil() {
				continue
			}
			wr.buf = ojg.AppendJSONString(wr.buf, kv.String(), !wr.HTMLUnsafe)
			wr.buf = append(wr.buf, ':')
			wr.tightMap(rm, st)
		default:
			wr.buf = ojg.AppendJSONString(wr.buf, kv.String(), !wr.HTMLUnsafe)
			wr.buf = append(wr.buf, ':')
			wr.tightJSON(rm.Interface())
		}
		wr.buf = append(wr.buf, ',')
		comma = true
	}
	if comma {
		wr.buf[len(wr.buf)-1] = '}'
	} else {
		wr.buf = append(wr.buf, '}')
	}
}
