// Copyright (c) 2021, Peter Ohler, All rights reserved.

package asm_test

import (
	"testing"

	"github.com/ohler55/ojg/asm"
	"github.com/ohler55/ojg/sen"
	"github.com/ohler55/ojg/tt"
)

func TestAppend(t *testing.T) {
	root := testPlan(t,
		`[
           [set $.asm.a [append [] a]]
           [set $.asm.b [append [a] 1]]
         ]`,
		"{src: []}",
	)
	tt.Equal(t,
		`{a:[a] b:[a 1]}`, sen.String(root["asm"], &sopt))
}

func TestAppendArgCount(t *testing.T) {
	p := asm.NewPlan([]any{
		[]any{"append", []any{}, 1, 2},
	})
	err := p.Execute(map[string]any{})
	tt.NotNil(t, err)
}

func TestAppendArgType(t *testing.T) {
	p := asm.NewPlan([]any{
		[]any{"append", 1, "x"},
	})
	err := p.Execute(map[string]any{})
	tt.NotNil(t, err)
}
