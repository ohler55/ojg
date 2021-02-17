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
			if i == 0 {
				ip = ii
			} else if isFloat {
				fp *= float64(ii)
			} else {
				ip *= ii
			}
		case float32, float64:
			f, _ := asFloat(v)
			if i == 0 {
				fp = f
				isFloat = true
			} else if isFloat {
				fp *= f
			} else {
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
