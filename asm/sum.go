// Copyright (c) 2021, Peter Ohler, All rights reserved.

package asm

import (
	"fmt"
)

func init() {
	Define(&Fn{
		Name: "sum",
		Eval: sum,
		Desc: `Returns the sum of all arguments. All arguments must be numbers
or strings. If any argument is a string then the result will be
a string otherwise the result will be a number. If any of the
arguments are not a number or a string an error is raised.`,
	})
	Define(&Fn{
		Name: "+",
		Eval: sum,
		Desc: `Returns the sum of all arguments. All arguments must be numbers
or strings. If any argument is a string then the result will be
a string otherwise the result will be a number. If any of the
arguments are not a number or a string an error is raised.`,
	})
}

const (
	intSum = iota
	floatSum
	strSum
)

func sum(root map[string]interface{}, at interface{}, args ...interface{}) interface{} {
	kind := intSum
	var ssum string
	var isum int64
	var fsum float64
	for i, arg := range args {
		switch v := evalArg(root, at, arg).(type) {
		case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
			ii, _ := asInt(v)
			if i == 0 {
				kind = intSum
			}
			switch kind {
			case intSum:
				isum += ii
			case floatSum:
				fsum += float64(ii)
			case strSum:
				ssum = fmt.Sprintf("%s%d", ssum, ii)
			}
		case float32, float64:
			if i == 0 {
				kind = floatSum
			}
			f, _ := asFloat(v)
			switch kind {
			case intSum:
				kind = floatSum
				fsum = float64(isum) + f
			case floatSum:
				fsum += f
			case strSum:
				ssum = fmt.Sprintf("%s%g", ssum, f)
			}
		case string:
			if i == 0 {
				kind = strSum
			}
			switch kind {
			case intSum:
				kind = strSum
				ssum = fmt.Sprintf("%d%s", isum, v)
			case floatSum:
				kind = strSum
				ssum = fmt.Sprintf("%g%s", fsum, v)
			case strSum:
				ssum += v
			}
		default:
			panic(fmt.Errorf("a %T argument can not be summed", v))
		}
	}
	switch kind {
	case intSum:
		return isum
	case floatSum:
		return fsum
	}
	return ssum
}
