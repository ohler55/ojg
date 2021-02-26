// Copyright (c) 2020, Peter Ohler, All rights reserved.

package oj

const (
	skipChar    = 'a'
	skipNewline = 'b'
	valNull     = 'c'
	valTrue     = 'd'
	valFalse    = 'e'
	valNeg      = 'f'
	val0        = 'g'
	valDigit    = 'h'
	valQuote    = 'i'
	openArray   = 'k'
	openObject  = 'l'
	closeArray  = 'm'
	closeObject = 'n'
	afterComma  = 'o'
	keyQuote    = 'p'
	colonColon  = 'q'
	numSpc      = 'r'
	numNewline  = 's'
	numDot      = 't'
	numComma    = 'u'
	numFrac     = 'v'
	fracE       = 'w'
	expSign     = 'x'
	expDigit    = 'y'
	strQuote    = 'z'
	negDigit    = '-'
	strSlash    = 'A'
	escOk       = 'B'
	uOk         = 'E'
	tokenOk     = 'F'
	numDigit    = 'N'
	numZero     = 'O'
	strOk       = 'R'
	escU        = 'U'
	charErr     = '.'

	//   0123456789abcdef0123456789abcdef
	valueMap = "" +
		".........ab..a.................." + // 0x00
		"a.i..........f..ghhhhhhhhh......" + // 0x20
		"...........................k.m.." + // 0x40
		"......e.......c.....d......l.n.." + // 0x60
		"................................" + // 0x80
		"................................" + // 0xa0
		"................................" + // 0xc0
		"................................v" //  0xe0
	//   0123456789abcdef0123456789abcdef
	nullMap = "" +
		"................................" + // 0x00
		"............o..................." + // 0x20
		"................................" + // 0x40
		"............F........F.........." + // 0x60
		"................................" + // 0x80
		"................................" + // 0xa0
		"................................" + // 0xc0
		"................................" //   0xe0
	//   0123456789abcdef0123456789abcdef
	trueMap = "" +
		"................................" + // 0x00
		"............o..................." + // 0x20
		"................................" + // 0x40
		".....F............F..F.........." + // 0x60
		"................................" + // 0x80
		"................................" + // 0xa0
		"................................" + // 0xc0
		"................................" //   0xe0
	//   0123456789abcdef0123456789abcdef
	falseMap = "" +
		"................................" + // 0x00
		"............o..................." + // 0x20
		"................................" + // 0x40
		".F...F......F......F............" + // 0x60
		"................................" + // 0x80
		"................................" + // 0xa0
		"................................" + // 0xc0
		"................................" //   0xe0
	//   0123456789abcdef0123456789abcdef
	commaMap = "" +
		".........ab..a.................." + // 0x00
		"a.i..........f..ghhhhhhhhh......" + // 0x20
		"...........................k...." + // 0x40
		"......e.......c.....d......l...." + // 0x60
		"................................" + // 0x80
		"................................" + // 0xa0
		"................................" + // 0xc0
		"................................" //   0xe0
	//   0123456789abcdef0123456789abcdef
	afterMap = "" +
		".........ab..a.................." + // 0x00
		"a...........o..................." + // 0x20
		".............................m.." + // 0x40
		".............................n.." + // 0x60
		"................................" + // 0x80
		"................................" + // 0xa0
		"................................" + // 0xc0
		"................................a" //  0xe0
	//   0123456789abcdef0123456789abcdef
	key1Map = "" +
		".........ab..a.................." + // 0x00
		"a.p............................." + // 0x20
		"................................" + // 0x40
		".............................n.." + // 0x60
		"................................" + // 0x80
		"................................" + // 0xa0
		"................................" + // 0xc0
		"................................" //   0xe0
	//   0123456789abcdef0123456789abcdef
	keyMap = "" +
		".........ab..a.................." + // 0x00
		"a.p............................." + // 0x20
		"................................" + // 0x40
		"................................" + // 0x60
		"................................" + // 0x80
		"................................" + // 0xa0
		"................................" + // 0xc0
		"................................" //   0xe0
	//   0123456789abcdef0123456789abcdef
	colonMap = "" +
		".........ab..a.................." + // 0x00
		"a.........................q....." + // 0x20
		"................................" + // 0x40
		"................................" + // 0x60
		"................................" + // 0x80
		"................................" + // 0xa0
		"................................" + // 0xc0
		"................................" //   0xe0
	//   0123456789abcdef0123456789abcdef
	negMap = "" +
		"................................" + // 0x00
		"................O---------......" + // 0x20
		"................................" + // 0x40
		"................................" + // 0x60
		"................................" + // 0x80
		"................................" + // 0xa0
		"................................" + // 0xc0
		"................................" //   0xe0
	//   0123456789abcdef0123456789abcdef
	zeroMap = "" +
		".........rs..r.................." + // 0x00
		"r...........u.t................." + // 0x20
		".............................m.." + // 0x40
		".............................n.." + // 0x60
		"................................" + // 0x80
		"................................" + // 0xa0
		"................................" + // 0xc0
		"................................n" //  0xe0
	//   0123456789abcdef0123456789abcdef
	digitMap = "" +
		".........rs..r.................." + // 0x00
		"r...........u.t.NNNNNNNNNN......" + // 0x20
		".............................m.." + // 0x40
		".............................n.." + // 0x60
		"................................" + // 0x80
		"................................" + // 0xa0
		"................................" + // 0xc0
		"................................n" //  0xe0
	//   0123456789abcdef0123456789abcdef
	dotMap = "" +
		"................................" + // 0x00
		"................vvvvvvvvvv......" + // 0x20
		"................................" + // 0x40
		"................................" + // 0x60
		"................................" + // 0x80
		"................................" + // 0xa0
		"................................" + // 0xc0
		"................................" //   0xe0
	//   0123456789abcdef0123456789abcdef
	fracMap = "" +
		".........rs..r.................." + // 0x00
		"r...........u...vvvvvvvvvv......" + // 0x20
		".....w.......................m.." + // 0x40
		".....w.......................n.." + // 0x60
		"................................" + // 0x80
		"................................" + // 0xa0
		"................................" + // 0xc0
		"................................n" //  0xe0
	//   0123456789abcdef0123456789abcdef
	expSignMap = "" +
		"................................" + // 0x00
		"...........x.x..yyyyyyyyyy......" + // 0x20
		"................................" + // 0x40
		"................................" + // 0x60
		"................................" + // 0x80
		"................................" + // 0xa0
		"................................" + // 0xc0
		"................................" //   0xe0
	//   0123456789abcdef0123456789abcdef
	expZeroMap = "" +
		"................................" + // 0x00
		"................yyyyyyyyyy......" + // 0x20
		"................................" + // 0x40
		"................................" + // 0x60
		"................................" + // 0x80
		"................................" + // 0xa0
		"................................" + // 0xc0
		"................................" //   0xe0
	//   0123456789abcdef0123456789abcdef
	expMap = "" +
		".........rs..r.................." + // 0x00
		"r...........u...yyyyyyyyyy......" + // 0x20
		".............................m.." + // 0x40
		".............................n.." + // 0x60
		"................................" + // 0x80
		"................................" + // 0xa0
		"................................" + // 0xc0
		"................................n" //  0xe0
	//   0123456789abcdef0123456789abcdef
	stringMap = "" +
		"................................" + // 0x00
		"RRzRRRRRRRRRRRRRRRRRRRRRRRRRRRRR" + // 0x20
		"RRRRRRRRRRRRRRRRRRRRRRRRRRRRARRR" + // 0x40
		"RRRRRRRRRRRRRRRRRRRRRRRRRRRRRRRR" + // 0x60
		"RRRRRRRRRRRRRRRRRRRRRRRRRRRRRRRR" + // 0x80
		"RRRRRRRRRRRRRRRRRRRRRRRRRRRRRRRR" + // 0xR0
		"RRRRRRRRRRRRRRRRRRRRRRRRRRRRRRRR" + // 0xc0
		"RRRRRRRRRRRRRRRRRRRRRRRRRRRRRRRR" //   0xe0
	//   0123456789abcdef0123456789abcdef
	escMap = "" +
		"................................" + // 0x00
		"..B............B................" + // 0x20
		"............................B..." + // 0x40
		"..B...B.......B...B.BU.........." + // 0x60
		"................................" + // 0x80
		"................................" + // 0xa0
		"................................" + // 0xc0
		"................................" //   0xe0
	//   0123456789abcdef0123456789abcdef
	escByteMap = "" +
		"................................" + // 0x00
		"..\"............/................" + // 0x20
		"............................\\..." + // 0x40
		"..\b...\f.......\n...\r.\t.........." + // 0x60
		"................................" + // 0x80
		"................................" + // 0xa0
		"................................" + // 0xc0
		"................................" //   0xe0
	//   0123456789abcdef0123456789abcdef
	uMap = "" +
		"................................" + // 0x00
		"................EEEEEEEEEE......" + // 0x20
		".EEEEEE........................." + // 0x40
		".EEEEEE........................." + // 0x60
		"................................" + // 0x80
		"................................" + // 0xa0
		"................................" + // 0xc0
		"................................" //   0xe0
	//   0123456789abcdef0123456789abcdef
	spaceMap = "" +
		".........ab..a.................." + // 0x00
		"a..............................." + // 0x20
		"................................" + // 0x40
		"................................" + // 0x60
		"................................" + // 0x80
		"................................" + // 0xa0
		"................................" + // 0xc0
		"................................s" //   0xe0
)
