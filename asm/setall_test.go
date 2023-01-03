// Copyright (c) 2021, Peter Ohler, All rights reserved.

package asm_test

import (
	"testing"

	"github.com/ohler55/ojg/asm"
	"github.com/ohler55/ojg/jp"
	"github.com/ohler55/ojg/sen"
	"github.com/ohler55/ojg/tt"
)

func TestSetall(t *testing.T) {
	root := testPlan(t,
		`[
           {x: 2 y: 3}
           [setall "@.*" 1]
           [set $.asm @]
         ]`,
		"{src: [1 2 3]}",
	)
	tt.Equal(t, "{x:1 y:1}", sen.String(root["asm"], &sopt))
}

func TestSetallExprError(t *testing.T) {
	p := asm.NewPlan([]any{
		map[string]any{}, // Sets @
		[]any{"setall", jp.D(), 1},
		[]any{"setall", "$.asm", "@"},
	})
	err := p.Execute(map[string]any{})
	tt.NotNil(t, err)
}

func TestSetallArgCount(t *testing.T) {
	p := asm.NewPlan([]any{
		[]any{"setall", "@.x"},
	})
	err := p.Execute(map[string]any{})
	tt.NotNil(t, err)
}

func TestSetallArgNotExpr(t *testing.T) {
	p := asm.NewPlan([]any{
		[]any{"setall", 1, 2},
	})
	err := p.Execute(map[string]any{})
	tt.NotNil(t, err)
}

func TestSetallArgType(t *testing.T) {
	p := asm.NewPlan([]any{
		[]any{"setall", []any{"sum"}, 1},
	})
	err := p.Execute(map[string]any{})
	tt.NotNil(t, err)
}
