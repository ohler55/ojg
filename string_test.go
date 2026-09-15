// Copyright (c) 2021, Peter Ohler, All rights reserved.

package ojg_test

import (
	"reflect"
	"testing"

	"github.com/ohler55/ojg"
	"github.com/ohler55/ojg/tt"
)

func TestStringJSON(t *testing.T) {
	type Data struct {
		src      string
		expect   string
		htmlSafe bool
	}
	for i, td := range []*Data{
		{src: "abc", expect: `"abc"`},
		{src: "a\tbc", expect: `"a\tbc"`},
		{src: "a<b>c", expect: `"a<b>c"`},
		{src: "a<b>c", expect: `"a\u003cb\u003ec"`, htmlSafe: true},
		{src: "a 𝄢 note", expect: `"a 𝄢 note"`},
		{src: "a\u001ec", expect: `"a\u001ec"`},
		{src: "a\u2028b\u2029c", expect: `"a\u2028b\u2029c"`},
		{src: "abc\ufffd", expect: `"abc\ufffd"`},
	} {
		var buf []byte
		buf = ojg.AppendJSONString(buf, td.src, td.htmlSafe)
		tt.Equal(t, td.expect, string(buf), i, ": ", td.src)
	}
}

func TestStringSEN(t *testing.T) {
	type Data struct {
		src      string
		expect   string
		htmlSafe bool
	}
	for i, td := range []*Data{
		{src: "abc", expect: `abc`},
		{src: "", expect: `""`},
		{src: "a\tbc", expect: "\"a\tbc\""},
		{src: "a\"bc", expect: `"a\"bc"`},
		{src: "a<b>c", expect: `a<b>c`},
		{src: "a<b>c", expect: `"a\u003cb\u003ec"`, htmlSafe: true},
		{src: "a 𝄢 note", expect: `"a 𝄢 note"`},
		{src: "a\u001ec", expect: `"a\u001ec"`},
		{src: "a\u2028b\u2029c", expect: `"a\u2028b\u2029c"`},
		{src: "abc\ufffd", expect: `"abc\ufffd"`},
	} {
		var buf []byte
		buf = ojg.AppendSENString(buf, td.src, td.htmlSafe)
		tt.Equal(t, td.expect, string(buf), i, ": ", td.src)
	}
}

func TestKeyString(t *testing.T) {
	tt.Equal(t, "quux", ojg.KeyString(reflect.ValueOf("quux")))
	tt.Equal(t, "123", ojg.KeyString(reflect.ValueOf(123)))
}
