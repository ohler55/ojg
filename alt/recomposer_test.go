// Copyright (c) 2020, Peter Ohler, All rights reserved.

package alt_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/ohler55/ojg/alt"
	"github.com/ohler55/ojg/gen"
	"github.com/ohler55/ojg/jp"
	"github.com/ohler55/ojg/tt"
)

type WithList struct {
	List []int
	Fun  func() bool
}

type Parent struct {
	Child
	Num      int
	Children []*Child
	Friends  []fmt.Stringer
}

type Child struct {
	Name string
}

func (c *Child) String() string {
	return c.Name
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
		s.a = alt.Int(val)
	case "b":
		s.b, _ = val.(string)
	default:
		return fmt.Errorf("%s is not an attribute of Setter", attr)
	}
	return nil
}

func sillyRecompose(data map[string]interface{}) (interface{}, error) {
	s := silly{}
	i, _ := data["val"].(int)
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

	src = map[string]interface{}{"type": "Setter", "a": 3, "b": "bee", "c": 5}
	_, err = r.Recompose(src)
	tt.NotNil(t, err, "Recompose from bad source")
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

func TestRecomposeArrayResult(t *testing.T) {
	src := gen.Array{
		gen.Object{"type": gen.String("Dummy"), "val": gen.Int(1)},
		gen.Object{"type": gen.String("Dummy"), "val": gen.Int(2)},
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

	src = gen.Array{
		gen.Object{"type": gen.String("Dummy"), "nest": gen.Object{"type": gen.String("Dummy"), "val": gen.String("x")}},
	}
	_, err = r.Recompose(src, []*Dummy{})
	tt.NotNil(t, err, "Recompose from bad source")
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

func TestRecomposeNested(t *testing.T) {
	src := Parent{
		Child: Child{Name: "Pat"},
		Num:   3,
		Children: []*Child{
			{Name: "Andy"},
			{Name: "Robin"},
		},
		Friends: []fmt.Stringer{
			&Child{Name: "Ash"},
			&Child{Name: "Riley"},
		},
	}
	simple := alt.Decompose(&src, &alt.Options{})

	jp.C("child").Del(simple)
	jp.C("name").Set(simple, "Pat")

	// Since friends is a slice of interfaces a hint is needed to determine
	// the type. Use ^ as an example.
	jp.C("friends").W().C("^").Set(simple, "Child")
	// Make sure the recomposer knows about the Child type so the hint has
	// something to refer to.
	r, err := alt.NewRecomposer("^", map[interface{}]alt.RecomposeFunc{&Parent{}: nil})
	tt.Nil(t, err, "NewRecomposer")

	var v interface{}
	v, err = r.Recompose(simple, &Parent{})
	tt.Nil(t, err, "Recompose")
	p, _ := v.(*Parent)
	tt.NotNil(t, p, "check type - %"+"T", v)

	diff := alt.Compare(&src, p)
	tt.Equal(t, 0, len(diff), "compare diff - ", diff)
}

func TestRecomposeSlice(t *testing.T) {
	src := []Child{
		{Name: "Andy"},
		{Name: "Robin"},
	}
	simple := alt.Decompose(&src, &alt.Options{})

	r, err := alt.NewRecomposer("", nil)
	tt.Nil(t, err, "NewRecomposer")
	var slice []Child
	var v interface{}
	v, err = r.Recompose(simple, &slice)
	tt.Nil(t, err, "Recompose")

	diff := alt.Compare(&src, v)
	tt.Equal(t, 0, len(diff), "compare to source: diff - ", diff)

	diff = alt.Compare(slice, v)
	tt.Equal(t, 0, len(diff), "compare target and return: diff - ", diff)
}

func TestRecomposePtrSlice(t *testing.T) {
	src := []*Child{
		{Name: "Andy"},
		{Name: "Robin"},
	}
	simple := alt.Decompose(&src, &alt.Options{})
	r, err := alt.NewRecomposer("", nil)
	tt.Nil(t, err, "NewRecomposer")

	var slice []*Child
	var v interface{}
	v, err = r.Recompose(simple, &slice)
	tt.Nil(t, err, "Recompose")

	diff := alt.Compare(&src, v)
	tt.Equal(t, 0, len(diff), "compare to source: diff - ", diff)

	diff = alt.Compare(slice, v)
	tt.Equal(t, 0, len(diff), "compare target and return: diff - ", diff)
}

func TestRecomposePtrMap(t *testing.T) {
	src := map[string]*Child{
		"a": {Name: "Andy"},
		"r": {Name: "Robin"},
	}
	simple := alt.Decompose(&src, &alt.Options{})
	r, err := alt.NewRecomposer("", nil)
	tt.Nil(t, err, "NewRecomposer")

	var out map[string]*Child
	var v interface{}
	v, err = r.Recompose(simple, &out)
	tt.Nil(t, err, "Recompose")

	diff := alt.Compare(&src, v)
	tt.Equal(t, 0, len(diff), "compare to source: diff - ", diff)

	diff = alt.Compare(out, v)
	tt.Equal(t, 0, len(diff), "compare target and return: diff - ", diff)

	// Try with allocated map.
	out = map[string]*Child{}
	v, err = r.Recompose(simple, &out)
	tt.Nil(t, err, "Recompose")

	diff = alt.Compare(&src, v)
	tt.Equal(t, 0, len(diff), "compare to source: diff - ", diff)

	diff = alt.Compare(out, v)
	tt.Equal(t, 0, len(diff), "compare target and return: diff - ", diff)
}
