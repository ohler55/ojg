// Copyright (c) 2021, Peter Ohler, All rights reserved.

package asm

import (
	"fmt"
	"time"

	"github.com/ohler55/ojg/sen"
)

func init() {
	Define(&Fn{
		Name: "string?",
		Eval: stringCheck,
		Desc: `Returns true if the single required argumement is a string
otherwise false is returned.`,
	})
	Define(&Fn{
		Name: "string",
		Eval: stringConv,
		Desc: `Converts a value into a string.`,
	})
}

func stringCheck(root map[string]any, at any, args ...any) any {
	if len(args) != 1 {
		panic(fmt.Errorf("string? expects exactly one argument. %d given", len(args)))
	}
	_, ok := evalArg(root, at, args[0]).(string)

	return ok
}

func stringConv(root map[string]any, at any, args ...any) (s any) {
	if len(args) < 1 || 2 < len(args) {
		panic(fmt.Errorf("string? expects one or two arguments. %d given", len(args)))
	}
	format := ""
	if 1 < len(args) {
		if x, _ := evalArg(root, at, args[1]).(string); 0 < len(x) {
			format = x
		} else {
			panic("string format argument must be a string")
		}
	}
	switch v := evalArg(root, at, args[0]).(type) {
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
		i, _ := asInt(v)
		if len(format) == 0 {
			format = "%d"
		}
		s = fmt.Sprintf(format, i)
	case float32, float64:
		f, _ := asFloat(v)
		if len(format) == 0 {
			format = "%g"
		}
		s = fmt.Sprintf(format, f)
	case string:
		s = v
	case time.Time:
		if len(format) == 0 {
			format = time.RFC3339Nano
		}
		s = v.Format(format)
	case []any, map[string]any:
		s = sen.String(v, &sen.Options{Sort: true})
	default:
		s = fmt.Sprintf("%v", v)
	}
	return
}
