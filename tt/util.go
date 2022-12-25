// Copyright (c) 2020, Peter Ohler, All rights reserved.

package tt

import (
	"encoding/json"
	"fmt"
	"runtime"
	"strings"
	"testing"
	"unsafe"

	"github.com/ohler55/ojg/gen"
)

type call struct {
	fn   string
	file string
	line int
}

func finishFail(t *testing.T, b *strings.Builder, args []any) {
	stackFill(b)
	if 0 < len(args) {
		if format, _ := args[0].(string); 0 < len(format) {
			b.WriteString(fmt.Sprintf(format, args[1:]...))
		} else {
			b.WriteString(fmt.Sprint(args...))
		}
	}
	t.Fatal(b.String())
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

func isNil(v any) bool {
	return (*[2]uintptr)(unsafe.Pointer(&v))[1] == 0
}

func asInt(v any) (i int64, ok bool) {
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
	case gen.Int:
		i = int64(tv)
	default:
		ok = false
	}
	return
}

func asFloat(v any) (f float64, ok bool) {
	ok = true
	switch tv := v.(type) {
	case float32:
		f = float64(tv)
	case float64:
		f = tv
	case gen.Float:
		f = float64(tv)
	default:
		ok = false
	}
	return
}

func asString(v any) (s string, ok bool) {
	ok = true
	switch tv := v.(type) {
	case string:
		s = tv
	case gen.String:
		s = string(tv)
	case json.Number:
		s = string(tv)
	default:
		ok = false
	}
	return
}
