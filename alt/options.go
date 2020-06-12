// Copyright (c) 2020, Peter Ohler, All rights reserved.

package alt

// Options are the options available to Decompose() function.
type Options struct {

	// CreateKey is the map element used to identify the type of a decomposed
	// object.
	CreateKey string

	// FullTypePath if true will use the package and type name as the
	// CreateKey value.
	FullTypePath bool

	// OmitNil if true omits object members that have nil values.
	OmitNil bool
}

// DefaultOptions are the default options for decompsing.
var DefaultOptions = Options{
	CreateKey:    "type",
	FullTypePath: false,
	OmitNil:      true,
}
