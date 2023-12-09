// Copyright (c) 2023, Peter Ohler, All rights reserved.

package jp_test

import (
	"fmt"
	"testing"

	"github.com/ohler55/ojg/jp"
)

func TestLocateDev(t *testing.T) {
	data := []any{map[string]any{"b": 1, "c": 2}, []any{1, 2, 3}}
	x := jp.MustParseString("$[?(@[1] == 2)].*")
	for _, ep := range x.Locate(data, 0) {
		fmt.Printf("*** %s\n", ep.BracketString())
	}

}
