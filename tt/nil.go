// Copyright (c) 2020, Peter Ohler, All rights reserved.

package tt

import (
	"fmt"
	"strings"
	"testing"
)

// Nil check.
func Nil(t *testing.T, actual interface{}, args ...interface{}) {
	if !isNil(actual) {
		var b strings.Builder
		b.WriteString(fmt.Sprintf("\nexpect: nil\nactual: (%T) %v\n", actual, actual))
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
}

// NotNil check.
func NotNil(t *testing.T, actual interface{}, args ...interface{}) {
	if isNil(actual) {
		var b strings.Builder
		b.WriteString("\nexpect: not nil\nactual: nil\n")
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
}
