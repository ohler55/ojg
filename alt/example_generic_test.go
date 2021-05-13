// Copyright (c) 2021, Peter Ohler, All rights reserved.

package alt_test

import (
	"fmt"

	"github.com/ohler55/ojg/alt"
	"github.com/ohler55/ojg/gen"
	"github.com/ohler55/ojg/oj"
)

type genny struct {
	val int
}

func (g *genny) Generic() gen.Node {
	return gen.Object{"type": gen.String("genny"), "val": gen.Int(g.val)}
}

func ExampleGenerify() {
	// Non public types can be encoded with the Genericer interface which
	// should decompose into a gen.Node.
	ga := []*genny{&genny{val: 3}}
	v := alt.Generify(ga)
	// Encode to JSON after decomposing using the Genericer interface.
	fmt.Println(oj.JSON(v, &oj.Options{Sort: true}))

	// Output: [{"type":"genny","val":3}]
}
