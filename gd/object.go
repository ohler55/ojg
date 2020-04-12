// Copyright (c) 2020, Peter Ohler, All rights reserved.

package gd

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

	b.WriteString("{")
	if Sort {
		keys := make([]string, 0, len(n))
		for k := range n {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		for i, k := range keys {
			if 0 < i {
				b.WriteString(",")
			}
			b.WriteString(`"`)
			b.WriteString(k)
			b.WriteString(`":`)
			b.WriteString(n[k].String())
		}
	} else {
		for k, m := range n {
			if first {
				first = false
			} else {
				b.WriteString(",")
			}
			b.WriteString(`"`)
			b.WriteString(k)
			b.WriteString(`":`)
			b.WriteString(m.String())
		}
	}
	b.WriteString("}")

	return b.String()
}

func (n Object) Alter() interface{} {
	var native map[string]interface{}

	if n != nil {
		native = *(*map[string]interface{})(unsafe.Pointer(&n))
		for k, m := range n {
			native[k] = m.Alter()
		}
	}
	return native
}

func (n Object) Native() interface{} {
	var dup map[string]interface{}

	if n != nil {
		dup = map[string]interface{}{}
		for k, m := range n {
			dup[k] = m.Native()
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
