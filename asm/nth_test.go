// Copyright (c) 2021, Peter Ohler, All rights reserved.

package asm_test

import (
	"testing"

	"github.com/ohler55/ojg/asm"
	"github.com/ohler55/ojg/sen"
	"github.com/ohler55/ojg/tt"
)

func TestNth(t *testing.T) {
	root := testPlan(t,
		`[
           [set $.asm.a [nth [a b c] 1]]
           [set $.asm.b [nth [a b c] -1]]
           [set $.asm.c [nth [a b c] 3]]
         ]`,
		"{src: []}",
	)
	opt := sopt
	opt.Indent = 2
	tt.Equal(t,
		`{
  a: b
  b: c
  c: null
}`, sen.String(root["asm"], &opt))
}

func TestNthArgCount(t *testing.T) {
	p := asm.NewPlan([]any{
		[]any{"nth", []any{}, 1, 2},
	})
	err := p.Execute(map[string]any{})
	tt.NotNil(t, err)
}

func TestNthArgType(t *testing.T) {
	p := asm.NewPlan([]any{
		[]any{"nth", 1, "x"},
	})
	err := p.Execute(map[string]any{})
	tt.NotNil(t, err)
}

func TestNthArgType2(t *testing.T) {
	p := asm.NewPlan([]any{
		[]any{"nth", []any{}, true},
	})
	err := p.Execute(map[string]any{})
	tt.NotNil(t, err)
}
