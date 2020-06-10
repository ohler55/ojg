// Copyright (c) 2020, Peter Ohler, All rights reserved.

/*

Package alt contains functions and type for altering values.

Conversions

Simple conversion from one to to another include converting to string, bool,
int64, float64, and time.Time. Each of these functions takes between one and
three arguments. The first is the value to convert. The second argument is the
value to return if the value can not be converted. For example, if the value
is an array then the second argument, the first default would be returned. If
the third argument is present then any input that is not the correct type will
cause the third default to be returned. The conversion functions are Int(),
FLoat(), Bool(), String(), and Time(). The reason for the defaults are to
allow a single return from a conversion unlike a type assertion.

  i := alt.Int("123", 0)

Generify

It is often useful to work with generic values that can be converted to JSON
and also provide type safety so that code can be checked at compile
time. Those value types are defined in the gen package. The Genericer
interface defines the Generic() function as

  Generic() gen.Node

A Generify() function is used to convert values to gen.Node types.

  // TBD example

Decompose

The Decompose() functions creates a simple type converting non simple to
simple types using either the Simplify() interface or reflection. Unlike
Alter() a deep copy is returned leaving the original data unchanged.

  // TBD decompose example with type

Recompose

Recompose simple data into more complex go types using either the Recompose()
function or the Recomposer struct that adds some efficiency by reusing
buffers.

  // TBD recompose something

Alter

The GenAlter() function converts a simple go data element into Node compliant
data. A best effort is made to convert values that are not simple into generic
Nodes. It modifies the values inplace if possible by altering the original.

  // TBD GenAlter(v interface{}, options ...*Options) (n gen.Node) {

*/
package alt
