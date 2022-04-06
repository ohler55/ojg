// Copyright (c) 2020, Peter Ohler, All rights reserved.

package jp_test

import (
	"testing"

	"github.com/ohler55/ojg/jp"
	"github.com/ohler55/ojg/tt"
)

func TestEquation(t *testing.T) {
	eq := jp.Neq(jp.ConstInt(3), jp.ConstFloat(1.5))
	tt.Equal(t, "(3 != 1.5)", eq.String())

	eq = jp.Eq(jp.ConstBool(true), jp.ConstNil())
	tt.Equal(t, "(true == null)", eq.String())

	eq = jp.Get(jp.A().C("xyz"))
	tt.Equal(t, "(@.xyz)", eq.String())

	eq = jp.Lt(jp.ConstInt(3), jp.ConstInt(4))
	tt.Equal(t, "(3 < 4)", eq.String())

	eq = jp.Lte(jp.ConstInt(3), jp.ConstInt(4))
	tt.Equal(t, "(3 <= 4)", eq.String())

	eq = jp.Gt(jp.ConstInt(3), jp.ConstInt(4))
	tt.Equal(t, "(3 > 4)", eq.String())

	eq = jp.Gte(jp.ConstInt(3), jp.ConstInt(4))
	tt.Equal(t, "(3 >= 4)", eq.String())

	eq = jp.Or(jp.ConstBool(true), jp.ConstBool(false))
	tt.Equal(t, "(true || false)", eq.String())

	eq = jp.And(jp.ConstBool(true), jp.ConstBool(false))
	tt.Equal(t, "(true && false)", eq.String())

	eq = jp.Not(jp.ConstBool(true))
	tt.Equal(t, "(!true)", eq.String())

	eq = jp.Add(jp.ConstInt(3), jp.ConstInt(4))
	tt.Equal(t, "(3 + 4)", eq.String())

	eq = jp.Sub(jp.ConstInt(3), jp.ConstInt(4))
	tt.Equal(t, "(3 - 4)", eq.String())

	eq = jp.Multiply(jp.ConstInt(3), jp.ConstInt(4))
	tt.Equal(t, "(3 * 4)", eq.String())

	eq = jp.Divide(jp.ConstInt(3), jp.ConstInt(4))
	tt.Equal(t, "(3 / 4)", eq.String())

	eq = jp.In(jp.ConstInt(3), jp.ConstList([]interface{}{int64(1), int64(2), int64(3)}))
	tt.Equal(t, "(3 in [1,2,3])", eq.String())
}

func TestEquationScript(t *testing.T) {
	eq := jp.And(nil, nil)
	tt.Equal(t, "(null && null)", eq.Script().String())

	eq = jp.Not(nil)
	tt.Equal(t, "(!null)", eq.Script().String())
}
