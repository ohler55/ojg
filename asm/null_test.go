// Copyright (c) 2021, Peter Ohler, All rights reserved.

package asm_test

import (
	"testing"

	"github.com/ohler55/ojg/asm"
	"github.com/ohler55/ojg/sen"
	"github.com/ohler55/ojg/tt"
)

func TestNull(t *testing.T) {
	root := testPlan(t,
		`[
           [set $.asm.a [null? null]]
           [set $.asm.b [null? a_string]]
         ]`,
		"{src: []}",
	)
	tt.Equal(t, "{a:true b:false}", sen.String(root["asm"], &sopt))
}

func TestNullArgCountT(t *testing.T) {
	p := asm.NewPlan([]interface{}{
		[]interface{}{"null?", 1, 2},
	})
	err := p.Execute(map[string]interface{}{})
	tt.NotNil(t, err)
}
