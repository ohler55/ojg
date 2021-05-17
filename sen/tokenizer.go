// Copyright (c) 2021, Peter Ohler, All rights reserved.

package sen

import (
	"encoding/json"
	"fmt"
	"io"
	"math"
	"unicode/utf8"

	"github.com/ohler55/ojg"
	"github.com/ohler55/ojg/gen"
	"github.com/ohler55/ojg/oj"
)

const (
	objectStart = '{'
	arrayStart  = '['
)

// Tokenizer is a reusable JSON tokenizer. It can be reused for multiple parsings
// which allows buffer reuse for a performance advantage.
type Tokenizer struct {
	tmp       []byte // used for numbers and strings
	runeBytes []byte
	starts    []byte
	handler   oj.TokenHandler
	line      int
	noff      int // Offset of last newline from start of buf. Can be negative when using a reader.
	ri        int // read index for null, false, and true
	mi        int
	num       gen.Number
	rn        rune
	mode      string
	exkey     bool

	// OnlyOne returns an error if more than one JSON is in the string or stream.
	OnlyOne bool
}

// TokenizeString the provided JSON and call the handler functions for each
// token in the JSON.
func TokenizeString(data string, handler oj.TokenHandler) error {
	t := Tokenizer{}
	return t.Parse([]byte(data), handler)
}

// Parse a JSON string in to simple types. An error is returned if not valid JSON.
func Tokenize(data []byte, handler oj.TokenHandler) error {
	t := Tokenizer{}
	return t.Parse(data, handler)
}

// TokenizeLoad a JSON io.Reader. An error is returned if not valid JSON.
func TokenizeLoad(r io.Reader, handler oj.TokenHandler) error {
	t := Tokenizer{}
	return t.Load(r, handler)
}

// Parse a JSON string in to simple types. An error is returned if not valid JSON.
func (t *Tokenizer) Parse(buf []byte, handler oj.TokenHandler) (err error) {
	t.handler = handler
	if t.starts == nil {
		t.tmp = make([]byte, 0, tmpInitSize)
		t.starts = make([]byte, 0, 16)
	} else {
		t.tmp = t.tmp[:0]
		t.starts = t.starts[:0]
	}
	t.noff = -1
	t.line = 1
	t.mode = valueMap
	t.mi = 0
	defer func() {
		if r := recover(); r != nil {
			err = ojg.NewError(r)
		}
	}()
	// Skip BOM if present.
	if 3 < len(buf) && buf[0] == 0xEF {
		if buf[1] == 0xBB && buf[2] == 0xBF {
			t.tokenizeBuffer(buf[3:], true)
		} else {
			return fmt.Errorf("expected BOM at 1:3")
		}
	} else {
		t.tokenizeBuffer(buf, true)
	}
	return
}

// Load a JSON io.Reader. An error is returned if not valid JSON.
func (t *Tokenizer) Load(r io.Reader, handler oj.TokenHandler) (err error) {
	t.handler = handler
	if t.starts == nil {
		t.tmp = make([]byte, 0, tmpInitSize)
		t.starts = make([]byte, 0, 16)
	} else {
		t.tmp = t.tmp[:0]
		t.starts = t.starts[:0]
	}
	t.noff = -1
	t.line = 1
	t.mi = 0
	buf := make([]byte, readBufSize)
	eof := false
	defer func() {
		if r := recover(); r != nil {
			err = ojg.NewError(r)
		}
	}()
	var cnt int
	cnt, err = r.Read(buf)
	buf = buf[:cnt]
	t.mode = valueMap
	if err != nil {
		if err != io.EOF {
			return
		}
		eof = true
		err = nil
	}
	var skip int
	// Skip BOM if present.
	if 3 < len(buf) && buf[0] == 0xEF && buf[1] == 0xBB && buf[2] == 0xBF {
		skip = 3
	}
	for {
		if 0 < skip {
			t.tokenizeBuffer(buf[skip:], eof)
		} else {
			t.tokenizeBuffer(buf, eof)
		}
		skip = 0
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
			err = nil
		}
	}
	return
}

