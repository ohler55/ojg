# [![{}j](http://www.ohler.com/dev/images/ojg_comet.jpg)](https://github.com/ohler55/ojg)

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

See [file:CHANGELOG.md](CHANGELOG.md)

## Links

- *Documentation*: [https://ohler55.github.io/ojg](https://ohler55.github.io/ojg)

- *GitHub* *repo*: https://github.com/ohler55/ojg

#### Links of Interest

 - *Oj, a Ruby JSON parser*: https://www.ohler.com/oj also at https://github.com/ohler55/oj

 - *OjC, a C JSON parser*: https://www.ohler.com/ojc also at https://github.com/ohler55/ojc

 - *Fast XML parser and marshaller on GitHub*: https://github.com/ohler55/ox

 - *Agoo, a high performance Ruby web server supporting GraphQL on GitHub*: https://github.com/ohler55/agoo

 - *Agoo-C, a high performance C web server supporting GraphQL on GitHub*: https://github.com/ohler55/agoo-c

#### Contributing

+ Provide a Pull Request off the `develop` branch.
+ Report a bug
+ Suggest an idea
