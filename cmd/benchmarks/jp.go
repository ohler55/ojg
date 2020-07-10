// Copyright (c) 2020, Peter Ohler, All rights reserved.

package main

import (
	"testing"

	"github.com/ohler55/ojg/jp"
)

func jpGet(b *testing.B) {
	p := jp.R().D().C("a").N(2).C("c")
	data := buildTree(10, 4, 0)
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		_ = p.Get(data)
		//x := p.Get(data)
		//fmt.Printf("*** %s\n", oj.JSON(x))
	}
}

func jpFirst(b *testing.B) {
	p := jp.R().D().C("a").N(2).C("c")
	data := buildTree(10, 4, 0)
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		_ = p.First(data)
		//fmt.Printf("*** %v\n", z)
	}
}
