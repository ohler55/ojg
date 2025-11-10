// Copyright (c) 2025, Peter Ohler, All rights reserved.

package discover_test

import (
	"os"
	"testing"

	"github.com/ohler55/ojg/discover"
	"github.com/ohler55/ojg/tt"
)

func TestFindNil(t *testing.T) {
	var found []byte
	discover.Find(nil, func(f []byte) (back, stop bool) {
		found = f
		return false, false
	})
	tt.Equal(t, "", string(found))
}

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

func TestFindMapNested(t *testing.T) {
	var found []byte
	discover.Find([]byte(`  { x:{y: {zz : {}} }} `), func(f []byte) (back, stop bool) {
		found = append(found, f...)
		return false, false
	})
	tt.Equal(t, `{ x:{y: {zz : {}} }}`, string(found))
}

func TestFindMapArray(t *testing.T) {
	var found []byte
	discover.Find([]byte(`  [{ x:[{y:2} z]}] `), func(f []byte) (back, stop bool) {
		found = append(found, f...)
		return false, false
	})
	tt.Equal(t, `[{ x:[{y:2} z]}]`, string(found))
}

func TestFindBackup(t *testing.T) {
	var found []byte
	discover.Find([]byte(`  [{ x:1}] `), func(f []byte) (back, stop bool) {
		if f[0] == '[' {
			return true, false
		}
		found = append(found, f...)
		return false, false
	})
	tt.Equal(t, `{ x:1}`, string(found))
}

func TestFindQuoteError(t *testing.T) {
	var found []byte
	discover.Find([]byte(`  [ a"bc" ["123"] `), func(f []byte) (back, stop bool) {
		found = f
		return false, false
	})
	tt.Equal(t, `["123"]`, string(found))

	found = found[:0]
	discover.Find([]byte(`  [ "ab"c ["123"] `), func(f []byte) (back, stop bool) {
		found = f
		return false, false
	})
	tt.Equal(t, `["123"]`, string(found))
}

// Tests the handling of reading multiple blocks of data.
func TestReadSplit(t *testing.T) {
	var found []byte
	r, err := os.Open("testdata/with-sen.txt")
	defer func() { _ = r.Close() }()
	tt.Nil(t, err)
	discover.Read(r, func(f []byte) (back, stop bool) {
		found = f
		return false, false
	})
	tt.Equal(t, `{
  abc: [
    1
    2
    3
    4
  ]

  xyz: [
    5
    6
    7
    8
  ]
}`, string(found))
}
