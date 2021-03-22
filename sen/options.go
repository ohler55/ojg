// Copyright (c) 2020, Peter Ohler, All rights reserved.

package sen

import (
	"io"
	"unicode/utf8"
)

const (
	Normal        = "\x1b[m"
	Black         = "\x1b[30m"
	Red           = "\x1b[31m"
	Green         = "\x1b[32m"
	Yellow        = "\x1b[33m"
	Blue          = "\x1b[34m"
	Magenta       = "\x1b[35m"
	Cyan          = "\x1b[36m"
	White         = "\x1b[37m"
	Gray          = "\x1b[90m"
	BrightRed     = "\x1b[91m"
	BrightGreen   = "\x1b[92m"
	BrightYellow  = "\x1b[93m"
	BrightBlue    = "\x1b[94m"
	BrightMagenta = "\x1b[95m"
	BrightCyan    = "\x1b[96m"
	BrightWhite   = "\x1b[97m"
)

// Options for writing data to JSON.
type Options struct {

	// Indent for the output.
	Indent int

	// Tab if true will indent using tabs and ignore the Indent member.
	Tab bool

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

	// TimeMap if true will encode time as a map with a create key and a
	// 'value' member formatted according to the TimeFormat options.
	TimeMap bool

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

	// Color if true will colorize the output.
	Color bool

	// SyntaxColor is the color for syntax in the JSON output.
	SyntaxColor string

	// KeyColor is the color for a key in the JSON output.
	KeyColor string

	// NullColor is the color for a null in the JSON output.
	NullColor string

	// BoolColor is the color for a bool in the JSON output.
	BoolColor string

	// NumberColor is the color for a number in the JSON output.
	NumberColor string

	// StringColor is the color for a string in the JSON output.
	StringColor string

	// TimeColor is the color for a time.Time in the JSON output.
	TimeColor string

	// NoColor turns the color off.
	NoColor string

	// HTMLSafe if true turns on escaping of &, <, and >.
	HTMLSafe bool

	Buf []byte
	Utf []byte
	W   io.Writer
}

// DefaultOptions are common defaults.
var DefaultOptions = Options{
	InitSize:    256,
	SyntaxColor: Normal,
	KeyColor:    Blue,
	NullColor:   Red,
	BoolColor:   Yellow,
	NumberColor: Cyan,
	StringColor: Green,
	TimeColor:   Magenta,
	NoColor:     Normal,
	Buf:         make([]byte, 0, 256),
	Utf:         make([]byte, utf8.UTFMax),
}

// BrightOptions bright color options.
var BrightOptions = Options{
	InitSize:    256,
	SyntaxColor: Normal,
	KeyColor:    BrightBlue,
	NullColor:   BrightRed,
	BoolColor:   BrightYellow,
	NumberColor: BrightCyan,
	StringColor: BrightGreen,
	TimeColor:   BrightMagenta,
	NoColor:     Normal,
	Buf:         make([]byte, 0, 256),
	Utf:         make([]byte, utf8.UTFMax),
}

// HTMLOptions defines color options for generating colored HTML. The encoding
// is suitable for use in a <pre> element.
var HTMLOptions = Options{
	InitSize:    256,
	SyntaxColor: "<span>",
	KeyColor:    `<span style="color:#44f">`,
	NullColor:   `<span style="color:red">`,
	BoolColor:   `<span style="color:#a40">`,
	NumberColor: `<span style="color:#04a">`,
	StringColor: `<span style="color:green">`,
	TimeColor:   `<span style="color:#f0f">`,
	NoColor:     "</span>",
	Buf:         make([]byte, 0, 256),
	Utf:         make([]byte, utf8.UTFMax),
	HTMLSafe:    true,
}
