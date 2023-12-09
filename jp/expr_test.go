// Copyright (c) 2020, Peter Ohler, All rights reserved.

package jp_test

import (
	"testing"

	"github.com/ohler55/ojg/jp"
	"github.com/ohler55/ojg/tt"
)

func TestExprBuild(t *testing.T) {
	x := jp.X().D().C("abc").W().N(3).U(2, "x").S(1, 5, 2).S(1, 5).S(1)
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

	x = jp.D()
	tt.Equal(t, "..", x.String())

	x = jp.U(1, "a")
	tt.Equal(t, "[1,'a']", x.String())

	x = jp.Expr{jp.Slice{}}
	tt.Equal(t, "[:]", x.String())

	x = jp.R().Child("'")
	tt.Equal(t, `$['\'']`, x.String())

	x = jp.R().Child("").Child("a")
	tt.Equal(t, `$[''].a`, x.String())
}

func TestExprFilter(t *testing.T) {
	f, err := jp.NewFilter("[?(@.x == 3)]")
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
