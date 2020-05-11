// Copyright (c) 2020, Peter Ohler, All rights reserved.

package simple_test

import (
	"fmt"
	"testing"

	"github.com/ohler55/ojg/simple"
	"github.com/ohler55/ojg/tt"
)

type WithList struct {
	List []int
	Fun  func() bool
}

func sillyRecompose(data map[string]interface{}) (interface{}, error) {
	s := silly{}
	i, _ := data["val"].(int64)
	s.val = int(i)
	return &s, nil
}

func TestSimpleRecomposeBasic(t *testing.T) {
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
	r, err := simple.NewRecomposer("type", map[interface{}]simple.RecomposeFunc{&Dummy{}: nil})
	tt.Nil(t, err, "NewRecomposer")
	var v interface{}
	v, err = r.Recompose(src)
	tt.Nil(t, err, "Recompose")
	d, _ := v.(*Dummy)
	tt.NotNil(t, d, "Dummy")
	tt.Equal(t, []interface{}{-8, -16, -32, 0, 8, 16, 32, 64, 1.2, map[string]interface{}{}}, d.Nest)
}

func TestSimpleRecomposeFunc(t *testing.T) {
	src := map[string]interface{}{"type": "silly", "val": 3}
	r, err := simple.NewRecomposer("type", map[interface{}]simple.RecomposeFunc{&silly{}: sillyRecompose})
	tt.Nil(t, err, "NewRecomposer")
	var v interface{}
	v, err = r.Recompose(src)
	tt.Nil(t, err, "Recompose")
	s, _ := v.(*silly)
	tt.NotNil(t, s, "silly")
	tt.Equal(t, 3, s.val)
}

func TestSimpleRecomposeReflect(t *testing.T) {
	src := map[string]interface{}{"type": "Dummy", "val": 3, "extra": true, "fun": true}
	r, err := simple.NewRecomposer("type", map[interface{}]simple.RecomposeFunc{&Dummy{}: nil})
	tt.Nil(t, err, "NewRecomposer")
	var v interface{}
	v, err = r.Recompose(src)
	tt.Nil(t, err, "Recompose")
	d, _ := v.(*Dummy)
	tt.NotNil(t, d, "check type")
	tt.Equal(t, 3, d.Val)
}

func TestSimpleRecomposeReflectList(t *testing.T) {
	src := map[string]interface{}{"type": "WithList", "list": []interface{}{1, 2, 3}}
	r, err := simple.NewRecomposer("type", map[interface{}]simple.RecomposeFunc{&WithList{}: nil})
	tt.Nil(t, err, "NewRecomposer")
	var v interface{}
	v, err = r.Recompose(src)
	tt.Nil(t, err, "Recompose")
	wl, _ := v.(*WithList)
	tt.NotNil(t, wl, "check type")
	tt.Equal(t, "[]int [1 2 3]", fmt.Sprintf("%T %v", wl.List, wl.List))
}

func TestSimpleRecomposeBadMap(t *testing.T) {
	_, err := simple.NewRecomposer("type", map[interface{}]simple.RecomposeFunc{3: nil})
	tt.NotNil(t, err, "NewRecomposer")
}

func TestSimpleRecomposeBadField(t *testing.T) {
	src := map[string]interface{}{"type": "Dummy", "val": true}
	r, err := simple.NewRecomposer("type", map[interface{}]simple.RecomposeFunc{&Dummy{}: nil})
	tt.Nil(t, err, "NewRecomposer")
	_, err = r.Recompose(src)
	tt.NotNil(t, err, "Recompose")
}

func TestSimpleRecomposeReflectListBad(t *testing.T) {
	src := map[string]interface{}{"type": "WithList", "list": []interface{}{1, true, 3}}
	r, err := simple.NewRecomposer("type", map[interface{}]simple.RecomposeFunc{&WithList{}: nil})
	tt.Nil(t, err, "NewRecomposer")
	_, err = r.Recompose(src)
	tt.NotNil(t, err, "Recompose")
}

func TestSimpleRecomposeBadListItem(t *testing.T) {
	src := map[string]interface{}{
		"type": "Dummy",
		"val":  3,
		"nest": []interface{}{func() {}},
	}
	r, err := simple.NewRecomposer("type", map[interface{}]simple.RecomposeFunc{&Dummy{}: nil})
	tt.Nil(t, err, "NewRecomposer")
	_, err = r.Recompose(src)
	tt.NotNil(t, err, "Recompose")
}

func TestSimpleRecomposeListResult(t *testing.T) {
	src := []interface{}{
		map[string]interface{}{"type": "Dummy", "val": 1},
		map[string]interface{}{"type": "Dummy", "val": 2},
	}
	r, err := simple.NewRecomposer("type", map[interface{}]simple.RecomposeFunc{&Dummy{}: nil})
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

func TestSimpleRecomposeListBadResult(t *testing.T) {
	src := []interface{}{true}
	r, err := simple.NewRecomposer("type", map[interface{}]simple.RecomposeFunc{})
	tt.Nil(t, err, "NewRecomposer")
	_, err = r.Recompose(src, []*Dummy{})
	tt.NotNil(t, err, "Recompose")
}

func TestSimpleRecomposeListBadTarget(t *testing.T) {
	r, err := simple.NewRecomposer("type", map[interface{}]simple.RecomposeFunc{})
	tt.Nil(t, err, "NewRecomposer")
	_, err = r.Recompose("[]", 7)
	tt.NotNil(t, err, "Recompose")
}
