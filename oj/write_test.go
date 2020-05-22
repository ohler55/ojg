// Copyright (c) 2020, Peter Ohler, All rights reserved.

package oj_test

import (
	"fmt"
	"strings"
	"testing"
	"time"

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

func TestOjString(t *testing.T) {
	opt := &oj.Options{}
	tm := time.Date(2020, time.May, 7, 19, 29, 19, 123456789, time.UTC)
	for i, d := range []data{
		{value: nil, expect: "null"},
		{value: true, expect: "true"},
		{value: false, expect: "false"},
		{value: "string", expect: `"string"`},
		{value: "\\\t\n\r\b\f\"&<>\u2028\u2029\x07\U0001D122", expect: `"\\\t\n\r\b\f\"\u0026\u003c\u003e\u2028\u2029\u0007𝄢"`},
		{value: oj.String("string"), expect: `"string"`},
		{value: []interface{}{true, false}, expect: "[true,false]"},
		{value: oj.Array{oj.Bool(true), oj.Bool(false)}, expect: "[true,false]"},
		{value: []interface{}{true, false}, indent: 2, expect: "[\n  true,\n  false\n]"},
		{value: oj.Array{oj.True, oj.False}, indent: 2, expect: "[\n  true,\n  false\n]"},
		{value: oj.Object{"t": oj.True, "f": oj.False}, expect: `{"f":false,"t":true}`, options: &oj.Options{Sort: true}},
		{value: map[string]interface{}{"t": true, "f": false}, expect: `{"f":false,"t":true}`, options: &oj.Options{Sort: true}},
		{value: oj.Array{oj.True, oj.False}, expect: "[true,false]", options: opt},
		{value: oj.Array{oj.False, oj.True}, expect: "[false,true]", options: opt},
		{value: []interface{}{-1, int8(2), int16(-3), int32(4), int64(-5)}, expect: "[-1,2,-3,4,-5]"},
		{value: []interface{}{uint(1), 'A', uint8(2), uint16(3), uint32(4), uint64(5)}, expect: "[1,65,2,3,4,5]"},
		{value: oj.Array{oj.Int(1), oj.Float(1.2)}, expect: "[1,1.2]"},
		{value: []interface{}{float32(1.2), float64(2.1)}, expect: "[1.2,2.1]"},
		{value: []interface{}{tm}, expect: "[1588879759123456789]"},
		{value: oj.Array{oj.Time(tm)}, expect: "[1588879759123456789]"},
		{value: oj.Array{oj.Time(tm)}, expect: `["2020-05-07T19:29:19.123456789Z"]`, options: &oj.Options{TimeFormat: time.RFC3339Nano}},
		{value: oj.Array{oj.Time(tm)}, expect: "[1588879759.123456789]", options: &oj.Options{TimeFormat: "second"}},
		{value: oj.Array{oj.Time(tm)}, expect: `[{"@":1588879759123456789}]`, options: &oj.Options{TimeWrap: "@"}},
		{value: map[string]interface{}{"t": true, "x": nil}, expect: "{\"t\":true}", options: &oj.Options{OmitNil: true}},
		{value: map[string]interface{}{"t": true, "f": false}, expect: "{\n  \"f\": false,\n  \"t\": true\n}", options: &oj.Options{Sort: true, Indent: 2}},
		{value: map[string]interface{}{"t": true}, expect: "{\n  \"t\": true\n}", options: &oj.Options{Indent: 2}},
		{value: oj.Object{"t": oj.True, "x": nil}, expect: "{\"t\":true}", options: &oj.Options{OmitNil: true}},
		{value: oj.Object{"t": oj.True}, expect: "{\n  \"t\": true\n}", options: &oj.Options{Indent: 2}},
		{value: oj.Object{"t": oj.True}, expect: "{\n  \"t\": true\n}", options: &oj.Options{Indent: 2, Sort: true}},

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

func TestOjWrite(t *testing.T) {
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

func TestOjColor(t *testing.T) {
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
		{value: oj.String("string"), expect: `q"string"` + oj.Normal},
		{value: []interface{}{true, false}, expect: "s[btrues,bfalses]" + oj.Normal},
		{value: oj.Array{oj.Bool(true), oj.Bool(false)}, expect: "s[btrues,bfalses]" + oj.Normal},
		{value: oj.Object{"f": oj.False}, expect: `s{k"f"s:bfalses}` + oj.Normal},
		//{value: oj.Object{"f": oj.False}, expect: `s{k"f"s:bfalses}`, options: &oj.Options{Sort: true}},
		//{value: map[string]interface{}{"t": true, "f": false}, expect: `{"f":false,"t":true}`, options: &oj.Options{Sort: true}},
		//{value: oj.Array{oj.True, oj.False}, expect: "[true,false]" + oj.Normal},
		//{value: oj.Array{oj.False, oj.True}, expect: "[false,true]" + oj.Normal},
		{value: []interface{}{-1, int8(2), int16(-3), int32(4), int64(-5)}, expect: "s[0-1s,02s,0-3s,04s,0-5s]" + oj.Normal},
		{value: []interface{}{uint(1), 'A', uint8(2), uint16(3), uint32(4), uint64(5)}, expect: "s[01s,065s,02s,03s,04s,05s]" + oj.Normal},
		{value: oj.Array{oj.Int(1), oj.Float(1.2)}, expect: "s[01s,01.2s]" + oj.Normal},
		{value: []interface{}{float32(1.2), float64(2.1)}, expect: "s[01.2s,02.1s]" + oj.Normal},
		{value: []interface{}{tm}, expect: "s[q1588879759123456789s]" + oj.Normal},
		{value: oj.Array{oj.Time(tm)}, expect: "s[q1588879759123456789s]" + oj.Normal},

		//{value: map[string]interface{}{"t": true, "x": nil}, expect: "{\"t\":true}", options: &oj.Options{OmitNil: true}},
		//{value: map[string]interface{}{"t": true, "f": false}, expect: "{\n  \"f\":false,\n  \"t\":true\n}", options: &oj.Options{Sort: true, Indent: 2}},
		//{value: map[string]interface{}{"t": true}, expect: "{\n  \"t\":true\n}", options: &oj.Options{Indent: 2}},
		//{value: oj.Object{"t": oj.True, "x": nil}, expect: "{\"t\":true}", options: &oj.Options{OmitNil: true}},
		//{value: oj.Object{"t": oj.True}, expect: "{\n  \"t\":true\n}", options: &oj.Options{Indent: 2}},
		//{value: oj.Object{"t": oj.True}, expect: "{\n  \"t\":true\n}", options: &oj.Options{Indent: 2, Sort: true}},

		//{value: &simon{x: 3}, expect: `{"type":"simon","x":3}`, options: &oj.Options{Sort: true}},
	} {
		var b strings.Builder
		err := oj.Write(&b, d.value, opt)
		tt.Nil(t, err)
		tt.Equal(t, d.expect, b.String(), fmt.Sprintf("%d: %v", i, d.value))
	}
}
