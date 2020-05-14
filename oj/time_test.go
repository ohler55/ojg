// Copyright (c) 2020, Peter Ohler, All rights reserved.

package oj_test

import (
	"testing"
	"time"

	"github.com/ohler55/ojg/oj"
	"github.com/ohler55/ojg/tt"
)

func TestTimeString(t *testing.T) {
	n := oj.Time(time.Date(2020, time.April, 12, 16, 34, 04, 123456789, time.UTC))

	oj.TimeWrap = "@"
	oj.TimeFormat = "nano"
	tt.Equal(t, `{"@":1586709244123456789}`, n.String())

	oj.TimeFormat = time.RFC3339Nano
	tt.Equal(t, `{"@":"2020-04-12T16:34:04.123456789Z"}`, n.String())

	oj.TimeFormat = "second"
	tt.Equal(t, `{"@":1586709244.123456789}`, n.String())

	oj.TimeWrap = ""
	oj.TimeFormat = "nano"
	tt.Equal(t, "1586709244123456789", n.String())

	oj.TimeFormat = time.RFC3339Nano
	tt.Equal(t, `"2020-04-12T16:34:04.123456789Z"`, n.String())

	oj.TimeFormat = "second"
	tt.Equal(t, "1586709244.123456789", n.String())
}
