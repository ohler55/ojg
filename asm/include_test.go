// Copyright (c) 2021, Peter Ohler, All rights reserved.

package asm_test

import (
	"testing"

	"github.com/ohler55/ojg/asm"
	"github.com/ohler55/ojg/sen"
	"github.com/ohler55/ojg/tt"
)

func TestInclude(t *testing.T) {
	root := testPlan(t,
		`[
           [set $.asm.a [include [a b c] b]]
           [set $.asm.b [include [a b c] x]]
           [set $.asm.c [include [1 2 3] 3]]
           [set $.asm.d [include abcdef cd]]
           [set $.asm.e [include abcdef cx]]
         ]`,
		"{src: []}",
	)
	opt := sopt
	opt.Indent = 2
	tt.Equal(t,
		`{
  a: true
  b: false
  c: true
  d: true
  e: false
}`, sen.String(root["asm"], &opt))
}

func TestIncludeArgCount(t *testing.T) {
	p := asm.NewPlan([]any{
		[]any{"include", []any{}, "x", 1},
	})
	err := p.Execute(map[string]any{})
	tt.NotNil(t, err)
}

func TestIncludeArgType(t *testing.T) {
	p := asm.NewPlan([]any{
		[]any{"include", 1, "x"},
	})
	err := p.Execute(map[string]any{})
	tt.NotNil(t, err)
}

func TestIncludeArgType2(t *testing.T) {
	p := asm.NewPlan([]any{
		[]any{"include", "abc", 1},
	})
	err := p.Execute(map[string]any{})
	tt.NotNil(t, err)
}
