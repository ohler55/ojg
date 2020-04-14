// Copyright (c) 2020, Peter Ohler, All rights reserved.

package gd_test

import (
	"testing"
	"time"

	"github.com/ohler55/ojg/gd"
	"github.com/ohler55/ojg/tt"
)

func TestTimeString(t *testing.T) {
	n := gd.Time(time.Date(2020, time.April, 12, 16, 34, 04, 123456789, time.UTC))

	gd.TimeWrap = "@"
	gd.TimeFormat = "nano"
	tt.Equal(t, `{"@":1586709244123456789}`, n.String())

	gd.TimeFormat = time.RFC3339Nano
	tt.Equal(t, `{"@":"2020-04-12T16:34:04.123456789Z"}`, n.String())

	gd.TimeFormat = "second"
	tt.Equal(t, `{"@":1586709244.123456789}`, n.String())

	gd.TimeWrap = ""
	gd.TimeFormat = "nano"
	tt.Equal(t, "1586709244123456789", n.String())

	gd.TimeFormat = time.RFC3339Nano
	tt.Equal(t, `"2020-04-12T16:34:04.123456789Z"`, n.String())

	gd.TimeFormat = "second"
	tt.Equal(t, "1586709244.123456789", n.String())
}
