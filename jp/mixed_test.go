// Copyright (c) 2026, Peter Ohler, All rights reserved.

package jp_test

import (
	"fmt"
	"sort"
	"strings"
	"testing"

	"github.com/ohler55/ojg/alt"
	"github.com/ohler55/ojg/gen"
	"github.com/ohler55/ojg/jp"
	"github.com/ohler55/ojg/oj"
	"github.com/ohler55/ojg/sen"
	"github.com/ohler55/ojg/tt"
)

const mixedJSON = `{"a":[10,11,12,13,14,15,16,17],"o":{"x":1,"y":2},"n":[{"v":1},{"v":2},{"v":3}],"m":15}`

// The expected results are from RFC 9535 section 2.5.1 and the JSONPath
// Compliance Test Suite: the values selected by each member are joined in
// member order and duplicates are kept.
var mixedTestData = []struct {
	path     string
	str      string
	expect   string
	anyOrder bool // object member order is not fixed
}{
	{path: "$.a[1,5:7]", str: "$.a[1,5:7]", expect: "[11 15 16]"},
	{path: "$.a[0:2,5]", str: "$.a[0:2,5]", expect: "[10 11 15]"},
	{path: "$.a[0:3,2:4]", str: "$.a[0:3,2:4]", expect: "[10 11 12 12 13]"},
	{path: "$.a[ 6 : , :1 ]", str: "$.a[6:,:1]", expect: "[16 17 10]"},
	{path: "$.a[::-3,-1]", str: "$.a[::-3,-1]", expect: "[17 14 11 17]"},
	{path: "$.a[*,1]", str: "$.a[*,1]", expect: "[10 11 12 13 14 15 16 17 11]"},
	{path: "$.a[1,*]", str: "$.a[1,*]", expect: "[11 10 11 12 13 14 15 16 17]"},
	{path: "$.a[?@>15,0]", str: "$.a[?(@ > 15),0]", expect: "[16 17 10]"},
	{path: "$.a[0,?@>15]", str: "$.a[0,?(@ > 15)]", expect: "[10 16 17]"},
	{path: "$.a[?@<11,?@>16]", str: "$.a[?(@ < 11),?(@ > 16)]", expect: "[10 17]"},
	{path: "$.a[?(@ > 15), 1:3]", str: "$.a[?(@ > 15),1:3]", expect: "[16 17 11 12]"},
	{path: "$.a[?@>$.m,-1]", str: "$.a[?(@ > $.m),-1]", expect: "[16 17 17]"},
	{path: "$.o['x',*]", str: "$.o['x',*]", expect: "[1 1 2]", anyOrder: true},
	{path: "$['a','o'][0,'y']", str: "$['a','o'][0,'y']", expect: "[10 2]"},
	{path: "$.n[0:2,?@.v==3].v", str: "$.n[0:2,?(@.v == 3)].v", expect: "[1 2 3]"},
	{path: "$..[1:2,'v']", str: "$..[1:2,'v']", expect: "[1 11 2 3 {v:2}]", anyOrder: true},
	{path: "$.o[1:2,'x']", str: "$.o[1:2,'x']", expect: "[1]"},
	{path: "$.a['x',?@.v]", str: "$.a['x',?(@.v)]", expect: "[]"},
}

func sortedString(list []any) string {
	strs := make([]string, len(list))
	for i, v := range list {
		strs[i] = sen.String(v)
	}
	sort.Strings(strs)
	return "[" + strings.Join(strs, " ") + "]"
}

func TestMixedUnionGet(t *testing.T) {
	data := oj.MustParseString(mixedJSON)
	for i, d := range mixedTestData {
		x := jp.MustParseString(d.path)
		tt.Equal(t, d.str, x.String(), "%d: %s string", i, d.path)
		// The string form must parse back to the same expression.
		tt.Equal(t, d.str, jp.MustParseString(x.String()).String(), "%d: %s reparse", i, d.path)

		got := x.Get(data)
		if d.anyOrder {
			tt.Equal(t, d.expect, sortedString(got), "%d: %s get", i, d.path)
			continue
		}
		tt.Equal(t, d.expect, sen.String(got), "%d: %s get", i, d.path)

		var first any
		if 0 < len(got) {
			first = got[0]
		}
		tt.Equal(t, first, x.First(data), "%d: %s first", i, d.path)
		tt.Equal(t, 0 < len(got), x.Has(data), "%d: %s has", i, d.path)

		n := alt.Generify(data)
		nodes := x.GetNodes(n)
		tt.Equal(t, d.expect, sen.String(nodes), "%d: %s get nodes", i, d.path)
		var firstNode gen.Node
		if 0 < len(nodes) {
			firstNode = nodes[0]
		}
		tt.Equal(t, firstNode, x.FirstNode(n), "%d: %s first node", i, d.path)
	}
}

