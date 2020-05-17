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
	obj := map[string]interface{}{}
	for i := 0; i < size; i++ {
		k := string([]byte{'a' + byte(i)})
		nv := iv*10 + i + 1
		if 0 < depth {
			obj[k] = buildTree(size, depth-1, nv)
		} else {
			obj[k] = nv
		}
	}
	return obj
}

func TestOjExprBuild(t *testing.T) {
	x := oj.X().D().C("abc").W()
	tt.Equal(t, "..abc.*", x.String())

	x = oj.R().Descent().Child("abc").Wildcard()
	tt.Equal(t, "$..abc.*", x.String())

	x = oj.B().Descent().Child("abc").Wildcard()
	tt.Equal(t, "[..]['abc'][*]", x.String())

	x = oj.R().B().Descent().Child("abc").Wildcard()
	tt.Equal(t, "$[..]['abc'][*]", x.String())
}

func TestOjExprGet(t *testing.T) {
	data := buildTree(4, 2, 0)
	x := oj.R().C("a").W().C("b")
	result := x.Get(data)
	sort.Slice(result, func(i, j int) bool {
		iv, _ := result[i].(int)
		jv, _ := result[j].(int)
		return iv < jv
	})
	tt.Equal(t, []interface{}{112, 122, 132, 142}, result)

	x = oj.R().D().C("b").W().C("c")
	result = x.Get(data)
	sort.Slice(result, func(i, j int) bool {
		iv, _ := result[i].(int)
		jv, _ := result[j].(int)
		return iv < jv
	})
	fmt.Printf("*** %s\n", oj.JSON(result, 2))
	tt.Equal(t, []interface{}{213, 223, 233, 243}, result)

	x = oj.X().D().C("a").W().C("c").C("d")
	data = buildTree(4, 3, 0)
	//fmt.Printf("*** %s\n", oj.JSON(data, 2))
	result = x.Get(data)
	fmt.Printf("*** %s\n", oj.JSON(result, 2))

}
