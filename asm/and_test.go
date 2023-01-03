// Copyright (c) 2021, Peter Ohler, All rights reserved.

package asm_test

import (
	"testing"

	"github.com/ohler55/ojg/asm"
	"github.com/ohler55/ojg/sen"
	"github.com/ohler55/ojg/tt"
)

func TestAndTrue(t *testing.T) {
	root := testPlan(t,
		`[
           [set $.asm [and true "$.src[0]"]]
         ]`,
		"{src: [true false]}",
	)
	tt.Equal(t, "true", sen.String(root["asm"]))
}

func TestAndFalse(t *testing.T) {
	root := testPlan(t,
		`[
           [set $.asm [and true "$.src[0]" false]]
         ]`,
		"{src: [true false]}",
	)
	tt.Equal(t, "false", sen.String(root["asm"]))
}

func TestAndNull(t *testing.T) {
	root := testPlan(t,
		`[
           [set $.asm [and "$.src[2]"]]
         ]`,
		"{src: [true false]}",
	)
	tt.Equal(t, "false", sen.String(root["asm"]))
}

func TestAndNotBool(t *testing.T) {
	p := asm.NewPlan([]any{
		[]any{"and", 1, 2},
	})
	err := p.Execute(map[string]any{})
	tt.NotNil(t, err)
}
