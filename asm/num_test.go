// Copyright (c) 2021, Peter Ohler, All rights reserved.

package asm_test

import (
	"testing"

	"github.com/ohler55/ojg/asm"
	"github.com/ohler55/ojg/sen"
	"github.com/ohler55/ojg/tt"
)

func TestNum(t *testing.T) {
	root := testPlan(t,
		`[
           [set $.asm.a [num? 2]]
           [set $.asm.b [num? 2.3]]
           [set $.asm.c [num? true]]
         ]`,
		"{src: []}",
	)
	tt.Equal(t, "{a:true b:true c:false}", sen.String(root["asm"], &sopt))
}

func TestNumArgCountT(t *testing.T) {
	p := asm.NewPlan([]interface{}{
		[]interface{}{"num?", 1, 2},
	})
	err := p.Execute(map[string]interface{}{})
	tt.NotNil(t, err)
}
