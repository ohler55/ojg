// Copyright (c) 2020, Peter Ohler, All rights reserved.

package jp

import (
	"strconv"
)

// Equation represents JSON Path script and filter equations. They are used to
// build a script. The purpose of the Equation is to allow scripts or filters
// to be created without using a parser which could return an error if an
// invalid string representation of the script is provided.
type Equation struct {
	o      *op
	result interface{}
	left   *Equation
	right  *Equation
}

// Script creates and returns a Script that implements the equation.
func (e *Equation) Script() (s *Script) {
	s = &Script{template: e.buildScript([]interface{}{})}
	s.stack = make([]interface{}, len(s.template))
	return
}

// Filter creates and returns a Script that implements the equation.
func (e *Equation) Filter() (f *Filter) {
	f = &Filter{Script: Script{template: e.buildScript([]interface{}{})}}
	f.stack = make([]interface{}, len(f.template))
	return
}

// ConstNil creates and returns an Equation for a constant of nil.
func ConstNil() *Equation {
	return &Equation{result: nil}
}

// ConstBool creates and returns an Equation for a bool constant.
func ConstBool(b bool) *Equation {
	return &Equation{result: b}
}

// ConstInt creates and returns an Equation for an int64 constant.
func ConstInt(i int64) *Equation {
	return &Equation{result: i}
}

// ConstFloat creates and returns an Equation for a float64 constant.
func ConstFloat(f float64) *Equation {
	return &Equation{result: f}
}

// ConstString creates and returns an Equation for a string constant.
func ConstString(s string) *Equation {
	return &Equation{result: s}
}

// Get creates and returns an Equation for an expression get of the form
// @.child.
func Get(x Expr) *Equation {
	return &Equation{o: get, left: &Equation{result: x}}
}

// Eq creates and returns an Equation for an == operator.
func Eq(left, right *Equation) *Equation {
	return &Equation{o: eq, left: left, right: right}
}

// Neq creates and returns an Equation for a != operator.
func Neq(left, right *Equation) *Equation {
	return &Equation{o: neq, left: left, right: right}
}

// Lt creates and returns an Equation for a < operator.
func Lt(left, right *Equation) *Equation {
	return &Equation{o: lt, left: left, right: right}
}

// Gt creates and returns an Equation for a > operator.
func Gt(left, right *Equation) *Equation {
	return &Equation{o: gt, left: left, right: right}
}

// Lte creates and returns an Equation for a <= operator.
func Lte(left, right *Equation) *Equation {
	return &Equation{o: lte, left: left, right: right}
}

// Gte creates and returns an Equation for a >= operator.
func Gte(left, right *Equation) *Equation {
	return &Equation{o: gte, left: left, right: right}
}

// Or creates and returns an Equation for a || operator.
func Or(left, right *Equation) *Equation {
	return &Equation{o: or, left: left, right: right}
}

// And creates and returns an Equation for a && operator.
func And(left, right *Equation) *Equation {
	return &Equation{o: and, left: left, right: right}
}

// Not creates and returns an Equation for a ! operator.
func Not(arg *Equation) *Equation {
	return &Equation{o: not, left: arg}
}

// Add creates and returns an Equation for a + operator.
func Add(left, right *Equation) *Equation {
	return &Equation{o: add, left: left, right: right}
}

// Sub creates and returns an Equation for a - operator.
func Sub(left, right *Equation) *Equation {
	return &Equation{o: sub, left: left, right: right}
}

// Multiply creates and returns an Equation for a * operator.
func Multiply(left, right *Equation) *Equation {
	return &Equation{o: mult, left: left, right: right}
}

// Divide creates and returns an Equation for a / operator.
func Divide(left, right *Equation) *Equation {
	return &Equation{o: divide, left: left, right: right}
}

// Append a fragment string representation of the fragment to the buffer
// then returning the expanded buffer.
func (s *Equation) Append(buf []byte, parens bool) []byte {
	if parens {
		buf = append(buf, '(')
	}
	if s.o == nil {
		buf = s.appendValue(buf, s.result)
	} else {
		switch s.o.code {
		case not.code:
			buf = append(buf, '!')
			if s.left != nil {
				buf = s.left.Append(buf, s.left.o != nil && s.left.o.prec >= s.o.prec)
			}
		case get.code:
			if s.left != nil {
				buf = s.appendValue(buf, s.left.result)
			}
		default:
			if s.left != nil {
				buf = s.left.Append(buf, s.left.o != nil && s.left.o.prec >= s.o.prec)
			}
			buf = append(buf, ' ')
			buf = append(buf, s.o.name...)
			buf = append(buf, ' ')
			if s.right != nil {
				buf = s.right.Append(buf, s.left.o != nil && s.left.o.prec >= s.o.prec)
			}
		}
	}
	if parens {
		buf = append(buf, ')')
	}
	return buf
}

func (s *Equation) appendValue(buf []byte, v interface{}) []byte {
	switch tv := v.(type) {
	case nil:
		buf = append(buf, "null"...)
	case string:
		buf = append(buf, '\'')
		buf = append(buf, tv...)
		buf = append(buf, '\'')
	case int64:
		buf = append(buf, strconv.FormatInt(tv, 10)...)
	case float64:
		buf = append(buf, strconv.FormatFloat(tv, 'g', -1, 64)...)
	case bool:
		if tv {
			buf = append(buf, "true"...)
		} else {
			buf = append(buf, "false"...)
		}
	case Expr:
		buf = tv.Append(buf)
	}
	return buf
}

// String representation of the equation.
func (s *Equation) String() string {
	return string(s.Append([]byte{}, true))
}

func (e *Equation) buildScript(stack []interface{}) []interface{} {
	if e.o == nil {
		stack = append(stack, e.result)
		return stack
	}
	switch e.o.code {
	case get.code:
		if e.left != nil {
			stack = append(stack, e.left.result) // should always be an Expr
		}
	case not.code:
		stack = append(stack, e.o)
		if e.left == nil {
			stack = append(stack, nil)
		} else {
			stack = e.left.buildScript(stack)
		}
	default:
		stack = append(stack, e.o)
		if e.left == nil {
			stack = append(stack, nil)
		} else {
			stack = e.left.buildScript(stack)
		}
		if e.right == nil {
			stack = append(stack, nil)
		} else {
			stack = e.right.buildScript(stack)
		}
	}
	return stack
}
