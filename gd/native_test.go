// Copyright (c) 2020, Peter Ohler, All rights reserved.

package gd_test

import (
	"testing"
	"time"

	"github.com/ohler55/ojg/gd"
	"github.com/ohler55/ojg/tt"

	"gitlab.com/uhn/core/pkg/tree"
)

func TestAlterNative(t *testing.T) {
	gd.Sort = true
	tm := time.Date(2020, time.April, 12, 16, 34, 04, 123456789, time.UTC)
	native := map[string]interface{}{
		"a": []interface{}{1, 2, true, tm},
		"b": 2.3,
		"c": map[string]interface{}{
			"x": "xxx",
		},
	}
	n, err := gd.AlterNative(native)
	tt.Nil(t, err)
	tt.Equal(t, `{"a":[1,2,true,"2020-04-12T16:34:04.123456789Z"],"b":2.3,"c":{"x":"xxx"}}`, n.String())
}

func BenchmarkAlterNative(b *testing.B) {
	t := time.Date(2020, time.April, 12, 16, 34, 04, 123456789, time.UTC)
	for n := 0; n < b.N; n++ {
		native := map[string]interface{}{
			"a": []interface{}{1, 2, true, t},
			"b": 2.3,
			"c": map[string]interface{}{
				"x": "xxx",
			},
		}
		_, _ = gd.AlterNative(native)
	}
}

func BenchmarkNative(b *testing.B) {
	t := time.Date(2020, time.April, 12, 16, 34, 04, 123456789, time.UTC)
	for n := 0; n < b.N; n++ {
		native := map[string]interface{}{
			"a": []interface{}{1, 2, true, t},
			"b": 2.3,
			"c": map[string]interface{}{
				"x": "xxx",
			},
		}
		_, _ = gd.FromNative(native)
	}
}

func BenchmarkTree(b *testing.B) {
	t := time.Date(2020, time.April, 12, 16, 34, 04, 123456789, time.UTC)
	for n := 0; n < b.N; n++ {
		native := map[string]interface{}{
			"a": []interface{}{1, 2, true, t},
			"b": 2.3,
			"c": map[string]interface{}{
				"x": "xxx",
			},
		}
		_, _ = tree.FromNative(native)
	}
}
