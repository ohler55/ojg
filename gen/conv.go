// Copyright (c) 2020, Peter Ohler, All rights reserved.

package gen

import (
	"time"
	"unsafe"

	"github.com/ohler55/ojg/simple"
)

func From(v interface{}) (n Node) {
	if v != nil {
		switch tv := v.(type) {
		case bool:
			n = Bool(tv)
		case Bool:
			n = tv
		case int:
			n = Int(int64(tv))
		case int8:
			n = Int(int64(tv))
		case int16:
			n = Int(int64(tv))
		case int32:
			n = Int(int64(tv))
		case int64:
			n = Int(tv)
		case uint:
			n = Int(int64(tv))
		case uint8:
			n = Int(int64(tv))
		case uint16:
			n = Int(int64(tv))
		case uint32:
			n = Int(int64(tv))
		case uint64:
			n = Int(int64(tv))
		case Int:
			n = tv
		case float32:
			n = Float(float64(tv))
		case float64:
			n = Float(tv)
		case Float:
			n = tv
		case string:
			n = String(tv)
		case String:
			n = tv
		case time.Time:
			n = Time(tv)
		case Time:
			n = tv
		case []interface{}:
			a := make(Array, len(tv))
			for i, m := range tv {
				a[i] = From(m)
			}
			n = a
		case map[string]interface{}:
			o := Object{}
			for k, m := range tv {
				o[k] = From(m)
			}
			n = o
		default:
			if g, _ := n.(Genericer); g != nil {
				return g.Generic()
			}
			if simp, _ := n.(simple.Simplifier); simp != nil {
				return From(simp.Simplify())
			}
			// TBD from

			// TBD always succeed with something
			// err = fmt.Errorf("can not convert a %T to a Node", v)
		}
	}
	return
}

func Alter(v interface{}) (n Node) {
	if v != nil {
		switch tv := v.(type) {
		case bool:
			n = Bool(tv)
		case Bool:
			n = tv
		case int:
			n = Int(int64(tv))
		case int8:
			n = Int(int64(tv))
		case int16:
			n = Int(int64(tv))
		case int32:
			n = Int(int64(tv))
		case int64:
			n = Int(tv)
		case uint:
			n = Int(int64(tv))
		case uint8:
			n = Int(int64(tv))
		case uint16:
			n = Int(int64(tv))
		case uint32:
			n = Int(int64(tv))
		case uint64:
			n = Int(int64(tv))
		case Int:
			n = tv
		case float32:
			n = Float(float64(tv))
		case float64:
			n = Float(tv)
		case Float:
			n = tv
		case string:
			n = String(tv)
		case String:
			n = tv
		case time.Time:
			n = Time(tv)
		case Time:
			n = tv
		case []interface{}:
			a := *(*Array)(unsafe.Pointer(&tv))
			for i, m := range tv {
				a[i] = Alter(m)
			}
			n = a
		case Array:
			n = tv
		case map[string]interface{}:
			o := *(*Object)(unsafe.Pointer(&tv))
			for k, m := range tv {
				o[k] = Alter(m)
			}
			n = o
		case Object:
			n = tv
		default:
			if g, _ := n.(Genericer); g != nil {
				return g.Generic()
			}
			if simp, _ := n.(simple.Simplifier); simp != nil {
				return From(simp.Simplify())
			}
			// TBD
			//err = fmt.Errorf("can not convert a %T to a Node", v)
		}
	}
	return
}
