// Copyright (c) 2025, Peter Ohler, All rights reserved.

package discover

import "fmt"

// Find occurrence of SEN documents that are either maps or arrays. This is a
// best effort search to find potential JSON or SEN documents. It is possible
// that document will not parse without errors. The callback function should
// return a true back return value to back up to the next open character after
// the current start. If back is false scanning continues after the end of the
// found section. If stop is true then no further scanning is attempted and
// the function returns.
func Find(buf []byte, cb func(found []byte) (back, stop bool)) {
	var (
		b     byte
		start int
		modes []string
		i     int
		ucnt  int
	)
	mode := scanMap
	reset := func() {
		i = start
		modes = modes[:0]
		mode = scanMap
	}
	for i = 0; i < len(buf); i++ {
		b = buf[i]
		// fmt.Printf("%d: '%c' 0x%02x - 0x%02x in %c\n", i, b, b, mode[b], mode[256])
		switch mode[b] {
		case skip:
			// no change
		case openArray:
			if len(modes) == 0 {
				start = i
			}
			modes = append(modes, mode)
			mode = senArrayMap
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
		case openObject:
			if len(modes) == 0 {
				start = i
			}
			// TBD append sendValueMap
			modes = append(modes, mode)
			mode = senPreKeyMap
			// TBD keyMode and map - keyMap->colonMap->valueMapp
			//  looks end of token with : or space terminator
			//  after key, need :
			//  read value next with space or }
			//    value can also be [] or {}
		case keyChar:
			mode = senKeyMap
		case keyDoneChar:
			fmt.Printf("*** key done at %q\n", buf[:i])
			// TBD colonMap
			return
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

			///////// TBD
			// case openObject:
			// 	if len(modes) == 0 {
			// 		start = i
			// 	}
			// 	modes = append(modes, mode)
			// 	mode = keyMap
			// case closeObject:
			// 	mode = modes[len(modes)-1]
			// 	modes = modes[:len(modes)-1]
			// 	if len(modes) == 0 {
			// 		back, stop := cb(buf[start : i+1])
			// 		if stop {
			// 			return
			// 		}
			// 		if back {
			// 			i = start
			// 		}
			// 	}

			// do nothing
			// TBD if starts is not empty
			//  must match starts
			//  shorten starts
			//  if starts is empty then cb
			// TBD if no match then back to start+1 or maybe keep track of next

			// TBD note in quotes, single or double
			// do we care about : and , - probably for better results
			// white space is important as a separator

			// TBD maybe , and required quotes for JSON and not SEN

		}
	}
}
