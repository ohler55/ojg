// Copyright (c) 2020, Peter Ohler, All rights reserved.

package alt

import (
	"strings"

	"github.com/ohler55/ojg/gen"
)

// Bool convert the value provided to a bool. If conversion is not possible
// such as if the provided value is an array then the first option default
// value is returned or if not provided false is returned. If the type is not
// a bool nor a gen.Bool and there is a second optional default then that
// second default value is returned. This approach keeps the return as a
// single value and gives the caller the choice of how to indicate a bad
// value.
func Bool(v interface{}, defaults ...bool) (b bool) {
	switch tv := v.(type) {
	case nil:
		if 1 < len(defaults) {
			b = defaults[1]
		}
	case bool:
		b = tv
	case string:
		if 1 < len(defaults) {
			b = defaults[1]
		} else if strings.EqualFold(tv, "true") {
			b = true
		} else if strings.EqualFold(tv, "false") {
			b = false
		} else if 0 < len(defaults) {
			b = defaults[0]
		}
	case gen.Bool:
		b = bool(tv)
	case gen.String:
		if 1 < len(defaults) {
			b = defaults[1]
		} else if strings.EqualFold(string(tv), "true") {
			b = true
		} else if strings.EqualFold(string(tv), "false") {
			b = false
		} else if 0 < len(defaults) {
			b = defaults[0]
		}
	default:
		if 0 < len(defaults) {
			b = defaults[0]
		}
	}
	return
}
