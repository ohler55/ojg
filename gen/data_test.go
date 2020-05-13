// Copyright (c) 2020, Peter Ohler, All rights reserved.

package gen_test

import "github.com/ohler55/ojg/oj"

type data struct {
	src string
	// Empty means no error expected while non empty should be compared
	// err.Error().
	expect    string
	value     interface{}
	onlyOne   bool
	noComment bool
	options   *oj.Options
	indent    int
}
