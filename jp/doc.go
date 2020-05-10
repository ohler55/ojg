// Copyright (c) 2020, Peter Ohler, All rights reserved.

/* Package jp is an implementation of JSON Path.

The JSONPath expressions operate on both generic data from the ojg/gen package
as well as on interface{} based data. interface{} data can include:

  nil
  bool
  int64 (or other sizes of int and uint)
  float64 and float32
  string
  time.Time
  []interface{}
  map[string]interface{}

Separate but equivalent functions are provided for gen and interface{}
navigation. The reasoning is that the gen based functions are more strongly
types so issues are caught at compile time while interface{} issues can only
be caught at run time. The interface{} functions are prefixed with an `I`.

*/
package jp
