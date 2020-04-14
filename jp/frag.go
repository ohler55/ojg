// Copyright (c) 2020, Peter Ohler, All rights reserved.

package jp

import (
	"strings"

	"github.com/ohler55/ojg/gd"
)

// Frag represents a JSONPath fragment. A JSONPath expression is composed of
// fragments (Frag) linked together to form a full path expression.
type Frag interface {
	fill(b *strings.Builder)
	bracketFill(b *strings.Builder)

	get(n interface{}, rest []Frag) []gd.Node
	first(n interface{}, rest []Frag) gd.Node

	set(n, value interface{}, rest []Frag) error
	setOne(n, value interface{}, rest []Frag) error

	remove(n interface{}, rest []Frag) []interface{}
	removeOne(n interface{}, rest []Frag) interface{}
}
