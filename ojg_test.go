// Copyright (c) 2020, Peter Ohler, All rights reserved.

package ojg_test

import (
	"testing"

	"github.com/ohler55/ojg"
	"github.com/ohler55/ojg/tt"
)

func TestValidateString(t *testing.T) {
	type data struct {
		src string
		// Empty means no error expected while non empty should be compared
		// err.Error().
		expect    string
		onlyOne   bool
		noComment bool
	}
	for _, d := range []data{
		{src: "null", expect: ""},
		{src: "true", expect: ""},
		{src: "false", expect: ""},
		{src: "[]", expect: ""},
		{src: "[true]", expect: ""},
		{src: "[true,false]", expect: ""},
		{src: "[[],[true],false]", expect: ""},
		{src: "[[],[true]false]", expect: "expected a comma or close, not 'f' at 1:11"},

		{src: "123", expect: ""},
		{src: "-1.23", expect: ""},
		{src: "[1,2]", expect: ""},

		{src: "[]", expect: ""},
		{src: "null {}", expect: ""},
		{src: "null {}", expect: "extra characters after close, '{' at 1:6", onlyOne: true},

		{src: "-1.23", expect: ""},
		{src: "+1.23", expect: "unexpected character '+' at 1:1"},
		{src: "1.23e+3", expect: ""},
		{src: "1.23e-3", expect: ""},
		{src: "1.23e3", expect: ""},
		{src: "1.2e3e3", expect: "invalid number '1.2e3e' at 1:6"},
		{src: "0.3", expect: ""},
		{src: "03", expect: "invalid number '03' at 1:2"},

		{src: `""`, expect: ""},
		{src: `"abc"`, expect: ""},
		{src: `"a\tb\nc\b\"\\d\f\r"`, expect: ""},
		{src: "\"bass \U0001D122\"", expect: ""},
		{src: `"a \u2669"`, expect: ""},
		{src: `"bad \uabcz"`, expect: "invalid JSON unicode character 'z' at 1:11"},

		{src: "[\n  // a comment\n]", expect: ""},
		{src: "[\n  // a comment\n]", expect: "comments not allowed at 2:3", noComment: true},
		{src: "[\n  / a comment\n]", expect: "unexpected character ' ' at 2:4"},

		{src: "{}", expect: ""},
		{src: `{"a":3}`, expect: ""},
		{src: `{"a": 3, "b": true}`, expect: ""},
		{src: `{"a":{"b":{"c":true}}}`, expect: ""},
		{src: `{x}`, expect: "expected a string start or object close, not 'x' at 1:2"},
	} {
		var err error
		if d.onlyOne || d.noComment {
			p := ojg.Parser{OnlyOne: d.onlyOne, NoComment: d.noComment}
			err = p.Validate(d.src)
		} else {
			err = ojg.Validate(d.src)
		}
		if 0 < len(d.expect) {
			tt.NotNil(t, err, d.src)
			tt.Equal(t, d.expect, err.Error(), d.src)
		} else {
			tt.Nil(t, err, d.src)
		}
	}
}
