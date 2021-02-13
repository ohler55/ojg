// Copyright (c) 2021, Peter Ohler, All rights reserved.

package asm

func init() {
	Define(&Fn{Name: "+", Eval: add})
}

func add(root map[string]interface{}, at interface{}, args ...interface{}) interface{} {

	return nil
}
