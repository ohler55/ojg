// Copyright (c) 2026, Peter Ohler, All rights reserved.

package bom_test

import (
	"io"
	"strings"
	"testing"
	"testing/iotest"

	"github.com/ohler55/ojg/internal/bom"
	"github.com/ohler55/ojg/tt"
)

type zeroReader struct{}

func (zeroReader) Read([]byte) (int, error) { return 0, nil }

func TestReadRest(t *testing.T) {
	buf := make([]byte, 16)
	buf[0] = 0xEF

	n, err := bom.ReadRest(iotest.OneByteReader(strings.NewReader("\xbb\xbf{}")), buf, 1)
	tt.Nil(t, err)
	tt.Equal(t, 3, n)
	tt.Equal(t, []byte{0xEF, 0xBB, 0xBF}, buf[:n])

	buf[0] = 0xEF
	n, err = bom.ReadRest(iotest.OneByteReader(strings.NewReader("\xbb")), buf, 1)
	tt.Equal(t, io.EOF, err)
	tt.Equal(t, 2, n)

	buf[0] = 0xEF
	n, err = bom.ReadRest(strings.NewReader("ab"), buf, 1)
	tt.Nil(t, err)
	tt.Equal(t, 3, n)

	buf[0] = 0xEF
	n, err = bom.ReadRest(zeroReader{}, buf, 1)
	tt.Nil(t, err)
	tt.Equal(t, 1, n)
}
