// Copyright (c) 2020, Peter Ohler, All rights reserved.

package sen_test

import (
	"encoding/json"
	"fmt"
	"strings"
	"testing"
	"testing/iotest"

	"github.com/ohler55/ojg"
	"github.com/ohler55/ojg/sen"
	"github.com/ohler55/ojg/tt"
)

const callbackSEN = `
1
[2]
{x:3}
true false 123`

const tokenSEN = `
abc
def`

type rdata struct {
	src    string
	expect string
	value  any
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
		{src: `'ab"cd'`, value: `ab"cd`},
		{src: `"ab'cd"`, value: `ab'cd`},
		{src: `'ab\'cd'`, value: `ab'cd`},
		{src: "[one two]", value: []any{"one", "two"}},
		{src: "123", value: 123},
		{src: "-12.3", value: -12.3},
		{src: "2e-7", value: 2e-7},
		{src: "-12.5e-2", value: -0.125},
		{src: "0", value: 0},
		{src: "0\n ", value: 0},
		{src: "-12.3 ", value: -12.3},
		{src: "-12.3\n", value: -12.3},
		{src: "-12.3e-5", value: -12.3e-5},
		{src: "12.3e+5 ", value: 12.3e+5},
		{src: "12.3e+5\n ", value: 12.3e+5},
		{src: "12.3e+05\n ", value: 12.3e+5},
		{src: "12.3e-05\n ", value: 12.3e-5},
		{src: `12345678901234567890`, value: "12345678901234567890"},
		{src: `9223372036854775807`, value: "9223372036854775807"},   // max int
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

		{src: "[]", value: []any{}},
		{src: "[0,\ntrue , false,null]", value: []any{0, true, false, nil}},
		{src: `[0.1e3,"x",-1,{}]`, value: []any{100.0, "x", -1, map[string]any{}}},
		{src: "[1.2,0]", value: []any{1.2, 0}},
		{src: "[1.2e2,0.1]", value: []any{1.2e2, 0.1}},
		{src: "[1.2e2,0]", value: []any{1.2e2, 0}},
		{src: "[true]", value: []any{true}},
		{src: "[true,false]", value: []any{true, false}},
		{src: "[[]]", value: []any{[]any{}}},
		{src: "[[true]]", value: []any{[]any{true}}},
		{src: `"x\t\n\"\b\f\r\u0041\\\/y"`, value: "x\t\n\"\b\f\r\u0041\\/y"},
		{src: `"x\u004a\u004Ay"`, value: "xJJy"},
		{src: `"x\ry"`, value: "x\ry"},

		{src: "{}", value: map[string]any{}},
		{src: `{"a\tbc":true}`, value: map[string]any{"a\tbc": true}},
		{src: `{x:null}`, value: map[string]any{"x": nil}},
		{src: `{x:true}`, value: map[string]any{"x": true}},
		{src: `{x:false}`, value: map[string]any{"x": false}},
		{src: "{\"z\":0,\n\"z2\":0}", value: map[string]any{"z": 0, "z2": 0}},
		{src: `{"z":1.2,"z2":0}`, value: map[string]any{"z": 1.2, "z2": 0}},
		{src: `{"abc":{"def" :3}}`, value: map[string]any{"abc": map[string]any{"def": 3}}},
		{src: `{"x":1.2e3,"y":true}`, value: map[string]any{"x": 1200.0, "y": true}},

		{src: `{"abc": [{"x": {"y": [{"b": true}]},"z": 7}]}`,
			value: map[string]any{
				"abc": []any{
					map[string]any{
						"x": map[string]any{
							"y": []any{
								map[string]any{
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
		{src: "{]}", expect: "unexpected array close at 1:2"},
		{src: "[}]", expect: "unexpected object close at 1:2"},
		{src: "{\"a\" \n : 1]}", expect: "unexpected array close at 2:5"},
		{src: `[1}]`, expect: "unexpected object close at 1:3"},
		{src: `1]`, expect: "unexpected array close at 1:2"},
		{src: `1}`, expect: "unexpected object close at 1:2"},
		{src: `]`, expect: "unexpected array close at 1:1"},
		{src: `[null x`, expect: "not closed at 1:8"},
		{src: "{\n\"x\":1 ]", expect: "unexpected array close at 2:7"},
		{src: `[1 }`, expect: "unexpected object close at 1:4"},
		{src: `{"x"x}`, expect: "expected a colon, not 'x' at 1:5"},
		{src: `-x`, expect: "invalid number at 1:2"},
		{src: `0]`, expect: "unexpected array close at 1:2"},
		{src: `0\n $`, expect: "invalid number at 1:2"},
		{src: `0}`, expect: "unexpected object close at 1:2"},
		{src: `0x`, expect: "invalid number at 1:2"},
		{src: `1x`, expect: "invalid number at 1:2"},
		{src: `1.x`, expect: "invalid number at 1:3"},
		{src: `1.2x`, expect: "invalid number at 1:4"},
		{src: `1.2ex`, expect: "invalid number at 1:5"},
		{src: `1.2e+x`, expect: "invalid number at 1:6"},
		{src: `1.2e`, expect: "incomplete JSON at 1:5"},
		{src: "1.2\n]", expect: "extra characters after close, ']' at 2:1"},
		{src: `1.2]`, expect: "unexpected array close at 1:4"},
		{src: `1.2}`, expect: "unexpected object close at 1:4"},
		{src: `1.2e2]`, expect: "unexpected array close at 1:6"},
		{src: `1.2e2}`, expect: "unexpected object close at 1:6"},
		{src: `1.2e2x`, expect: "invalid number at 1:6"},
		{src: "\"x\fy\"", expect: "invalid JSON character 0x0c at 1:3"},
		{src: `"x\zy"`, expect: "invalid JSON escape character '\\z' at 1:4"},
		{src: `"x\u004z"`, expect: "invalid JSON unicode character 'z' at 1:8"},
		{src: "\xef\xbb[]", expect: "expected BOM at 1:3"},
		{src: "#x", expect: "unexpected character '#' at 1:1"},
		{src: "x]", expect: "extra characters after close, ']' at 1:2"},
		{src: "x}", expect: "extra characters after close, '}' at 1:2"},
		{src: "x#", expect: "extra characters after close, '#' at 1:2"},
		{src: "{x#:1}", expect: "expected a colon, not '#' at 1:3"},
		{src: `{123}`, expect: "expected a key at 1:5"},
		{src: `{123 }`, expect: "expected a key at 1:5"},
		{src: `{123[}`, expect: "expected a key at 1:5"},
		{src: `{{}}`, expect: "expected a key at 1:3"},
		{src: `{123`, expect: "not closed at 1:5"},
		{src: "{123\n", expect: "expected a key at 1:5"},
		{src: `{123// comment`, expect: "expected a key at 1:5"},
		{src: `{123{}}`, expect: "expected a key at 1:5"},
		{src: `{[]]}`, expect: "expected a key at 1:3"},

		{src: "[0 1,2]", value: []any{0, 1, 2}},
		{src: "[0[1[2]]]", value: []any{0, []any{1, []any{2}}}},
		{src: `{aaa:0 "bbb":"one" , c:2}`, value: map[string]any{"aaa": 0, "bbb": "one", "c": 2}},
		{src: "{aaa\n:one b:\ntwo}", value: map[string]any{"aaa": "one", "b": "two"}},
		{src: "[abc[x]]", value: []any{"abc", []any{"x"}}},
		{src: "[abc{x:1}]", value: []any{"abc", map[string]any{"x": 1}}},
		{src: `{aaa:"bbb" "bbb":"one" , c:2}`, value: map[string]any{"aaa": "bbb", "bbb": "one", "c": 2}},
		{src: `{aaa:"b\tb" x:2}`, value: map[string]any{"aaa": "b\tb", "x": 2}},

		{src: "[0{x:1}]", value: []any{0, map[string]any{"x": 1}}},
		{src: "[1{x:1}]", value: []any{1, map[string]any{"x": 1}}},
		{src: "[1.5{x:1}]", value: []any{1.5, map[string]any{"x": 1}}},
		{src: "[1.5[1]]", value: []any{1.5, []any{1}}},
		{src: "[1.5e2{x:1}]", value: []any{150., map[string]any{"x": 1}}},
		{src: "[1.5e2[1]]", value: []any{150.0, []any{1}}},

		{src: "[abc// a comment\n]", value: []any{"abc"}},
		{src: "[123// a comment\n]", value: []any{123}},
		{src: "[ // a comment\n  true\n]", value: []any{true}},
		{src: "[\n  null // a comment\n  true\n]", value: []any{nil, true}},
		{src: "[\n  null / a comment\n  true\n]", expect: "unexpected character ' ' at 2:9"},
	} {
		if testing.Verbose() {
			fmt.Printf("... %d: %q\n", i, d.src)
		}
		v, err := sen.Parse([]byte(d.src))
		if 0 < len(d.expect) {
			tt.NotNil(t, err, d.src)
			tt.Equal(t, d.expect, err.Error(), i, ": ", d.src)
		} else {
			tt.Nil(t, err, d.src)
			tt.Equal(t, d.value, v, i, ": '", d.src, "'")
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
		{src: strings.Repeat(" ", 4092) + "{x:null} ", value: map[string]any{"x": nil}},
		{src: strings.Repeat(" ", 4092) + "{x:true} ", value: map[string]any{"x": true}},
		{src: strings.Repeat(" ", 4092) + "{x:false} ", value: map[string]any{"x": false}},
		{src: strings.Repeat(" ", 4090) + "{abc:def} ", value: map[string]any{"abc": "def"}},
		{src: strings.Repeat(" ", 4093) + "{abc:def} ", value: map[string]any{"abc": "def"}},
		{src: strings.Repeat(" ", 4093) + "[abc[def]]", value: []any{"abc", []any{"def"}}},
		{src: strings.Repeat(" ", 4093) + "[abc]", value: []any{"abc"}},
		{src: strings.Repeat(" ", 4093) + "[abc// comment\n]", value: []any{"abc"}},
		{src: strings.Repeat(" ", 4093) + "[abc{x:1}]", value: []any{"abc", map[string]any{"x": 1}}},

		{src: strings.Repeat(" ", 4094) + "abc#", expect: "unexpected character '#' at 1:2"},
		{src: strings.Repeat(" ", 4094) + "hello\n #", expect: "extra characters after close, '#' at 2:2"},
		{src: strings.Repeat(" ", 4094) + "hello]", expect: "unexpected array close at 1:4"},
		{src: strings.Repeat(" ", 4094) + "hello}", expect: "unexpected object close at 1:4"},
		{src: strings.Repeat(" ", 4095) + `"x"`, value: "x"},
	} {
		if testing.Verbose() {
			fmt.Printf("... %d: %q\n", i, d.src)
		}
		var err error
		var v any
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
	cb := func(n any) {
		if 0 < len(results) {
			results = append(results, ' ')
		}
		results = append(results, fmt.Sprintf("%v", n)...)
	}
	p := sen.Parser{Reuse: true}
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

	results = results[:0]
	v, err = p.Parse([]byte("123"), cb)
	tt.Nil(t, err)
	tt.Nil(t, v)
	tt.Equal(t, `123`, string(results))

	results = results[:0]
	v, err = p.Parse([]byte("abc"), cb)
	tt.Nil(t, err)
	tt.Nil(t, v)
	tt.Equal(t, "abc", string(results))
}

func TestParserParseCallbackAlt(t *testing.T) {
	var results []byte
	cb := func(n any) bool {
		if 0 < len(results) {
			results = append(results, ' ')
		}
		results = append(results, fmt.Sprintf("%v", n)...)
		return false
	}
	p := sen.Parser{Reuse: true}
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

	results = results[:0]
	v, err = p.Parse([]byte("123"), cb)
	tt.Nil(t, err)
	tt.Nil(t, v)
	tt.Equal(t, `123`, string(results))

	results = results[:0]
	v, err = p.Parse([]byte("abc"), cb)
	tt.Nil(t, err)
	tt.Nil(t, v)
	tt.Equal(t, "abc", string(results))
}

func TestParserParseReaderCallback(t *testing.T) {
	var results []byte
	cb := func(n any) {
		if 0 < len(results) {
			results = append(results, ' ')
		}
		results = append(results, fmt.Sprintf("%v", n)...)
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

func TestParserParseReaderCallbackAlt(t *testing.T) {
	var results []byte
	cb := func(n any) bool {
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

func TestParserParseChan(t *testing.T) {
	var results []byte
	rc := make(chan any, 10)
	var p sen.Parser
	_, err := p.Parse([]byte(callbackSEN), rc)
	tt.Nil(t, err)
	rc <- nil
	for {
		n := <-rc
		if n == nil {
			break
		}
		if 0 < len(results) {
			results = append(results, ' ')
		}
		results = append(results, fmt.Sprintf("%v", n)...)
	}
	tt.Equal(t, `1 [2] map[x:3] true false 123`, string(results))

	results = results[:0]
	_, err = p.Parse([]byte(tokenSEN), rc)
	tt.Nil(t, err)
	rc <- nil
	for {
		n := <-rc
		if n == nil {
			break
		}
		if 0 < len(results) {
			results = append(results, ' ')
		}
		results = append(results, fmt.Sprintf("%v", n)...)
	}
	tt.Equal(t, `abc def`, string(results))
}

func TestParserParseReaderChan(t *testing.T) {
	var results []byte
	rc := make(chan any, 10)
	_, err := sen.ParseReader(strings.NewReader(callbackSEN), rc)
	tt.Nil(t, err)
	rc <- nil
	for {
		n := <-rc
		if n == nil {
			break
		}
		if 0 < len(results) {
			results = append(results, ' ')
		}
		results = append(results, fmt.Sprintf("%v", n)...)
	}
	tt.Equal(t, `1 [2] map[x:3] true false 123`, string(results))
}

func TestMustParsePanic(t *testing.T) {
	tt.Panic(t, func() { _ = sen.MustParse([]byte("[1 2}")) })
}

func TestMustParseReaderPanic(t *testing.T) {
	tt.Panic(t, func() { _ = sen.MustParseReader(strings.NewReader("[1 2}")) })
}

func TestParserMustParsePanic(t *testing.T) {
	var p sen.Parser
	tt.Panic(t, func() { _ = p.MustParse([]byte("[1 2}")) })
}

func TestParserMustParseReader(t *testing.T) {
	var p sen.Parser
	tt.Panic(t, func() { _ = p.MustParseReader(iotest.DataErrReader(strings.NewReader("[1 2}"))) })
}

func TestParserPlus(t *testing.T) {
	src := `['abc' + "def" + 'ghi']`
	v := sen.MustParse([]byte(src))
	tt.Equal(t, []any{"abcdefghi"}, v)

	src = `{a: abc + "def" + 'ghi'}`
	v = sen.MustParse([]byte(src))
	tt.Equal(t, map[string]any{"a": "abcdefghi"}, v)
}

func TestParserTokenFunc(t *testing.T) {
	v := sen.MustParse([]byte("fun(123)"))
	tt.Equal(t, 123, v)

	v = sen.MustParse([]byte(`[fun("xyz")]`))
	tt.Equal(t, []any{"xyz"}, v)

	p := sen.Parser{}
	p.AddTokenFunc("fun", func(args ...any) any {
		var sum int64
		for _, a := range args {
			i, _ := a.(int64)
			sum += i
		}
		return sum
	})
	v = p.MustParse([]byte("fun(1,2,3)"))
	tt.Equal(t, 6, v)

	_, err := sen.Parse([]byte("[1,2,3)"))
	tt.NotNil(t, err)

	_, err = sen.Parse([]byte("3)"))
	tt.NotNil(t, err)

	src := strings.Repeat(" ", 4094) + "fun(1,3,5)"
	v = p.MustParseReader(strings.NewReader(src))
	tt.Equal(t, 9, v)

	src = strings.Repeat(" ", 4090) + "fun(abc)"
	v = sen.MustParseReader(strings.NewReader(src))
	tt.Equal(t, "abc", v)
}

func TestParserNumConv(t *testing.T) {
	v := sen.MustParse([]byte("0.1234567890123456789"))
	tt.Equal(t, json.Number("0.1234567890123456789"), v)

	v = sen.MustParse([]byte("0.1234567890123456789"), ojg.NumConvFloat64)
	tt.Equal(t, 0.123456789012345678, v)

	v = sen.MustParse([]byte("0.1234567890123456789"), ojg.NumConvString)
	tt.Equal(t, "0.1234567890123456789", v)

	v = sen.MustParseReader(strings.NewReader("0.1234567890123456789"), ojg.NumConvFloat64)
	tt.Equal(t, 0.123456789012345678, v)
}
