// Copyright (c) 2020, Peter Ohler, All rights reserved.

package alt_test

import (
	"fmt"
	"testing"

	"github.com/ohler55/ojg/alt"
	"github.com/ohler55/ojg/tt"
)

func TestBuilderArray(t *testing.T) {
	var b alt.Builder

	err := b.Array()
	tt.Nil(t, err, "b.Array()")
	b.Pop()
	v := b.Result()
	tt.Equal(t, []interface{}{}, v)

	b.Reset()
	tt.Nil(t, b.Result(), "b.Result() after reset")

	err = b.Array()
	tt.Nil(t, err, "first b.Array()")
	err = b.Value(true)
	tt.Nil(t, err, "b.Value(true)")
	err = b.Array()
	tt.Nil(t, err, "second b.Array()")
	err = b.Value(false)
	tt.Nil(t, err, "b.Value(false)")

	b.Pop()
	err = b.Value(nil)
	tt.Nil(t, err, "b.Value(nil)")
	b.PopAll()

	v = b.Result()
	tt.Equal(t, []interface{}{true, []interface{}{false}, nil}, v)
}

func TestBuilderObject(t *testing.T) {
	var b alt.Builder

	err := b.Object()
	tt.Nil(t, err, "b.Object()")
	b.Pop()
	v := b.Result()
	tt.Equal(t, map[string]interface{}{}, v)

	b.Reset()
	tt.Nil(t, b.Result(), "b.Result() after reset")

	err = b.Object()
	tt.Nil(t, err, "first b.Object()")
	err = b.Value(true, "a")
	tt.Nil(t, err, "b.Value(true, a)")

	err = b.Object("b")
	tt.Nil(t, err, "second b.Object()")
	err = b.Value(false, "c")
	tt.Nil(t, err, "b.Value(false, c)")

	b.Pop()
	err = b.Value(nil, "d")
	tt.Nil(t, err, "b.Value(nil, d)")
	b.PopAll()

	v = b.Result()
	tt.Equal(t, map[string]interface{}{"a": true, "b": map[string]interface{}{"c": false}, "d": nil}, v)
}

func TestBuilderMixed(t *testing.T) {
	var b alt.Builder

	b.Reset() // not needed, just making sure there are not issues

	err := b.Object()
	tt.Nil(t, err, "b.Object()")
	err = b.Array("a")
	tt.Nil(t, err, "b.Array(a)")
	err = b.Value(true)
	tt.Nil(t, err, "b.Value(true)")
	err = b.Object()
	tt.Nil(t, err, "b.Object()")
	err = b.Value(123, "x")
	tt.Nil(t, err, "b.Value(123, x)")
	b.Pop()
	err = b.Value(nil)
	tt.Nil(t, err, "b.Value(nil)")
	b.PopAll()

	v := b.Result()
	tt.Equal(t, map[string]interface{}{"a": []interface{}{true, map[string]interface{}{"x": 123}, nil}}, v)
}

func TestBuilderErrors(t *testing.T) {
	var b alt.Builder

	err := b.Object("bad")
	tt.Equal(t, "can not use a key when pushing to an array", err.Error())

	err = b.Array("bad")
	tt.Equal(t, "can not use a key when pushing to an array", err.Error())

	err = b.Value(true, "bad")
	tt.Equal(t, "can not use a key when pushing to an array", err.Error())

	err = b.Object()
	tt.Nil(t, err, "b.Object()")

	err = b.Object()
	tt.Equal(t, "must have a key when pushing to an object", err.Error())

	err = b.Array()
	tt.Equal(t, "must have a key when pushing to an object", err.Error())

	err = b.Value(true)
	tt.Equal(t, "must have a key when pushing to an object", err.Error())
}

func valueFromFuzzingInput(input []byte) interface{} {
	switch {
	case len(input) == 0 || input[0] == 0:
		return nil
	case input[0] == 1:
		return true
	case input[0] == 2:
		return false
	case 'a' <= input[0] && input[0] <= 'z':
		return string(input[0])
	case 'A' <= input[0] && input[0] <= 'Z':
		return string(input[0])
	case input[0] <= 160:
		l := int(input[0] / 16)
		array := make([]interface{}, 0, l)
		for i := 0; i < l; i++ {
			if len(input) > 0 {
				input = input[1:]
			}
			item := valueFromFuzzingInput(input)
			array = append(array, item)
		}
		return array
	default:
		l := int(input[0]/16 - 10)
		obj := make(map[string]interface{})
		for i := 0; i < l; i++ {
			if len(input) > 0 {
				input = input[1:]
			}
			item := valueFromFuzzingInput(input)
			obj[fmt.Sprintf("%d", i)] = item
		}
		return obj
	}
}

func buildFromValue(t *testing.T, b *alt.Builder, value interface{}, key ...string) {
	switch value := value.(type) {
	case map[string]interface{}:
		err := b.Object(key...)
		tt.Nil(t, err, "b.Object()")
		for k, v := range value {
			buildFromValue(t, b, v, k)
		}
		b.Pop()
	case []interface{}:
		err := b.Array(key...)
		for _, v := range value {
			buildFromValue(t, b, v)
		}
		b.Pop()
		tt.Nil(t, err, "b.Array()")
	case string, nil, bool:
		err := b.Value(value, key...)
		tt.Nil(t, err, "b.Value()")
	default:
		t.Fatal(fmt.Sprintf("invalid type: %T", value))
	}
}

// This test can be runned with this command:
//
//     $ go test -fuzz=FuzzBuilder -fuzztime=1s
func FuzzBuilder(f *testing.F) {
	f.Fuzz(func(t *testing.T, input []byte) {
		var b alt.Builder
		value := valueFromFuzzingInput(input)
		buildFromValue(t, &b, value)
		tt.Equal(t, value, b.Result())
	})
}
