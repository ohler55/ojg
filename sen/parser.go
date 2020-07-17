// Copyright (c) 2020, Peter Ohler, All rights reserved.

package sen

import (
	"fmt"
	"io"
	"unicode/utf8"

	"github.com/ohler55/ojg/gen"
	"github.com/ohler55/ojg/oj"
)

const (
	tmpMinSize  = 32 // for tokens and numbers
	readBufSize = 4096

	bomMode          = 'b'
	valueMode        = 'v'
	negMode          = '-'
	zeroMode         = '0'
	digitMode        = 'd'
	dotMode          = '.'
	fracMode         = 'F'
	expSignMode      = '+'
	expZeroMode      = 'X'
	expMode          = 'x'
	strMode          = 's'
	escMode          = 'e'
	uMode            = 'u'
	key1Mode         = 'k'
	tokenMode        = 't'
	tokenKeyMode     = 'K'
	colonMode        = ':'
	spaceMode        = ' '
	commentStartMode = '/'
	commentMode      = 'c'

	//   0123456789abcdef0123456789abcdef
	tokenMap = "" +
		"................................" + // 0x00
		"..............o.oooooooooo......" + // 0x20
		".oooooooooooooooooooooooooo...oo" + // 0x40
		".oooooooooooooooooooooooooo...o." + // 0x60
		"oooooooooooooooooooooooooooooooo" + // 0x80
		"oooooooooooooooooooooooooooooooo" + // 0xa0
		"oooooooooooooooooooooooooooooooo" + // 0xc0
		"oooooooooooooooooooooooooooooooo" //   0xe0

	//   0123456789abcdef0123456789abcdef
	strMap = "" +
		"................................" + // 0x00
		"oo.ooooooooooooooooooooooooooooo" + // 0x20
		"oooooooooooooooooooooooooooo.ooo" + // 0x40
		"ooooooooooooooooooooooooooooooo." + // 0x60
		"oooooooooooooooooooooooooooooooo" + // 0x80
		"oooooooooooooooooooooooooooooooo" + // 0xa0
		"oooooooooooooooooooooooooooooooo" + // 0xc0
		"oooooooooooooooooooooooooooooooo" //   0xe0

	//   0123456789abcdef0123456789abcdef
	charTypeMap = "" +
		".........ss....................." + // 0x00
		"s...........s..xdddddddddd......" + // 0x20
		"...........................x.x.." + // 0x40
		"...........................x.x.." + // 0x60
		"................................" + // 0x80
		"................................" + // 0xa0
		"................................" + // 0xc0
		"................................" //   0xe0
)

var emptySlice = []interface{}{}

// Parser a JSON parser. It can be reused for multiple parsings which allows
// buffer reuse for a performance advantage.
type Parser struct {
	tmp       []byte // used for numbers and strings
	runeBytes []byte
	stack     []interface{}
	starts    []int
	cb        func(interface{}) bool
	ri        int // read index for null, false, and true
	line      int
	noff      int // Offset of last newline from start of buf. Can be negative when using a reader.
	num       gen.Number
	rn        rune
	mode      byte
	nextMode  byte
	onlyOne   bool
}

// Parse a JSON string in to simple types. An error is returned if not valid JSON.
func (p *Parser) Parse(buf []byte, args ...interface{}) (data interface{}, err error) {
	var callback func(interface{}) bool

	for _, a := range args {
		switch ta := a.(type) {
		case func(interface{}) bool:
			callback = ta
			p.onlyOne = false
		default:
			return nil, fmt.Errorf("a %T is not a valid option type", a)
		}
	}
	if callback == nil {
		p.onlyOne = true
		callback = func(n interface{}) bool {
			data = n
			return false // tells the parser to stop
		}
	}
	p.cb = callback
	if cap(p.tmp) < tmpMinSize { // indicates not initialized
		p.tmp = make([]byte, 0, tmpMinSize)
		p.stack = make([]interface{}, 0, 64)
		p.starts = make([]int, 0, 16)
	} else {
		p.tmp = p.tmp[0:0]
		p.stack = p.stack[:0]
		p.starts = p.starts[:0]
	}
	p.noff = -1
	p.line = 1
	p.mode = valueMode
	// Skip BOM if present.
	if 0 < len(buf) && buf[0] == 0xEF {
		p.mode = bomMode
		p.ri = 0
	}
	err = p.parseBuffer(buf, true)
	for i := len(p.stack) - 1; 0 <= i; i-- {
		p.stack = nil
	}
	p.stack = p.stack[:0]
	return
}

