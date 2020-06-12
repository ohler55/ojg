// Copyright (c) 2020, Peter Ohler, All rights reserved.

package alt_test

import (
	"fmt"

	"github.com/ohler55/ojg/alt"
	"github.com/ohler55/ojg/gen"
	"github.com/ohler55/ojg/oj"
)

func ExampleDecompose() {
	type Sample struct {
		Int int
		Str string
	}
	sample := Sample{Int: 3, Str: "three"}
	simple := alt.Decompose(&sample, &alt.Options{CreateKey: "^", FullTypePath: true})

	fmt.Println(oj.JSON(simple, &oj.Options{Sort: true}))

	// Output: {"^":"github.com/ohler55/ojg/alt_test/Sample","int":3,"str":"three"}
}

func ExampleRecomposer() {
	type Sample struct {
		Int int
		Str string
	}
	r, err := alt.NewRecomposer("^", map[interface{}]alt.RecomposeFunc{&Sample{}: nil})
	var v interface{}
	if err == nil {
		v, err = r.Recompose(map[string]interface{}{"^": "Sample", "int": 3, "str": "three"})
	}
	if err == nil {
		fmt.Printf("type: %T\n", v)
		if sample, _ := v.(*Sample); sample != nil {
			fmt.Printf("sample: {Int: %d, Str: %q}\n", sample.Int, sample.Str)
		}
	} else {
		fmt.Println(err.Error())
	}
	// Output:
	// type: *alt_test.Sample
	// sample: {Int: 3, Str: "three"}
}

func ExampleInt() {
	for _, src := range []interface{}{1, "1", "x", 1.5, []interface{}{}} {
		fmt.Printf("alt.Int(%T(%v)) = %d  alt.Int(%T(%v), 2) = %d   alt.Int(%T(%v), 2, 3) = %d\n",
			src, src, alt.Int(src),
			src, src, alt.Int(src, 2),
			src, src, alt.Int(src, 2, 3))
	}
	// Output:
	// alt.Int(int(1)) = 1  alt.Int(int(1), 2) = 1   alt.Int(int(1), 2, 3) = 1
	// alt.Int(string(1)) = 1  alt.Int(string(1), 2) = 1   alt.Int(string(1), 2, 3) = 3
	// alt.Int(string(x)) = 0  alt.Int(string(x), 2) = 2   alt.Int(string(x), 2, 3) = 3
	// alt.Int(float64(1.5)) = 1  alt.Int(float64(1.5), 2) = 1   alt.Int(float64(1.5), 2, 3) = 3
	// alt.Int([]interface {}([])) = 0  alt.Int([]interface {}([]), 2) = 2   alt.Int([]interface {}([]), 2, 3) = 2
}

type Genny struct {
	val int
}

func (g *Genny) Generic() gen.Node {
	return gen.Object{"type": gen.String("Genny"), "val": gen.Int(g.val)}
}

func ExampleGenerify() {
	// type Genny struct {
	// 	val int
	// }
	//
	// func (g *Genny) Generic() gen.Node {
	// 	return gen.Object{"type": gen.String("genny"), "val": gen.Int(g.val)}
	// }
	ga := []*Genny{&Genny{val: 3}}
	v := alt.Generify(ga)
	fmt.Println(oj.JSON(v, &oj.Options{Sort: true}))

	// Output: [{"type":"Genny","val":3}]
}

func ExampleAlter() {
	m := map[string]interface{}{"a": 1, "b": 4, "c": 9}
	v := alt.GenAlter(m)
	fmt.Println(oj.JSON(v, &oj.Options{Sort: true}))

	// Output: {"a":1,"b":4,"c":9}
}
