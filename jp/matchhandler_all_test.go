package jp_test

import (
	"testing"

	"github.com/ohler55/ojg/jp"
	"github.com/ohler55/ojg/sen"
)

func TestMatchHandlerFilterAllMatches(t *testing.T) {
	// Every element a filter matches is reported, with its own path and value.
	got := map[string]any{}
	h := jp.NewMatchHandler(func(path jp.Expr, data any) {
		got[path.String()] = data
	}, jp.MustParseString("$[?(@ != 1)]"))
	if err := sen.TokenizeString(`["a" "b" "c"]`, h); err != nil {
		t.Fatal(err)
	}
	want := map[string]any{"$[0]": "a", "$[1]": "b", "$[2]": "c"}
	if len(got) != len(want) {
		t.Fatalf("callbacks = %v, want %v", got, want)
	}
	for p, v := range want {
		if got[p] != v {
			t.Errorf("%s = %v, want %v", p, got[p], v)
		}
	}
}
