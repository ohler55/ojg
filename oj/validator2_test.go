// Copyright (c) 2020, Peter Ohler, All rights reserved.

package oj_test

import (
	"testing"

	"github.com/ohler55/ojg/oj"
	"github.com/ohler55/ojg/tt"
)

func TestValidator2ValidateString(t *testing.T) {
	for i, d := range []data{
		{src: `{
  "resourceType": "Patient",
  "id": "example",
"x":2
}`},

		{src: "null"},
		{src: "true"},
		{src: "false"},
		{src: "false \n "},
		{src: "123"},
		{src: "-321"},
		{src: "12.3"},
		{src: "0 "},
		{src: "12\n"},
		{src: "[]"},
		{src: "0\n"},
		{src: "-12.3 "},
		{src: "-12.3\n"},
		{src: "-12.3e-5"},
		{src: "12.3e+5 "},
		{src: "12.3e+5\n"},
		{src: `12345678901234567890`},
		{src: `9223372036854775807`},
		{src: `9223372036854775808`},
		{src: `-9223372036854775807`},
		{src: `-9223372036854775808`},
		{src: `0.9223372036854775808`},
		{src: `-0.9223372036854775808`},
		{src: `1.2e1025`},
		{src: `-1.2e-1025`},

		{src: "\xef\xbb\xbf\"xyz\"", value: "xyz"},

		{src: "[]"},
		{src: "[0,\ntrue , false,null]"},
		{src: `[0.1e3,"x",-1,{}]`},
		{src: "[1.2,0]"},
		{src: "[1.2e2,0.1]"},
		{src: "[1.2e2,0]"},
		{src: "[true]"},
		{src: "[true,false]"},
		{src: "[[]]"},
		{src: "[[true]]"},
		{src: `"x\t\n\"\b\f\r\u0041\\\/y"`},
		{src: `"x\u004a\u004Ay"`},
		{src: "\"bass \U0001D122\""},
		{src: `""`},
		{src: `[1,"a\tb"]`},
		{src: `{"a\tb":1}`},
		{src: `{"x":1,"a\tb":2}`},
		{src: "[0\n,3\n,5.0e2\n]"},

		{src: "{}"},
		{src: `{"abc":true}`},
		{src: "{\"z\":0,\n\"z2\":0}"},
		{src: `{"z":1.2,"z2":0}`},
		{src: `{"abc":{"def":3}}`},
		{src: `{"x":1.2e3,"y":true}`},
		{src: `{"abc": [{"x": {"y": [{"b": true}]},"z": 7}]}`},

		{src: "null {}"},

		{src: "{}}", expect: "too many closes at 1:3"},
		{src: "{ \n", expect: "incomplete JSON at 2:1"},
		{src: "{]}", expect: "expected a string start or object close, not ']' at 1:2"},
		{src: "[}]", expect: "unexpected object close at 1:2"},
		{src: "{\"a\" \n : 1]}", expect: "unexpected array close at 2:5"},
		{src: `[1}]`, expect: "unexpected object close at 1:3"},
		{src: `1]`, expect: "too many closes at 1:2"},
		{src: `1}`, expect: "too many closes at 1:2"},
		{src: `]`, expect: "too many closes at 1:1"},
		{src: `x`, expect: "unexpected character 'x' at 1:1"},
		{src: `[1,]`, expect: "unexpected character ']' at 1:4"},
		{src: `[null x`, expect: "expected a comma or close, not 'x' at 1:7"},
		{src: "{\n\"x\":1 ]", expect: "unexpected array close at 2:7"},
		{src: `[1 }`, expect: "unexpected object close at 1:4"},
		{src: "{\n\"x\":1,}", expect: "expected a string start, not '}' at 2:7"},
		{src: `{"x"x}`, expect: "expected a colon, not 'x' at 1:5"},
		{src: `nuul`, expect: "expected null at 1:3"},
		{src: `fasle`, expect: "expected false at 1:3"},
		{src: `ture`, expect: "expected true at 1:2"},
		{src: `[0,nuul]`, expect: "expected null at 1:6"},
		{src: `[0,fail]`, expect: "expected false at 1:6"},
		{src: `-x`, expect: "invalid number at 1:2"},
		{src: `0]`, expect: "too many closes at 1:2"},
		{src: `0}`, expect: "too many closes at 1:2"},
		{src: `0x`, expect: "invalid number at 1:2"},
		{src: `1x`, expect: "invalid number at 1:2"},
		{src: `1.x`, expect: "invalid number at 1:3"},
		{src: `1.2x`, expect: "invalid number at 1:4"},
		{src: `1.2ex`, expect: "invalid number at 1:5"},
		{src: `1.2e+x`, expect: "invalid number at 1:6"},
		{src: "1.2\n]", expect: "too many closes at 2:1"},
		{src: `1.2]`, expect: "too many closes at 1:4"},
		{src: `1.2}`, expect: "too many closes at 1:4"},
		{src: `1.2e2]`, expect: "too many closes at 1:6"},
		{src: `1.2e2}`, expect: "too many closes at 1:6"},
		{src: `1.2e2x`, expect: "invalid number at 1:6"},
		{src: "\"x\ty\"", expect: "invalid JSON character 0x09 at 1:3"},
		{src: `"x\zy"`, expect: "invalid JSON escape character '\\z' at 1:4"},
		{src: `"x\u004z"`, expect: "invalid JSON unicode character 'z' at 1:8"},
		{src: "\xef\xbb[]", expect: "expected BOM at 1:3"},
		{src: "null \n {}", expect: "extra characters after close, '{' at 2:2", onlyOne: true},

		{src: "[ // a comment\n  true\n]"},
		{src: "[ // a comment\n  true\n]", expect: "comments not allowed at 1:3", noComment: true},
		{src: "[\n  null, // a comment\n  true\n]"},
		{src: "[\n  null, / a comment\n  true\n]", expect: "unexpected character ' ' at 2:10", noComment: false},
		//{src: "[\n  null, // a comment\n  true\n]", expect: "comments not allowed at 2:9", noComment: true},
	} {
		var err error
		if d.onlyOne || d.noComment {
			p := oj.Validator2{OnlyOne: d.onlyOne, NoComment: d.noComment}
			err = p.Validate([]byte(d.src))
		} else {
			err = oj.Validate([]byte(d.src))
		}
		if 0 < len(d.expect) {
			tt.NotNil(t, err, d.src)
			tt.Equal(t, d.expect, err.Error(), i, ": ", d.src)
		} else {
			tt.Nil(t, err, i, ": ", d.src)
		}
	}
}
