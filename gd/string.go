// Copyright (c) 2020, Peter Ohler, All rights reserved.

package gd

import (
	"strconv"
	"strings"
)

type String string

func (n String) String() string {
	return `"` + string(n) + `"`
}

func (n String) Alter() interface{} {
	return string(n)
}

func (n String) Native() interface{} {
	return string(n)
}

func (n String) Dup() Node {
	return n
}

func (n String) Empty() bool {
	return len(string(n)) == 0
}

func (n String) AsBool() (v Bool, ok bool) {
	switch strings.ToLower(string(n)) {
	case "false":
		v = Bool(false)
		ok = true
	case "true":
		v = Bool(true)
		ok = true
	}
	return
}

func (n String) AsInt() (Int, bool) {
	i, err := strconv.ParseInt(string(n), 10, 64)
	return Int(i), err == nil
}

func (n String) AsFloat() (Float, bool) {
	f, err := strconv.ParseFloat(string(n), 64)
	return Float(f), err == nil
}

func (n String) JSON(_ ...int) string {
	var b strings.Builder

	n.BuildJSON(&b, 0, 0)

	return b.String()
}

func (n String) BuildJSON(b *strings.Builder, _, _ int) {
	b.WriteString(`"`)
	// TBD convert special
	b.WriteString(string(n))
	b.WriteString(`"`)
}