// ParseReader a JSON io.Reader. An error is returned if not valid JSON.
func (p *Parser) ParseReader(r io.Reader, args ...interface{}) (node interface{}, err error) {
	var callback func(interface{}) bool

	for _, a := range args {
		switch ta := a.(type) {
		case func(interface{}) bool:
			callback = ta
			p.onlyOne = false
		default:
			return nil, fmt.Errorf("a %T is not a valid option type", a)
		}
	}
	if callback == nil {
		p.onlyOne = true
		callback = func(n interface{}) bool {
			node = n
			return false // tells the parser to stop
		}
	}
	p.cb = callback
	if cap(p.tmp) < tmpMinSize { // indicates not initialized
		p.tmp = make([]byte, 0, tmpMinSize)
		p.stack = make([]interface{}, 0, 64)
		p.starts = make([]int, 0, 16)
	} else {
		p.tmp = p.tmp[0:0]
		p.stack = p.stack[:0]
		p.starts = p.starts[:0]
	}
	p.noff = -1
	p.line = 1
	p.mode = valueMode
	buf := make([]byte, readBufSize)
	eof := false
	var cnt int
	cnt, err = r.Read(buf)
	buf = buf[:cnt]
	if err != nil {
		if err != io.EOF {
			return
		}
		eof = true
	}
	// Skip BOM if present.
	if 0 < len(buf) && buf[0] == 0xEF {
		p.mode = bomMode
		p.ri = 0
	}
	for {
		if err = p.parseBuffer(buf, eof); err != nil {
			for i := len(p.stack) - 1; 0 <= i; i-- {
				p.stack = nil
			}
			p.stack = p.stack[:0]

			return
		}
		if eof {
			break
		}
		buf = buf[:cap(buf)]
		cnt, err = r.Read(buf)
		buf = buf[:cnt]
		if err != nil {
			if err != io.EOF {
				return
			}
			eof = true
		}
	}
	return
}

