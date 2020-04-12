// Copyright (c) 2020, Peter Ohler, All rights reserved.

package gd

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

var TimeFormat = time.RFC3339Nano

type Time time.Time

func (n Time) String() string {
	var b strings.Builder

	n.BuildJSON(&b)

	return b.String()
}

func (n Time) Alter() interface{} {
	return time.Time(n)
}

func (n Time) Native() interface{} {
	return time.Time(n)
}

func (n Time) Dup() Node {
	return n
}

func (n Time) JSON(_ ...int) string {
	var b strings.Builder

	n.BuildJSON(&b)

	return b.String()
}

func (n Time) BuildJSON(b *strings.Builder) {
	switch TimeFormat {
	case "", "nano":
		b.WriteString(strconv.FormatInt(time.Time(n).UnixNano(), 10))
	case "second":
		// Decimal format but float is not accurate enough so build the output
		// in two parts.
		nano := time.Time(n).UnixNano()
		secs := nano / int64(time.Second)
		if 0 < nano {
			b.WriteString(fmt.Sprintf("%d.%09d", secs, nano-(secs*int64(time.Second))))
		} else {
			b.WriteString(fmt.Sprintf("%d.%09d", secs, -nano-(secs*int64(time.Second))))
		}
	default:
		b.WriteString(`"`)
		b.WriteString(time.Time(n).Format(TimeFormat))
		b.WriteString(`"`)
	}
}
