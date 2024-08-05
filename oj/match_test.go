// Copyright (c) 2024, Peter Ohler, All rights reserved.

package oj_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/ohler55/ojg/jp"
	"github.com/ohler55/ojg/oj"
	"github.com/ohler55/ojg/pretty"
	"github.com/ohler55/ojg/tt"
)

func TestMatch(t *testing.T) {
	var buf []byte
	err := oj.Match([]byte(`{"a":1, "b":2}`), jp.C("a"), func(path jp.Expr, data any) {
		buf = fmt.Appendf(buf, "%s: %v", path, pretty.SEN(data))
	})
	tt.Nil(t, err)
	tt.Equal(t, "$.a: 1", string(buf))
}

func TestMatchString(t *testing.T) {
	var buf []byte
	err := oj.MatchString(`{"a":1, "b":2}`, jp.C("a"), func(path jp.Expr, data any) {
		buf = fmt.Appendf(buf, "%s: %v", path, pretty.SEN(data))
	})
	tt.Nil(t, err)
	tt.Equal(t, "$.a: 1", string(buf))
}

func TestMatchLoad(t *testing.T) {
	var buf []byte
	err := oj.MatchLoad(strings.NewReader(`{"a":1, "b":2}`), jp.C("a"), func(path jp.Expr, data any) {
		buf = fmt.Appendf(buf, "%s: %v", path, pretty.SEN(data))
	})
	tt.Nil(t, err)
	tt.Equal(t, "$.a: 1", string(buf))
}
