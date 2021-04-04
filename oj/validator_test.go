// Copyright (c) 2020, Peter Ohler, All rights reserved.

package oj_test

import (
	"bytes"
	"strings"
	"testing"
	"testing/iotest"

	"github.com/ohler55/ojg/oj"
	"github.com/ohler55/ojg/tt"
)

func TestValidatorValidateString(t *testing.T) {
	for i, d := range []data{
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
		{src: "2e-7"},
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

		{src: "{}}", expect: "unexpected object close at 1:3"},
		{src: "{ \n", expect: "incomplete JSON at 2:1"},
		{src: "{]}", expect: "expected a string start or object close, not ']' at 1:2"},
		{src: "[}]", expect: "unexpected object close at 1:2"},
		{src: "{\"a\" \n : 1]}", expect: "unexpected array close at 2:5"},
		{src: `[1}]`, expect: "unexpected object close at 1:3"},
		{src: `1]`, expect: "unexpected array close at 1:2"},
		{src: `1}`, expect: "unexpected object close at 1:2"},
		{src: `]`, expect: "unexpected array close at 1:1"},
		{src: `x`, expect: "unexpected character 'x' at 1:1"},
		{src: `[1,]`, expect: "unexpected character ']' at 1:4"},
		{src: `[null x`, expect: "expected a comma or close, not 'x' at 1:7"},
		{src: "{\n\"x\":1 ]", expect: "unexpected array close at 2:7"},
		{src: `[1 }`, expect: "unexpected object close at 1:4"},
		{src: "{\n\"x\":1,}", expect: "expected a string start, not '}' at 2:7"},
		{src: `{"x"x}`, expect: "expected a colon, not 'x' at 1:5"},
		{src: `nuul`, expect: "expected null at 1:3"},
		{src: `nxul`, expect: "expected null at 1:2"},
		{src: `fasle`, expect: "expected false at 1:3"},
		{src: `fxsle`, expect: "expected false at 1:2"},
		{src: `ture`, expect: "expected true at 1:2"},
		{src: `trxe`, expect: "expected true at 1:3"},
		{src: `[0,nuts]`, expect: "expected null at 1:6"},
		{src: `[0,fail]`, expect: "expected false at 1:6"},
		{src: `-x`, expect: "invalid number at 1:2"},
		{src: `0]`, expect: "unexpected array close at 1:2"},
		{src: `0}`, expect: "unexpected object close at 1:2"},
		{src: `0x`, expect: "invalid number at 1:2"},
		{src: `1x`, expect: "invalid number at 1:2"},
		{src: `1.x`, expect: "invalid number at 1:3"},
		{src: `1.2x`, expect: "invalid number at 1:4"},
		{src: `1.2ex`, expect: "invalid number at 1:5"},
		{src: `1.2e+x`, expect: "invalid number at 1:6"},
		{src: "1.2\n]", expect: "unexpected array close at 2:1"},
		{src: `1.2]`, expect: "unexpected array close at 1:4"},
		{src: `1.2}`, expect: "unexpected object close at 1:4"},
		{src: `1.2e2]`, expect: "unexpected array close at 1:6"},
		{src: `1.2e2}`, expect: "unexpected object close at 1:6"},
		{src: `1.2e2x`, expect: "invalid number at 1:6"},
		{src: "\"x\ty\"", expect: "invalid JSON character 0x09 at 1:3"},
		{src: `"x\zy"`, expect: "invalid JSON escape character '\\z' at 1:4"},
		{src: `"x\u004z"`, expect: "invalid JSON unicode character 'z' at 1:8"},
		{src: "\xef\xbb[]", expect: "expected BOM at 1:3"},
		{src: "null \n {}", expect: "extra characters after close, '{' at 2:2", onlyOne: true},
		{src: "[ // a comment\n  true\n]", expect: "unexpected character '/' at 1:3"},
	} {
		var err error
		if d.onlyOne {
			p := oj.Validator{OnlyOne: d.onlyOne}
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

func TestValidatorValidateReaderMany(t *testing.T) {
	for i, d := range []data{
		{src: "null"},
		// The read buffer is 4096 so force a buffer read in the middle of
		// reading a token.
		{src: strings.Repeat(" ", 4094) + "null "},
		{src: strings.Repeat(" ", 4094) + "true "},
		{src: strings.Repeat(" ", 4094) + "false "},
	} {
		r := strings.NewReader(d.src)
		err := oj.ValidateReader(r)
		tt.Nil(t, err, i, ": ", d.src)
	}
}

func TestValidatorValidateReaderBasic(t *testing.T) {
	r := strings.NewReader("[true,[false,[null],123],456]")
	err := oj.ValidateReader(r)
	tt.Nil(t, err)

	var buf []byte
	buf = append(buf, "[\n"...)
	for i := 0; i < 1000; i++ {
		buf = append(buf, "  true,\n"...)
	}
	buf = append(buf, "  false\n]\n"...)
	br := bytes.NewReader(buf)
	err = oj.ValidateReader(br)
	tt.Nil(t, err)
}

func TestValidatorValidateResuse(t *testing.T) {
	var v oj.Validator
	err := v.Validate([]byte("[true,[false,[null],123],456]"))
	tt.Nil(t, err)
	// a second time
	err = v.Validate([]byte("[true,[false,[null],123],456]"))
	tt.Nil(t, err)

	r := strings.NewReader("[true,[false,[null],123],456]")
	err = v.ValidateReader(r)
	tt.Nil(t, err)
}

func TestValidatorValidateBOM(t *testing.T) {
	var v oj.Validator
	err := v.Validate([]byte("\xef\xbb\xbf[true]"))
	tt.Nil(t, err)
}

func TestValidatorValidateReaderBOM(t *testing.T) {
	var v oj.Validator
	err := v.ValidateReader(strings.NewReader("\xef\xbb\xbf[true]"))
	tt.Nil(t, err)
}

func TestValidatorValidateReaderErrRead(t *testing.T) {
	var v oj.Validator
	r := tt.ShortReader{Max: 20, Content: []byte(callbackJSON)}
	err := v.ValidateReader(&r)
	tt.NotNil(t, err)
}

func TestValidatorValidateReaderEOF(t *testing.T) {
	var v oj.Validator
	err := v.ValidateReader(iotest.DataErrReader(strings.NewReader("[1,2]")))
	tt.Nil(t, err)
}

func TestValidatorValidateReaderErr(t *testing.T) {
	var v oj.Validator
	err := v.ValidateReader(iotest.DataErrReader(strings.NewReader("[1,2}")))
	tt.NotNil(t, err)

	r := tt.ShortReader{Max: 5000, Content: []byte("[ 123" + strings.Repeat(",  123", 120) + "]")}
	err = v.ValidateReader(&r)
	tt.NotNil(t, err)
}
