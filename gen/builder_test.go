// Copyright (c) 2020, Peter Ohler, All rights reserved.

package gen_test

import (
	"testing"

	"github.com/ohler55/ojg/gen"
	"github.com/ohler55/ojg/tt"
)

func TestGenBuilderArray(t *testing.T) {
	var b gen.Builder

	err := b.Array()
	tt.Nil(t, err, "b.Array()")
	b.Pop()
	v := b.Result()
	tt.Equal(t, gen.Array{}, v)

	b.Reset()
	tt.Nil(t, b.Result(), "b.Result() after reset")

	err = b.Array()
	tt.Nil(t, err, "first b.Array()")
	err = b.Value(gen.True)
	tt.Nil(t, err, "b.Value(true)")
	err = b.Array()
	tt.Nil(t, err, "second b.Array()")
	err = b.Value(gen.False)
	tt.Nil(t, err, "b.Value(false)")
	b.PopAll()

	v = b.Result()
	tt.Equal(t, gen.Array{gen.True, gen.Array{gen.False}}, v)
}

func TestGenBuilderObject(t *testing.T) {
	var b gen.Builder

	err := b.Object()
	tt.Nil(t, err, "b.Object()")
	b.Pop()
	v := b.Result()
	tt.Equal(t, map[string]interface{}{}, v)

	b.Reset()
	tt.Nil(t, b.Result(), "b.Result() after reset")

	err = b.Object()
	tt.Nil(t, err, "first b.Object()")
	err = b.Value(gen.True, "a")
	tt.Nil(t, err, "b.Value(true, a)")

	err = b.Object("b")
	tt.Nil(t, err, "second b.Object()")
	err = b.Value(gen.False, "c")
	tt.Nil(t, err, "b.Value(false, c)")
	b.PopAll()

	v = b.Result()
	tt.Equal(t, gen.Object{"a": gen.True, "b": gen.Object{"c": gen.False}}, v)
}

func TestGenBuilderMixed(t *testing.T) {
	var b gen.Builder

	b.Reset() // not needed, just making sure there are not issues

	err := b.Object()
	tt.Nil(t, err, "b.Object()")
	err = b.Array("a")
	tt.Nil(t, err, "b.Array(a)")
	err = b.Value(gen.True)
	tt.Nil(t, err, "b.Value(true)")
	err = b.Object()
	tt.Nil(t, err, "b.Object()")
	err = b.Value(gen.Int(123), "x")
	tt.Nil(t, err, "b.Value(123, x)")
	b.Pop()
	err = b.Value(nil)
	tt.Nil(t, err, "b.Value(nil)")
	b.PopAll()

	v := b.Result()
	tt.Equal(t, gen.Object{"a": gen.Array{gen.True, gen.Object{"x": gen.Int(123)}, nil}}, v)
}

func TestGenBuilderErrors(t *testing.T) {
	var b gen.Builder

	err := b.Object("bad")
	tt.Equal(t, "can not use a key when pushing to an array", err.Error())

	err = b.Array("bad")
	tt.Equal(t, "can not use a key when pushing to an array", err.Error())

	err = b.Value(gen.True, "bad")
	tt.Equal(t, "can not use a key when pushing to an array", err.Error())

	err = b.Object()
	tt.Nil(t, err, "b.Object()")

	err = b.Object()
	tt.Equal(t, "must have a key when pushing to an object", err.Error())

	err = b.Array()
	tt.Equal(t, "must have a key when pushing to an object", err.Error())

	err = b.Value(gen.True)
	tt.Equal(t, "must have a key when pushing to an object", err.Error())
}

func TestGenBuilderPanic(t *testing.T) {
	var b gen.Builder

	tt.Panic(t, func() { b.MustObject("bad") })
	tt.Panic(t, func() { b.MustArray("bad") })
	tt.Panic(t, func() { b.MustValue(gen.True, "bad") })
	b.MustObject()
	tt.Panic(t, func() { b.MustObject() })
	tt.Panic(t, func() { b.MustArray() })
	tt.Panic(t, func() { b.MustValue(gen.True) })
}
