// Copyright (c) 2020, Peter Ohler, All rights reserved.

package gen

import (
	"sort"
	"strings"
	"unsafe"
)

var Sort = false

type Object map[string]Node

func (n Object) String() string {
	var b strings.Builder
	first := true

	b.WriteByte('{')
	if Sort {
		keys := make([]string, 0, len(n))
		for k := range n {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		for i, k := range keys {
			if 0 < i {
				b.WriteByte(',')
			}
			b.WriteByte('"')
			b.WriteString(k)
			b.WriteString(`":`)
			if m := n[k]; m == nil {
				b.WriteString("null")
			} else {
				b.WriteString(m.String())
			}
		}
	} else {
		for k, m := range n {
			if first {
				first = false
			} else {
				b.WriteByte(',')
			}
			b.WriteByte('"')
			b.WriteString(k)
			b.WriteString(`":`)
			if m == nil {
				b.WriteString("null")
			} else {
				b.WriteString(m.String())
			}
		}
	}
	b.WriteByte('}')

	return b.String()
}

func (n Object) Alter() interface{} {
	var simple map[string]interface{}

	if n != nil {
		simple = *(*map[string]interface{})(unsafe.Pointer(&n))
		for k, m := range n {
			simple[k] = m.Alter()
		}
	}
	return simple
}

func (n Object) Simplify() interface{} {
	var dup map[string]interface{}

	if n != nil {
		dup = map[string]interface{}{}
		for k, m := range n {
			dup[k] = m.Simplify()
		}
	}
	return dup
}

func (n Object) Dup() Node {
	var o Object

	if n != nil {
		o = Object{}
		for k, m := range n {
			o[k] = m.Dup()
		}
	}
	return o
}

func (n Object) Empty() bool {
	return len(n) == 0
}
