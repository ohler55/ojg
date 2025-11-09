// Copyright (c) 2025, Peter Ohler, All rights reserved.

package discover

import (
	"errors"
	"io"
)

const readBufSize = 4096

func readMore(r io.Reader, b []byte, start, i int) ([]byte, int, int, bool) {
	bc := cap(b)
	used := i - start
	orig := b
	if bc < used+readBufSize {
		b = make([]byte, used+readBufSize)
		if start == 0 {
			copy(b, orig)
		}
	}
	if 0 < start {
		copy(b, orig[start:i])
		i = used
		start = 0
	}

	cnt, err := r.Read(b[i:])
	if err != nil {
		if !errors.Is(err, io.EOF) {
			panic(err)
		}
	}
	b = b[:used+cnt]

	return b, start, i, cnt == 0
}
