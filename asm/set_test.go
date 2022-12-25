// Copyright (c) 2021, Peter Ohler, All rights reserved.

package asm_test

import (
	"testing"

	"github.com/ohler55/ojg/asm"
	"github.com/ohler55/ojg/jp"
	"github.com/ohler55/ojg/sen"
	"github.com/ohler55/ojg/tt"
)

func TestSet(t *testing.T) {
	root := testPlan(t,
		`[
           {}
           [set @.one 1]
           [set $.asm @]
         ]`,
		"{src: [1 2 3]}",
	)
	tt.Equal(t, "{one:1}", sen.String(root["asm"]))
}

func TestSetFn(t *testing.T) {
	root := testPlan(t,
		`[
           {}
           [set [at one two] 2]
           [set [root asm] @]
         ]`,
		"{src: []}",
	)
	tt.Equal(t, "{one:{two:2}}", sen.String(root["asm"]))
}

func TestSetExprError(t *testing.T) {
	p := asm.NewPlan([]any{
		map[string]any{}, // Sets @
		[]any{"set", jp.D(), 1},
		[]any{"set", "$.asm", "@"},
	})
	err := p.Execute(map[string]any{})
	tt.NotNil(t, err)
}

func TestSetArgCount(t *testing.T) {
	p := asm.NewPlan([]any{
		[]any{"set", "@.x"},
	})
	err := p.Execute(map[string]any{})
	tt.NotNil(t, err)
}

func TestSetArgNotExpr(t *testing.T) {
	p := asm.NewPlan([]any{
		[]any{"set", 1, 2},
	})
	err := p.Execute(map[string]any{})
	tt.NotNil(t, err)
}

func TestSetArgType(t *testing.T) {
	p := asm.NewPlan([]any{
		[]any{"set", []any{"sum"}, 1},
	})
	err := p.Execute(map[string]any{})
	tt.NotNil(t, err)
}
