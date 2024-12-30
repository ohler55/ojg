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
	jp.CompileScript = nil
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
		{src: "@['x']", expect: "@.x"},
		{src: "@['@x']", expect: "@['@x']"},
		{src: "[?@['@type'] == 'something']", expect: "[?(@['@type'] == 'something')]"},
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
		{src: "a[:].b", expect: "a[:].b"},
		{src: "abc[::2]", expect: "abc[::2]"},
		{src: "abc[1:5:2]", expect: "abc[1:5:2]"},
		{src: "abc[:-1]", expect: "abc[:-1]"},
		{src: "$['abc']", expect: "$.abc"},
		{src: "$['a b']", expect: "$['a b']"},
		{src: "$['ぴーたー']", expect: "$.ぴーたー"},
		{src: "$[1,2]", expect: "$[1,2]"},
		{src: "$[1,2,3]", expect: "$[1,2,3]"},
		{src: "$['a','b']", expect: "$['a','b']"},
		{src: "$[1,'a']", expect: "$[1,'a']"},
		{src: "$[1,'a',2,'b']", expect: "$[1,'a',2,'b']"},
		{src: "$[ 1, 'a' , 2 ,'b' ]", expect: "$[1,'a',2,'b']"},
		{src: "$[?(@.x == true)]", expect: "$[?(@.x == true)]"},
		{src: "$[?(@.x == false)]", expect: "$[?(@.x == false)]"},
		{src: "$[?(@.x == Nothing)]", expect: "$[?(@.x == Nothing)]"},
		{src: "$[?(@.x == null)]", expect: "$[?(@.x == null)]"},
		{src: "$[?(@.x == 'abc')]", expect: "$[?(@.x == 'abc')]"},
		{src: "$[?(1==1)]", expect: "$[?(1 == 1)]"},
		{src: "$[?(@.x)]", expect: "$[?(@.x)]"},
		{src: "$[?@.x]", expect: "$[?(@.x)]"},
		{src: `['a\\b']`, expect: `['a\\b']`},
		{src: `[:]`, expect: `[:]`},
		{src: `[::]`, expect: `[:]`},
		{src: `[1:2:]`, expect: `[1:2]`},
		{src: `[01:02:02]`, expect: `[1:2:2]`},
		{src: "[?@.x == 'abc']", expect: "[?(@.x == 'abc')]"},

		{src: "$[1,'a']  ", err: "parse error at 9 in $[1,'a']  "},
		{src: "abc.", err: "not terminated at 5 in abc."},
		{src: "abc.+", err: "an expression fragment can not start with a '+' at 6 in abc.+"},
		{src: "abc..+", err: "parse error at 6 in abc..+"},
		{src: `['\z']`, err: "0x7a (z) is not a valid escaped character"},
		{src: `['\xx']`, err: "0x78 (x) is not a valid hexadecimal character"},
		{src: `['\x1`, err: "0x31 (1) is not a valid escaped character"},
		{src: `['\x`, err: "0x78 (x) is not a valid escaped character"},
		{src: `['\u0a`, err: "0x61 (a) is not a valid escaped character"},
		{src: "[", err: "not terminated at 2 in ["},
		{src: "]", err: "parse error at 1 in ]"},
		{src: "[]", err: "parse error at 2 in []"},
		{src: "[**", err: "not terminated at 4 in [**"},
		{src: "['x'z]", err: "invalid bracket fragment at 6 in ['x'z]"},
		{src: "[(x)]", err: "jp.CompileScript has not been set"},
		{src: "[(x)", err: "not terminated at 3 in [(x)"},
		{src: "[-x]", err: "expected a number at 4 in [-x]"},
		{src: "[0x]", err: "invalid bracket fragment at 4 in [0x]"},
		{src: "[x]", err: "parse error at 2 in [x]"},
		{src: "[?(@.x == 1.2e", err: "expected a number at 15 in [?(@.x == 1.2e"},
		{src: "[?(@.x == 1e", err: "expected a number at 13 in [?(@.x == 1e"},
		{src: "[?(@.x == 1e+", err: "expected a number at 14 in [?(@.x == 1e+"},
		{src: "[-", err: "expected a number at 3 in [-"},
		{src: "[1", err: "expected a number at 3 in [1"},
		{src: "[1,", err: "not terminated at 4 in [1,"},
		{src: "[:", err: "not terminated at 3 in [:"},
		{src: "[::", err: "not terminated at 4 in [::"},
		{src: "[:-x", err: "expected a number at 5 in [:-x"},
		{src: "[1:-x", err: "expected a number at 6 in [1:-x"},
		{src: "[1::-x", err: "expected a number at 7 in [1::-x"},
		{src: "[1:2:", err: "not terminated at 6 in [1:2:"},
		{src: "[1:2:-x", err: "expected a number at 8 in [1:2:-x"},
		{src: "[2,3:", err: "invalid union syntax at 6 in [2,3:"},
		{src: "[2,3x", err: "invalid union syntax at 6 in [2,3x"},
		{src: "[2,-", err: "expected a number at 5 in [2,-"},
		{src: "[2,x", err: "invalid union syntax at 5 in [2,x"},
		{src: "[?", err: "not terminated at 3 in [?"},
		{src: "[?(", err: "'' is not a value or function at 4 in [?("},
		{src: "[?x", err: "'x' is not a value or function at 3 in [?x"},
		{src: "[?(@.x == 3)", err: "not terminated at 13 in [?(@.x == 3)"},
		{src: "[?(!(@.x == -x)", err: `strconv.ParseInt: parsing "-": invalid syntax at 14 in [?(!(@.x == -x)`},
		{src: "[?(!(@.x == 1)]", err: "not terminated at 15 in [?(!(@.x == 1)]"},
		{src: "[?(- == 1)]", err: `strconv.ParseInt: parsing "-": invalid syntax at 5 in [?(- == 1)]`},
		{src: "[?(2 ++ 1)]", err: "'++' is not a valid operation at 8 in [?(2 ++ 1)]"},
		{src: "[?(2 + -)]", err: `strconv.ParseInt: parsing "-": invalid syntax at 9 in [?(2 + -)]`},
		{src: "[?(2 + 1 ++)]", err: `'++' is not a valid operation at 12 in [?(2 + 1 ++)]`},
		{src: "[?(2 + 1 + -)]", err: `strconv.ParseInt: parsing "-": invalid syntax at 13 in [?(2 + 1 + -)]`},
		{src: "[?(2 + 1 * -)]", err: `strconv.ParseInt: parsing "-": invalid syntax at 13 in [?(2 + 1 * -)]`},
		{src: "[?(@.x == trux)]", err: "'trux' is not a value or function at 11 in [?(@.x == trux)]"},
		{src: "[?(@.x == fx)]", err: "'fx' is not a value or function at 11 in [?(@.x == fx)]"},
		{src: "[?(@.x == nulx)]", err: "'nulx' is not a value or function at 11 in [?(@.x == nulx)]"},
		{src: "[?(@.x == x)]", err: "'x' is not a value or function at 11 in [?(@.x == x)]"},
		{src: "[?(@.x -- x)]", err: "'--' is not a valid operation at 9 in [?(@.x -- x)]"},
		{src: "[?(@.x =", err: "equation not terminated at 9 in [?(@.x ="},
		{src: "[?(@.x in [1 2])]", err: "'' is not a valid operation at 14 in [?(@.x in [1 2])]"},
		{src: "$[?(@.x == North)]", err: "expected Nothing at 14 in $[?(@.x == North)]"},
		{src: "[?length]", err: "expected a length function at 9 in [?length]"},
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
		_ = jp.MustParse([]byte("@.abc.*[2,3]..xyz[2]"))
	}
}
