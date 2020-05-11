// Copyright (c) 2020, Peter Ohler, All rights reserved.

/*

Package simple contains functions and types to support building simple types
where simple types are:

  nil
  bool
  int64
  float64
  string
  time.Time
  []interface{}
  map[string]interface{}

Supporting functionality include decomposing structs and recomposing them. A
builder is also included that allows a simple way to build complex data using
a stack based model.

Builder

An example of building simple data is:

  var b simple.Builder

  b.Object()
  b.Value(1, "a")
  b.Array("b")
  b.Value(2)
  b.Pop()
  b.Pop()

  // creates map[string]interface{}{"a": 1, "b": []interface{}{2}}

Decompose and Recompose

// TBD

*/
package simple
