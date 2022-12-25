// Copyright (c) 2021, Peter Ohler, All rights reserved.

package asm

func init() {
	Define(&Fn{
		Name:    "quote",
		Eval:    quote,
		Compile: func(*Fn) {},
		Desc: `Does not evaluate arguments. One argument is expected. Null is
returned if no arguments are given while any arguments other
than the first are ignored. An example for use would be to
treats "@.x" as a string instead of as a path.`,
	})
}

func quote(root map[string]any, at any, args ...any) (val any) {
	if 0 < len(args) {
		val = args[0]
	}
	return
}
