// Copyright (c) 2020, Peter Ohler, All rights reserved.

package conv

import (
	"fmt"
	"math"
	"reflect"
	"time"
	"unsafe"

	"github.com/ohler55/ojg/gen"
)

// 23 for fraction in IEEE 754 which amounts to 7 significant digits. Use base
// 10 so that numbers look correct when displayed in base 10.
const fracMax = 10000000.0

const (
	Exact = Status(iota)
	Ok
	Fail
)

type Status int

func (status Status) String() (s string) {
	switch status {
	case Exact:
		s = "Exact"
	case Ok:
		s = "Ok"
	case Fail:
		s = "Fail"
	}
	return
}

// ConvOptions are the options available to Decompose() function.
type ConvOptions struct {

	// CreateKey is the map element used to identify the type of a decomposed
	// object.
	CreateKey string

	// FullTypePath if true will use the package and type name as the
	// CreateKey value.
	FullTypePath bool

	// OmitNil if true omits object members that have nil values.
	OmitNil bool
}

// DefaultConvOptions are the default options for decompsing.
var DefaultConvOptions = ConvOptions{
	CreateKey:    "type",
	FullTypePath: false,
	OmitNil:      true,
}

// Decompose creates a simple type converting non simple to simple types using
// either the Simplify() interface or reflection. Unlike Alter() a deep copy
// is returned leaving the original data unchanged.
func Decompose(v interface{}, options ...*ConvOptions) interface{} {
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
				a[i] = Decompose(m)
			}
			v = a
		case map[string]interface{}:
			o := map[string]interface{}{}
			for k, m := range tv {
				if mv := Decompose(m); mv != nil || !opt.OmitNil {
					o[k] = mv
				}
			}
			v = o
		default:
			if simp, _ := v.(Simplifier); simp != nil {
				return Decompose(simp.Simplify())
			}
			return Decompose(reflectData(v, opt))
		}
	}
	return v
}

// Alter the data into all simple types converting non simple to simple types
// using either the Simplify() interface or reflection. Unlike Decompose() map and
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
		obj[string(name)] = Decompose(rv.Field(i).Interface(), opt)
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
			obj[ks] = Decompose(it.Value().Interface())
		} else {
			obj[fmt.Sprint(k)] = Decompose(it.Value().Interface(), opt)
		}
	}
	return obj
}

func reflectArray(rv reflect.Value, opt *ConvOptions) interface{} {
	size := rv.Len()
	a := make([]interface{}, size)
	for i := size - 1; 0 <= i; i-- {
		a[i] = Decompose(rv.Index(i).Interface(), opt)
	}
	return a
}

// Generify converts a value into Node compliant data. A best effort is made
// to convert values that are not simple into generic Nodes.
func Generify(v interface{}, options ...*ConvOptions) (n gen.Node) {
	opt := &DefaultConvOptions
	if 0 < len(options) {
		opt = options[0]
	}
	if v != nil {
		switch tv := v.(type) {
		case bool:
			n = gen.Bool(tv)
		case gen.Bool:
			n = tv
		case int:
			n = gen.Int(int64(tv))
		case int8:
			n = gen.Int(int64(tv))
		case int16:
			n = gen.Int(int64(tv))
		case int32:
			n = gen.Int(int64(tv))
		case int64:
			n = gen.Int(tv)
		case uint:
			n = gen.Int(int64(tv))
		case uint8:
			n = gen.Int(int64(tv))
		case uint16:
			n = gen.Int(int64(tv))
		case uint32:
			n = gen.Int(int64(tv))
		case uint64:
			n = gen.Int(int64(tv))
		case gen.Int:
			n = tv
		case float32:
			n = gen.Float(float64(tv))
		case float64:
			n = gen.Float(tv)
		case gen.Float:
			n = tv
		case string:
			n = gen.String(tv)
		case gen.String:
			n = tv
		case time.Time:
			n = gen.Time(tv)
		case gen.Time:
			n = tv
		case []interface{}:
			a := make(gen.Array, len(tv))
			for i, m := range tv {
				a[i] = Generify(m, opt)
			}
			n = a
		case map[string]interface{}:
			o := gen.Object{}
			for k, m := range tv {
				o[k] = Generify(m, opt)
			}
			n = o
		default:
			if g, _ := n.(Genericer); g != nil {
				return g.Generic()
			}
			if simp, _ := n.(Simplifier); simp != nil {
				return Generify(simp.Simplify(), opt)
			}
			return Generify(reflectData(v, opt))
		}
	}
	return
}

// GenAlter converts a simple go data element into Node compliant data. A best
// effort is made to convert values that are not simple into generic Nodes. It
// modifies the values inplace if possible by altering the original.
func GenAlter(v interface{}, options ...*ConvOptions) (n gen.Node) {
	opt := &DefaultConvOptions
	if 0 < len(options) {
		opt = options[0]
	}
	if v != nil {
		switch tv := v.(type) {
		case bool:
			n = gen.Bool(tv)
		case gen.Bool:
			n = tv
		case int:
			n = gen.Int(int64(tv))
		case int8:
			n = gen.Int(int64(tv))
		case int16:
			n = gen.Int(int64(tv))
		case int32:
			n = gen.Int(int64(tv))
		case int64:
			n = gen.Int(tv)
		case uint:
			n = gen.Int(int64(tv))
		case uint8:
			n = gen.Int(int64(tv))
		case uint16:
			n = gen.Int(int64(tv))
		case uint32:
			n = gen.Int(int64(tv))
		case uint64:
			n = gen.Int(int64(tv))
		case gen.Int:
			n = tv
		case float32:
			n = gen.Float(float64(tv))
		case float64:
			n = gen.Float(tv)
		case gen.Float:
			n = tv
		case string:
			n = gen.String(tv)
		case gen.String:
			n = tv
		case time.Time:
			n = gen.Time(tv)
		case gen.Time:
			n = tv
		case []interface{}:
			a := *(*gen.Array)(unsafe.Pointer(&tv))
			for i, m := range tv {
				a[i] = GenAlter(m)
			}
			n = a
		case gen.Array:
			n = tv
		case map[string]interface{}:
			o := *(*gen.Object)(unsafe.Pointer(&tv))
			for k, m := range tv {
				o[k] = GenAlter(m, opt)
			}
			n = o
		case gen.Object:
			n = tv
		default:
			if g, _ := n.(Genericer); g != nil {
				return g.Generic()
			}
			if simp, _ := n.(Simplifier); simp != nil {
				return GenAlter(simp.Simplify(), opt)
			}
			return GenAlter(reflectData(v, opt))
		}
	}
	return
}
