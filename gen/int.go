// Copyright (c) 2020, Peter Ohler, All rights reserved.

package gen

import (
	"strconv"
)

type Int int64

func (n Int) String() string {
	return strconv.FormatInt(int64(n), 10)
}

func (n Int) Alter() interface{} {
	return int64(n)
}

func (n Int) Simplify() interface{} {
	return int64(n)
}

func (n Int) Dup() Node {
	return n
}

func (n Int) Empty() bool {
	return false
}
