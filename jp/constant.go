// Copyright (c) 2020, Peter Ohler, All rights reserved.

package jp

import (
	"fmt"
	"strings"

	"github.com/ohler55/ojg/gd"
)

// Constant
type Constant struct {
	Value gd.Node
}

func (f *Constant) fill(b *strings.Builder) {
	// TBD
}

func (f *Constant) bracketFill(b *strings.Builder) {
	// TBD
}

func (f *Constant) get(_ interface{}, _ []Frag) (result []gd.Node) {
	return append(result, f.Value)
}

func (f *Constant) first(_ interface{}, _ []Frag) gd.Node {
	return f.Value
}

func (f *Constant) set(_, _ interface{}, _ []Frag) error {
	return fmt.Errorf("can not set a constant path element")
}

func (f *Constant) setOne(_, _ interface{}, _ []Frag) error {
	return fmt.Errorf("can not set a constant path element")
}

func (f *Constant) remove(_ interface{}, _ []Frag) (removed []interface{}) {
	return
}

func (f *Constant) removeOne(_ interface{}, _ []Frag) (removed interface{}) {
	return
}
