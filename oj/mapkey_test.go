package oj_test

import (
	"encoding/json"
	"testing"

	"github.com/ohler55/ojg/oj"
	"github.com/ohler55/ojg/tt"
)

type mapKey string

func (k mapKey) String() string { return "stringer:" + string(k) }

func TestMarshalNonStringMapKey(t *testing.T) {
	// Matches encoding/json, which renders integer map keys as their decimal form.
	src := map[int]string{1: "a", 2: "b"}
	std, err := json.Marshal(src)
	tt.Nil(t, err)
	b, err := oj.Marshal(src)
	tt.Nil(t, err)
	tt.Equal(t, `{"1":"a","2":"b"}`, sorted(string(b)))
	tt.Equal(t, string(std), sorted(string(b)))

	// Indented writer path.
	tt.Equal(t, "{\n  \"1\": \"a\"\n}", oj.JSON(map[int]string{1: "a"}, 2))

	// Sort must order by the rendered key, not by the empty reflect string.
	tt.Equal(t, `{"1":"a","2":"b"}`, oj.JSON(src, &oj.Options{Sort: true}))

	// A string-derived key keeps its own value even when it is a fmt.Stringer,
	// which is what encoding/json does.
	ks := map[mapKey]int{"k": 1}
	std, err = json.Marshal(ks)
	tt.Nil(t, err)
	tt.Equal(t, `{"k":1}`, oj.JSON(ks))
	tt.Equal(t, string(std), oj.JSON(ks))
}

// sorted normalizes Go's random map iteration order for the two-entry case.
func sorted(s string) string {
	if s == `{"2":"b","1":"a"}` {
		return `{"1":"a","2":"b"}`
	}
	return s
}
