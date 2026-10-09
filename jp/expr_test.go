// Copyright (c) 2020, Peter Ohler, All rights reserved.

package jp_test

import (
	"testing"

	"github.com/ohler55/ojg/jp"
	"github.com/ohler55/ojg/tt"
)

func TestExprBuild(t *testing.T) {
	x := jp.X().D().C("abc").W().N(3).U(2, "x").S(1, 5, 2, 100).S(1, 5).S(1)
	tt.Equal(t, "..abc.*[3][2,'x'][1:5:2][1:5][1:]", x.String())

	x = jp.R().Descent().Child("abc").Wildcard().Nth(3).Union(int64(2), "x").Slice(1, 5, 2).Slice(1, 5).Slice(1)
	tt.Equal(t, "$..abc.*[3][2,'x'][1:5:2][1:5][1:]", x.String())

	x = jp.B().Descent().Child("abc").Wildcard()
	tt.Equal(t, "[..]['abc'][*]", x.String())

	x = jp.R().B().Descent().Child("abc").Wildcard()
	tt.Equal(t, "$[..]['abc'][*]", x.String())

	eq := jp.Lt(jp.Get(jp.A().C("a")), jp.ConstInt(52))
	x = jp.F(eq)
	tt.Equal(t, "[?(@.a < 52)]", x.String())

	x = jp.W().F(eq)
	tt.Equal(t, "*[?(@.a < 52)]", x.String())

	x = jp.B().R().Filter(eq)
	tt.Equal(t, "$[?(@.a < 52)]", x.String())

	x = jp.B().Root().W()
	tt.Equal(t, "$[*]", x.String())

	x = jp.B().A().W()
	tt.Equal(t, "@[*]", x.String())

	x = jp.B().At().W()
	tt.Equal(t, "@[*]", x.String())

	x = jp.N(3)
	tt.Equal(t, "[3]", x.String())

	x = jp.S(3, 4)
	tt.Equal(t, "[3:4]", x.String())

	x = jp.S(3, 4, 5, 6) // values after 3 args are ignored
	tt.Equal(t, "[3:4:5]", x.String())

	x = jp.R().Slice(3, 4, 5, 6) // values after 3 args are ignored
	tt.Equal(t, "$[3:4:5]", x.String())

	x = jp.D()
	tt.Equal(t, "..", x.String())

	x = jp.U(1, "a")
	tt.Equal(t, "[1,'a']", x.String())

	x = jp.Expr{jp.NewSlice()}
	tt.Equal(t, "[:]", x.String())

	x = jp.R().Child("'")
	tt.Equal(t, `$['\'']`, x.String())

	x = jp.R().Child("").Child("a")
	tt.Equal(t, `$[''].a`, x.String())

	x = jp.R().Child("a & b")
	tt.Equal(t, `$['a & b']`, x.String())

	x = jp.R().Child("cognito:username")
	tt.Equal(t, `$['cognito:username']`, x.String())

	x = jp.R().Child(":")
	tt.Equal(t, `$[':']`, x.String())

	x = jp.R().Child(":start")
	tt.Equal(t, `$[':start']`, x.String())

	x = jp.R().Child("end:")
	tt.Equal(t, `$['end:']`, x.String())

	x = jp.R().Child("a::b")
	tt.Equal(t, `$['a::b']`, x.String())
}

func TestExprUnionEscaping(t *testing.T) {
	for _, tc := range []struct {
		name string
		key  string
		want string
	}{
		{"plain", "plain", `['plain','other']`},
		{"empty", "", `['','other']`},
		{"quote", "a'b", `['a\'b','other']`},
		{"double quote", `a"b`, `['a\"b','other']`},
		{"backslash", `a\b`, `['a\\b','other']`},
		{"trailing backslash", `a\`, `['a\\','other']`},
		{"backslash quote", `a\'b`, `['a\\\'b','other']`},
		{"newline", "a\nb", `['a\nb','other']`},
		{"tab", "a\tb", `['a\tb','other']`},
		{"carriage return", "a\rb", `['a\rb','other']`},
		{"backspace", "a\bb", `['a\bb','other']`},
		{"form feed", "a\fb", `['a\fb','other']`},
		{"null", "a\x00b", `['a\u0000b','other']`},
		{"control", "a\x1fb", `['a\u001fb','other']`},
		{"line separator", "a\u2028b", `['a\u2028b','other']`},
		{"paragraph separator", "a\u2029b", `['a\u2029b','other']`},
		{"unicode", "東京🎾", `['東京🎾','other']`},
	} {
		t.Run(tc.name, func(t *testing.T) {
			data := map[string]any{tc.key: "hit", "other": "second", "ab": "wrong"}
			for _, rooted := range []bool{false, true} {
				x := jp.U(tc.key, "other")
				want := tc.want
				if rooted {
					x = jp.R().U(tc.key, "other")
					want = "$" + want
				}
				tt.Equal(t, want, x.String())
				tt.Equal(t, want, x.BracketString())
				parsed, err := jp.ParseString(x.String())
				tt.Nil(t, err)
				tt.Equal(t, x, parsed)
				tt.Equal(t, x.Get(data), parsed.Get(data))
			}
			tt.Equal(t, "prefix"+tc.want, string(jp.NewUnion(tc.key, "other").Append([]byte("prefix"), false, false)))
		})
	}
}

func TestExprFilter(t *testing.T) {
	f, err := jp.NewFilter("[?(@.x == 3)]")
	tt.Nil(t, err)
	tt.Equal(t, "[?(@.x == 3)]", f.String())

	f, err = jp.NewFilter("[?@.x == 3]")
	tt.Nil(t, err)
	tt.Equal(t, "[?(@.x == 3)]", f.String())

	_, err = jp.NewFilter("[(@.x == 3)]")
	tt.NotNil(t, err)

	_, err = jp.NewFilter("[?(@.x ++ 3)]")
	tt.NotNil(t, err)
}

func TestExprBracket(t *testing.T) {
	br := jp.Bracket('x')
	tt.Equal(t, 0, len(br.Append([]byte{}, true, true)))
}

func TestExprBracketString(t *testing.T) {
	x := jp.R().C("abc").N(1).C("def")
	tt.Equal(t, "$['abc'][1]['def']", x.BracketString())
}
