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
	c.indexes = indexType(c.rtype)
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
func (r *Recomposer) Recompose(v interface{}, tv ...interface{}) (out interface{}, err error) {
	defer func() {
		if rec := recover(); rec != nil {
			if err, _ = rec.(error); err == nil {
				err = fmt.Errorf("%v", rec)
			}
			out = nil
		}
	}()
	if 0 < len(tv) {
		out = tv[0]
		rv := reflect.ValueOf(tv[0])
		switch rv.Kind() {
		case reflect.Array, reflect.Slice:
			rv = reflect.New(rv.Type())
			r.recomp(v, rv)
			out = rv.Elem().Interface()
		case reflect.Map:
			r.recomp(v, rv)
		case reflect.Ptr:
			r.recomp(v, rv)
		default:
			return nil, fmt.Errorf("only a slice, map, or pointer is allowed as an optional argument")
		}
	} else {
		out = r.recompAny(v)
	}
	return
}

func (r *Recomposer) recompAny(v interface{}) interface{} {
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
			a[i] = r.recompAny(m)
		}
		v = a
	case map[string]interface{}:
		if cv := tv[r.CreateKey]; cv != nil {
			tn, _ := cv.(string)
			if c := r.composers[tn]; c != nil {
				if c.fun != nil {
					val, err := c.fun(tv)
					if err != nil {
						panic(err)
					}
					return val
				}
				rv := reflect.New(c.rtype)
				r.recomp(v, rv)
				return rv.Interface()
			}
		}
		o := map[string]interface{}{}
		for k, m := range tv {
			o[k] = r.recompAny(m)
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
		for i, m := range tv {
			a[i] = r.recompAny(m)
		}
		v = a
	case gen.Object:
		if cv := tv[r.CreateKey]; cv != nil {
			gn, _ := cv.(gen.String)
			tn := string(gn)
			if c := r.composers[tn]; c != nil {
				simple, _ := tv.Simplify().(map[string]interface{})
				if c.fun != nil {
					val, err := c.fun(simple)
					if err != nil {
						panic(err)
					}
					return val
				}
				rv := reflect.New(c.rtype)
				r.recomp(simple, rv)
				return rv.Interface()
			}
		}
		o := map[string]interface{}{}
		for k, m := range tv {
			o[k] = r.recompAny(m)
		}
		v = o

	default:
		panic(fmt.Errorf("can not recompose a %T", v))
	}
	return v
}

func (r *Recomposer) recomp(v interface{}, rv reflect.Value) {
	as, _ := rv.Interface().(AttrSetter)
	if rv.Kind() == reflect.Ptr {
		rv = rv.Elem()
	}
	switch rv.Kind() {
	case reflect.Array, reflect.Slice:
		va, ok := (v).([]interface{})
		if !ok {
			vv := reflect.ValueOf(v)
			if vv.Kind() != reflect.Slice {
				panic(fmt.Errorf("can only recompose a %s from a []interface{}, not a %T", rv.Type(), v))
			}
			va = make([]interface{}, vv.Len())
			for i := len(va) - 1; 0 <= i; i-- {
				va[i] = vv.Index(i).Interface()
			}
		}
		size := len(va)
		av := reflect.MakeSlice(rv.Type(), size, size)
		et := av.Type().Elem()
		if et.Kind() == reflect.Ptr {
			et = et.Elem()
			for i := 0; i < size; i++ {
				ev := reflect.New(et)
				r.recomp(va[i], ev)
				av.Index(i).Set(ev)
			}
		} else {
			for i := 0; i < size; i++ {
				r.setValue(va[i], av.Index(i))
			}
		}
		rv.Set(av)
	case reflect.Map:
		et := rv.Type().Elem()
		vm, ok := (v).(map[string]interface{})
		if !ok {
			vv := reflect.ValueOf(v)
			if vv.Kind() != reflect.Map {
				panic(fmt.Errorf("can only recompose a map from a map[string]interface{}, not a %T", v))
			}
			vm = map[string]interface{}{}
			iter := vv.MapRange()
			for iter.Next() {
				k := iter.Key().Interface().(string)
				vm[k] = iter.Value().Interface()
			}
		}
		if rv.IsNil() {
			rv.Set(reflect.MakeMap(rv.Type()))
		}
		if et.Kind() == reflect.Interface {
			for k, m := range vm {
				rv.SetMapIndex(reflect.ValueOf(k), reflect.ValueOf(r.recompAny(m)))
			}
		} else if et.Kind() == reflect.Ptr {
			et = et.Elem()
			for k, m := range vm {
				ev := reflect.New(et)
				r.recomp(m, ev)
				rv.SetMapIndex(reflect.ValueOf(k), ev)
			}
		} else {
			for k, m := range vm {
				ev := reflect.New(et)
				r.recomp(m, ev)
				rv.SetMapIndex(reflect.ValueOf(k), ev)
			}
		}
	case reflect.Struct:
		vm, ok := (v).(map[string]interface{})
		if !ok {
			vv := reflect.ValueOf(v)
			if vv.Kind() != reflect.Map {
				panic(fmt.Errorf("can only recompose a %s from a map[string]interface{}, not a %T", rv.Type(), v))
			}
			vm = map[string]interface{}{}
			iter := vv.MapRange()
			for iter.Next() {
				k := iter.Key().Interface().(string)
				vm[k] = iter.Value().Interface()
			}
		}
		if as != nil {
			for k, m := range vm {
				if r.CreateKey == k {
					continue
				}
				if err := as.SetAttr(k, m); err != nil {
					panic(err)
				}
			}
			return
		}
		var im map[string]reflect.StructField
		if c := r.composers[rv.Type().Name()]; c != nil {
			im = c.indexes
		} else {
			im = indexType(rv.Type())
		}
		for k, sf := range im {
			f := rv.FieldByIndex(sf.Index)
			if m, has := vm[k]; has {
				r.setValue(m, f)
			} else if m, has = vm[sf.Name]; has {
				r.setValue(m, f)
			} else {
				name := []byte(sf.Name)
				name[0] = name[0] | 0x20
				if m, has = vm[string(name)]; has {
					r.setValue(m, f)
				}
			}
		}
	case reflect.Interface:
		v = r.recompAny(v)
		rv.Set(reflect.ValueOf(v))
	default:
		panic(fmt.Errorf("can not convert (%T)%v to a %s", v, v, rv.Type()))
	}
}

func (r *Recomposer) setValue(v interface{}, rv reflect.Value) {
	if !rv.IsValid() || !rv.CanSet() {
		return
	}
	switch rv.Kind() {
	case reflect.Bool, reflect.String,
		reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
		reflect.Float32, reflect.Float64:

		rv.Set(reflect.ValueOf(v).Convert(rv.Type()))
	case reflect.Interface:
		v = r.recompAny(v)
		rv.Set(reflect.ValueOf(v))
	default:
		r.recomp(v, rv)
	}
}
