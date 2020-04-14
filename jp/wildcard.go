// Copyright (c) 2020, Peter Ohler, All rights reserved.

package jp

import (
	"strings"

	"github.com/ohler55/ojg/gd"
)

// Wildcard
type Wildcard int

func (f *Wildcard) fill(b *strings.Builder) {
	// TBD "*" or is it ".*"
}

func (f *Wildcard) bracketFill(b *strings.Builder) {
	// TBD "[*]"
}

func (f *Wildcard) get(n interface{}, rest []Frag) (result []gd.Node) {
	if 0 < len(rest) {
		// TBD match all
	}
	return
}

func (f *Wildcard) first(n interface{}, rest []Frag) (result gd.Node) {
	if 0 < len(rest) {
		// TBD match all
	}
	return
}

func (f *Wildcard) set(n, value interface{}, rest []Frag) error {
	if 0 < len(rest) {
		// TBD match all
	}
	return nil
}

func (f *Wildcard) setOne(n, value interface{}, rest []Frag) error {
	if 0 < len(rest) {
		// TBD match all
	}
	return nil
}

func (f *Wildcard) remove(n interface{}, rest []Frag) (removed []interface{}) {
	if 0 < len(rest) {
		// TBD match all
	}
	return
}

func (f *Wildcard) removeOne(n interface{}, rest []Frag) (removed interface{}) {
	if 0 < len(rest) {
		// TBD match all
	}
	return
}
