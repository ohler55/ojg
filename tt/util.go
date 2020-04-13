// Copyright (c) 2020, Peter Ohler, All rights reserved.

package tt

import (
	"fmt"
	"runtime"
	"strings"
	"unsafe"

	"github.com/ohler55/ojg/gd"
)

type call struct {
	fn   string
	file string
	line int
}

func stackFill(b *strings.Builder) {
	pc := make([]uintptr, 40)
	cnt := runtime.Callers(3, pc) - 2
	stack := make([]call, cnt)

	var fn *runtime.Func
	var c *call

	for i := 0; i < cnt; i++ {
		c = &stack[i]
		fn = runtime.FuncForPC(pc[i])
		c.file, c.line = fn.FileLine(pc[i])
		c.fn = fn.Name()
		b.WriteString(fmt.Sprintf("%s @ %s:%d", c.fn, c.file, c.line))
		b.WriteByte('\n')
	}
}

func isNil(v interface{}) bool {
	return (*[2]uintptr)(unsafe.Pointer(&v))[1] == 0
}

func asInt(v interface{}) (i int64, ok bool) {
	ok = true
	switch tv := v.(type) {
	case int:
		i = int64(tv)
	case int8:
		i = int64(tv)
	case int16:
		i = int64(tv)
	case int32:
		i = int64(tv)
	case int64:
		i = tv
	case uint:
		i = int64(tv)
	case uint8:
		i = int64(tv)
	case uint16:
		i = int64(tv)
	case uint32:
		i = int64(tv)
	case uint64:
		i = int64(tv)
	case gd.Int:
		i = int64(tv)
	default:
		ok = false
	}
	return
}

func asFloat(v interface{}) (f float64, ok bool) {
	ok = true
	switch tv := v.(type) {
	case float32:
		f = float64(tv)
	case float64:
		f = tv
	case gd.Float:
		f = float64(tv)
	default:
		ok = false
	}
	return
}

func asString(v interface{}) (s string, ok bool) {
	ok = true
	switch tv := v.(type) {
	case string:
		s = tv
	case gd.String:
		s = string(tv)
	default:
		ok = false
	}
	return
}
