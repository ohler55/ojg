// Copyright (c) 2020, Peter Ohler, All rights reserved.

package ojg

// ParseOptions for parsing JSON.
type ParseOptions struct {

	// NoComments returns an error if a comment is encountered.
	NoComment bool

	// OnlyOne returns an error if more than one JSON is in the string or
	// stream.
	OnlyOne bool

	// Handler for parsing.
	Handler interface{}

	Error       func(handler interface{}, err error, line, col int64)
	ObjectStart func(handler interface{})
	ObjectEnd   func(handler interface{})
	ArrayStart  func(handler interface{})
	ArrayEnd    func(handler interface{})
	Null        func(handler interface{})
	Int         func(handler interface{}, value int64)
	Float       func(handler interface{}, value float64)
	Bool        func(handler interface{}, value bool)
	Str         func(handler interface{}, key string)
	Key         func(handler interface{}, key string)
	Call        func(handler interface{})
}
