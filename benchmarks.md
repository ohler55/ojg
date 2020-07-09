# OjG Benchmarks

Benchmarks were run on a MacBook Pro with a 2.8 GHz Quad-core
I7. Higher numbers in parenthesis are better in all cases.

```
 The number in parenthesis are the ratio of results between the reference and
 the listed. Higher values are better.

Parse string/[]byte
 json.Unmarshal      7298 ns/op (1.00x)    4808 B/op (1.00x)      90 allocs/op (1.00x)
   oj.Parse          4414 ns/op (1.65x)    3984 B/op (1.21x)      86 allocs/op (1.05x)
  gen.Parse          4453 ns/op (1.64x)    3984 B/op (1.21x)      86 allocs/op (1.05x)
  sen.Parse          4360 ns/op (1.67x)    3984 B/op (1.21x)      86 allocs/op (1.05x)

Parse io.Reader
 json.Decode        52582 ns/op (1.00x)   32657 B/op (1.00x)     346 allocs/op (1.00x)
   oj.ParseReader   27830 ns/op (1.89x)   22913 B/op (1.43x)     434 allocs/op (0.80x)
  gen.ParseReder    27796 ns/op (1.89x)   22912 B/op (1.43x)     434 allocs/op (0.80x)
  sen.ParseReader   28031 ns/op (1.88x)   22833 B/op (1.43x)     428 allocs/op (0.81x)

Validate string/[]byte
 json.Valid          1407 ns/op (1.00x)       0 B/op (NaNx)       0 allocs/op (NaNx)
   oj.Valdate         668 ns/op (2.11x)      64 B/op (0.00x)       2 allocs/op (0.00x)

Validate io.Reader
 json.Decode        52831 ns/op (1.00x)   32655 B/op (1.00x)     346 allocs/op (1.00x)
   oj.Valdate        2177 ns/op (24.27x)    4161 B/op (7.85x)       3 allocs/op (115.33x)

to JSON
 json.Marshal        2684 ns/op (1.00x)     992 B/op (1.00x)      22 allocs/op (1.00x)
   oj.JSON            447 ns/op (6.00x)     131 B/op (7.57x)       4 allocs/op (5.50x)
   oj.Write           464 ns/op (5.78x)     131 B/op (7.57x)       4 allocs/op (5.50x)
  sen.String          447 ns/op (6.00x)     131 B/op (7.57x)       4 allocs/op (5.50x)

to JSON with indentation
 json.Marshal        3849 ns/op (1.00x)    1488 B/op (1.00x)      25 allocs/op (1.00x)
   oj.JSON            553 ns/op (6.96x)     179 B/op (8.31x)       4 allocs/op (6.25x)
   oj.Write           576 ns/op (6.68x)     179 B/op (8.31x)       4 allocs/op (6.25x)
  sen.String          540 ns/op (7.13x)     163 B/op (9.13x)       4 allocs/op (6.25x)

to JSON with indentation and sorted keys
   oj.JSON            897 ns/op (1.00x)     323 B/op (1.00x)       8 allocs/op (1.00x)
   oj.Write           918 ns/op (0.98x)     323 B/op (1.00x)       8 allocs/op (1.00x)
  sen.String          866 ns/op (1.04x)     307 B/op (1.05x)       8 allocs/op (1.00x)

Convert or Alter
  alt.Generify       1220 ns/op (1.00x)    1712 B/op (1.00x)      27 allocs/op (1.00x)
  alt.Alter           916 ns/op (1.33x)     960 B/op (1.78x)      19 allocs/op (1.42x)

JSONPath Get $..a[2].c
   jp.Get            3814 ns/op (1.00x)    1992 B/op (1.00x)      65 allocs/op (1.00x)

JSONPath First  $..a[2].c
   jp.First          1939 ns/op (1.00x)    1272 B/op (1.00x)      32 allocs/op (1.00x)
```
