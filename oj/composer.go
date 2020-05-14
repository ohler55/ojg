// Copyright (c) 2020, Peter Ohler, All rights reserved.

package oj

import (
	"fmt"
	"reflect"
	"strings"
)

// RecomposeFunc should build an object from data in a map returning the
// recomposed object or an error.
type RecomposeFunc func(map[string]interface{}) (interface{}, error)

type composer struct {
	fun   RecomposeFunc
	short string
	full  string
	rtype reflect.Type
}

func newComposer(v interface{}, fun RecomposeFunc) (*composer, error) {
	rt := reflect.TypeOf(v)
	if rt.Kind() == reflect.Ptr {
		rt = rt.Elem()
	}
	if rt.Kind() != reflect.Struct {
		return nil, fmt.Errorf("only structs can be recomposed. %s is not a struct type", rt)
	}
	c := composer{
		fun:   fun,
		short: rt.Name(),
		full:  rt.PkgPath() + "/" + rt.Name(),
		rtype: rt,
	}
	return &c, nil
}

func (c *composer) compose(obj map[string]interface{}, createKey string) (interface{}, error) {
	if c.fun != nil {
		return c.fun(obj)
	}
	nvp := reflect.New(c.rtype)
	nv := nvp.Elem()
	for key, v := range obj {
		if createKey == key {
			continue
		}
		f, ok := c.rtype.FieldByNameFunc(func(s string) bool { return strings.EqualFold(s, key) })
		if !ok {
			continue
		}
		fv := nv.FieldByIndex(f.Index)
		if fv.CanSet() {
			ft := fv.Type()
			vv := reflect.ValueOf(v)
			if vv.Type().ConvertibleTo(ft) {
				fv.Set(vv.Convert(ft))
			} else if (fv.Kind() == reflect.Slice || fv.Kind() == reflect.Array) &&
				(vv.Kind() == reflect.Slice || vv.Kind() == reflect.Array) {

				size := vv.Len()
				av := reflect.MakeSlice(ft, size, size)
				at := av.Type().Elem()
				for i := 0; i < size; i++ {
					// Index(i) returns interface{} type so get the value then
					// create an Value from that.
					vi := reflect.ValueOf(vv.Index(i).Interface())
					fi := av.Index(i)
					if vi.Type().ConvertibleTo(at) {
						fi.Set(vi.Convert(at))
					} else {
						return nil, fmt.Errorf("can not convert (%s)%v to a %s for field %s", vi.Type(), vi, fi.Type(), f.Name)
					}
				}
				if av.Type().AssignableTo(ft) {
					fv.Set(av)
				}
			} else {
				return nil, fmt.Errorf("can not convert (%T)%v to a %s for field %s", v, v, ft, f.Name)
			}
		}
	}
	return nvp.Interface(), nil
}
