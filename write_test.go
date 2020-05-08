// Copyright (c) 2020, Peter Ohler, All rights reserved.

package ojg_test

import (
	"strings"
	"testing"
	"time"

	"github.com/ohler55/ojg"
	"github.com/ohler55/ojg/gd"
	"github.com/ohler55/ojg/tt"
)

func TestOjgString(t *testing.T) {
	opt := &ojg.Options{}
	tm := time.Date(2020, time.May, 7, 19, 29, 19, 123456789, time.UTC)
	for _, d := range []data{
		{value: nil, expect: "null"},
		{value: true, expect: "true"},
		{value: false, expect: "false"},
		{value: "string", expect: `"string"`},
		{value: "\\\t\n\r\b\f\"&<>\u2028\u2029", expect: `"\\\t\n\r\b\f\"\u0026\u003c\u003e\u2028\u2029"`},
		{value: gd.String("string"), expect: `"string"`},
		{value: []interface{}{true, false}, expect: "[true,false]"},
		{value: gd.Array{gd.Bool(true), gd.Bool(false)}, expect: "[true,false]"},
		{value: []interface{}{true, false}, indent: 2, expect: "[\n  true,\n  false\n]"},
		{value: gd.Array{gd.True, gd.False}, indent: 2, expect: "[\n  true,\n  false\n]"},
		{value: gd.Object{"t": gd.True, "f": gd.False}, expect: `{"f":false,"t":true}`, options: &ojg.Options{Sort: true}},
		{value: map[string]interface{}{"t": true, "f": false}, expect: `{"f":false,"t":true}`, options: &ojg.Options{Sort: true}},
		{value: gd.Array{gd.True, gd.False}, expect: "[true,false]", options: opt},
		{value: gd.Array{gd.False, gd.True}, expect: "[false,true]", options: opt},
		{value: []interface{}{-1, int8(2), int16(-3), int32(4), int64(-5)}, expect: "[-1,2,-3,4,-5]"},
		{value: []interface{}{uint(1), 'A', uint8(2), uint16(3), uint32(4), uint64(5)}, expect: "[1,65,2,3,4,5]"},
		{value: gd.Array{gd.Int(1), gd.Float(1.2)}, expect: "[1,1.2]"},
		{value: []interface{}{float32(1.2), float64(2.1)}, expect: "[1.2,2.1]"},
		{value: []interface{}{tm}, expect: "[1588879759123456789]"},
		{value: gd.Array{gd.Time(tm)}, expect: "[1588879759123456789]"},
		{value: gd.Array{gd.Time(tm)}, expect: `["2020-05-07T19:29:19.123456789Z"]`, options: &ojg.Options{TimeFormat: time.RFC3339Nano}},
		{value: gd.Array{gd.Time(tm)}, expect: "[1588879759.123456789]", options: &ojg.Options{TimeFormat: "second"}},
		{value: gd.Array{gd.Time(tm)}, expect: `[{"@":1588879759123456789}]`, options: &ojg.Options{TimeWrap: "@"}},
	} {
		var s string
		if d.options == nil {
			if 0 < d.indent {
				s = ojg.String(d.value, d.indent)
			} else {
				s = ojg.String(d.value)
			}
		} else {
			s = ojg.String(d.value, d.options)
		}
		tt.Equal(t, d.expect, s, d.src)
	}
}

func TestOjgWrite(t *testing.T) {
	var b strings.Builder

	err := ojg.Write(&b, []interface{}{true, false})
	tt.Nil(t, err)
	tt.Equal(t, "[true,false]", b.String())

	opt := ojg.Options{}
	b.Reset()
	err = ojg.Write(&b, []interface{}{true, false}, &opt)
	tt.Nil(t, err)
	tt.Equal(t, "[true,false]", b.String())

	// A second time.
	b.Reset()
	err = ojg.Write(&b, []interface{}{true, false}, &opt)
	tt.Nil(t, err)
	tt.Equal(t, "[true,false]", b.String())

	b.Reset()
	err = ojg.Write(&b, []interface{}{false, true}, 2)
	tt.Nil(t, err)
	tt.Equal(t, "[\n  false,\n  true\n]", b.String())
}
