// Copyright (c) 2020, Peter Ohler, All rights reserved.

package ojg_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/ohler55/ojg"
	"github.com/ohler55/ojg/tt"
)

func TestValidateString(t *testing.T) {
	for i, d := range []data{
		{src: "null"},
		{src: "true"},
		{src: "false"},
		{src: "[]"},
		{src: "[true]"},
		{src: "[true,false]"},
		{src: "[[],[true],false]"},
		{src: "[[],[true]false]", expect: "expected a comma or close, not 'f' at 1:11"},

		{src: "123"},
		{src: "-1.23"},
		{src: "[1,2]"},

		{src: "[]"},
		{src: "null {}"},
		{src: "null {}", expect: "extra characters after close, '{' at 1:6", onlyOne: true},

		{src: "-1.23"},
		{src: "+1.23", expect: "unexpected character '+' at 1:1"},
		{src: "1.23e+3"},
		{src: "1.23e-3"},
		{src: "1.23e3"},
		{src: "1.2e3e3", expect: "invalid number at 1:6"},
		{src: "0.3"},
		{src: "03", expect: "invalid number at 1:2"},

		{src: `""`},
		{src: `"abc"`},
		{src: `"a\tb\nc\b\"\\d\f\r"`},
		{src: "\"bass \U0001D122\""},
		{src: `"a \u2669"`},
		{src: `"bad \uabcz"`, expect: "invalid JSON unicode character 'z' at 1:11"},

		{src: "[\n  // a comment\n]"},
		{src: "[\n  // a comment\n]", expect: "comments not allowed at 2:3", noComment: true},
		{src: "[\n  / a comment\n]", expect: "unexpected character ' ' at 2:4"},

		{src: "{}"},
		{src: `{"a":3}`},
		{src: `{"a": 3, "b": true}`},
		{src: `{"a":{"b":{"c":true}}}`},
		{src: ` { "a"
:
true
}`},
		{src: `{x}`, expect: "expected a string start or object close, not 'x' at 1:2"},
		{src: "{}}", expect: "too many closes at 1:3"},
		{src: "[]]", expect: "too many closes at 1:3"},
		{src: `{"x":2]`, expect: "unexpected array close at 1:7"},
		{src: "[}", expect: "unexpected object close at 1:2"},
	} {
		var err error
		if d.onlyOne || d.noComment {
			p := ojg.Validator{OnlyOne: d.onlyOne, NoComment: d.noComment}
			err = p.Validate([]byte(d.src))
		} else {
			err = ojg.Validate([]byte(d.src))
		}
		if 0 < len(d.expect) {
			tt.NotNil(t, err, d.src)
			tt.Equal(t, d.expect, err.Error(), i, ": ", d.src)
		} else {
			tt.Nil(t, err, i, ": ", d.src)
		}
	}
}

func TestValidateReader(t *testing.T) {
	r := strings.NewReader("[true,[false,[null],123],456]")
	err := ojg.ValidateReader(r)
	tt.Nil(t, err)

	var buf []byte
	buf = append(buf, "[\n"...)
	for i := 0; i < 1000; i++ {
		buf = append(buf, "  true,\n"...)
	}
	buf = append(buf, "  false\n]\n"...)
	br := bytes.NewReader(buf)
	err = ojg.ValidateReader(br)
	tt.Nil(t, err)
}

func TestValidateResuse(t *testing.T) {
	var v ojg.Validator
	err := v.Validate([]byte("[true,[false,[null],123],456]"))
	tt.Nil(t, err)
	// a second time
	err = v.Validate([]byte("[true,[false,[null],123],456]"))
	tt.Nil(t, err)

	r := strings.NewReader("[true,[false,[null],123],456]")
	err = v.ValidateReader(r)
	tt.Nil(t, err)
}

func TestValidateBOM(t *testing.T) {
	var v ojg.Validator
	err := v.Validate([]byte("\xef\xbb\xbf[true]"))
	tt.Nil(t, err)
}
