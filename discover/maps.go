// Copyright (c) 2025, Peter Ohler, All rights reserved.

package discover

const (
	skip    = '.'
	errChar = 'e'
	popMode = 'p'

	quote2   = 'Q'
	quote1   = 'q'
	quoteEnd = 'z'
	escape   = 'b'
	escU     = 'U'
	uOk      = 'u'

	openArray  = '['
	closeArray = ']'

	openObject    = '{'
	closeObject   = '}'
	keyChar       = 'k'
	keyDoneChar   = 'K'
	colonChar     = ':'
	valueChar     = 'v'
	valueDoneChar = 'V'

	//   0123456789abcdef0123456789abcdef
	scanMap = "" +
		"................................" + // 0x00
		"................................" + // 0x20
		"...........................[...." + // 0x40
		"...........................{...." + // 0x60
		"................................" + // 0x80
		"................................" + // 0xa0
		"................................" + // 0xc0
		"................................S" //  0xe0

	//   0123456789abcdef0123456789abcdef
	senArrayMap = "" +
		"eeeeeeeee..ee.eeeeeeeeeeeeeeeeee" + // 0x00
		".eQe...qee................ee...." + // 0x20
		"...........................[e].." + // 0x40
		"e..........................{.e.e" + // 0x60
		"................................" + // 0x80
		"................................" + // 0xa0
		"................................" + // 0xc0
		"................................A" //  0xe0

	//   0123456789abcdef0123456789abcdef
	senObjectMap = "" +
		"eeeeeeeee..ee.eeeeeeeeeeeeeeeeee" + // 0x00
		".eQekkkqeekk.kkkeeeeeeeeeeeekkkk" + // 0x20
		"kkkkkkkkkkkkkkkkkkkkkkkkkkkeeekk" + // 0x40
		"ekkkkkkkkkkkkkkkkkkkkkkkkkkek}ke" + // 0x60
		"kkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkk" + // 0x80
		"kkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkk" + // 0xa0
		"kkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkk" + // 0xc0
		"kkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkO" //  0xe0

	//   0123456789abcdef0123456789abcdef
	senKeyMap = "" +
		"eeeeeeeeeKKeeKeeeeeeeeeeeeeeeeee" + // 0x00
		"Keee...eee..K.............Ke...." + // 0x20
		"...........................eee.." + // 0x40
		"e..........................e.e.e" + // 0x60
		"................................" + // 0x80
		"................................" + // 0xa0
		"................................" + // 0xc0
		"................................K" //  0xe0

	//   0123456789abcdef0123456789abcdef
	senColonMap = "" +
		"eeeeeeeee..ee.eeeeeeeeeeeeeeeeee" + // 0x00
		".eeeeeeeeeeeeeeeeeeeeeeeee:eeeee" + // 0x20
		"eeeeeeeeeeeeeeeeeeeeeeeeeeeeeeee" + // 0x40
		"eeeeeeeeeeeeeeeeeeeeeeeeeeeeeeee" + // 0x60
		"eeeeeeeeeeeeeeeeeeeeeeeeeeeeeeee" + // 0x80
		"eeeeeeeeeeeeeeeeeeeeeeeeeeeeeeee" + // 0xa0
		"eeeeeeeeeeeeeeeeeeeeeeeeeeeeeeee" + // 0xc0
		"eeeeeeeeeeeeeeeeeeeeeeeeeeeeeeee:" //  0xe0

	//   0123456789abcdef0123456789abcdef
	senPreValueMap = "" +
		"eeeeeeeee..ee.eeeeeeeeeeeeeeeeee" + // 0x00
		".eQevvvqeevvevvvvvvvvvvvvveevvvv" + // 0x20
		"vvvvvvvvvvvvvvvvvvvvvvvvvvv[eevv" + // 0x40
		"evvvvvvvvvvvvvvvvvvvvvvvvvv{veve" + // 0x60
		"vvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvv" + // 0x80
		"vvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvv" + // 0xa0
		"vvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvv" + // 0xc0
		"vvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvv" //  0xe0

	//   0123456789abcdef0123456789abcdef
	senValueMap = "" +
		"eeeeeeeeeVVeeVeeeeeeeeeeeeeeeeee" + // 0x00
		"Veee...eee..V.............ee...." + // 0x20
		"...........................eee.." + // 0x40
		"e..........................e.}.e" + // 0x60
		"................................" + // 0x80
		"................................" + // 0xa0
		"................................" + // 0xc0
		"................................V" //  0xe0

	//   0123456789abcdef0123456789abcdef
	quoteOkMap = "" +
		".........oo..o.................." + // 0x00
		"o...........o.............o....." + // 0x20
		"...........................o.o.." + // 0x40
		"...........................o.o.." + // 0x60
		"................................" + // 0x80
		"................................" + // 0xa0
		"................................" + // 0xc0
		"................................" //   0xe0

	//   0123456789abcdef0123456789abcdef
	quote1Map = "" +
		"eeeeeeeee..ee.eeeeeeeeeeeeeeeeee" + // 0x00
		".......z........................" + // 0x20
		"............................b..." + // 0x40
		"...............................e" + // 0x60
		"................................" + // 0x80
		"................................" + // 0xa0
		"................................" + // 0xc0
		"................................q" //  0xe0

	//   0123456789abcdef0123456789abcdef
	quote2Map = "" +
		"eeeeeeeee..ee.eeeeeeeeeeeeeeeeee" + // 0x00
		"..z............................." + // 0x20
		"............................b..." + // 0x40
		"...............................e" + // 0x60
		"................................" + // 0x80
		"................................" + // 0xa0
		"................................" + // 0xc0
		"................................Q" //  0xe0

	//   0123456789abcdef0123456789abcdef
	escMap = "" +
		"eeeeeeeeeeeeeeeeeeeeeeeeeeeeeeee" + // 0x00
		"eepeeeepeeeeeeepeeeeeeeeeeeeeeee" + // 0x20
		"eeeeeeeeeeeeeeeeeeeeeeeeeeeepeee" + // 0x40
		"eepeeepeeeeeeepeeepepUeeeeeeeeee" + // 0x60
		"eeeeeeeeeeeeeeeeeeeeeeeeeeeeeeee" + // 0x80
		"eeeeeeeeeeeeeeeeeeeeeeeeeeeeeeee" + // 0xa0
		"eeeeeeeeeeeeeeeeeeeeeeeeeeeeeeee" + // 0xc0
		"eeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeb" //  0xe0

	//   0123456789abcdef0123456789abcdef
	uMap = "" +
		"eeeeeeeeeeeeeeeeeeeeeeeeeeeeeeee" + // 0x00
		"eeeeeeeeeeeeeeeeuuuuuuuuuueeeeee" + // 0x20
		"euuuuuueeeeeeeeeeeeeeeeeeeeeeeee" + // 0x40
		"euuuuuueeeeeeeeeeeeeeeeeeeeeeeee" + // 0x60
		"eeeeeeeeeeeeeeeeeeeeeeeeeeeeeeee" + // 0x80
		"eeeeeeeeeeeeeeeeeeeeeeeeeeeeeeee" + // 0xa0
		"eeeeeeeeeeeeeeeeeeeeeeeeeeeeeeee" + // 0xc0
		"eeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeu" //   0xe0

)
