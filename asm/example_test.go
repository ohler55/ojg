// Copyright (c) 2021, Peter Ohler, All rights reserved.

package asm_test

import (
	"fmt"

	"github.com/ohler55/ojg/asm"
	"github.com/ohler55/ojg/jp"
)

func ExamplePlan() {
	p := asm.NewPlan([]any{
		"asm",
		[]any{"inspect", 0, "one", []any{1}, "@", jp.C("src").N(1), "test", "$"},
	})
	root := map[string]any{
		"src": []any{1, int64(2), 3},
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
