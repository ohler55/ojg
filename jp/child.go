// Copyright (c) 2020, Peter Ohler, All rights reserved.

package jp

import (
	"strings"

	"github.com/ohler55/ojg/gen"
)

// Child
type Child string

func (f *Child) fill(b *strings.Builder) {
	// TBD "." + string(f)
}

func (f *Child) bracketFill(b *strings.Builder) {
	// TBD "['" + string(f) + "']"
}

func (f *Child) get(n interface{}, rest []Frag) (result []gen.Node) {
	if 0 < len(rest) {
		// TBD match name is a map then rest
	}
	return
}

func (f *Child) first(n interface{}, rest []Frag) (result gen.Node) {
	if 0 < len(rest) {
		// TBD match name is a map then rest
	}
	return
}

func (f *Child) set(n, value interface{}, rest []Frag) error {
	if 0 < len(rest) {
		// TBD match name is a map then rest
	}
	return nil
}

func (f *Child) setOne(n, value interface{}, rest []Frag) error {
	if 0 < len(rest) {
		// TBD match name is a map then rest
	}
	return nil
}

func (f *Child) remove(n interface{}, rest []Frag) (removed []interface{}) {
	if 0 < len(rest) {
		// TBD match name is a map then rest
	}
	return
}

func (f *Child) removeOne(n interface{}, rest []Frag) (removed interface{}) {
	if 0 < len(rest) {
		// TBD match name is a map then rest
	}
	return
}
