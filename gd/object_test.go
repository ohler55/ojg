// Copyright (c) 2020, Peter Ohler, All rights reserved.

package gd_test

import (
	"fmt"
	"testing"

	"github.com/ohler55/ojg/gd"
	"github.com/ohler55/ojg/tt"
)

func TestObjectString(t *testing.T) {
	gd.Sort = true
	o := gd.Object{"a": gd.Int(3), "b": gd.Object{"c": gd.Int(5)}, "d": gd.Int(7)}

	tt.Equal(t, `{"a":3,"b":{"c":5},"d":7}`, o.String())
}

func TestObjectNative(t *testing.T) {
	o := gd.Object{"a": gd.Int(3), "b": gd.Int(7)}
	native := o.Native()

	tt.Equal(t, "map[string]interface {}", fmt.Sprintf("%T", native))
	no := native.(map[string]interface{})
	tt.Equal(t, "int64 3  int64 7", fmt.Sprintf("%T %v  %T %v", no["a"], no["a"], no["b"], no["b"]))

}

func TestObjectAlter(t *testing.T) {
	o := gd.Object{"a": gd.Int(3), "b": gd.Int(7)}
	alt := o.Alter()

	tt.Equal(t, "map[string]interface {}", fmt.Sprintf("%T", alt))

	ao := alt.(map[string]interface{})
	tt.Equal(t, "int64 3  int64 7", fmt.Sprintf("%T %v  %T %v", ao["a"], ao["a"], ao["b"], ao["b"]))
}
