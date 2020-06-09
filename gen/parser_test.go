// Copyright (c) 2020, Peter Ohler, All rights reserved.

package gen_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/ohler55/ojg/gen"
	"github.com/ohler55/ojg/tt"
)

const callbackJSON = `
1
[2]
{"x":3}
true false 123
`

func TestParseString(t *testing.T) {
	for i, d := range []data{
		{src: "null", value: nil},
		{src: "true", value: true},
		{src: "false", value: false},
		{src: "123", value: 123},
		{src: "-321", value: -321},
		{src: "12.3", value: 12.3},
		{src: "-12.345e2", value: -1234.5},
		{src: "12.5e-1", value: 1.25},
		{src: `12345678901234567890`, value: gen.Big("12345678901234567890")},
		{src: `9223372036854775807`, value: 9223372036854775807},              // max int
		{src: `9223372036854775808`, value: gen.Big("9223372036854775808")},   // max int + 1
		{src: `-9223372036854775807`, value: -9223372036854775807},            // min int
		{src: `-9223372036854775808`, value: gen.Big("-9223372036854775808")}, // min int -1
		{src: `0.9223372036854775808`, value: gen.Big("0.9223372036854775808")},
		{src: `123456789012345678901234567890`, value: gen.Big("123456789012345678901234567890")},
		{src: `0.123456789012345678901234567890`, value: gen.Big("0.123456789012345678901234567890")},
		{src: `[12345678901234567890,12345678901234567891]`,
			value: gen.Array{gen.Big("12345678901234567890"), gen.Big("12345678901234567891")}},
		{src: `0.1e20000`, value: gen.Big("0.1e20000")},
		{src: `1.2e1025`, value: gen.Big("1.2e1025")},
		{src: `-1.2e-1025`, value: gen.Big("-1.2e-1025")},

		{src: `"xyz"`, value: "xyz"},

		{src: "[]", value: []interface{}{}},
		{src: "[true]", value: []interface{}{true}},
		{src: "[true,false]", value: []interface{}{true, false}},
		{src: "[[]]", value: []interface{}{[]interface{}{}}},
		{src: "[[true]]", value: []interface{}{[]interface{}{true}}},

		{src: "{}", value: map[string]interface{}{}},
		{src: `{"abc":true}`, value: map[string]interface{}{"abc": true}},
		{src: `{"abc":{"def":3}}`, value: map[string]interface{}{"abc": map[string]interface{}{"def": 3}}},

		{src: `{"abc": [{"x": {"y": [{"b": true}]},"z": 7}]}`,
			value: map[string]interface{}{
				"abc": []interface{}{
					map[string]interface{}{
						"x": map[string]interface{}{
							"y": []interface{}{
								map[string]interface{}{
									"b": true,
								},
							},
						},
						"z": 7,
					},
				},
			}},
		{src: "{}}", expect: "extra characters after close, '}' at 1:3"},
		{src: "{]}", expect: "expected a string start or object close, not ']' at 1:2"},
		{src: "[}]", expect: "unexpected object close at 1:2"},
		{src: `{"a":1]}`, expect: "unexpected array close at 1:7"},
		{src: `[1}]`, expect: "unexpected object close at 1:3"},
		{src: `1]`, expect: "too many closes at 1:2"},
		{src: `1}`, expect: "too many closes at 1:2"},

		{src: "[\n  null, // a comment\n  true\n]", value: []interface{}{nil, true}, noComment: false},
		{src: "[\n  null, // a comment\n  true\n]", expect: "comments not allowed at 2:9", noComment: true},
		//		{src: "[null // a comment\n]", expect: "xxx", noComment: true},
	} {
		if testing.Verbose() {
			fmt.Printf("... %s\n", d.src)
		}
		var err error
		var v interface{}
		if d.onlyOne || d.noComment {
			p := gen.Parser{NoComment: d.noComment}
			v, err = p.Parse([]byte(d.src))
		} else {
			var p gen.Parser
			v, err = p.Parse([]byte(d.src))
		}
		if 0 < len(d.expect) {
			tt.NotNil(t, err, d.src)
			tt.Equal(t, d.expect, err.Error(), i, ": ", d.src)
		} else {
			tt.Nil(t, err, d.src)
			tt.Equal(t, d.value, v, ": ", d.src)
		}
	}
}

func TestParseCallback(t *testing.T) {
	var results []byte
	cb := func(n gen.Node) bool {
		if 0 < len(results) {
			results = append(results, ' ')
		}
		results = append(results, n.String()...)
		return false
	}
	var p gen.Parser
	v, err := p.Parse([]byte(callbackJSON), cb, true)
	tt.Nil(t, err)
	tt.Nil(t, v)
	tt.Equal(t, `1 [2] {"x":3} true false 123`, string(results))

	results = results[:0]
	v, err = p.Parse([]byte(callbackJSON), cb)
	tt.Nil(t, err)
	tt.Nil(t, v)
	tt.Equal(t, `1 [2] {"x":3} true false 123`, string(results))
}

func TestParseReaderCallback(t *testing.T) {
	var results []byte
	cb := func(n gen.Node) bool {
		if 0 < len(results) {
			results = append(results, ' ')
		}
		results = append(results, n.String()...)
		return false
	}
	var p gen.Parser
	v, err := p.ParseReader(strings.NewReader(callbackJSON), cb, true)
	tt.Nil(t, err)
	tt.Nil(t, v)
	tt.Equal(t, `1 [2] {"x":3} true false 123`, string(results))

	results = results[:0]
	v, err = p.ParseReader(strings.NewReader(callbackJSON), cb, true)
	tt.Nil(t, err)
	tt.Nil(t, v)
	tt.Equal(t, `1 [2] {"x":3} true false 123`, string(results))
}
