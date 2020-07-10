# OjG Benchmarks

Benchmarks were run on a MacBook Pro with a 2.8 GHz Quad-core
I7. Higher numbers in parenthesis are better in all cases.

```
 The number in parenthesis are the ratio of results between the reference and
 the listed. Higher values are better.

 The Benchmarks reflect a use case where JSON is either provided as a string or
 read from a file (io.Reader) then parsed into simple go types of nil, bool, int64
 float64, string, []interface{}, or map[string]interface{}. When supported, an
 io.Writer benchmark is also included along with some miscellaneous operations.

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

Validate string/[]byte
 json.Valid         12700 ns/op (1.00x)       0 B/op (NaNx)       0 allocs/op (NaNx)
   oj.Valdate        5905 ns/op (2.15x)       0 B/op (NaNx)       0 allocs/op (NaNx)

Validate io.Reader
 json.Decode        53228 ns/op (1.00x)   32656 B/op (1.00x)     346 allocs/op (1.00x)
   oj.Valdate        8892 ns/op (5.99x)    4096 B/op (7.97x)       1 allocs/op (346.00x)

to JSON
 json.Marshal       42127 ns/op (1.00x)   17908 B/op (1.00x)     345 allocs/op (1.00x)
   oj.JSON          10387 ns/op (4.06x)    3072 B/op (5.83x)       1 allocs/op (345.00x)
   oj.Write          9579 ns/op (4.40x)       0 B/op (+Infx)       0 allocs/op (+Infx)
  sen.String        11453 ns/op (3.68x)    2688 B/op (6.66x)       1 allocs/op (345.00x)

to JSON with indentation
 json.Marshal       74958 ns/op (1.00x)   27325 B/op (1.00x)     352 allocs/op (1.00x)
   oj.JSON          12552 ns/op (5.97x)    4096 B/op (6.67x)       1 allocs/op (352.00x)
  sen.String        13750 ns/op (5.45x)    4096 B/op (6.67x)       1 allocs/op (352.00x)

to JSON with indentation and sorted keys
   oj.JSON          19466 ns/op (1.00x)    6560 B/op (1.00x)      63 allocs/op (1.00x)
  sen.String        20496 ns/op (0.95x)    6560 B/op (1.00x)      63 allocs/op (1.00x)

Write indented JSON
 json.Encode        74690 ns/op (1.00x)   28384 B/op (1.00x)     353 allocs/op (1.00x)
   oj.Write         11274 ns/op (6.62x)       0 B/op (+Infx)       0 allocs/op (+Infx)

Convert or Alter
  alt.Generify       1234 ns/op (1.00x)    1712 B/op (1.00x)      27 allocs/op (1.00x)
  alt.Alter           927 ns/op (1.33x)     960 B/op (1.78x)      19 allocs/op (1.42x)

JSONPath Get $..a[2].c
   jp.Get            3778 ns/op (1.00x)    1992 B/op (1.00x)      65 allocs/op (1.00x)

JSONPath First  $..a[2].c
   jp.First          2094 ns/op (1.00x)    1272 B/op (1.00x)      32 allocs/op (1.00x)
```
