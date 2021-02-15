// Copyright (c) 2021, Peter Ohler, All rights reserved.

package asm

import (
	"time"
)

func init() {
	Define(&Fn{
		Name: "equal",
		Eval: equal,
		Desc: `Returns true if all the argument are equal. Aliases are eq, ==,
and equal.`,
	})
	Define(&Fn{
		Name: "eq",
		Eval: equal,
		Desc: `Returns true if all the argument are equal. Aliases are eq, ==,
and equal.`,
	})
	Define(&Fn{
		Name: "==",
		Eval: equal,
		Desc: `Returns true if all the argument are equal. Aliases are eq, ==,
and equal.`,
	})
}

func equal(root map[string]interface{}, at interface{}, args ...interface{}) interface{} {
	if 0 < len(args) {
		v0 := evalArg(root, at, args[0])
		for _, v := range args[1:] {
			v = evalArg(root, at, v)
			if !equalVals(v0, v) {
				return false
			}
		}
	}
	return true
}

func equalVals(v0, v1 interface{}) (eq bool) {
	switch t0 := v0.(type) {
	case nil:
		eq = nil == v1
	case bool:
		if b1, ok := v1.(bool); ok {
			eq = b1 == t0
		}
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
		x, _ := asInt(v0)
		a, ok := asInt(v1)
		eq = ok && x == a
	case float32, float64:
		x, _ := asFloat(v0)
		a, ok := asFloat(v1)
		eq = ok && x == a
	case string:
		a, ok := v1.(string)
		eq = v0 == a && ok
	case time.Time:
		tm, _ := v1.(time.Time)
		eq = tm == t0
	case []interface{}:
		if t1, ok := v1.([]interface{}); ok && len(t0) == len(t1) {
			eq = true
			for i, m0 := range t0 {
				if eq = equalVals(m0, t1[i]); !eq {
					break
				}
			}
		}
	case map[string]interface{}:
		if t1, ok := v1.(map[string]interface{}); ok && len(t0) == len(t1) {
			eq = true
			for k, m0 := range t0 {
				m1, has := t1[k]
				if eq = has && equalVals(m0, m1); !eq {
					break
				}
			}
		}
	}
	return
}

func asInt(v interface{}) (i int64, ok bool) {
	ok = true
	switch tv := v.(type) {
	case int:
		i = int64(tv)
	case int8:
		i = int64(tv)
	case int16:
		i = int64(tv)
	case int32:
		i = int64(tv)
	case int64:
		i = tv
	case uint:
		i = int64(tv)
	case uint8:
		i = int64(tv)
	case uint16:
		i = int64(tv)
	case uint32:
		i = int64(tv)
	case uint64:
		i = int64(tv)
	default:
		ok = false
	}
	return
}

func asFloat(v interface{}) (f float64, ok bool) {
	ok = true
	switch tv := v.(type) {
	case float32:
		f = float64(tv)
	case float64:
		f = tv
	case int:
		f = float64(tv)
	case int8:
		f = float64(tv)
	case int16:
		f = float64(tv)
	case int32:
		f = float64(tv)
	case int64:
		f = float64(tv)
	case uint:
		f = float64(tv)
	case uint8:
		f = float64(tv)
	case uint16:
		f = float64(tv)
	case uint32:
		f = float64(tv)
	case uint64:
		f = float64(tv)
	default:
		ok = false
	}
	return
}