func TestMixedUnionParse(t *testing.T) {
	x := jp.MustParseString("$[1,2]")
	_, ok := x[1].(jp.Union)
	tt.Equal(t, true, ok, "a name and index list stays a Union")

	x = jp.MustParseString("$[1,2:3,*,?@.x,'y']")
	tt.Equal(t, "jp.MixedUnion", fmt.Sprintf("%T", x[1]))
	tt.Equal(t, 5, len(x[1].(jp.MixedUnion)))

	for _, d := range []struct {
		src string
		err string
	}{
		{src: "$[1,2:3", err: "expected a number at 8 in $[1,2:3"},
		{src: "$[*,", err: "not terminated at 5 in $[*,"},
		{src: "$[*x]", err: "not terminated at 5 in $[*x]"},
		{src: "$[1:2 3,4]", err: "invalid slice syntax at 8 in $[1:2 3,4]"},
		{src: "$[1,2:3 4]", err: "invalid slice syntax at 10 in $[1,2:3 4]"},
		{src: "$[?@.x,", err: "not terminated at 8 in $[?@.x,"},
		{src: "$[1,(@.x)]", err: "invalid union syntax at 6 in $[1,(@.x)]"},
		{src: "$[?1,2]", err: "a literal must be used in a comparison at 5 in $[?1,2]"},
	} {
		_, err := jp.ParseString(d.src)
		tt.NotNil(t, err, d.src)
		tt.Equal(t, d.err, err.Error(), d.src)
	}
}

func TestMixedUnionNew(t *testing.T) {
	f := jp.MustNewFilter("[?@.x == 1]")
	x := append(jp.R().C("a"), jp.NewMixedUnion("b", 1, int64(-1), jp.Slice{1, 3, 1}, jp.Wildcard('*'), f, 2.5))
	tt.Equal(t, "$.a['b',1,-1,1:3:1,*,?(@.x == 1)]", x.String())
}

func TestMixedUnionLocate(t *testing.T) {
	data := oj.MustParseString(mixedJSON)
	for _, d := range []struct {
		path   string
		max    int
		expect string
	}{
		{path: "$.a[1,5:7]", expect: "[$.a[1] $.a[5] $.a[6]]"},
		{path: "$.a[1,5:7]", max: 2, expect: "[$.a[1] $.a[5]]"},
		{path: "$.a[6:,0]", max: 3, expect: "[$.a[6] $.a[7] $.a[0]]"},
		{path: "$.o['y',*]", max: 1, expect: "[$.o.y]"},
		{path: "$.n[2,0:1].v", expect: "[$.n[2].v $.n[0].v]"},
	} {
		locs := jp.MustParseString(d.path).Locate(data, d.max)
		tt.Equal(t, d.expect, fmt.Sprintf("%v", locs), "%s max %d", d.path, d.max)
	}
}

func TestMixedUnionWalk(t *testing.T) {
	data := oj.MustParseString(mixedJSON)
	var paths []string
	jp.MustParseString("$.n[2,0:2].v").Walk(data, func(path jp.Expr, nodes []any) {
		paths = append(paths, fmt.Sprintf("%s=%v", path, nodes[len(nodes)-1]))
	})
	tt.Equal(t, "[n[2].v=3 n[0].v=1 n[1].v=2]", fmt.Sprintf("%v", paths))
}

func TestMixedUnionPathMatch(t *testing.T) {
	target := jp.MustParseString("$.a[1,3:5,?@.x].b")
	tt.Equal(t, true, jp.PathMatch(target, jp.MustParseString("$.a[1].b")))
	tt.Equal(t, true, jp.PathMatch(target, jp.MustParseString("$.a[7].b")), "a filter matches any path")
	tt.Equal(t, false, jp.PathMatch(target, jp.MustParseString("$.a[1].c")))

	target = jp.MustParseString("$.a['x',2:3]")
	tt.Equal(t, true, jp.PathMatch(target, jp.MustParseString("$.a.x")))
	tt.Equal(t, true, jp.PathMatch(target, jp.MustParseString("$.a[0]")), "a slice matches any index")
	tt.Equal(t, false, jp.PathMatch(target, jp.MustParseString("$.a.y")))
	tt.Equal(t, true, jp.PathMatch(jp.MustParseString("$['x',*]"), jp.MustParseString("$.y")))
}

func TestMixedUnionSet(t *testing.T) {
	data := oj.MustParseString(mixedJSON)
	x := jp.MustParseString("$.n[0,?@.v>2].v")
	tt.Nil(t, x.Set(data, int64(9)))
	tt.Equal(t, "[{v:9}{v:2}{v:9}]", sen.String(jp.C("n").Get(data)[0]))

	err := jp.MustParseString("$.a[1,2:3]").Set(data, 0)
	tt.NotNil(t, err)
	tt.Equal(t, "can not set with an expression ending with a MixedUnion", err.Error())

	_, err = jp.MustParseString("$.a[1,2:3]").Modify(data, func(v any) (any, bool) { return v, false })
	tt.NotNil(t, err)

	_, err = jp.MustParseString("$.a[1,2:3]").Remove(data)
	tt.NotNil(t, err)

	result, err := jp.MustParseString("$.n[0:1,2].v").Modify(data, func(v any) (any, bool) {
		return v.(int64) + 1, true
	})
	tt.Nil(t, err)
	tt.Equal(t, "[{v:10}{v:2}{v:10}]", sen.String(jp.C("n").Get(result)[0]))
}
