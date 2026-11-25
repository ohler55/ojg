// Copyright (c) 2026, Peter Ohler, All rights reserved.

package bom

import "io"

// ReadRest finishes reading up to 3 leading bytes into buf after an initial
// short read that began with the UTF-8 BOM lead byte 0xEF.
func ReadRest(r io.Reader, buf []byte, n int) (int, error) {
	for n < 3 {
		cnt, err := r.Read(buf[n:])
		n += cnt
		if err != nil || cnt == 0 || (1 < n && buf[1] != 0xBB) || (2 < n && buf[2] != 0xBF) {
			return n, err
		}
	}
	return n, nil
}
