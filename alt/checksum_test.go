// Copyright (c) 2026, Peter Ohler, All rights reserved.

package alt_test

import (
	"testing"
	"time"

	"github.com/ohler55/ojg/alt"
	"github.com/ohler55/ojg/tt"
)

type checker struct {
	a int
	b int
}

func (c *checker) Simplify() any {
	return map[string]any{"a": c.a, "b": c.b}
}

func TestChecksumNil(t *testing.T) {
	tt.Equal(t, uint64(2282658103124508505), alt.Checksum(nil))
}

func TestChecksumBool(t *testing.T) {
	tt.Equal(t, uint64(3761427544677127395), alt.Checksum(true))
	tt.Equal(t, uint64(17263442184195172173), alt.Checksum(false))
}

func TestChecksumInt(t *testing.T) {
	tt.Equal(t, uint64(13144445341178972864), alt.Checksum(0))
	tt.Equal(t, uint64(9948673712954436551), alt.Checksum(7))
	tt.Equal(t, uint64(9948673712954436551), alt.Checksum(int8(7)))
	tt.Equal(t, uint64(9948673712954436551), alt.Checksum(int16(7)))
	tt.Equal(t, uint64(9948673712954436551), alt.Checksum(int32(7)))
	tt.Equal(t, uint64(9948673712954436551), alt.Checksum(int64(7)))
	tt.Equal(t, uint64(9948673712954436551), alt.Checksum(uint(7)))
	tt.Equal(t, uint64(9948673712954436551), alt.Checksum(uint8(7)))
	tt.Equal(t, uint64(9948673712954436551), alt.Checksum(uint16(7)))
	tt.Equal(t, uint64(9948673712954436551), alt.Checksum(uint32(7)))
	tt.Equal(t, uint64(9948673712954436551), alt.Checksum(uint64(7)))
}

func TestChecksumFloat(t *testing.T) {
	tt.Equal(t, uint64(3849879868877111371), alt.Checksum(2.5))
	tt.Equal(t, uint64(3849879868877111371), alt.Checksum(float32(2.5)))
}

func TestChecksumString(t *testing.T) {
	tt.Equal(t, uint64(0), alt.Checksum(""))
	tt.Equal(t, uint64(11642063096747747405), alt.Checksum("quux"))
	tt.Equal(t, uint64(11642063096747747405), alt.Checksum([]byte("quux")))
}

func TestChecksumTime(t *testing.T) {
	tt.Equal(t, uint64(11794479198235629926),
		alt.Checksum(time.Date(2026, time.January, 21, 40, 29, 04, 123456789, time.UTC)))
}

func TestChecksumArray(t *testing.T) {
	tt.Equal(t, uint64(1400374585248166553), alt.Checksum([]any{1, []any{true, false}, 3}))
	tt.Equal(t, uint64(12716749365223810667), alt.Checksum([]any{1, []any{false, true}, 3}))
}

func TestChecksumMap(t *testing.T) {
	tt.Equal(t, uint64(18395434359605210485),
		alt.Checksum(map[string]any{"a": 1, "b": map[string]any{"c": true, "d": false}, "e": 3}))
}

func TestChecksumSimplifier(t *testing.T) {
	tt.Equal(t, uint64(15053623628309097403), alt.Checksum(&checker{a: 1, b: 2}))
}

func TestChecksumDecompose(t *testing.T) {
	type quux struct {
		a int
		b int
	}
	tt.Equal(t, uint64(10091151871759098234), alt.Checksum(&quux{a: 1, b: 2}))
}
