// Copyright (c) 2021, Peter Ohler, All rights reserved.

package pretty_test

import (
	"fmt"
	"testing"

	"github.com/ohler55/ojg/pretty"
	"github.com/ohler55/ojg/sen"
	"github.com/ohler55/ojg/tt"
)

func TestSEN(t *testing.T) {
	p := sen.Parser{}
	val, err := p.Parse([]byte(`[true false [3 2 1] [1 2 3 [x y z []]]]`))
	tt.Nil(t, err)
	opt := sen.DefaultOptions
	opt.Color = true
	s := pretty.JSON(val, 32, &opt)

	fmt.Printf("*** %s\n", s)

	s = pretty.SEN(val, 60)

	fmt.Printf("*** %s\n", s)

}

func TestSEN2(t *testing.T) {
	p := sen.Parser{}
	val, err := p.Parse([]byte(`[true {abc: 123 def: true}]`))
	tt.Nil(t, err)
	opt := sen.DefaultOptions
	opt.Color = true
	s := pretty.JSON(val, 25, &opt)

	fmt.Printf("*** %s\n", s)
	//fmt.Printf("*** % x\n", s)

}
