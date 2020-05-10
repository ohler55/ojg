// Copyright (c) 2020, Peter Ohler, All rights reserved.

package simple

import (
	"fmt"
	"math"
	"reflect"
	"time"
)

// 23 for fraction in IEEE 754 which amounts to 7 significant digits. Use base
// 10 so that numbers look correct when displayed in base 10.
const fracMax = 10000000.0

type ConvOptions struct {
	CreateKey    string
	FullTypePath bool
	OmitNil      bool
}

var DefaultConvOptions = ConvOptions{
	CreateKey:    "type",
	FullTypePath: false,
	OmitNil:      true,
}

// Decompose is an alias for From(). It exists for symmetry in function names
// with Decompose vs Recompose.
func Decompose(v interface{}, options ...*ConvOptions) interface{} {
	return From(v, options...)
}

// From creates a simple type converting non simple to simple types using
// either the Simplify() interface or reflection. Unlike Alter() a deep copy
// is returned leaving the original data unchanged.
func From(v interface{}, options ...*ConvOptions) interface{} {
	opt := &DefaultConvOptions
	if 0 < len(options) {
		opt = options[0]
	}
	if v != nil {
		switch tv := v.(type) {
		case bool, int64, float64, string, time.Time:
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
				a[i] = From(m)
			}
			v = a
		case map[string]interface{}:
			o := map[string]interface{}{}
			for k, m := range tv {
				if mv := From(m); mv != nil || !opt.OmitNil {
					o[k] = mv
				}
			}
			v = o
		default:
			if simp, _ := v.(Simplifier); simp != nil {
				return From(simp.Simplify())
			}
			return From(reflectData(v, opt))
		}
	}
	return v
}

// Alter the data into all simple types converting non simple to simple types
// using either the Simplify() interface or reflection. Unlike From() map and
// slices members are modified if necessary to assure all elements are simple
// types.
func Alter(v interface{}, options ...*ConvOptions) interface{} {
	opt := &DefaultConvOptions
	if 0 < len(options) {
		opt = options[0]
	}
	if v != nil {
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
				tv[i] = Alter(m)
			}
		case map[string]interface{}:
			for k, m := range tv {
				if mv := Alter(m); mv != nil || !opt.OmitNil {
					tv[k] = mv
				}
			}
		default:
			if simp, _ := v.(Simplifier); simp != nil {
				return Alter(simp.Simplify())
			}
			return Alter(reflectData(v, opt), opt)
		}
	}
	return v
}

func reflectData(data interface{}, opt *ConvOptions) interface{} {
	return reflectValue(reflect.ValueOf(data), opt)
}

func reflectValue(rv reflect.Value, opt *ConvOptions) (v interface{}) {
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
	}
	return
}

func reflectStruct(rv reflect.Value, opt *ConvOptions) interface{} {
	obj := map[string]interface{}{}

	t := rv.Type()
	if 0 < len(opt.CreateKey) {
		if opt.FullTypePath {
			obj[opt.CreateKey] = t.PkgPath() + "/" + t.Name()
		} else {
			obj[opt.CreateKey] = t.Name()
		}
	}
	for i := rv.NumField() - 1; 0 <= i; i-- {
		name := []byte(t.Field(i).Name)
		name[0] = name[0] | 0x20
		obj[string(name)] = From(rv.Field(i).Interface(), opt)
	}
	return obj
}

func reflectComplex(rv reflect.Value, opt *ConvOptions) interface{} {
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

func reflectMap(rv reflect.Value, opt *ConvOptions) interface{} {
	obj := map[string]interface{}{}
	it := rv.MapRange()
	for it.Next() {
		k := it.Key().Interface()
		if ks, ok := k.(string); ok {
			obj[ks] = From(it.Value().Interface())
		} else {
			obj[fmt.Sprint(k)] = From(it.Value().Interface(), opt)
		}
	}
	return obj
}

func reflectArray(rv reflect.Value, opt *ConvOptions) interface{} {
	size := rv.Len()
	a := make([]interface{}, size)
	for i := size - 1; 0 <= i; i-- {
		a[i] = From(rv.Index(i).Interface(), opt)
	}
	return a
}
