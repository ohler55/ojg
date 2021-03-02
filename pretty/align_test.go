// Copyright (c) 2021, Peter Ohler, All rights reserved.

package pretty_test

import (
	"fmt"
	"testing"

	"github.com/ohler55/ojg/pretty"
	"github.com/ohler55/ojg/tt"
)

func TestWriteAlignSENArray(t *testing.T) {
	w := pretty.Writer{
		Width:    20,
		MaxDepth: 3,
		Align:    true,
	}
	out, err := w.Marshal([]interface{}{
		[]interface{}{1, 2, 3},
		[]interface{}{100, 200, 300},
	})
	tt.Nil(t, err)
	fmt.Printf("***\n%s\n", out)

	w.Width = 40
	out = w.Encode([]interface{}{
		[]interface{}{1, 2, 3, []interface{}{100, 200, 300}},
		[]interface{}{10, 20, 30, []interface{}{1, 20, 300}},
	})
	fmt.Printf("***\n%s\n", out)
}

func TestWriteAlignSENMap(t *testing.T) {
	w := pretty.Writer{
		Width:    50,
		MaxDepth: 3,
		Align:    true,
	}
	out := w.Encode([]interface{}{
		map[string]interface{}{"x": 1, "y": 2, "z": 3},
		map[string]interface{}{"x": 100, "y": 200, "z": 300},
		map[string]interface{}{"x": 10, "z": 30},
	})
	fmt.Printf("***\n%s\n", out)
}
