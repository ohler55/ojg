// Copyright (c) 2024, Peter Ohler, All rights reserved.
package sen_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/ohler55/ojg/jp"
	"github.com/ohler55/ojg/pretty"
	"github.com/ohler55/ojg/sen"
	"github.com/ohler55/ojg/tt"
)

func TestMatch(t *testing.T) {
	var buf []byte
	err := sen.Match([]byte(`{a:1 b:2}`), func(path jp.Expr, data any) {
		buf = fmt.Appendf(buf, "%s: %v", path, pretty.SEN(data))
	}, jp.C("a"))
	tt.Nil(t, err)
	tt.Equal(t, "$.a: 1", string(buf))
}

func TestMatchString(t *testing.T) {
	var buf []byte
	err := sen.MatchString(`{a:1 b:2}`, func(path jp.Expr, data any) {
		buf = fmt.Appendf(buf, "%s: %v", path, pretty.SEN(data))
	}, jp.C("a"))
	tt.Nil(t, err)
	tt.Equal(t, "$.a: 1", string(buf))
}

func TestMatchLoad(t *testing.T) {
	var buf []byte
	err := sen.MatchLoad(strings.NewReader(`{a:1 b:2}`), func(path jp.Expr, data any) {
		buf = fmt.Appendf(buf, "%s: %v", path, pretty.SEN(data))
	}, jp.C("a"))
	tt.Nil(t, err)
	tt.Equal(t, "$.a: 1", string(buf))
}
