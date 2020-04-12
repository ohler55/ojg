// Copyright (c) 2020, Peter Ohler, All rights reserved.

package gd

import (
	"strings"
	"unsafe"
)

type Array []Node

func (n Array) String() string {
	var b strings.Builder

	b.WriteString("[")
	for i, m := range n {
		if 0 < i {
			b.WriteString(",")
		}
		b.WriteString(m.String())
	}
	b.WriteString("]")

	return b.String()
}

func (n Array) Alter() interface{} {
	var native []interface{}

	if n != nil {
		native = *(*[]interface{})(unsafe.Pointer(&n))
		for i, m := range n {
			native[i] = m.Alter()
		}
	}
	return native
}

func (n Array) Native() interface{} {
	var dup []interface{}

	if n != nil {
		dup = make([]interface{}, 0, len(n))
		for _, m := range n {
			dup = append(dup, m.Native())
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
