// Copyright (c) 2022, Peter Ohler, All rights reserved.

package tt

import (
	"fmt"
	"reflect"
	"strings"
	"testing"
)

// SameType returns true if the actual and expected values are of the same type.
func SameType(t *testing.T, expect, actual any, args ...any) (eq bool) {
	eq = reflect.TypeOf(expect) == reflect.TypeOf(actual)
	if !eq {
		var b strings.Builder
		b.WriteString(fmt.Sprintf("\nexpect: (%T) %v\nactual: (%T) %v\n", expect, expect, actual, actual))
		finishFail(t, &b, args)
	}
	return
}
