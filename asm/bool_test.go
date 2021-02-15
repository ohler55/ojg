// Copyright (c) 2021, Peter Ohler, All rights reserved.

package asm_test

import (
	"testing"

	"github.com/ohler55/ojg/asm"
	"github.com/ohler55/ojg/sen"
	"github.com/ohler55/ojg/tt"
)

func TestBool(t *testing.T) {
	root := testPlan(t,
		`[
           [set $.asm.a [bool? true]]
           [set $.asm.b [bool? false]]
           [set $.asm.c [bool? 3]]
         ]`,
		"{src: []}",
	)
	tt.Equal(t, "{a:true b:true c:false}", sen.String(root["asm"], &sopt))
}

func TestBoolArgCountT(t *testing.T) {
	p := asm.NewPlan([]interface{}{
		[]interface{}{"bool?", 1, 2},
	})
	err := p.Execute(map[string]interface{}{})
	tt.NotNil(t, err)
}
