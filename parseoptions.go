// Copyright (c) 2020, Peter Ohler, All rights reserved.

package ojg

// ParseOptions for parsing JSON.
type ParseOptions struct {

	// NoComments returns an error if a comment is encountered.
	NoComment bool

	// OnlyOne returns an error if more than one JSON is in the string or
	// stream.
	OnlyOne bool
}
