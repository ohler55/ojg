// Copyright (c) 2020, Peter Ohler, All rights reserved.

package oj_test

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/ohler55/ojg"
	"github.com/ohler55/ojg/gen"
	"github.com/ohler55/ojg/oj"
	"github.com/ohler55/ojg/tt"
)

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
		TimeColor:   "t",
		NoColor:     "x",
	}
	tm := time.Date(2020, time.May, 7, 19, 29, 19, 123456789, time.UTC)
	for i, d := range []data{
		{value: nil, expect: "nnullx"},
		{value: true, expect: "btruex"},
		{value: false, expect: "bfalsex"},
		{value: "string", expect: `q"string"x`},
		{value: gen.String("string"), expect: `q"string"x`},
		{value: []interface{}{true, false}, expect: "s[xbtruexs,xbfalsexs]x"},
		{value: gen.Array{gen.Bool(true), gen.Bool(false)}, expect: "s[xbtruexs,xbfalsexs]x"},
		{value: gen.Object{"f": gen.False}, expect: `s{xk"f"xs:xbfalsexs}x`},
		{value: gen.Object{"f": gen.False}, expect: `s{xk"f"xs:xbfalsexs}x`, options: &oj.Options{Sort: true}},
		{value: map[string]interface{}{"t": true, "f": false},
			expect: `s{xk"f"xs:xbfalsexs,xk"t"xs:xbtruexs}x`, options: &oj.Options{Sort: true}},
		{value: gen.Array{gen.True, gen.False}, expect: "s[xbtruexs,xbfalsexs]x"},
		{value: gen.Array{gen.False, gen.True}, expect: "s[xbfalsexs,xbtruexs]x"},
		{value: []interface{}{-1, int8(2), int16(-3), int32(4), int64(-5)},
			expect: "s[x0-1xs,x02xs,x0-3xs,x04xs,x0-5xs]x"},
		{value: []interface{}{uint(1), 'A', uint8(2), uint16(3), uint32(4), uint64(5)},
			expect: "s[x01xs,x065xs,x02xs,x03xs,x04xs,x05xs]x"},
		{value: gen.Array{gen.Int(1), gen.Float(1.2)}, expect: "s[x01xs,x01.2xs]x"},
		{value: []interface{}{float32(1.2), float64(2.1)}, expect: "s[x01.2xs,x02.1xs]x"},
		{value: []interface{}{tm}, expect: "s[xt1588879759123456789xs]x"},
		{value: gen.Array{gen.Time(tm)}, expect: "s[xt1588879759123456789xs]x"},

		{value: map[string]interface{}{"t": true, "x": nil}, expect: "s{xk\"t\"xs:xbtruexs}x",
			options: &oj.Options{OmitNil: true}},
		{value: map[string]interface{}{"t": true, "x": nil}, expect: "s{xk\"t\"xs:xbtruexs}x",
			options: &oj.Options{OmitNil: true, Sort: true}},
		{value: map[string]interface{}{"t": true, "f": false},
			expect:  "s{x\n  k\"f\"xs:x bfalsexs,x\n  k\"t\"xs:x btruex\ns}x",
			options: &oj.Options{Sort: true, Indent: 2}},
		{value: map[string]interface{}{"t": true},
			expect: "s{x\n  k\"t\"xs:x btruex\ns}x", options: &oj.Options{Indent: 2}},
		{value: gen.Object{"t": gen.True, "x": nil}, expect: "s{xk\"t\"xs:xbtruexs}x",
			options: &oj.Options{OmitNil: true}},
		{value: gen.Object{"t": gen.True, "x": nil}, expect: "s{xk\"t\"xs:xbtruexs}x",
			options: &oj.Options{OmitNil: true, Sort: true}},
		{value: gen.Object{"t": gen.True}, expect: "s{x\n  k\"t\"xs:x btruex\ns}x", options: &oj.Options{Indent: 2}},
		{value: gen.Object{"t": gen.True}, expect: "s{x\n  k\"t\"xs:x btruex\ns}x",
			options: &oj.Options{Indent: 2, Sort: true}},

		{value: &simon{x: 3}, expect: `s{xk"type"xs:xq"simon"xs,xk"x"xs:x03xs}x`, options: &oj.Options{Sort: true}},
		{value: &genny{val: 3}, expect: `s{xk"type"xs:xq"genny"xs,xk"val"xs:x03xs}x`, options: &oj.Options{Sort: true}},
		{value: &Dummy{Val: 3}, expect: `s{xk"^"xs:xq"Dummy"xs,xk"val"xs:x03xs}x`,
			options: &oj.Options{Sort: true, CreateKey: "^"}},
		{value: &Dummy{Val: 3}, expect: `s{xk"val"xs:x03xs}x`},
	} {
		if testing.Verbose() {
			fmt.Printf("... %d: %v\n", i, d.value)
		}
		var b strings.Builder
		var err error
		if d.options != nil {
			d.options.Color = true
			d.options.SyntaxColor = "s"
			d.options.KeyColor = "k"
			d.options.NullColor = "n"
			d.options.BoolColor = "b"
			d.options.NumberColor = "0"
			d.options.StringColor = "q"
			d.options.NoColor = "x"
			err = oj.Write(&b, d.value, d.options)
		} else {
			err = oj.Write(&b, d.value, opt)
		}
		tt.Nil(t, err)
		tt.Equal(t, d.expect, b.String(), fmt.Sprintf("%d: %v", i, d.value))
	}
}

