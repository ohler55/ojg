// Copyright (c) 2021, Peter Ohler, All rights reserved.

package asm_test

import (
	"testing"
	"time"

	"github.com/ohler55/ojg/asm"
	"github.com/ohler55/ojg/sen"
	"github.com/ohler55/ojg/tt"
)

func TestEqualNull(t *testing.T) {
	root := testPlan(t,
		`[
           [set $.asm.true [equal null "$.src[3]"]]
           [set $.asm.false [equal null "$.src[1]"]]
         ]`,
		"{src: [1 2 3]}",
	)
	tt.Equal(t, "{false:false true:true}", sen.String(root["asm"], &sopt))
}

func TestEqualBool(t *testing.T) {
	root := testPlan(t,
		`[
           [set $.asm.true [equal true "$.src[0]"]]
           [set $.asm.false [equal true "$.src[1]"]]
         ]`,
		"{src: [true false]}",
	)
	tt.Equal(t, "{false:false true:true}", sen.String(root["asm"], &sopt))
}

func TestEqualString(t *testing.T) {
	root := testPlan(t,
		`[
           [set $.asm.true [equal abc "$.src[0]"]]
           [set $.asm.false [equal abc "$.src[1]"]]
         ]`,
		"{src: [abc xyz]}",
	)
	tt.Equal(t, "{false:false true:true}", sen.String(root["asm"], &sopt))
}

func TestEqualInt(t *testing.T) {
	root := testPlan(t,
		`[
           [set $.asm.true [equal 1 "$.src[0]"]]
           [set $.asm.false [equal 1 1.0 "$.src[1]"]]
         ]`,
		"{src: [1 2 3]}",
	)
	tt.Equal(t, "{false:false true:true}", sen.String(root["asm"], &sopt))
}

func TestEqualFloat(t *testing.T) {
	root := testPlan(t,
		`[
           [set $.asm.a [equal 1.1 "$.src[0]"]]
           [set $.asm.b [equal 1.1 "$.src[1]"]]
           [set $.asm.c [equal 1 "$.src[2]"]]
           [set $.asm.d [equal 1 "$.src[0]"]]
         ]`,
		"{src: [1.1 2.2 1.0 2.0]}",
	)
	tt.Equal(t, "{a:true b:false c:true d:false}", sen.String(root["asm"], &sopt))
}

func TestEqualArray(t *testing.T) {
	root := testPlan(t,
		`[
           [set $.asm.a [equal [1 2 3] $.src]]
           [set $.asm.b [equal [1 2] $.src]]
           [set $.asm.c [equal [1 2 3 4] $.src]]
           [set $.asm.d [equal [1 2 4] $.src]]
         ]`,
		"{src: [1 2 3]}",
	)
	tt.Equal(t, "{a:true b:false c:false d:false}", sen.String(root["asm"], &sopt))
}

func TestEqualMap(t *testing.T) {
	root := testPlan(t,
		`[
           [set $.asm.a [eq {b:2 a:1} $.src]]
           [set $.asm.b [equal {a:1} $.src]]
           [set $.asm.c ["==" {a:1 b:2 c:3} $.src]]
           [set $.asm.d [equal {a:1 b:3} $.src]]
         ]`,
		"{src: {a:1 b:2}}",
	)
	tt.Equal(t, "{a:true b:false c:false d:false}", sen.String(root["asm"], &sopt))
}

func TestEqualTime(t *testing.T) {
	t1 := time.Now().UTC()
	t2 := t1.Add(time.Second)
	p := asm.NewPlan([]interface{}{
		[]interface{}{"set", "$.asm.true", []interface{}{"eq", t1, t1}},
		[]interface{}{"set", "$.asm.false", []interface{}{"eq", t1, t2}},
	})
	root := map[string]interface{}{
		"src": []interface{}{},
	}
	err := p.Execute(root)
	tt.Nil(t, err)

	tt.Equal(t, "{false:false true:true}", sen.String(root["asm"], &sopt))
}

func TestEqualIntOthers(t *testing.T) {
	p := asm.NewPlan([]interface{}{
		[]interface{}{"set", "$.asm.i", []interface{}{"eq", 1, int8(1), int16(1), int32(1), int64(1)}},
		[]interface{}{"set", "$.asm.u", []interface{}{"eq", uint(1), uint8(1), uint16(1), uint32(1), uint64(1)}},
	})
	root := map[string]interface{}{
		"src": []interface{}{},
	}
	err := p.Execute(root)
	tt.Nil(t, err)

	tt.Equal(t, "{i:true u:true}", sen.String(root["asm"], &sopt))
}

func TestEqualFloatOthers(t *testing.T) {
	p := asm.NewPlan([]interface{}{
		// Float32 almost never matches float64 due to rounding errors
		[]interface{}{"set", "$.asm.float", []interface{}{"eq", float32(1.1), float32(1.1)}},
	})
	root := map[string]interface{}{
		"src": []interface{}{},
	}
	err := p.Execute(root)
	tt.Nil(t, err)

	tt.Equal(t, "{float:true}", sen.String(root["asm"], &sopt))
}
