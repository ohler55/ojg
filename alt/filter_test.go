// Copyright (c) 2023, Peter Ohler, All rights reserved.

package alt_test

import (
	"testing"
	"time"

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

func TestNewFilterFlat(t *testing.T) {
	f := alt.NewFilter(map[string]any{
		"top.one":          1,
		"top.two.son":      2,
		"top.three.child":  3,
		"top.two.daughter": 22,
	})
	tt.Equal(t, `{top: {one: 1 three: {child: 3} two: {daughter: 22 son: 2}}}`, string(prettyWriter.Encode(f)))
}

func TestNewFilterNested(t *testing.T) {
	f := alt.NewFilter(map[string]any{
		"top.one": 1,
		"top": map[string]any{
			"two": map[string]any{
				"son":      2,
				"daughter": 22,
			},
			"three": map[string]any{
				"child": "x",
			},
		},
	})
	tt.Equal(t, `{top: {one: 1 three: {child: x} two: {daughter: 22 son: 2}}}`, string(prettyWriter.Encode(f)))
}

func TestFilterMatchMap(t *testing.T) {
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

	tt.Equal(t, false, f.Match(7))
}

func TestFilterMatchList1(t *testing.T) {
	f := alt.NewFilter(map[string]any{
		"a.b": 2,
		"a.c": 3,
	})
	tt.Equal(t, true, f.Match(map[string]any{
		"a": map[string]any{
			"a": 1,
			"b": []any{1, 2, 3},
			"c": 3,
		},
	}))
	tt.Equal(t, false, f.Match(map[string]any{
		"a": map[string]any{
			"a": 1,
			"b": []any{1, 3, 5},
			"c": 3,
		},
	}))
}

func TestFilterMatchList2(t *testing.T) {
	f := alt.NewFilter(map[string]any{
		"a.b": 2,
		"a.c": 3,
	})
	tt.Equal(t, true, f.Match(map[string]any{
		"a": []any{
			map[string]any{
				"x": 1,
			},
			map[string]any{
				"a": 1,
				"b": []any{1, 2, 3},
				"c": 3,
			},
			map[string]any{
				"y": 1,
			},
		},
	}))
}

func TestFilterMatchFloat(t *testing.T) {
	f := alt.NewFilter(map[string]any{
		"a.b": 2.5,
		"a.c": 3.0,
	})
	tt.Equal(t, true, f.Match(map[string]any{
		"a": map[string]any{
			"a": 1.0,
			"b": 2.5,
			"c": 3,
		},
	}))
	tt.Equal(t, false, f.Match(map[string]any{
		"a": map[string]any{
			"a": 1.0,
			"b": 2.4,
			"c": 3,
		},
	}))
}

func TestFilterMatchBool(t *testing.T) {
	f := alt.NewFilter(map[string]any{
		"a.b": true,
		"a.c": false,
	})
	tt.Equal(t, true, f.Match(map[string]any{
		"a": map[string]any{
			"a": false,
			"b": true,
			"c": false,
		},
	}))
	tt.Equal(t, false, f.Match(map[string]any{
		"a": map[string]any{
			"a": false,
			"b": false,
			"c": false,
		},
	}))
}

func TestFilterMatchString(t *testing.T) {
	f := alt.NewFilter(map[string]any{
		"a.b": "xyz",
		"a.c": "abc",
	})
	tt.Equal(t, true, f.Match(map[string]any{
		"a": map[string]any{
			"a": "zzz",
			"b": "xyz",
			"c": "abc",
		},
	}))
	tt.Equal(t, false, f.Match(map[string]any{
		"a": map[string]any{
			"a": "zzz",
			"b": "abc",
			"c": "abc",
		},
	}))
}

func TestFilterMatchSimplifier(t *testing.T) {
	f := alt.NewFilter(map[string]any{
		"type": "silly",
		"val":  3,
	})
	tt.Equal(t, true, f.Match(&silly{val: 3}))
	tt.Equal(t, false, f.Match(&silly{val: 2}))
}

func TestFilterMatchReflect(t *testing.T) {
	f := alt.NewFilter(map[string]any{
		"val":  3,
		"nest": 2,
	})
	tt.Equal(t, true, f.Match(&Dummy{Val: 3, Nest: []any{1, 2, 4}}))
	tt.Equal(t, false, f.Match(&Dummy{Val: 3, Nest: []any{1, 3, 5}}))
}

func TestFilterMatchTime(t *testing.T) {
	now := time.Now()
	f := alt.NewFilter(map[string]any{
		"when": now,
	})
	tt.Equal(t, true, f.Match(map[string]any{"when": now}))
	tt.Equal(t, false, f.Match(map[string]any{"when": time.Now()}))
}
