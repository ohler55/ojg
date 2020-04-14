// Copyright (c) 2020, Peter Ohler, All rights reserved.

package gd

import "strings"

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

func (n Bool) Empty() bool {
	return false
}

func (n Bool) AsBool() (Bool, bool) {
	return n, true
}

func (n Bool) AsInt() (Int, bool) {
	var i int64
	if n {
		i = 1
	}
	return Int(i), false
}

func (n Bool) AsFloat() (Float, bool) {
	var f float64
	if n {
		f = 1.0
	}
	return Float(f), false
}

func (n Bool) JSON(_ ...int) string {
	return n.String()
}

func (n Bool) BuildJSON(b *strings.Builder, _, _ int) {
	if n {
		b.WriteString("true")
	} else {
		b.WriteString("false")
	}
}
