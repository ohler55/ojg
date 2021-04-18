// Copyright (c) 2021, Peter Ohler, All rights reserved.

package sen_test

import (
	"fmt"
	"io"
	"strings"
	"testing"

	"github.com/ohler55/ojg/gen"
	"github.com/ohler55/ojg/oj"
	"github.com/ohler55/ojg/sen"
	"github.com/ohler55/ojg/tt"
)

type testHandler struct {
	buf []byte
}

func (h *testHandler) Null() {
	h.buf = append(h.buf, "null "...)
}

func (h *testHandler) Bool(v bool) {
	h.buf = append(h.buf, fmt.Sprintf("%t ", v)...)
}

func (h *testHandler) Int(v int64) {
	h.buf = append(h.buf, fmt.Sprintf("%d ", v)...)
}

func (h *testHandler) Float(v float64) {
	h.buf = append(h.buf, fmt.Sprintf("%g ", v)...)
}

func (h *testHandler) Number(v gen.Big) {
	h.buf = append(h.buf, fmt.Sprintf("%s ", v)...)
}

func (h *testHandler) String(v string) {
	h.buf = append(h.buf, fmt.Sprintf("%s ", v)...)
}

func (h *testHandler) ObjectStart() {
	h.buf = append(h.buf, '{')
	h.buf = append(h.buf, ' ')
}

func (h *testHandler) ObjectEnd() {
	h.buf = append(h.buf, '}')
	h.buf = append(h.buf, ' ')
}

func (h *testHandler) ArrayStart() {
	h.buf = append(h.buf, '[')
	h.buf = append(h.buf, ' ')
}

func (h *testHandler) ArrayEnd() {
	h.buf = append(h.buf, ']')
	h.buf = append(h.buf, ' ')
}

type tokeTest struct {
	src    string
	expect string
	err    string
}

func TestTokenizerParseBasic(t *testing.T) {
	toker := sen.Tokenizer{}
	h := testHandler{}
	src := `[true,null,123,12.3]{x:12345678901234567890}`
	err := toker.Parse([]byte(src), &h)
	tt.Nil(t, err)
	tt.Equal(t, "[ true null 123 12.3 ] { x 12345678901234567890 } ", string(h.buf))

	h.buf = h.buf[:0]
	err = sen.Tokenize([]byte(src), &h)
	tt.Nil(t, err)
	tt.Equal(t, "[ true null 123 12.3 ] { x 12345678901234567890 } ", string(h.buf))

	h.buf = h.buf[:0]
	toker.OnlyOne = true
	err = toker.Parse([]byte("[1, 2, 3]  "), &h)
	tt.Nil(t, err)
	tt.Equal(t, "[ 1 2 3 ] ", string(h.buf))

	err = toker.Parse([]byte("[1, 2, 3]  4"), &h)
	tt.NotNil(t, err)
}

func TestTokenizerLoad(t *testing.T) {
	toker := sen.Tokenizer{}
	h := testHandler{}
	err := toker.Load(strings.NewReader("\xef\xbb\xbf"+`[true,null,123,12.3]{x:3}`), &h)
	tt.Nil(t, err)
	tt.Equal(t, "[ true null 123 12.3 ] { x 3 } ", string(h.buf))
}

func TestTokenizerLoadErrRead(t *testing.T) {
	h := oj.ZeroHandler{}
	r := tt.ShortReader{Max: 5, Content: []byte("[1, 2, 3, true, false]")}
	err := sen.TokenizeLoad(&r, &h)
	tt.NotNil(t, err)

	r = tt.ShortReader{Max: 5000, Content: []byte("[ 123" + strings.Repeat(",  123", 120) + "]")}
	err = sen.TokenizeLoad(&r, &h)
	tt.NotNil(t, err)
}

type eofReader int

func (r eofReader) Read(b []byte) (int, error) {
	b[0] = '['
	return 1, io.EOF
}

func TestTokenizerLoadEOF(t *testing.T) {
	h := oj.ZeroHandler{}
	toker := sen.Tokenizer{}
	err := toker.Load(eofReader(0), &h)
	tt.NotNil(t, err)

	err = toker.Load(eofReader(0), &h)
	tt.NotNil(t, err)
}

func TestTokenizerLoadMany(t *testing.T) {
	h := oj.ZeroHandler{}
	for i, s := range []string{
		// The read buffer is 4096 so force a buffer read in the middle of
		// reading a token.
		strings.Repeat(" ", 4094) + "null ",
		strings.Repeat(" ", 4094) + "true ",
		strings.Repeat(" ", 4094) + "false ",
		strings.Repeat(" ", 4094) + `{x:1}`,
		strings.Repeat(" ", 4095) + `"x"`,
		strings.Repeat(" ", 4095) + "xyz",
		strings.Repeat(" ", 4094) + "[xyz]",
		strings.Repeat(" ", 4094) + "[xyz[]]",
		strings.Repeat(" ", 4094) + "[xyz{}]",
		strings.Repeat(" ", 4092) + "{x:abc}",
		strings.Repeat(" ", 4094) + "abc// comment\n",
		strings.Repeat(" ", 4094) + "[abc\n  def]",
	} {
		toker := sen.Tokenizer{}
		err := toker.Load(strings.NewReader(s), &h)
		tt.Nil(t, err, i)
	}
}

