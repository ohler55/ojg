// Copyright (c) 2020, Peter Ohler, All rights reserved.

package jp

import (
	"strings"

	"github.com/ohler55/ojg/gd"
)

// Union
type Union struct {
	// TBD [0,1]
}

func (f *Union) fill(b *strings.Builder) {
	// TBD
}

func (f *Union) bracketFill(b *strings.Builder) {
	// TBD
}

func (f *Union) get(n interface{}, rest []Frag) (result []gd.Node) {
	if 0 < len(rest) {
		// TBD union
	}
	return
}

func (f *Union) first(n interface{}, rest []Frag) (result gd.Node) {
	if 0 < len(rest) {
		// TBD union
	}
	return
}

func (f *Union) set(n, value interface{}, rest []Frag) error {
	if 0 < len(rest) {
		// TBD union
	}
	return nil
}

func (f *Union) setOne(n, value interface{}, rest []Frag) error {
	if 0 < len(rest) {
		// TBD union
	}
	return nil
}

func (f *Union) remove(n interface{}, rest []Frag) (removed []interface{}) {
	if 0 < len(rest) {
		// TBD union
	}
	return
}

func (f *Union) removeOne(n interface{}, rest []Frag) (removed interface{}) {
	if 0 < len(rest) {
		// TBD union
	}
	return
}
