# OjG Benchmarks

Benchmarks were run from the ojg/cmd/benchmark directory with the command:

```
go run *.go -cat
```

```
Parse string/[]byte
       json.Unmarshal        11364753 ns/op  5122778 B/op   95372 allocs/op
         oj.Parse            7597726 ns/op  5343960 B/op  112400 allocs/op
   oj-reuse.Parse            4573241 ns/op  1512773 B/op   90348 allocs/op
        gen.Parse            8096291 ns/op  5344006 B/op  112400 allocs/op
  gen-reuse.Parse            4769426 ns/op  1513386 B/op   90351 allocs/op
        sen.Parse            8363065 ns/op  5349046 B/op  113662 allocs/op
  sen-reuse.Parse            5045270 ns/op  1520293 B/op   91623 allocs/op

   oj-reuse.Parse        ██████████████▉ 2.49
  gen-reuse.Parse        ██████████████▎ 2.38
  sen-reuse.Parse        █████████████▌ 2.25
         oj.Parse        ████████▉ 1.50
        gen.Parse        ████████▍ 1.40
        sen.Parse        ████████▏ 1.36
       json.Unmarshal    ▓▓▓▓▓▓ 1.00

Unmarshal []byte to type
       json.Unmarshal        12621718 ns/op  1214596 B/op   44157 allocs/op
         oj.Unmarshal        9380795 ns/op  2309873 B/op  130351 allocs/op
        sen.Unmarshal        9670015 ns/op  2315838 B/op  131620 allocs/op

         oj.Unmarshal    ████████  1.35
        sen.Unmarshal    ███████▊ 1.31
       json.Unmarshal    ▓▓▓▓▓▓ 1.00

Tokenize
       json.Decode           20282177 ns/op  5333749 B/op  323430 allocs/op
         oj.Tokenize         2040074 ns/op  377010 B/op   40996 allocs/op
        sen.Tokenize         2144064 ns/op  381970 B/op   42259 allocs/op

         oj.Tokenize     ███████████████████████████████████████████████████████████▋ 9.94
        sen.Tokenize     ████████████████████████████████████████████████████████▊ 9.46
       json.Decode       ▓▓▓▓▓▓ 1.00

Parse io.Reader
       json.Decode           19287966 ns/op  9315221 B/op   95390 allocs/op
         oj.ParseReader      6833773 ns/op  5348105 B/op  112401 allocs/op
   oj-reuse.ParseReader      4896066 ns/op  1517777 B/op   90354 allocs/op
        gen.ParseReder       7045802 ns/op  5348007 B/op  112401 allocs/op
  gen-reuse.ParseReder       5156753 ns/op  1518624 B/op   90358 allocs/op
        sen.ParseReader      7140175 ns/op  5353030 B/op  113664 allocs/op
  sen-reuse.ParseReader      5373338 ns/op  1523618 B/op   91622 allocs/op
         oj.TokenizeLoad     2267153 ns/op  381105 B/op   40997 allocs/op
        sen.TokenizeLoad     2334417 ns/op  386065 B/op   42260 allocs/op

         oj.TokenizeLoad ███████████████████████████████████████████████████  8.51
        sen.TokenizeLoad █████████████████████████████████████████████████▌ 8.26
   oj-reuse.ParseReader  ███████████████████████▋ 3.94
  gen-reuse.ParseReder   ██████████████████████▍ 3.74
  sen-reuse.ParseReader  █████████████████████▌ 3.59
         oj.ParseReader  ████████████████▉ 2.82
        gen.ParseReder   ████████████████▍ 2.74
        sen.ParseReader  ████████████████▏ 2.70
       json.Decode       ▓▓▓▓▓▓ 1.00

Parse chan interface{}
       json.Parse-chan       11778186 ns/op  5122740 B/op   95373 allocs/op
         oj.Parse            8636017 ns/op  5344067 B/op  112400 allocs/op
        gen.Parse            7922710 ns/op  5343875 B/op  112399 allocs/op
        sen.Parse            8320907 ns/op  5349322 B/op  113663 allocs/op

        gen.Parse        ████████▉ 1.49
        sen.Parse        ████████▍ 1.42
         oj.Parse        ████████▏ 1.36
       json.Parse-chan   ▓▓▓▓▓▓ 1.00

Validate string/[]byte
       json.Valid            3099171 ns/op       4 B/op       0 allocs/op
         oj.Valdate          1098024 ns/op       0 B/op       0 allocs/op

         oj.Valdate      ████████████████▉ 2.82
       json.Valid        ▓▓▓▓▓▓ 1.00

Validate io.Reader
       json.Decode           14468919 ns/op  9315299 B/op   95390 allocs/op
         oj.Valdate          1302676 ns/op    4096 B/op       1 allocs/op

         oj.Valdate      ██████████████████████████████████████████████████████████████████▋ 11.11
       json.Decode       ▓▓▓▓▓▓ 1.00

to JSON
       json.Marshal          11394357 ns/op  5391228 B/op  117351 allocs/op
         oj.JSON             1480108 ns/op    2963 B/op       0 allocs/op
        sen.SEN              1685107 ns/op    3318 B/op       0 allocs/op

         oj.JSON         ██████████████████████████████████████████████▏ 7.70
        sen.SEN          ████████████████████████████████████████▌ 6.76
       json.Marshal      ▓▓▓▓▓▓ 1.00

to JSON with indentation
       json.Marshal          17593380 ns/op  10105276 B/op  117377 allocs/op
         oj.JSON             1823519 ns/op    9104 B/op       0 allocs/op
        sen.Bytes            2022058 ns/op   10377 B/op       0 allocs/op
     pretty.JSON             9610741 ns/op  10995943 B/op  150410 allocs/op
     pretty.SEN              9178997 ns/op  9692717 B/op  129021 allocs/op

         oj.JSON         █████████████████████████████████████████████████████████▉ 9.65
        sen.Bytes        ████████████████████████████████████████████████████▏ 8.70
     pretty.SEN          ███████████▌ 1.92
     pretty.JSON         ██████████▉ 1.83
       json.Marshal      ▓▓▓▓▓▓ 1.00

to JSON with indentation and sorted keys
         oj.JSON             3611825 ns/op  694320 B/op   21872 allocs/op
        sen.Bytes            3893618 ns/op  695653 B/op   21872 allocs/op
     pretty.JSON             9639054 ns/op  10995949 B/op  150410 allocs/op
     pretty.SEN              9163758 ns/op  9692712 B/op  129021 allocs/op

         oj.JSON         ▓▓▓▓▓▓ 1.00
        sen.Bytes        █████▌ 0.93
     pretty.SEN          ██▎ 0.39
     pretty.JSON         ██▏ 0.37

Write indented JSON
       json.Encode           17836094 ns/op  10385101 B/op  117372 allocs/op
         oj.Write            1850944 ns/op       4 B/op       0 allocs/op
        sen.Write            2038636 ns/op       5 B/op       0 allocs/op
     pretty.WriteJSON        8613925 ns/op  5570019 B/op  150384 allocs/op
     pretty.WriteSEN         8414559 ns/op  5307171 B/op  128996 allocs/op

         oj.Write        █████████████████████████████████████████████████████████▊ 9.64
        sen.Write        ████████████████████████████████████████████████████▍ 8.75
     pretty.WriteSEN     ████████████▋ 2.12
     pretty.WriteJSON    ████████████▍ 2.07
       json.Encode       ▓▓▓▓▓▓ 1.00

Marshal Struct
       json.Marshal          1688290 ns/op  539891 B/op     453 allocs/op
         oj.Marshal          1490603 ns/op  282262 B/op   12238 allocs/op

         oj.Marshal      ██████▊ 1.13
       json.Marshal      ▓▓▓▓▓▓ 1.00

Convert or Alter
        alt.Generify           2340 ns/op    1664 B/op      25 allocs/op
        alt.Alter              1830 ns/op     912 B/op      17 allocs/op

        alt.Alter        ███████▋ 1.28
        alt.Generify     ▓▓▓▓▓▓ 1.00

JSONPath Get $..a[2].c
         jp.Get              250864 ns/op   19288 B/op    2227 allocs/op

         jp.Get          ▓▓▓▓▓▓ 1.00

JSONPath First  $..a[2].c
         jp.First             23026 ns/op    2880 B/op     233 allocs/op

         jp.First        ▓▓▓▓▓▓ 1.00

 Higher values (longer bars) are better in all cases. The bar graph compares the
 parsing performance. The lighter colored bar is the reference, usually the go
 json package.

 The Benchmarks reflect a use case where JSON is either provided as a string or
 read from a file (io.Reader) then parsed into simple go types of nil, bool, int64
 float64, string, []interface{}, or map[string]interface{}. When supported, an
 io.Writer benchmark is also included along with some miscellaneous operations.

Tests run on:
 OS:              Ubuntu 20.04.2 LTS
 Processor:       Intel(R) Core(TM) i7-8700 CPU
 Cores:           12
 Processor Speed: 3.20GHz
```
