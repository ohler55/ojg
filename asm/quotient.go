// Copyright (c) 2021, Peter Ohler, All rights reserved.

package asm

import (
	"fmt"
)

func init() {
	Define(&Fn{
		Name: "quotient",
		Eval: quotient,
		Desc: `Returns the quotient of all arguments. All arguments must be
numbers. If any of the arguments are not a number an error is
raised. If an attempt is made to divide by zero and error will
be raised.`,
	})
	Define(&Fn{
		Name: "/",
		Eval: quotient,
		Desc: `Returns the quotient of all arguments. All arguments must be
numbers. If any of the arguments are not a number an error is
raised. If an attempt is made to divide by zero and error will
be raised.`,
	})
}

func quotient(root map[string]interface{}, at interface{}, args ...interface{}) interface{} {
	var iq int64
	var fq float64
	isFloat := false
	for i, arg := range args {
		switch v := evalArg(root, at, arg).(type) {
		case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
			ii, _ := asInt(v)
			switch {
			case i == 0:
				iq = ii
			case isFloat:
				fq /= float64(ii)
			default:
				iq /= ii
			}
		case float32, float64:
			f, _ := asFloat(v)
			switch {
			case i == 0:
				fq = f
				isFloat = true
			case isFloat:
				fq /= f
			default:
				isFloat = true
				fq = float64(iq) / f
			}
		default:
			panic(fmt.Errorf("a %T argument can not be an argument to quotient", v))
		}
	}
	if isFloat {
		return fq
	}
	return iq
}
