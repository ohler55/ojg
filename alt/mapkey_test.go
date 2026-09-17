package alt_test

import (
	"encoding/json"
	"testing"

	"github.com/ohler55/ojg/alt"
	"github.com/ohler55/ojg/oj"
	"github.com/ohler55/ojg/pretty"
	"github.com/ohler55/ojg/tt"
)

type mapKey string

func (k mapKey) String() string { return "stringer:" + string(k) }

func TestMapKeyStringerNotCalled(t *testing.T) {
	src := map[mapKey]int{"k": 1}
	std, err := json.Marshal(src)
	tt.Nil(t, err)
	tt.Equal(t, `{"k":1}`, string(std))

	// Decompose must key on the string value, not on String().
	tt.Equal(t, map[string]any{"k": 1}, alt.Decompose(src))

	// Generify takes the same key through a second copy of the same code.
	tt.Equal(t, `{"k":1}`, alt.Generify(src).String())

	// pretty writes through alt.Decompose, and oj writes through ojg.KeyString;
	// the two must agree with each other and with encoding/json.
	tt.Equal(t, `{"k": 1}`, pretty.JSON(src))
	tt.Equal(t, oj.JSON(src), string(std))

	// A plain named string key was already correct and must stay correct.
	type plainKey string
	tt.Equal(t, map[string]any{"k": 1}, alt.Decompose(map[plainKey]int{"k": 1}))

	// A genuinely non-string key still renders through fmt.
	tt.Equal(t, map[string]any{"1": "a"}, alt.Decompose(map[int]string{1: "a"}))
	tt.Equal(t, `{"1":"a"}`, alt.Generify(map[int]string{1: "a"}).String())
}
