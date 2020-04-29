// Copyright (c) 2020, Peter Ohler, All rights reserved.

package ojg_test

import (
	"testing"

	"github.com/ohler55/ojg"
	"github.com/ohler55/ojg/tt"
)

type data struct {
	src string
	// Empty means no error expected while non empty should be compared
	// err.Error().
	expect    string
	value     interface{}
	onlyOne   bool
	noComment bool
}

func TestValidateString(t *testing.T) {
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

func TestParseString(t *testing.T) {
	for _, d := range []data{

		{src: "null", expect: "", value: nil},
		{src: "true", expect: "", value: true},
		{src: "false", expect: "", value: false},
		{src: "123", expect: "", value: 123},
		{src: "-321", expect: "", value: -321},
		{src: "12.3", expect: "", value: 12.3},
		{src: `"xyz"`, expect: "", value: "xyz"},

		{src: "[]", expect: "", value: []interface{}{}},
		{src: "[true]", expect: "", value: []interface{}{true}},
		{src: "[true,false]", expect: "", value: []interface{}{true, false}},
		{src: "[[]]", expect: "", value: []interface{}{[]interface{}{}}},
		{src: "[[true]]", expect: "", value: []interface{}{[]interface{}{true}}},

		{src: "{}", expect: "", value: map[string]interface{}{}},
		{src: `{"abc":true}`, expect: "", value: map[string]interface{}{"abc": true}},
		{src: `{"abc":{"def":3}}`, expect: "", value: map[string]interface{}{"abc": map[string]interface{}{"def": 3}}},
	} {
		var err error
		var v interface{}
		if d.onlyOne || d.noComment {
			p := ojg.Parser{OnlyOne: d.onlyOne, NoComment: d.noComment}
			v, err = p.Parse(d.src)
		} else {
			v, err = ojg.Parse(d.src)
		}
		if 0 < len(d.expect) {
			tt.NotNil(t, err, d.src)
			tt.Equal(t, d.expect, err.Error(), d.src)
		} else {
			tt.Nil(t, err, d.src)
			tt.Equal(t, d.value, v, d.src)
		}
	}
}
