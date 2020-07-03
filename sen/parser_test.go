// Copyright (c) 2020, Peter Ohler, All rights reserved.

package sen_test

import (
	"fmt"
	"testing"

	"github.com/ohler55/ojg/sen"
	"github.com/ohler55/ojg/tt"
)

type rdata struct {
	src string
	// Empty means no error expected while non empty should be compared
	// err.Error().
	expect string
	value  interface{}
	//onlyOne bool
}

func TestParserParseString(t *testing.T) {
	for i, d := range []rdata{
		{src: "null", value: nil},
		{src: "true", value: true},
		{src: "false", value: false},
		{src: "false \n ", value: false},
		{src: "hello ", value: "hello"},
		{src: `"hello"`, value: "hello"},
		{src: "[one two]", value: []interface{}{"one", "two"}},
		{src: "123", value: 123},
		{src: "-12.3", value: -12.3},
		{src: "-12.5e-2", value: -0.125},
		{src: "0", value: 0},
		{src: "[0 1,2]", value: []interface{}{0, 1, 2}},
		{src: "[0[1[2]]]", value: []interface{}{0, []interface{}{1, []interface{}{2}}}},
		{src: `{a:0 "bbb":1 , c:2}`, value: map[string]interface{}{"a": 0, "bbb": 1, "c": 2}},
		{src: `{a:one "b": "two" c : see}`, value: map[string]interface{}{"a": "one", "b": "two", "c": "see"}},
	} {
		if testing.Verbose() {
			fmt.Printf("... %q\n", d.src)
		}
		var err error
		var v interface{}
		var p sen.Parser
		v, err = p.Parse([]byte(d.src))

		if 0 < len(d.expect) {
			tt.NotNil(t, err, d.src)
			tt.Equal(t, d.expect, err.Error(), i, ": ", d.src)
		} else {
			tt.Nil(t, err, d.src)
			tt.Equal(t, d.value, v, i, ": ", d.src)
		}
	}
}
