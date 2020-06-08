// Copyright (c) 2020, Peter Ohler, All rights reserved.

package gen

import (
	"strconv"
)

type Float float64

func (n Float) String() string {
	return strconv.FormatFloat(float64(n), 'g', -1, 64)
}

func (n Float) Alter() interface{} {
	return float64(n)
}

func (n Float) Simplify() interface{} {
	return float64(n)
}

func (n Float) Dup() Node {
	return n
}

func (n Float) Empty() bool {
	return false
}
