// Copyright (c) 2020, Peter Ohler, All rights reserved.

package ojg_test

import (
	"testing"

	"github.com/ohler55/ojg"
	"github.com/ohler55/ojg/tt"
)

func TestValidateString(t *testing.T) {
	type data struct {
		src string
		// Empty means no error expected while non empty should be compared
		// err.Error().
		expect  string
		options *ojg.ParseOptions
	}
	for _, d := range []data{
		{src: "null", expect: ""},
		{src: "true", expect: ""},
		{src: "false", expect: ""},
		{src: "[]", expect: ""},
		{src: "[true]", expect: ""},
		{src: "[true,false]", expect: ""},
		{src: "[[],[true],false]", expect: ""},
		{src: "[[],[true]false]", expect: "expected a comma or close, not 'f' at 1:11"},

		{src: "123", expect: ""},
		{src: "-1.23", expect: ""},

		{src: "[]", expect: ""},
		{src: "null {}", expect: ""},
		{src: "null {}", expect: "extra characters after close, '{' at 1:6", options: &ojg.ParseOptions{OnlyOne: true}},

		{src: "-1.23", expect: ""},
		{src: "+1.23", expect: "unexpected character '+' at 1:1"},
		{src: "1.23e+3", expect: ""},
		{src: "1.23e-3", expect: ""},
		{src: "1.23e3", expect: ""},
		{src: "1.2e3e3", expect: "invalid number '1.2e3e' at 1:6"},
		{src: "0.3", expect: ""},
		{src: "03", expect: "invalid number '03' at 1:2"},
		/*
			{src: "{}", expect: ""},
			{src: " { \t }  ", expect: ""},
			{src: "{\n  // a comment\n}", expect: ""},
			{src: "{\n  // a comment\n}", expect: "did not expect '/' at 2:3", strict: true},
			{src: `{x}`, expect: "did not expect 'x' at 1:2"},
		*/
	} {
		var err error
		if d.options != nil {
			err = ojg.Validate(d.src, d.options)
		} else {
			err = ojg.Validate(d.src)
		}
		if 0 < len(d.expect) {
			tt.NotNil(t, err, d.src)
			tt.Equal(t, d.expect, err.Error(), d.src)
		} else {
			tt.Nil(t, err, d.src)
		}
	}
}
