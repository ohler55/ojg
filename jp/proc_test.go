// Copyright (c) 2024, Peter Ohler, All rights reserved.

package jp_test

import (
	"fmt"
	"testing"

	"github.com/ohler55/ojg/jp"
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

func TestProc(t *testing.T) {
	jp.CompileScript = compileMathProc

	p := jp.MustNewProc([]byte("(+ 0 1)"))
	tt.Equal(t, "[(+ 0 1)]", p.String())

	x := jp.MustParseString("[(+ 0 1)]")
	tt.Equal(t, "[(+ 0 1)]", x.String())

	data := []any{2, 3, 4}
	result := x.First(data)
	tt.Equal(t, 5, result)

	// TBD Get
}
