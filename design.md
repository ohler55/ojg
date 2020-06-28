# A Journey building a fast JSON parser and full JSONPath for Go

I had a dream. I'd write a fast JSON parser, generic data, and
JSONPath implementation and it would be beautiful, well organized, and
something to be admired. Well, reality kicked in laughed at those
dreams. A Go JSON parser and tools could be high performance but to
get that performance compromised in beauty would have to be made. This
is a tale of journey that ended with a Parser that leaves the Go JSON
parser in the dust and resulted in some useful tools including an
efficient JSONPath implemenation.

In all fairness I did embark on with some previous experience having
written two JSON parser before. Both the Ruby
[Oj](https://github.com/ohler55/oj) and the C parser
[OjC](https://github.com/ohler55/ojc) are the best performers in their
respective languages. Why not an [OjG](https://github.com/ohler55/ojg)
for go.

## Planning

Like any journey it starts with the planning. Yeah, I know, its called
requirement gathering but casting it as planning a journey sounds like
more fun and this was all about enjoying the discoveries on the
journey. The journey takes place in the land of OjG which stands for
Oj for Go. [Oj](https://github.com/ohler55/oj) or Optimized JSON being
a popular gem I wrote for Ruby.

First, JSON parsing and any frequently used operations such as
JSONPath evaluation had to be fast over everything else. With the
luxury of not having to follow the existing Go json package API the
API could be designed for the best performance.

The journey would visit several areas. In each area the planning would
vary from the the others.

### Generic Data

The first area to visit was generic data. Not to be confused with the
propose Go generics. That a completely different animal and has
nothing to do with whats being referred to as generic data here. In
building tool or packages for reuse the data acted on by those tools
needs to be navigatable.

Reflection can be used but that gets a bit tricky when dealing with
private fields or field that can't be converted to something that can
say be written as a JSON element.

Another approach is to use simple Go types such as bool, int64,
[]interface{}, and other types that map directly on to JSON or some
other subset of all possible Go types. If too open, such as with
[]interface{} it is still possible for the user to put unsupported
types into the data. Not to pick out any package specifically but it
is frustrating to see an argument type of interface{} and then no
documentation describing that the supported types are.

There is another approach though. Define a set of types that can be in
a collection and use those types. With this approach, the generic data
implementation has to support the basic JSON types of `null`, boolean,
number, string, array, and object. In addition time should be
supported. In both JSON use in Ruby and Go time has alway been
needed. Time is just too much apart of any set of data to leave it
out.

The generic data had to be type safe. It would not do to have an
element that could not be endcoded as JSON in the data.

A frequent operation for generic data is to store that data into a
JSON database or similar. That meant converting to simple Go types of
nil, bool, int64, float64, string, []interface{}, and
map[string]interface{} had to be fast.

Also planned for this part of the journey was methods on the types to
support getting, setting, and deleting elements using JSONPath. The
hope was to have an object based approach to the generic nodes so
something like the following could be used but keeping generic data,
JONPath, and parsing in separate packages.

```golang
    var n gen.Node
    n = gen.Int(123)
    i, ok := n.AsInt()
```

Unfortunately that part of the journey had to be cancelled as the Go
travel guide refuses to let packages talk back and forth. Imports are
one way only. After trying to put all the code in one package it
eventually got unwieldy function names started being prefixed with
what should really have been package names so the object and method
approach was dropped. A change in API but the journey would continue.

### JSON Parser and Validotor

The next stop was the parser and validator. After some consideration
it seemed like starting with the validator would be best way to become
familiar with the territory. The JSON parser and validator need not be
the same and each should be as performant as possible. The parsers
needed to support parsing into simple Go types as well as the generic
data types.

When parsing files that include millions or more JSON elements in
files that might be over 100GB a streaming parser is necessary. It
would be nice to share some code with both the streaming and string
parsers of course. Its easier to pack light when the areas are
similar.

The parser must also allow parsing into Go types. Furthermore
interfaces must be supported. Go unmarshalling does not support
interface fields. Since many data types make use of interfaces that
limitation was not acceptable for the OjG parser. A different approach
to support interfaces was possible.

JSON document of any non trivial size, especially if hand edited are
likely to have errors at some point. Parse errors must identify where
in the document the error occured.

### JSONPath

Saving the most interesting part of the trip for last, the JSONPath
implementation promised to have all sorts of interesting problems to
solve with descents, wildcards, and especially filters.

A JSONPath is used to extract elements from data. That part of the
implementation had to be fast. Parsing really didn't have to be fast
but it would be nice to have a way of building a JSONPath in a
performant manner even if it was not as convenient as parsing a
string.

The JSONPath implementaiton had implement all the features described
by the [Goessner
article](https://goessner.net/articles/JsonPath). There are other
descriptions of JSONPath but the Goessner description is the most
referenced. Since the implementation is in Go the scripting feature
described could be left out as long as similar functionality could be
provided for array indexes based on the length of the array. Borrowing
from Ruby, using negative indexes would provide that functionality.

## The Journey

The journey unfolded as planned to a degree. There were some false
starts and revisits but eventually each destination was reached and
the journey completed.

### Generic Data (`gen` package)

What better way to make generic type fast than to just define generic
types from simple types and then define methods on those types. A
`gen.Int` is just an `int64` and a `gen.Array` is just a
`[]gen.Node`. With that approach there are no extra allocations.

```golang
type Node interface{}
type Int int64
type Array []Node
```

Since generic arrays and objects restrict the type of the values in
each collection to `gen.Node` types the collection are assured to
contain only elements that can be encoded as JSON.

Since methods on the `Node` could not be implemented without import
loops the number of functions in the `Node` interface were limited. It
was clear a parser would be needed but that would have to wait until
the next part of the journey was completed. Then the generic data
package could be revisited and the parser explored.

Let just ahead to the generic data parser revisit. It was not very
interesting after the deep dive into the simple data parser. The
parser for generic types is a copy of the oj package parser but
instead of simple types being created instances that support the
`gen.Node` interface are created.

### Simple Parser (`oj` package)

Looking back it hard to say what was the most interesting part of the
journey, the parser or JSONPath. Each had their own unique set of
issues. The parser was the best place to start though as some valuable
lessons were learned about what to avoid and what to gravitate toward
in trying to achieve high performance code.

--- left off here -----------
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

** TBD redo **

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
(the language) evaluation stack approach is used. Not exactly Forth
but a similar approach. Each fragment takes it's matches and puts them
on the stack. Then the next fragment evaluates each in turn. This
continue until the stack shrinks back to one element indicating the
evaluation is complete. The last fragment also puts any matches on a
results list which is returned.

 | Stack  | Frag  |
 | ------ | ----- |
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

Benchmarking was instrumental to tuning and picking the most favorable
approach to the implemenation. Through those benchmarks a number of
lessons were learned.  The final benchmarks results can be viewed by
running the `cmd/benchmark` command. See the results at
[benchmarks.md](benchmarks.md).

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

Theres alway something new ready to be explored. For OjG there are a few things in the planning stage.

 - A short trip to Regex filters for JSONPath.
 - A construction project to add JSON building to the **oj** command which is an alternative to jq but using JSONPath.
 - Explore new territory by implementing a Simple Encoding Notation which mixes GraphQL symtax with JSON for simplier more forgiving format.
