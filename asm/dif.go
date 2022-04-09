// Copyright (c) 2021, Peter Ohler, All rights reserved.

package asm

import (
	"fmt"
)

func init() {
	Define(&Fn{
		Name: "dif",
		Eval: dif,
		Desc: `Returns the difference of all arguments. All arguments must be
numbers. If any of the arguments are not a number an error is
raised.`,
	})
	Define(&Fn{
		Name: "-",
		Eval: dif,
		Desc: `Returns the difference of all arguments. All arguments must be
numbers. If any of the arguments are not a number an error is
raised.`,
	})
}

func dif(root map[string]interface{}, at interface{}, args ...interface{}) interface{} {
	var idif int64
	var fdif float64
	isFloat := false
	for i, arg := range args {
		switch v := evalArg(root, at, arg).(type) {
		case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
			ii, _ := asInt(v)
			switch {
			case i == 0:
				idif = ii
			case isFloat:
				fdif -= float64(ii)
			default:
				idif -= ii
			}
		case float32, float64:
			f, _ := asFloat(v)
			switch {
			case i == 0:
				fdif = f
				isFloat = true
			case isFloat:
				fdif -= f
			default:
				isFloat = true
				fdif = float64(idif) - f
			}
		default:
			panic(fmt.Errorf("a %T argument can not be an argument to dif", v))
		}
	}
	if isFloat {
		return fdif
	}
	return idif
}
