// Copyright (c) 2020, Peter Ohler, All rights reserved.

package oj

const (
	skipChar    = 0x00
	skipNewline = 0x01

	valNull  = 0x02
	valTrue  = 0x03
	valFalse = 0x04
	valNeg   = 0x05
	val0     = 0x06
	valDigit = 0x07
	valQuote = 0x08
	valSlash = 0x09

	openArray   = 0x10
	openObject  = 0x11
	closeArray  = 0x12
	closeObject = 0x13

	afterComma = 0x20
	keyQuote   = 0x21
	colonColon = 0x22

	numSpc     = 0x30
	numNewline = 0x31
	numDot     = 0x32
	numComma   = 0x33
	numFrac    = 0x34
	fracE      = 0x35
	expSign    = 0x36
	expDigit   = 0x37

	strQuote = 0x40
	strSlash = 0x41
	escOk    = 0x42
	escU     = 0x43
	uOk      = 0x44

	nullOk  = 0x60
	trueOk  = 0x61
	falseOk = 0x62

	commentStart = 0x70
	commentEnd   = 0x71

	bomBB = 0xe0
	bomBF = 0xe1
	bomEF = 0xe2

	bomErr     = 0xf1
	valErr     = 0xf2
	trueErr    = 0xf3
	nullErr    = 0xf4
	falseErr   = 0xf5
	afterErr   = 0xf6
	key1Err    = 0xf7
	keyErr     = 0xf8
	colonErr   = 0xf9
	numErr     = 0xfa
	strLowErr  = 0xfb
	strErr     = 0xfc
	escErr     = 0xfd
	spcErr     = 0xfe
	commentErr = 0xff
)

var (
	bomBBmap        = [257]byte{}
	bomBFmap        = [257]byte{}
	valueMap        = [257]byte{}
	afterMap        = [257]byte{}
	nullMap         = [257]byte{}
	trueMap         = [257]byte{}
	falseMap        = [257]byte{}
	negMap          = [257]byte{}
	zeroMap         = [257]byte{}
	digitMap        = [257]byte{}
	dotMap          = [257]byte{}
	fracMap         = [257]byte{}
	expSignMap      = [257]byte{}
	expZeroMap      = [257]byte{}
	expMap          = [257]byte{}
	stringMap       = [257]byte{}
	escMap          = [257]byte{}
	uMap            = [257]byte{}
	key1Map         = [257]byte{}
	keyMap          = [257]byte{}
	colonMap        = [257]byte{}
	commaMap        = [257]byte{}
	spaceMap        = [257]byte{}
	commentStartMap = [257]byte{}
	commentMap      = [257]byte{}
)

