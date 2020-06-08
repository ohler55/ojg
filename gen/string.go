// Copyright (c) 2020, Peter Ohler, All rights reserved.

package gen

type String string

func (n String) String() string {
	return `"` + string(n) + `"`
}

func (n String) Alter() interface{} {
	return string(n)
}

func (n String) Simplify() interface{} {
	return string(n)
}

func (n String) Dup() Node {
	return n
}

func (n String) Empty() bool {
	return len(string(n)) == 0
}
