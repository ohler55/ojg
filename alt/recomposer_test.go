// Copyright (c) 2020, Peter Ohler, All rights reserved.

package alt_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/ohler55/ojg/alt"
	"github.com/ohler55/ojg/gen"
	"github.com/ohler55/ojg/tt"
)

type WithList struct {
	List []int
	Fun  func() bool
}

type Setter struct {
	a int64
	b string
	//s *Setter
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

func TestRecomposeBasic(t *testing.T) {
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
	r, err := alt.NewRecomposer("type", map[interface{}]alt.RecomposeFunc{&Dummy{}: nil})
	tt.Nil(t, err, "NewRecomposer")
	var v interface{}
	v, err = r.Recompose(src)
	tt.Nil(t, err, "Recompose")
	d, _ := v.(*Dummy)
	tt.NotNil(t, d, "Dummy")
	tt.Equal(t, []interface{}{-8, -16, -32, 0, 8, 16, 32, 64, 1.2, map[string]interface{}{}}, d.Nest)
}

func TestRecomposeNode(t *testing.T) {
	tm := time.Date(2020, time.April, 12, 16, 34, 04, 123456789, time.UTC)
	src := map[string]interface{}{
		"type": "Dummy",
		"val":  gen.Int(3),
		"nest": gen.Array{gen.Int(-8), gen.Bool(true), gen.Float(1.2), gen.String("abc"),
			gen.Object{"big": gen.Big("123"), "time": gen.Time(tm)},
		},
	}
	r, err := alt.NewRecomposer("type", map[interface{}]alt.RecomposeFunc{&Dummy{}: nil})
	tt.Nil(t, err, "NewRecomposer")
	var v interface{}
	v, err = r.Recompose(src)
	tt.Nil(t, err, "Recompose")
	d, _ := v.(*Dummy)
	tt.NotNil(t, d, "Dummy")
	tt.Equal(t, []interface{}{-8, true, 1.2, "abc", map[string]interface{}{"big": "123", "time": tm}}, d.Nest)
}

func TestRecomposeFunc(t *testing.T) {
	src := map[string]interface{}{"type": "silly", "val": 3}
	r, err := alt.NewRecomposer("type", map[interface{}]alt.RecomposeFunc{&silly{}: sillyRecompose})
	tt.Nil(t, err, "NewRecomposer")
	var v interface{}
	v, err = r.Recompose(src)
	tt.Nil(t, err, "Recompose")
	s, _ := v.(*silly)
	tt.NotNil(t, s, "silly")
	tt.Equal(t, 3, s.val)
}

func TestRecomposeReflect(t *testing.T) {
	src := map[string]interface{}{"type": "Dummy", "val": 3, "extra": true, "fun": true}
	r, err := alt.NewRecomposer("type", map[interface{}]alt.RecomposeFunc{&Dummy{}: nil})
	tt.Nil(t, err, "NewRecomposer")
	var v interface{}
	v, err = r.Recompose(src)
	tt.Nil(t, err, "Recompose")
	d, _ := v.(*Dummy)
	tt.NotNil(t, d, "check type")
	tt.Equal(t, 3, d.Val)
}

func TestRecomposeAttrSetter(t *testing.T) {
	src := map[string]interface{}{"type": "Setter", "a": 3, "b": "bee"}
	r, err := alt.NewRecomposer("type", map[interface{}]alt.RecomposeFunc{&Setter{}: nil})
	tt.Nil(t, err, "NewRecomposer")
	var v interface{}
	v, err = r.Recompose(src)
	tt.Nil(t, err, "Recompose")
	s, _ := v.(*Setter)
	tt.NotNil(t, s, "check type")
	tt.Equal(t, "Setter{a:3,b:bee}", s.String())
}

func TestRecomposeReflectList(t *testing.T) {
	src := map[string]interface{}{"type": "WithList", "list": []interface{}{1, 2, 3}}
	r, err := alt.NewRecomposer("type", map[interface{}]alt.RecomposeFunc{&WithList{}: nil})
	tt.Nil(t, err, "NewRecomposer")
	var v interface{}
	v, err = r.Recompose(src)
	tt.Nil(t, err, "Recompose")
	wl, _ := v.(*WithList)
	tt.NotNil(t, wl, "check type")
	tt.Equal(t, "[]int [1 2 3]", fmt.Sprintf("%T %v", wl.List, wl.List))
}

func TestRecomposeBadMap(t *testing.T) {
	_, err := alt.NewRecomposer("type", map[interface{}]alt.RecomposeFunc{3: nil})
	tt.NotNil(t, err, "NewRecomposer")
}

func TestRecomposeBadField(t *testing.T) {
	src := map[string]interface{}{"type": "Dummy", "val": true}
	r, err := alt.NewRecomposer("type", map[interface{}]alt.RecomposeFunc{&Dummy{}: nil})
	tt.Nil(t, err, "NewRecomposer")
	_, err = r.Recompose(src)
	tt.NotNil(t, err, "Recompose")
}

func TestRecomposeReflectListBad(t *testing.T) {
	src := map[string]interface{}{"type": "WithList", "list": []interface{}{1, true, 3}}
	r, err := alt.NewRecomposer("type", map[interface{}]alt.RecomposeFunc{&WithList{}: nil})
	tt.Nil(t, err, "NewRecomposer")
	_, err = r.Recompose(src)
	tt.NotNil(t, err, "Recompose")
}

func TestRecomposeBadListItem(t *testing.T) {
	src := map[string]interface{}{
		"type": "Dummy",
		"val":  3,
		"nest": []interface{}{func() {}},
	}
	r, err := alt.NewRecomposer("type", map[interface{}]alt.RecomposeFunc{&Dummy{}: nil})
	tt.Nil(t, err, "NewRecomposer")
	_, err = r.Recompose(src)
	tt.NotNil(t, err, "Recompose")
}

func TestRecomposeListResult(t *testing.T) {
	src := []interface{}{
		map[string]interface{}{"type": "Dummy", "val": 1},
		map[string]interface{}{"type": "Dummy", "val": 2},
	}
	r, err := alt.NewRecomposer("type", map[interface{}]alt.RecomposeFunc{&Dummy{}: nil})
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

func TestRecomposeListBadResult(t *testing.T) {
	src := []interface{}{true}
	r, err := alt.NewRecomposer("type", map[interface{}]alt.RecomposeFunc{})
	tt.Nil(t, err, "NewRecomposer")
	_, err = r.Recompose(src, []*Dummy{})
	tt.NotNil(t, err, "Recompose")
}

func TestRecomposeListBadTarget(t *testing.T) {
	r, err := alt.NewRecomposer("type", map[interface{}]alt.RecomposeFunc{})
	tt.Nil(t, err, "NewRecomposer")
	_, err = r.Recompose("[]", 7)
	tt.NotNil(t, err, "Recompose")
}
