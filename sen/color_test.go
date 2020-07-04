// Copyright (c) 2020, Peter Ohler, All rights reserved.

package sen_test

import (
	"fmt"
	"strings"
	"testing"
	"time"

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
	}
	tm := time.Date(2020, time.May, 7, 19, 29, 19, 123456789, time.UTC)
	for i, d := range []wdata{
		{value: nil, expect: "nnull" + sen.Normal},
		{value: true, expect: "btrue" + sen.Normal},
		{value: false, expect: "bfalse" + sen.Normal},
		{value: "string", expect: `qstring` + sen.Normal},
		{value: gen.String("string"), expect: `qstring` + sen.Normal},
		{value: []interface{}{true, false}, expect: "s[btrues bfalses]" + sen.Normal},
		{value: gen.Array{gen.Bool(true), gen.Bool(false)}, expect: "s[btrues bfalses]" + sen.Normal},
		{value: gen.Object{"f": gen.False}, expect: `s{kfs:bfalses}` + sen.Normal},
		{value: gen.Object{"f": gen.False}, expect: `s{kfs:bfalses}` + sen.Normal, options: &sen.Options{Sort: true}},
		{value: map[string]interface{}{"t": true, "f": false},
			expect: `s{kfs:bfalses kts:btrues}` + sen.Normal, options: &sen.Options{Sort: true}},
		{value: gen.Array{gen.True, gen.False}, expect: "s[btrues bfalses]" + sen.Normal},
		{value: gen.Array{gen.False, gen.True}, expect: "s[bfalses btrues]" + sen.Normal},
		{value: []interface{}{-1, int8(2), int16(-3), int32(4), int64(-5)}, expect: "s[0-1s 02s 0-3s 04s 0-5s]" + sen.Normal},
		{value: []interface{}{uint(1), 'A', uint8(2), uint16(3), uint32(4), uint64(5)}, expect: "s[01s 065s 02s 03s 04s 05s]" + sen.Normal},
		{value: gen.Array{gen.Int(1), gen.Float(1.2)}, expect: "s[01s 01.2s]" + sen.Normal},
		{value: []interface{}{float32(1.2), float64(2.1)}, expect: "s[01.2s 02.1s]" + sen.Normal},
		{value: []interface{}{tm}, expect: "s[q1588879759123456789s]" + sen.Normal},
		{value: gen.Array{gen.Time(tm)}, expect: "s[q1588879759123456789s]" + sen.Normal},

		{value: map[string]interface{}{"t": true, "x": nil}, expect: "s{kts:btrues}" + sen.Normal,
			options: &sen.Options{OmitNil: true}},
		{value: map[string]interface{}{"t": true, "x": nil}, expect: "s{kts:btrues}" + sen.Normal,
			options: &sen.Options{OmitNil: true, Sort: true}},
		{value: map[string]interface{}{"t": true, "f": false}, expect: "s{\n  kfs: bfalse\n  kts: btrue\ns}" + sen.Normal,
			options: &sen.Options{Sort: true, Indent: 2}},
		{value: map[string]interface{}{"t": true}, expect: "s{\n  kts: btrue\ns}" + sen.Normal, options: &sen.Options{Indent: 2}},
		{value: gen.Object{"t": gen.True, "x": nil}, expect: "s{kts:btrues}" + sen.Normal,
			options: &sen.Options{OmitNil: true}},
		{value: gen.Object{"t": gen.True, "x": nil}, expect: "s{kts:btrues}" + sen.Normal,
			options: &sen.Options{OmitNil: true, Sort: true}},
		{value: gen.Object{"t": gen.True}, expect: "s{\n  kts: btrue\ns}" + sen.Normal, options: &sen.Options{Indent: 2}},
		{value: gen.Object{"t": gen.True}, expect: "s{\n  kts: btrue\ns}" + sen.Normal, options: &sen.Options{Indent: 2, Sort: true}},

		{value: &simon{x: 3}, expect: `s{ktypes:qsimons kxs:03s}` + sen.Normal, options: &sen.Options{Sort: true}},
		{value: &genny{val: 3}, expect: `s{ktypes:qgennys kvals:03s}` + sen.Normal, options: &sen.Options{Sort: true}},
		{value: &Dummy{Val: 3}, expect: `s{k^s:qDummys kvals:03s}` + sen.Normal, options: &sen.Options{Sort: true, CreateKey: "^"}},
		{value: &Dummy{Val: 3}, expect: `"{val: 3}"` + sen.Normal, options: &sen.Options{CreateKey: ""}},
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
		Indent:      300,
		WriteLimit:  2,
	}
	err := sen.Write(&b, []interface{}{[]interface{}{true, nil}}, &opt)
	tt.Nil(t, err)
	tt.Equal(t, 538, len(b.String()))

	b.Reset()
	err = sen.Write(&b, gen.Array{gen.Array{gen.True, nil}}, &opt)
	tt.Nil(t, err)
	tt.Equal(t, 538, len(b.String()))

	b.Reset()
	err = sen.Write(&b, map[string]interface{}{"x": map[string]interface{}{"y": true, "z": nil}}, &opt)
	tt.Nil(t, err)
	tt.Equal(t, 553, len(b.String()))

	b.Reset()
	err = sen.Write(&b, gen.Object{"x": gen.Object{"y": gen.True, "z": nil}}, &opt)
	tt.Nil(t, err)
	tt.Equal(t, 553, len(b.String()))
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