func (p *Parser) parseBuffer(buf []byte, last bool) error {
	var b byte
	var i int
	var off int
	for off = 0; off < len(buf); off++ {
		b = buf[off]
		switch p.mode {
		case valueMode:
			switch b {
			case ' ', '\t', '\r', ',':
				// ignore and continue
			case '\n':
				p.line++
				p.noff = off
				for i, b = range buf[off+1:] {
					if charTypeMap[b] != 's' {
						break
					}
				}
				off += i
			case '-':
				p.mode = negMode
				p.num.Reset()
				p.num.Neg = true
			case '0':
				p.mode = zeroMode
				p.num.Reset()
			case '1', '2', '3', '4', '5', '6', '7', '8', '9':
				p.mode = digitMode
				p.num.Reset()
				p.num.I = uint64(b - '0')
			case '"':
				start := off + 1
				if len(buf) <= start {
					p.tmp = p.tmp[:0]
					p.mode = strMode
					if 0 < len(p.starts) && p.starts[len(p.starts)-1] == -1 {
						p.nextMode = key1Mode
					} else {
						p.nextMode = valueMode
					}
					continue
				}
				for i, b = range buf[start:] {
					if strMap[b] != 'o' {
						break
					}
				}
				off += i
				if b == '"' {
					off++
					p.iadd(string(buf[start:off]))
					if 0 < len(p.starts) && p.starts[len(p.starts)-1] == -1 {
						p.mode = key1Mode
					} else {
						p.mode = valueMode
					}
				} else {
					p.tmp = p.tmp[:0]
					p.tmp = append(p.tmp, buf[start:off+1]...)
					p.mode = strMode
					if 0 < len(p.starts) && p.starts[len(p.starts)-1] == -1 {
						p.nextMode = key1Mode
					} else {
						p.nextMode = valueMode
					}
				}
			case '[':
				p.starts = append(p.starts, len(p.stack))
				p.stack = append(p.stack, emptySlice)
			case ']':
				if err := p.arrayEnd(off); err != nil {
					return err
				}
			case '{':
				p.starts = append(p.starts, -1)
				p.mode = key1Mode
				n := map[string]interface{}{}
				p.stack = append(p.stack, n)
			case '}':
				if err := p.objectEnd(off); err != nil {
					return err
				}
			case '/':
				p.nextMode = p.mode
				p.mode = commentStartMode
			default:
				start := off
				for i, b = range buf[start:] {
					if tokenMap[b] != 'o' {
						break
					}
				}
				off += i
				ct := charTypeMap[b]
				if ct == 's' || ct == 'x' {
					switch {
					case i == 4 && buf[start] == 'n' && buf[start+1] == 'u' && buf[start+2] == 'l' && buf[start+3] == 'l':
						p.iadd(nil)
					case i == 4 && buf[start] == 't' && buf[start+1] == 'r' && buf[start+2] == 'u' && buf[start+3] == 'e':
						p.iadd(true)
					case i == 5 && buf[start] == 'f' && buf[start+1] == 'a' && buf[start+2] == 'l' && buf[start+3] == 's' && buf[start+4] == 'e':
						p.iadd(false)
					default:
						p.iadd(string(buf[start:off]))
					}
					switch b {
					case '[':
						p.starts = append(p.starts, len(p.stack))
						p.stack = append(p.stack, emptySlice)
						p.mode = valueMode
					case ']':
						if err := p.arrayEnd(off); err != nil {
							return err
						}
					case '{':
						p.starts = append(p.starts, -1)
						p.mode = key1Mode
						n := map[string]interface{}{}
						p.stack = append(p.stack, n)
					case '}':
						if err := p.objectEnd(off); err != nil {
							return err
						}
					case '/':
						p.nextMode = p.mode
						p.mode = commentStartMode
					default:
						if 0 < len(p.starts) && p.starts[len(p.starts)-1] == -1 {
							p.mode = key1Mode
						} else {
							p.mode = valueMode
						}
					}
				} else if tokenMap[b] == 'o' {
					// Must be end of buffer.
					p.tmp = p.tmp[:0]
					p.tmp = append(p.tmp, buf[start:off+1]...)
					p.mode = tokenMode
				} else if tokenMap[b] != 'o' {
					return p.newError(off, "expected a token, not '%c'", b)
				}
			}
		case tokenMode:
			if tokenMap[b] == 'o' {
				p.tmp = append(p.tmp, b)
			} else {
				switch {
				case len(p.tmp) == 4 && p.tmp[0] == 'n' && p.tmp[1] == 'u' && p.tmp[2] == 'l' && p.tmp[3] == 'l':
					p.iadd(nil)
				case len(p.tmp) == 4 && p.tmp[0] == 't' && p.tmp[1] == 'r' && p.tmp[2] == 'u' && p.tmp[3] == 'e':
					p.iadd(true)
				case len(p.tmp) == 5 && p.tmp[0] == 'f' && p.tmp[1] == 'a' && p.tmp[2] == 'l' && p.tmp[3] == 's' && p.tmp[4] == 'e':
					p.iadd(false)
				default:
					p.iadd(string(p.tmp))
				}
				if 0 < len(p.starts) && p.starts[len(p.starts)-1] == -1 {
					p.mode = key1Mode
				} else {
					p.mode = valueMode
				}
				switch b {
				case ' ', '\t', '\r', ',':
				case '\n':
					p.line++
					p.noff = off
					for i, b = range buf[off+1:] {
						if charTypeMap[b] != 's' {
							break
						}
					}
					off += i
				case '[':
					p.starts = append(p.starts, len(p.stack))
					p.stack = append(p.stack, emptySlice)
				case '{':
					p.starts = append(p.starts, -1)
					p.mode = key1Mode
					n := map[string]interface{}{}
					p.stack = append(p.stack, n)
				case ']':
					if err := p.arrayEnd(off); err != nil {
						return err
					}
				case '}':
					if err := p.objectEnd(off); err != nil {
						return err
					}
				case '/':
					p.nextMode = p.mode
					p.mode = commentStartMode
				default:
					return p.newError(off, "expected a value, not '%c'", b)
				}
			}
		case key1Mode:
			switch b {
			case ' ', '\t', '\r', ',':
				continue
			case '\n':
				p.line++
				p.noff = off
				for i, b = range buf[off+1:] {
					if charTypeMap[b] != 's' {
						break
					}
				}
				off += i
			case '"':
				start := off + 1
				if len(buf) <= start {
					p.tmp = p.tmp[:0]
					p.mode = strMode
					p.nextMode = colonMode
					continue
				}
				for i, b = range buf[start:] {
					if strMap[b] != 'o' {
						break
					}
				}
				off += i
				if b == '"' {
					off++
					p.stack = append(p.stack, gen.Key(buf[start:off]))
					p.mode = colonMode
				} else {
					p.tmp = p.tmp[:0]
					p.tmp = append(p.tmp, buf[start:off+1]...)
					p.mode = strMode
					p.nextMode = colonMode
				}
			case '}':
				// If in key mode } is always okay.
				_ = p.objectEnd(off)
			default:
				start := off
				for i, b = range buf[start:] {
					if tokenMap[b] != 'o' {
						break
					}
				}
				off += i
				if b == ':' {
					p.mode = valueMode
					p.stack = append(p.stack, gen.Key(buf[start:off]))
				} else if charTypeMap[b] == 's' {
					p.mode = colonMode
					p.stack = append(p.stack, gen.Key(buf[start:off]))
				} else if tokenMap[b] == 'o' {
					// Must be end of buffer.
					p.tmp = p.tmp[:0]
					p.tmp = append(p.tmp, buf[start:off+1]...)
					p.mode = tokenKeyMode
				} else if tokenMap[b] != 'o' {
					return p.newError(off, "expected a token followed by a ':', not '%c'", b)
				}
			}
		case tokenKeyMode:
			if tokenMap[b] == 'o' {
				p.tmp = append(p.tmp, b)
			} else {
				p.stack = append(p.stack, gen.Key(p.tmp))
				p.mode = colonMode
				switch b {
				case ' ', '\t', '\r', ',':
				case '\n':
					p.line++
					p.noff = off
					for i, b = range buf[off+1:] {
						if charTypeMap[b] != 's' {
							break
						}
					}
					off += i
				case ':':
					p.mode = valueMode
				default:
					return p.newError(off, "expected a token character, not '%c'", b)
				}
			}
		case colonMode:
			switch b {
			case ' ', '\t', '\r':
				continue
			case '\n':
				p.line++
				p.noff = off
				for i, b = range buf[off+1:] {
					if charTypeMap[b] != 's' {
						break
					}
				}
				off += i
			case ':':
				p.mode = valueMode
			default:
				return p.newError(off, "expected a colon, not '%c'", b)
			}
		case negMode:
			switch b {
			case '0':
				p.mode = zeroMode
			case '1', '2', '3', '4', '5', '6', '7', '8', '9':
				p.mode = digitMode
				p.num.AddDigit(b)
			default:
				return p.newError(off, "invalid number")
			}
		case zeroMode:
			switch b {
			case '.':
				p.mode = dotMode
			case ' ', '\t', '\r', ',':
				p.appendNum()
			case '\n':
				p.line++
				p.noff = off
				p.appendNum()
				for i, b = range buf[off+1:] {
					if charTypeMap[b] != 's' {
						break
					}
				}
				off += i
			case '[':
				p.appendNum()
				p.starts = append(p.starts, len(p.stack))
				p.stack = append(p.stack, emptySlice)
			case ']':
				p.appendNum()
				if err := p.arrayEnd(off); err != nil {
					return err
				}
			case '{':
				p.appendNum()
				p.starts = append(p.starts, -1)
				p.mode = key1Mode
				n := map[string]interface{}{}
				p.stack = append(p.stack, n)
			case '}':
				p.appendNum()
				if err := p.objectEnd(off); err != nil {
					return err
				}
			default:
				return p.newError(off, "invalid number")
			}
		case digitMode:
			switch b {
			case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
				p.num.AddDigit(b)
			case '.':
				p.mode = dotMode
				if 0 < len(p.num.BigBuf) {
					p.num.BigBuf = append(p.num.BigBuf, b)
				}
			case ' ', '\t', '\r', ',':
				p.appendNum()
			case '\n':
				p.line++
				p.noff = off
				p.appendNum()
				for i, b = range buf[off+1:] {
					if charTypeMap[b] != 's' {
						break
					}
				}
				off += i
			case '[':
				p.appendNum()
				p.starts = append(p.starts, len(p.stack))
				p.stack = append(p.stack, emptySlice)
			case ']':
				p.appendNum()
				if err := p.arrayEnd(off); err != nil {
					return err
				}
			case '{':
				p.appendNum()
				p.starts = append(p.starts, -1)
				p.mode = key1Mode
				n := map[string]interface{}{}
				p.stack = append(p.stack, n)
			case '}':
				p.appendNum()
				if err := p.objectEnd(off); err != nil {
					return err
				}
			default:
				return p.newError(off, "invalid number")
			}
		case dotMode:
			if '0' <= b && b <= '9' {
				p.mode = fracMode
				p.num.AddFrac(b)
			} else {
				return p.newError(off, "invalid number")
			}
		case fracMode:
			switch b {
			case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
				p.num.AddFrac(b)
			case 'e', 'E':
				p.mode = expSignMode
				if 0 < len(p.num.BigBuf) {
					p.num.BigBuf = append(p.num.BigBuf, b)
				}
			case ' ', '\t', '\r', ',':
				p.appendNum()
			case '\n':
				p.line++
				p.noff = off
				p.appendNum()
				for i, b = range buf[off+1:] {
					if charTypeMap[b] != 's' {
						break
					}
				}
				off += i
			case '[':
				p.appendNum()
				p.starts = append(p.starts, len(p.stack))
				p.stack = append(p.stack, emptySlice)
			case ']':
				p.appendNum()
				if err := p.arrayEnd(off); err != nil {
					return err
				}
			case '{':
				p.appendNum()
				p.starts = append(p.starts, -1)
				p.mode = key1Mode
				n := map[string]interface{}{}
				p.stack = append(p.stack, n)
			case '}':
				p.appendNum()
				if err := p.objectEnd(off); err != nil {
					return err
				}
			default:
				return p.newError(off, "invalid number")
			}
		case expSignMode:
			switch b {
			case '-':
				p.mode = expZeroMode
				p.num.NegExp = true
			case '+':
				p.mode = expZeroMode
			case '1', '2', '3', '4', '5', '6', '7', '8', '9':
				p.mode = expMode
				p.num.AddExp(b)
			default:
				return p.newError(off, "invalid number")
			}
		case expZeroMode:
			if '0' <= b && b <= '9' {
				p.mode = expMode
				p.num.AddExp(b)
			} else {
				return p.newError(off, "invalid number")
			}
		case expMode:
			switch b {
			case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
				p.num.AddExp(b)
			case ' ', '\t', '\r', ',':
				p.appendNum()
			case '\n':
				p.line++
				p.noff = off
				p.appendNum()
				for i, b = range buf[off+1:] {
					if charTypeMap[b] != 's' {
						break
					}
				}
				off += i
			case '[':
				p.appendNum()
				p.starts = append(p.starts, len(p.stack))
				p.stack = append(p.stack, emptySlice)
			case ']':
				p.appendNum()
				if err := p.arrayEnd(off); err != nil {
					return err
				}
			case '{':
				p.appendNum()
				p.starts = append(p.starts, -1)
				p.mode = key1Mode
				n := map[string]interface{}{}
				p.stack = append(p.stack, n)
			case '}':
				p.appendNum()
				if err := p.objectEnd(off); err != nil {
					return err
				}
			default:
				return p.newError(off, "invalid number")
			}
		case strMode:
			if b < 0x20 {
				return p.newError(off, "invalid JSON character 0x%02x", b)
			}
			switch b {
			case '\\':
				p.mode = escMode
			case '"':
				p.mode = p.nextMode
				if p.mode == colonMode {
					p.stack = append(p.stack, gen.Key(p.tmp))
				} else {
					p.iadd(string(p.tmp))
				}
			default:
				p.tmp = append(p.tmp, b)
			}
		case escMode:
			p.mode = strMode
			switch b {
			case 'n':
				p.tmp = append(p.tmp, '\n')
			case '"':
				p.tmp = append(p.tmp, '"')
			case '\\':
				p.tmp = append(p.tmp, '\\')
			case '/':
				p.tmp = append(p.tmp, '/')
			case 'b':
				p.tmp = append(p.tmp, '\b')
			case 'f':
				p.tmp = append(p.tmp, '\f')
			case 'r':
				p.tmp = append(p.tmp, '\r')
			case 't':
				p.tmp = append(p.tmp, '\t')
			case 'u':
				p.mode = uMode
				p.rn = 0
				p.ri = 0
			default:
				return p.newError(off, "invalid JSON escape character '\\%c'", b)
			}
		case uMode:
			p.ri++
			switch b {
			case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
				p.rn = p.rn<<4 | rune(b-'0')
			case 'a', 'b', 'c', 'd', 'e', 'f':
				p.rn = p.rn<<4 | rune(b-'a'+10)
			case 'A', 'B', 'C', 'D', 'E', 'F':
				p.rn = p.rn<<4 | rune(b-'A'+10)
			default:
				return p.newError(off, "invalid JSON unicode character '%c'", b)
			}
			if p.ri == 4 {
				if len(p.runeBytes) < 6 {
					p.runeBytes = make([]byte, 6)
				}
				n := utf8.EncodeRune(p.runeBytes, p.rn)
				p.tmp = append(p.tmp, p.runeBytes[:n]...)
				p.mode = strMode
			}
		case spaceMode:
			switch b {
			case ' ', '\t', '\r':
				continue
			case '\n':
				p.line++
				p.noff = off
				for i, b = range buf[off+1:] {
					if charTypeMap[b] != 's' {
						break
					}
				}
				off += i
			default:
				return p.newError(off, "extra characters after close, '%c'", b)
			}
		case commentStartMode:
			if b != '/' {
				return p.newError(off, "unexpected character '%c'", b)
			}
			p.mode = commentMode
		case commentMode:
			if b == '\n' {
				p.line++
				p.noff = off
				p.mode = p.nextMode
			}
		case bomMode:
			if []byte{0xEF, 0xBB, 0xBF}[p.ri] != b {
				return p.newError(off, "expected BOM")
			}
			p.ri++
			if 3 <= p.ri {
				p.mode = valueMode
			}
		}
		if len(p.starts) == 0 && p.mode == valueMode && 0 < len(p.stack) {
			p.cb(p.stack[0])
			p.stack[0] = nil
			p.stack = p.stack[:0]
			if p.onlyOne {
				p.mode = spaceMode
			} else {
				p.mode = valueMode
			}
		}
	}
	if last {
		if 0 < len(p.starts) {
			return p.newError(off, "not closed")
		}
		switch p.mode {
		case tokenMode:
			switch {
			case len(p.tmp) == 4 && p.tmp[0] == 'n' && p.tmp[1] == 'u' && p.tmp[2] == 'l' && p.tmp[3] == 'l':
				p.iadd(nil)
			case len(p.tmp) == 4 && p.tmp[0] == 't' && p.tmp[1] == 'r' && p.tmp[2] == 'u' && p.tmp[3] == 'e':
				p.iadd(true)
			case len(p.tmp) == 5 && p.tmp[0] == 'f' && p.tmp[1] == 'a' && p.tmp[2] == 'l' && p.tmp[3] == 's' && p.tmp[4] == 'e':
				p.iadd(false)
			default:
				p.iadd(string(p.tmp))
			}
			if 0 < len(p.stack) {
				p.cb(p.stack[0])
			}
		case valueMode:
			// normal completion
		case zeroMode, digitMode, fracMode, expMode:
			p.appendNum()
			if 0 < len(p.stack) {
				p.cb(p.stack[0])
			}
		case spaceMode:
			// just reading white space
		default:
			//fmt.Printf("*** final mode: %c\n", p.mode)
			return p.newError(off, "incomplete JSON")
		}
	}
	return nil
}

