# OjG Benchmarks

Benchmarks were run on a MacBook Pro with a 2.8 GHz Quad-core I7.

```
Parse string/[]byte
       json.Unmarshal       40464 ns/op   17984 B/op     336 allocs/op
         oj.Parse           22845 ns/op   18783 B/op     431 allocs/op
   oj-reuse.Parse           17546 ns/op    5984 B/op     366 allocs/op
        gen.Parse           23126 ns/op   18784 B/op     431 allocs/op
  gen-reuse.Parse           17789 ns/op    5984 B/op     366 allocs/op
        sen.Parse           24824 ns/op   18736 B/op     427 allocs/op
  sen-reuse.Parse           19475 ns/op    5920 B/op     362 allocs/op

   oj-reuse ███████████████████████  2.31
  gen-reuse ██████████████████████▋ 2.27
  sen-reuse ████████████████████▊ 2.08
         oj █████████████████▋ 1.77
        gen █████████████████▍ 1.75
        sen ████████████████▎ 1.63
       json ▓▓▓▓▓▓▓▓▓▓ 1.00

Parse io.Reader
       json.Decode          51431 ns/op   32657 B/op     346 allocs/op
         oj.ParseReader     26502 ns/op   22881 B/op     432 allocs/op
   oj-reuse.ParseReader     20803 ns/op   10080 B/op     367 allocs/op
        gen.ParseReder      26497 ns/op   22880 B/op     432 allocs/op
  gen-reuse.ParseReder      21229 ns/op   10080 B/op     367 allocs/op
        sen.ParseReader     28020 ns/op   22881 B/op     434 allocs/op
  sen-reuse.ParseReader     22274 ns/op   10080 B/op     369 allocs/op

   oj-reuse ████████████████████████▋ 2.47
  gen-reuse ████████████████████████▏ 2.42
  sen-reuse ███████████████████████  2.31
        gen ███████████████████▍ 1.94
         oj ███████████████████▍ 1.94
        sen ██████████████████▎ 1.84
       json ▓▓▓▓▓▓▓▓▓▓ 1.00

Parse chan interface{}
       json.Parse-chan      50370 ns/op   18000 B/op     337 allocs/op
         oj.Parse           29183 ns/op   18784 B/op     431 allocs/op
        gen.Parse           29555 ns/op   18784 B/op     431 allocs/op
        sen.Parse           31553 ns/op   18736 B/op     427 allocs/op

         oj █████████████████▎ 1.73
        gen █████████████████  1.70
        sen ███████████████▉ 1.60
       json ▓▓▓▓▓▓▓▓▓▓ 1.00

Validate string/[]byte
       json.Valid           12726 ns/op       0 B/op       0 allocs/op
         oj.Valdate          4362 ns/op       0 B/op       0 allocs/op

         oj █████████████████████████████▏ 2.92
       json ▓▓▓▓▓▓▓▓▓▓ 1.00

Validate io.Reader
       json.Decode          51417 ns/op   32656 B/op     346 allocs/op
         oj.Valdate          7095 ns/op    4096 B/op       1 allocs/op

         oj ████████████████████████████████████████████████████████████████████████▍ 7.25
       json ▓▓▓▓▓▓▓▓▓▓ 1.00

to JSON
       json.Marshal         40063 ns/op   17908 B/op     345 allocs/op
         oj.JSON             9964 ns/op    3072 B/op       1 allocs/op
        sen.String          10730 ns/op    2688 B/op       1 allocs/op

         oj ████████████████████████████████████████▏ 4.02
        sen █████████████████████████████████████▎ 3.73
       json ▓▓▓▓▓▓▓▓▓▓ 1.00

to JSON with indentation
       json.Marshal         71587 ns/op   27327 B/op     352 allocs/op
         oj.JSON            11950 ns/op    4096 B/op       1 allocs/op
        sen.String          12524 ns/op    4096 B/op       1 allocs/op

         oj ███████████████████████████████████████████████████████████▉ 5.99
        sen █████████████████████████████████████████████████████████▏ 5.72
       json ▓▓▓▓▓▓▓▓▓▓ 1.00

to JSON with indentation and sorted keys
         oj.JSON            18104 ns/op    6560 B/op      63 allocs/op
        sen.String          19254 ns/op    6560 B/op      63 allocs/op

         oj ▓▓▓▓▓▓▓▓▓▓ 1.00
        sen █████████▍ 0.94

Write indented JSON
       json.Encode          71947 ns/op   28384 B/op     353 allocs/op
         oj.Write           11170 ns/op       0 B/op       0 allocs/op

         oj ████████████████████████████████████████████████████████████████▍ 6.44
       json ▓▓▓▓▓▓▓▓▓▓ 1.00

Convert or Alter
        alt.Generify         1418 ns/op    1712 B/op      27 allocs/op
        alt.Alter             885 ns/op     960 B/op      19 allocs/op

        alt ████████████████  1.60
        alt ▓▓▓▓▓▓▓▓▓▓ 1.00

JSONPath Get $..a[2].c
         jp.Get              3603 ns/op    1992 B/op      65 allocs/op

         jp ▓▓▓▓▓▓▓▓▓▓ 1.00

JSONPath First  $..a[2].c
         jp.First            2010 ns/op    1272 B/op      32 allocs/op

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
