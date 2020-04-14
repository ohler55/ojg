// Copyright (c) 2020, Peter Ohler, All rights reserved.

package jp

import (
	"strings"

	"github.com/ohler55/ojg/gd"
)

// Bracketed is an indicator that the expression should be written in bracketed format.
type Bracketed int

func (f *Bracketed) fill(b *strings.Builder) {
	// dummy
}

func (f *Bracketed) bracketFill(b *strings.Builder) {
	// dummy
}

func (f *Bracketed) get(n interface{}, rest []Frag) (result []gd.Node) {
	if 0 < len(rest) {
		result = rest[0].get(n, rest[1:])
	}
	return
}

func (f *Bracketed) first(n interface{}, rest []Frag) (result gd.Node) {
	if 0 < len(rest) {
		result = rest[0].first(n, rest[1:])
	}
	return
}

func (f *Bracketed) set(n, value interface{}, rest []Frag) error {
	if 0 < len(rest) {
		return rest[0].set(n, value, rest[1:])
	}
	return nil
}

func (f *Bracketed) setOne(n, value interface{}, rest []Frag) error {
	if 0 < len(rest) {
		return rest[0].setOne(n, value, rest[1:])
	}
	return nil
}

func (f *Bracketed) remove(n interface{}, rest []Frag) (removed []interface{}) {
	if 0 < len(rest) {
		removed = rest[0].remove(n, rest[1:])
	}
	return
}

func (f *Bracketed) removeOne(n interface{}, rest []Frag) (removed interface{}) {
	if 0 < len(rest) {
		removed = rest[0].removeOne(n, rest[1:])
	}
	return
}
