// Copyright (c) 2020, Peter Ohler, All rights reserved.

package gen

import (
	"sort"
	"unsafe"
)

var Sort = false

type Object map[string]Node

func (n Object) String() string {
	b := []byte{'{'}
	first := true

	if Sort {
		keys := make([]string, 0, len(n))
		for k := range n {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		for i, k := range keys {
			if 0 < i {
				b = append(b, ',')
			}
			b = append(b, '"')
			b = append(b, k...)
			b = append(b, '"')
			b = append(b, ':')
			if m := n[k]; m == nil {
				b = append(b, "null"...)
			} else {
				b = append(b, m.String()...)
			}
		}
	} else {
		for k, m := range n {
			if first {
				first = false
			} else {
				b = append(b, ',')
			}
			b = append(b, '"')
			b = append(b, k...)
			b = append(b, '"')
			b = append(b, ':')
			if m == nil {
				b = append(b, "null"...)
			} else {
				b = append(b, m.String()...)
			}
		}
	}
	b = append(b, '}')

	return string(b)
}

func (n Object) Alter() interface{} {
	var simple map[string]interface{}

	if n != nil {
		simple = *(*map[string]interface{})(unsafe.Pointer(&n))
		for k, m := range n {
			if m == nil {
				simple[k] = nil
			} else {
				simple[k] = m.Alter()
			}
		}
	}
	return simple
}

func (n Object) Simplify() interface{} {
	var dup map[string]interface{}

	if n != nil {
		dup = map[string]interface{}{}
		for k, m := range n {
			if m == nil {
				dup[k] = m
			} else {
				dup[k] = m.Simplify()
			}
		}
	}
	return dup
}

func (n Object) Dup() Node {
	var o Object

	if n != nil {
		o = Object{}
		for k, m := range n {
			if m == nil {
				o[k] = nil
			} else {
				o[k] = m.Dup()
			}
		}
	}
	return o
}

func (n Object) Empty() bool {
	return len(n) == 0
}
