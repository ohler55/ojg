// Copyright (c) 2020, Peter Ohler, All rights reserved.

package alt

import "github.com/ohler55/ojg"

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

	// UseTags if true will use the json annotation tags when marhsalling,
	// writing, or decomposing an struct. If no tag is present then the
	// KeyExact flag is referenced to determine the key.
	UseTags bool

	// KeyExact if true will use the exact field name for an encoded struct
	// field. If false the key style most often seen in JSON files where the
	// first character of the object keys is lowercase.
	KeyExact bool

	// NestEmbed if true will generate an element for each anonymous embedded
	// field.
	NestEmbed bool

	// Converter to use when decomposing or altering if non nil.
	Converter *Converter

	// BytesAs indicates how []byte fields should be encoded. Choices are
	// BytesAsString, BytesAsBase64 (the go json package default), or
	// BytesAsArray.
	BytesAs int
}

// DefaultOptions are the default options for decompsing.
var DefaultOptions = Options{
	CreateKey:    "type",
	FullTypePath: false,
	OmitNil:      true,
	UseTags:      false,
	Converter:    nil,
}

// GoOptions are the options closest to the go json package if it were
// decomposing.
var GoOptions = Options{
	CreateKey:    "",
	FullTypePath: false,
	OmitNil:      false,
	UseTags:      true,
	KeyExact:     true,
	NestEmbed:    false,
	Converter:    nil,
	BytesAs:      ojg.BytesAsBase64,
}
