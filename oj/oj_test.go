// Copyright (c) 2020, Peter Ohler, All rights reserved.

package oj_test

import (
	"strings"
	"testing"

	"github.com/ohler55/ojg/oj"
	"github.com/ohler55/ojg/tt"
)

func TestParseString(t *testing.T) {
	v, err := oj.ParseString("true")
	tt.Nil(t, err)
	tt.Equal(t, true, v)
}

func TestLoad(t *testing.T) {
	v, err := oj.Load(strings.NewReader("true"))
	tt.Nil(t, err)
	tt.Equal(t, true, v)
}

func TestValidateString(t *testing.T) {
	err := oj.ValidateString("true")
	tt.Nil(t, err)
}

/*
func TestDev(t *testing.T) {
	for _, d := range []data{
		{src: `1.2e200`, value: gen.Big("0.9223372036854775808")},
	} {
		var err error
		var v any
		if d.onlyOne || d.noComment {
			p := oj.Parser{NoComment: d.noComment}
			v, err = p.Parse([]byte(d.src))
		} else {
			v, err = oj.Parse([]byte(d.src))
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
