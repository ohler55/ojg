// Copyright (c) 2020, Peter Ohler, All rights reserved.

package gen

import (
	"strings"
	"unsafe"
)

type Array []Node

func (n Array) String() string {
	var b strings.Builder

	b.WriteByte('[')
	for i, m := range n {
		if 0 < i {
			b.WriteByte(',')
		}
		if m == nil {
			b.WriteString("null")
		} else {
			b.WriteString(m.String())
		}
	}
	b.WriteByte(']')

	return b.String()
}

func (n Array) Alter() interface{} {
	var simple []interface{}

	if n != nil {
		simple = *(*[]interface{})(unsafe.Pointer(&n))
		for i, m := range n {
			simple[i] = m.Alter()
		}
	}
	return simple
}

func (n Array) Simplify() interface{} {
	var dup []interface{}

	if n != nil {
		dup = make([]interface{}, 0, len(n))
		for _, m := range n {
			dup = append(dup, m.Simplify())
		}
	}
	return dup
}

func (n Array) Dup() Node {
	var a Array

	if n != nil {
		a = make(Array, 0, len(n))
		for _, m := range n {
			a = append(a, m.Dup())
		}
	}
	return a
}

func (n Array) Empty() bool {
	return len(n) == 0
}
