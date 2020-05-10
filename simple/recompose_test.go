// Copyright (c) 2020, Peter Ohler, All rights reserved.

package simple_test

import (
	"testing"

	"github.com/ohler55/ojg/simple"
	"github.com/ohler55/ojg/tt"
)

func sillyRecompose(data map[string]interface{}) (interface{}, error) {
	s := silly{}
	i, _ := data["val"].(int64)
	s.val = int(i)
	return &s, nil
}

func TestSimpleRecompose(t *testing.T) {
	src := map[string]interface{}{"type": "Dummy", "val": 3}
	r, err := simple.NewRecomposer("type", map[interface{}]simple.RecomposeFunc{&Dummy{}: sillyRecompose})
	tt.Nil(t, err, "NewRecomposer")
	var v interface{}
	v, err = r.Recompose(src)
	tt.Nil(t, err, "Recompose")
	s, _ := v.(*silly)
	tt.NotNil(t, s, "silly")
	tt.Equal(t, 3, s.val)
}
