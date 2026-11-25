// Copyright (c) 2026, Peter Ohler, All rights reserved.

package ojg

import "io"

// ReadBOMStart reads the start of r into buf, continuing short reads only
// while they can still form a UTF-8 BOM.
func ReadBOMStart(r io.Reader, buf []byte) (n int, err error) {
	for n < 3 {
		var cnt int
		cnt, err = r.Read(buf[n:])
		n += cnt
		if err != nil || cnt == 0 || buf[0] != 0xEF || (1 < n && buf[1] != 0xBB) || (2 < n && buf[2] != 0xBF) {
			return
		}
	}
	return
}
