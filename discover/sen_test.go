// Copyright (c) 2025, Peter Ohler, All rights reserved.

package discover_test

import (
	"os"
	"testing"

	"github.com/ohler55/ojg/discover"
	"github.com/ohler55/ojg/pretty"
	"github.com/ohler55/ojg/tt"
)

func TestSENbytesNil(t *testing.T) {
	var found []byte
	discover.SENbytes(nil, func(f []byte) (back, stop bool) {
		found = f
		return false, false
	})
	tt.Equal(t, "", string(found))
}

func TestSENbytesArrayEmpty(t *testing.T) {
	var found []byte
	discover.SENbytes([]byte("  [  ] "), func(f []byte) (back, stop bool) {
		found = f
		return false, false
	})
	tt.Equal(t, "[  ]", string(found))
	found = found[:0]
	discover.SENbytes([]byte("  [] "), func(f []byte) (back, stop bool) {
		found = f
		return false, true
	})
	tt.Equal(t, "[]", string(found))
}

func TestSENbytesArrayNested(t *testing.T) {
	var found []byte
	discover.SENbytes([]byte("  [ [ ] ] "), func(f []byte) (back, stop bool) {
		found = f
		return false, false
	})
	tt.Equal(t, "[ [ ] ]", string(found))
}

func TestSENbytesArrayValues(t *testing.T) {
	var found []byte
	discover.SENbytes([]byte("  [ abc 123 [ ] true] "), func(f []byte) (back, stop bool) {
		found = f
		return false, false
	})
	tt.Equal(t, "[ abc 123 [ ] true]", string(found))
}

func TestSENbytesArrayBackup(t *testing.T) {
	var found []byte
	discover.SENbytes([]byte("  [ a#b [x] ] "), func(f []byte) (back, stop bool) {
		found = f
		return false, false
	})
	tt.Equal(t, "[x]", string(found))
}

func TestSENbytesArrayQuote1(t *testing.T) {
	var found []byte
	discover.SENbytes([]byte(`  [ 'ab"c' [ ] '123'] `), func(f []byte) (back, stop bool) {
		found = f
		return false, false
	})
	tt.Equal(t, `[ 'ab"c' [ ] '123']`, string(found))
}

func TestSENbytesArrayQuote2(t *testing.T) {
	var found []byte
	discover.SENbytes([]byte(`  [ "ab'c" [ ] "123"] `), func(f []byte) (back, stop bool) {
		found = f
		return false, false
	})
	tt.Equal(t, `[ "ab'c" [ ] "123"]`, string(found))
}

func TestSENbytesArrayEscape1(t *testing.T) {
	var found []byte
	discover.SENbytes([]byte(`  [ '\t\n\b\f' '\"\'\\'] `), func(f []byte) (back, stop bool) {
		found = f
		return false, false
	})
	tt.Equal(t, `[ '\t\n\b\f' '\"\'\\']`, string(found))
}

func TestSENbytesArrayEscape2(t *testing.T) {
	var found []byte
	discover.SENbytes([]byte(`  [ "\t\n\b\f" "\"\'\\"] `), func(f []byte) (back, stop bool) {
		found = f
		return false, false
	})
	tt.Equal(t, `[ "\t\n\b\f" "\"\'\\"]`, string(found))
}

func TestSENbytesArrayEscape1Unicode(t *testing.T) {
	var found []byte
	discover.SENbytes([]byte(`  [ '\u0065' '\u12ab'] `), func(f []byte) (back, stop bool) {
		found = f
		return false, false
	})
	tt.Equal(t, `[ '\u0065' '\u12ab']`, string(found))
}

func TestSENbytesArrayEscape2Unicode(t *testing.T) {
	var found []byte
	discover.SENbytes([]byte(`  [ "\u0065" "\u12ab"] `), func(f []byte) (back, stop bool) {
		found = f
		return false, false
	})
	tt.Equal(t, `[ "\u0065" "\u12ab"]`, string(found))
}

func TestSENbytesMapEmpty(t *testing.T) {
	var found []byte
	discover.SENbytes([]byte(`  {} `), func(f []byte) (back, stop bool) {
		found = f
		return false, false
	})
	tt.Equal(t, `{}`, string(found))

	found = found[:0]
	discover.SENbytes([]byte(`{ }`), func(f []byte) (back, stop bool) {
		found = f
		return false, false
	})
	tt.Equal(t, `{ }`, string(found))
}

func TestSENbytesMapOneTight(t *testing.T) {
	var found []byte
	discover.SENbytes([]byte(`  {abc:123} `), func(f []byte) (back, stop bool) {
		found = f
		return false, false
	})
	tt.Equal(t, `{abc:123}`, string(found))
}

