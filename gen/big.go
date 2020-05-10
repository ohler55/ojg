// Copyright (c) 2020, Peter Ohler, All rights reserved.

package gen

type Big string

func (n Big) String() string {
	return string(n)
}

func (n Big) Alter() interface{} {
	return string(n)
}

func (n Big) Simplify() interface{} {
	return string(n)
}

func (n Big) Dup() Node {
	return n
}

func (n Big) Empty() bool {
	return len(string(n)) == 0
}
