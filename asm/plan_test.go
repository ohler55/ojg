// Copyright (c) 2021, Peter Ohler, All rights reserved.

package asm_test

import (
	"testing"

	"github.com/ohler55/ojg/asm"
	"github.com/ohler55/ojg/sen"
	"github.com/ohler55/ojg/tt"
)

func TestPlanSimplify(t *testing.T) {
	p := asm.Plan{
		Fn: asm.Fn{
			Name: "fun",
			Args: []interface{}{
				&asm.Fn{Name: "+", Args: []interface{}{3, 4}},
				&asm.Fn{Name: "list", Args: []interface{}{1, 2, 3}},
			},
		},
	}
	tt.Equal(t, "[fun [+ 3 4][list 1 2 3]]", sen.String(&p), "plan simplify")
}

func TestPlanNil(t *testing.T) {
	p := asm.NewPlan(nil)
	tt.Nil(t, p)

	p = asm.NewPlan([]interface{}{})
	tt.Nil(t, p)
}

func TestPlanNonAsm(t *testing.T) {
	p := asm.NewPlan([]interface{}{"set", "$.asm.x", 1})
	tt.NotNil(t, p)

	root := map[string]interface{}{"src": []interface{}{}}
	err := p.Execute(root)
	tt.Nil(t, err)
	tt.Equal(t, "{asm:{x:1} src:[]}", sen.String(root, &sopt))
}

func TestPlanImpliedAsm(t *testing.T) {
	p := asm.NewPlan([]interface{}{
		[]interface{}{"set", "$.asm.x", 1},
	})
	tt.NotNil(t, p)

	root := map[string]interface{}{"src": []interface{}{}}
	err := p.Execute(root)
	tt.Nil(t, err)
	tt.Equal(t, "{asm:{x:1} src:[]}", sen.String(root, &sopt))
}

func TestPlanPanic(t *testing.T) {
	asm.Define(&asm.Fn{
		Name: "panic",
		Eval: func(_ map[string]interface{}, _ interface{}, _ ...interface{}) interface{} {
			panic("abort")
		},
	})
	p := asm.NewPlan([]interface{}{
		[]interface{}{"panic"},
	})
	tt.NotNil(t, p)

	root := map[string]interface{}{"src": []interface{}{}}
	err := p.Execute(root)
	tt.NotNil(t, err)
}
