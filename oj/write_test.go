// Copyright (c) 2020, Peter Ohler, All rights reserved.

package oj_test

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/ohler55/ojg/gen"
	"github.com/ohler55/ojg/oj"
	"github.com/ohler55/ojg/tt"
)

// Used to test Simplifier objects in simple data.
type simon struct {
	x int
}

func (s *simon) Simplify() interface{} {
	return map[string]interface{}{
		"type": "simon",
		"x":    s.x,
	}
}

func TestString(t *testing.T) {
	opt := &oj.Options{}
	tm := time.Date(2020, time.May, 7, 19, 29, 19, 123456789, time.UTC)
	for i, d := range []data{
		{value: nil, expect: "null"},
		{value: true, expect: "true"},
		{value: false, expect: "false"},
		{value: "string", expect: `"string"`},
		{value: "\\\t\n\r\b\f\"&<>\u2028\u2029\x07\U0001D122", expect: `"\\\t\n\r\b\f\"\u0026\u003c\u003e\u2028\u2029\u0007𝄢"`},
		{value: gen.String("string"), expect: `"string"`},
		{value: []interface{}{true, false}, expect: "[true,false]"},
		{value: gen.Array{gen.Bool(true), gen.Bool(false)}, expect: "[true,false]"},
		{value: []interface{}{true, false}, indent: 2, expect: "[\n  true,\n  false\n]"},
		{value: gen.Array{gen.True, gen.False}, indent: 2, expect: "[\n  true,\n  false\n]"},
		{value: gen.Object{"t": gen.True, "f": gen.False}, expect: `{"f":false,"t":true}`, options: &oj.Options{Sort: true}},
		{value: map[string]interface{}{"t": true, "f": false}, expect: `{"f":false,"t":true}`, options: &oj.Options{Sort: true}},
		{value: gen.Array{gen.True, gen.False}, expect: "[true,false]", options: opt},
		{value: gen.Array{gen.False, gen.True}, expect: "[false,true]", options: opt},
		{value: []interface{}{-1, int8(2), int16(-3), int32(4), int64(-5)}, expect: "[-1,2,-3,4,-5]"},
		{value: []interface{}{uint(1), 'A', uint8(2), uint16(3), uint32(4), uint64(5)}, expect: "[1,65,2,3,4,5]"},
		{value: gen.Array{gen.Int(1), gen.Float(1.2)}, expect: "[1,1.2]"},
		{value: []interface{}{float32(1.2), float64(2.1)}, expect: "[1.2,2.1]"},
		{value: []interface{}{tm}, expect: "[1588879759123456789]"},
		{value: gen.Array{gen.Time(tm)}, expect: "[1588879759123456789]"},
		{value: gen.Array{gen.Time(tm)}, expect: `["2020-05-07T19:29:19.123456789Z"]`, options: &oj.Options{TimeFormat: time.RFC3339Nano}},
		{value: gen.Array{gen.Time(tm)}, expect: "[1588879759.123456789]", options: &oj.Options{TimeFormat: "second"}},
		{value: gen.Array{gen.Time(tm)}, expect: `[{"@":1588879759123456789}]`, options: &oj.Options{TimeWrap: "@"}},
		{value: map[string]interface{}{"t": true, "x": nil}, expect: "{\"t\":true}", options: &oj.Options{OmitNil: true}},
		{value: map[string]interface{}{"t": true, "f": false}, expect: "{\n  \"f\": false,\n  \"t\": true\n}", options: &oj.Options{Sort: true, Indent: 2}},
		{value: map[string]interface{}{"t": true}, expect: "{\n  \"t\": true\n}", options: &oj.Options{Indent: 2}},
		{value: gen.Object{"t": gen.True, "x": nil}, expect: "{\"t\":true}", options: &oj.Options{OmitNil: true}},
		{value: gen.Object{"t": gen.True}, expect: "{\n  \"t\": true\n}", options: &oj.Options{Indent: 2}},
		{value: gen.Object{"t": gen.True}, expect: "{\n  \"t\": true\n}", options: &oj.Options{Indent: 2, Sort: true}},

		{value: &simon{x: 3}, expect: `{"type":"simon","x":3}`, options: &oj.Options{Sort: true}},
	} {
		var s string
		if d.options == nil {
			if 0 < d.indent {
				s = oj.JSON(d.value, d.indent)
			} else {
				s = oj.JSON(d.value)
			}
		} else {
			s = oj.JSON(d.value, d.options)
		}
		tt.Equal(t, d.expect, s, fmt.Sprintf("%d: %v", i, d.value))
	}
}

