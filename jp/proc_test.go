// Copyright (c) 2024, Peter Ohler, All rights reserved.

package jp_test

import (
	"fmt"
	"testing"

	"github.com/ohler55/ojg/jp"
	"github.com/ohler55/ojg/pretty"
	"github.com/ohler55/ojg/tt"
)

type mathProc struct {
	op    rune
	left  int
	right int
}

func (mp *mathProc) Get(data any) []any {
	return []any{mp.First(data)}
}

func (mp *mathProc) First(data any) any {
	a, ok := data.([]any)
	if ok {
		var (
			left  int
			right int
		)
		if mp.left < len(a) {
			left, _ = a[mp.left].(int)
		}
		if mp.right < len(a) {
			right, _ = a[mp.right].(int)
		}
		switch mp.op {
		case '+':
			return left + right
		case '-':
			return left - right
		}
		return 0
	}
	return nil
}

func compileMathProc(code []byte) jp.Procedure {
	var mp mathProc
	_, _ = fmt.Sscanf(string(code), "(%c %d %d)", &mp.op, &mp.left, &mp.right)

	return &mp
}

type mapProc struct{}

func (mp mapProc) Get(data any) (result []any) {
	a, _ := data.([]any)
	for i, v := range a {
		result = append(result, map[string]any{"i": i, "v": v})
	}
	return
}

func (mp mapProc) First(data any) any {
	if a, _ := data.([]any); 0 < len(a) {
		return map[string]any{"i": 0, "v": a[0]}
	}
	return nil
}

func compileMapProc(code []byte) jp.Procedure {
	return mapProc{}
}

type mapIntProc struct{}

func (mip mapIntProc) Get(data any) (result []any) {
	a, _ := data.([]int)
	for i, v := range a {
		result = append(result, map[string]int{"i": i, "v": v})
	}
	return
}

func (mip mapIntProc) First(data any) any {
	if a, _ := data.([]int); 0 < len(a) {
		return map[string]int{"i": 0, "v": a[0]}
	}
	return nil
}

func compileMapIntProc(code []byte) jp.Procedure {
	return mapIntProc{}
}

func TestProcLast(t *testing.T) {
	jp.CompileScript = compileMathProc

	p := jp.MustNewProc([]byte("(+ 0 1)"))
	tt.Equal(t, "[(+ 0 1)]", p.String())

	x := jp.MustParseString("[(+ 0 1)]")
	tt.Equal(t, "[(+ 0 1)]", x.String())

	data := []any{2, 3, 4}
	result := x.First(data)
	tt.Equal(t, 5, result)

	got := x.Get(data)
	tt.Equal(t, []any{5}, got)

	locs := x.Locate(data, 1)
	tt.Equal(t, "[[0]]", pretty.SEN(locs))

	var buf []byte
	x.Walk(data, func(path jp.Expr, nodes []any) {
		buf = fmt.Appendf(buf, "%s : %v\n", path, nodes)
	})
	tt.Equal(t, "[0] : [[2 3 4] 5]\n", string(buf))
}

func TestProcNotLast(t *testing.T) {
	jp.CompileScript = compileMapProc

	x := jp.MustParseString("[(quux)].v")
	tt.Equal(t, "[(quux)].v", x.String())

	data := []any{2, 3, 4}
	result := x.First(data)
	tt.Equal(t, 2, result)

	got := x.Get(data)
	tt.Equal(t, []any{2, 3, 4}, got)

	locs := x.Locate(data, 2)
	tt.Equal(t, "[[0 v] [1 v]]", pretty.SEN(locs))

	var buf []byte
	x.Walk(data, func(path jp.Expr, nodes []any) {
		buf = fmt.Appendf(buf, "%s : %v\n", path, nodes)
	})
	tt.Equal(t, `[0].v : [[2 3 4] map[i:0 v:2] 2]
[1].v : [[2 3 4] map[i:1 v:3] 3]
[2].v : [[2 3 4] map[i:2 v:4] 4]
`, string(buf))
}

func TestProcNotLastReflect(t *testing.T) {
	jp.CompileScript = compileMapIntProc

	x := jp.MustParseString("[(quux)].v")
	tt.Equal(t, "[(quux)].v", x.String())

	data := []int{2, 3, 4}
	result := x.First(data)
	tt.Equal(t, 2, result)

	got := x.Get(data)
	tt.Equal(t, []any{2, 3, 4}, got)
}
