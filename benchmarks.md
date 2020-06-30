# OjG Benchmarks

Benchmarks were run on a MacBook Pro with a 2.8 GHz Quad-core
I7. Higher number in parenthesis are better in all cases.

```
JSON Path Get
  oj.Expr.Get:            7386 ns/op (1.00x)    5568 B/op (1.00x)      67 allocs/op (1.00x)

JSON Path First
  oj.Expr.First:          2065 ns/op (1.00x)    1248 B/op (1.00x)      29 allocs/op (1.00x)

Parse JSON
json.Unmarshal:           7104 ns/op (1.00x)    4808 B/op (1.00x)      90 allocs/op (1.00x)
  oj.Parse:               4518 ns/op (1.57x)    3984 B/op (1.21x)      86 allocs/op (1.05x)
  oj.GenParse:            4623 ns/op (1.54x)    3984 B/op (1.21x)      86 allocs/op (1.05x)

Parse io.Reader JSON
json.Decoder:            49962 ns/op (1.00x)   32655 B/op (1.00x)     346 allocs/op (1.00x)
  oj.ParseReader:        30380 ns/op (1.64x)   22912 B/op (1.43x)     434 allocs/op (0.80x)
  oj.GenParseReader:     31831 ns/op (1.57x)   22913 B/op (1.43x)     434 allocs/op (0.80x)

Validate JSON
json.Valid:               1434 ns/op (1.00x)       0 B/op (1.00x)       0 allocs/op (1.00x)
  oj.Validate:            1241 ns/op (1.16x)       0 B/op ( NaNx)       0 allocs/op ( NaNx)

Validate io.Reader JSON
json.Decoder:            50213 ns/op (1.00x)   32658 B/op (1.00x)     346 allocs/op (1.00x)
  oj.ValidateReader:     12740 ns/op (3.94x)    4096 B/op (7.97x)       1 allocs/op (346.00x)

JSON() benchmarks, indent: false, sort: false
json.Marshal:             2616 ns/op (1.00x)     992 B/op (1.00x)      22 allocs/op (1.00x)
  oj.JSON:                 436 ns/op (6.00x)     131 B/op (7.57x)       4 allocs/op (5.50x)
  oj.Write:                455 ns/op (5.75x)     131 B/op (7.57x)       4 allocs/op (5.50x)

JSON() benchmarks, indent: true, sort: false
json.Marshal:             3710 ns/op (1.00x)    1488 B/op (1.00x)      25 allocs/op (1.00x)
  oj.JSON:                 553 ns/op (6.71x)     179 B/op (8.31x)       4 allocs/op (6.25x)
  oj.Write:                572 ns/op (6.49x)     179 B/op (8.31x)       4 allocs/op (6.25x)

JSON() benchmarks, sort: true
  oj.JSON:                 886 ns/op (1.00x)     323 B/op (10.00x)       8 allocs/op (1.00x)
  oj.Write:                916 ns/op (0.97x)     323 B/op (1.00x)       8 allocs/op (1.00x)

Converting from simple to canonical types benchmarks
  oj.Generify:             908 ns/op (1.0x)    1248 B/op (1.0x)      22 allocs/op (1.0x)
  oj.GenAlter:             546 ns/op (1.7x)     480 B/op (2.6x)      13 allocs/op (1.7x)
```
