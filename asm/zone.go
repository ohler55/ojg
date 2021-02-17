// Copyright (c) 2021, Peter Ohler, All rights reserved.

package asm

import (
	"fmt"
	"time"
)

func init() {
	Define(&Fn{
		Name: "zone",
		Eval: zone,
		Desc: `Changes the timezone on a time to the location specified in the
second argument. Raises an error if the first argument does not
evaluate to a time or the location can not be determined.
Location can be either a string or the number of minutes offset
from UTC.`,
	})
}

func zone(root map[string]interface{}, at interface{}, args ...interface{}) interface{} {
	if len(args) != 2 {
		panic(fmt.Errorf("zone expects exactly two arguments. %d given", len(args)))
	}
	v := evalArg(root, at, args[0])
	t, ok := v.(time.Time)
	if !ok {
		panic(fmt.Errorf("zone requires a time argument, not a %T", v))
	}
	var loc *time.Location
	switch v := evalArg(root, at, args[1]).(type) {
	case string:
		var err error
		if loc, err = time.LoadLocation(v); err != nil {
			loc = time.FixedZone(v, 0)
		}
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
		i, _ := asInt(v)
		loc = time.FixedZone("", int(i))
	default:
		panic(fmt.Errorf("zone location must be a string or number, not a %T", v))
	}
	t = t.In(loc)

	return t
}
