# How I implemented high performance JSON parsers, JSONPath, and other tools

The JSON parser has gotten gotten faster over the years but was there
still some room for improvement? That was part of the challenge I
undertook. That wasn't the initial goal. The initial goal was to
implement a set of fast generic data that matched JSON. As part of
that effort a high performance JSONPath implementation was needed. The
JSONPath implementations for Go were limited and not at all
complete.

I suppose I knew all along a new JSON parser would be needed. Having
written two JSON parser before it didn't seem like a very daunting
task. Both the Ruby [OJ](https://github.com/ohler55/oj) and the C
parser [OjC](https://github.com/ohler55/ojc) are the best performers
in their respective languages. Why not an
[OjG](https://github.com/ohler55/ojg) for go.

## Requirements

The first requirement was that JSON parsing and any frequently used
operation such as JSONPath evaluation had to be fast over everything
else. With the luxury of not having to follow the existing Go json
package API the API could be designed for the best performance.

Additional requirements varied depending on the JSON tool.

### Generic Data

The generic data implementation had to support the basic JSON types of
`null`, boolean, number, string, array, and object. In addition time
should be supported. In both JSON use in Ruby and Go time has alway
been needed. Time is just too much apart of any set of data to leave
it out.

The generic data had to be type safe. It would not do to have an
element that could not be endcoded as JSON in the data.

A frequent operation for generic data is to store that data into a
JSON database or similar. That meant converting to simple Go types of
nil, bool, int64, float64, string, []interface{}, and
map[string]interface{} had to be fast.

### JSONPath

A JSONPath is used to extract elements from data. That part of the
implementation had to be fast. Parsing really didn't have to be fast
but it would be nice to have a way of building a JSONPath in a
performant manner even if it was not as convenient as parsing a
string.

The JSONPath implementaiton had implement all the features described
by [the Goessner
article](https://goessner.net/articles/JsonPath). There are other
descriptions of JSONPath but the Goessner description is the most
referenced. Since the implementation is in Go the scripting feature
described could be left out as long as similar functionality could be
provided for array indexes based on the length of the array. Borrowing
from Ruby, using negative indexes provided that functionality.

### JSON Parser and Validotor

A JSON parser and validator need not be the same and each should be as
performant as possible. The parsers needed to support parsing into
simple Go types as well as the generic data types.

When parsing files that include millions or more JSON elements that
might be over 100GB a streaming parser is necessary.

The parser must also allow parsing into Go types. Furthermore
interfaces must be supported. Go unmarshalling does not support
interface fields. Since many data types make use of interfaces that
limitation is not acceptable for the OjG parser. A different approach
to support interfaces was possible.

JSON document of any non trivial size, especially if hand edited are
likely to have errors at some point. Parse errors must identify where
in the document the error occured.

## The Gory Details

Each sub-package of the OjG package had their own approaches to
achieve the optimal performance.

### Generic Data (`gen` package)

What better way to generic type fast than to just define types from
simple types and then define methods on those type. A `gen.Int` is
just an `int64` and a `gen.Array` is just a `[]gen.Node`. With that
approach there are no extra allocations.

Since generic arrays and objects restrict the type of the values in
each collection to `gen.Node` types the collection are assured to
contain only elements that can be encoded as JSON.

The parser for generic types is a copy of the oj package parser but
instead of simple types being created instances that support the
`gen.Node` interface are created.

### Simple Parser (`oj` package)

 - single pass
  - no going back (except one character)

 - modes to large state diagram
 - opt for code duplication dues to overhead of function calls
  - not pretty but the price for performance

 - streaming
  - callbacks

 - minimize allocations
 - reuse buffers
 - no tokens
 - build numbers during parsing, not string first
 - build string in reusable buffer
 - reuse when possible
  - true and false in gen parser
  - building a string map was more expensive
   - might be something to consider next
 - character maps (really a long string/[]byte)


### JSONPath (`jp` package)

A JSONPath is represented by a `jp.Expr` which is composed of
fragments or `jp.Frag` objects. Keeping with the guideline of
minimizing allocations the `jp.Expr` is just a slice of `jp.Frag`. In
most cases expressions are defined statically so the parser need not
be fast so no special care was taken to make that fast. Instead
functions are used in an approach that is easier to understand.

If the need exists to create expressions at run time then functions
are used that allow them to be constructed more easily.

Evaluating an expression against data involves walking down the data
tree to find one or more elements. Conceptually each fragment of a
path sets up zero or more paths to follow through the data. When the
last fragment is reached the search is done. A recursive approach
would be ideal where the evaluation of one fragment then invokes the
next fragment eval function with as many paths it matches. Great on
paper but for something like a descent fragment (`..`) that is a lot
of function calls.

Given that function calls are expensive and slices are cheap a Forth
(the language) evaluation stack approach is used. Each fragment taks
it's matches and puts them on the stack. Then the next fragment
evaluates each in turn. This continue until the stack shrinks back to
one element indicating the evaluation is complete. The last fragment
also puts any matches on a results list which is returned.

 | Stack  | Frag  |
 | ------ | ----- |
 | $      | Root  |
 | {a:3}  | data  |
 | 'a'    | Child |

One fragment type is a filter which looks like `[?(@.x == 3)]`. This
requires a script or function evaluation. A similar stack based
approach is used for evaluating scripts. Note that scripts can and
almost always contain a JSONPath expression starting with a `@`
character. An interesting aspect of this is that a filter can contain
other filters. OjG supports nested filters.

### Converting or Altering Data (`alt` package)

 - parse first then convert, more general, does add overhead

 - not constrained by existing encoder which forces a use pattern that is slower

 - wanted to handle interfaces
  - approach used in ruby oj and other ruby json parsers

 - use exist but assert or replace values

## Interesting Tidbits

Through lots of benchmarking various approach to the implemenation a
few lessons were learned.

### Functions Add Overhead

Sure we all know a function call add some overhead in any language. In
C that overhead is pretty small or nonexistent with inline
functions. That is not true for Go. There is considerable overhead in
making a function call and if that functional call included any kind
on context such as being the function of a type the overhead is even
higher. That observation while disappointing drove a lot of the parser
and JSONPath evaluation code. For nice looking and well organized code
using functions are highly recommended but for high perfomance find a
way to reduce function calls.

The implementation of the parser included a lot of duplicate code to
reduce function calls and it did make a significant difference in
performance.

The JSONPath evaluation takes an additional approach. It includes a
fair amount of code duplication but it also implement its own stack to
avoid nested functional calls even though the nature of the evaluation
is a better match for a recursive implementation.

### Slices are Nice

Slice are implemented very efficiently in Go. Appending to a slice has
very little overhead. Reusing slice by collapsing then to zero length
is a great way to avoid allocating additional memory. Care has to be
taken when collapsing though as any cells in the slice that point to
object will now leave those objects dangling and they will never be
garbage collected. Simply setting the slice slot to `nil` will avoid
memory leaks.

### Memory Allocation

Like most languages, memory allocation adds overhead. Best to avoid
when possible. A good example of that is in the `alt` package. The
`Alter()` function replaces slice and map members instead of
allocating a new slice or map when possible.

Parsers take advantage by reusing buffers and avoiding allocation of
token during possible when possible.

### Range Has Been Optimized

Using a `for` `range` loop is better than incrementing an index. The
difference was not huge but was noticable.

### APIs Matter

It's important to define an API that is easy to use as well as one
that allows for the best performance. The parser as well as the
JSONPath builders attempt to do both. An even better example is the
[GGql](https://github.com/uhn/ggql) GraphQL package. It provides a
very simple AI when compared to previous Go GraphQL packages and it
many times
[faster](https://github.com/the-benchmarker/graphql-benchmarks).


## Whats Next?

Theres alway more to come. For OjG there are a few things in the works.

 - Regex filters for JSONPath
 - Add JSON building to the oj command which is an alternative to jq but uses JSONPath.
 - Implement a Simple Encoding Notation which mixes GraphQL symtax with JSON for simplier more forgiving format.
