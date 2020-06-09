// Copyright (c) 2020, Peter Ohler, All rights reserved.

package gen

import (
	"unsafe"
)

type Array []Node

var EmptyArray = Array{}

func (n Array) String() string {
	b := []byte{'['}
	for i, m := range n {
		if 0 < i {
			b = append(b, ',')
		}
		if m == nil {
			b = append(b, "null"...)
		} else {
			b = append(b, m.String()...)
		}
	}
	b = append(b, ']')

	return string(b)
}

func (n Array) Alter() interface{} {
	var simple []interface{}

	if n != nil {
		simple = *(*[]interface{})(unsafe.Pointer(&n))
		for i, m := range n {
			if m == nil {
				simple[i] = nil
			} else {
				simple[i] = m.Alter()
			}
		}
	}
	return simple
}

func (n Array) Simplify() interface{} {
	var dup []interface{}

	if n != nil {
		dup = make([]interface{}, 0, len(n))
		for _, m := range n {
			if m == nil {
				dup = append(dup, nil)
			} else {
				dup = append(dup, m.Simplify())
			}
		}
	}
	return dup
}

func (n Array) Dup() Node {
	var a Array

	if n != nil {
		a = make(Array, 0, len(n))
		for _, m := range n {
			if m == nil {
				a = append(a, nil)
			} else {
				a = append(a, m.Dup())
			}
		}
	}
	return a
}

func (n Array) Empty() bool {
	return len(n) == 0
}
