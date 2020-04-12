// Copyright (c) 2020, Peter Ohler, All rights reserved.

package gd

type Bool bool

func (n Bool) String() (s string) {
	if n {
		s = "true"
	} else {
		s = "false"
	}
	return
}

func (n Bool) Alter() interface{} {
	return bool(n)
}

func (n Bool) Native() interface{} {
	return bool(n)
}

func (n Bool) Dup() Node {
	return n
}

func (n Bool) JSON(_ ...int) string {
	return n.String()
}
