// Copyright (c) 2020, Peter Ohler, All rights reserved.

package gen_test

import (
	"testing"
	"time"

	"github.com/ohler55/ojg/gen"
	"github.com/ohler55/ojg/tt"
)

func TestAlterSimple(t *testing.T) {
	gen.Sort = true
	gen.TimeFormat = time.RFC3339Nano
	tm := time.Date(2020, time.April, 12, 16, 34, 04, 123456789, time.UTC)
	simple := map[string]interface{}{
		"a": []interface{}{1, 2, true, tm},
		"b": 2.3,
		"c": map[string]interface{}{
			"x": "xxx",
		},
		"d": nil,
	}
	n := gen.Alter(simple)
	tt.Equal(t, `{"a":[1,2,true,"2020-04-12T16:34:04.123456789Z"],"b":2.3,"c":{"x":"xxx"},"d":null}`, n.String())
}
