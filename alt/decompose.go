// Copyright (c) 2020, Peter Ohler, All rights reserved.

package alt

import (
	"fmt"
	"math"
	"reflect"
	"strings"
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
			a[i] = Decompose(m)
		}
		v = a
	case map[string]interface{}:
		o := map[string]interface{}{}
		for k, m := range tv {
			if mv := Decompose(m); mv != nil || !opt.OmitNil {
				if mv != nil || !opt.OmitNil {
					o[k] = mv
				}
			}
		}
		v = o
	default:
		if simp, _ := v.(Simplifier); simp != nil {
			return Decompose(simp.Simplify())
		}
		return Decompose(reflectData(v, opt))
	}
	if opt.Converter != nil {
		v, _ = opt.Converter.convert(v)
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
				if mv != nil || !opt.OmitNil {
					tv[k] = mv
				}
			}
		}
	default:
		if simp, _ := v.(Simplifier); simp != nil {
			return Alter(simp.Simplify())
		}
		return Alter(reflectData(v, opt), opt)
	}
	if opt.Converter != nil {
		v, _ = opt.Converter.convert(v)
	}
	return v
}

func reflectData(data interface{}, opt *Options) interface{} {
	return reflectValue(reflect.ValueOf(data), opt)
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
	if opt.NestEmbed {
		for i := rv.NumField() - 1; 0 <= i; i-- {
			reflectStructMap(obj, rv.Field(i), t.Field(i), opt)
		}
	} else {
		im := allFields(t)
		for _, sf := range im {
			fv := rv.FieldByIndex(sf.Index)
			reflectStructMap(obj, fv, sf, opt)
		}
	}
	return obj
}

func allFields(rt reflect.Type) (im []reflect.StructField) {
	for i := rt.NumField() - 1; 0 <= i; i-- {
		f := rt.Field(i)
		if f.Anonymous {
			// prepend index and add to im
			for _, ff := range allFields(f.Type) {
				ff.Index = append([]int{i}, ff.Index...)
				im = append(im, ff)
			}
		} else {
			im = append(im, f)
		}
	}
	return
}

func reflectStructMap(obj map[string]interface{}, rv reflect.Value, f reflect.StructField, opt *Options) {
	name := []byte(f.Name)
	if len(name) == 0 || 'a' <= name[0] {
		// not a public field
		return
	}
	var g interface{}
	if !isNil(rv) {
		g = Decompose(rv.Interface(), opt)
	}

	if !opt.KeyExact {
		name[0] = name[0] | 0x20
	}
	if opt.UseTags {
		if tag, ok := f.Tag.Lookup("json"); ok && 0 < len(tag) {
			parts := strings.Split(tag, ",")
			switch parts[0] {
			case "":
				name = []byte(f.Name)
			case "-":
				if 1 < len(parts) {
					name = []byte{'-'}
				} else {
					// skip
					return
				}
			default:
				name = []byte(parts[0])
			}
			for _, p := range parts[1:] {
				switch p {
				case "omitempty":
					if g == nil || reflect.ValueOf(g).IsZero() {
						return
					}
				case "string":
					g = fmt.Sprintf("%v", g)
				}
			}
		}
	}
	if g != nil || !opt.OmitNil {
		obj[string(name)] = g
	}
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
			g = Decompose(vv.Interface(), opt)
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
		a[i] = Decompose(rv.Index(i).Interface(), opt)
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
