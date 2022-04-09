// Copyright (c) 2021, Peter Ohler, All rights reserved.

package asm

import (
	"fmt"
)

func init() {
	Define(&Fn{
		Name: "product",
		Eval: product,
		Desc: `Returns the product of all arguments. All arguments must be
numbers. If any of the arguments are not a number an error is
raised.`,
	})
	Define(&Fn{
		Name: "*",
		Eval: product,
		Desc: `Returns the product of all arguments. All arguments must be
numbers. If any of the arguments are not a number an error is
raised.`,
	})
}

func product(root map[string]interface{}, at interface{}, args ...interface{}) interface{} {
	var ip int64
	var fp float64
	isFloat := false
	for i, arg := range args {
		switch v := evalArg(root, at, arg).(type) {
		case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
			ii, _ := asInt(v)
			switch {
			case i == 0:
				ip = ii
			case isFloat:
				fp *= float64(ii)
			default:
				ip *= ii
			}
		case float32, float64:
			f, _ := asFloat(v)
			switch {
			case i == 0:
				fp = f
				isFloat = true
			case isFloat:
				fp *= f
			default:
				isFloat = true
				fp = float64(ip) * f
			}
		default:
			panic(fmt.Errorf("a %T argument can not be multiplied", v))
		}
	}
	if isFloat {
		return fp
	}
	return ip
}