func init() {
	for i := 0; i < 256; i++ {
		bomBBmap[i] = bomErr
		bomBFmap[i] = bomErr
		valueMap[i] = valErr
		nullMap[i] = nullErr
		trueMap[i] = trueErr
		falseMap[i] = falseErr
		afterMap[i] = afterErr
		key1Map[i] = key1Err
		keyMap[i] = keyErr
		colonMap[i] = colonErr
		negMap[i] = numErr
		zeroMap[i] = numErr
		digitMap[i] = numErr
		dotMap[i] = numErr
		fracMap[i] = numErr
		expSignMap[i] = numErr
		expZeroMap[i] = numErr
		expMap[i] = numErr
		escMap[i] = escErr
		uMap[i] = strErr
		spaceMap[i] = spcErr
		commentStartMap[i] = commentErr
		commentMap[i] = skipChar
	}
	bomBBmap[0xBF] = bomBF
	bomBFmap[0xBB] = bomBB

	nullMap['u'] = nullOk
	nullMap['l'] = nullOk
	trueMap['r'] = trueOk
	trueMap['u'] = trueOk
	trueMap['e'] = trueOk
	falseMap['r'] = falseOk
	falseMap['u'] = falseOk
	falseMap['e'] = falseOk

	valueMap[' '] = skipChar
	valueMap['\t'] = skipChar
	valueMap['\r'] = skipChar
	valueMap['\n'] = skipNewline
	valueMap['n'] = valNull
	valueMap['t'] = valTrue
	valueMap['f'] = valFalse
	valueMap['-'] = valNeg
	valueMap['0'] = val0
	for b := '1'; b <= '9'; b++ {
		valueMap[b] = valDigit
	}
	valueMap['"'] = valQuote
	valueMap['['] = openArray
	valueMap['{'] = openObject
	valueMap[']'] = closeArray
	valueMap['}'] = closeObject
	valueMap['/'] = valSlash

	commaMap = valueMap
	commaMap[']'] = valErr
	commaMap['}'] = valErr

	afterMap[' '] = skipChar
	afterMap['\t'] = skipChar
	afterMap['\r'] = skipChar
	afterMap['\n'] = skipNewline
	afterMap[','] = afterComma
	afterMap[']'] = closeArray
	afterMap['}'] = closeObject

	keyMap[' '] = skipChar
	keyMap['\t'] = skipChar
	keyMap['\r'] = skipChar
	keyMap['\n'] = skipNewline
	keyMap['"'] = keyQuote
	key1Map = keyMap
	key1Map['}'] = closeObject

	colonMap[' '] = skipChar
	colonMap['\t'] = skipChar
	colonMap['\r'] = skipChar
	colonMap['\n'] = skipNewline
	colonMap[':'] = colonColon

	negMap['0'] = val0
	for b := '1'; b <= '9'; b++ {
		negMap[b] = valDigit
	}

	zeroMap[' '] = numSpc
	zeroMap['\t'] = numSpc
	zeroMap['\r'] = numSpc
	zeroMap['\n'] = numNewline
	zeroMap['.'] = numDot
	zeroMap[','] = numComma
	zeroMap[']'] = closeArray
	zeroMap['}'] = closeObject

	digitMap = zeroMap
	for b := '0'; b <= '9'; b++ {
		digitMap[b] = skipChar
	}

	for b := '0'; b <= '9'; b++ {
		dotMap[b] = numFrac
	}

	fracMap = zeroMap
	for b := '0'; b <= '9'; b++ {
		fracMap[b] = skipChar
	}
	fracMap['e'] = fracE

	expSignMap['-'] = expSign
	expSignMap['+'] = expSign
	for b := '0'; b <= '9'; b++ {
		expSignMap[b] = expDigit
	}
	for b := '0'; b <= '9'; b++ {
		expZeroMap[b] = expDigit
	}
	expMap = zeroMap
	expMap['.'] = numErr

	for i := 0; i < 0x20; i++ {
		stringMap[i] = strLowErr
	}
	for i := 0x20; i <= 0xff; i++ {
		stringMap[i] = skipChar
	}
	stringMap['"'] = strQuote
	stringMap['\\'] = strSlash

	escMap['u'] = escU
	escMap['n'] = escOk
	escMap['"'] = escOk
	escMap['\\'] = escOk
	escMap['/'] = escOk
	escMap['b'] = escOk
	escMap['f'] = escOk
	escMap['r'] = escOk
	escMap['t'] = escOk

	for b := '0'; b <= '9'; b++ {
		uMap[b] = uOk
	}
	for b := 'a'; b <= 'f'; b++ {
		uMap[b] = uOk
	}
	for b := 'A'; b <= 'F'; b++ {
		uMap[b] = uOk
	}

	spaceMap[' '] = skipChar
	spaceMap['\t'] = skipChar
	spaceMap['\r'] = skipChar
	spaceMap['\n'] = skipNewline

	commentStartMap['/'] = commentStart
	commentMap['\n'] = commentEnd

	initMapIDs()
}

func initMapIDs() {
	valueMap[256] = 'v'
	afterMap[256] = 'a'
	nullMap[256] = 'n'
	trueMap[256] = 't'
	falseMap[256] = 'f'
	negMap[256] = '-'
	zeroMap[256] = '0'
	digitMap[256] = 'd'
	dotMap[256] = '.'
	fracMap[256] = 'F'
	expSignMap[256] = '+'
	expZeroMap[256] = 'X'
	expMap[256] = 'x'
	stringMap[256] = 's'
	escMap[256] = 'e'
	uMap[256] = 'u'
	key1Map[256] = 'K'
	keyMap[256] = 'k'
	colonMap[256] = ':'
	commaMap[256] = ','
	spaceMap[256] = ' '
	commentStartMap[256] = '/'
	commentMap[256] = 'c'
}
