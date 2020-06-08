// Copyright (c) 2020, Peter Ohler, All rights reserved.

package jp_test

import (
	"sort"
	"testing"

	"github.com/ohler55/ojg/gen"
	"github.com/ohler55/ojg/jp"
	"github.com/ohler55/ojg/tt"
)

func buildTree(size, depth, iv int) interface{} {
	if depth%2 == 0 {
		list := []interface{}{}
		for i := 0; i < size; i++ {
			nv := iv*10 + i + 1
			if 1 < depth {
				list = append(list, buildTree(size, depth-1, nv))
			} else {
				list = append(list, nv)
			}
		}
		return list
	}
	obj := map[string]interface{}{}
	for i := 0; i < size; i++ {
		k := string([]byte{'a' + byte(i)})
		nv := iv*10 + i + 1
		if 1 < depth {
			obj[k] = buildTree(size, depth-1, nv)
		} else {
			obj[k] = nv
		}
	}
	return obj
}

func buildNodeTree(size, depth, iv int) gen.Node {
	if depth%2 == 0 {
		list := gen.Array{}
		for i := 0; i < size; i++ {
			nv := iv*10 + i + 1
			if 1 < depth {
				list = append(list, buildNodeTree(size, depth-1, nv))
			} else {
				list = append(list, gen.Int(nv))
			}
		}
		return list
	}
	obj := gen.Object{}
	for i := 0; i < size; i++ {
		k := string([]byte{'a' + byte(i)})
		nv := iv*10 + i + 1
		if 1 < depth {
			obj[k] = buildNodeTree(size, depth-1, nv)
		} else {
			obj[k] = gen.Int(nv)
		}
	}
	return obj
}

func TestExprBuild(t *testing.T) {
	x := jp.X().D().C("abc").W().N(3).U(2, "x").S(1, 5, 2).S(1, 5).S(1)
	tt.Equal(t, "..abc.*[3][2,'x'][1:5:2][1:5][1:]", x.String())

	x = jp.R().Descent().Child("abc").Wildcard().Nth(3).Union(2, "x").Slice(1, 5, 2).Slice(1, 5).Slice(1)
	tt.Equal(t, "$..abc.*[3][2,'x'][1:5:2][1:5][1:]", x.String())

	x = jp.B().Descent().Child("abc").Wildcard()
	tt.Equal(t, "[..]['abc'][*]", x.String())

	x = jp.R().B().Descent().Child("abc").Wildcard()
	tt.Equal(t, "$[..]['abc'][*]", x.String())
}

func TestExprGet(t *testing.T) {
	data := buildTree(4, 3, 0)
	x := jp.R().C("a").W().C("b")
	result := x.Get(data)
	sort.Slice(result, func(i, j int) bool {
		iv, _ := result[i].(int)
		jv, _ := result[j].(int)
		return iv < jv
	})
	tt.Equal(t, []interface{}{112, 122, 132, 142}, result)

	x = jp.R().C("b").N(1).C("c")
	result = x.Get(data)
	tt.Equal(t, []interface{}{223}, result)

	x = jp.D().N(1).C("b")
	result = x.Get(data)
	sort.Slice(result, func(i, j int) bool {
		iv, _ := result[i].(int)
		jv, _ := result[j].(int)
		return iv < jv
	})
	tt.Equal(t, []interface{}{122, 222, 322, 422}, result)

	x = jp.U(1, "a").U("b", 2).U("c", 3)
	result = x.Get(data)
	tt.Equal(t, []interface{}{133}, result)

	x = jp.C("a").S(1, -1, 2).C("a")
	result = x.Get(data)
	sort.Slice(result, func(i, j int) bool {
		iv, _ := result[i].(int)
		jv, _ := result[j].(int)
		return iv < jv
	})
	tt.Equal(t, []interface{}{121, 141}, result)

	x = jp.C("a").F(jp.Gt(jp.Get(jp.A().C("a")), jp.ConstInt(135))).C("b")
	result = x.Get(data)
	tt.Equal(t, []interface{}{142}, result)
}

func TestExprGetNodes(t *testing.T) {
	data := buildNodeTree(4, 3, 0)
	x := jp.R().C("a").W().C("b")
	result := x.GetNodes(data)
	sort.Slice(result, func(i, j int) bool {
		iv, _ := result[i].(gen.Int)
		jv, _ := result[j].(gen.Int)
		return iv < jv
	})
	tt.Equal(t, []gen.Node{gen.Int(112), gen.Int(122), gen.Int(132), gen.Int(142)}, result)

	x = jp.R().C("b").N(1).C("c")
	result = x.GetNodes(data)
	tt.Equal(t, []gen.Node{gen.Int(223)}, result)

	x = jp.D().N(1).C("b")
	result = x.GetNodes(data)
	sort.Slice(result, func(i, j int) bool {
		iv, _ := result[i].(gen.Int)
		jv, _ := result[j].(gen.Int)
		return iv < jv
	})
	tt.Equal(t, []gen.Node{gen.Int(122), gen.Int(222), gen.Int(322), gen.Int(422)}, result)
}

func TestExprFirst(t *testing.T) {
	data := buildTree(4, 3, 0)
	x := jp.R().C("a").W().C("b")
	result := x.First(data)
	i, _ := result.(int)
	tt.Equal(t, 1, i/100)
	tt.Equal(t, 2, i%10)

	x = jp.R().C("b").N(1).C("c")
	result = x.First(data)
	tt.Equal(t, 223, result)
}

func TestExprFirstNode(t *testing.T) {
	data := buildNodeTree(4, 3, 0)
	x := jp.R().C("a").W().C("b")
	result := x.FirstNode(data)
	i, _ := result.(gen.Int)
	tt.Equal(t, 1, int(i)/100)
	tt.Equal(t, 2, int(i)%10)

	x = jp.R().C("b").N(1).C("c")
	result = x.FirstNode(data)
	tt.Equal(t, 223, result)
}
