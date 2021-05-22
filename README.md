# [![{}j](assets/ojg_comet.svg)](https://github.com/ohler55/ojg)

[![Build Status](https://img.shields.io/travis/ohler55/ojg/master.svg?logo=travis)](http://travis-ci.org/ohler55/ojg?branch=master)[![Coverage Status](https://coveralls.io/repos/github/ohler55/ojg/badge.svg?branch=master)](https://coveralls.io/github/ohler55/ojg?branch=master)

Optimized JSON for Go is a high performance parser with a variety of
additional JSON tools. OjG is optimized to processing huge data sets
where data does not necessarily conform to a fixed structure.

## Features

 - Fast JSON parser. Check out the cmd/benchmarks app in this repo.
 - Full JSONPath implemenation that operates on simple types as well as structs.
 - Generic types. Not the proposed golang generics but type safe JSON elements.
 - Fast JSON validator (7 times faster with io.Reader).
 - Fast JSON writer with a sort option (4 times faster).
 - JSON builder from JSON sources using a simple assembly plan.
 - Simple data builders using a push and pop approach.
 - Object encoding and decoding using an approach similar to that used with Oj for Ruby.
 - [Simple Encoding Notation](sen.md), a lazy way to write JSON omitting commas and quotes.

## Using

A basic Parse:

```golang
    obj, err := oj.ParseString(`{
        "a":[
            {"x":1,"y":2,"z":3},
            {"x":2,"y":4,"z":6}
        ]
    }`)
```

Using JSONPath expressions:

```golang
    x, err := jp.ParseString("a[?(@.x > 1)].y")
    ys := x.Get(obj)
    // returns [4]
```

The **oj** command (cmd/oj) uses JSON path for filtering and
extracting JSON elements. It also includes sorting, reformatting, and
colorizing options.

```
$ oj -m "(@.name == 'Pete')" myfile.json

```

More complete examples are available in the go docs for most
functions. The example for [Unmarshalling
interfaces](oj/example_interface_test.go) demonstrates a feature that
allows interfaces to be marshalled and unmarshalled.

## Installation
```
go get github.com/ohler55/ojg
go get github.com/ohler55/ojg/cmd/oj

```

or just import in your `.go` files.

```
import (
    "github.com/ohler55/ojg/alt"
    "github.com/ohler55/ojg/asm"
    "github.com/ohler55/ojg/gen"
    "github.com/ohler55/ojg/jp"
    "github.com/ohler55/ojg/oj"
    "github.com/ohler55/ojg/sen"
)
```

To build and install the `oj` application:

```
go install ./...
```

## Benchmarks

Higher numbers (longer bars) are better.

```
Parse string/[]byte
     json.Unmarshal         47204 ns/op   17777 B/op     334 allocs/op
       oj.Parse             27983 ns/op   18489 B/op     429 allocs/op
 oj-reuse.Parse             18573 ns/op    5691 B/op     364 allocs/op

 oj-reuse.Parse        █████████████████████████▍ 2.54
       oj.Parse        ████████████████▊ 1.69
     json.Unmarshal    ▓▓▓▓▓▓▓▓▓▓ 1.00

Parse io.Reader
     json.Decode            68799 ns/op   32448 B/op     344 allocs/op
       oj.ParseReader       37805 ns/op   22584 B/op     430 allocs/op
 oj-reuse.ParseReader       23222 ns/op    9788 B/op     365 allocs/op
       oj.TokenizeLoad      11804 ns/op    6072 B/op     157 allocs/op

       oj.TokenizeLoad ██████████████████████████████████████████████▋ 5.83
 oj-reuse.ParseReader  ███████████████████████▋ 2.96
       oj.ParseReader  ██████████████▌ 1.82
     json.Decode       ▓▓▓▓▓▓▓▓ 1.00

to JSON with indentation
     json.Marshal           73578 ns/op   26983 B/op     352 allocs/op
       oj.JSON               7532 ns/op       0 B/op       0 allocs/op
      sen.Bytes              8982 ns/op       0 B/op       0 allocs/op

       oj.JSON         ██████████████████████████████████████████████████████████████████████████████▏ 9.77
      sen.Bytes        █████████████████████████████████████████████████████████████████▌ 8.19
     json.Marshal      ▓▓▓▓▓▓▓▓ 1.00
```

See [all benchmarks](benchmarks.md)

[Compare Go JSON parsers](https://github.com/ohler55/compare-go-json)

## Releases

See [CHANGELOG.md](CHANGELOG.md)

## Links

- *Documentation*: [https://pkg.go.dev/github.com/ohler55/ojg](https://pkg.go.dev/github.com/ohler55/ojg)

- *GitHub* *repo*: https://github.com/ohler55/ojg

- *JSONPath* description: https://goessner.net/articles/JsonPath

- *JSONPath Comparisons* https://cburgmer.github.io/json-path-comparison


#### Links of Interest

 - *Oj, a Ruby JSON parser*: http://www.ohler.com/oj/doc/index.html also at https://github.com/ohler55/oj

 - *OjC, a C JSON parser*: http://www.ohler.com/ojc/doc/index.html also at https://github.com/ohler55/ojc

 - *Fast XML parser and marshaller on GitHub*: https://github.com/ohler55/ox

 - *Agoo, a high performance Ruby web server supporting GraphQL on GitHub*: https://github.com/ohler55/agoo

 - *Agoo-C, a high performance C web server supporting GraphQL on GitHub*: https://github.com/ohler55/agoo-c

#### Contributing

+ Provide a Pull Request off the `develop` branch.
+ Report a bug
+ Suggest an idea
