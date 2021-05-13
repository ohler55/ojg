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
	Spouse   *Parent
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
	i, ok := data["val"].(int)
	if !ok {
		return nil, fmt.Errorf("val is not an int")
	}
	return &silly{val: int(i)}, nil
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
	type SillyWrap struct {
		Silly *silly
	}
	src := map[string]interface{}{
		"silly": map[string]interface{}{"type": "silly", "val": 3},
	}
	r, err := alt.NewRecomposer("type", map[interface{}]alt.RecomposeFunc{&silly{}: sillyRecompose})
	tt.Nil(t, err, "NewRecomposer")
	var v interface{}
	var wrap SillyWrap
	v, err = r.Recompose(src, &wrap)
	tt.Nil(t, err, "Recompose")
	w, _ := v.(*SillyWrap)
	tt.NotNil(t, w, "silly wrap")
	tt.Equal(t, 3, w.Silly.val)

	v, err = r.Recompose(src)
	tt.Nil(t, err, "Recompose")
	m, _ := v.(map[string]interface{})
	tt.NotNil(t, m["silly"])

	src = map[string]interface{}{
		"silly": map[string]interface{}{"type": "silly", "val": true},
	}
	_, err = r.Recompose(src, &wrap)
	tt.NotNil(t, err, "Recompose should return and error")
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
		Spouse: &Parent{Child: Child{Name: "Bobby"}},
	}
	simple := alt.Decompose(&src, &alt.Options{OmitNil: true})

	// Since friends is a slice of interfaces a hint is needed to determine
	// the type. Use ^ as an example.
	_ = jp.C("friends").W().C("^").Set(simple, "Child")
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

	diff := alt.Compare(src, v)
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

	diff := alt.Compare(src, v)
	tt.Equal(t, 0, len(diff), "compare to source: diff - ", diff)

	diff = alt.Compare(slice, v)
	tt.Equal(t, 0, len(diff), "compare target and return: diff - ", diff)
}

func TestRecomposePtrMap(t *testing.T) {
	src := map[string]*Child{
		"a": {Name: "Andy"},
		"r": {Name: "Robin"},
	}
	simple := alt.Decompose(src, &alt.Options{})
	r, err := alt.NewRecomposer("", nil)
	tt.Nil(t, err, "NewRecomposer")

	var out map[string]*Child
	var v interface{}
	v, err = r.Recompose(simple, &out)
	tt.Nil(t, err, "Recompose")

	diff := alt.Compare(src, v)
	tt.Equal(t, 0, len(diff), "compare to source: diff - ", diff)

	diff = alt.Compare(out, v)
	tt.Equal(t, 0, len(diff), "compare target and return: diff - ", diff)

	// Try with allocated map.
	out = map[string]*Child{}
	v, err = r.Recompose(simple, &out)
	tt.Nil(t, err, "Recompose")

	diff = alt.Compare(src, v)
	tt.Equal(t, 0, len(diff), "compare to source: diff - ", diff)

	diff = alt.Compare(out, v)
	tt.Equal(t, 0, len(diff), "compare target and return: diff - ", diff)
}

func TestRecomposeNotSettable(t *testing.T) {
	type NotSet struct {
		X chan bool
	}
	src := map[string]interface{}{"x": 3}
	var out NotSet
	_, err := alt.Recompose(src, &out)
	tt.NotNil(t, err, "Recompose")
}

func TestRecomposeMap(t *testing.T) {
	src := map[string]interface{}{"x": 3}
	var out map[string]interface{}
	var v interface{}
	v, err := alt.Recompose(src, &out)
	tt.Nil(t, err, "Recompose")

	diff := alt.Compare(src, v)
	tt.Equal(t, 0, len(diff), "compare to source: diff - ", diff)

	diff = alt.Compare(out, v)
	tt.Equal(t, 0, len(diff), "compare target and return: diff - ", diff)

	out = map[string]interface{}{}
	v, err = alt.Recompose(src, out)
	tt.Nil(t, err, "Recompose")

	diff = alt.Compare(src, v)
	tt.Equal(t, 0, len(diff), "compare to source: diff - ", diff)

	diff = alt.Compare(out, v)
	tt.Equal(t, 0, len(diff), "compare target and return: diff - ", diff)

	v, err = alt.Recompose(src)
	tt.Nil(t, err, "Recompose")

	diff = alt.Compare(src, v)
	tt.Equal(t, 0, len(diff), "compare to source: diff - ", diff)
}

