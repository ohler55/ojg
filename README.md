# [![{}j](assets/ojg_comet.svg)](https://github.com/ohler55/ojg)

[![Build Status](https://github.com/ohler55/ojg/actions/workflows/CI.yml/badge.svg)](https://github.com/ohler55/ojg/actions)
[![Go Report Card](https://goreportcard.com/badge/github.com/ohler55/ojg)](https://goreportcard.com/report/github.com/ohler55/ojg)
![](assets/coverage-badge.svg)

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
 - [JSON and SEN Discovery](discover.md) as a package or option to the **oj** application.

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
extracting JSON elements. It also includes sorting, reformatting,
discovery of JSON in documents, and colorizing options.

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
    "github.com/ohler55/ojg/discover"
)
```

To build and install the `oj` application:

```
go install ./...
```

The `oj` application can be installed with brew.

```
brew install oj
```

## Benchmarks

Higher numbers (longer bars) are better. The run was on 2026-11-19.

```
Parse string/[]byte
       json.Unmarshal           18330 ns/op    17812 B/op    389 allocs/op
         oj.Parse                7676 ns/op    17336 B/op    432 allocs/op
   oj-reuse.Parse                6225 ns/op     5688 B/op    364 allocs/op

   oj-reuse.Parse        ████████████████████▌ 2.94
         oj.Parse        ████████████████▋ 2.39
       json.Unmarshal    ▓▓▓▓▓▓▓ 1.00

Parse io.Reader
       json.Decode              24667 ns/op    26379 B/op    401 allocs/op
         oj.ParseReader          9014 ns/op    21432 B/op    433 allocs/op
   oj-reuse.ParseReader          7429 ns/op     9784 B/op    365 allocs/op
         oj.TokenizeLoad         5137 ns/op     6072 B/op    157 allocs/op

         oj.TokenizeLoad █████████████████████████████████▌ 4.80
   oj-reuse.ParseReader  ███████████████████████▏ 3.32
         oj.ParseReader  ███████████████████▏ 2.74
       json.Decode       ▓▓▓▓▓▓▓ 1.00

to JSON
       json.Marshal              6547 ns/op     3573 B/op     34 allocs/op
         oj.JSON                 3928 ns/op        0 B/op      0 allocs/op
        sen.SEN                  4473 ns/op        0 B/op      0 allocs/op

         oj.JSON         ███████████▋ 1.67
        sen.SEN          ██████████▏ 1.46
       json.Marshal      ▓▓▓▓▓▓▓ 1.00
```

See [all benchmarks](benchmarks.md)

[Compare Go JSON parsers](https://github.com/ohler55/compare-go-json)

## Releases

See [CHANGELOG.md](CHANGELOG.md)

## Links

- *Documentation*: [https://pkg.go.dev/github.com/ohler55/ojg](https://pkg.go.dev/github.com/ohler55/ojg)

- *GitHub* *repo*: https://github.com/ohler55/ojg

- *JSONPath* draft specification: https://datatracker.ietf.org/doc/draft-ietf-jsonpath-base

- *JSONPath Comparisons*: https://cburgmer.github.io/json-path-comparison

- *Go Report Card*: https://goreportcard.com/report/github.com/ohler55/ojg

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
