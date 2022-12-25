// Copyright (c) 2021, Peter Ohler, All rights reserved.

package asm

import (
	"fmt"
	"strconv"
	"time"
)

func init() {
	Define(&Fn{
		Name: "float",
		Eval: floatEval,
		Desc: `Converts a value into a float if possible. I no conversion is
possible nil is returned.`,
	})
}

func floatEval(root map[string]any, at any, args ...any) (f any) {
	if len(args) != 1 {
		panic(fmt.Errorf("float expects exactly one argument. %d given", len(args)))
	}
	switch v := evalArg(root, at, args[0]).(type) {
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64,
		float32, float64:
		f, _ = asFloat(v)
	case string:
		var err error
		if f, err = strconv.ParseFloat(v, 64); err != nil {
			f = nil
		}
	case time.Time:
		f = float64(v.UnixNano()) / 1000000000.0
	}
	return
}
