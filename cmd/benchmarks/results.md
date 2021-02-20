## Benchmark Results

```
Parse string/[]byte
       json.Unmarshal       40143 ns/op   17921 B/op     334 allocs/op
         oj.Parse           22221 ns/op   18768 B/op     429 allocs/op
   oj-reuse.Parse           16864 ns/op    5968 B/op     364 allocs/op
        gen.Parse           22158 ns/op   18769 B/op     429 allocs/op
  gen-reuse.Parse           16959 ns/op    5968 B/op     364 allocs/op
        sen.Parse           22646 ns/op   18784 B/op     431 allocs/op
  sen-reuse.Parse           17544 ns/op    5968 B/op     366 allocs/op

   oj-reuse.Parse        ███████████████████████▊ 2.38
  gen-reuse.Parse        ███████████████████████▋ 2.37
  sen-reuse.Parse        ██████████████████████▉ 2.29
        gen.Parse        ██████████████████  1.81
         oj.Parse        ██████████████████  1.81
        sen.Parse        █████████████████▋ 1.77
       json.Unmarshal    ▓▓▓▓▓▓▓▓▓▓ 1.00

Parse io.Reader
       json.Decode          50819 ns/op   32593 B/op     344 allocs/op
         oj.ParseReader     26042 ns/op   22863 B/op     430 allocs/op
   oj-reuse.ParseReader     20696 ns/op   10064 B/op     365 allocs/op
        gen.ParseReder      26141 ns/op   22865 B/op     430 allocs/op
  gen-reuse.ParseReder      20796 ns/op   10064 B/op     365 allocs/op
        sen.ParseReader     26748 ns/op   22863 B/op     432 allocs/op
  sen-reuse.ParseReader     21336 ns/op   10064 B/op     367 allocs/op

   oj-reuse.ParseReader  ████████████████████████▌ 2.46
  gen-reuse.ParseReder   ████████████████████████▍ 2.44
  sen-reuse.ParseReader  ███████████████████████▊ 2.38
         oj.ParseReader  ███████████████████▌ 1.95
        gen.ParseReder   ███████████████████▍ 1.94
        sen.ParseReader  ██████████████████▉ 1.90
       json.Decode       ▓▓▓▓▓▓▓▓▓▓ 1.00

Parse chan interface{}
       json.Parse-chan      44783 ns/op   17936 B/op     335 allocs/op
         oj.Parse           25568 ns/op   18768 B/op     429 allocs/op
        gen.Parse           25684 ns/op   18767 B/op     429 allocs/op
        sen.Parse           26480 ns/op   18768 B/op     431 allocs/op

         oj.Parse        █████████████████▌ 1.75
        gen.Parse        █████████████████▍ 1.74
        sen.Parse        ████████████████▉ 1.69
       json.Parse-chan   ▓▓▓▓▓▓▓▓▓▓ 1.00

Validate string/[]byte
       json.Valid           11287 ns/op       0 B/op       0 allocs/op
         oj.Valdate          3774 ns/op       0 B/op       0 allocs/op

         oj.Valdate      █████████████████████████████▉ 2.99
       json.Valid        ▓▓▓▓▓▓▓▓▓▓ 1.00

Validate io.Reader
       json.Decode          51106 ns/op   32592 B/op     344 allocs/op
         oj.Valdate          7220 ns/op    4096 B/op       1 allocs/op

         oj.Valdate      ██████████████████████████████████████████████████████████████████████▊ 7.08
       json.Decode       ▓▓▓▓▓▓▓▓▓▓ 1.00

to JSON
       json.Marshal         39872 ns/op   17909 B/op     345 allocs/op
         oj.JSON             9525 ns/op    3072 B/op       1 allocs/op
        sen.String           9745 ns/op    2304 B/op       1 allocs/op

         oj.JSON         █████████████████████████████████████████▊ 4.19
        sen.String       ████████████████████████████████████████▉ 4.09
       json.Marshal      ▓▓▓▓▓▓▓▓▓▓ 1.00

to JSON with indentation
       json.Marshal         68700 ns/op   27329 B/op     352 allocs/op
         oj.JSON            11274 ns/op    4096 B/op       1 allocs/op
        sen.String          12029 ns/op    3456 B/op       1 allocs/op
     pretty.JSON            35249 ns/op   36928 B/op     445 allocs/op
     pretty.SEN             32162 ns/op   29792 B/op     396 allocs/op

         oj.JSON         ████████████████████████████████████████████████████████████▉ 6.09
        sen.String       █████████████████████████████████████████████████████████  5.71
     pretty.SEN          █████████████████████▎ 2.14
     pretty.JSON         ███████████████████▍ 1.95
       json.Marshal      ▓▓▓▓▓▓▓▓▓▓ 1.00

to JSON with indentation and sorted keys
         oj.JSON            17844 ns/op    6560 B/op      63 allocs/op
        sen.String          17503 ns/op    5920 B/op      63 allocs/op
     pretty.JSON            35168 ns/op   36928 B/op     445 allocs/op
     pretty.SEN             32729 ns/op   29792 B/op     396 allocs/op

        sen.String       ██████████▏ 1.02
         oj.JSON         ▓▓▓▓▓▓▓▓▓▓ 1.00
     pretty.SEN          █████▍ 0.55
     pretty.JSON         █████  0.51

Write indented JSON
       json.Encode          69805 ns/op   28384 B/op     353 allocs/op
         oj.Write           10189 ns/op       0 B/op       0 allocs/op
        sen.Write            8674 ns/op       0 B/op       0 allocs/op
     pretty.WriteJSON       31392 ns/op   17984 B/op     436 allocs/op
     pretty.WriteSEN        29628 ns/op   16736 B/op     388 allocs/op

        sen.Write        ████████████████████████████████████████████████████████████████████████████████▍ 8.05
         oj.Write        ████████████████████████████████████████████████████████████████████▌ 6.85
     pretty.WriteSEN     ███████████████████████▌ 2.36
     pretty.WriteJSON    ██████████████████████▏ 2.22
       json.Encode       ▓▓▓▓▓▓▓▓▓▓ 1.00

Convert or Alter
        alt.Generify         1215 ns/op    1696 B/op      25 allocs/op
        alt.Alter             900 ns/op     944 B/op      17 allocs/op

        alt.Alter        █████████████▌ 1.35
        alt.Generify     ▓▓▓▓▓▓▓▓▓▓ 1.00

JSONPath Get $..a[2].c
         jp.Get            264494 ns/op   19288 B/op    2227 allocs/op

         jp.Get          ▓▓▓▓▓▓▓▓▓▓ 1.00

JSONPath First  $..a[2].c
         jp.First           24476 ns/op    2880 B/op     233 allocs/op

         jp.First        ▓▓▓▓▓▓▓▓▓▓ 1.00

 Higher values (longer bars) are better in all cases. The bar graph compares the
 parsing performance. The lighter colored bar is the reference, usually the go
 json package.

 The Benchmarks reflect a use case where JSON is either provided as a string or
 read from a file (io.Reader) then parsed into simple go types of nil, bool, int64
 float64, string, []interface{}, or map[string]interface{}. When supported, an
 io.Writer benchmark is also included along with some miscellaneous operations.

Tests run on:
 Machine:         MacBookPro15,2
 OS:              macOS 11.2.1
 Processor:       Quad-Core Intel Core i7
 Cores:           4
 Processor Speed: 2.8 GHz
```
