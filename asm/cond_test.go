// Copyright (c) 2021, Peter Ohler, All rights reserved.

package asm_test

import (
	"testing"

	"github.com/ohler55/ojg/asm"
	"github.com/ohler55/ojg/sen"
	"github.com/ohler55/ojg/tt"
)

func TestCond(t *testing.T) {
	root := testPlan(t,
		`[
           [set $.asm.a [cond [true abc]]]
           [set $.asm.b [cond [false abc][true def]]]
           [set $.asm.c [cond [1 abc][true def]]]
           [set $.asm.d [cond [1 abc][false def]]]
           [set $.asm.e [cond]]
         ]`,
		"{src: []}",
	)
	opt := sopt
	opt.Indent = 2
	tt.Equal(t,
		`{
  a: abc
  b: def
  c: def
  d: null
  e: null
}`, sen.String(root["asm"], &opt))
}

func TestCondArgType(t *testing.T) {
	p := asm.NewPlan([]any{
		[]any{"cond", 1, "x"},
	})
	err := p.Execute(map[string]any{})
	tt.NotNil(t, err)
}

func TestCondArgElementCount(t *testing.T) {
	p := asm.NewPlan([]any{
		[]any{"cond", []any{true, 1, 2}},
	})
	err := p.Execute(map[string]any{})
	tt.NotNil(t, err)
}
