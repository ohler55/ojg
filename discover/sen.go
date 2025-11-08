// Copyright (c) 2025, Peter Ohler, All rights reserved.

package discover

import (
	"errors"
	"fmt"
	"io"

	"github.com/ohler55/ojg/sen"
)

// TBD uncomment for debugging
// func modesString(modes []string) string {
// 	var b []byte
// 	b = append(b, '[')
// 	for i, m := range modes {
// 		if 0 < i {
// 			b = append(b, ' ')
// 		}
// 		b = append(b, m[256])
// 	}
// 	return string(append(b, ']'))
// }

// SENbytes finds potential occurrence of SEN documents that are either maps
// or arrays. This is a best effort search to find potential SEN documents. It
// is possible that document will not parse without errors. The callback
// function should return a true back return value to back up to the next open
// character after the current start. If back is false scanning continues
// after the end of the found section. If stop is true then no further
// scanning is attempted and the function returns.
func SENbytes(buf []byte, cb func(found []byte) (back, stop bool)) {
	senBytes(buf, cb, nil)
}

func senBytes(
	buf []byte,
	cb func(found []byte) (back, stop bool),
	more func(buf []byte, start, i int) ([]byte, int, int, bool)) {

	var (
		b     byte
		start int
		modes []string
		ucnt  int
		i     int
	)
	mode := scanMap
	reset := func() {
		i = start
		modes = modes[:0]
		mode = scanMap
	}
retry:
	for i = start; i < len(buf); i++ {
		b = buf[i]
		// fmt.Printf("%d: '%c' 0x%02x - %c in %c in %s\n", i, b, b, mode[b], mode[256], modesString(modes))
		switch mode[b] {
		case skip:
			// no change
		case openArray:
			if len(modes) == 0 {
				start = i
			}
			if mode == senPreValueMap || mode == senValueMap {
				modes = append(modes, senObjectMap)
			} else {
				modes = append(modes, mode)
			}
			mode = senArrayMap
		case openObject:
			if len(modes) == 0 {
				start = i
			}
			if mode == senPreValueMap || mode == senValueMap {
				modes = append(modes, senObjectMap)
			} else {
				modes = append(modes, mode)
			}
			mode = senObjectMap
		case closeArray, closeObject:
			mode = modes[len(modes)-1]
			modes = modes[:len(modes)-1]
			if len(modes) == 0 {
				back, stop := cb(buf[start : i+1])
				if stop {
					return
				}
				if back {
					i = start
				}
			}
		case keyChar:
			mode = senKeyMap
		case keyDoneChar:
			if b == ':' {
				mode = senPreValueMap
			} else {
				mode = senColonMap
			}
		case colonChar:
			mode = senPreValueMap

		case quote1, quote2:
			if quoteOkMap[buf[i-1]] != 'o' {
				reset()
				break
			}
			modes = append(modes, mode)
			if b == '"' {
				mode = quote2Map
			} else {
				mode = quote1Map
			}
		case quoteEnd:
			mode = modes[len(modes)-1]
			modes = modes[:len(modes)-1]
			switch mode {
			case senObjectMap:
				mode = senColonMap
			case senPreValueMap:
				mode = senObjectMap
			}
		case valueChar:
			mode = senValueMap
		case valueDoneChar:
			mode = senObjectMap
		case popMode:
			mode = modes[len(modes)-1]
			modes = modes[:len(modes)-1]
		case escape:
			modes = append(modes, mode)
			mode = escMap
		case escU:
			ucnt = 0
			modes = append(modes, mode)
			mode = uMap
		case uOk:
			if ucnt == 3 {
				mode = modes[len(modes)-2]
				modes = modes[:len(modes)-2]
			} else {
				ucnt++
			}
		case errChar:
			reset()
		}
	}
	// fmt.Printf("*** out of loop\n")

	if more != nil {
		var eof bool
		buf, start, i, eof = more(buf, start, i)
		if !eof {
			goto retry
		}
	}
	if 0 < len(modes) {
		start++
		if start < len(buf) {
			reset()
			goto retry
		}
	}
}

// SEN finds occurrence of SEN documents that are either maps or arrays. The
// callback function should return true to stop discovering.
func SEN(buf []byte, cb func(value any) (stop bool)) {
	SENbytes(buf, func(found []byte) (bool, bool) {
		if value, err := sen.Parse(found); err == nil {
			return false, cb(value)
		}
		return true, false
	})
}

// ReadSENbytes finds potential occurrence of SEN documents that are either
// maps or arrays in a stream. This is a best effort search to find potential
// SEN documents. It is possible that document will not parse without
// errors. The callback function should return a true back return value to
// back up to the next open character after the current start. If back is
// false scanning continues after the end of the found section. If stop is
// true then no further scanning is attempted and the function returns.
func ReadSENbytes(r io.Reader, cb func(b []byte) (back, stop bool)) {
	senBytes(nil, cb, func(buf []byte, start, i int) ([]byte, int, int, bool) {
		return readMore(r, buf, start, i)
	})
}

// ReadSEN finds occurrence of SEN documents that are either maps or arrays in
// a stream. The callback function should return true to stop discovering.
func ReadSEN(r io.Reader, f func(v any) bool) {
	// TBD
}

const readBufSize = 4096

func readMore(r io.Reader, b []byte, start, i int) ([]byte, int, int, bool) {
	fmt.Printf("*** readMore %d %d\n", start, i)
	bc := cap(b)
	used := i - start
	orig := b
	if used+readBufSize < bc {
		b = make([]byte, used+readBufSize)
		// if 0 < used {
		// 	// shift existing
		// 	copy(b, orig[start:i])
		// 	i = used
		// 	start = 0
		// }
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
