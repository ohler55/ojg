// Copyright (c) 2021, Peter Ohler, All rights reserved.

package asm

func init() {
	Define(&Fn{
		Name: "list",
		Eval: list,
		Desc: `Creates a list from all the argument and return that list.`,
	})
}

func list(root map[string]interface{}, at interface{}, args ...interface{}) interface{} {
	var a []interface{}
	for _, arg := range args {
		a = append(a, evalArg(root, at, arg))
	}
	return a
}