func TestRecomposeBadComposerFunc(t *testing.T) {
	src := map[string]interface{}{"^": "Dummy", "x": 3}
	r, err := alt.NewRecomposer("^",
		map[interface{}]alt.RecomposeFunc{&Dummy{}: func(_ map[string]interface{}) (interface{}, error) {
			return nil, fmt.Errorf("failed")
		}})
	tt.Nil(t, err, "NewRecomposer")
	_, err = r.Recompose(src)
	tt.NotNil(t, err, "Recompose")
}

func TestRecomposeGenComposerFunc(t *testing.T) {
	src := gen.Object{"^": gen.String("Dummy"), "val": gen.Int(3)}
	r, err := alt.NewRecomposer("^",
		map[interface{}]alt.RecomposeFunc{&Dummy{}: func(_ map[string]interface{}) (interface{}, error) {
			return nil, fmt.Errorf("failed")
		}})
	tt.Nil(t, err, "NewRecomposer")
	_, err = r.Recompose(src)
	tt.NotNil(t, err, "Recompose")

	r, err = alt.NewRecomposer("^",
		map[interface{}]alt.RecomposeFunc{&Dummy{}: func(data map[string]interface{}) (interface{}, error) {
			return &Dummy{Val: int(alt.Int(jp.C("val").First(data)))}, nil
		}})
	tt.Nil(t, err, "NewRecomposer")
	var v interface{}
	v, err = r.Recompose(src)
	tt.Nil(t, err, "Recompose")
	d, _ := v.(*Dummy)
	tt.NotNil(t, d)
	tt.Equal(t, 3, d.Val)
}

func TestRecomposeNotSlice(t *testing.T) {
	src := map[string]interface{}{"x": 3}
	var out []interface{}
	_, err := alt.Recompose(src, &out)
	tt.NotNil(t, err, "Recompose")
}

func TestRecomposeNotMap(t *testing.T) {
	src := []interface{}{3}
	var out map[string]interface{}
	_, err := alt.Recompose(src, &out)
	tt.NotNil(t, err, "Recompose")
}

func TestRecomposeOtherMap(t *testing.T) {
	src := map[string]int{"x": 3}
	var out map[string]interface{}
	v, err := alt.Recompose(src, &out)
	tt.Nil(t, err, "Recompose")
	vo, _ := v.(map[string]interface{})
	tt.NotNil(t, vo)
	tt.Equal(t, 3, vo["x"].(int64))
}

func TestRecomposeSimpleMap(t *testing.T) {
	src := map[string]interface{}{"x": map[string]interface{}{"val": 3}}
	var out map[string]Dummy
	v, err := alt.Recompose(src, &out)
	tt.Nil(t, err, "Recompose")
	vo, _ := v.(map[string]Dummy)
	tt.NotNil(t, vo)
	tt.Equal(t, 3, vo["x"].Val)
}

func TestRecomposeAlternateKeys(t *testing.T) {
	src := map[string]interface{}{"Val": 3}
	var out Anno
	v, err := alt.Recompose(src, &out)
	tt.Nil(t, err, "Recompose")
	a, _ := v.(*Anno)
	tt.NotNil(t, a)
	tt.Equal(t, 3, a.Val)

	src = map[string]interface{}{"val": 3}
	v, err = alt.Recompose(src, &out)
	tt.Nil(t, err, "Recompose")
	a, _ = v.(*Anno)
	tt.NotNil(t, a)
	tt.Equal(t, 3, a.Val)

	src = map[string]interface{}{"v": 3}
	v, err = alt.Recompose(src, &out)
	tt.Nil(t, err, "Recompose")
	a, _ = v.(*Anno)
	tt.NotNil(t, a)
	tt.Equal(t, 3, a.Val)
}

func TestRecomposeInterface(t *testing.T) {
	src := map[string]interface{}{"^": "Child", "name": "Pat"}

	r, err := alt.NewRecomposer("^", map[interface{}]alt.RecomposeFunc{&Child{}: nil})
	tt.Nil(t, err, "NewRecomposer")

	var out fmt.Stringer
	v, err := r.Recompose(src, &out)
	tt.Nil(t, err, "Recompose")
	s, _ := v.(fmt.Stringer)
	tt.NotNil(t, s)
	tt.Equal(t, "Pat", s.String())
}

type BooFlu struct {
	Boo bool    `json:",string"`
	Flu float64 `json:",string"`
}

