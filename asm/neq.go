// Copyright (c) 2021, Peter Ohler, All rights reserved.

package asm

func init() {
	Define(&Fn{
		Name: "neq",
		Eval: neq,
		Desc: `Returns true if any the argument are not equal. An alias is !==.`,
	})
	Define(&Fn{
		Name: "!=",
		Eval: neq,
		Desc: `Returns true if any the argument are not equal. An alias is !==.`,
	})
}

func neq(root map[string]interface{}, at interface{}, args ...interface{}) interface{} {
	b, _ := equal(root, at, args...).(bool)

	return !b
}
