// Copyright (c) 2026, Peter Ohler, All rights reserved.

package oj

import "io"

func readBOMStart(r io.Reader, buf []byte) (n int, err error) {
	n, err = r.Read(buf)
	if n == 0 || err != nil || buf[0] != 0xEF {
		return
	}
	if 1 < n && buf[1] != 0xBB {
		return
	}
	for n < 3 {
		var cnt int
		cnt, err = r.Read(buf[n:])
		n += cnt
		if err != nil {
			return
		}
		if 1 < n && buf[1] != 0xBB {
			return
		}
		if 2 < n && buf[2] != 0xBF {
			return
		}
	}
	return
}
