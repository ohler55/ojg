# [![{}j](http://www.ohler.com/dev/images/ojg_comet.jpg)](https://github.com/ohler55/ojg)

[![Build Status](https://img.shields.io/travis/ohler55/ojg/master.svg?logo=travis)](http://travis-ci.org/ohler55/ojg?branch=master)[![Coverage Status](https://coveralls.io/repos/github/ohler55/ojg/badge.svg?branch=master)](https://coveralls.io/github/ohler55/ojg?branch=master)

# OjG

Optimized JSON for Go is a high performance parser with a variety of
additional JSON tools.

## Features

 - Fast JSON parser. Check out the cmd/benchmarks app in this repo.
 - Full JSONPath implemenation that operates on simple types as well as structs.
 - Generic types. Not the proposed golang generics but type safe JSON elements.
 - Fast JSON validator (4 times faster with io.Reader).
 - Fast JSON writer with a sort option (1.6 times faster).
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

The **oj** command (cmd/oj) which uses JSON path for filtering and
extracting JSON elements. It also includes sorting, reformatting, and
colorizing options.

```
$ oj -m "(@.name == 'Pete')" myfile.json

```

## Installation
```
go get github.com/ohler55/ojg
go get github.com/ohler55/ojg/cmd/oj

```

or just import in your `.go` files.

```
import (
    "github.com/ohler55/ojg/alt"
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

Higher numbers in parenthesis are better.

```
Parse string/[]byte
 json.Unmarshal     42649 ns/op (1.00x)   17984 B/op (1.00x)     336 allocs/op (1.00x)
   oj.Parse         24584 ns/op (1.73x)   18816 B/op (0.96x)     433 allocs/op (0.78x)
  gen.Parse         24315 ns/op (1.75x)   18816 B/op (0.96x)     433 allocs/op (0.78x)
  sen.Parse         26164 ns/op (1.63x)   18752 B/op (0.96x)     427 allocs/op (0.79x)

Parse io.Reader
 json.Decode        52987 ns/op (1.00x)   32656 B/op (1.00x)     346 allocs/op (1.00x)
   oj.ParseReader   28079 ns/op (1.89x)   22913 B/op (1.43x)     434 allocs/op (0.80x)
  gen.ParseReder    28058 ns/op (1.89x)   22912 B/op (1.43x)     434 allocs/op (0.80x)
  sen.ParseReader   28238 ns/op (1.88x)   22913 B/op (1.43x)     434 allocs/op (0.80x)

to JSON
 json.Marshal       42127 ns/op (1.00x)   17908 B/op (1.00x)     345 allocs/op (1.00x)
   oj.JSON          10387 ns/op (4.06x)    3072 B/op (5.83x)       1 allocs/op (345.00x)
   oj.Write          9579 ns/op (4.40x)       0 B/op (+Infx)       0 allocs/op (+Infx)
  sen.String        11453 ns/op (3.68x)    2688 B/op (6.66x)       1 allocs/op (345.00x)
```

See [all benchmarks](benchmarks.md)

## Releases

See [CHANGELOG.md](CHANGELOG.md)

## Links

- *Documentation*: [https://pkg.go.dev/github.com/ohler55/ojg](https://pkg.go.dev/github.com/ohler55/ojg)

- *GitHub* *repo*: https://github.com/ohler55/ojg

- *JSONPath* description: https://goessner.net/articles/JsonPath

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
