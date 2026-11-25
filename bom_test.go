// Copyright (c) 2026, Peter Ohler, All rights reserved.

package ojg_test

import (
	"io"
	"strings"
	"testing"
	"testing/iotest"

	"github.com/ohler55/ojg"
	"github.com/ohler55/ojg/tt"
)

type zeroReader struct{}

func (zeroReader) Read([]byte) (int, error) { return 0, nil }

func TestReadBOMStart(t *testing.T) {
	buf := make([]byte, 16)

	n, err := ojg.ReadBOMStart(iotest.OneByteReader(strings.NewReader("\xef\xbb\xbf{}")), buf)
	tt.Nil(t, err)
	tt.Equal(t, 3, n)
	tt.Equal(t, []byte{0xEF, 0xBB, 0xBF}, buf[:n])

	n, err = ojg.ReadBOMStart(iotest.OneByteReader(strings.NewReader("\xef\xbb")), buf)
	tt.Equal(t, io.EOF, err)
	tt.Equal(t, 2, n)
	tt.Equal(t, []byte{0xEF, 0xBB}, buf[:n])

	n, err = ojg.ReadBOMStart(strings.NewReader("abc"), buf)
	tt.Nil(t, err)
	tt.Equal(t, 3, n)
	tt.Equal(t, []byte("abc"), buf[:n])

	n, err = ojg.ReadBOMStart(zeroReader{}, buf)
	tt.Nil(t, err)
	tt.Equal(t, 0, n)
}
