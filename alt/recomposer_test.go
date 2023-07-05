// Copyright (c) 2020, Peter Ohler, All rights reserved.

package alt_test

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/ohler55/ojg"
	"github.com/ohler55/ojg/alt"
	"github.com/ohler55/ojg/gen"
	"github.com/ohler55/ojg/jp"
	"github.com/ohler55/ojg/oj"
	"github.com/ohler55/ojg/pretty"
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

type PickANumber struct {
	AsString string
	AsFloat  float64
	AsInt    int64
	AsNumber json.Number
}

func (c *Child) String() string {
	return c.Name
}

func sillyRecompose(data map[string]any) (any, error) {
	i, ok := data["val"].(int)
	if !ok {
		return nil, fmt.Errorf("val is not an int")
	}
	return &silly{val: i}, nil
}

func TestRecomposeBasic(t *testing.T) {
	src := map[string]any{
		"type": "Dummy",
		"val":  3,
		"nest": []any{
			int8(-8), int16(-16), int32(-32),
			uint(0), uint8(8), uint16(16), uint32(32), uint64(64),
			float32(1.2),
			map[string]any{},
		},
	}
	r, err := alt.NewRecomposer("type", map[any]alt.RecomposeFunc{&Dummy{}: nil})
	tt.Nil(t, err, "NewRecomposer")
	var v any
	v, err = r.Recompose(src)
	tt.Nil(t, err, "Recompose")
	d, _ := v.(*Dummy)
	tt.NotNil(t, d, "Dummy")
	tt.Equal(t, []any{-8, -16, -32, 0, 8, 16, 32, 64, 1.2, map[string]any{}}, d.Nest)
}

func TestRecomposeNode(t *testing.T) {
	tm := time.Date(2020, time.April, 12, 16, 34, 04, 123456789, time.UTC)
	src := map[string]any{
		"type": "Dummy",
		"val":  gen.Int(3),
		"nest": gen.Array{gen.Int(-8), gen.Bool(true), gen.Float(1.2), gen.String("abc"),
			gen.Object{"big": gen.Big("123"), "time": gen.Time(tm)},
		},
	}
	r, err := alt.NewRecomposer("type", map[any]alt.RecomposeFunc{&Dummy{}: nil})
	tt.Nil(t, err, "NewRecomposer")
	var v any
	v, err = r.Recompose(src)
	tt.Nil(t, err, "Recompose")
	d, _ := v.(*Dummy)
	tt.NotNil(t, d, "Dummy")
	tt.Equal(t, []any{-8, true, 1.2, "abc", map[string]any{"big": "123", "time": tm}}, d.Nest)
}

func TestRecomposeFunc(t *testing.T) {
	type SillyWrap struct {
		Silly *silly
	}
	src := map[string]any{
		"silly": map[string]any{"type": "silly", "val": 3},
	}
	r, err := alt.NewRecomposer("type", map[any]alt.RecomposeFunc{&silly{}: sillyRecompose})
	tt.Nil(t, err, "NewRecomposer")
	var v any
	var wrap SillyWrap
	v, err = r.Recompose(src, &wrap)
	tt.Nil(t, err, "Recompose")
	w, _ := v.(*SillyWrap)
	tt.NotNil(t, w, "silly wrap")
	tt.Equal(t, 3, w.Silly.val)

	v, err = r.Recompose(src)
	tt.Nil(t, err, "Recompose")
	m, _ := v.(map[string]any)
	tt.NotNil(t, m["silly"])

	src = map[string]any{
		"silly": map[string]any{"type": "silly", "val": true},
	}
	_, err = r.Recompose(src, &wrap)
	tt.NotNil(t, err, "Recompose should return and error")
}

func TestRecomposeReflect(t *testing.T) {
	src := map[string]any{"type": "Dummy", "val": 3, "extra": true, "fun": true}
	r, err := alt.NewRecomposer("type", map[any]alt.RecomposeFunc{&Dummy{}: nil})
	tt.Nil(t, err, "NewRecomposer")
	var v any
	v, err = r.Recompose(src)
	tt.Nil(t, err, "Recompose")
	d, _ := v.(*Dummy)
	tt.NotNil(t, d, "check type")
	tt.Equal(t, 3, d.Val)
}

