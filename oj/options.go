// Copyright (c) 2020, Peter Ohler, All rights reserved.

package oj

import (
	"io"
)

// Options for writing data to JSON.
type Options struct {

	// Indent for the output.
	Indent int

	// Sort object members if true.
	Sort bool

	// OmitNil skips the writing of nil values in an object.
	OmitNil bool

	// InitSize is the initial buffer size.
	InitSize int

	// WriteLimit is the size of the buffer that will trigger a write when
	// using a writer.
	WriteLimit int

	// TimeFormat defines how time is encoded. Options are to use a time. layout
	// string format such as time.RFC3339Nano, "second" for a decimal
	// representation, "nano" for a an integer.
	TimeFormat string

	// TimeWrap if not empty encoded time as an object with a single member. For
	// example if set to "@" then and TimeFormat is RFC3339Nano then the encoded
	// time will look like '{"@":"2020-04-12T16:34:04.123456789Z"}'
	TimeWrap string

	// CreateKey if set is the key to use when encoding objects that can later
	// be reconstituted with an Unmarshall call. This is only use when writing
	// simple types where one of the object in an array or map is not a
	// Simplifier. Reflection is used to encode all public members of the
	// object if possible. For example, is CreateKey is set to "type" this
	// might be the encoding.
	//
	//   { "type": "MyType", "a": 3, "b": true }
	//
	CreateKey string

	// FullTypePath if true includes the full type name and path when used
	// with the CreateKey.
	FullTypePath bool

	buf []byte
	w   io.Writer
}

var defaultOptions = Options{
	InitSize: 256,
	buf:      make([]byte, 0, 256),
}
