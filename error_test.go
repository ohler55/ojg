// Copyright (c) 2021, Peter Ohler, All rights reserved.

package ojg_test

import (
	"strings"
	"testing"

	"github.com/ohler55/ojg"
	"github.com/ohler55/ojg/tt"
)

func TestError(t *testing.T) {
	err := ojg.NewError("some error")
	tt.Equal(t, "some error", err.Error())

	lines := strings.Split(string(err.Stack()), "\n")
	tt.Equal(t, true, strings.Contains(lines[0], "goroutine"))
	tt.Equal(t, true, strings.Contains(lines[len(lines)-2], "testing.go"))

	ojg.ErrorWithStack = true
	lines = strings.Split(err.Error(), "\n")
	tt.Equal(t, true, strings.Contains(lines[0], "some error"))
	tt.Equal(t, true, strings.Contains(lines[len(lines)-2], "testing.go"))
}
