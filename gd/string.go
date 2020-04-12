// Copyright (c) 2020, Peter Ohler, All rights reserved.

package gd

type String string

func (n String) String() string {
	return `"` + string(n) + `"`
}

func (n String) Alter() interface{} {
	return string(n)
}

func (n String) Native() interface{} {
	return string(n)
}

func (n String) Dup() Node {
	return n
}
