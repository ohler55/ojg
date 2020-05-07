// Copyright (c) 2020, Peter Ohler, All rights reserved.

package ojg_test

import (
	"bytes"
	"fmt"
	"strings"
	"testing"

	"github.com/ohler55/ojg"
	"github.com/ohler55/ojg/gd"
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
			tt.Equal(t, d.expect, err.Error(), d.src)
		} else {
			tt.Nil(t, err, d.src)
		}
	}
}

func TestParseString(t *testing.T) {
	for _, d := range []data{
		{src: "null", value: nil},
		{src: "true", value: true},
		{src: "false", value: false},
		{src: "123", value: 123},
		{src: "-321", value: -321},
		{src: "12.3", value: 12.3},
		{src: `12345678901234567890`, value: gd.Big("12345678901234567890")},
		{src: `9223372036854775807`, value: 9223372036854775807},             // max int
		{src: `9223372036854775808`, value: gd.Big("9223372036854775808")},   // max int + 1
		{src: `-9223372036854775807`, value: -9223372036854775807},           // min int
		{src: `-9223372036854775808`, value: gd.Big("-9223372036854775808")}, // min int -1
		{src: `0.9223372036854775808`, value: gd.Big("0.9223372036854775808")},
		{src: `123456789012345678901234567890`, value: gd.Big("123456789012345678901234567890")},
		{src: `0.123456789012345678901234567890`, value: gd.Big("0.123456789012345678901234567890")},
		{src: `0.1e20000`, value: gd.Big("0.1e20000")},
		{src: `1.2e1025`, value: gd.Big("1.2e1025")},
		{src: `-1.2e-1025`, value: gd.Big("-1.2e-1025")},

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
	} {
		var err error
		var v interface{}
		if d.onlyOne || d.noComment {
			p := ojg.Parser{NoComment: d.noComment}
			v, err = p.Parse([]byte(d.src))
		} else {
			v, err = ojg.Parse([]byte(d.src))
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

func TestParseSimpleString(t *testing.T) {
	for _, d := range []data{
		{src: "null", value: nil},
		{src: "true", value: true},
		{src: "false", value: false},
		{src: "123", value: 123},
		{src: "-321", value: -321},
		{src: "12.3", value: 12.3},
		{src: "-12.3", value: -12.3},
		{src: "-12.3e-5", value: -12.3e-5},
		{src: `12345678901234567890`, value: "12345678901234567890"},
		{src: `9223372036854775807`, value: 9223372036854775807},     // max int
		{src: `9223372036854775808`, value: "9223372036854775808"},   // max int + 1
		{src: `-9223372036854775807`, value: -9223372036854775807},   // min int
		{src: `-9223372036854775808`, value: "-9223372036854775808"}, // min int -1
		{src: `0.9223372036854775808`, value: "0.9223372036854775808"},
		{src: `-0.9223372036854775808`, value: "-0.9223372036854775808"},
		{src: `1.2e1025`, value: "1.2e1025"},
		{src: `-1.2e-1025`, value: "-1.2e-1025"},

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
	} {
		var err error
		var v interface{}
		if d.onlyOne || d.noComment {
			p := ojg.Parser{NoComment: d.noComment}
			v, err = p.ParseSimple([]byte(d.src))
		} else {
			v, err = ojg.ParseSimple([]byte(d.src))
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
	// TBD try a longer JSON generate
}

const callbackJSON = `
1
[2]
{"x":3}
true false 123
`

func TestParseCallback(t *testing.T) {
	var results []byte
	cb := func(n gd.Node) bool {
		if 0 < len(results) {
			results = append(results, ' ')
		}
		results = append(results, n.String()...)
		return false
	}
	var p ojg.Parser
	v, err := p.Parse([]byte(callbackJSON), cb)
	tt.Nil(t, err)
	tt.Nil(t, v)
	tt.Equal(t, `1 [2] {"x":3} true false 123`, string(results))
}

func TestParseReaderCallback(t *testing.T) {
	var results []byte
	cb := func(n gd.Node) bool {
		if 0 < len(results) {
			results = append(results, ' ')
		}
		results = append(results, n.String()...)
		return false
	}
	var p ojg.Parser
	v, err := p.ParseReader(strings.NewReader(callbackJSON), cb)
	tt.Nil(t, err)
	tt.Nil(t, v)
	tt.Equal(t, `1 [2] {"x":3} true false 123`, string(results))
}

func TestParseSimpleCallback(t *testing.T) {
	var results []byte
	cb := func(n interface{}) bool {
		if 0 < len(results) {
			results = append(results, ' ')
		}
		results = append(results, fmt.Sprintf("%v", n)...)
		return false
	}
	var p ojg.Parser
	v, err := p.ParseSimple([]byte(callbackJSON), cb)
	tt.Nil(t, err)
	tt.Nil(t, v)
	tt.Equal(t, `1 [2] map[x:3] true false 123`, string(results))
}

func TestParseSimpleReaderCallback(t *testing.T) {
	var results []byte
	cb := func(n interface{}) bool {
		if 0 < len(results) {
			results = append(results, ' ')
		}
		results = append(results, fmt.Sprintf("%v", n)...)
		return false
	}
	var p ojg.Parser
	v, err := p.ParseSimpleReader(strings.NewReader(callbackJSON), cb)
	tt.Nil(t, err)
	tt.Nil(t, v)
	tt.Equal(t, `1 [2] map[x:3] true false 123`, string(results))
}

func TestNumberReset(t *testing.T) {
	var p ojg.Parser
	_, err := p.Parse([]byte("123456789012345678901234567890 1234567890"), func(gd.Node) bool { return false })
	tt.Nil(t, err)
}

func TestOjgParseString(t *testing.T) {
	v, err := ojg.ParseString("true")
	tt.Nil(t, err)
	tt.Equal(t, true, v)
}

func TestOjgParseSimpleString(t *testing.T) {
	v, err := ojg.ParseSimpleString("true")
	tt.Nil(t, err)
	tt.Equal(t, true, v)
}

func TestOjgLoad(t *testing.T) {
	v, err := ojg.Load(strings.NewReader("true"))
	tt.Nil(t, err)
	tt.Equal(t, true, v)
}

func TestOjgLoadSimple(t *testing.T) {
	v, err := ojg.LoadSimple(strings.NewReader("true"))
	tt.Nil(t, err)
	tt.Equal(t, true, v)
}

func TestOjgValidateString(t *testing.T) {
	err := ojg.ValidateString("true")
	tt.Nil(t, err)
}

/*
func TestDev(t *testing.T) {
	for _, d := range []data{
		{src: `1.2e200`, value: gd.Big("0.9223372036854775808")},
	} {
		var err error
		var v interface{}
		if d.onlyOne || d.noComment {
			p := ojg.Parser{NoComment: d.noComment}
			v, err = p.Parse([]byte(d.src))
		} else {
			v, err = ojg.Parse([]byte(d.src))
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
*/
