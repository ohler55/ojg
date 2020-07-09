// Copyright (c) 2020, Peter Ohler, All rights reserved.

package sen_test

import (
	"fmt"
	"strings"
	"testing"
	"testing/iotest"

	"github.com/ohler55/ojg/sen"
	"github.com/ohler55/ojg/tt"
)

const callbackSEN = `
1
[2]
{x:3}
true false 123
`

type rdata struct {
	src    string
	expect string
	value  interface{}
}

func TestParserParseString(t *testing.T) {
	for i, d := range []rdata{
		{src: "null", value: nil},
		{src: "true", value: true},
		{src: "false", value: false},
		{src: "false \n ", value: false},
		{src: "hello", value: "hello"},
		{src: "hello ", value: "hello"},
		{src: `"hello"`, value: "hello"},
		{src: "[one two]", value: []interface{}{"one", "two"}},
		{src: "123", value: 123},
		{src: "-12.3", value: -12.3},
		{src: "-12.5e-2", value: -0.125},
		{src: "0", value: 0},
		{src: "0\n ", value: 0},
		{src: "-12.3 ", value: -12.3},
		{src: "-12.3\n", value: -12.3},
		{src: "-12.3e-5", value: -12.3e-5},
		{src: "12.3e+5 ", value: 12.3e+5},
		{src: "12.3e+5\n ", value: 12.3e+5},
		{src: `12345678901234567890`, value: "12345678901234567890"},
		{src: `9223372036854775807`, value: 9223372036854775807},     // max int
		{src: `9223372036854775808`, value: "9223372036854775808"},   // max int + 1
		{src: `-9223372036854775807`, value: -9223372036854775807},   // min int
		{src: `-9223372036854775808`, value: "-9223372036854775808"}, // min int -1
		{src: `0.9223372036854775808`, value: "0.9223372036854775808"},
		{src: `-0.9223372036854775808`, value: "-0.9223372036854775808"},
		{src: `1.2e1025`, value: "1.2e1025"},
		{src: `-1.2e-1025`, value: "-1.2e-1025"},
		{src: `12345678901234567890.321e66`, value: "12345678901234567890.321e66"},
		{src: `321.12345678901234567890e66`, value: "321.12345678901234567890e66"},
		{src: `321.123e2345`, value: "321.123e2345"},

		{src: "\xef\xbb\xbf\"xyz\"", value: "xyz"},

		{src: "[]", value: []interface{}{}},
		{src: "[0,\ntrue , false,null]", value: []interface{}{0, true, false, nil}},
		{src: `[0.1e3,"x",-1,{}]`, value: []interface{}{100.0, "x", -1, map[string]interface{}{}}},
		{src: "[1.2,0]", value: []interface{}{1.2, 0}},
		{src: "[1.2e2,0.1]", value: []interface{}{1.2e2, 0.1}},
		{src: "[1.2e2,0]", value: []interface{}{1.2e2, 0}},
		{src: "[true]", value: []interface{}{true}},
		{src: "[true,false]", value: []interface{}{true, false}},
		{src: "[[]]", value: []interface{}{[]interface{}{}}},
		{src: "[[true]]", value: []interface{}{[]interface{}{true}}},
		{src: `"x\t\n\"\b\f\r\u0041\\\/y"`, value: "x\t\n\"\b\f\r\u0041\\/y"},
		{src: `"x\u004a\u004Ay"`, value: "xJJy"},

		{src: "{}", value: map[string]interface{}{}},
		{src: `{"a\tbc":true}`, value: map[string]interface{}{"a\tbc": true}},
		{src: "{\"z\":0,\n\"z2\":0}", value: map[string]interface{}{"z": 0, "z2": 0}},
		{src: `{"z":1.2,"z2":0}`, value: map[string]interface{}{"z": 1.2, "z2": 0}},
		{src: `{"abc":{"def" :3}}`, value: map[string]interface{}{"abc": map[string]interface{}{"def": 3}}},
		{src: `{"x":1.2e3,"y":true}`, value: map[string]interface{}{"x": 1200.0, "y": true}},

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
		{src: "{ \n", expect: "not closed at 2:1"},
		{src: "{}\n }", expect: "extra characters after close, '}' at 2:2"},
		{src: "{]}", expect: "expected a token followed by a ':', not ']' at 1:2"},
		{src: "[}]", expect: "unexpected object close at 1:2"},
		{src: "{\"a\" \n : 1]}", expect: "unexpected array close at 2:5"},
		{src: `[1}]`, expect: "unexpected object close at 1:3"},
		{src: `1]`, expect: "too many closes at 1:2"},
		{src: `1}`, expect: "too many closes at 1:2"},
		{src: `]`, expect: "too many closes at 1:1"},
		{src: `[null x`, expect: "not closed at 1:8"},
		{src: "{\n\"x\":1 ]", expect: "expected a token followed by a ':', not ']' at 2:7"},
		{src: `[1 }`, expect: "unexpected object close at 1:4"},
		{src: `{"x"x}`, expect: "expected a colon, not 'x' at 1:5"},
		{src: `-x`, expect: "invalid number at 1:2"},
		{src: `0]`, expect: "too many closes at 1:2"},
		{src: `0\n $`, expect: "invalid number at 1:2"},
		{src: `0}`, expect: "too many closes at 1:2"},
		{src: `0x`, expect: "invalid number at 1:2"},
		{src: `1x`, expect: "invalid number at 1:2"},
		{src: `1.x`, expect: "invalid number at 1:3"},
		{src: `1.2x`, expect: "invalid number at 1:4"},
		{src: `1.2ex`, expect: "invalid number at 1:5"},
		{src: `1.2e+x`, expect: "invalid number at 1:6"},
		{src: `1.2e`, expect: "incomplete JSON at 1:5"},
		{src: "1.2\n]", expect: "extra characters after close, ']' at 2:1"},
		{src: `1.2]`, expect: "too many closes at 1:4"},
		{src: `1.2}`, expect: "too many closes at 1:4"},
		{src: `1.2e2]`, expect: "too many closes at 1:6"},
		{src: `1.2e2}`, expect: "too many closes at 1:6"},
		{src: `1.2e2x`, expect: "invalid number at 1:6"},
		{src: "\"x\ty\"", expect: "invalid JSON character 0x09 at 1:3"},
		{src: `"x\zy"`, expect: "invalid JSON escape character '\\z' at 1:4"},
		{src: `"x\u004z"`, expect: "invalid JSON unicode character 'z' at 1:8"},
		{src: "\xef\xbb[]", expect: "expected BOM at 1:3"},
		{src: "$x", expect: "expected a token, not '$' at 1:1"},
		{src: "x]", expect: "too many closes at 1:2"},
		{src: "x}", expect: "too many closes at 1:2"},
		{src: "x$", expect: "expected a token, not '$' at 1:2"},
		{src: "{x$:1}", expect: "expected a token followed by a ':', not '$' at 1:3"},

		{src: "[0 1,2]", value: []interface{}{0, 1, 2}},
		{src: "[0[1[2]]]", value: []interface{}{0, []interface{}{1, []interface{}{2}}}},
		{src: `{aaa:0 "bbb":"one" , c:2}`, value: map[string]interface{}{"aaa": 0, "bbb": "one", "c": 2}},
		{src: "{aaa\n:one b:\ntwo}", value: map[string]interface{}{"aaa": "one", "b": "two"}},
		{src: "[abc[x]]", value: []interface{}{"abc", []interface{}{"x"}}},
		{src: "[abc{x:1}]", value: []interface{}{"abc", map[string]interface{}{"x": 1}}},
		{src: `{aaa:"bbb" "bbb":"one" , c:2}`, value: map[string]interface{}{"aaa": "bbb", "bbb": "one", "c": 2}},
		{src: `{aaa:"b\tb" x:2}`, value: map[string]interface{}{"aaa": "b\tb", "x": 2}},

		{src: "[0{x:1}]", value: []interface{}{0, map[string]interface{}{"x": 1}}},
		{src: "[1{x:1}]", value: []interface{}{1, map[string]interface{}{"x": 1}}},
		{src: "[1.5{x:1}]", value: []interface{}{1.5, map[string]interface{}{"x": 1}}},
		{src: "[1.5[1]]", value: []interface{}{1.5, []interface{}{1}}},
		{src: "[1.5e2{x:1}]", value: []interface{}{150., map[string]interface{}{"x": 1}}},
		{src: "[1.5e2[1]]", value: []interface{}{150.0, []interface{}{1}}},
		{src: "[abc// a comment\n]", value: []interface{}{"abc"}},

		{src: "[ // a comment\n  true\n]", value: []interface{}{true}},
		{src: "[\n  null // a comment\n  true\n]", value: []interface{}{nil, true}},
		{src: "[\n  null / a comment\n  true\n]", expect: "unexpected character ' ' at 2:9"},
	} {
		if testing.Verbose() {
			fmt.Printf("... %d: %q\n", i, d.src)
		}
		var err error
		var v interface{}
		var p sen.Parser
		v, err = p.Parse([]byte(d.src))

		if 0 < len(d.expect) {
			tt.NotNil(t, err, d.src)
			tt.Equal(t, d.expect, err.Error(), i, ": ", d.src)
		} else {
			tt.Nil(t, err, d.src)
			tt.Equal(t, d.value, v, i, ": ", d.src)
		}
	}
}

