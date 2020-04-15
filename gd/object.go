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

func (n Object) Empty() bool {
	return len(n) == 0
}

func (n Object) AsBool() (Bool, bool) {
	return Bool(len(n) == 0), false
}

func (n Object) AsInt() (Int, bool) {
	return 0, false
}

func (n Object) AsFloat() (Float, bool) {
	return Float(0.0), false
}

func (n Object) JSON(indent ...int) string {
	var b strings.Builder

	if 0 < len(indent) {
		n.BuildJSON(&b, indent[0], 0)
	} else {
		n.BuildJSON(&b, 0, 0)
	}
	return b.String()
}

func (n Object) BuildJSON(b *strings.Builder, indent, depth int) {
	b.WriteByte('{')
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
				b.WriteString(cs)
				String(k).BuildJSON(b, 0, 0)
				b.WriteByte(':')
				if m := n[k]; m == nil {
					b.WriteString("null")
				} else {
					m.BuildJSON(b, indent, d2)
				}
			}
		} else {
			first := true
			for k, m := range n {
				if first {
					first = false
				} else {
					b.WriteByte(',')
				}
				b.WriteString(cs)
				String(k).BuildJSON(b, 0, 0)
				b.WriteByte(':')
				if m == nil {
					b.WriteString("null")
				} else {
					m.BuildJSON(b, indent, d2)
				}
			}
		}
		b.WriteString(is)
	} else {
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
				String(k).BuildJSON(b, 0, 0)
				b.WriteByte(':')
				if m := n[k]; m == nil {
					b.WriteString("null")
				} else {
					m.BuildJSON(b, 0, 0)
				}
			}
		} else {
			first := true
			for k, m := range n {
				if first {
					first = false
				} else {
					b.WriteByte(',')
				}
				String(k).BuildJSON(b, 0, 0)
				b.WriteByte(':')
				if m == nil {
					b.WriteString("null")
				} else {
					m.BuildJSON(b, 0, 0)
				}
			}
		}
	}
	b.WriteByte('}')
}
