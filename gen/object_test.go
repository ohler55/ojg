// Copyright (c) 2020, Peter Ohler, All rights reserved.

package gen_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/ohler55/ojg/gen"
	"github.com/ohler55/ojg/tt"
)

func TestObjectString(t *testing.T) {
	gen.Sort = true
	o := gen.Object{"a": gen.Int(3), "b": gen.Object{"c": gen.Int(5)}, "d": gen.Int(7), "n": nil}
	tt.Equal(t, `{"a":3,"b":{"c":5},"d":7,"n":null}`, o.String())

	gen.Sort = false
	o = gen.Object{"a": nil, "b": gen.Int(7)}
	tt.Equal(t, true, strings.Contains(o.String(), "null"))
	tt.Equal(t, true, strings.Contains(o.String(), "7"))
}

func TestObjectSimplify(t *testing.T) {
	o := gen.Object{"a": gen.Int(3), "b": gen.Int(7), "n": nil}
	simple := o.Simplify()

	tt.Equal(t, "map[string]interface {}", fmt.Sprintf("%T", simple))
	no := simple.(map[string]interface{})
	tt.Equal(t, "int64 3  int64 7", fmt.Sprintf("%T %v  %T %v", no["a"], no["a"], no["b"], no["b"]))
}

func TestObjectAlter(t *testing.T) {
	o := gen.Object{"a": gen.Int(3), "b": gen.Int(7), "n": nil}
	alt := o.Alter()

	tt.Equal(t, "map[string]interface {}", fmt.Sprintf("%T", alt))

	ao := alt.(map[string]interface{})
	tt.Equal(t, "int64 3  int64 7", fmt.Sprintf("%T %v  %T %v", ao["a"], ao["a"], ao["b"], ao["b"]))
}

func TestObjectDup(t *testing.T) {
	gen.Sort = true
	o := gen.Object{"a": gen.Int(3), "b": gen.Int(7), "n": nil}

	dup := o.Dup()
	tt.NotNil(t, dup)
	tt.Equal(t, `{"a":3,"b":7,"n":null}`, dup.String())
}

func TestObjectEmpty(t *testing.T) {
	tt.Equal(t, true, gen.Object{}.Empty())
}
