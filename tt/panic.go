// Copyright (c) 2020, Peter Ohler, All rights reserved.

package tt

import (
	"strings"
	"testing"
)

// Panic verifies that a function panics..
func Panic(t *testing.T, fun func(), args ...interface{}) {
	ff := func() {
		var b strings.Builder
		b.WriteString("\nexpect: panic\nactual: no panic\n")
		finishFail(t, &b, args)
	}
	defer func() {
		if r := recover(); r == nil {
			ff()
		}
	}()
	fun()
	ff()
}
