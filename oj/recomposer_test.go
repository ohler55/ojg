// Copyright (c) 2020, Peter Ohler, All rights reserved.

package oj_test

import (
	"fmt"
	"testing"

	"github.com/ohler55/ojg/oj"
	"github.com/ohler55/ojg/tt"
)

type WithList struct {
	List []int
	Fun  func() bool
}

type Setter struct {
	a int64
	b string
}

func (s *Setter) String() string {
	return fmt.Sprintf("Setter{a:%d,b:%s}", s.a, s.b)
}

func (s *Setter) SetAttr(attr string, val interface{}) error {
	switch attr {
	case "a":
		s.a, _ = val.(int64)
	case "b":
		s.b, _ = val.(string)
	default:
		return fmt.Errorf("%s is not an attribute of Setter", attr)
	}
	return nil
}

func sillyRecompose(data map[string]interface{}) (interface{}, error) {
	s := silly{}
	i, _ := data["val"].(int64)
	s.val = int(i)
	return &s, nil
}

func TestOjRecomposeBasic(t *testing.T) {
	src := map[string]interface{}{
		"type": "Dummy",
		"val":  3,
		"nest": []interface{}{
			int8(-8), int16(-16), int32(-32),
			uint(0), uint8(8), uint16(16), uint32(32), uint64(64),
			float32(1.2),
			map[string]interface{}{},
		},
	}
	r, err := oj.NewRecomposer("type", map[interface{}]oj.RecomposeFunc{&Dummy{}: nil})
	tt.Nil(t, err, "NewRecomposer")
	var v interface{}
	v, err = r.Recompose(src)
	tt.Nil(t, err, "Recompose")
	d, _ := v.(*Dummy)
	tt.NotNil(t, d, "Dummy")
	tt.Equal(t, []interface{}{-8, -16, -32, 0, 8, 16, 32, 64, 1.2, map[string]interface{}{}}, d.Nest)
}

func TestOjRecomposeFunc(t *testing.T) {
	src := map[string]interface{}{"type": "silly", "val": 3}
	r, err := oj.NewRecomposer("type", map[interface{}]oj.RecomposeFunc{&silly{}: sillyRecompose})
	tt.Nil(t, err, "NewRecomposer")
	var v interface{}
	v, err = r.Recompose(src)
	tt.Nil(t, err, "Recompose")
	s, _ := v.(*silly)
	tt.NotNil(t, s, "silly")
	tt.Equal(t, 3, s.val)
}

func TestOjRecomposeReflect(t *testing.T) {
	src := map[string]interface{}{"type": "Dummy", "val": 3, "extra": true, "fun": true}
	r, err := oj.NewRecomposer("type", map[interface{}]oj.RecomposeFunc{&Dummy{}: nil})
	tt.Nil(t, err, "NewRecomposer")
	var v interface{}
	v, err = r.Recompose(src)
	tt.Nil(t, err, "Recompose")
	d, _ := v.(*Dummy)
	tt.NotNil(t, d, "check type")
	tt.Equal(t, 3, d.Val)
}

func TestOjRecomposeAttrSetter(t *testing.T) {
	src := map[string]interface{}{"type": "Setter", "a": 3, "b": "bee"}
	r, err := oj.NewRecomposer("type", map[interface{}]oj.RecomposeFunc{&Setter{}: nil})
	tt.Nil(t, err, "NewRecomposer")
	var v interface{}
	v, err = r.Recompose(src)
	tt.Nil(t, err, "Recompose")
	s, _ := v.(*Setter)
	tt.NotNil(t, s, "check type")
	tt.Equal(t, "Setter{a:3,b:bee}", s.String())
}

func TestOjRecomposeReflectList(t *testing.T) {
	src := map[string]interface{}{"type": "WithList", "list": []interface{}{1, 2, 3}}
	r, err := oj.NewRecomposer("type", map[interface{}]oj.RecomposeFunc{&WithList{}: nil})
	tt.Nil(t, err, "NewRecomposer")
	var v interface{}
	v, err = r.Recompose(src)
	tt.Nil(t, err, "Recompose")
	wl, _ := v.(*WithList)
	tt.NotNil(t, wl, "check type")
	tt.Equal(t, "[]int [1 2 3]", fmt.Sprintf("%T %v", wl.List, wl.List))
}

func TestOjRecomposeBadMap(t *testing.T) {
	_, err := oj.NewRecomposer("type", map[interface{}]oj.RecomposeFunc{3: nil})
	tt.NotNil(t, err, "NewRecomposer")
}

func TestOjRecomposeBadField(t *testing.T) {
	src := map[string]interface{}{"type": "Dummy", "val": true}
	r, err := oj.NewRecomposer("type", map[interface{}]oj.RecomposeFunc{&Dummy{}: nil})
	tt.Nil(t, err, "NewRecomposer")
	_, err = r.Recompose(src)
	tt.NotNil(t, err, "Recompose")
}

func TestOjRecomposeReflectListBad(t *testing.T) {
	src := map[string]interface{}{"type": "WithList", "list": []interface{}{1, true, 3}}
	r, err := oj.NewRecomposer("type", map[interface{}]oj.RecomposeFunc{&WithList{}: nil})
	tt.Nil(t, err, "NewRecomposer")
	_, err = r.Recompose(src)
	tt.NotNil(t, err, "Recompose")
}

func TestOjRecomposeBadListItem(t *testing.T) {
	src := map[string]interface{}{
		"type": "Dummy",
		"val":  3,
		"nest": []interface{}{func() {}},
	}
	r, err := oj.NewRecomposer("type", map[interface{}]oj.RecomposeFunc{&Dummy{}: nil})
	tt.Nil(t, err, "NewRecomposer")
	_, err = r.Recompose(src)
	tt.NotNil(t, err, "Recompose")
}

func TestOjRecomposeListResult(t *testing.T) {
	src := []interface{}{
		map[string]interface{}{"type": "Dummy", "val": 1},
		map[string]interface{}{"type": "Dummy", "val": 2},
	}
	r, err := oj.NewRecomposer("type", map[interface{}]oj.RecomposeFunc{&Dummy{}: nil})
	tt.Nil(t, err, "NewRecomposer")
	var v interface{}
	v, err = r.Recompose(src, []*Dummy{})
	tt.Nil(t, err, "Recompose")
	da, _ := v.([]*Dummy)
	tt.NotNil(t, da, "check type")
	tt.Equal(t, 2, len(da))
	for i, d := range da {
		tt.Equal(t, i+1, d.Val)
	}
}

func TestOjRecomposeListBadResult(t *testing.T) {
	src := []interface{}{true}
	r, err := oj.NewRecomposer("type", map[interface{}]oj.RecomposeFunc{})
	tt.Nil(t, err, "NewRecomposer")
	_, err = r.Recompose(src, []*Dummy{})
	tt.NotNil(t, err, "Recompose")
}

func TestOjRecomposeListBadTarget(t *testing.T) {
	r, err := oj.NewRecomposer("type", map[interface{}]oj.RecomposeFunc{})
	tt.Nil(t, err, "NewRecomposer")
	_, err = r.Recompose("[]", 7)
	tt.NotNil(t, err, "Recompose")
}
