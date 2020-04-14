// Copyright (c) 2020, Peter Ohler, All rights reserved.

package jp

import (
	"strings"
)

// Current represents the current position. It is returns the current element
// when evaluated.
type Current Bracketed

func (f *Current) fill(b *strings.Builder) {
	b.WriteString("@")
}

func (f *Current) bracketFill(b *strings.Builder) {
	b.WriteString("@")
}
