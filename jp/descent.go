// Copyright (c) 2020, Peter Ohler, All rights reserved.

package jp

import (
	"strings"

	"github.com/ohler55/ojg/gen"
)

// Descent
type Descent int

func (f *Descent) fill(b *strings.Builder) {
	// TBD ".."
}

func (f *Descent) bracketFill(b *strings.Builder) {
	// TBD "[.]" or "[..]" ??
}

func (f *Descent) get(n interface{}, rest []Frag) (result []gen.Node) {
	if 0 < len(rest) {
		// TBD recursive descent
	}
	return
}

func (f *Descent) first(n interface{}, rest []Frag) (result gen.Node) {
	if 0 < len(rest) {
		// TBD recursive descent
	}
	return
}

func (f *Descent) set(n, value interface{}, rest []Frag) error {
	if 0 < len(rest) {
		// TBD recursive descent
	}
	return nil
}

func (f *Descent) setOne(n, value interface{}, rest []Frag) error {
	if 0 < len(rest) {
		// TBD recursive descent
	}
	return nil
}

func (f *Descent) remove(n interface{}, rest []Frag) (removed []interface{}) {
	if 0 < len(rest) {
		// TBD recursive descent
	}
	return
}

func (f *Descent) removeOne(n interface{}, rest []Frag) (removed interface{}) {
	if 0 < len(rest) {
		// TBD recursive descent
	}
	return
}
