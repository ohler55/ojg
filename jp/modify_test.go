// Copyright (c) 2020, Peter Ohler, All rights reserved.

package jp_test

import (
	"sort"
	"testing"

	"github.com/ohler55/ojg/jp"
	"github.com/ohler55/ojg/tt"
)

func TestExprModify(t *testing.T) {
	x := jp.MustParseString("[*]")

	data := []any{[]any{1, 3, 2, 4}, []any{4, 3, 2, 1}}
	result, err := x.Modify(data, func(element any) (any, bool) {
		if a, ok := element.([]any); ok {
			sort.Slice(a, func(i, j int) bool { return a[i].(int) < a[j].(int) })
			element = a
		}
		return element, true
	})
	tt.Nil(t, err)
	tt.Equal(t, "[[1 2 3 4] [1 2 3 4]]", string(pw.Encode(result)))
}

func TestExprMustModify(t *testing.T) {
	x := jp.MustParseString("[*]")

	data := []any{[]any{1, 3, 2, 4}, []any{4, 3, 2, 1}}
	result := x.MustModify(data, func(element any) (any, bool) {
		if a, ok := element.([]any); ok {
			sort.Slice(a, func(i, j int) bool { return a[i].(int) < a[j].(int) })
			element = a
		}
		return element, true
	})
	tt.Equal(t, "[[1 2 3 4] [1 2 3 4]]", string(pw.Encode(result)))
}

func TestExprModifyOne(t *testing.T) {
	x := jp.MustParseString("[*]")

	data := []any{[]any{1, 3, 2, 4}, []any{4, 3, 2, 1}}
	result, err := x.ModifyOne(data, func(element any) (any, bool) {
		if a, ok := element.([]any); ok {
			sort.Slice(a, func(i, j int) bool { return a[i].(int) < a[j].(int) })
			element = a
		}
		return element, true
	})
	tt.Nil(t, err)
	tt.Equal(t, "[[1 2 3 4] [4 3 2 1]]", string(pw.Encode(result)))
}

func TestExprMustModifyOne(t *testing.T) {
	x := jp.MustParseString("[*]")

	data := []any{[]any{1, 3, 2, 4}, []any{4, 3, 2, 1}}
	result := x.MustModifyOne(data, func(element any) (any, bool) {
		if a, ok := element.([]any); ok {
			sort.Slice(a, func(i, j int) bool { return a[i].(int) < a[j].(int) })
			element = a
		}
		return element, true
	})
	tt.Equal(t, "[[1 2 3 4] [4 3 2 1]]", string(pw.Encode(result)))
}

func TestExprModifyOneEmpty(t *testing.T) {
	x := jp.MustParseString("")

	data := []any{}
	_, err := x.ModifyOne(data, func(element any) (any, bool) {
		return element, false
	})
	tt.NotNil(t, err)
}

func TestExprModifyEmpty(t *testing.T) {
	x := jp.MustParseString("")

	data := []any{}
	_, err := x.Modify(data, func(element any) (any, bool) {
		return element, false
	})
	tt.NotNil(t, err)
}
