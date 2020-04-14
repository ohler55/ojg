// Copyright (c) 2020, Peter Ohler, All rights reserved.

package jp

import (
	"strings"

	"github.com/ohler55/ojg/gd"
)

// Slice
type Slice struct {
	Start int
	End   int
	Step  int
}

func (f *Slice) fill(b *strings.Builder) {
	// TBD
}

func (f *Slice) bracketFill(b *strings.Builder) {
	// TBD
}

func (f *Slice) get(n interface{}, rest []Frag) (result []gd.Node) {
	if 0 < len(rest) {
		// TBD slice
	}
	return
}

func (f *Slice) first(n interface{}, rest []Frag) (result gd.Node) {
	if 0 < len(rest) {
		// TBD slice
	}
	return
}

func (f *Slice) set(n, value interface{}, rest []Frag) error {
	if 0 < len(rest) {
		// TBD slice
	}
	return nil
}

func (f *Slice) setOne(n, value interface{}, rest []Frag) error {
	if 0 < len(rest) {
		// TBD slice
	}
	return nil
}

func (f *Slice) remove(n interface{}, rest []Frag) (removed []interface{}) {
	if 0 < len(rest) {
		// TBD slice
	}
	return
}

func (f *Slice) removeOne(n interface{}, rest []Frag) (removed interface{}) {
	if 0 < len(rest) {
		// TBD slice
	}
	return
}
