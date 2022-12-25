// Copyright (c) 2021, Peter Ohler, All rights reserved.

package asm_test

import (
	"testing"

	"github.com/ohler55/ojg/asm"
	"github.com/ohler55/ojg/jp"
	"github.com/ohler55/ojg/sen"
	"github.com/ohler55/ojg/tt"
)

func TestDelall(t *testing.T) {
	root := testPlan(t,
		`[
           [{one:1 two:2 three:3}{one:4 two:5 three:6}]
           [delall "@.*.one"]
           [set $.asm @]
           [delall "$.asm.*.three"]
         ]`,
		"{src: [1 2 3]}",
	)
	tt.Equal(t, "[{two:2}{two:5}]", sen.String(root["asm"]))
}

func TestDelallExprError(t *testing.T) {
	p := asm.NewPlan([]any{
		map[string]any{},
		[]any{"delall", jp.D()},
	})
	err := p.Execute(map[string]any{})
	tt.NotNil(t, err)
}

func TestDelallArgCount(t *testing.T) {
	p := asm.NewPlan([]any{
		[]any{"delall"},
	})
	err := p.Execute(map[string]any{})
	tt.NotNil(t, err)
}

func TestDelallArgNotExpr(t *testing.T) {
	p := asm.NewPlan([]any{
		[]any{"delall", 1},
	})
	err := p.Execute(map[string]any{})
	tt.NotNil(t, err)
}
