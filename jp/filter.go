// Copyright (c) 2020, Peter Ohler, All rights reserved.

package jp

import (
	"strings"

	"github.com/ohler55/ojg/gd"
)

// Filter
type Filter struct {
	Left  Expr
	Right Expr
	Op    string
	// TBD handle && and ||
}

func (f *Filter) fill(b *strings.Builder) {
	// TBD
}

func (f *Filter) bracketFill(b *strings.Builder) {
	// TBD
}

func (f *Filter) get(n interface{}, rest []Frag) (result []gd.Node) {
	if 0 < len(rest) {
		// TBD filter
	}
	return
}

func (f *Filter) first(n interface{}, rest []Frag) (result gd.Node) {
	if 0 < len(rest) {
		// TBD filter
	}
	return
}

func (f *Filter) set(n, value interface{}, rest []Frag) error {
	if 0 < len(rest) {
		// TBD filter
	}
	return nil
}

func (f *Filter) setOne(n, value interface{}, rest []Frag) error {
	if 0 < len(rest) {
		// TBD filter
	}
	return nil
}

func (f *Filter) remove(n interface{}, rest []Frag) (removed []interface{}) {
	if 0 < len(rest) {
		// TBD filter
	}
	return
}

func (f *Filter) removeOne(n interface{}, rest []Frag) (removed interface{}) {
	if 0 < len(rest) {
		// TBD filter
	}
	return
}
