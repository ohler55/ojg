// Copyright (c) 2020, Peter Ohler, All rights reserved.

package alt

import (
	"encoding/base64"
	"fmt"
	"math"
	"reflect"
	"time"

	"github.com/ohler55/ojg"
)

// 23 for fraction in IEEE 754 which amounts to 7 significant digits. Use base
// 10 so that numbers look correct when displayed in base 10.
const fracMax = 10000000.0

func decompose(v interface{}, opt *Options, embedded bool) interface{} {
	switch tv := v.(type) {
	case nil, bool, int64, float64, string, time.Time:
	case int:
		v = int64(tv)
	case int8:
		v = int64(tv)
	case int16:
		v = int64(tv)
	case int32:
		v = int64(tv)
	case uint:
		v = int64(tv)
	case uint8:
		v = int64(tv)
	case uint16:
		v = int64(tv)
	case uint32:
		v = int64(tv)
	case uint64:
		v = int64(tv)
	case float32:
		// This small rounding makes the conversion from 32 bit to 64 bit
		// display nicer.
		f, i := math.Frexp(float64(tv))
		f = float64(int64(f*fracMax)) / fracMax
		v = math.Ldexp(f, i)
	case []interface{}:
		a := make([]interface{}, len(tv))
		for i, m := range tv {
			a[i] = decompose(m, opt, false)
		}
		v = a
	case map[string]interface{}:
		o := map[string]interface{}{}
		for k, m := range tv {
			if mv := decompose(m, opt, false); mv != nil || !opt.OmitNil {
				if mv != nil || !opt.OmitNil {
					o[k] = mv
				}
			}
		}
		v = o
	case []byte:
		switch opt.BytesAs {
		case ojg.BytesAsBase64:
			v = base64.StdEncoding.EncodeToString(tv)
		case ojg.BytesAsArray:
			a := make([]interface{}, len(tv))
			for i, m := range tv {
				a[i] = decompose(m, opt, false)
			}
			v = a
		default:
			v = string(tv)
		}
	default:
		if simp, _ := v.(Simplifier); simp != nil {
			return decompose(simp.Simplify(), opt, false)
		}
		return reflectValue(reflect.ValueOf(v), v, opt, embedded)
	}
	return v
}

func alter(v interface{}, opt *Options, embedded bool) interface{} {
	switch tv := v.(type) {
	case bool, nil, int64, float64, string, time.Time:
	case int:
		v = int64(tv)
	case int8:
		v = int64(tv)
	case int16:
		v = int64(tv)
	case int32:
		v = int64(tv)
	case uint:
		v = int64(tv)
	case uint8:
		v = int64(tv)
	case uint16:
		v = int64(tv)
	case uint32:
		v = int64(tv)
	case uint64:
		v = int64(tv)
	case float32:
		// This small rounding makes the conversion from 32 bit to 64 bit
		// display nicer.
		f, i := math.Frexp(float64(tv))
		f = float64(int64(f*fracMax)) / fracMax
		v = math.Ldexp(f, i)
	case []interface{}:
		for i, m := range tv {
			tv[i] = alter(m, opt, false)
		}
	case map[string]interface{}:
		for k, m := range tv {
			if mv := alter(m, opt, false); mv != nil || !opt.OmitNil {
				if mv != nil || !opt.OmitNil {
					tv[k] = mv
				}
			}
		}
	case []byte:
		switch opt.BytesAs {
		case ojg.BytesAsBase64:
			v = base64.StdEncoding.EncodeToString(tv)
		case ojg.BytesAsArray:
			a := make([]interface{}, len(tv))
			for i, m := range tv {
				a[i] = decompose(m, opt, false)
			}
			v = a
		default:
			v = string(tv)
		}
	default:
		if simp, _ := v.(Simplifier); simp != nil {
			return alter(simp.Simplify(), opt, false)
		}
		return reflectValue(reflect.ValueOf(v), v, opt, embedded)
	}
	return v
}

