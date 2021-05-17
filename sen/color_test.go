// Copyright (c) 2020, Peter Ohler, All rights reserved.

package sen_test

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/ohler55/ojg"
	"github.com/ohler55/ojg/gen"
	"github.com/ohler55/ojg/sen"
	"github.com/ohler55/ojg/tt"
)

func TestColor(t *testing.T) {
	opt := &sen.Options{
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
	for i, d := range []wdata{
		{value: nil, expect: "nnullx"},
		{value: true, expect: "btruex"},
		{value: false, expect: "bfalsex"},
		{value: "string", expect: `qstringx`},
		{value: gen.String("string"), expect: `qstringx`},
		{value: []interface{}{true, false}, expect: "s[xbtruex bfalsexs]x"},
		{value: gen.Array{gen.Bool(true), gen.Bool(false)}, expect: "s[xbtruex bfalsexs]x"},
		{value: gen.Object{"f": gen.False}, expect: `s{xkfxs:xbfalsexs}x`},
		{value: gen.Object{"f": gen.False}, expect: `s{xkfxs:xbfalsexs}x`, options: &sen.Options{Sort: true}},
		{value: map[string]interface{}{"t": true, "f": false},
			expect: `s{xkfxs:xbfalsex ktxs:xbtruexs}x`, options: &sen.Options{Sort: true}},
		{value: gen.Array{gen.True, gen.False}, expect: "s[xbtruex bfalsexs]x"},
		{value: gen.Array{gen.False, gen.True}, expect: "s[xbfalsex btruexs]x"},
		{value: []interface{}{-1, int8(2), int16(-3), int32(4), int64(-5)}, expect: "s[x0-1x 02x 0-3x 04x 0-5xs]x"},
		{value: []interface{}{uint(1), 'A', uint8(2), uint16(3), uint32(4), uint64(5)}, expect: "s[x01x 065x 02x 03x 04x 05xs]x"},
		{value: gen.Array{gen.Int(1), gen.Float(1.2)}, expect: "s[x01x 01.2xs]x"},
		{value: []interface{}{float32(1.2), float64(2.1)}, expect: "s[x01.2x 02.1xs]x"},
		{value: []interface{}{tm}, expect: "s[xt1588879759123456789xs]x"},
		{value: gen.Array{gen.Time(tm)}, expect: "s[xt1588879759123456789xs]x"},

		{value: map[string]interface{}{"t": true, "x": nil}, expect: "s{xktxs:xbtruexs}x",
			options: &sen.Options{OmitNil: true}},
		{value: map[string]interface{}{"t": true, "x": nil}, expect: "s{xktxs:xbtruexs}x",
			options: &sen.Options{OmitNil: true, Sort: true}},
		{value: map[string]interface{}{"t": true, "f": false}, expect: "s{x\n  kfxs:x bfalsex\n  ktxs:x btruex\ns}x",
			options: &sen.Options{Sort: true, Indent: 2}},
		{value: map[string]interface{}{"t": true}, expect: "s{x\n  ktxs:x btruex\ns}x", options: &sen.Options{Indent: 2}},
		{value: gen.Object{"t": gen.True, "x": nil}, expect: "s{xktxs:xbtruexs}x",
			options: &sen.Options{OmitNil: true}},
		{value: gen.Object{"t": gen.True, "x": nil}, expect: "s{xktxs:xbtruexs}x",
			options: &sen.Options{OmitNil: true, Sort: true}},
		{value: gen.Object{"t": gen.True}, expect: "s{x\n  ktxs:x btruex\ns}x", options: &sen.Options{Indent: 2}},
		{value: gen.Object{"t": gen.True}, expect: "s{x\n  ktxs:x btruex\ns}x", options: &sen.Options{Indent: 2, Sort: true}},

		{value: &simon{x: 3}, expect: `s{xktypexs:xqsimonx kxxs:x03xs}x`, options: &sen.Options{Sort: true}},
		{value: &genny{val: 3}, expect: `s{xktypexs:xqgennyx kvalxs:x03xs}x`, options: &sen.Options{Sort: true}},
		{value: &Dummy{Val: 3}, expect: `s{xk^xs:xqDummyx kvalxs:x03xs}x`, options: &sen.Options{Sort: true, CreateKey: "^"}},
		{value: &Dummy{Val: 3}, expect: `s{xkvalxs:x03xs}xx`, options: &sen.Options{CreateKey: ""}},
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
			err = sen.Write(&b, d.value, d.options)
		} else {
			err = sen.Write(&b, d.value, opt)
		}
		tt.Nil(t, err)
		tt.Equal(t, d.expect, b.String(), fmt.Sprintf("%d: %v", i, d.value))
	}
}

