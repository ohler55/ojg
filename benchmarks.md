# OjG Benchmarks

Benchmarks were run on a MacBook Pro with a 2.8 GHz Quad-core
I7. Higher numbers in parenthesis are better in all cases.

```

Parse string/[]byte
     json.Unmarshal     42144 ns/op   17982 B/op     336 allocs/op
       oj.Parse         24508 ns/op   18816 B/op     433 allocs/op
      gen.Parse         23930 ns/op   18815 B/op     433 allocs/op
      sen.Parse         25666 ns/op   18770 B/op     427 allocs/op

      gen █████████████████▌ 1.76
       oj █████████████████▏ 1.72
      sen ████████████████▍ 1.64
     json ▓▓▓▓▓▓▓▓▓▓ 1.00

Parse io.Reader
     json.Decode        52568 ns/op   32658 B/op     346 allocs/op
       oj.ParseReader   27764 ns/op   22911 B/op     434 allocs/op
      gen.ParseReder    27554 ns/op   22912 B/op     434 allocs/op
      sen.ParseReader   27654 ns/op   22912 B/op     434 allocs/op

      gen ███████████████████  1.91
      sen ███████████████████  1.90
       oj ██████████████████▉ 1.89
     json ▓▓▓▓▓▓▓▓▓▓ 1.00

Validate string/[]byte
     json.Valid         12007 ns/op       0 B/op       0 allocs/op
       oj.Valdate        5801 ns/op       0 B/op       0 allocs/op

       oj ████████████████████▋ 2.07
     json ▓▓▓▓▓▓▓▓▓▓ 1.00

Validate io.Reader
     json.Decode        52527 ns/op   32656 B/op     346 allocs/op
       oj.Valdate        8835 ns/op    4096 B/op       1 allocs/op

       oj ███████████████████████████████████████████████████████████▍ 5.95
     json ▓▓▓▓▓▓▓▓▓▓ 1.00

to JSON
     json.Marshal       40550 ns/op   17908 B/op     345 allocs/op
       oj.JSON          10082 ns/op    3072 B/op       1 allocs/op
      sen.String        11388 ns/op    2688 B/op       1 allocs/op

       oj ████████████████████████████████████████▏ 4.02
      sen ███████████████████████████████████▌ 3.56
     json ▓▓▓▓▓▓▓▓▓▓ 1.00

to JSON with indentation
     json.Marshal       72788 ns/op   27326 B/op     352 allocs/op
       oj.JSON          12143 ns/op    4096 B/op       1 allocs/op
      sen.String        13615 ns/op    4096 B/op       1 allocs/op

       oj ███████████████████████████████████████████████████████████▉ 5.99
      sen █████████████████████████████████████████████████████▍ 5.35
     json ▓▓▓▓▓▓▓▓▓▓ 1.00

to JSON with indentation and sorted keys
       oj.JSON          19393 ns/op    6560 B/op      63 allocs/op
      sen.String        20641 ns/op    6560 B/op      63 allocs/op

       oj ▓▓▓▓▓▓▓▓▓▓ 1.00
      sen █████████▍ 0.94

Write indented JSON
     json.Encode        73027 ns/op   28383 B/op     353 allocs/op
       oj.Write         11360 ns/op       0 B/op       0 allocs/op

       oj ████████████████████████████████████████████████████████████████▎ 6.43
     json ▓▓▓▓▓▓▓▓▓▓ 1.00

Convert or Alter
      alt.Generify       1229 ns/op    1712 B/op      27 allocs/op
      alt.Alter           932 ns/op     960 B/op      19 allocs/op

      alt █████████████▏ 1.32
      alt ▓▓▓▓▓▓▓▓▓▓ 1.00

JSONPath Get $..a[2].c
       jp.Get            3968 ns/op    1992 B/op      65 allocs/op

       jp ▓▓▓▓▓▓▓▓▓▓ 1.00

JSONPath First  $..a[2].c
       jp.First          2090 ns/op    1272 B/op      32 allocs/op

       jp ▓▓▓▓▓▓▓▓▓▓ 1.00

 Higher values (longer bars) are better in all cases. The bar graph compares the
 parsing performance. The lighter colored bar is the reference, usually the go
 json package.

 The Benchmarks reflect a use case where JSON is either provided as a string or
 read from a file (io.Reader) then parsed into simple go types of nil, bool, int64
 float64, string, []interface{}, or map[string]interface{}. When supported, an
 io.Writer benchmark is also included along with some miscellaneous operations.
```
