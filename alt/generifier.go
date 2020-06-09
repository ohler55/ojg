// Copyright (c) 2020, Peter Ohler, All rights reserved.

package alt

import (
	"time"
	"unsafe"

	"github.com/ohler55/ojg/gen"
)

// Genericer is the interface for the Generic() function that converts types
// to generic types.
type Genericer interface {

	// Generic should return a Node that represents the object. Generally this
	// includes the use of a creation key consistent with call to the
	// reflection based Generic() function.
	Generic() gen.Node
}

// Generify converts a value into Node compliant data. A best effort is made
// to convert values that are not simple into generic Nodes.
func Generify(v interface{}, options ...*Options) (n gen.Node) {
	opt := &DefaultOptions
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
func GenAlter(v interface{}, options ...*Options) (n gen.Node) {
	opt := &DefaultOptions
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
