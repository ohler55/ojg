// Copyright (c) 2020, Peter Ohler, All rights reserved.

package jp

import (
	"strings"

	"github.com/ohler55/ojg/gen"
)

// Nth
type Nth int

func (f *Nth) fill(b *strings.Builder) {
	// TBD "[n]"
}

func (f *Nth) bracketFill(b *strings.Builder) {
	// TBD "[n]"
}

func (f *Nth) get(n interface{}, rest []Frag) (result []gen.Node) {
	if 0 < len(rest) {
		// TBD match all
	}
	return
}

func (f *Nth) first(n interface{}, rest []Frag) (result gen.Node) {
	if 0 < len(rest) {
		// TBD match all
	}
	return
}

func (f *Nth) set(n, value interface{}, rest []Frag) error {
	if 0 < len(rest) {
		// TBD match all
	}
	return nil
}

func (f *Nth) setOne(n, value interface{}, rest []Frag) error {
	if 0 < len(rest) {
		// TBD match all
	}
	return nil
}

func (f *Nth) remove(n interface{}, rest []Frag) (removed []interface{}) {
	if 0 < len(rest) {
		// TBD match all
	}
	return
}

func (f *Nth) removeOne(n interface{}, rest []Frag) (removed interface{}) {
	if 0 < len(rest) {
		// TBD match all
	}
	return
}