func (t *Tokenizer) tokenizeBuffer(buf []byte, last bool) {
	var b byte
	var i int
	var off int
	depth := len(t.starts)
	for off = 0; off < len(buf); off++ {
		b = buf[off]
		switch t.mode[b] {
		case skipNewline:
			t.line++
			t.noff = off
			for i, b = range buf[off+1:] {
				if spaceMap[b] != skipChar {
					break
				}
			}
			off += i
			continue
		case tokenStart:
			start := off
			for i, b = range buf[off:] {
				if tokenMap[b] != tokenOk {
					break
				}
			}
			off += i
			if tokenMap[b] == tokenOk { // end of buf reached
				t.tmp = t.tmp[:0]
				t.tmp = append(t.tmp, buf[start:off+1]...)
				t.mode = tokenMap
				continue
			}
			t.addToken(string(buf[start:off]))
			off--
		case strOk:
			t.tmp = append(t.tmp, b)
		case colonColon:
			t.mode = valueMap
			continue
		case skipChar: // skip and continue
			continue
		case openObject:
			if 256 < len(t.mode) {
				switch t.mode[256] {
				case 'n':
					t.handleNum(off)
				case 't':
					t.addToken(string(t.tmp))
				}
			}
			if t.exkey {
				t.newError(off, "expected a key")
			}
			t.starts = append(t.starts, objectStart)
			t.handler.ObjectStart()
			t.exkey = true
			depth++
			continue
		case closeObject:
			depth--
			if depth < 0 || t.starts[depth] != objectStart {
				t.newError(off, "unexpected object close")
			}
			if 256 < len(t.mode) {
				switch t.mode[256] {
				case 'n':
					t.handleNum(off)
				case 't':
					t.addToken(string(t.tmp))
				}
			}
			t.starts = t.starts[0:depth]
			t.handler.ObjectEnd()
			t.exkey = 0 < len(t.starts) && t.starts[len(t.starts)-1] == objectStart
		case valDigit:
			t.num.Reset()
			t.mode = digitMap
			t.num.I = uint64(b - '0')
			for i, b = range buf[off+1:] {
				if digitMap[b] != numDigit {
					break
				}
				t.num.I = t.num.I*10 + uint64(b-'0')
				if math.MaxInt64 < t.num.I {
					t.num.FillBig()
					break
				}
			}
			if digitMap[b] == numDigit {
				off++
			}
			off += i
		case valQuote:
			start := off + 1
			if len(buf) <= start {
				t.tmp = t.tmp[:0]
				t.mode = stringMap
				continue
			}
			for i, b = range buf[off+1:] {
				if stringMap[b] != strOk {
					break
				}
			}
			off += i
			if b == '"' {
				off++
				t.addString(string(buf[start:off]))
			} else {
				t.tmp = t.tmp[:0]
				t.tmp = append(t.tmp, buf[start:off+1]...)
				t.mode = stringMap
				continue
			}
		case numSpc:
			t.handleNum(off)
		case strSlash:
			t.mode = escMap
			continue
		case escOk:
			t.tmp = append(t.tmp, escByteMap[b])
			t.mode = stringMap
			continue
		case val0:
			t.mode = zeroMap
			t.num.Reset()
		case valNeg:
			t.mode = negMap
			t.num.Reset()
			t.num.Neg = true
			continue
		case escU:
			t.mode = uMap
			t.rn = 0
			t.ri = 0
			continue
		case openArray:
			if 256 < len(t.mode) {
				switch t.mode[256] {
				case 'n':
					t.handleNum(off)
				case 't':
					t.addToken(string(t.tmp))
				}
			}
			if t.exkey {
				t.newError(off, "expected a key")
			}
			t.starts = append(t.starts, arrayStart)
			t.handler.ArrayStart()
			t.mode = valueMap
			depth++
			continue
		case closeArray:
			depth--
			if depth < 0 || t.starts[depth] != arrayStart {
				t.newError(off, "unexpected array close")
			}
			// Only modes with a close array are value, token, and numbers
			// which are all over 256 long.
			switch t.mode[256] {
			case 'n':
				t.handleNum(off)
			case 't':
				t.addToken(string(t.tmp))
			}
			t.starts = t.starts[:len(t.starts)-1]
			t.handler.ArrayEnd()
			t.exkey = 0 < len(t.starts) && t.starts[len(t.starts)-1] == objectStart
			t.mode = valueMap
		case numDot:
			if 0 < len(t.num.BigBuf) {
				t.num.BigBuf = append(t.num.BigBuf, b)
				t.mode = dotMap
				continue
			}
			for i, b = range buf[off+1:] {
				if digitMap[b] != numDigit {
					break
				}
				t.num.Frac = t.num.Frac*10 + uint64(b-'0')
				t.num.Div *= 10.0
				if math.MaxInt64 < t.num.Frac {
					t.num.FillBig()
					break
				}
			}
			off += i
			if digitMap[b] == numDigit {
				off++
			}
			t.mode = fracMap
		case numFrac:
			t.num.AddFrac(b)
			t.mode = fracMap
		case fracE:
			if 0 < len(t.num.BigBuf) {
				t.num.BigBuf = append(t.num.BigBuf, b)
			}
			t.mode = expSignMap
			continue
		case tokenOk:
			t.tmp = append(t.tmp, b)
		case tokenSpc:
			t.addToken(string(t.tmp))
		case tokenColon:
			t.addToken(string(t.tmp))
			t.mode = valueMap
		case tokenNlColon:
			t.addToken(string(t.tmp))
			t.line++
			t.noff = off
			for i, b = range buf[off+1:] {
				if spaceMap[b] != skipChar {
					break
				}
			}
			off += i
		case strQuote:
			t.addString(string(t.tmp))
		case numZero:
			t.mode = zeroMap
		case numDigit:
			t.num.AddDigit(b)
		case negDigit:
			t.num.AddDigit(b)
			t.mode = digitMap
		case numNewline:
			t.handleNum(off)
			t.line++
			t.noff = off
			t.mode = valueMap
			for i, b = range buf[off+1:] {
				if spaceMap[b] != skipChar {
					break
				}
			}
			off += i
		case expSign:
			t.mode = expZeroMap
			if b == '-' {
				t.num.NegExp = true
			}
			continue
		case expDigit:
			t.num.AddExp(b)
			t.mode = expMap
		case uOk:
			t.ri++
			switch b {
			case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
				t.rn = t.rn<<4 | rune(b-'0')
			case 'a', 'b', 'c', 'd', 'e', 'f':
				t.rn = t.rn<<4 | rune(b-'a'+10)
			case 'A', 'B', 'C', 'D', 'E', 'F':
				t.rn = t.rn<<4 | rune(b-'A'+10)
			}
			if t.ri == 4 {
				if len(t.runeBytes) < 6 {
					t.runeBytes = make([]byte, 6)
				}
				n := utf8.EncodeRune(t.runeBytes, t.rn)
				t.tmp = append(t.tmp, t.runeBytes[:n]...)
				t.mode = stringMap
			}
			continue
		case valSlash:
			if 256 < len(t.mode) {
				switch t.mode[256] {
				case 'n':
					t.handleNum(off)
				case 't':
					t.addToken(string(t.tmp))
				}
			}
			t.mode = commentStartMap
		case commentStart:
			t.mode = commentMap
		case commentEnd:
			t.mode = valueMap
		case charErr:
			t.byteError(off, t.mode, b)
		}
		if depth == 0 && 256 < len(t.mode) && t.mode[256] == 'v' {
			t.mi = 0
			if t.OnlyOne {
				t.mode = spaceMap
			} else {
				t.mode = valueMap
			}
		}
	}
	if last {
		if 0 < len(t.starts) {
			t.newError(off, "not closed")
		}
		if len(t.mode) == 256 { // valid finishing maps are one byte longer
			t.newError(off, "incomplete JSON")
		}
		switch t.mode[256] {
		case 'n': // number
			t.handleNum(off)
		case 't': // token
			t.addToken(string(t.tmp))
		}
	}
}