func reflectValue(rv reflect.Value, val interface{}, opt *Options, embedded bool) (v interface{}) {
	switch rv.Kind() {
	case reflect.Interface:
		fmt.Printf("*** interface %T %v\n", v, v)
		v = nil
	case reflect.Invalid, reflect.Uintptr, reflect.UnsafePointer, reflect.Chan, reflect.Func:
		v = nil
	case reflect.Complex64, reflect.Complex128:
		v = reflectComplex(rv, opt)
	case reflect.Map:
		v = reflectMap(rv, opt)
	case reflect.Ptr:
		elem := rv.Elem()
		if elem.IsValid() && elem.CanInterface() {
			v = reflectValue(elem, elem.Interface(), opt, false)
		} else {
			v = nil
		}
	case reflect.Slice, reflect.Array:
		v = reflectArray(rv, opt)
	case reflect.Struct:
		v = reflectStruct(rv, val, opt, embedded)
	default:
		v = val
	}
	return
}

func reflectStruct(rv reflect.Value, val interface{}, opt *Options, embedded bool) interface{} {
	obj := map[string]interface{}{}
	st := ojg.GetStruct(val)
	t := st.Type
	if 0 < len(opt.CreateKey) {
		if opt.FullTypePath {
			obj[opt.CreateKey] = t.PkgPath() + "/" + t.Name()
		} else {
			obj[opt.CreateKey] = t.Name()
		}
	}
	var fields []*ojg.Field
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
	for _, fi := range fields {
		// TBD is using a blank reflect.Value better than a pointer?
		if v, fv, omit := fi.Value(fi, rv); !omit {
			if fv == nil {
				if !opt.OmitNil || v != nil {
					obj[fi.Key] = v
				}
			} else {
				if simp, _ := v.(Simplifier); simp != nil {
					v = simp.Simplify()
				} else if _, ok := v.([]byte); ok {
					v = decompose(v, opt, false)
				} else if opt.NestEmbed {
					v = reflectEmbed(reflect.ValueOf(v), v, opt)
				} else {
					v = decompose(v, opt, false)
				}
				if !opt.OmitNil || v != nil {
					obj[fi.Key] = v
				}
			}
		}
	}
	return obj
}

func reflectEmbed(rv reflect.Value, val interface{}, opt *Options) interface{} {
	obj := map[string]interface{}{}
	st := ojg.GetStruct(val)
	t := st.Type
	if 0 < len(opt.CreateKey) {
		if opt.FullTypePath {
			obj[opt.CreateKey] = t.PkgPath() + "/" + t.Name()
		} else {
			obj[opt.CreateKey] = t.Name()
		}
	}
	var fields []*ojg.Field
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
	for _, fi := range fields {
		fv := rv.FieldByIndex(fi.Index)
		var v interface{}
		if fv.Kind() == reflect.Struct {
			v = reflectEmbed(fv, fv.Interface(), opt)
		} else {
			v = decompose(fv.Interface(), opt, false)
		}
		if !opt.OmitNil || v != nil {
			obj[fi.Key] = v
		}
	}
	return obj
}

func reflectComplex(rv reflect.Value, opt *Options) interface{} {
	c := rv.Complex()
	obj := map[string]interface{}{
		"real": real(c),
		"imag": imag(c),
	}
	if 0 < len(opt.CreateKey) {
		obj[opt.CreateKey] = "complex"
	}
	return obj
}

func reflectMap(rv reflect.Value, opt *Options) interface{} {
	obj := map[string]interface{}{}
	it := rv.MapRange()
	for it.Next() {
		k := it.Key().Interface()
		var g interface{}
		vv := it.Value()
		if !isNil(vv) {
			g = decompose(vv.Interface(), opt, false)
		}
		if g != nil || !opt.OmitNil {
			if ks, ok := k.(string); ok {
				obj[ks] = g
			} else {
				obj[fmt.Sprint(k)] = g
			}
		}
	}
	return obj
}

func reflectArray(rv reflect.Value, opt *Options) interface{} {
	size := rv.Len()
	a := make([]interface{}, size)
	for i := size - 1; 0 <= i; i-- {
		a[i] = decompose(rv.Index(i).Interface(), opt, false)
	}
	return a
}

func isNil(rv reflect.Value) bool {
	switch rv.Kind() {
	case reflect.Interface, reflect.Map, reflect.Ptr, reflect.Slice:
		return rv.IsNil()
	}
	return false
}