func TestSENbytesMapOneLoose(t *testing.T) {
	var found []byte
	discover.SENbytes([]byte(`  { abc : 123 } `), func(f []byte) (back, stop bool) {
		found = f
		return false, false
	})
	tt.Equal(t, `{ abc : 123 }`, string(found))
}

func TestSENbytesMapMultipleTight(t *testing.T) {
	var found []byte
	discover.SENbytes([]byte(`  {abc:123 d : e, f:2} `), func(f []byte) (back, stop bool) {
		found = f
		return false, false
	})
	tt.Equal(t, `{abc:123 d : e, f:2}`, string(found))
}

func TestSENbytesMapQuote2(t *testing.T) {
	var found []byte
	discover.SENbytes([]byte(`  { "abc" : "123" } `), func(f []byte) (back, stop bool) {
		found = f
		return false, false
	})
	tt.Equal(t, `{ "abc" : "123" }`, string(found))

	found = found[:0]
	discover.SENbytes([]byte(`  {"abc":"123"} `), func(f []byte) (back, stop bool) {
		found = f
		return false, false
	})
	tt.Equal(t, `{"abc":"123"}`, string(found))
}

func TestSENbytesMapQuote1(t *testing.T) {
	var found []byte
	discover.SENbytes([]byte(`  { 'abc' : '123' } `), func(f []byte) (back, stop bool) {
		found = f
		return false, false
	})
	tt.Equal(t, `{ 'abc' : '123' }`, string(found))

	found = found[:0]
	discover.SENbytes([]byte(`  {'abc':'123'} `), func(f []byte) (back, stop bool) {
		found = f
		return false, false
	})
	tt.Equal(t, `{'abc':'123'}`, string(found))
}

func TestSENbytesMapNested(t *testing.T) {
	var found []byte
	discover.SENbytes([]byte(`  { x:{y: {zz : {}} }} `), func(f []byte) (back, stop bool) {
		found = append(found, f...)
		return false, false
	})
	tt.Equal(t, `{ x:{y: {zz : {}} }}`, string(found))
}

func TestSENbytesMapArray(t *testing.T) {
	var found []byte
	discover.SENbytes([]byte(`  [{ x:[{y:2} z]}] `), func(f []byte) (back, stop bool) {
		found = append(found, f...)
		return false, false
	})
	tt.Equal(t, `[{ x:[{y:2} z]}]`, string(found))
}

func TestSENbytesBackup(t *testing.T) {
	var found []byte
	discover.SENbytes([]byte(`  [{ x:1}] `), func(f []byte) (back, stop bool) {
		if f[0] == '[' {
			return true, false
		}
		found = append(found, f...)
		return false, false
	})
	tt.Equal(t, `{ x:1}`, string(found))
}

func TestSENbytesQuoteError(t *testing.T) {
	var found []byte
	discover.SENbytes([]byte(`  [ a"bc" ["123"] `), func(f []byte) (back, stop bool) {
		found = f
		return false, false
	})
	tt.Equal(t, `["123"]`, string(found))

	found = found[:0]
	discover.SENbytes([]byte(`  [ "ab"c ["123"] `), func(f []byte) (back, stop bool) {
		found = f
		return false, false
	})
	tt.Equal(t, `["123"]`, string(found))
}

func TestSENOkay(t *testing.T) {
	var found any
	discover.SEN([]byte(` [ [[] ] ] `), func(f any) (stop bool) {
		found = f
		return false
	})
	tt.Equal(t, `[[[]]]`, pretty.SEN(found))
}

func TestSENBack(t *testing.T) {
	var found any
	discover.SEN([]byte(` [ }[[] ] ] `), func(f any) (stop bool) {
		found = f
		return false
	})
	tt.Equal(t, `[[]]`, pretty.SEN(found))
}

func TestSENBadSEN(t *testing.T) {
	var found any
	discover.SEN([]byte(` [ 12x3 [4 ]] `), func(f any) (stop bool) {
		found = f
		return false
	})
	tt.Equal(t, `[4]`, pretty.SEN(found))
}

func TestSENShort(t *testing.T) {
	var found any
	discover.SEN([]byte(` [ 123 [4 ] `), func(f any) (stop bool) {
		found = f
		return false
	})
	tt.Equal(t, `[4]`, pretty.SEN(found))
}

func TestReadSENSplit(t *testing.T) {
	var found []byte
	r, err := os.Open("testdata/with-sen.txt")
	tt.Nil(t, err)
	discover.ReadSENbytes(r, func(f []byte) (back, stop bool) {
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
