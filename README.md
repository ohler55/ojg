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
)
```

To build and install the `oj` application:

```
go install ./...
```

## Benchmarks

Higher numbers in parenthesis are better.

```
Parse JSON
json.Unmarshal:       7104 ns/op (1.00x)    4808 B/op (1.00x)      90 allocs/op (1.00x)
  oj.Parse:           4518 ns/op (1.57x)    3984 B/op (1.21x)      86 allocs/op (1.05x)
  oj.GenParse:        4623 ns/op (1.54x)    3984 B/op (1.21x)      86 allocs/op (1.05x)

Validate io.Reader JSON
json.Decoder:        50213 ns/op (1.00x)   32658 B/op (1.00x)     346 allocs/op (1.00x)
  oj.ValidateReader: 12740 ns/op (3.94x)    4096 B/op (7.97x)       1 allocs/op (346.00x)

JSON()
json.Marshal:         2616 ns/op (1.00x)     992 B/op (1.00x)      22 allocs/op (1.00x)
  oj.JSON:             436 ns/op (6.00x)     131 B/op (7.57x)       4 allocs/op (5.50x)
  oj.Write:            455 ns/op (5.75x)     131 B/op (7.57x)       4 allocs/op (5.50x)
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
