// Copyright (c) 2020, Peter Ohler, All rights reserved.

package jp_test

import (
	"fmt"
	"testing"

	"github.com/ohler55/ojg/jp"
	"github.com/ohler55/ojg/tt"
)

type xdata struct {
	src    string
	expect string
	err    string
}

func TestParse(t *testing.T) {
	for i, d := range []xdata{
		{src: "@", expect: "@"},
		{src: "$", expect: "$"},
		{src: "@.abc", expect: "@.abc"},
		{src: "@.a.b.c", expect: "@.a.b.c"},
		{src: "$.abc", expect: "$.abc"},
		{src: "$.a.b.c", expect: "$.a.b.c"},
		{src: "abc", expect: "abc"},
		{src: "abc.def", expect: "abc.def"},
		{src: "abc.*.def", expect: "abc.*.def"},
		{src: "@..", expect: "@.."},
		{src: "@..x.y", expect: "@..x.y"},
		{src: "@.*", expect: "@.*"},
		{src: "[1,2]", expect: "[1,2]"},
		{src: "abc..def", expect: "abc..def"},
		{src: "abc[*].def", expect: "abc[*].def"},
		{src: "abc[0].def", expect: "abc[0].def"},
		{src: "abc[-1].def", expect: "abc[-1].def"},
		{src: "abc[2].def", expect: "abc[2].def"},
		{src: "abc[ -2 ].def", expect: "abc[-2].def"},
		{src: "abc[1:3]", expect: "abc[1:3]"},
		{src: "abc[1:]", expect: "abc[1:]"},
		{src: "abc[0:]", expect: "abc[:]"},
		{src: "abc[:]", expect: "abc[:]"},
		{src: "abc[:3]", expect: "abc[:3]"},
		{src: "abc[::2]", expect: "abc[::2]"},
		{src: "abc[1:5:2]", expect: "abc[1:5:2]"},
		{src: "abc[:-1]", expect: "abc[:]"},
		{src: "$['abc']", expect: "$.abc"},
		{src: "$['a b']", expect: "$['a b']"},
		{src: "$['ぴーたー']", expect: "$.ぴーたー"},
		{src: "$[1,2]", expect: "$[1,2]"},
		{src: "$[1,2,3]", expect: "$[1,2,3]"},
		{src: "$['a','b']", expect: "$['a','b']"},
		{src: "$[1,'a']", expect: "$[1,'a']"},
		{src: "$[1,'a',2,'b']", expect: "$[1,'a',2,'b']"},
		{src: "$[ 1, 'a' , 2 ,'b' ]", expect: "$[1,'a',2,'b']"},
		{src: "$[?(@.x == 'abc')]", expect: "$[?(@.x == 'abc')]"},
		{src: `['a\\b']`, expect: `['a\\b']`},

		{src: "$[1,'a']  ", err: "parse error at 9 in $[1,'a']  "},
		{src: "abc.", err: "not terminated at 5 in abc."},
		{src: "abc.+", err: "an expression fragment can not start with a '+' at 6 in abc.+"},
		{src: "abc..+", err: "parse error at 6 in abc..+"},
		{src: "[", err: "not terminated at 2 in ["},
		{src: "[**", err: "not terminated at 4 in [**"},
		{src: "['x'z]", err: "invalid bracket fragment at 6 in ['x'z]"},
		{src: "[(x)]", err: "scripts not implemented yet at 3 in [(x)]"},
		{src: "[-x]", err: "expected a number at 4 in [-x]"},
		{src: "[0x]", err: "invalid bracket fragment at 4 in [0x]"},
		{src: "[x]", err: "parse error at 3 in [x]"},
		{src: "[?(@.x == 1.2e", err: "expected a number at 15 in [?(@.x == 1.2e"},
		{src: "[?(@.x == 1e", err: "expected a number at 13 in [?(@.x == 1e"},
		{src: "[?(@.x == 1e+", err: "expected a number at 14 in [?(@.x == 1e+"},
		{src: "[-", err: "expected a number at 3 in [-"},
		{src: "[1", err: "invalid bracket fragment at 3 in [1"},
		{src: "[1,", err: "not terminated at 4 in [1,"},
		{src: "[:", err: "not terminated at 3 in [:"},
		{src: "[::", err: "not terminated at 4 in [::"},
		{src: "[:-x", err: "invalid slice syntax at 5 in [:-x"},
		{src: "[1:-x", err: "invalid slice syntax at 6 in [1:-x"},
		{src: "[1::-x", err: "expected a number at 7 in [1::-x"},
		{src: "[1:2:", err: "not terminated at 6 in [1:2:"},
		{src: "[1:2:-x", err: "expected a number at 8 in [1:2:-x"},
		{src: "[2,3:", err: "invalid union syntax at 6 in [2,3:"},
		{src: "[2,3x", err: "invalid union syntax at 6 in [2,3x"},
		{src: "[2,-", err: "expected a number at 5 in [2,-"},
		{src: "[2,x", err: "invalid union syntax at 5 in [2,x"},
		{src: "[?", err: "not terminated at 3 in [?"},
		{src: "[?(", err: "not terminated at 4 in [?("},
		{src: "[?x", err: "expected a '(' in filter at 4 in [?x"},
		{src: "[?(@.x == 3)", err: "not terminated at 13 in [?(@.x == 3)"},
		{src: "[?(!(@.x == -x)", err: `strconv.ParseInt: parsing "-": invalid syntax at 14 in [?(!(@.x == -x)`},
		{src: "[?(!(@.x == 1)]", err: "not terminated at 15 in [?(!(@.x == 1)]"},
		{src: "[?(- == 1)]", err: `strconv.ParseInt: parsing "-": invalid syntax at 5 in [?(- == 1)]`},
		{src: "[?(2 ++ 1)]", err: "'++' is not a valid operation at 8 in [?(2 ++ 1)]"},
		{src: "[?(2 + -)]", err: `strconv.ParseInt: parsing "-": invalid syntax at 9 in [?(2 + -)]`},
		{src: "[?(2 + 1 ++)]", err: `'++' is not a valid operation at 12 in [?(2 + 1 ++)]`},
		{src: "[?(2 + 1 + -)]", err: `strconv.ParseInt: parsing "-": invalid syntax at 13 in [?(2 + 1 + -)]`},
		{src: "[?(2 + 1 * -)]", err: `strconv.ParseInt: parsing "-": invalid syntax at 13 in [?(2 + 1 * -)]`},
		{src: "[?(@.x == trux)]", err: "expected true at 14 in [?(@.x == trux)]"},
		{src: "[?(@.x == fx)]", err: "expected false at 12 in [?(@.x == fx)]"},
		{src: "[?(@.x == nulx)]", err: "expected null at 14 in [?(@.x == nulx)]"},
		{src: "[?(@.x == x)]", err: "expected a value at 11 in [?(@.x == x)]"},
		{src: "[?(@.x -- x)]", err: "'--' is not a valid operation at 9 in [?(@.x -- x)]"},
		{src: "[?(@.x =", err: "equation not terminated at 9 in [?(@.x ="},
	} {
		if testing.Verbose() {
			fmt.Printf("... %s\n", d.src)
		}
		x, err := jp.ParseString(d.src)
		if 0 < len(d.err) {
			tt.NotNil(t, err, d.src)
			tt.Equal(t, d.err, err.Error(), i, ": ", d.src)
		} else {
			tt.Nil(t, err, d.src)
			tt.NotNil(t, x)
			tt.Equal(t, d.expect, x.String(), i, ": ", d.src)
		}
	}
}

func BenchmarkParse(b *testing.B) {
	b.ReportAllocs()
	for n := 0; n < b.N; n++ {
		_, _ = jp.Parse([]byte("@.abc.*[2,3]..xyz[2]"))
		//fmt.Printf("*** x: %s\n", x)
	}
}
