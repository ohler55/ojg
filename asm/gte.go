// Copyright (c) 2021, Peter Ohler, All rights reserved.

package asm

import "fmt"

func init() {
	Define(&Fn{
		Name: "gte",
		Eval: gte,
		Desc: `Returns true if each argument is greater than or equal to any
subsequent argument. An alias is >=.`,
	})
	Define(&Fn{
		Name: ">=",
		Eval: gte,
		Desc: `Returns true if each argument is greater than or equal to any
subsequent argument. An alias is gte.`,
	})
}

func gte(root map[string]interface{}, at interface{}, args ...interface{}) interface{} {
	answer := true
	if 0 < len(args) {
		switch t0 := args[0].(type) {
		case float32, float64,
			int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
			f0, _ := asFloat(t0)
			for _, arg := range args[1:] {
				v := evalArg(root, at, arg)
				if f, ok := asFloat(v); !ok {
					panic(fmt.Errorf("gte of a number must be another number, not %T", v))
				} else if f0 < f {
					answer = false
					break
				} else {
					f0 = f
				}
			}
		case string:
			for _, arg := range args[1:] {
				v := evalArg(root, at, arg)
				if s, _ := v.(string); t0 < s {
					answer = false
					break
				} else {
					t0 = s
				}
			}
		default:
			panic(fmt.Errorf("gte only applies to ints, floats, and strings, not %T", t0))
		}
	}
	return answer
}
