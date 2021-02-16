// Copyright (c) 2021, Peter Ohler, All rights reserved.

package asm

import (
	"fmt"
)

func init() {
	Define(&Fn{
		Name: "substr",
		Eval: substr,
		Desc: `Returns a substring of the input string. The second argument
must be an integer that marks the start of the substring. The
third integer argument indicates the length of the substring
if provided. If the length argument is not provided the end of
the substring is the end of the input string.`,
	})
}

func substr(root map[string]interface{}, at interface{}, args ...interface{}) interface{} {
	if len(args) < 1 || 3 < len(args) {
		panic(fmt.Errorf("substr expects two or three arguments. %d given", len(args)))
	}
	v := evalArg(root, at, args[0])
	s, ok := v.(string)
	if !ok {
		panic(fmt.Errorf("substr expects a string argument, not a %T", v))
	}
	v = evalArg(root, at, args[1])
	var start int64
	if start, ok = asInt(v); !ok {
		panic(fmt.Errorf("substr expects an integer second argument, not a %T", v))
	}
	if start < 0 {
		start = int64(len(s)) + start
		if start < 0 {
			start = 0
		}
	}
	var count int64
	if 2 < len(args) {
		v = evalArg(root, at, args[2])
		if count, ok = asInt(v); !ok {
			panic(fmt.Errorf("substr expects an integer third argument, not a %T", v))
		}
		if count < 0 {
			return ""
		}
		if int64(len(s)) < start+count {
			s = s[start:]
		} else {
			s = s[start : start+count]
		}
	} else {
		s = s[start:]
	}
	return s
}