func TestParserParseReader(t *testing.T) {
	for i, d := range []rdata{
		{src: "null", value: nil},
		// The read buffer is 4096 so force a buffer read in the middle of
		// reading a token.
		{src: strings.Repeat(" ", 4094) + "null ", value: nil},
		{src: strings.Repeat(" ", 4094) + "true ", value: true},
		{src: strings.Repeat(" ", 4094) + "false ", value: false},
		{src: strings.Repeat(" ", 4094) + "hello\n  ", value: "hello"},
		{src: strings.Repeat(" ", 4090) + "{abc:def} ", value: map[string]interface{}{"abc": "def"}},
		{src: strings.Repeat(" ", 4093) + "{abc:def} ", value: map[string]interface{}{"abc": "def"}},
		{src: strings.Repeat(" ", 4093) + "[abc[def]]", value: []interface{}{"abc", []interface{}{"def"}}},
		{src: strings.Repeat(" ", 4093) + "[abc]", value: []interface{}{"abc"}},
		{src: strings.Repeat(" ", 4093) + "[abc// comment\n]", value: []interface{}{"abc"}},
		{src: strings.Repeat(" ", 4093) + "[abc{x:1}]", value: []interface{}{"abc", map[string]interface{}{"x": 1}}},

		{src: strings.Repeat(" ", 4094) + "abc$", expect: "expected a value, not '$' at 1:2"},
		{src: strings.Repeat(" ", 4094) + "hello\n $", expect: "extra characters after close, '$' at 2:2"},
		{src: strings.Repeat(" ", 4094) + "hello]", expect: "too many closes at 1:4"},
		{src: strings.Repeat(" ", 4094) + "hello}", expect: "too many closes at 1:4"},
	} {
		if testing.Verbose() {
			fmt.Printf("... %d: %q\n", i, d.src)
		}
		var err error
		var v interface{}
		var p sen.Parser
		v, err = p.ParseReader(strings.NewReader(d.src))

		if 0 < len(d.expect) {
			tt.NotNil(t, err, d.src)
			tt.Equal(t, d.expect, err.Error(), i, ": ", d.src)
		} else {
			tt.Nil(t, err, d.src)
			tt.Equal(t, d.value, v, i, ": ", d.src)
		}
	}
}

