// Copyright (c) 2020, Peter Ohler, All rights reserved.

package jp

import (
	"strings"
)

// Root represents the root position. It is returns the root element when
// evaluated. Note this can only be the first element in an expression making
// it effectively useless.
type Root Bracketed

func (f *Root) fill(b *strings.Builder) {
	b.WriteString("$")
}

func (f *Root) bracketFill(b *strings.Builder) {
	b.WriteString("$")
}
