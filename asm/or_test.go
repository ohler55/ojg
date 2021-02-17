// Copyright (c) 2021, Peter Ohler, All rights reserved.

package asm_test

import (
	"testing"

	"github.com/ohler55/ojg/asm"
	"github.com/ohler55/ojg/sen"
	"github.com/ohler55/ojg/tt"
)

func TestOrTrue(t *testing.T) {
	root := testPlan(t,
		`[
           [set $.asm [or false "$.src[0]"]]
         ]`,
		"{src: [true false]}",
	)
	tt.Equal(t, "true", sen.String(root["asm"]))
}

func TestOrFalse(t *testing.T) {
	root := testPlan(t,
		`[
           [set $.asm [or false "$.src[1]" false]]
         ]`,
		"{src: [true false]}",
	)
	tt.Equal(t, "false", sen.String(root["asm"]))
}

func TestOrNull(t *testing.T) {
	root := testPlan(t,
		`[
           [set $.asm [or "$.src[2]"]]
         ]`,
		"{src: [true false]}",
	)
	tt.Equal(t, "false", sen.String(root["asm"]))
}

func TestOrNotBool(t *testing.T) {
	p := asm.NewPlan([]interface{}{
		[]interface{}{"or", 1, 2},
	})
	err := p.Execute(map[string]interface{}{})
	tt.NotNil(t, err)
}
