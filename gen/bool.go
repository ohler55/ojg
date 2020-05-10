// Copyright (c) 2020, Peter Ohler, All rights reserved.

package gen

type Bool bool

var True = Bool(true)
var False = Bool(false)

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

func (n Bool) Simplify() interface{} {
	return bool(n)
}

func (n Bool) Dup() Node {
	return n
}

func (n Bool) Empty() bool {
	return false
}
