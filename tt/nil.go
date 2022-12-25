// Copyright (c) 2020, Peter Ohler, All rights reserved.

package tt

import (
	"fmt"
	"strings"
	"testing"
)

// Nil check.
func Nil(t *testing.T, actual any, args ...any) {
	if !isNil(actual) {
		var b strings.Builder
		b.WriteString(fmt.Sprintf("\nexpect: nil\nactual: (%T) %v\n", actual, actual))
		finishFail(t, &b, args)
	}
}

// NotNil check.
func NotNil(t *testing.T, actual any, args ...any) {
	if isNil(actual) {
		var b strings.Builder
		b.WriteString("\nexpect: not nil\nactual: nil\n")
		finishFail(t, &b, args)
	}
}
