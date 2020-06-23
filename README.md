# [![{}j](http://www.ohler.com/dev/images/ojg_comet.jpg)](https://github.com/ohler55/ojg)

[![Build Status](https://img.shields.io/travis/ohler55/ojg/master.svg?logo=travis)](http://travis-ci.org/ohler55/ojg?branch=master)[![Coverage Status](https://coveralls.io/repos/github/ohler55/ojg/badge.svg?branch=master)](https://coveralls.io/github/ohler55/ojg?branch=master)

# OjG

Optimized JSON for Go is a high performance parser with a variety of
additional JSON tools including a JSONPath implemenation that will
operation on golang structs as well as simple types.

## Using

```golang

    v, err := oj.ParseString("[true,[false,[null],123],456]")

```

## Installation
```
go get github.com/ohler55/ojg

```

or just import

```
import (
    "github.com/ohler55/ojg/alt"
    "github.com/ohler55/ojg/gen"
    "github.com/ohler55/ojg/jp"
    "github.com/ohler55/ojg/oj"
)
```

## Releases

See [CHANGELOG.md](CHANGELOG.md)

## Links

- *Documentation*: [https://ohler55.github.com/ojg](https://ohler55.github.com/ojg)

- *GitHub* *repo*: https://github.com/ohler55/ojg

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
