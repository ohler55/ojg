// Copyright (c) 2020, Peter Ohler, All rights reserved.

package gd

import (
	"fmt"
	"time"
	"unsafe"
)

func FromNative(v interface{}) (n Node, err error) {
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
				if a[i], err = FromNative(m); err != nil {
					break
				}
			}
			n = a
		case map[string]interface{}:
			o := Object{}
			for k, m := range tv {
				if o[k], err = FromNative(m); err != nil {
					break
				}
			}
			n = o
		default:
			err = fmt.Errorf("can not convert a %T to a Node", v)
		}
	}
	return nil, nil
}

func AlterNative(v interface{}) (n Node, err error) {
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
				if a[i], err = AlterNative(m); err != nil {
					break
				}
			}
			n = a
		case Array:
			n = tv
		case map[string]interface{}:
			o := *(*Object)(unsafe.Pointer(&tv))
			for k, m := range tv {
				if o[k], err = AlterNative(m); err != nil {
					break
				}
			}
			n = o
		case Object:
			n = tv
		default:
			err = fmt.Errorf("can not convert a %T to a Node", v)
		}
	}
	return
}