func TestRecomposeAttrSetter(t *testing.T) {
	src := map[string]any{"type": "Setter", "a": 3, "b": "bee"}
	r, err := alt.NewRecomposer("type", map[any]alt.RecomposeFunc{&Setter{}: nil})
	tt.Nil(t, err, "NewRecomposer")
	var v any
	v, err = r.Recompose(src)
	tt.Nil(t, err, "Recompose")
	s, _ := v.(*Setter)
	tt.NotNil(t, s, "check type")
	tt.Equal(t, "Setter{a:3,b:bee}", s.String())

	src = map[string]any{"type": "Setter", "a": 3, "b": "bee", "c": 5}
	_, err = r.Recompose(src)
	tt.NotNil(t, err, "Recompose from bad source")
}

func TestRecomposeReflectList(t *testing.T) {
	src := map[string]any{"type": "WithList", "list": []any{1, 2, 3}}
	r, err := alt.NewRecomposer("type", map[any]alt.RecomposeFunc{&WithList{}: nil})
	tt.Nil(t, err, "NewRecomposer")
	var v any
	v, err = r.Recompose(src)
	tt.Nil(t, err, "Recompose")
	wl, _ := v.(*WithList)
	tt.NotNil(t, wl, "check type")
	tt.Equal(t, "[]int [1 2 3]", fmt.Sprintf("%T %v", wl.List, wl.List))
}

func TestRecomposeBadMap(t *testing.T) {
	_, err := alt.NewRecomposer("type", map[any]alt.RecomposeFunc{3: nil})
	tt.NotNil(t, err, "NewRecomposer")
}

func TestRecomposeBadField(t *testing.T) {
	src := map[string]any{"type": "Dummy", "val": true}
	r, err := alt.NewRecomposer("type", map[any]alt.RecomposeFunc{&Dummy{}: nil})
	tt.Nil(t, err, "NewRecomposer")
	_, err = r.Recompose(src)
	tt.NotNil(t, err, "Recompose")
}

func TestRecomposeReflectListBad(t *testing.T) {
	src := map[string]any{"type": "WithList", "list": []any{1, true, 3}}
	r, err := alt.NewRecomposer("type", map[any]alt.RecomposeFunc{&WithList{}: nil})
	tt.Nil(t, err, "NewRecomposer")
	_, err = r.Recompose(src)
	tt.NotNil(t, err, "Recompose")
}

func TestRecomposeBadListItem(t *testing.T) {
	src := map[string]any{
		"type": "Dummy",
		"val":  3,
		"nest": []any{func() {}},
	}
	r, err := alt.NewRecomposer("type", map[any]alt.RecomposeFunc{&Dummy{}: nil})
	tt.Nil(t, err, "NewRecomposer")
	_, err = r.Recompose(src)
	tt.NotNil(t, err, "Recompose")
}

func TestRecomposeListResult(t *testing.T) {
	src := []any{
		map[string]any{"type": "Dummy", "val": 1},
		map[string]any{"type": "Dummy", "val": 2},
	}
	r, err := alt.NewRecomposer("type", map[any]alt.RecomposeFunc{&Dummy{}: nil})
	tt.Nil(t, err, "NewRecomposer")
	var v any
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
	r, err := alt.NewRecomposer("type", map[any]alt.RecomposeFunc{&Dummy{}: nil})
	tt.Nil(t, err, "NewRecomposer")
	var v any
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
	src := []any{true}
	r, err := alt.NewRecomposer("type", map[any]alt.RecomposeFunc{})
	tt.Nil(t, err, "NewRecomposer")
	_, err = r.Recompose(src, []*Dummy{})
	tt.NotNil(t, err, "Recompose")
}

