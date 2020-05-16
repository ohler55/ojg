// Copyright (c) 2020, Peter Ohler, All rights reserved.

package oj_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/ohler55/ojg/oj"
	"github.com/ohler55/ojg/tt"
)

func TestObjectString(t *testing.T) {
	oj.Sort = true
	o := oj.Object{"a": oj.Int(3), "b": oj.Object{"c": oj.Int(5)}, "d": oj.Int(7)}
	tt.Equal(t, `{"a":3,"b":{"c":5},"d":7}`, o.String())

	oj.Sort = false
	o = oj.Object{"a": nil, "b": oj.Int(7)}
	tt.Equal(t, true, strings.Contains(o.String(), "null"))
	tt.Equal(t, true, strings.Contains(o.String(), "7"))
}

func TestObjectSimplify(t *testing.T) {
	o := oj.Object{"a": oj.Int(3), "b": oj.Int(7)}
	simple := o.Simplify()

	tt.Equal(t, "map[string]interface {}", fmt.Sprintf("%T", simple))
	no := simple.(map[string]interface{})
	tt.Equal(t, "int64 3  int64 7", fmt.Sprintf("%T %v  %T %v", no["a"], no["a"], no["b"], no["b"]))
}

func TestObjectAlter(t *testing.T) {
	o := oj.Object{"a": oj.Int(3), "b": oj.Int(7)}
	alt := o.Alter()

	tt.Equal(t, "map[string]interface {}", fmt.Sprintf("%T", alt))

	ao := alt.(map[string]interface{})
	tt.Equal(t, "int64 3  int64 7", fmt.Sprintf("%T %v  %T %v", ao["a"], ao["a"], ao["b"], ao["b"]))
}

func TestObjectDup(t *testing.T) {
	oj.Sort = true
	o := oj.Object{"a": oj.Int(3), "b": oj.Int(7)}

	dup := o.Dup()
	tt.NotNil(t, dup)
	tt.Equal(t, `{"a":3,"b":7}`, dup.String())
}

func TestObjectEmpty(t *testing.T) {
	tt.Equal(t, true, oj.Object{}.Empty())
}
