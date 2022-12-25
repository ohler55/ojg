// Copyright (c) 2021, Peter Ohler, All rights reserved.

package asm

var asmFn = Fn{
	Name: "asm",
	Eval: asmEval,
	Desc: `Processes all arguments in order using the return of each as
input for the next.`,
}

func init() {
	Define(&asmFn)
}

func asmEval(root map[string]any, at any, args ...any) any {
	for _, a := range args {
		at = evalArg(root, at, a)
	}
	return at
}