func TestTokenizerMany(t *testing.T) {
	for i, d := range []tokeTest{
		{src: "null", expect: "null"},
		{src: "true", expect: "true"},
		{src: "false", expect: "false"},
		{src: "false \n ", expect: "false"},
		{src: "hello", expect: "hello"},
		{src: "hello ", expect: "hello"},
		{src: `"hello"`, expect: "hello"},
		{src: "[one two]", expect: "[ one two ]"},
		{src: "123", expect: "123"},
		{src: "-12.3", expect: "-12.3"},
		{src: "2e-7", expect: "2e-07"},
		{src: "-12.5e-2", expect: "-0.125"},
		{src: "0", expect: "0"},
		{src: "0\n ", expect: "0"},
		{src: "-12.3 ", expect: "-12.3"},
		{src: "-12.3\n", expect: "-12.3"},
		{src: "-12.3e-5", expect: "-0.000123"},
		{src: "12.3e+5 ", expect: "1.23e+06"},
		{src: "12.3e+5\n ", expect: "1.23e+06"},
		{src: "12.3e+05\n ", expect: "1.23e+06"},
		{src: "12.3e-05\n ", expect: "0.000123"},
		{src: `12345678901234567890`, expect: "12345678901234567890"},
		{src: `9223372036854775807`, expect: "9223372036854775807"},
		{src: `9223372036854775808`, expect: "9223372036854775808"},
		{src: `-9223372036854775807`, expect: "-9223372036854775807"},
		{src: `-9223372036854775808`, expect: "-9223372036854775808"},
		{src: `0.9223372036854775808`, expect: "0.9223372036854775808"},
		{src: `-0.9223372036854775808`, expect: "-0.9223372036854775808"},
		{src: `1.2e1025`, expect: "1.2e1025"},
		{src: `-1.2e-1025`, expect: "-1.2e-1025"},
		{src: `12345678901234567890.321e66`, expect: "12345678901234567890.321e66"},
		{src: `321.12345678901234567890e66`, expect: "321.12345678901234567890e66"},
		{src: `321.123e2345`, expect: "321.123e2345"},

		{src: "\xef\xbb\xbf\"xyz\"", expect: "xyz"},

		{src: "[]", expect: "[ ]"},
		{src: "[0,\ntrue , false,null]", expect: "[ 0 true false null ]"},
		{src: `[0.1e3,"x",-1,{}]`, expect: "[ 100 x -1 { } ]"},
		{src: "[1.2,0]", expect: "[ 1.2 0 ]"},
		{src: "[1.2e2,0.1]", expect: "[ 120 0.1 ]"},
		{src: "[1.2e2,0]", expect: "[ 120 0 ]"},
		{src: "[true]", expect: "[ true ]"},
		{src: "[true,false]", expect: "[ true false ]"},
		{src: "[[]]", expect: "[ [ ] ]"},
		{src: "[[true]]", expect: "[ [ true ] ]"},
		{src: `"x\t\n\"\b\f\r\u0041\\\/y"`, expect: "x\t\n\"\b\f\r\u0041\\/y"},
		{src: `"x\u004a\u004Ay"`, expect: "xJJy"},
		{src: `"x\ry"`, expect: "x\ry"},

		{src: "{}", expect: "{ }"},
		{src: `{"a\tbc":true}`, expect: "{ a\tbc true }"},
		{src: `{x:null}`, expect: "{ x null }"},
		{src: `{x:true}`, expect: "{ x true }"},
		{src: `{x:false}`, expect: "{ x false }"},
		{src: "{\"z\":0,\n\"z2\":0}", expect: "{ z 0 z2 0 }"},
		{src: `{"z":1.2,"z2":0}`, expect: "{ z 1.2 z2 0 }"},
		{src: `{"abc":{"def" :3}}`, expect: "{ abc { def 3 } }"},
		{src: `{"x":1.2e3,"y":true}`, expect: "{ x 1200 y true }"},
		{src: `{"abc": [{"x": {"y": [{"b": true}]},"z": 7}]}`, expect: "{ abc [ { x { y [ { b true } ] } z 7 } ] }"},

		{src: "{}}", err: "unexpected object close at 1:3"},
		{src: "{ \n", err: "not closed at 2:1"},
		{src: "{}\n }", err: "unexpected object close at 2:2"},
		{src: "{]}", err: "unexpected array close at 1:2"},
		{src: "[}]", err: "unexpected object close at 1:2"},
		{src: "{\"a\" \n : 1]}", err: "unexpected array close at 2:5"},
		{src: `[1}]`, err: "unexpected object close at 1:3"},
		{src: `1]`, err: "unexpected array close at 1:2"},
		{src: `1}`, err: "unexpected object close at 1:2"},
		{src: `]`, err: "unexpected array close at 1:1"},
		{src: `[null x`, err: "not closed at 1:8"},
		{src: "{\n\"x\":1 ]", err: "unexpected array close at 2:7"},
		{src: `[1 }`, err: "unexpected object close at 1:4"},
		{src: `{"x"x}`, err: "expected a colon, not 'x' at 1:5"},
		{src: `-x`, err: "invalid number at 1:2"},
		{src: `0]`, err: "unexpected array close at 1:2"},
		{src: `0\n $`, err: "invalid number at 1:2"},
		{src: `0}`, err: "unexpected object close at 1:2"},
		{src: `0x`, err: "invalid number at 1:2"},
		{src: `1x`, err: "invalid number at 1:2"},
		{src: `1.x`, err: "invalid number at 1:3"},
		{src: `1.2x`, err: "invalid number at 1:4"},
		{src: `1.2ex`, err: "invalid number at 1:5"},
		{src: `1.2e+x`, err: "invalid number at 1:6"},
		{src: `1.2e`, err: "incomplete JSON at 1:5"},
		{src: "1.2\n]", err: "unexpected array close at 2:1"},
		{src: `1.2]`, err: "unexpected array close at 1:4"},
		{src: `1.2}`, err: "unexpected object close at 1:4"},
		{src: `1.2e2]`, err: "unexpected array close at 1:6"},
		{src: `1.2e2}`, err: "unexpected object close at 1:6"},
		{src: `1.2e2x`, err: "invalid number at 1:6"},
		{src: "\"x\fy\"", err: "invalid JSON character 0x0c at 1:3"},
		{src: `"x\zy"`, err: "invalid JSON escape character '\\z' at 1:4"},
		{src: `"x\u004z"`, err: "invalid JSON unicode character 'z' at 1:8"},
		{src: "\xef\xbb[]", err: "expected BOM at 1:3"},
		{src: "#x", err: "unexpected character '#' at 1:1"},
		{src: "x]", err: "unexpected array close at 1:2"},
		{src: "x}", err: "unexpected object close at 1:2"},
		{src: "x#", err: "unexpected character '#' at 1:2"},
		{src: "{x#:1}", err: "expected a colon, not '#' at 1:3"},
		{src: `{123}`, err: "expected a key at 1:5"},
		{src: `{123 }`, err: "expected a key at 1:5"},
		{src: `{123[}`, err: "expected a key at 1:5"},
		{src: `{{}}`, err: "expected a key at 1:2"},
		{src: `{123`, err: "not closed at 1:5"},
		{src: "{123\n", err: "expected a key at 1:5"},
		{src: `{123// comment`, err: "expected a key at 1:5"},
		{src: `{123{}}`, err: "expected a key at 1:5"},
		{src: `{[]]}`, err: "expected a key at 1:2"},

		{src: "[0 1,2]", expect: "[ 0 1 2 ]"},
		{src: "[0[1[2]]]", expect: "[ 0 [ 1 [ 2 ] ] ]"},

		{src: `{aaa:0 "bbb":"one" , c:2}`, expect: "{ aaa 0 bbb one c 2 }"},
		{src: "{aaa\n:one b:\ntwo}", expect: "{ aaa one b two }"},
		{src: "[abc[x]]", expect: "[ abc [ x ] ]"},
		{src: "[abc{x:1}]", expect: "[ abc { x 1 } ]"},
		{src: `{aaa:"bbb" "bbb":"one" , c:2}`, expect: "{ aaa bbb bbb one c 2 }"},
		{src: `{aaa:"b\tb" x:2}`, expect: "{ aaa b\tb x 2 }"},

		{src: "[0{x:1}]", expect: "[ 0 { x 1 } ]"},
		{src: "[1{x:1}]", expect: "[ 1 { x 1 } ]"},
		{src: "[1.5{x:1}]", expect: "[ 1.5 { x 1 } ]"},
		{src: "[1.5[1]]", expect: "[ 1.5 [ 1 ] ]"},
		{src: "[1.5e2{x:1}]", expect: "[ 150 { x 1 } ]"},
		{src: "[1.5e2[1]]", expect: "[ 150 [ 1 ] ]"},

		{src: "[abc// a comment\n]", expect: "[ abc ]"},
		{src: "[123// a comment\n]", expect: "[ 123 ]"},
		{src: "[ // a comment\n  true\n]", expect: "[ true ]"},
		{src: "[\n  null // a comment\n  true\n]", expect: "[ null true ]"},
		{src: "[\n  null / a comment\n  true\n]", err: "unexpected character ' ' at 2:9"},
	} {
		if testing.Verbose() {
			fmt.Printf("... %d: %q\n", i, d.src)
		}
		h := testHandler{}
		err := sen.TokenizeString(d.src, &h)
		if 0 < len(d.expect) {
			tt.Nil(t, err, d.src)
			tt.Equal(t, d.expect, strings.TrimSpace(string(h.buf)), i, ": ", d.src)
		} else {
			tt.NotNil(t, err, d.src)
			tt.Equal(t, d.err, err.Error(), i, ": ", d.src)
		}
	}
}
