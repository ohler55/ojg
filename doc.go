// Copyright (c) 2020, Peter Ohler, All rights reserved.

/*

package ojg is a collection of JSON tool including a validator and parser. It
supports parsing into either simple types which are nil, bool, int64, float64,
string, []interface{}, and map[string]interface{} or parsing into general data
types which are used to enforce type safety on JSON compatible set of types.


oj

- parser
- builder
- value as simple type
- validate

gen

- node, compiler checked values
- parser
- builder

jp

- JSONPath
- get
- set
- reflection

conv

- decompose
- recompose
- alter
- generify

cmd/oj

- parse and reformat
- extract
- filter

*/
package ojg
