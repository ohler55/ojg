// Copyright (c) 2023, Peter Ohler, All rights reserved.

package alt_test

import (
	"testing"

	"github.com/ohler55/ojg"
	"github.com/ohler55/ojg/alt"
	"github.com/ohler55/ojg/pretty"
	"github.com/ohler55/ojg/tt"
)

var prettyWriter = pretty.Writer{
	Options:  ojg.Options{Sort: true},
	Width:    100,
	MaxDepth: 4,
	SEN:      true,
}

func TestNewFilter(t *testing.T) {
	f := alt.NewFilter(map[string]any{
		"top.one":          1,
		"top.two.son":      2,
		"top.three.child":  3,
		"top.two.daughter": 22,
	})
	tt.Equal(t, `{top: {one: 1 three: {child: 3} two: {daughter: 22 son: 2}}}`, string(prettyWriter.Encode(f)))
}

func TestFilterMatch(t *testing.T) {
	f := alt.NewFilter(map[string]any{
		"a.b": 2,
		"a.c": 3,
	})
	tt.Equal(t, true, f.Match(map[string]any{
		"a": map[string]any{
			"a": 1,
			"b": 2,
			"c": 3,
		},
	}))

	tt.Equal(t, false, f.Match(map[string]any{
		"a": map[string]any{
			"a": 1,
			"b": 0,
			"c": 3,
		},
	}))
}
