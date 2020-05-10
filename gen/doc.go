// Copyright (c) 2020, Peter Ohler, All rights reserved.

// Package gen defines OjG Generic types which is a thin wrapper around native
// or simple go types that can be represented by JSON. In addition to the JSON
// types a Time type is included that can be written in several forms. The
// types in this package enforce the type safety of the datain contrast to the
// simple interface{} based types which can hold any values whether that can
// be converted to JSON or not.
//
// The package also includes support for converting from go type to the
// generic types.
package gen