func TestRecomposeListBadTarget(t *testing.T) {
	r, err := alt.NewRecomposer("type", map[any]alt.RecomposeFunc{})
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
	r, err := alt.NewRecomposer("^", map[any]alt.RecomposeFunc{&Parent{}: nil})
	tt.Nil(t, err, "NewRecomposer")

	var v any
	v, err = r.Recompose(simple, &Parent{})
	tt.Nil(t, err, "Recompose")
	p, _ := v.(*Parent)
	format := "check type - %T"
	tt.NotNil(t, p, format, v)

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
	var v any
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
	var v any
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
	var v any
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
	src := map[string]any{"x": 3}
	var out NotSet
	_, err := alt.Recompose(src, &out)
	tt.NotNil(t, err, "Recompose")
}

func TestRecomposeMap(t *testing.T) {
	src := map[string]any{"x": 3}
	var out map[string]any
	var v any
	v, err := alt.Recompose(src, &out)
	tt.Nil(t, err, "Recompose")

	diff := alt.Compare(src, v)
	tt.Equal(t, 0, len(diff), "compare to source: diff - ", diff)

	diff = alt.Compare(out, v)
	tt.Equal(t, 0, len(diff), "compare target and return: diff - ", diff)

	out = map[string]any{}
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
	src := map[string]any{"^": "Dummy", "x": 3}
	r, err := alt.NewRecomposer("^",
		map[any]alt.RecomposeFunc{&Dummy{}: func(_ map[string]any) (any, error) {
			return nil, fmt.Errorf("failed")
		}})
	tt.Nil(t, err, "NewRecomposer")
	_, err = r.Recompose(src)
	tt.NotNil(t, err, "Recompose")
}

func TestRecomposeGenComposerFunc(t *testing.T) {
	src := gen.Object{"^": gen.String("Dummy"), "val": gen.Int(3)}
	r, err := alt.NewRecomposer("^",
		map[any]alt.RecomposeFunc{&Dummy{}: func(_ map[string]any) (any, error) {
			return nil, fmt.Errorf("failed")
		}})
	tt.Nil(t, err, "NewRecomposer")
	_, err = r.Recompose(src)
	tt.NotNil(t, err, "Recompose")

	r, err = alt.NewRecomposer("^",
		map[any]alt.RecomposeFunc{&Dummy{}: func(data map[string]any) (any, error) {
			return &Dummy{Val: int(alt.Int(jp.C("val").First(data)))}, nil
		}})
	tt.Nil(t, err, "NewRecomposer")
	var v any
	v, err = r.Recompose(src)
	tt.Nil(t, err, "Recompose")
	d, _ := v.(*Dummy)
	tt.NotNil(t, d)
	tt.Equal(t, 3, d.Val)
}

func TestRecomposeNotSlice(t *testing.T) {
	src := map[string]any{"x": 3}
	var out []any
	_, err := alt.Recompose(src, &out)
	tt.NotNil(t, err, "Recompose")
}

func TestRecomposeNotMap(t *testing.T) {
	src := []any{3}
	var out map[string]any
	_, err := alt.Recompose(src, &out)
	tt.NotNil(t, err, "Recompose")
}

func TestRecomposeOtherMap(t *testing.T) {
	src := map[string]int{"x": 3}
	var out map[string]any
	v, err := alt.Recompose(src, &out)
	tt.Nil(t, err, "Recompose")
	vo, _ := v.(map[string]any)
	tt.NotNil(t, vo)
	tt.Equal(t, 3, vo["x"].(int64))
}

func TestRecomposeOtherMap2(t *testing.T) {
	src := map[string]int64{"x": 3}
	var out map[string]any
	v, err := alt.Recompose(src, &out)
	tt.Nil(t, err, "Recompose")
	vo, _ := v.(map[string]any)
	tt.NotNil(t, vo)
	tt.Equal(t, 3, vo["x"].(int64))
}

func TestRecomposeSimpleMap(t *testing.T) {
	src := map[string]any{"x": map[string]any{"val": 3}}
	var out map[string]Dummy
	v, err := alt.Recompose(src, &out)
	tt.Nil(t, err, "Recompose")
	vo, _ := v.(map[string]Dummy)
	tt.NotNil(t, vo)
	tt.Equal(t, 3, vo["x"].Val)
}

