# OjG Benchmarks

Benchmarks were run on a MacBook Pro with a 2.8 GHz Quad-core I7.

```

Parse string/[]byte
     json.Unmarshal     41858 ns/op   17985 B/op     336 allocs/op
       oj.Parse         24407 ns/op   18815 B/op     433 allocs/op
      gen.Parse         23998 ns/op   18816 B/op     433 allocs/op
      sen.Parse         26097 ns/op   18737 B/op     427 allocs/op

      gen █████████████████▍ 1.74
       oj █████████████████▏ 1.71
      sen ████████████████  1.60
     json ▓▓▓▓▓▓▓▓▓▓ 1.00

Parse io.Reader
     json.Decode        52997 ns/op   32657 B/op     346 allocs/op
       oj.ParseReader   28645 ns/op   22912 B/op     434 allocs/op
      gen.ParseReder    27330 ns/op   22912 B/op     434 allocs/op
      sen.ParseReader   28103 ns/op   22912 B/op     434 allocs/op

      gen ███████████████████▍ 1.94
      sen ██████████████████▊ 1.89
       oj ██████████████████▌ 1.85
     json ▓▓▓▓▓▓▓▓▓▓ 1.00

Validate string/[]byte
     json.Valid         11974 ns/op       0 B/op       0 allocs/op
       oj.Valdate        5095 ns/op       0 B/op       0 allocs/op

       oj ███████████████████████▌ 2.35
     json ▓▓▓▓▓▓▓▓▓▓ 1.00

Validate io.Reader
     json.Decode        53027 ns/op   32658 B/op     346 allocs/op
       oj.Valdate        8171 ns/op    4096 B/op       1 allocs/op

       oj ████████████████████████████████████████████████████████████████▉ 6.49
     json ▓▓▓▓▓▓▓▓▓▓ 1.00

to JSON
     json.Marshal       41151 ns/op   17908 B/op     345 allocs/op
       oj.JSON          11473 ns/op    3072 B/op       1 allocs/op
      sen.String        11193 ns/op    2688 B/op       1 allocs/op

      sen ████████████████████████████████████▊ 3.68
       oj ███████████████████████████████████▊ 3.59
     json ▓▓▓▓▓▓▓▓▓▓ 1.00

to JSON with indentation
     json.Marshal       73681 ns/op   27329 B/op     352 allocs/op
       oj.JSON          13021 ns/op    4096 B/op       1 allocs/op
      sen.String        13109 ns/op    4096 B/op       1 allocs/op

       oj ████████████████████████████████████████████████████████▌ 5.66
      sen ████████████████████████████████████████████████████████▏ 5.62
     json ▓▓▓▓▓▓▓▓▓▓ 1.00

to JSON with indentation and sorted keys
       oj.JSON          19846 ns/op    6560 B/op      63 allocs/op
      sen.String        19911 ns/op    6560 B/op      63 allocs/op

       oj ▓▓▓▓▓▓▓▓▓▓ 1.00
      sen █████████▉ 1.00

Write indented JSON
     json.Encode        74459 ns/op   28384 B/op     353 allocs/op
       oj.Write         12131 ns/op       0 B/op       0 allocs/op

       oj █████████████████████████████████████████████████████████████▍ 6.14
     json ▓▓▓▓▓▓▓▓▓▓ 1.00

Convert or Alter
      alt.Generify       1246 ns/op    1712 B/op      27 allocs/op
      alt.Alter           919 ns/op     960 B/op      19 allocs/op

      alt █████████████▌ 1.36
      alt ▓▓▓▓▓▓▓▓▓▓ 1.00

JSONPath Get $..a[2].c
       jp.Get            3703 ns/op    1992 B/op      65 allocs/op

       jp ▓▓▓▓▓▓▓▓▓▓ 1.00

JSONPath First  $..a[2].c
       jp.First          2058 ns/op    1272 B/op      32 allocs/op

       jp ▓▓▓▓▓▓▓▓▓▓ 1.00

 Higher values (longer bars) are better in all cases. The bar graph compares the
 parsing performance. The lighter colored bar is the reference, usually the go
 json package.

 The Benchmarks reflect a use case where JSON is either provided as a string or
 read from a file (io.Reader) then parsed into simple go types of nil, bool, int64
 float64, string, []interface{}, or map[string]interface{}. When supported, an
 io.Writer benchmark is also included along with some miscellaneous operations.

Tests run on:
 Machine:         MacBookPro15,2
 OS:              Mac OS X 10.15.5
 Processor:       Quad-Core Intel Core i7
 Cores:           4
 Processor Speed: 2.8 GHz
```
