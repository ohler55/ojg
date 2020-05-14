// Copyright (c) 2020, Peter Ohler, All rights reserved.

package oj_test

import (
	"testing"

	"github.com/ohler55/ojg/oj"
	"github.com/ohler55/ojg/tt"
)

func TestGenBuilderArray(t *testing.T) {
	var b oj.NodeBuilder

	err := b.Array()
	tt.Nil(t, err, "b.Array()")
	b.Pop()
	v := b.Result()
	tt.Equal(t, oj.Array{}, v)

	b.Reset()
	tt.Nil(t, b.Result(), "b.Result() after reset")

	err = b.Array()
	tt.Nil(t, err, "first b.Array()")
	err = b.Value(oj.True)
	tt.Nil(t, err, "b.Value(true)")
	err = b.Array()
	tt.Nil(t, err, "second b.Array()")
	err = b.Value(oj.False)
	tt.Nil(t, err, "b.Value(false)")
	b.PopAll()

	v = b.Result()
	tt.Equal(t, oj.Array{oj.True, oj.Array{oj.False}}, v)
}

func TestGenBuilderObject(t *testing.T) {
	var b oj.NodeBuilder

	err := b.Object()
	tt.Nil(t, err, "b.Object()")
	b.Pop()
	v := b.Result()
	tt.Equal(t, map[string]interface{}{}, v)

	b.Reset()
	tt.Nil(t, b.Result(), "b.Result() after reset")

	err = b.Object()
	tt.Nil(t, err, "first b.Object()")
	err = b.Value(oj.True, "a")
	tt.Nil(t, err, "b.Value(true, a)")

	err = b.Object("b")
	tt.Nil(t, err, "second b.Object()")
	err = b.Value(oj.False, "c")
	tt.Nil(t, err, "b.Value(false, c)")
	b.PopAll()

	v = b.Result()
	tt.Equal(t, oj.Object{"a": oj.True, "b": oj.Object{"c": oj.False}}, v)
}

func TestGenBuilderMixed(t *testing.T) {
	var b oj.NodeBuilder

	b.Reset() // not needed, just making sure there are not issues

	err := b.Object()
	tt.Nil(t, err, "b.Object()")
	err = b.Array("a")
	tt.Nil(t, err, "b.Array(a)")
	err = b.Value(oj.True)
	tt.Nil(t, err, "b.Value(true)")
	err = b.Object()
	tt.Nil(t, err, "b.Object()")
	err = b.Value(oj.Int(123), "x")
	tt.Nil(t, err, "b.Value(123, x)")
	b.Pop()
	err = b.Value(nil)
	tt.Nil(t, err, "b.Value(nil)")
	b.PopAll()

	v := b.Result()
	tt.Equal(t, oj.Object{"a": oj.Array{oj.True, oj.Object{"x": oj.Int(123)}, nil}}, v)
}

func TestGenBuilderErrors(t *testing.T) {
	var b oj.NodeBuilder

	err := b.Object("bad")
	tt.Equal(t, "can not use a key when pushing to an array", err.Error())

	err = b.Array("bad")
	tt.Equal(t, "can not use a key when pushing to an array", err.Error())

	err = b.Value(oj.True, "bad")
	tt.Equal(t, "can not use a key when pushing to an array", err.Error())

	err = b.Object()
	tt.Nil(t, err, "b.Object()")

	err = b.Object()
	tt.Equal(t, "must have a key when pushing to an object", err.Error())

	err = b.Array()
	tt.Equal(t, "must have a key when pushing to an object", err.Error())

	err = b.Value(oj.True)
	tt.Equal(t, "must have a key when pushing to an object", err.Error())
}
