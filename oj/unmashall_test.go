// Copyright (c) 2021, Peter Ohler, All rights reserved.

package oj_test

import (
	"testing"

	"github.com/ohler55/ojg/oj"
	"github.com/ohler55/ojg/tt"
)

func TestUnmarshal(t *testing.T) {
	var obj map[string]interface{}
	src := `{"x":3}`
	err := oj.Unmarshal([]byte(src), &obj)
	tt.Nil(t, err)
	tt.Equal(t, src, oj.JSON(obj))

	obj = nil
	p := oj.Parser{}
	err = p.Unmarshal([]byte(src), &obj)
	tt.Nil(t, err)
	tt.Equal(t, src, oj.JSON(obj))
}