func TestParserParseCallback(t *testing.T) {
	var results []byte
	cb := func(n interface{}) bool {
		if 0 < len(results) {
			results = append(results, ' ')
		}
		results = append(results, fmt.Sprintf("%v", n)...)
		return false
	}
	var p sen.Parser
	v, err := p.Parse([]byte(callbackSEN), cb)
	tt.Nil(t, err)
	tt.Nil(t, v)
	tt.Equal(t, `1 [2] map[x:3] true false 123`, string(results))

	_, _ = p.Parse([]byte("[1,[2,[3}]]")) // fail to leave stack not cleaned up

	results = results[:0]
	v, err = p.Parse([]byte(callbackSEN), cb)
	tt.Nil(t, err)
	tt.Nil(t, v)
	tt.Equal(t, `1 [2] map[x:3] true false 123`, string(results))
}

func TestParserParseReaderCallback(t *testing.T) {
	var results []byte
	cb := func(n interface{}) bool {
		if 0 < len(results) {
			results = append(results, ' ')
		}
		results = append(results, fmt.Sprintf("%v", n)...)
		return false
	}
	var p sen.Parser
	v, err := p.ParseReader(strings.NewReader("\xef\xbb\xbf"+callbackSEN), cb)
	tt.Nil(t, err)
	tt.Nil(t, v)
	tt.Equal(t, `1 [2] map[x:3] true false 123`, string(results))

	results = results[:0]
	v, err = p.ParseReader(strings.NewReader(callbackSEN), cb)
	tt.Nil(t, err)
	tt.Nil(t, v)
	tt.Equal(t, `1 [2] map[x:3] true false 123`, string(results))
}

func TestParseBadArg(t *testing.T) {
	var p sen.Parser
	_, err := p.Parse([]byte(callbackSEN), "bad")
	tt.NotNil(t, err)

	_, err = p.ParseReader(strings.NewReader(callbackSEN), "bad")
	tt.NotNil(t, err)
}

func TestParserParseReaderErrRead(t *testing.T) {
	var p sen.Parser
	r := tt.ShortReader{Max: 20, Content: []byte(callbackSEN)}
	_, err := p.ParseReader(&r)
	tt.NotNil(t, err)
}

func TestParserParseReaderEOF(t *testing.T) {
	var p sen.Parser
	_, err := p.ParseReader(iotest.DataErrReader(strings.NewReader("[1 2]")))
	tt.Nil(t, err)
}

func TestParserParseReaderErr(t *testing.T) {
	var p sen.Parser
	_, err := p.ParseReader(iotest.DataErrReader(strings.NewReader("[1 2}")))
	tt.NotNil(t, err)

	r := tt.ShortReader{Max: 5000, Content: []byte("[ 123" + strings.Repeat(",  123", 120) + "]")}
	_, err = p.ParseReader(&r)
	tt.NotNil(t, err)
}
