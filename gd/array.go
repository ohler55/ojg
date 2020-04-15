// Copyright (c) 2020, Peter Ohler, All rights reserved.

package gd

import (
	"strings"
	"unsafe"
)

const spaces = "\n                                                                                                                                "

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

func (n Array) Empty() bool {
	return len(n) == 0
}

func (n Array) AsBool() (Bool, bool) {
	return Bool(len(n) == 0), false
}

func (n Array) AsInt() (Int, bool) {
	return 0, false
}

func (n Array) AsFloat() (Float, bool) {
	return Float(0.0), false
}

func (n Array) JSON(indent ...int) string {
	var b strings.Builder

	if 0 < len(indent) {
		n.BuildJSON(&b, indent[0], 0)
	} else {
		n.BuildJSON(&b, 0, 0)
	}
	return b.String()
}

func (n Array) BuildJSON(b *strings.Builder, indent, depth int) {
	b.WriteByte('[')
	if 0 < indent {
		x := depth*indent + 1
		if len(spaces) < x {
			x = depth*indent + 1
		}
		is := spaces[0:x]
		d2 := depth + 1
		x = d2*indent + 1
		if len(spaces) < x {
			x = depth*indent + 1
		}
		cs := spaces[0:x]

		for j, m := range n {
			if 0 < j {
				b.WriteByte(',')
			}
			b.WriteString(cs)
			if m == nil {
				b.WriteString("null")
			} else {
				m.BuildJSON(b, indent, d2)
			}
		}
		b.WriteString(is)
	} else {
		for j, m := range n {
			if 0 < j {
				b.WriteByte(',')
			}
			if m == nil {
				b.WriteString("null")
			} else {
				m.BuildJSON(b, 0, 0)
			}
		}
	}
	b.WriteByte(']')
}
