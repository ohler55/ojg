// Copyright (c) 2021, Peter Ohler, All rights reserved.

package sen_test

import (
	"testing"

	"github.com/ohler55/ojg/sen"
	"github.com/ohler55/ojg/tt"
)

func TestUnmarshal(t *testing.T) {
	var obj map[string]interface{}
	src := `{x:3}`
	err := sen.Unmarshal([]byte(src), &obj)
	tt.Nil(t, err)
	tt.Equal(t, src, sen.String(obj))

	obj = nil
	p := sen.Parser{}
	err = p.Unmarshal([]byte(src), &obj)
	tt.Nil(t, err)
	tt.Equal(t, src, sen.String(obj))
}
