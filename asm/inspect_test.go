// Copyright (c) 2021, Peter Ohler, All rights reserved.

package asm_test

import (
	"fmt"
	"testing"

	"github.com/ohler55/ojg/asm"
	"github.com/ohler55/ojg/sen"
	"github.com/ohler55/ojg/tt"
)

func TestInspect(t *testing.T) {
	p := asm.NewPlan([]interface{}{
		"asm",
		[]interface{}{"inspect", "test", "$"},
	})
	tt.Equal(t, "[asm [inspect test $]]", sen.String(p), "inspect plan simplify")
}

func ExampleInspect() {
	p := asm.NewPlan([]interface{}{
		"asm",
		[]interface{}{"inspect", "test", "$"},
	})
	root := map[string]interface{}{
		"src": []interface{}{1, 2, 3},
	}
	if err := p.Execute(root); err != nil {
		fmt.Println(err.Error())
	}
	// Output:
	// test: {"src":[1,2,3]}
}
