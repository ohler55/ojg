// Copyright (c) 2020, Peter Ohler, All rights reserved.

package oj_test

import (
	"fmt"
	"sort"
	"testing"

	"github.com/ohler55/ojg/oj"
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

func buildNodeTree(size, depth, iv int) oj.Node {
	if depth%2 == 0 {
		list := oj.Array{}
		for i := 0; i < size; i++ {
			nv := iv*10 + i + 1
			if 1 < depth {
				list = append(list, buildNodeTree(size, depth-1, nv))
			} else {
				list = append(list, oj.Int(nv))
			}
		}
		return list
	}
	obj := oj.Object{}
	for i := 0; i < size; i++ {
		k := string([]byte{'a' + byte(i)})
		nv := iv*10 + i + 1
		if 1 < depth {
			obj[k] = buildNodeTree(size, depth-1, nv)
		} else {
			obj[k] = oj.Int(nv)
		}
	}
	return obj
}

func TestOjExprBuild(t *testing.T) {
	x := oj.X().D().C("abc").W().N(3)
	tt.Equal(t, "..abc.*[3]", x.String())

	x = oj.R().Descent().Child("abc").Wildcard().Nth(3)
	tt.Equal(t, "$..abc.*[3]", x.String())

	x = oj.B().Descent().Child("abc").Wildcard()
	tt.Equal(t, "[..]['abc'][*]", x.String())

	x = oj.R().B().Descent().Child("abc").Wildcard()
	tt.Equal(t, "$[..]['abc'][*]", x.String())
}

func TestOjExprGet(t *testing.T) {
	data := buildTree(4, 3, 0)
	x := oj.R().C("a").W().C("b")
	result := x.Get(data)
	sort.Slice(result, func(i, j int) bool {
		iv, _ := result[i].(int)
		jv, _ := result[j].(int)
		return iv < jv
	})
	tt.Equal(t, []interface{}{112, 122, 132, 142}, result)

	x = oj.R().C("b").N(1).C("c")
	result = x.Get(data)
	tt.Equal(t, []interface{}{223}, result)

	x = oj.D().N(1).C("b")
	result = x.Get(data)
	sort.Slice(result, func(i, j int) bool {
		iv, _ := result[i].(int)
		jv, _ := result[j].(int)
		return iv < jv
	})
	tt.Equal(t, []interface{}{122, 222, 322, 422}, result)
}

func TestOjExprGetNodes(t *testing.T) {
	data := buildNodeTree(4, 3, 0)
	x := oj.R().C("a").W().C("b")
	result := x.GetNodes(data)
	sort.Slice(result, func(i, j int) bool {
		iv, _ := result[i].(oj.Int)
		jv, _ := result[j].(oj.Int)
		return iv < jv
	})
	tt.Equal(t, []oj.Node{oj.Int(112), oj.Int(122), oj.Int(132), oj.Int(142)}, result)

	x = oj.R().C("b").N(1).C("c")
	result = x.GetNodes(data)
	tt.Equal(t, []oj.Node{oj.Int(223)}, result)

	x = oj.D().N(1).C("b")
	result = x.GetNodes(data)
	sort.Slice(result, func(i, j int) bool {
		iv, _ := result[i].(oj.Int)
		jv, _ := result[j].(oj.Int)
		return iv < jv
	})
	tt.Equal(t, []oj.Node{oj.Int(122), oj.Int(222), oj.Int(322), oj.Int(422)}, result)
}

func TestOjExprFirst(t *testing.T) {
	data := buildTree(4, 3, 0)
	x := oj.R().C("a").W().C("b")
	result := x.First(data)
	i, _ := result.(int)
	tt.Equal(t, 1, i/100)
	tt.Equal(t, 2, i%10)

	x = oj.R().C("b").N(1).C("c")
	result = x.First(data)
	tt.Equal(t, 223, result)
}

func TestOjExprFirstNode(t *testing.T) {
	data := buildNodeTree(4, 3, 0)
	x := oj.R().C("a").W().C("b")
	result := x.FirstNode(data)
	i, _ := result.(oj.Int)
	tt.Equal(t, 1, int(i)/100)
	tt.Equal(t, 2, int(i)%10)

	x = oj.R().C("b").N(1).C("c")
	result = x.FirstNode(data)
	tt.Equal(t, 223, result)
}

func xTestOjExprDev(t *testing.T) {
	data := buildTree(4, 3, 0)
	/*
		data := map[string]interface{}{
			"a": 1,
			"b": map[string]interface{}{
				"x": 2,
			},
			"c": 3,
		}
	*/
	x := oj.D().N(1).C("b")

	result := x.First(data)
	fmt.Printf("*** data: %s\n", oj.JSON(data, 2))
	fmt.Printf("*** results: %s\n", oj.JSON(result))
}