func TestWrite(t *testing.T) {
	var b strings.Builder

	err := oj.Write(&b, []interface{}{true, false})
	tt.Nil(t, err)
	tt.Equal(t, "[true,false]", b.String())

	opt := oj.Options{}
	b.Reset()
	err = oj.Write(&b, []interface{}{true, false}, &opt)
	tt.Nil(t, err)
	tt.Equal(t, "[true,false]", b.String())

	// A second time.
	b.Reset()
	err = oj.Write(&b, []interface{}{true, false}, &opt)
	tt.Nil(t, err)
	tt.Equal(t, "[true,false]", b.String())

	b.Reset()
	err = oj.Write(&b, []interface{}{false, true}, 2)
	tt.Nil(t, err)
	tt.Equal(t, "[\n  false,\n  true\n]", b.String())
}

func TestColor(t *testing.T) {
	opt := &oj.Options{
		Color: true,
		// use visible character to make it easier to verify
		SyntaxColor: "s",
		KeyColor:    "k",
		NullColor:   "n",
		BoolColor:   "b",
		NumberColor: "0",
		StringColor: "q",
	}
	tm := time.Date(2020, time.May, 7, 19, 29, 19, 123456789, time.UTC)
	for i, d := range []data{
		{value: nil, expect: "nnull" + oj.Normal},
		{value: true, expect: "btrue" + oj.Normal},
		{value: false, expect: "bfalse" + oj.Normal},
		{value: "string", expect: `q"string"` + oj.Normal},
		{value: gen.String("string"), expect: `q"string"` + oj.Normal},
		{value: []interface{}{true, false}, expect: "s[btrues,bfalses]" + oj.Normal},
		{value: gen.Array{gen.Bool(true), gen.Bool(false)}, expect: "s[btrues,bfalses]" + oj.Normal},
		{value: gen.Object{"f": gen.False}, expect: `s{k"f"s:bfalses}` + oj.Normal},
		//{value: gen.Object{"f": gen.False}, expect: `s{k"f"s:bfalses}`, options: &oj.Options{Sort: true}},
		//{value: map[string]interface{}{"t": true, "f": false}, expect: `{"f":false,"t":true}`, options: &oj.Options{Sort: true}},
		//{value: gen.Array{gen.True, gen.False}, expect: "[true,false]" + oj.Normal},
		//{value: gen.Array{gen.False, gen.True}, expect: "[false,true]" + oj.Normal},
		{value: []interface{}{-1, int8(2), int16(-3), int32(4), int64(-5)}, expect: "s[0-1s,02s,0-3s,04s,0-5s]" + oj.Normal},
		{value: []interface{}{uint(1), 'A', uint8(2), uint16(3), uint32(4), uint64(5)}, expect: "s[01s,065s,02s,03s,04s,05s]" + oj.Normal},
		{value: gen.Array{gen.Int(1), gen.Float(1.2)}, expect: "s[01s,01.2s]" + oj.Normal},
		{value: []interface{}{float32(1.2), float64(2.1)}, expect: "s[01.2s,02.1s]" + oj.Normal},
		{value: []interface{}{tm}, expect: "s[q1588879759123456789s]" + oj.Normal},
		{value: gen.Array{gen.Time(tm)}, expect: "s[q1588879759123456789s]" + oj.Normal},

		//{value: map[string]interface{}{"t": true, "x": nil}, expect: "{\"t\":true}", options: &oj.Options{OmitNil: true}},
		//{value: map[string]interface{}{"t": true, "f": false}, expect: "{\n  \"f\":false,\n  \"t\":true\n}", options: &oj.Options{Sort: true, Indent: 2}},
		//{value: map[string]interface{}{"t": true}, expect: "{\n  \"t\":true\n}", options: &oj.Options{Indent: 2}},
		//{value: gen.Object{"t": gen.True, "x": nil}, expect: "{\"t\":true}", options: &oj.Options{OmitNil: true}},
		//{value: gen.Object{"t": gen.True}, expect: "{\n  \"t\":true\n}", options: &oj.Options{Indent: 2}},
		//{value: gen.Object{"t": gen.True}, expect: "{\n  \"t\":true\n}", options: &oj.Options{Indent: 2, Sort: true}},

		//{value: &simon{x: 3}, expect: `{"type":"simon","x":3}`, options: &oj.Options{Sort: true}},
	} {
		var b strings.Builder
		err := oj.Write(&b, d.value, opt)
		tt.Nil(t, err)
		tt.Equal(t, d.expect, b.String(), fmt.Sprintf("%d: %v", i, d.value))
	}
}
