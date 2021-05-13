// Copyright (c) 2020, Peter Ohler, All rights reserved.

package alt_test

import (
	"fmt"
	"time"

	"github.com/ohler55/ojg/alt"
)

func ExampleBool() {
	for _, src := range []interface{}{true, "tRuE", "x", 1, nil} {
		fmt.Printf("alt.Bool(%T(%v)) = %t  alt.Bool(%T(%v), false) = %t   alt.Bool(%T(%v), false, true) = %t\n",
			src, src, alt.Bool(src),
			src, src, alt.Bool(src, false),
			src, src, alt.Bool(src, false, true))
	}
	// Output:
	// alt.Bool(bool(true)) = true  alt.Bool(bool(true), false) = true   alt.Bool(bool(true), false, true) = true
	// alt.Bool(string(tRuE)) = true  alt.Bool(string(tRuE), false) = true   alt.Bool(string(tRuE), false, true) = true
	// alt.Bool(string(x)) = false  alt.Bool(string(x), false) = false   alt.Bool(string(x), false, true) = true
	// alt.Bool(int(1)) = false  alt.Bool(int(1), false) = false   alt.Bool(int(1), false, true) = false
	// alt.Bool(<nil>(<nil>)) = false  alt.Bool(<nil>(<nil>), false) = false   alt.Bool(<nil>(<nil>), false, true) = true
}

func ExampleInt() {
	for _, src := range []interface{}{1, "1", "x", 1.5, []interface{}{}} {
		fmt.Printf("alt.Int(%T(%v)) = %d  alt.Int(%T(%v), 2) = %d   alt.Int(%T(%v), 2, 3) = %d\n",
			src, src, alt.Int(src),
			src, src, alt.Int(src, 2),
			src, src, alt.Int(src, 2, 3))
	}
	// Output:
	// alt.Int(int(1)) = 1  alt.Int(int(1), 2) = 1   alt.Int(int(1), 2, 3) = 1
	// alt.Int(string(1)) = 1  alt.Int(string(1), 2) = 1   alt.Int(string(1), 2, 3) = 3
	// alt.Int(string(x)) = 0  alt.Int(string(x), 2) = 2   alt.Int(string(x), 2, 3) = 3
	// alt.Int(float64(1.5)) = 1  alt.Int(float64(1.5), 2) = 1   alt.Int(float64(1.5), 2, 3) = 3
	// alt.Int([]interface {}([])) = 0  alt.Int([]interface {}([]), 2) = 2   alt.Int([]interface {}([]), 2, 3) = 2
}

func ExampleFloat() {
	for _, src := range []interface{}{1, "1,5", "x", 1.5, true} {
		fmt.Printf("alt.Float(%T(%v)) = %.1f  alt.Float(%T(%v), 2.5) = %.1f   alt.Float(%T(%v), 2.5, 3.5) = %.1f\n",
			src, src, alt.Float(src),
			src, src, alt.Float(src, 2.5),
			src, src, alt.Float(src, 2.5, 3.5))
	}
	// Output:
	// alt.Float(int(1)) = 1.0  alt.Float(int(1), 2.5) = 1.0   alt.Float(int(1), 2.5, 3.5) = 3.5
	// alt.Float(string(1,5)) = 0.0  alt.Float(string(1,5), 2.5) = 2.5   alt.Float(string(1,5), 2.5, 3.5) = 3.5
	// alt.Float(string(x)) = 0.0  alt.Float(string(x), 2.5) = 2.5   alt.Float(string(x), 2.5, 3.5) = 3.5
	// alt.Float(float64(1.5)) = 1.5  alt.Float(float64(1.5), 2.5) = 1.5   alt.Float(float64(1.5), 2.5, 3.5) = 1.5
	// alt.Float(bool(true)) = 0.0  alt.Float(bool(true), 2.5) = 2.5   alt.Float(bool(true), 2.5, 3.5) = 3.5
}

