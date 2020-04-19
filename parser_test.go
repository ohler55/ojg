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
		expect string
		strict bool
		limit  int
	}
	for _, d := range []data{
		{src: "{}", expect: ""},
		{src: " { \t }  ", expect: ""},
		{src: "{\n  // a comment\n}", expect: ""},
		{src: "{\n  // a comment\n}", expect: "did not expect '/' at 2:3", strict: true},
		{src: `{x}`, expect: "did not expect 'x' at 1:2"},

		{src: "[]", expect: ""},
		{src: "null {}", expect: ""},
		{src: "null {}", expect: "extra characters at 1:6", limit: 1},

		{src: "[true]", expect: ""},
		{src: "[true,false]", expect: ""},
		{src: "[[],[true],false]", expect: ""},
		{src: "[[],[true]false]", expect: "expected a comma or close, not 'f' at 1:11"},

		{src: "123", expect: ""},
		{src: "-1.23", expect: ""},
		{src: "+1.23", expect: ""},
		{src: "+1.23", expect: "numbers can not start with a '+' at 1:1", strict: true},
		{src: "1.23e+3", expect: ""},
		{src: "1.23e-3", expect: ""},
		{src: "1.23e3", expect: ""},
		{src: "1.2e3e3", expect: `strconv.ParseFloat: parsing "1.2e3e3": invalid syntax at 1:1`},
		{src: "0.3", expect: ""},
		{src: "03", expect: "numbers can not start with a '0' if not 0 at 1:1", strict: true},
	} {
		err := ojg.Validate(d.src, d.strict, d.limit)
		if 0 < len(d.expect) {
			tt.NotNil(t, err, d.src)
			tt.Equal(t, d.expect, err.Error(), d.src)
		} else {
			tt.Nil(t, err, d.src)
		}
	}
}
