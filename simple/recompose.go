// Copyright (c) 2020, Peter Ohler, All rights reserved.

package simple

import (
	"fmt"
	"math"
	"reflect"
	"time"
)

type RecomposeFunc func(map[string]interface{}) (interface{}, error)

type builder struct {
	build RecomposeFunc
	short string
	full  string
	rtype reflect.Type
}

type Recomposer struct {
	CreateKey string
	builders  map[string]*builder
}

func NewRecomposer(createKey string, builders map[interface{}]RecomposeFunc) (*Recomposer, error) {
	r := Recomposer{
		CreateKey: createKey,
		builders:  map[string]*builder{},
	}
	for v, fun := range builders {
		rt := reflect.TypeOf(v)
		if rt.Kind() == reflect.Ptr {
			rt = rt.Elem()
		}
		b := builder{
			build: fun,
			short: rt.Name(),
			full:  rt.PkgPath() + "/" + rt.Name(),
			rtype: rt,
		}
		r.builders[b.short] = &b
		r.builders[b.full] = &b
	}
	return &r, nil
}

//
func (r *Recomposer) Recompose(v interface{}, tv ...interface{}) (interface{}, error) {
	// TBD only use option for list and map types
	/*
		var xt reflect.Type
		if 0 < len(tv) && tv[0] != nil {
			xt = reflect.TypeOf(tv[0])
		}
	*/
	return r.recompose(v)
}

func (r *Recomposer) recompose(v interface{}) (interface{}, error) {
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
		var err error
		for i, m := range tv {
			if a[i], err = r.recompose(m); err != nil {
				return nil, err
			}
		}
		v = a
	case map[string]interface{}:
		o := map[string]interface{}{}
		for k, m := range tv {
			if mv, err := r.recompose(m); err == nil {
				o[k] = mv
			} else {
				return nil, err
			}
		}
		if cv := o[r.CreateKey]; cv != nil {
			tn, _ := cv.(string)
			if b := r.builders[tn]; b != nil {
				if b.build != nil {
					return b.build(o)
				}
				// TBD use reflection
				//  handle embedded as well
			}
		}
		v = o
	default:
		return nil, fmt.Errorf("%T is not a valid simple type", v)
	}
	return v, nil
}
