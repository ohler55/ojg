// Copyright (c) 2020, Peter Ohler, All rights reserved.

package gd

import (
	"strconv"
	"strings"
)

type Float float64

func (n Float) String() string {
	return strconv.FormatFloat(float64(n), 'g', -1, 64)
}

func (n Float) Alter() interface{} {
	return float64(n)
}

func (n Float) Native() interface{} {
	return float64(n)
}

func (n Float) Dup() Node {
	return n
}

func (n Float) Empty() bool {
	return false
}

func (n Float) AsBool() (Bool, bool) {
	// Not really a bool but just in case return a value but ok as false.
	return Bool(float64(n) != 0.0), false
}

func (n Float) AsInt() (Int, bool) {
	return Int(int64(n)), true
}

func (n Float) AsFloat() (Float, bool) {
	return n, true
}

func (n Float) JSON(_ ...int) string {
	return strconv.FormatFloat(float64(n), 'g', -1, 64)
}

func (n Float) BuildJSON(b *strings.Builder, _, _ int) {
	b.WriteString(strconv.FormatFloat(float64(n), 'g', -1, 64))
}
