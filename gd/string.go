// Copyright (c) 2020, Peter Ohler, All rights reserved.

package gd

import (
	"strconv"
	"strings"
)

const hex = "0123456789abcdef"

type String string

func (n String) String() string {
	return `"` + string(n) + `"`
}

func (n String) Alter() interface{} {
	return string(n)
}

func (n String) Simplify() interface{} {
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
	b.WriteByte('"')
	for _, r := range string(n) {
		switch r {
		case '\\':
			b.WriteString(`\\`)
		case '"':
			b.WriteString(`\"`)
		case '\b':
			b.WriteString(`\b`)
		case '\f':
			b.WriteString(`\f`)
		case '\n':
			b.WriteString(`\n`)
		case '\r':
			b.WriteString(`\r`)
		case '\t':
			b.WriteString(`\t`)
		case '&', '<', '>': // prefectly okay for JSON but commonly escaped
			b.WriteString(`\u00`)
			b.WriteByte(hex[r>>4])
			b.WriteByte(hex[r&0x0f])
		case '\u2028':
			b.WriteString(`\u2028`)
		case '\u2029':
			b.WriteString(`\u2029`)
		default:
			if r < ' ' {
				b.WriteString(`\u`)
				b.WriteByte(hex[r>>12])
				b.WriteByte(hex[(r>>8)&0x0f])
				b.WriteByte(hex[(r>>4)&0x0f])
				b.WriteByte(hex[r&0x0f])
			} else if r < 0x80 {
				b.WriteByte(byte(r))
			} else {
				b.WriteRune(r)
			}
		}
	}
	b.WriteByte('"')
}
