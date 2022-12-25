// Copyright (c) 2020, Peter Ohler, All rights reserved.

package alt_test

import (
	"fmt"

	"github.com/ohler55/ojg/alt"
)

type Setter struct {
	a int64
	b string
}

func (s *Setter) String() string {
	return fmt.Sprintf("Setter{a:%d,b:%s}", s.a, s.b)
}

func (s *Setter) SetAttr(attr string, val any) error {
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

func ExampleAttrSetter() {
	src := map[string]any{"a": 3, "b": "bee"}
	r := alt.MustNewRecomposer("", nil)
	var setter Setter
	_ = r.MustRecompose(src, &setter)
	fmt.Println(setter.String())

	// Output:
	// Setter{a:3,b:bee}
}
