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
		{src: `{
  // a comment
}`, expect: ""},
		{src: `{
  // a comment
}`, expect: "did not expect '/' at 2:3", strict: true},
		{src: `{x}`, expect: "did not expect 'x' at 1:2"},

		{src: "[]", expect: ""},
		{src: "null {}", expect: ""},
		{src: "null {}", expect: "extra characters at 1:6", limit: 1},
		{src: "[true]", expect: ""},
		{src: "[true,false]", expect: ""},
	} {
		err := ojg.Validate(d.src, d.strict, d.limit)
		if 0 < len(d.expect) {
			tt.NotNil(t, err)
			tt.Equal(t, d.expect, err.Error(), d.src)
		} else {
			tt.Nil(t, err)
		}
	}
}
