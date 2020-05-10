// Copyright (c) 2020, Peter Ohler, All rights reserved.

package gen_test

import (
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
}
