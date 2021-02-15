// Copyright (c) 2021, Peter Ohler, All rights reserved.

package asm

import (
	"fmt"
	"strconv"
	"time"
)

func init() {
	Define(&Fn{
		Name: "int",
		Eval: intEval,
		Desc: `Converts a value into a integer if possible. I no conversion is
possible nil is returned.`,
	})
}

func intEval(root map[string]interface{}, at interface{}, args ...interface{}) (i interface{}) {
	if len(args) != 1 {
		panic(fmt.Errorf("int expects exactly one argument. %d given", len(args)))
	}
	switch v := evalArg(root, at, args[0]).(type) {
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
		i, _ = asInt(v)
	case float32, float64:
		f, _ := asFloat(v)
		i = int64(f)
	case string:
		var err error
		if i, err = strconv.ParseInt(v, 10, 64); err != nil {
			i = nil
		}
	case time.Time:
		i = v.UnixNano()
	}
	return
}
