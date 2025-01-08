// Copyright (c) 2021, Peter Ohler, All rights reserved.

package oj_test

import (
	"fmt"
	"io"
	"strings"
	"testing"

	"github.com/ohler55/ojg/oj"
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

func (h *testHandler) Number(v string) {
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

func (h *testHandler) Key(v string) {
	h.buf = append(h.buf, fmt.Sprintf("%s: ", v)...)
}

func (h *testHandler) ArrayStart() {
	h.buf = append(h.buf, '[')
	h.buf = append(h.buf, ' ')
}

func (h *testHandler) ArrayEnd() {
	h.buf = append(h.buf, ']')
	h.buf = append(h.buf, ' ')
}

func TestTokenizerParseBasic(t *testing.T) {
	toker := oj.Tokenizer{}
	h := testHandler{}
	src := `[true,null,123,12.3]{"x":12345678901234567890}`
	err := toker.Parse([]byte(src), &h)
	tt.Nil(t, err)
	tt.Equal(t, "[ true null 123 12.3 ] { x: 12345678901234567890 } ", string(h.buf))

	h.buf = h.buf[:0]
	err = toker.Parse([]byte(src), &h)
	tt.Nil(t, err)
	tt.Equal(t, "[ true null 123 12.3 ] { x: 12345678901234567890 } ", string(h.buf))

	h.buf = h.buf[:0]
	toker.OnlyOne = true
	err = toker.Parse([]byte("[1, 2, 3]  "), &h)
	tt.Nil(t, err)
	tt.Equal(t, "[ 1 2 3 ] ", string(h.buf))

}

func TestTokenizerLoad(t *testing.T) {
	toker := oj.Tokenizer{}
	h := testHandler{}
	err := toker.Load(strings.NewReader("\xef\xbb\xbf"+`[true,null,123,12.3]{"x":3}`), &h)
	tt.Nil(t, err)
	tt.Equal(t, "[ true null 123 12.3 ] { x: 3 } ", string(h.buf))
}

func TestZeroHandler(t *testing.T) {
	h := oj.ZeroHandler{}
	src := `[true,null,123,12.3]{"x":12345678901234567890}`
	err := oj.TokenizeString(src, &h)
	tt.Nil(t, err)

	err = oj.Tokenize([]byte(src), &h)
	tt.Nil(t, err)

	err = oj.TokenizeLoad(strings.NewReader(src), &h)
	tt.Nil(t, err)
}

func TestTokenizerLoadErrRead(t *testing.T) {
	h := oj.ZeroHandler{}
	r := tt.ShortReader{Max: 5, Content: []byte("[1, 2, 3, true, false]")}
	err := oj.TokenizeLoad(&r, &h)
	tt.NotNil(t, err)

	r = tt.ShortReader{Max: 5000, Content: []byte("[ 123" + strings.Repeat(",  123", 120) + "]")}
	err = oj.TokenizeLoad(&r, &h)
	tt.NotNil(t, err)
}

type eofReader int

func (r eofReader) Read(b []byte) (int, error) {
	b[0] = 'X'
	return 1, io.EOF
}

func TestTokenizerLoadEOF(t *testing.T) {
	h := oj.ZeroHandler{}
	toker := oj.Tokenizer{}
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
		strings.Repeat(" ", 4094) + `{"x":1}`,
		strings.Repeat(" ", 4095) + `"x"`,
	} {
		toker := oj.Tokenizer{}
		err := toker.Load(strings.NewReader(s), &h)
		tt.Nil(t, err, i)
	}
}

type tokeTest struct {
	src    string
	expect string
	err    string
}

func TestTokenizerMany(t *testing.T) {
	for i, d := range []tokeTest{
		{src: "null", expect: "null"},
		{src: "true", expect: "true"},
		{src: "false", expect: "false"},
		{src: "false \n ", expect: "false"},
		{src: "123", expect: "123"},
		{src: "-321", expect: "-321"},
		{src: "12.3", expect: "12.3"},
		{src: "0 ", expect: "0"},
		{src: "0\n", expect: "0"},
		{src: "2e-7 ", expect: "2e-07"},
		{src: "-12.3 ", expect: "-12.3"},
		{src: "-12.3\n", expect: "-12.3"},
		{src: "-12.3e-5", expect: "-0.000123"},
		{src: "12.3e+5", expect: "1.23e+06"},
		{src: "12.3e5", expect: "1.23e+06"},
		{src: "12.3e05", expect: "1.23e+06"},
		{src: "12.3e-05", expect: "0.000123"},
		{src: "12.3e+5\n", expect: "1.23e+06"},
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
		{src: "8.26e-05", expect: "8.26e-05"},

		{src: "\xef\xbb\xbf\"xyz\"", expect: "xyz"},
		{src: `"Bénédicte"`, expect: "Bénédicte"},

		{src: "[]", expect: "[ ]"},
		{src: "[0,\ntrue , false,null]", expect: "[ 0 true false null ]"},
		{src: `[0.1e3,"x",-1,{}]`, expect: "[ 100 x -1 { } ]"},
		{src: "[1.2,0]", expect: "[ 1.2 0 ]"},
		{src: "[1.2e2,0.1]", expect: "[ 120 0.1 ]"},
		{src: "[1.2e2,0]", expect: "[ 120 0 ]"},
		{src: "[true]", expect: "[ true ]"},
		{src: "[true,false]", expect: "[ true false ]"},
		{src: "[[]]", expect: "[ [ ] ]"},
		{src: "[true,[]]", expect: "[ true [ ] ]"},
		{src: "[[true]]", expect: "[ [ true ] ]"},
		{src: `"x\t\n\"\b\f\r\u0041\\\/y"`, expect: "x\t\n\"\b\f\r\u0041\\/y"},
		{src: `"x\u004a\u004Ay"`, expect: "xJJy"},
		{src: `[1,"a\tb"]`, expect: "[ 1 a\tb ]"},
		{src: `{"a\tb":1}`, expect: "{ a\tb: 1 }"},
		{src: `{"x":1,"a\tb":2}`, expect: "{ x: 1 a\tb: 2 }"},
		{src: "[0\n,3\n,5.0e2\n]", expect: "[ 0 3 500 ]"},

		{src: "{}", expect: "{ }"},
		{src: `{"abc":true}`, expect: "{ abc: true }"},
		{src: "{\"z\":0,\n\"z2\":0}", expect: "{ z: 0 z2: 0 }"},
		{src: `{"z":1.2,"z2":0}`, expect: "{ z: 1.2 z2: 0 }"},
		{src: `{"abc":{"def":3}}`, expect: "{ abc: { def: 3 } }"},
		{src: `{"x":1.2e3,"y":true}`, expect: "{ x: 1200 y: true }"},
		{src: `{"abc": [{"x": {"y": [{"b": true}]},"z": 7}]}`,
			expect: "{ abc: [ { x: { y: [ { b: true } ] } z: 7 } ] }"},

		{src: "{}}", err: "unexpected object close at 1:3"},
		{src: "{}\n }", err: "unexpected object close at 2:2"},
		{src: "{ \n", err: "incomplete JSON at 2:1"},
		{src: "{]}", err: "expected a string start or object close, not ']' at 1:2"},
		{src: "[}]", err: "unexpected object close at 1:2"},
		{src: "{\"a\" \n : 1]}", err: "unexpected array close at 2:5"},
		{src: `[1}]`, err: "unexpected object close at 1:3"},
		{src: `1]`, err: "unexpected array close at 1:2"},
		{src: `1}`, err: "unexpected object close at 1:2"},
		{src: `]`, err: "unexpected array close at 1:1"},
		{src: `x`, err: "unexpected character 'x' at 1:1"},
		{src: `[1,]`, err: "unexpected character ']' at 1:4"},
		{src: `[null x`, err: "expected a comma or close, not 'x' at 1:7"},
		{src: "{\n\"x\":1 ]", err: "unexpected array close at 2:7"},
		{src: `[1 }`, err: "unexpected object close at 1:4"},
		{src: "{\n\"x\":1,}", err: "expected a string start, not '}' at 2:7"},
		{src: `{"x"x}`, err: "expected a colon, not 'x' at 1:5"},
		{src: `nuul`, err: "expected null at 1:3"},
		{src: `fasle`, err: "expected false at 1:3"},
		{src: `ture`, err: "expected true at 1:2"},
		{src: `[0,nuul]`, err: "expected null at 1:6"},
		{src: `[0,fail]`, err: "expected false at 1:6"},
		{src: `[0,truk]`, err: "expected true at 1:7"},
		{src: `-x`, err: "invalid number at 1:2"},
		{src: `0]`, err: "unexpected array close at 1:2"},
		{src: `0}`, err: "unexpected object close at 1:2"},
		{src: `0x`, err: "invalid number at 1:2"},
		{src: `1x`, err: "invalid number at 1:2"},
		{src: `1.x`, err: "invalid number at 1:3"},
		{src: `1.2x`, err: "invalid number at 1:4"},
		{src: `1.2ex`, err: "invalid number at 1:5"},
		{src: `1.2e+x`, err: "invalid number at 1:6"},
		{src: "1.2\n]", err: "unexpected array close at 2:1"},
		{src: `1.2]`, err: "unexpected array close at 1:4"},
		{src: `1.2}`, err: "unexpected object close at 1:4"},
		{src: `1.2e2]`, err: "unexpected array close at 1:6"},
		{src: `1.2e2}`, err: "unexpected object close at 1:6"},
		{src: `1.2e2x`, err: "invalid number at 1:6"},
		{src: "\"x\ty\"", err: "invalid JSON character 0x09 at 1:3"},
		{src: `"x\zy"`, err: "invalid JSON escape character '\\z' at 1:4"},
		{src: `"x\u004z"`, err: "invalid JSON unicode character 'z' at 1:8"},
		{src: "\xef\xbb[]", err: "expected BOM at 1:3"},
		{src: "[ // a comment\n  true\n]", err: "unexpected character '/' at 1:3"},
	} {
		if testing.Verbose() {
			fmt.Printf("... %d: %s\n", i, d.src)
		}
		h := testHandler{}
		err := oj.TokenizeString(d.src, &h)
		if 0 < len(d.expect) {
			tt.Nil(t, err, d.src)
			tt.Equal(t, d.expect, strings.TrimSpace(string(h.buf)), i, ": ", d.src)
		} else {
			tt.NotNil(t, err, d.src)
			tt.Equal(t, d.err, err.Error(), i, ": ", d.src)
		}
	}
}

func TestTokenizerNesting(t *testing.T) {
	var h testHandler
	err := oj.TokenizeString(`[{"a":[1, 2]}]`, &h)
	tt.Nil(t, err)
	tt.Equal(t, "[ { a: [ 1 2 ] } ] ", string(h.buf))
}
