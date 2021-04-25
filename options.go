// Copyright (c) 2020, Peter Ohler, All rights reserved.

package ojg

import (
	"fmt"
	"strconv"
	"time"
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

	BytesAsString = iota
	BytesAsBase64
	BytesAsArray
)

var (
	// DefaultOptions default options that can be set as desired.
	DefaultOptions = Options{
		InitSize:    256,
		SyntaxColor: Normal,
		KeyColor:    Blue,
		NullColor:   Red,
		BoolColor:   Yellow,
		NumberColor: Cyan,
		StringColor: Green,
		TimeColor:   Magenta,
	}

	// BrightOptions encoding options for color encoding.
	BrightOptions = Options{
		InitSize:    256,
		SyntaxColor: Normal,
		KeyColor:    BrightBlue,
		NullColor:   BrightRed,
		BoolColor:   BrightYellow,
		NumberColor: BrightCyan,
		StringColor: BrightGreen,
		TimeColor:   BrightMagenta,
	}

	// GoOptions are the options closest to the go json package.
	GoOptions = Options{
		InitSize:     256,
		SyntaxColor:  Normal,
		KeyColor:     Blue,
		NullColor:    Red,
		BoolColor:    Yellow,
		NumberColor:  Cyan,
		StringColor:  Green,
		TimeColor:    Magenta,
		CreateKey:    "",
		FullTypePath: false,
		OmitNil:      false,
		UseTags:      true,
		KeyExact:     true,
		NestEmbed:    false,
		BytesAs:      BytesAsBase64,
	}
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

	// NoReflect if true does not use reflection to encode an object. This is
	// only considered if the CreateKey is empty.
	NoReflect bool

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

	// UseTags if true will use the json annotation tags when marhsalling,
	// writing, or decomposing an struct. If no tag is present then the
	// KeyExact flag is referenced to determine the key.
	UseTags bool

	// KeyExact if true will use the exact field name for an encoded struct
	// field. If false the key style most often seen in JSON files where the
	// first character of the object keys is lowercase.
	KeyExact bool

	// HTMLUnsafe if true turns off escaping of &, <, and >.
	HTMLUnsafe bool

	// NestEmbed if true will generate an element for each anonymous embedded
	// field.
	NestEmbed bool

	// BytesAs indicates how []byte fields should be encoded. Choices are
	// BytesAsString, BytesAsBase64 (the go json package default), or
	// BytesAsArray.
	BytesAs int
}

// BuildTime appends a time string to the buffer.
func (o *Options) BuildTime(buf []byte, t time.Time) []byte {
	if o.TimeMap {
		buf = append(buf, []byte(`{"`)...)
		buf = append(buf, o.CreateKey...)
		buf = append(buf, []byte(`":`)...)
		if o.FullTypePath {
			buf = append(buf, []byte(`"time/Time"`)...)
		} else {
			buf = append(buf, []byte("Time")...)
		}
		buf = append(buf, []byte(` value:`)...)
	} else if 0 < len(o.TimeWrap) {
		buf = append(buf, []byte(`{"`)...)
		buf = append(buf, []byte(o.TimeWrap)...)
		buf = append(buf, []byte(`":`)...)
	}
	switch o.TimeFormat {
	case "", "nano":
		buf = append(buf, []byte(strconv.FormatInt(t.UnixNano(), 10))...)
	case "second":
		// Decimal format but float is not accurate enough so build the output
		// in two parts.
		nano := t.UnixNano()
		secs := nano / int64(time.Second)
		if 0 < nano {
			buf = append(buf, []byte(fmt.Sprintf("%d.%09d", secs, nano-(secs*int64(time.Second))))...)
		} else {
			buf = append(buf, []byte(fmt.Sprintf("%d.%09d", secs, -(nano-(secs*int64(time.Second)))))...)
		}
	default:
		buf = append(buf, '"')
		buf = append(buf, []byte(t.Format(o.TimeFormat))...)
		buf = append(buf, '"')
	}
	if 0 < len(o.TimeWrap) || o.TimeMap {
		buf = append(buf, '}')
	}
	return buf
}