func TestRecomposeAlternateKeys(t *testing.T) {
	src := map[string]any{"Val": 3}
	var out Anno
	v, err := alt.Recompose(src, &out)
	tt.Nil(t, err, "Recompose")
	a, _ := v.(*Anno)
	tt.NotNil(t, a)
	tt.Equal(t, 3, a.Val)

	src = map[string]any{"val": 3}
	v, err = alt.Recompose(src, &out)
	tt.Nil(t, err, "Recompose")
	a, _ = v.(*Anno)
	tt.NotNil(t, a)
	tt.Equal(t, 3, a.Val)

	src = map[string]any{"v": 3}
	v, err = alt.Recompose(src, &out)
	tt.Nil(t, err, "Recompose")
	a, _ = v.(*Anno)
	tt.NotNil(t, a)
	tt.Equal(t, 3, a.Val)
}

func TestRecomposeInterface(t *testing.T) {
	src := map[string]any{"^": "Child", "name": "Pat"}

	r, err := alt.NewRecomposer("^", map[any]alt.RecomposeFunc{&Child{}: nil})
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
	src := map[string]any{"v": 3, "Title": 2, "skip": 7, "-": 4, "str": "1"}
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

	src = map[string]any{"str": "1x"}
	_, err = alt.Recompose(src, &out)
	tt.NotNil(t, err, "Recompose tag str invalid")

	var bf BooFlu
	src = map[string]any{"boo": "true", "flu": "1.23"}
	_, err = alt.Recompose(src, &bf)
	tt.Nil(t, err, "Recompose tag string ok")

	src = map[string]any{"boo": true, "flu": 1.23}
	_, err = alt.Recompose(src, &bf)
	tt.Nil(t, err, "Recompose tag not string")

	src = map[string]any{"boo": "yes"}
	_, err = alt.Recompose(src, &bf)
	tt.NotNil(t, err, "Recompose tag invalid string")

	src = map[string]any{"boo": "true", "flu": "1x2"}
	_, err = alt.Recompose(src, &bf)
	tt.NotNil(t, err, "Recompose tag invalid string")
}

func TestRecomposeNil(t *testing.T) {
	r, err := alt.NewRecomposer("", nil)
	tt.Nil(t, err, "NewRecomposer")
	var v any
	var a []any
	v, err = r.Recompose(nil, &a)
	tt.Nil(t, err, "Recompose")
	tt.Equal(t, []any{}, v)

	var list WithList
	v, err = r.Recompose(map[string]any{"list": nil}, &list)
	tt.Nil(t, err, "Recompose")
	l2, _ := v.(*WithList)
	tt.NotNil(t, l2)

	m := map[string]any{}
	v, err = r.Recompose(nil, m)
	tt.Nil(t, err, "Recompose")
	tt.Equal(t, map[string]any{}, v)

	var d Dummy
	v, err = r.Recompose(nil, &d)
	tt.Nil(t, err, "Recompose")
	d2, _ := v.(*Dummy)
	tt.NotNil(t, d2)

}

func TestMustNewRecomposePanic(t *testing.T) {
	tt.Panic(t, func() {
		_ = alt.MustNewRecomposer("^", nil, map[any]alt.RecomposeAnyFunc{true: nil})
	})
}

