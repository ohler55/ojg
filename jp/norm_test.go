// Copyright (c) 2025, Peter Ohler, All rights reserved.

package jp

import "testing"

func TestNorm(t *testing.T) {
	n := norm('x')
	// Just make sure there are no panics.
	_ = n.locate(R(), nil, C("a"), 1)
	n.Walk(C("a"), R().C("a"), []any{map[string]any{"a": 1}}, func(path Expr, nodes []any) {})
}
