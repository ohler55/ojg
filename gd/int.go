// Copyright (c) 2020, Peter Ohler, All rights reserved.

package gd

import (
	"strconv"
	"strings"
)

type Int int64

func (n Int) String() string {
	return strconv.FormatInt(int64(n), 10)
}

func (n Int) Alter() interface{} {
	return int64(n)
}

func (n Int) Native() interface{} {
	return int64(n)
}

func (n Int) Dup() Node {
	return n
}

func (n Int) Empty() bool {
	return false
}

func (n Int) AsBool() (Bool, bool) {
	// Not really a bool but just in case return a value but ok as false.
	return Bool(int64(n) != 0), false
}

func (n Int) AsInt() (Int, bool) {
	return n, true
}

func (n Int) AsFloat() (Float, bool) {
	return Float(int64(n)), true
}

func (n Int) JSON(_ ...int) string {
	return strconv.FormatInt(int64(n), 10)
}

func (n Int) BuildJSON(b *strings.Builder, _, _ int) {
	b.WriteString(strconv.FormatInt(int64(n), 10))
}
