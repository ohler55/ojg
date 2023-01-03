// Copyright (c) 2021, Peter Ohler, All rights reserved.

package asm_test

import (
	"testing"

	"github.com/ohler55/ojg/asm"
	"github.com/ohler55/ojg/sen"
	"github.com/ohler55/ojg/tt"
)

func TestGet(t *testing.T) {
	root := testPlan(t,
		`[
           {one:1 two:2}
           [set $.at [get @.one]]
           [set $.root [get $.src.b]]
           [set $.arg [get @.x {x:1 y:2}]]
         ]`,
		"{src: {a:1 b:2 c:3}}",
	)
	tt.Equal(t, "1", sen.String(root["at"]))
	tt.Equal(t, "2", sen.String(root["root"]))
	tt.Equal(t, "1", sen.String(root["arg"]))
}

func TestGetArgCount(t *testing.T) {
	p := asm.NewPlan([]any{
		[]any{"get"},
	})
	err := p.Execute(map[string]any{})
	tt.NotNil(t, err)
}

func TestGetArgNotExpr(t *testing.T) {
	p := asm.NewPlan([]any{
		[]any{"get", 1},
	})
	err := p.Execute(map[string]any{})
	tt.NotNil(t, err)
}

func TestGetArgType(t *testing.T) {
	p := asm.NewPlan([]any{
		[]any{"get", []any{"sum"}},
	})
	err := p.Execute(map[string]any{})
	tt.NotNil(t, err)
}
