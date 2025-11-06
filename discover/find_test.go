// Copyright (c) 2025, Peter Ohler, All rights reserved.

package discover_test

import (
	"testing"

	"github.com/ohler55/ojg/discover"
	"github.com/ohler55/ojg/tt"
)

func TestFindArrayEmpty(t *testing.T) {
	var found []byte
	discover.Find([]byte("  [  ] "), func(f []byte) (back, stop bool) {
		found = f
		return false, false
	})
	tt.Equal(t, "[  ]", string(found))
	found = found[:0]
	discover.Find([]byte("  [] "), func(f []byte) (back, stop bool) {
		found = f
		return false, true
	})
	tt.Equal(t, "[]", string(found))
}

func TestFindArrayNested(t *testing.T) {
	var found []byte
	discover.Find([]byte("  [ [ ] ] "), func(f []byte) (back, stop bool) {
		found = f
		return false, false
	})
	tt.Equal(t, "[ [ ] ]", string(found))
}

func TestFindArrayValues(t *testing.T) {
	var found []byte
	discover.Find([]byte("  [ abc 123 [ ] true] "), func(f []byte) (back, stop bool) {
		found = f
		return false, false
	})
	tt.Equal(t, "[ abc 123 [ ] true]", string(found))
}

func TestFindArrayBackup(t *testing.T) {
	var found []byte
	discover.Find([]byte("  [ a#b [x] ] "), func(f []byte) (back, stop bool) {
		found = f
		return false, false
	})
	tt.Equal(t, "[x]", string(found))
}

func TestFindArrayQuote1(t *testing.T) {
	var found []byte
	discover.Find([]byte(`  [ 'ab"c' [ ] '123'] `), func(f []byte) (back, stop bool) {
		found = f
		return false, false
	})
	tt.Equal(t, `[ 'ab"c' [ ] '123']`, string(found))
}

func TestFindArrayQuote2(t *testing.T) {
	var found []byte
	discover.Find([]byte(`  [ "ab'c" [ ] "123"] `), func(f []byte) (back, stop bool) {
		found = f
		return false, false
	})
	tt.Equal(t, `[ "ab'c" [ ] "123"]`, string(found))
}

func TestFindArrayEscape1(t *testing.T) {
	var found []byte
	discover.Find([]byte(`  [ '\t\n\b\f' '\"\'\\'] `), func(f []byte) (back, stop bool) {
		found = f
		return false, false
	})
	tt.Equal(t, `[ '\t\n\b\f' '\"\'\\']`, string(found))
}

func TestFindArrayEscape2(t *testing.T) {
	var found []byte
	discover.Find([]byte(`  [ "\t\n\b\f" "\"\'\\"] `), func(f []byte) (back, stop bool) {
		found = f
		return false, false
	})
	tt.Equal(t, `[ "\t\n\b\f" "\"\'\\"]`, string(found))
}

func TestFindArrayEscape1Unicode(t *testing.T) {
	var found []byte
	discover.Find([]byte(`  [ '\u0065' '\u12ab'] `), func(f []byte) (back, stop bool) {
		found = f
		return false, false
	})
	tt.Equal(t, `[ '\u0065' '\u12ab']`, string(found))
}

func TestFindArrayEscape2Unicode(t *testing.T) {
	var found []byte
	discover.Find([]byte(`  [ "\u0065" "\u12ab"] `), func(f []byte) (back, stop bool) {
		found = f
		return false, false
	})
	tt.Equal(t, `[ "\u0065" "\u12ab"]`, string(found))
}

func TestFindMapEmpty(t *testing.T) {
	var found []byte
	discover.Find([]byte(`  {} `), func(f []byte) (back, stop bool) {
		found = f
		return false, false
	})
	tt.Equal(t, `{}`, string(found))

	found = found[:0]
	discover.Find([]byte(`{ }`), func(f []byte) (back, stop bool) {
		found = f
		return false, false
	})
	tt.Equal(t, `{ }`, string(found))
}

func TestFindMapOneTight(t *testing.T) {
	var found []byte
	discover.Find([]byte(`  {abc:123} `), func(f []byte) (back, stop bool) {
		found = f
		return false, false
	})
	tt.Equal(t, `{abc:123}`, string(found))
}

func TestFindMapOneLoose(t *testing.T) {
	var found []byte
	discover.Find([]byte(`  { abc : 123 } `), func(f []byte) (back, stop bool) {
		found = f
		return false, false
	})
	tt.Equal(t, `{ abc : 123 }`, string(found))
}

func TestFindMapMultipleTight(t *testing.T) {
	var found []byte
	discover.Find([]byte(`  {abc:123 d : e, f:2} `), func(f []byte) (back, stop bool) {
		found = f
		return false, false
	})
	tt.Equal(t, `{abc:123 d : e, f:2}`, string(found))
}

func TestFindMapQuote2(t *testing.T) {
	var found []byte
	discover.Find([]byte(`  { "abc" : "123" } `), func(f []byte) (back, stop bool) {
		found = f
		return false, false
	})
	tt.Equal(t, `{ "abc" : "123" }`, string(found))

	found = found[:0]
	discover.Find([]byte(`  {"abc":"123"} `), func(f []byte) (back, stop bool) {
		found = f
		return false, false
	})
	tt.Equal(t, `{"abc":"123"}`, string(found))
}

func TestFindMapQuote1(t *testing.T) {
	var found []byte
	discover.Find([]byte(`  { 'abc' : '123' } `), func(f []byte) (back, stop bool) {
		found = f
		return false, false
	})
	tt.Equal(t, `{ 'abc' : '123' }`, string(found))

	found = found[:0]
	discover.Find([]byte(`  {'abc':'123'} `), func(f []byte) (back, stop bool) {
		found = f
		return false, false
	})
	tt.Equal(t, `{'abc':'123'}`, string(found))
}

// func TestFindMapDev(t *testing.T) {
// 	var found []byte
// 	discover.Find([]byte(`  { x:{} } `), func(f []byte) (back, stop bool) {
// 		found = append(found, f...)
// 		found = append(found, '\n')
// 		return false, false
// 	})
// 	tt.Equal(t, `{abc:123}`, string(found))
// }
