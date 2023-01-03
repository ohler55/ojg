// Copyright (c) 2021, Peter Ohler, All rights reserved.

package asm_test

import (
	"testing"

	"github.com/ohler55/ojg/asm"
	"github.com/ohler55/ojg/sen"
	"github.com/ohler55/ojg/tt"
)

func TestNotTrue(t *testing.T) {
	root := testPlan(t,
		`[
           [set $.asm [not "$.src[0]"]]
         ]`,
		"{src: [true false]}",
	)
	tt.Equal(t, "false", sen.String(root["asm"]))
}

func TestNotFalse(t *testing.T) {
	root := testPlan(t,
		`[
           [set $.asm [not "$.src[1]"]]
         ]`,
		"{src: [true false]}",
	)
	tt.Equal(t, "true", sen.String(root["asm"]))
}

func TestNotArgCount(t *testing.T) {
	p := asm.NewPlan([]any{
		[]any{"not", true, false},
	})
	err := p.Execute(map[string]any{})
	tt.NotNil(t, err)
}

func TestNotNotBool(t *testing.T) {
	p := asm.NewPlan([]any{
		[]any{"not", 1},
	})
	err := p.Execute(map[string]any{})
	tt.NotNil(t, err)
}