func TestColorWide(t *testing.T) {
	var b strings.Builder
	opt := oj.Options{
		Color: true,
		// use visible character to make it easier to verify
		SyntaxColor: "s",
		KeyColor:    "k",
		NullColor:   "n",
		BoolColor:   "b",
		NumberColor: "0",
		StringColor: "q",
		NoColor:     "x",
		Indent:      300,
		WriteLimit:  2,
	}
	err := oj.Write(&b, []interface{}{[]interface{}{true, nil}}, &opt)
	tt.Nil(t, err)
	tt.Equal(t, 544, len(b.String()))

	b.Reset()
	err = oj.Write(&b, gen.Array{gen.Array{gen.True, nil}}, &opt)
	tt.Nil(t, err)
	tt.Equal(t, 544, len(b.String()))

	b.Reset()
	err = oj.Write(&b, map[string]interface{}{"x": map[string]interface{}{"y": true, "z": nil}}, &opt)
	tt.Nil(t, err)
	tt.Equal(t, 571, len(b.String()))

	b.Reset()
	err = oj.Write(&b, gen.Object{"x": gen.Object{"y": gen.True, "z": nil}}, &opt)
	tt.Nil(t, err)
	tt.Equal(t, 571, len(b.String()))
}

func TestColorDeep(t *testing.T) {
	var b strings.Builder
	opt := oj.Options{
		Color: true,
		// use visible character to make it easier to verify
		SyntaxColor: "s",
		KeyColor:    "k",
		NullColor:   "n",
		BoolColor:   "b",
		NumberColor: "0",
		StringColor: "q",
		NoColor:     "x",
		Tab:         true,
		WriteLimit:  2,
	}
	a := []interface{}{map[string]interface{}{"x": true}}
	for i := 40; 0 < i; i-- {
		a = []interface{}{a}
	}
	err := oj.Write(&b, a, &opt)
	tt.Nil(t, err)
	tt.Equal(t, 1970, len(b.String()))

	b.Reset()
	g := gen.Array{gen.Object{"x": gen.True}}
	for i := 40; 0 < i; i-- {
		g = gen.Array{g}
	}
	err = oj.Write(&b, g, &opt)
	tt.Nil(t, err)
	tt.Equal(t, 1970, len(b.String()))
}

func TestColorShort(t *testing.T) {
	opt := oj.Options{
		Color: true,
		// use visible character to make it easier to verify
		SyntaxColor: "s",
		KeyColor:    "k",
		NullColor:   "n",
		BoolColor:   "b",
		NumberColor: "0",
		StringColor: "q",
		NoColor:     "x",
		Indent:      2,
		WriteLimit:  2,
	}
	err := oj.Write(&shortWriter{max: 3}, []interface{}{true, nil}, &opt)
	tt.NotNil(t, err)
	err = oj.Write(&shortWriter{max: 3}, gen.Array{gen.True, nil}, &opt)
	tt.NotNil(t, err)

	opt.Indent = 0
	err = oj.Write(&shortWriter{max: 3}, []interface{}{true, nil}, &opt)
	tt.NotNil(t, err)
	err = oj.Write(&shortWriter{max: 3}, gen.Array{gen.True, nil}, &opt)
	tt.NotNil(t, err)

	obj := map[string]interface{}{"t": true, "n": nil}
	sobj := gen.Object{"t": gen.True, "n": nil}
	err = oj.Write(&shortWriter{max: 7}, obj, &opt)
	tt.NotNil(t, err)
	err = oj.Write(&shortWriter{max: 7}, sobj, &opt)
	tt.NotNil(t, err)

	opt.Sort = true
	err = oj.Write(&shortWriter{max: 7}, obj, &opt)
	tt.NotNil(t, err)
	err = oj.Write(&shortWriter{max: 7}, sobj, &opt)
	tt.NotNil(t, err)

	opt.Indent = 2
	err = oj.Write(&shortWriter{max: 11}, obj, &opt)
	tt.NotNil(t, err)
	err = oj.Write(&shortWriter{max: 11}, sobj, &opt)
	tt.NotNil(t, err)

	opt.Sort = false
	err = oj.Write(&shortWriter{max: 11}, obj, &opt)
	tt.NotNil(t, err)
	err = oj.Write(&shortWriter{max: 11}, sobj, &opt)
	tt.NotNil(t, err)
}

func TestColorMarshal(t *testing.T) {
	opt := oj.Options{
		Color: true,
		// use visible character to make it easier to verify
		SyntaxColor: "s",
		KeyColor:    "k",
		NullColor:   "n",
		BoolColor:   "b",
		NumberColor: "0",
		StringColor: "q",
		NoColor:     "x",
		NoReflect:   true,
	}
	var b strings.Builder

	err := oj.Write(&b, []interface{}{true, &Dummy{Val: 3}}, &opt)
	tt.Nil(t, err)
	tt.Equal(t, `s[xbtruexs,x"\u0026{3}"xs]x`, b.String())
}

func TestColorMustJSON(t *testing.T) {
	wr := oj.Writer{Options: ojg.Options{
		Color: true,
		// use visible character to make it easier to verify
		SyntaxColor: "s",
		KeyColor:    "k",
		NullColor:   "n",
		BoolColor:   "b",
		NumberColor: "0",
		StringColor: "q",
		TimeColor:   "t",
		NoColor:     "x",
	}}
	tt.Equal(t, "btruex", string(wr.MustJSON(true)))
}
