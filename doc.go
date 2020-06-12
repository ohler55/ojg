// Copyright (c) 2020, Peter Ohler, All rights reserved.

/*

Package ojg is a collection of JSON tools including a validator and parser.

Oj

Package oj contains functions and types for parsing JSON as well as support
for building building simple types. Included in the oj package are:

  Parser for parsing JSON strings and streams into simple types.

  Validator for validating JSON strings and streams.

  Builder for building simple types.

  Writer for writing data as JSON.

Gen

- node, compiler checked values
- parser
- builder

Jp

- JSONPath
- get
- set
- reflection

Alt

The alt package contains functions and types for altering values. It includes functions for:

  Decompose() a value into simple types of bool, int64, float64, string,
              time.Time, []interface{} and map[string]interface{}.

  Recompose() takes simple data type and converts it back into a complex type.

  Alter() is the same as decompose except it alters the value in place.

  Generify() converts a simple value into a gen.Node.

Cmd oj

- parse and reformat
- extract
- filter

*/
package ojg