func TestRecomposeTags(t *testing.T) {
	src := map[string]interface{}{"v": 3, "Title": 2, "skip": 7, "-": 4, "str": "1"}
	var out Anno
	v, err := alt.Recompose(src, &out)
	tt.Nil(t, err, "Recompose")
	a, _ := v.(*Anno)
	tt.NotNil(t, a)
	tt.Equal(t, 3, a.Val)
	tt.Equal(t, 2, a.Title)
	tt.Equal(t, 0, a.Skip)
	tt.Equal(t, 1, a.Str)
	tt.Equal(t, 4, a.Dash)

	src = map[string]interface{}{"str": "1x"}
	_, err = alt.Recompose(src, &out)
	tt.NotNil(t, err, "Recompose tag str invalid")

	var bf BooFlu
	src = map[string]interface{}{"boo": "true", "flu": "1.23"}
	_, err = alt.Recompose(src, &bf)
	tt.Nil(t, err, "Recompose tag string ok")

	src = map[string]interface{}{"boo": true, "flu": 1.23}
	_, err = alt.Recompose(src, &bf)
	tt.Nil(t, err, "Recompose tag not string")

	src = map[string]interface{}{"boo": "yes"}
	_, err = alt.Recompose(src, &bf)
	tt.NotNil(t, err, "Recompose tag invalid string")

	src = map[string]interface{}{"boo": "true", "flu": "1x2"}
	_, err = alt.Recompose(src, &bf)
	tt.NotNil(t, err, "Recompose tag invalid string")
}

func TestRecomposeNil(t *testing.T) {
	r, err := alt.NewRecomposer("", nil)
	tt.Nil(t, err, "NewRecomposer")
	var v interface{}
	var a []interface{}
	v, err = r.Recompose(nil, &a)
	tt.Nil(t, err, "Recompose")
	tt.Equal(t, []interface{}{}, v)

	var list WithList
	v, err = r.Recompose(map[string]interface{}{"list": nil}, &list)
	tt.Nil(t, err, "Recompose")
	l2, _ := v.(*WithList)
	tt.NotNil(t, l2)

	m := map[string]interface{}{}
	v, err = r.Recompose(nil, m)
	tt.Nil(t, err, "Recompose")
	tt.Equal(t, map[string]interface{}{}, v)

	var d Dummy
	v, err = r.Recompose(nil, &d)
	tt.Nil(t, err, "Recompose")
	d2, _ := v.(*Dummy)
	tt.NotNil(t, d2)

}

func TestMustNewRecomposePanic(t *testing.T) {
	tt.Panic(t, func() {
		_ = alt.MustNewRecomposer("^", nil, map[interface{}]alt.RecomposeAnyFunc{true: nil})
	})
}

func TestRecomposerRegister(t *testing.T) {
	type Sample struct {
		Int  int
		When time.Time
	}
	r := alt.MustNewRecomposer("^", nil)
	r.RegisterComposer(&Sample{}, nil)
	r.RegisterAnyComposer(time.Time{},
		func(v interface{}) (interface{}, error) {
			if secs, ok := v.(int); ok {
				return time.Unix(int64(secs), 0), nil
			}
			return nil, fmt.Errorf("can not convert a %T to a time.Time", v)
		})
	data := map[string]interface{}{"^": "Sample", "int": 3, "when": 1612872722}
	v := r.MustRecompose(data)
	sample, _ := v.(*Sample)
	tt.NotNil(t, sample)
	tt.Equal(t, 3, sample.Int)
	tt.Equal(t, int64(1612872722), sample.When.Unix())
}

func TestRecomposeReflectBool(t *testing.T) {
	type Sample struct {
		Boo bool
	}
	r := alt.MustNewRecomposer("^", map[interface{}]alt.RecomposeFunc{&Sample{}: nil})
	data := map[string]interface{}{"^": "Sample", "boo": true}
	var sample Sample
	v := r.MustRecompose(data, &sample)
	tt.NotNil(t, v)
}

func TestRecomposerAnyComposePtr(t *testing.T) {
	type Sample struct {
		When time.Time
	}
	r := alt.MustNewRecomposer("^", nil)
	r.RegisterAnyComposer(time.Time{},
		func(v interface{}) (interface{}, error) {
			if secs, ok := v.(int); ok {
				t := time.Unix(int64(secs), 0)
				return &t, nil
			}
			return nil, fmt.Errorf("can not convert a %T to a time.Time", v)
		})
	data := map[string]interface{}{"^": "Sample", "when": 1612872722}
	var sample Sample
	_ = r.MustRecompose(data, &sample)
	tt.Equal(t, int64(1612872722), sample.When.Unix())

	data = map[string]interface{}{"^": "Sample", "when": true}
	tt.Panic(t, func() {
		_ = r.MustRecompose(data, &sample)
	})
}