func TestRecomposerRegister(t *testing.T) {
	type Sample struct {
		Int  int
		When time.Time
	}
	r := alt.MustNewRecomposer("^", nil)
	err := r.RegisterComposer(&Sample{}, nil)
	tt.Nil(t, err)
	err = r.RegisterAnyComposer(time.Time{},
		func(v any) (any, error) {
			if secs, ok := v.(int); ok {
				return time.Unix(int64(secs), 0), nil
			}
			return nil, fmt.Errorf("can not convert a %T to a time.Time", v)
		})
	tt.Nil(t, err)
	data := map[string]any{"^": "Sample", "int": 3, "when": 1612872722}
	v := r.MustRecompose(data)
	sample, _ := v.(*Sample)
	tt.NotNil(t, sample)
	tt.Equal(t, 3, sample.Int)
	tt.Equal(t, int64(1612872722), sample.When.Unix())

	// Register two composers for time.
	r = alt.MustNewRecomposer("^", nil)
	err = r.RegisterComposer(&Sample{}, nil)
	tt.Nil(t, err)
	err = r.RegisterAnyComposer(time.Time{},
		func(v any) (any, error) {
			if secs, ok := v.(int); ok {
				return time.Unix(int64(secs), 0), nil
			}
			return nil, nil
		})
	tt.Nil(t, err)
	err = r.RegisterComposer(time.Time{},
		func(v map[string]any) (any, error) {
			for _, m := range v {
				if secs, ok := m.(int); ok {
					return time.Unix(int64(secs), 0), nil
				}
				break
			}
			return nil, fmt.Errorf("can not convert a %T to a time.Time", v)
		})
	tt.Nil(t, err)
	data = map[string]any{"^": "Sample", "int": 3, "when": map[string]any{"@": 1612872722}}
	v = r.MustRecompose(data)
	sample, _ = v.(*Sample)
	tt.NotNil(t, sample)
	tt.Equal(t, 3, sample.Int)
	tt.Equal(t, int64(1612872722), sample.When.Unix())

	data = map[string]any{"^": "Sample", "int": 3, "when": true}
	v = r.MustRecompose(data)
	sample, _ = v.(*Sample)
	tt.NotNil(t, sample)
	tt.Equal(t, time.Time{}.Unix(), sample.When.Unix())
}

func TestRecomposeReflectBool(t *testing.T) {
	type Sample struct {
		Boo bool
	}
	r := alt.MustNewRecomposer("^", map[any]alt.RecomposeFunc{&Sample{}: nil})
	data := map[string]any{"^": "Sample", "boo": true}
	var sample Sample
	v := r.MustRecompose(data, &sample)
	tt.NotNil(t, v)
}

func TestRecomposerAnyComposePtr(t *testing.T) {
	type Sample struct {
		When time.Time
	}
	r := alt.MustNewRecomposer("^", nil)
	err := r.RegisterAnyComposer(time.Time{},
		func(v any) (any, error) {
			if secs, ok := v.(int); ok {
				t := time.Unix(int64(secs), 0)
				return &t, nil
			}
			return nil, fmt.Errorf("can not convert a %T to a time.Time", v)
		})
	tt.Nil(t, err)
	data := map[string]any{"^": "Sample", "when": 1612872722}
	var sample Sample
	_ = r.MustRecompose(data, &sample)
	tt.Equal(t, int64(1612872722), sample.When.Unix())

	data = map[string]any{"^": "Sample", "when": true}
	tt.Panic(t, func() {
		_ = r.MustRecompose(data, &sample)
	})
}

func TestRecomposeReflectPrivate(t *testing.T) {
	type Sample struct {
		X int
		y int
	}
	r := alt.MustNewRecomposer("^", map[any]alt.RecomposeFunc{&Sample{}: nil})
	data := map[string]any{"^": "Sample", "x": 1, "y": 2}
	var sample Sample
	v := r.MustRecompose(data, &sample)
	tt.NotNil(t, v)
	tt.Equal(t, 0, sample.y)
	tt.Equal(t, 1, sample.X)
}

func TestRecomposeEmbeddedMap(t *testing.T) {
	type Sample struct {
		Ss map[string]string
		Sb map[string]bool
	}
	r := alt.MustNewRecomposer("^", nil)
	src := map[string]any{
		"ss": map[string]any{"one": "two"},
		"sb": map[string]any{"yes": true},
	}
	var sample Sample
	_ = r.MustRecompose(src, &sample)
	tt.Equal(t, "two", sample.Ss["one"])
	tt.Equal(t, true, sample.Sb["yes"])
}

type TagMap map[string]any