func (t *Tokenizer) addToken(s string) {
	t.mode = valueMap
	if t.exkey {
		t.handler.Key(s)
		t.mode = colonMap
		t.exkey = false
	} else {
		switch s {
		case "null":
			t.handler.Null()
		case "true":
			t.handler.Bool(true)
		case "false":
			t.handler.Bool(false)
		default:
			t.handler.String(s)
		}
		t.exkey = 0 < len(t.starts) && t.starts[len(t.starts)-1] == objectStart
	}
}

func (t *Tokenizer) addString(s string) {
	t.mode = valueMap
	if t.exkey {
		t.handler.Key(s)
		t.mode = colonMap
		t.exkey = false
	} else {
		t.handler.String(s)
		t.exkey = 0 < len(t.starts) && t.starts[len(t.starts)-1] == objectStart
	}
}

func (t *Tokenizer) newError(off int, format string, args ...interface{}) {
	panic(&oj.ParseError{
		Message: fmt.Sprintf(format, args...),
		Line:    t.line,
		Column:  off - t.noff,
	})
}

func (t *Tokenizer) byteError(off int, mode string, b byte) {
	err := &oj.ParseError{
		Line:   t.line,
		Column: off - t.noff,
	}
	switch mode {
	case colonMap:
		err.Message = fmt.Sprintf("expected a colon, not '%c'", b)
	case negMap, zeroMap, digitMap, dotMap, fracMap, expSignMap, expZeroMap, expMap:
		err.Message = "invalid number"
	case stringMap:
		err.Message = fmt.Sprintf("invalid JSON character 0x%02x", b)
	case escMap:
		err.Message = fmt.Sprintf("invalid JSON escape character '\\%c'", b)
	case uMap:
		err.Message = fmt.Sprintf("invalid JSON unicode character '%c'", b)
	case spaceMap:
		err.Message = fmt.Sprintf("extra characters after close, '%c'", b)
	default:
		err.Message = fmt.Sprintf("unexpected character '%c'", b)
	}
	panic(err)
}

func (t *Tokenizer) handleNum(off int) {
	if t.exkey {
		t.newError(off, "expected a key")
	}
	t.mode = valueMap
	t.exkey = 0 < len(t.starts) && t.starts[len(t.starts)-1] == objectStart
	switch tn := t.num.AsNum().(type) {
	case int64:
		t.handler.Int(tn)
	case float64:
		t.handler.Float(tn)
	case json.Number:
		t.handler.Number(string(tn))
	}
}
