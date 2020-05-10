// Copyright (c) 2020, Peter Ohler, All rights reserved.

package ojg_test

import (
	"strings"
	"testing"

	"github.com/ohler55/ojg"
	"github.com/ohler55/ojg/tt"
)

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
		{src: `1.2e200`, value: gen.Big("0.9223372036854775808")},
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
