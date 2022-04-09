// Copyright (c) 2021, Peter Ohler, All rights reserved.

package asm

import (
	"fmt"
	"time"
)

func init() {
	Define(&Fn{
		Name: "time?",
		Eval: timeCheck,
		Desc: `Returns true if the single required argumement is a time
otherwise false is returned.`,
	})
	Define(&Fn{
		Name: "time",
		Eval: timeConv,
		Desc: `Converts the first argument to a time if possible otherwise
an error is raised. The first argument can be a integer, float,
or string and are converted as follows:
  integer < 10^10:  time in seconds since 1970-01-01 UTC
  integer >= 10^10: time in nanoseconds 1970-01-01 UTC
  decimal (float):  time in seconds 1970-01-01 UTC
  string:           assumed to be formated as RFC3339 unless a
                    format argument is provided`,
	})
}

func timeCheck(root map[string]interface{}, at interface{}, args ...interface{}) interface{} {
	if len(args) != 1 {
		panic(fmt.Errorf("time? expects exactly one arguments. %d given", len(args)))
	}
	_, ok := evalArg(root, at, args[0]).(time.Time)

	return ok
}

func timeConv(root map[string]interface{}, at interface{}, args ...interface{}) (t interface{}) {
	if len(args) < 1 || 2 < len(args) {
		panic(fmt.Errorf("time expects one or two arguments. %d given", len(args)))
	}
	switch v := evalArg(root, at, args[0]).(type) {
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
		i, _ := asInt(v)
		if i < 10000000000 {
			t = time.Unix(i, 0).UTC()
		} else {
			t = time.Unix(0, i).UTC()
		}
	case float32, float64:
		f, _ := asFloat(v)
		sec := int64(f)
		nano := int64((f - float64(sec)) * 1000000000.0)
		t = time.Unix(sec, nano).UTC()
	case string:
		layout := time.RFC3339Nano
		if 1 < len(args) {
			v2 := evalArg(root, at, args[1])
			if s, ok := v2.(string); ok {
				layout = s
			} else {
				panic(fmt.Errorf("time format must be a string, not a %T", v2))
			}
		}
		var err error
		if t, err = time.Parse(layout, v); err != nil {
			panic(err)
		}
	}
	return
}
