// Copyright (c) 2021, Peter Ohler, All rights reserved.

package asm_test

import (
	"fmt"
	"testing"

	"github.com/ohler55/ojg/asm"
	"github.com/ohler55/ojg/jp"
	"github.com/ohler55/ojg/sen"
	"github.com/ohler55/ojg/tt"
)

func TestInspect(t *testing.T) {
	p := asm.NewPlan([]interface{}{
		"asm",
		[]interface{}{"inspect", "test", "$"},
	})
	tt.Equal(t, "[asm [inspect test $]]", sen.String(p), "inspect plan simplify")
	fn, _ := p.Args[0].(*asm.Fn)
	tt.NotNil(t, fn)
	tt.Equal(t, "[inspect test $]", fn.String(), "inspect string")
}

func Example_inspect() {
	p := asm.NewPlan([]interface{}{
		"asm",
		[]interface{}{"inspect", 0, "one", []interface{}{1}, "@", jp.C("src").N(1), "test", "$"},
	})
	root := map[string]interface{}{
		"src": []interface{}{1, int64(2), 3},
	}
	if err := p.Execute(root); err != nil {
		fmt.Println(err.Error())
	}
	// Output:
	// one: [1]
	// {"src":[1,2,3]}
	// test:
	// {
	//   "src": [
	//     1,
	//     2,
	//     3
	//   ]
	// }
}
