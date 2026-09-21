package alt_test

import (
	"testing"

	"github.com/ohler55/ojg/alt"
	"github.com/ohler55/ojg/oj"
)

type RKey string

func (k RKey) String() string { return "STRINGER-" + string(k) }

type rsample struct {
	Val int
}

// alt.Decompose keys any map with ojg.KeyString, so Recompose has to read back
// what Decompose writes. A struct target and a map target take different arms.
func TestRecomposeNonStringKeyedMap(t *testing.T) {
	t.Run("struct target, defined string key", func(t *testing.T) {
		v, err := alt.Recompose(map[RKey]any{"val": 3}, &rsample{})
		if err != nil {
			t.Fatalf("Recompose: %v", err)
		}
		if s, _ := v.(*rsample); s == nil || s.Val != 3 {
			t.Errorf("got %v, want &{3}", v)
		}
	})
	t.Run("struct target, int key", func(t *testing.T) {
		if _, err := alt.Recompose(map[int]any{7: 3}, &rsample{}); err != nil {
			t.Errorf("Recompose: %v", err)
		}
	})
	t.Run("map target, defined string key", func(t *testing.T) {
		v, err := alt.Recompose(map[RKey]any{"val": 3}, map[string]int{})
		if err != nil {
			t.Fatalf("Recompose: %v", err)
		}
		m, _ := v.(map[string]int)
		if m == nil || m["val"] != 3 {
			t.Errorf("got %v, want map[val:3]", v)
		}
	})
	t.Run("map target, int key", func(t *testing.T) {
		v, err := alt.Recompose(map[int]any{7: 3}, map[string]int{})
		if err != nil {
			t.Fatalf("Recompose: %v", err)
		}
		m, _ := v.(map[string]int)
		if m == nil || m["7"] != 3 {
			t.Errorf("got %v, want map[7:3]", v)
		}
	})
	t.Run("plain map[string]any is unchanged", func(t *testing.T) {
		v, err := alt.Recompose(map[string]any{"val": 3}, &rsample{})
		if err != nil {
			t.Fatalf("Recompose: %v", err)
		}
		if s, _ := v.(*rsample); s == nil || s.Val != 3 {
			t.Errorf("got %v, want &{3}", v)
		}
	})
	// The writers already key these the same way.
	t.Logf("oj.JSON(map[RKey]any) = %s", oj.JSON(map[RKey]any{"val": 3}))
	t.Logf("oj.JSON(map[int]any)  = %s", oj.JSON(map[int]any{7: 3}))
}
