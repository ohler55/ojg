# OjG Benchmarks

Benchmarks were run on a MacBook Pro with a 2.8 GHz Quad-core
I7. Higher numbers in parenthesis are better in all cases.

```
 The number in parenthesis are the ratio of results between the reference and
 the listed. Higher values are better.

Parse string/[]byte
 json.Unmarshal      7291 ns/op (1.00x)    4808 B/op (1.00x)      90 allocs/op (1.00x)
   oj.Parse          4775 ns/op (1.53x)    3984 B/op (1.21x)      86 allocs/op (1.05x)
  gen.Parse          4922 ns/op (1.48x)    3984 B/op (1.21x)      86 allocs/op (1.05x)
  sen.Parse          4713 ns/op (1.55x)    3984 B/op (1.21x)      86 allocs/op (1.05x)

Parse io.Reader
 json.Decode        52665 ns/op (1.00x)   32656 B/op (1.00x)     346 allocs/op (1.00x)
   oj.ParseReader   34116 ns/op (1.54x)   22913 B/op (1.43x)     434 allocs/op (0.80x)
  gen.ParseReder    33256 ns/op (1.58x)   22912 B/op (1.43x)     434 allocs/op (0.80x)
  sen.ParseReader   31718 ns/op (1.66x)   22832 B/op (1.43x)     428 allocs/op (0.81x)

Validate string/[]byte
 json.Valid          1409 ns/op (1.00x)       0 B/op (NaNx)       0 allocs/op (NaNx)
   oj.Valdate        1303 ns/op (1.08x)       0 B/op (NaNx)       0 allocs/op (NaNx)

Validate io.Reader
 json.Decode        52631 ns/op (1.00x)   32656 B/op (1.00x)     346 allocs/op (1.00x)
   oj.Valdate       15786 ns/op (3.33x)    4096 B/op (7.97x)       1 allocs/op (346.00x)

to JSON
 json.Marshal        2746 ns/op (1.00x)     992 B/op (1.00x)      22 allocs/op (1.00x)
   oj.JSON            445 ns/op (6.17x)     131 B/op (7.57x)       4 allocs/op (5.50x)
   oj.Write           460 ns/op (5.97x)     131 B/op (7.57x)       4 allocs/op (5.50x)
  sen.String          445 ns/op (6.17x)     131 B/op (7.57x)       4 allocs/op (5.50x)

to JSON with indentation
 json.Marshal        3820 ns/op (1.00x)    1488 B/op (1.00x)      25 allocs/op (1.00x)
   oj.JSON            555 ns/op (6.88x)     179 B/op (8.31x)       4 allocs/op (6.25x)
   oj.Write           570 ns/op (6.70x)     179 B/op (8.31x)       4 allocs/op (6.25x)
  sen.String          543 ns/op (7.03x)     163 B/op (9.13x)       4 allocs/op (6.25x)

to JSON with indentation and sorted keys
   oj.JSON            880 ns/op (1.00x)     323 B/op (1.00x)       8 allocs/op (1.00x)
   oj.Write           905 ns/op (0.97x)     323 B/op (1.00x)       8 allocs/op (1.00x)
  sen.String          879 ns/op (1.00x)     307 B/op (1.05x)       8 allocs/op (1.00x)

Convert or Alter
  alt.Generify       1245 ns/op (1.00x)    1712 B/op (1.00x)      27 allocs/op (1.00x)
  alt.Alter           912 ns/op (1.37x)     960 B/op (1.78x)      19 allocs/op (1.42x)

JSONPath Get $..a[2].c
   jp.Get            3704 ns/op (1.00x)    1992 B/op (1.00x)      65 allocs/op (1.00x)

JSONPath First  $..a[2].c
   jp.First          2077 ns/op (1.00x)    1272 B/op (1.00x)      32 allocs/op (1.00x)
```
