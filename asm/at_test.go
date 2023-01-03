// Copyright (c) 2021, Peter Ohler, All rights reserved.

package asm_test

import (
	"testing"

	"github.com/ohler55/ojg/asm"
	"github.com/ohler55/ojg/sen"
	"github.com/ohler55/ojg/tt"
)

func TestAt(t *testing.T) {
	root := testPlan(t,
		`[
           {x:3}
           [set $.asm [get [at x]]]
         ]`,
		"{src: []}",
	)
	tt.Equal(t, `3`, sen.String(root["asm"], &sopt))
}

func TestAtArgNotString(t *testing.T) {
	p := asm.NewPlan([]any{
		[]any{"at", 1},
	})
	err := p.Execute(map[string]any{})
	tt.NotNil(t, err)
}

func TestAtArgParseError(t *testing.T) {
	p := asm.NewPlan([]any{
		[]any{"at", "[[["},
	})
	err := p.Execute(map[string]any{})
	tt.NotNil(t, err)
}
