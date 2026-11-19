# OjG Benchmarks

Benchmarks were run from the ojg/cmd/benchmark directory on 2026-11-19 with the command:

```
go run *.go
```

```

Parse string/[]byte
       json.Unmarshal           18330 ns/op    17812 B/op    389 allocs/op
         oj.Parse                7676 ns/op    17336 B/op    432 allocs/op
   oj-reuse.Parse                6225 ns/op     5688 B/op    364 allocs/op
        gen.Parse                7938 ns/op    17336 B/op    432 allocs/op
  gen-reuse.Parse                6476 ns/op     5688 B/op    364 allocs/op
        sen.Parse                8500 ns/op    17336 B/op    434 allocs/op
  sen-reuse.Parse                7129 ns/op     5672 B/op    366 allocs/op

   oj-reuse.Parse        ████████████████████▌ 2.94
  gen-reuse.Parse        ███████████████████▊ 2.83
  sen-reuse.Parse        █████████████████▉ 2.57
         oj.Parse        ████████████████▋ 2.39
        gen.Parse        ████████████████▏ 2.31
        sen.Parse        ███████████████  2.16
       json.Unmarshal    ▓▓▓▓▓▓▓ 1.00

Unmarshal []byte to type
       json.Unmarshal           10528 ns/op     3626 B/op     46 allocs/op
         oj.Unmarshal           16036 ns/op    10266 B/op    488 allocs/op
        sen.Unmarshal           16453 ns/op    10250 B/op    488 allocs/op

       json.Unmarshal    ▓▓▓▓▓▓▓ 1.00
         oj.Unmarshal    ████▌ 0.66
        sen.Unmarshal    ████▍ 0.64

Tokenize
       json.Decode               7477 ns/op     9944 B/op    338 allocs/op
         oj.Tokenize             4036 ns/op     1976 B/op    156 allocs/op
        sen.Tokenize             4450 ns/op     1976 B/op    158 allocs/op

         oj.Tokenize     ████████████▉ 1.85
        sen.Tokenize     ███████████▊ 1.68
       json.Decode       ▓▓▓▓▓▓▓ 1.00

Parse io.Reader
       json.Decode              24667 ns/op    26379 B/op    401 allocs/op
         oj.ParseReader          9014 ns/op    21432 B/op    433 allocs/op
   oj-reuse.ParseReader          7429 ns/op     9784 B/op    365 allocs/op
        gen.ParseReder           9374 ns/op    21432 B/op    433 allocs/op
  gen-reuse.ParseReder           7691 ns/op     9784 B/op    365 allocs/op
        sen.ParseReader          9320 ns/op    21432 B/op    435 allocs/op
  sen-reuse.ParseReader          7728 ns/op     9784 B/op    367 allocs/op
         oj.TokenizeLoad         5137 ns/op     6072 B/op    157 allocs/op
        sen.TokenizeLoad         5482 ns/op     6072 B/op    159 allocs/op

         oj.TokenizeLoad █████████████████████████████████▌ 4.80
        sen.TokenizeLoad ███████████████████████████████▍ 4.50
   oj-reuse.ParseReader  ███████████████████████▏ 3.32
  gen-reuse.ParseReder   ██████████████████████▍ 3.21
  sen-reuse.ParseReader  ██████████████████████▎ 3.19
         oj.ParseReader  ███████████████████▏ 2.74
        sen.ParseReader  ██████████████████▌ 2.65
        gen.ParseReder   ██████████████████▍ 2.63
       json.Decode       ▓▓▓▓▓▓▓ 1.00

Parse chan any
       json.Parse-chan          26058 ns/op    17835 B/op    390 allocs/op
         oj.Parse               11476 ns/op    17336 B/op    432 allocs/op
        gen.Parse               11764 ns/op    17336 B/op    432 allocs/op
        sen.Parse               12382 ns/op    17320 B/op    434 allocs/op

         oj.Parse        ███████████████▉ 2.27
        gen.Parse        ███████████████▌ 2.22
        sen.Parse        ██████████████▋ 2.10
       json.Parse-chan   ▓▓▓▓▓▓▓ 1.00

Validate string/[]byte
       json.Valid                1990 ns/op        0 B/op      0 allocs/op
         oj.Valdate              2300 ns/op        0 B/op      0 allocs/op

       json.Valid        ▓▓▓▓▓▓▓ 1.00
         oj.Valdate      ██████  0.87

Validate io.Reader
       json.Decode              24603 ns/op    26382 B/op    401 allocs/op
         oj.Valdate              3463 ns/op     4096 B/op      1 allocs/op

         oj.Valdate      █████████████████████████████████████████████████▋ 7.10
       json.Decode       ▓▓▓▓▓▓▓ 1.00

to JSON
       json.Marshal              6547 ns/op     3573 B/op     34 allocs/op
         oj.JSON                 3928 ns/op        0 B/op      0 allocs/op
        sen.SEN                  4473 ns/op        0 B/op      0 allocs/op

         oj.JSON         ███████████▋ 1.67
        sen.SEN          ██████████▏ 1.46
       json.Marshal      ▓▓▓▓▓▓▓ 1.00

to JSON with indentation
       json.Marshal             10334 ns/op     7676 B/op     35 allocs/op
         oj.JSON                 4294 ns/op        0 B/op      0 allocs/op
        sen.Bytes                4759 ns/op        0 B/op      0 allocs/op
     pretty.JSON                12092 ns/op    32888 B/op    393 allocs/op
     pretty.SEN                 11765 ns/op    26912 B/op    344 allocs/op

         oj.JSON         ████████████████▊ 2.41
        sen.Bytes        ███████████████▏ 2.17
       json.Marshal      ▓▓▓▓▓▓▓ 1.00
     pretty.SEN          ██████▏ 0.88
     pretty.JSON         █████▉ 0.85

to JSON with indentation and sorted keys
         oj.JSON                 5542 ns/op     1040 B/op     11 allocs/op
        sen.Bytes                6098 ns/op     1040 B/op     11 allocs/op
     pretty.JSON                12060 ns/op    32888 B/op    393 allocs/op
     pretty.SEN                 11804 ns/op    26912 B/op    344 allocs/op

         oj.JSON         ▓▓▓▓▓▓▓ 1.00
        sen.Bytes        ██████▎ 0.91
     pretty.SEN          ███▎ 0.47
     pretty.JSON         ███▏ 0.46

Write indented JSON
       json.Encode              10551 ns/op    10768 B/op     35 allocs/op
         oj.Write                4372 ns/op        0 B/op      0 allocs/op
        sen.Write                4866 ns/op        0 B/op      0 allocs/op
     pretty.WriteJSON           11098 ns/op    21624 B/op    390 allocs/op
     pretty.WriteSEN            11160 ns/op    18720 B/op    341 allocs/op

         oj.Write        ████████████████▉ 2.41
        sen.Write        ███████████████▏ 2.17
       json.Encode       ▓▓▓▓▓▓▓ 1.00
     pretty.WriteJSON    ██████▋ 0.95
     pretty.WriteSEN     ██████▌ 0.95

Marshal Struct
       json.Marshal              5120 ns/op     3458 B/op      1 allocs/op
         oj.Marshal              3836 ns/op     5171 B/op     45 allocs/op

         oj.Marshal      █████████▎ 1.33
       json.Marshal      ▓▓▓▓▓▓▓ 1.00

Convert or Alter
        alt.Generify              411 ns/op     1664 B/op     25 allocs/op
        alt.Alter                 325 ns/op      912 B/op     17 allocs/op

        alt.Alter        ████████▊ 1.26
        alt.Generify     ▓▓▓▓▓▓▓ 1.00

JSONPath Get $..a[2].c
         jp.Get                115346 ns/op    21721 B/op   2228 allocs/op

         jp.Get          ▓▓▓▓▓▓▓ 1.00

JSONPath First  $..a[2].c
         jp.First               10829 ns/op     5312 B/op    234 allocs/op

         jp.First        ▓▓▓▓▓▓▓ 1.00

 Higher values (longer bars) are better in all cases. The bar graph compares the
 parsing performance. The lighter colored bar is the reference, usually the go
 json package.

 The Benchmarks reflect a use case where JSON is either provided as a string or
 read from a file (io.Reader) then parsed into simple go types of nil, bool, int64
 float64, string, []any, or map[string]any. When supported, an
 io.Writer benchmark is also included along with some miscellaneous operations.

Tests run on:
 Machine:         Mac16,5
 OS:              macOS 27.2.1
 Processor:
 Cores:           proc 16:0:12:4
 Processor Speed:

```

Previous benchmark results are [benchmarks-2021](benchmarks-2021.md)
