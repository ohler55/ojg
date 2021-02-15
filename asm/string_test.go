// Copyright (c) 2021, Peter Ohler, All rights reserved.

package asm_test

import (
	"testing"

	"github.com/ohler55/ojg/asm"
	"github.com/ohler55/ojg/sen"
	"github.com/ohler55/ojg/tt"
)

func TestString(t *testing.T) {
	root := testPlan(t,
		`[
           [set $.asm.a [string? abc]]
           [set $.asm.b [string? 123]]
           [set $.asm.c [string? true]]
         ]`,
		"{src: []}",
	)
	tt.Equal(t, "{a:true b:false c:false}", sen.String(root["asm"], &sopt))
}

func TestStringArgCountT(t *testing.T) {
	p := asm.NewPlan([]interface{}{
		[]interface{}{"string?", 1, 2},
	})
	err := p.Execute(map[string]interface{}{})
	tt.NotNil(t, err)
}
