// Copyright (c) 2020, Peter Ohler, All rights reserved.

package alt

import (
	"fmt"
	"math"
	"reflect"
	"time"

	"github.com/ohler55/ojg/gen"
)

// RecomposeFunc should build an object from data in a map returning the
// recomposed object or an error.
type RecomposeFunc func(map[string]interface{}) (interface{}, error)

// Recomposer is used to recompose simple data into structs.
type Recomposer struct {

	// CreateKey identifies the creation key in decomposed objects.
	CreateKey string

	composers map[string]*composer
}

// NewRecomposer creates a new instance. The composers are a map of objects
// expected and functions to recompose them. If no function is provided then
// reflection is used instead.
func NewRecomposer(createKey string, composers map[interface{}]RecomposeFunc) (*Recomposer, error) {
	r := Recomposer{
		CreateKey: createKey,
		composers: map[string]*composer{},
	}
	for v, fun := range composers {
		rt := reflect.TypeOf(v)
		if err := r.registerComposer(rt, fun); err != nil {
			return nil, err
		}
	}
	return &r, nil
}

func (r *Recomposer) registerComposer(rt reflect.Type, fun RecomposeFunc) error {
	if rt.Kind() == reflect.Ptr {
		rt = rt.Elem()
	}
	full := rt.PkgPath() + "/" + rt.Name()
	if _, has := r.composers[full]; has {
		return nil
	}
	// TBD could loosen this up and allow any type as long as a function is provided.
	if rt.Kind() != reflect.Struct {
		return fmt.Errorf("only structs can be recomposed. %s is not a struct type", rt)
	}
	c := composer{
		fun:   fun,
		short: rt.Name(),
		full:  full,
		rtype: rt,
	}
	r.composers[c.short] = &c
	r.composers[c.full] = &c
	for i := rt.NumField() - 1; 0 <= i; i-- {
		f := rt.Field(i)
		ft := f.Type
		switch ft.Kind() {
		case reflect.Array, reflect.Slice, reflect.Map, reflect.Ptr:
			ft = ft.Elem()
		}
		_ = r.registerComposer(ft, nil)
	}
	return nil
}

// Recompose simple data into more complex go types.
func (r *Recomposer) Recompose(v interface{}, tv ...interface{}) (interface{}, error) {
	var rt reflect.Type

	if 0 < len(tv) {
		rt = reflect.TypeOf(tv[0])
		if rt.Kind() != reflect.Slice && rt.Kind() != reflect.Array {
			return nil, fmt.Errorf("only a slice type can be provided as an optional argument")
		}
	}
	result, err := r.recompose(v)
	if err == nil && rt != nil {
		if ra, ok := result.([]interface{}); ok {
			av := reflect.MakeSlice(rt, len(ra), len(ra))
			et := rt.Elem()
			for i, v := range ra {
				vv := reflect.ValueOf(v)
				iv := av.Index(i)
				if vv.Type().ConvertibleTo(et) {
					iv.Set(vv.Convert(et))
				} else {
					return nil, fmt.Errorf("can not convert (%s)%v to a %s", iv.Type(), iv, et)
				}
			}
			result = av.Interface()
		}
	}
	return result, err
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
			if b := r.composers[tn]; b != nil {
				return b.compose(o, r.CreateKey)
			}
		}
		v = o

	case gen.Bool:
		v = bool(tv)
	case gen.Int:
		v = int64(tv)
	case gen.Float:
		v = float64(tv)
	case gen.String:
		v = string(tv)
	case gen.Time:
		v = time.Time(tv)
	case gen.Big:
		v = string(tv)
	case gen.Array:
		a := make([]interface{}, len(tv))
		var err error
		for i, m := range tv {
			if a[i], err = r.recompose(m); err != nil {
				return nil, err
			}
		}
		v = a
	case gen.Object:
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
			if b := r.composers[tn]; b != nil {
				return b.compose(o, r.CreateKey)
			}
		}
		v = o

	default:
		return nil, fmt.Errorf("%T is not a valid simple type", v)
	}
	return v, nil
}
