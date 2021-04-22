// Copyright (c) 2020, Peter Ohler, All rights reserved.

package alt

import (
	"encoding/base64"
	"fmt"
	"math"
	"reflect"
	"time"
)

// 23 for fraction in IEEE 754 which amounts to 7 significant digits. Use base
// 10 so that numbers look correct when displayed in base 10.
const fracMax = 10000000.0

// Dup is an alias for Decompose.
func Dup(v interface{}, options ...*Options) interface{} {
	return Decompose(v, options...)
}

// Decompose creates a simple type converting non simple to simple types using
// either the Simplify() interface or reflection. Unlike Alter() a deep copy
// is returned leaving the original data unchanged.
func Decompose(v interface{}, options ...*Options) interface{} {
	opt := &DefaultOptions
	if 0 < len(options) {
		opt = options[0]
	}
	if opt.Converter != nil {
		v, _ = opt.Converter.convert(v)
	}
	return decompose(v, opt)
}

func decompose(v interface{}, opt *Options) interface{} {
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
			a[i] = decompose(m, opt)
		}
		v = a
	case map[string]interface{}:
		o := map[string]interface{}{}
		for k, m := range tv {
			if mv := decompose(m, opt); mv != nil || !opt.OmitNil {
				if mv != nil || !opt.OmitNil {
					o[k] = mv
				}
			}
		}
		v = o
	case []byte:
		switch opt.BytesAs {
		case BytesAsBase64:
			v = base64.StdEncoding.EncodeToString(tv)
		case BytesAsArray:
			a := make([]interface{}, len(tv))
			for i, m := range tv {
				a[i] = decompose(m, opt)
			}
			v = a
		default:
			v = string(tv)
		}
	default:
		if simp, _ := v.(Simplifier); simp != nil {
			return decompose(simp.Simplify(), opt)
		}
		return reflectValue(reflect.ValueOf(v), opt)
	}
	return v
}

// Alter the data into all simple types converting non simple to simple types
// using either the Simplify() interface or reflection. Unlike Decompose() map and
// slices members are modified if necessary to assure all elements are simple
// types.
func Alter(v interface{}, options ...*Options) interface{} {
	opt := &DefaultOptions
	if 0 < len(options) {
		opt = options[0]
	}
	if opt.Converter != nil {
		v, _ = opt.Converter.convert(v)
	}
	return alter(v, opt)
}

func alter(v interface{}, opt *Options) interface{} {
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
			tv[i] = alter(m, opt)
		}
	case map[string]interface{}:
		for k, m := range tv {
			if mv := alter(m, opt); mv != nil || !opt.OmitNil {
				if mv != nil || !opt.OmitNil {
					tv[k] = mv
				}
			}
		}
	case []byte:
		switch opt.BytesAs {
		case BytesAsBase64:
			v = base64.StdEncoding.EncodeToString(tv)
		case BytesAsArray:
			a := make([]interface{}, len(tv))
			for i, m := range tv {
				a[i] = decompose(m, opt)
			}
			v = a
		default:
			v = string(tv)
		}
	default:
		if simp, _ := v.(Simplifier); simp != nil {
			return alter(simp.Simplify(), opt)
		}
		return reflectValue(reflect.ValueOf(v), opt)
	}
	return v
}

func reflectValue(rv reflect.Value, opt *Options) (v interface{}) {
	switch rv.Kind() {
	case reflect.Invalid, reflect.Uintptr, reflect.UnsafePointer, reflect.Chan, reflect.Func, reflect.Interface:
		v = nil
	case reflect.Complex64, reflect.Complex128:
		v = reflectComplex(rv, opt)
	case reflect.Map:
		v = reflectMap(rv, opt)
	case reflect.Ptr:
		v = reflectValue(rv.Elem(), opt)
	case reflect.Slice, reflect.Array:
		v = reflectArray(rv, opt)
	case reflect.Struct:
		v = reflectStruct(rv, opt)
	default:
		v = rv.Interface()
	}
	return
}

func reflectStruct(rv reflect.Value, opt *Options) interface{} {
	obj := map[string]interface{}{}
	t := rv.Type()
	if 0 < len(opt.CreateKey) {
		if opt.FullTypePath {
			obj[opt.CreateKey] = t.PkgPath() + "/" + t.Name()
		} else {
			obj[opt.CreateKey] = t.Name()
		}
	}
	dc := LookupDecomposer(t)
	var fields []*Field
	if opt.NestEmbed {
		if opt.UseTags {
			fields = dc.OutTag
		} else if opt.KeyExact {
			fields = dc.OutName
		} else {
			fields = dc.OutLow
		}
	} else {
		if opt.UseTags {
			fields = dc.ByTag
		} else if opt.KeyExact {
			fields = dc.ByName
		} else {
			fields = dc.ByLow
		}
	}
	for _, fi := range fields {
		// TBD change to return a reflect.Value
		//  use reflectValue instead of decompose
		//  switch on fv.Kind
		if v, omit := fi.Value(rv, opt); !omit {
			switch v.(type) {
			case bool, nil, string, time.Time:
			case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
			case float32, float64:
			default:
				v = decompose(v, opt)
			}
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
			g = decompose(vv.Interface(), opt)
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
		a[i] = decompose(rv.Index(i).Interface(), opt)
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