func ExampleString() {
	tm := time.Date(2021, time.February, 9, 12, 13, 14, 0, time.UTC)
	for _, src := range []interface{}{"xyz", 1, 1.5, true, tm, []interface{}{}} {
		fmt.Printf("alt.String(%T(%v)) = %s  alt.String(%T(%v), default) = %s   alt.String(%T(%v), default, picky) = %s\n",
			src, src, alt.String(src),
			src, src, alt.String(src, "default"),
			src, src, alt.String(src, "default", "picky"))
	}
	// Output:
	// alt.String(string(xyz)) = xyz  alt.String(string(xyz), default) = xyz   alt.String(string(xyz), default, picky) = xyz
	// alt.String(int(1)) = 1  alt.String(int(1), default) = 1   alt.String(int(1), default, picky) = picky
	// alt.String(float64(1.5)) = 1.5  alt.String(float64(1.5), default) = 1.5   alt.String(float64(1.5), default, picky) = picky
	// alt.String(bool(true)) = true  alt.String(bool(true), default) = true   alt.String(bool(true), default, picky) = picky
	// alt.String(time.Time(2021-02-09 12:13:14 +0000 UTC)) = 2021-02-09T12:13:14Z  alt.String(time.Time(2021-02-09 12:13:14 +0000 UTC), default) = 2021-02-09T12:13:14Z   alt.String(time.Time(2021-02-09 12:13:14 +0000 UTC), default, picky) = picky
	// alt.String([]interface {}([])) =   alt.String([]interface {}([]), default) = default   alt.String([]interface {}([]), default, picky) = picky
}

func ExampleTime() {
	tm := time.Date(2021, time.February, 9, 12, 13, 14, 0, time.UTC)
	td := time.Date(2021, time.January, 1, 0, 0, 0, 0, time.UTC) // default
	tp := time.Date(2000, time.January, 1, 0, 0, 0, 0, time.UTC) // picky default

	for _, src := range []interface{}{tm, 1612872711000000000, "2021-02-09T01:02:03Z", "x", 1612872722.0, true} {
		fmt.Printf("alt.Time(%T(%v)) = %s\nalt.Time(%T(%v), td) = %s\nalt.Time(%T(%v), td, tp) = %s\n",
			src, src, alt.Time(src).Format(time.RFC3339),
			src, src, alt.Time(src, td).Format(time.RFC3339),
			src, src, alt.Time(src, td, tp).Format(time.RFC3339))
	}
	// Output:
	// alt.Time(time.Time(2021-02-09 12:13:14 +0000 UTC)) = 2021-02-09T12:13:14Z
	// alt.Time(time.Time(2021-02-09 12:13:14 +0000 UTC), td) = 2021-02-09T12:13:14Z
	// alt.Time(time.Time(2021-02-09 12:13:14 +0000 UTC), td, tp) = 2021-02-09T12:13:14Z
	// alt.Time(int(1612872711000000000)) = 2021-02-09T12:11:51Z
	// alt.Time(int(1612872711000000000), td) = 2021-02-09T12:11:51Z
	// alt.Time(int(1612872711000000000), td, tp) = 2000-01-01T00:00:00Z
	// alt.Time(string(2021-02-09T01:02:03Z)) = 2021-02-09T01:02:03Z
	// alt.Time(string(2021-02-09T01:02:03Z), td) = 2021-02-09T01:02:03Z
	// alt.Time(string(2021-02-09T01:02:03Z), td, tp) = 2000-01-01T00:00:00Z
	// alt.Time(string(x)) = 0001-01-01T00:00:00Z
	// alt.Time(string(x), td) = 2021-01-01T00:00:00Z
	// alt.Time(string(x), td, tp) = 2000-01-01T00:00:00Z
	// alt.Time(float64(1.612872722e+09)) = 2021-02-09T12:12:02Z
	// alt.Time(float64(1.612872722e+09), td) = 2021-02-09T12:12:02Z
	// alt.Time(float64(1.612872722e+09), td, tp) = 2000-01-01T00:00:00Z
	// alt.Time(bool(true)) = 0001-01-01T00:00:00Z
	// alt.Time(bool(true), td) = 2021-01-01T00:00:00Z
	// alt.Time(bool(true), td, tp) = 2000-01-01T00:00:00Z
}
