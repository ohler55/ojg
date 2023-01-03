// Copyright (c) 2020, Peter Ohler, All rights reserved.

package gen_test

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/ohler55/ojg/gen"
	"github.com/ohler55/ojg/tt"
)

func TestNumber(t *testing.T) {
	for i, d := range []data{
		{src: "123", value: 123},
		{src: "-123", value: -123},
		{src: "1.25", value: 1.25},
		{src: "-1.25", value: -1.25},
		{src: "1.25e3", value: 1.25e3},
		{src: "-1.25e-1", value: -1.25e-1},
		{src: "12345678901234567890", value: json.Number("12345678901234567890")},
		{src: "0.12345678901234567890", value: "0.12345678901234567890"},
		{src: "0.9223372036854775808", value: "0.9223372036854775808"},
	} {
		if testing.Verbose() {
			fmt.Printf("... %d: %s\n", i, d.src)
		}
		var v any
		var n gen.Number
		n.Reset()
		frac := false
		exp := false
		for _, b := range []byte(d.src) {
			switch b {
			case '.':
				frac = true
			case '-':
				if exp {
					n.NegExp = true
				} else {
					n.Neg = true
				}
			case 'e':
				exp = true
			default:
				if frac {
					if exp {
						n.AddExp(b)
					} else {
						n.AddFrac(b)
					}
				} else {
					n.AddDigit(b)
				}
			}
		}
		v = n.AsNum()
		tt.Equal(t, d.value, v, ": ", d.src)
	}
}

func TestNumberForceFloat(t *testing.T) {
	var n gen.Number
	n.Reset()
	n.ForceFloat = true
	for _, b := range []byte("123") {
		n.AddDigit(b)
	}
	v := n.AsNum()
	tt.Equal(t, 123.0, v)
}
