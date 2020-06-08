// Copyright (c) 2020, Peter Ohler, All rights reserved.

package gen_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/ohler55/ojg/gen"
	"github.com/ohler55/ojg/tt"
)

func TestTimeString(t *testing.T) {
	n := gen.Time(time.Date(2020, time.April, 12, 16, 34, 04, 123456789, time.UTC))

	gen.TimeWrap = "@"
	gen.TimeFormat = "nano"
	tt.Equal(t, `{"@":1586709244123456789}`, n.String())

	gen.TimeFormat = time.RFC3339Nano
	tt.Equal(t, `{"@":"2020-04-12T16:34:04.123456789Z"}`, n.String())

	gen.TimeFormat = "second"
	tt.Equal(t, `{"@":1586709244.123456789}`, n.String())

	gen.TimeWrap = ""
	gen.TimeFormat = "nano"
	tt.Equal(t, "1586709244123456789", n.String())

	gen.TimeFormat = time.RFC3339Nano
	tt.Equal(t, `"2020-04-12T16:34:04.123456789Z"`, n.String())

	gen.TimeFormat = "second"
	tt.Equal(t, "1586709244.123456789", n.String())
	n = gen.Time(time.Date(1888, time.April, 12, 16, 34, 04, 123456789, time.UTC))
	tt.Equal(t, "-2578807555.876543211", n.String())
}

func TestTimeSimplify(t *testing.T) {
	n := gen.Time(time.Date(2020, time.April, 12, 16, 34, 04, 123456789, time.UTC))
	simple := n.Simplify()
	tt.Equal(t, "time.Time 2020-04-12 16:34:04.123456789 +0000 UTC", fmt.Sprintf("%T %v", simple, simple))
}

func TestTimeAlter(t *testing.T) {
	n := gen.Time(time.Date(2020, time.April, 12, 16, 34, 04, 123456789, time.UTC))
	alt := n.Alter()
	tt.Equal(t, "time.Time 2020-04-12 16:34:04.123456789 +0000 UTC", fmt.Sprintf("%T %v", alt, alt))
}

func TestTimeDup(t *testing.T) {
	gen.TimeFormat = time.RFC3339Nano
	n := gen.Time(time.Date(2020, time.April, 12, 16, 34, 04, 123456789, time.UTC))
	dup := n.Dup()
	tt.NotNil(t, dup)
	tt.Equal(t, `"2020-04-12T16:34:04.123456789Z"`, dup.String())
}

func TestTimeEmpty(t *testing.T) {
	n := gen.Time(time.Date(2020, time.April, 12, 16, 34, 04, 123456789, time.UTC))
	tt.Equal(t, false, n.Empty())
}