func (p *Parser) newError(off int, format string, args ...interface{}) error {
	return &oj.ParseError{
		Message: fmt.Sprintf(format, args...),
		Line:    p.line,
		Column:  off - p.noff,
	}
}

func (p *Parser) iadd(n interface{}) {
	if 2 <= len(p.stack) {
		if k, ok := p.stack[len(p.stack)-1].(gen.Key); ok {
			obj, _ := p.stack[len(p.stack)-2].(map[string]interface{})
			obj[string(k)] = n
			p.stack = p.stack[0 : len(p.stack)-1]

			return
		}
	}
	p.stack = append(p.stack, n)
}

func (p *Parser) appendNum() {
	if 0 < len(p.num.BigBuf) {
		p.iadd(string(p.num.AsBig()))
	} else if p.num.Frac == 0 && p.num.Exp == 0 {
		p.iadd(p.num.AsInt())
	} else {
		p.iadd(p.num.AsFloat())
	}
	if 0 < len(p.starts) && p.starts[len(p.starts)-1] == -1 {
		p.mode = key1Mode
	} else {
		p.mode = valueMode
	}
}

func (p *Parser) arrayEnd(off int) error {
	depth := len(p.starts)
	if depth == 0 {
		return p.newError(off, "too many closes")
	}
	depth--
	if p.starts[depth] < 0 {
		return p.newError(off, "unexpected array close")
	}
	p.mode = valueMode
	start := p.starts[len(p.starts)-1] + 1
	p.starts = p.starts[:len(p.starts)-1]
	size := len(p.stack) - start
	n := make([]interface{}, size)
	copy(n, p.stack[start:len(p.stack)])
	p.stack = p.stack[0 : start-1]
	p.iadd(n)
	if 0 < len(p.starts) && p.starts[len(p.starts)-1] == -1 {
		p.mode = key1Mode
	} else {
		p.mode = valueMode
	}
	return nil
}

func (p *Parser) objectEnd(off int) error {
	depth := len(p.starts)
	if depth == 0 {
		return p.newError(off, "too many closes")
	}
	depth--
	if 0 <= p.starts[depth] {
		return p.newError(off, "unexpected object close")
	}
	p.starts = p.starts[0:depth]
	p.mode = valueMode
	n := p.stack[len(p.stack)-1]
	p.stack = p.stack[:len(p.stack)-1]
	p.iadd(n)
	if 0 < len(p.starts) && p.starts[len(p.starts)-1] == -1 {
		p.mode = key1Mode
	} else {
		p.mode = valueMode
	}
	return nil
}
