// Copyright (c) 2020, Peter Ohler, All rights reserved.

package jp_test

import (
	"sort"
	"testing"

	"github.com/ohler55/ojg/gen"
	"github.com/ohler55/ojg/jp"
	"github.com/ohler55/ojg/sen"
	"github.com/stretchr/testify/require"
)

func TestWalk(t *testing.T) {
	data := map[string]interface{}{"a": []interface{}{1, 2, 3}, "b": nil}
	var paths []string
	jp.Walk(data, func(path jp.Expr, value interface{}) { paths = append(paths, path.String()) })
	sort.Strings(paths)
	require.Equal(t, `[$ $.a "$.a[0]" "$.a[1]" "$.a[2]" $.b]`, string(sen.Bytes(paths)))
}

func TestWalkNode(t *testing.T) {
	data := gen.Object{"a": gen.Array{gen.Int(1), gen.Int(2), gen.Int(3)}, "b": nil}
	var paths []string
	jp.Walk(data, func(path jp.Expr, value interface{}) { paths = append(paths, path.String()) })
	sort.Strings(paths)
	require.Equal(t, `[$ $.a "$.a[0]" "$.a[1]" "$.a[2]" $.b]`, string(sen.Bytes(paths)))
}