func TestColorWide(t *testing.T) {
	var b strings.Builder
	opt := sen.Options{
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
	err := sen.Write(&b, []interface{}{[]interface{}{true, nil}}, &opt)
	tt.Nil(t, err)
	tt.Equal(t, 541, len(b.String()))

	b.Reset()
	err = sen.Write(&b, gen.Array{gen.Array{gen.True, nil}}, &opt)
	tt.Nil(t, err)
	tt.Equal(t, 541, len(b.String()))

	b.Reset()
	err = sen.Write(&b, map[string]interface{}{"x": map[string]interface{}{"y": true, "z": nil}}, &opt)
	tt.Nil(t, err)
	tt.Equal(t, 562, len(b.String()))

	b.Reset()
	err = sen.Write(&b, gen.Object{"x": gen.Object{"y": gen.True, "z": nil}}, &opt)
	tt.Nil(t, err)
	tt.Equal(t, 562, len(b.String()))
}

func TestColorDeep(t *testing.T) {
	var b strings.Builder
	opt := sen.Options{
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
	a := []interface{}{map[string]interface{}{"x": true, "y": false}}
	for i := 40; 0 < i; i-- {
		a = []interface{}{a}
	}
	err := sen.Write(&b, a, &opt)
	tt.Nil(t, err)
	tt.Equal(t, 2012, len(b.String()))

	b.Reset()
	g := gen.Array{gen.Object{"x": gen.True, "y": gen.False}}
	for i := 40; 0 < i; i-- {
		g = gen.Array{g}
	}
	err = sen.Write(&b, g, &opt)
	tt.Nil(t, err)
	tt.Equal(t, 2012, len(b.String()))
}

func TestColorShort(t *testing.T) {
	opt := sen.Options{
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
	err := sen.Write(&shortWriter{max: 3}, []interface{}{true, nil}, &opt)
	tt.NotNil(t, err)
	err = sen.Write(&shortWriter{max: 3}, gen.Array{gen.True, nil}, &opt)
	tt.NotNil(t, err)

	opt.Indent = 0
	err = sen.Write(&shortWriter{max: 3}, []interface{}{true, nil}, &opt)
	tt.NotNil(t, err)
	err = sen.Write(&shortWriter{max: 3}, gen.Array{gen.True, nil}, &opt)
	tt.NotNil(t, err)

	obj := map[string]interface{}{"t": true, "n": nil}
	sobj := gen.Object{"t": gen.True, "n": nil}
	err = sen.Write(&shortWriter{max: 7}, obj, &opt)
	tt.NotNil(t, err)
	err = sen.Write(&shortWriter{max: 7}, sobj, &opt)
	tt.NotNil(t, err)

	opt.Sort = true
	err = sen.Write(&shortWriter{max: 7}, obj, &opt)
	tt.NotNil(t, err)
	err = sen.Write(&shortWriter{max: 7}, sobj, &opt)
	tt.NotNil(t, err)

	opt.Indent = 2
	err = sen.Write(&shortWriter{max: 11}, obj, &opt)
	tt.NotNil(t, err)
	err = sen.Write(&shortWriter{max: 11}, sobj, &opt)
	tt.NotNil(t, err)

	opt.Sort = false
	err = sen.Write(&shortWriter{max: 11}, obj, &opt)
	tt.NotNil(t, err)
	err = sen.Write(&shortWriter{max: 11}, sobj, &opt)
	tt.NotNil(t, err)
}

func TestColorObject(t *testing.T) {
	opt := sen.Options{
		Color: true,
		// use visible character to make it easier to verify
		SyntaxColor: "s",
		KeyColor:    "k",
		NullColor:   "n",
		BoolColor:   "b",
		NumberColor: "0",
		StringColor: "q",
		NoColor:     "x",
		Indent:      0,
		WriteLimit:  2,
	}
	var b strings.Builder
	err := sen.Write(&b, map[string]interface{}{"a": 1, "b": 3}, &opt)
	tt.Nil(t, err)
	tt.Equal(t, 25, len(b.String()))

	b.Reset()
	err = sen.Write(&b, gen.Object{"a": gen.True, "b": gen.False}, &opt)
	tt.Nil(t, err)
	tt.Equal(t, 32, len(b.String()))
}

func TestColorMustSen(t *testing.T) {
	wr := sen.Writer{Options: ojg.Options{
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
	tt.Equal(t, "btruex", string(wr.MustSEN(true)))
}
