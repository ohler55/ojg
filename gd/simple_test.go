// Copyright (c) 2020, Peter Ohler, All rights reserved.

package gd_test

import (
	"testing"
	"time"

	"github.com/ohler55/ojg/gd"
	"github.com/ohler55/ojg/tt"
)

func TestAlterSimple(t *testing.T) {
	gd.Sort = true
	gd.TimeFormat = time.RFC3339Nano
	tm := time.Date(2020, time.April, 12, 16, 34, 04, 123456789, time.UTC)
	simple := map[string]interface{}{
		"a": []interface{}{1, 2, true, tm},
		"b": 2.3,
		"c": map[string]interface{}{
			"x": "xxx",
		},
		"d": nil,
	}
	n, err := gd.AlterSimple(simple)
	tt.Nil(t, err)
	tt.Equal(t, `{"a":[1,2,true,"2020-04-12T16:34:04.123456789Z"],"b":2.3,"c":{"x":"xxx"},"d":null}`, n.String())
}
