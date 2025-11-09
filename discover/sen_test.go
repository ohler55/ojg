// Copyright (c) 2025, Peter Ohler, All rights reserved.

package discover_test

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/ohler55/ojg/discover"
	"github.com/ohler55/ojg/pretty"
	"github.com/ohler55/ojg/tt"
)

func TestSENOkay(t *testing.T) {
	var found any
	discover.SEN([]byte(` [ [[] ] ] `), func(f any) (stop bool) {
		found = f
		return false
	})
	tt.Equal(t, `[[[]]]`, pretty.SEN(found))
}

func TestSENBack(t *testing.T) {
	var found any
	discover.SEN([]byte(` [ }[[] ] ] `), func(f any) (stop bool) {
		found = f
		return false
	})
	tt.Equal(t, `[[]]`, pretty.SEN(found))
}

func TestSENBadSEN(t *testing.T) {
	var found any
	discover.SEN([]byte(` [ 12x3 [4 ]] `), func(f any) (stop bool) {
		found = f
		return false
	})
	tt.Equal(t, `[4]`, pretty.SEN(found))
}

func TestSENShort(t *testing.T) {
	var found any
	discover.SEN([]byte(` [ 123 [4 ] `), func(f any) (stop bool) {
		found = f
		return false
	})
	tt.Equal(t, `[4]`, pretty.SEN(found))
}

func TestReadSEN(t *testing.T) {
	var found any
	r, err := os.Open("testdata/with-sen.txt")
	defer func() { _ = r.Close() }()
	tt.Nil(t, err)
	discover.ReadSEN(r, func(f any) (stop bool) {
		found = f
		return false
	})
	tt.Equal(t, `{abc: [1 2 3 4] xyz: [5 6 7 8]}`, pretty.SEN(found))
}

func TestReadSENBack(t *testing.T) {
	var found any
	r := strings.NewReader("start here [1x2] [1 2 3] end")
	discover.ReadSEN(r, func(f any) (stop bool) {
		found = f
		return false
	})
	tt.Equal(t, `[1 2 3]`, pretty.SEN(found))
}

type badReader int

func (w badReader) Read([]byte) (int, error) {
	return 0, fmt.Errorf("oops")
}

func TestReadSENError(t *testing.T) {
	tt.Panic(t, func() { discover.ReadSEN(badReader(0), func(_ any) bool { return false }) })
}
