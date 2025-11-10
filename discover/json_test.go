// Copyright (c) 2025, Peter Ohler, All rights reserved.

package discover_test

import (
	"strings"
	"testing"

	"github.com/ohler55/ojg/discover"
	"github.com/ohler55/ojg/pretty"
	"github.com/ohler55/ojg/tt"
)

func TestJSONOkay(t *testing.T) {
	var found any
	discover.JSON([]byte(` [ [[] ] ] `), func(f any) (stop bool) {
		found = f
		return false
	})
	tt.Equal(t, `[[[]]]`, pretty.SEN(found))
}

func TestJSONBack(t *testing.T) {
	var found any
	discover.JSON([]byte(` [ }[[] ] ] `), func(f any) (stop bool) {
		found = f
		return false
	})
	tt.Equal(t, `[[]]`, pretty.SEN(found))
}

func TestJSONBadJSON(t *testing.T) {
	var found any
	discover.JSON([]byte(` [ 12x3 [4 ]] `), func(f any) (stop bool) {
		found = f
		return false
	})
	tt.Equal(t, `[4]`, pretty.SEN(found))
}

func TestJSONShort(t *testing.T) {
	var found any
	discover.JSON([]byte(` [ 123 [4 ] `), func(f any) (stop bool) {
		found = f
		return false
	})
	tt.Equal(t, `[4]`, pretty.SEN(found))
}

func TestReadJSON(t *testing.T) {
	var found any
	r := strings.NewReader(`start here [1,  "two", 3] end`)
	discover.ReadJSON(r, func(f any) (stop bool) {
		found = f
		return false
	})
	tt.Equal(t, `[1 two 3]`, pretty.SEN(found))
}

func TestReadJSONBack(t *testing.T) {
	var found any
	r := strings.NewReader("start here [1x2] [1, 2, 3] end")
	discover.ReadJSON(r, func(f any) (stop bool) {
		found = f
		return false
	})
	tt.Equal(t, `[1 2 3]`, pretty.SEN(found))
}

func TestReadJSONError(t *testing.T) {
	tt.Panic(t, func() { discover.ReadJSON(badReader(0), func(_ any) bool { return false }) })
}