func (tm *TagMap) UnmarshalJSON(data []byte) error {
	*tm = map[string]any{}
	simple, err := oj.Parse(data)
	if err != nil {
		return err
	}
	for _, kv := range simple.([]any) {
		k := jp.C("key").First(kv).(string)
		if k == "fail" {
			return fmt.Errorf("intentional fail")
		}
		(*tm)[k] = jp.C("value").First(kv)
	}
	return nil
}

func recomposeToJSON(v any) (any, error) {
	return []byte(oj.JSON(v)), nil
}

func TestRecomposeUnmarshaller(t *testing.T) {
	src := []any{map[string]any{"key": "k1", "value": 1}}

	r := alt.MustNewRecomposer("^", nil)
	r.RegisterUnmarshalerComposer(recomposeToJSON)

	var tags TagMap
	_ = alt.MustRecompose(src, &tags)
	tt.Equal(t, 1, len(tags))
	tt.Equal(t, 1, tags["k1"])

	src = []any{map[string]any{"key": "fail", "value": 1}}

	_, err := alt.Recompose(src, &tags)
	tt.NotNil(t, err)
}

func TestRecomposeUnmarshallerField(t *testing.T) {
	type Wrap struct {
		X TagMap
	}
	src := map[string]any{"X": []any{map[string]any{"key": "k1", "value": 1}}}

	r := alt.MustNewRecomposer("^", nil)
	r.RegisterUnmarshalerComposer(recomposeToJSON)

	var wrap Wrap
	_ = alt.MustRecompose(src, &wrap)

	tt.Equal(t, 1, len(wrap.X))
	tt.Equal(t, 1, wrap.X["k1"])

	src = map[string]any{"X": []any{map[string]any{"key": "fail", "value": 1}}}
	_, err := alt.Recompose(src, &wrap)
	tt.NotNil(t, err)
}

func TestRecomposeUnmarshallerList(t *testing.T) {
	src := []any{[]any{map[string]any{"key": "k1", "value": 1}}}

	r := alt.MustNewRecomposer("^", nil)
	r.RegisterUnmarshalerComposer(recomposeToJSON)

	var list []TagMap
	_ = alt.MustRecompose(src, &list)
	tt.Equal(t, 1, len(list))
	tt.Equal(t, 1, len(list[0]))
	tt.Equal(t, 1, list[0]["k1"])
}

func TestRecomposeNumber(t *testing.T) {
	src := json.Number("0.1234567890123456789")

	r, err := alt.NewRecomposer("type", nil)
	tt.Nil(t, err, "NewRecomposer")
	r.NumConvMethod = ojg.NumConvFloat64

	var v any
	v, err = r.Recompose(src)
	tt.Nil(t, err, "Recompose")
	tt.Equal(t, 0.123456789012345678, v)

	tt.Panic(t, func() { _ = r.MustRecompose(json.Number("1.2.3")) })

	r.NumConvMethod = ojg.NumConvString
	v, err = r.Recompose(src)
	tt.Nil(t, err, "Recompose")
	tt.Equal(t, "0.1234567890123456789", v)
}

func TestRecomposeReflectNumber(t *testing.T) {
	src := map[string]any{
		"asString": json.Number("0.1234567890123456789"),
		"asFloat":  json.Number("0.1234567890123456789"),
		"asInt":    json.Number("1234567890123456789"),
		"asNumber": json.Number("0.1234567890123456789"),
	}
	r, err := alt.NewRecomposer("type", nil)
	tt.Nil(t, err, "NewRecomposer")

	var pan PickANumber

	_, err = r.Recompose(src, &pan)
	tt.Nil(t, err, "Recompose")

	tt.Equal(t, `{
  asFloat: 0.12345678901234568
  asInt: 1234567890123456789
  asNumber: "0.1234567890123456789"
  asString: "0.1234567890123456789"
}`, pretty.SEN(&pan))

	src["asFloat"] = json.Number("1.2.3")
	tt.Panic(t, func() { _ = r.MustRecompose(src, &pan) })

	src["asFloat"] = json.Number("123.4")
	src["asInt"] = json.Number("123.4")
	tt.Panic(t, func() { _ = r.MustRecompose(src, &pan) })
}
