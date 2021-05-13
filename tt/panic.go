// Copyright (c) 2020, Peter Ohler, All rights reserved.

package tt

import (
	"fmt"
	"strings"
	"testing"
)

// Panic verifies that a function panics..
func Panic(t *testing.T, fun func(), args ...interface{}) {
	defer func() {
		if r := recover(); r != nil {
			// Pass, as expected
		}
	}()
	fun()
	var b strings.Builder
	b.WriteString(fmt.Sprintf("\nexpect: panic\nactual: no panic\n"))
	stackFill(&b)
	if 0 < len(args) {
		if format, _ := args[0].(string); 0 < len(format) {
			b.WriteString(fmt.Sprintf(format, args[1:]...))
		} else {
			b.WriteString(fmt.Sprint(args...))
		}
	}
	t.Fatal(b.String())
}
