package scanner

import (
	"fmt"
	"myGo/mytoken"
	"path/filepath"
)

// An ErrorHandler may be provided to Scanner.Init. If a syntax error is
// encountered and a handler was installed, the handler is called with a
// position and an error message. The position points to the beginning of
// the offending token.
//
type ErrorHandler func(pos mytoken.Position, msg string)

// A Scanner holds the scanner's internal state while processing
// a given text. It can be allocated as part of another data
// structure but must be initialized via Init before use.
//
type Scanner struct {
	file *mytoken.File // source file handle
	dir  string        // directory portion of file.Name()
	src  []byte        // source
	err  ErrorHandler  // error reporting; or null
	mode Mode          // scanning mode

	// scanning state
	ch         rune // current character
	offset     int  // character offset
	rdOffset   int  // reading offset(position after current character)
	lineOffset int  // current line offset
	insertSemi bool // insert a semicolon before next newline

	ErrorCount int // number of errors encountered
}

const bom = 0xFEFF //  byte order mark, only permitted as very first character

// Read the next Unicode char into s.ch
// s.ch < 0 means EOF
//
func (s *Scanner) next() {
	if s.rdOffset < len(s.src) {
		s.offset = s.rdOffset
		if s.ch == '\n' {
			s.lineOffset = s.offset
			s.file.AddLine(s.offset)
		}
		r, w := rune(s.src[s.offset]), 1
		if r == 0 {
			s.error(s.offset, "illegal character")
		}
		s.rdOffset += w
		s.ch = r
	} else {
		s.offset = len(s.src)
		if s.ch == '\n' {
			s.lineOffset = s.offset
			s.file.AddLine(s.offset)
		}
		s.ch = -1 //EOF
	}
}

// A mode value is a set of flags (or 0).
// They control scanner behavior.
//
type Mode uint

const (
	ScanComments Mode = 1 << iota
	dontInsertSemis
)

func (s *Scanner) Init(file *mytoken.File, src []byte, err ErrorHandler, mode Mode) {
	if file.Size() != len(src) {
		panic(fmt.Sprintf("file size (%d) does not match src len (%d)", file.Size(), len(src)))
	}
	s.file = file
	s.dir, _ = filepath.Split(file.Name())
	s.src = src
	s.err = err
	s.mode = mode

	s.ch = ' '
	s.offset = 0
	s.rdOffset = 0
	s.lineOffset = 0
	s.insertSemi = false
	s.ErrorCount = 0

	s.next()
	if s.ch == bom {
		s.next()
	}
}

func (s *Scanner) error(offs int, msg string) {
	if s.err != nil {
		s.err(s.file.Position(s.file.Pos(offs)), msg)
	}
	s.ErrorCount++
}

func isLetter(ch rune) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_'
}

func isDigit(ch rune) bool {
	return '0' <= ch && ch <= '9'
}

func (s *Scanner) skipWhitespace() {
	for s.ch == ' ' || s.ch == '\t' || s.ch == '\n' || s.ch == '\r' {
		s.next()
	}
}

func (s *Scanner) scanComment() string {
	// initial '/' already consumed; s.ch == '*'
	offs := s.offset - 1 // position of initial '/'
	s.next()
	for s.ch >= 0 {
		ch := s.ch
		s.next()
		if ch == '*' && s.ch == '/' {
			s.next()
			goto exit
		}
	}
	s.error(offs, "comment not terminated")
exit:
	lit := s.src[offs:s.offset]
	return string(lit)
}

func (s *Scanner) scanIdentifier() string {
	offset := s.offset
	for isLetter(s.ch) || isDigit(s.ch) {
		s.next()
	}
	return string(s.src[offset:s.offset])
}

func digitVal(ch rune) int {
	switch {
	case '0' <= ch && ch <= '9':
		return int(ch - '0')
	case 'a' <= ch && ch <= 'f':
		return int(ch - 'a' + 10)
	case 'A' <= ch && ch <= 'F':
		return int(ch - 'A' + 10)
	default:
		return 16
	}
}

func (s *Scanner) scanMantissa(base int) {
	for digitVal(s.ch) < base {
		s.next()
	}
}

func (s *Scanner) scanNumber(seenDecimalPoint bool) (mytoken.Token, string) {
	offs := s.offset
	tok := mytoken.INT

	if seenDecimalPoint {
		offs--
		tok = mytoken.FLOAT
		s.scanMantissa(10)
		goto exponent
	}
	if s.ch == '0' {
		// int or float
		offs := s.offset
		s.next()
		// hexadecimal int
		if s.ch == 'x' || s.ch == 'X' {
			s.next()
			s.scanMantissa(16)
			if s.offset-offs <= 2 {
				// only "0x" or "0X"
				s.error(offs, "illegal hexadecimal number")
			}
		} else {
			// octal int or float
			seenDecimalDigit := false
			s.scanMantissa(8)
			if s.ch == '8' || s.ch == '9' {
				// illegal octal int or float
				seenDecimalDigit = true
				s.scanMantissa(10)
			}
			if s.ch == '.' || s.ch == 'e' || s.ch == 'E' {
				goto fraction
			}
			if seenDecimalDigit {
				s.error(offs, "illegal octal number")
			}
		}
		goto exit
	}

	// decimal int or float
	s.scanMantissa(10)
fraction:
	if s.ch == '.' {
		tok = mytoken.FLOAT
		s.next()
		s.scanMantissa(10)
	}

exponent:
	if s.ch == 'e' || s.ch == 'E' {
		tok = mytoken.FLOAT
		s.next()
		if s.ch == '-' || s.ch == '+' {
			s.next()
		}
		if digitVal(s.ch) < 10 {
			s.scanMantissa(10)
		} else {
			s.error(offs, "illegal floating-point exponent")
		}
	}
exit:
	return tok, string(s.src[offs:s.offset])
}

// scanEscape parses an escape sequence where rune is the accepted
// escaped quote. In case of a syntax error, it stops at the offending
// character (without consuming it) and returns false. Otherwise
// it returns true.(diff go, it only support anscii)
func (s *Scanner) scanEscape(quote rune) bool {
	offs := s.offset
	switch s.ch {
	case 'a', 'b', 'f', 'n', 'r', 't', 'v', '\\', quote:
		s.next()
		return true
	default:
		s.error(offs, "unknown escape sequence")
		return false
	}
}

func (s *Scanner) scanString() string {
	// '"' consumed
	offs := s.offset
	for {
		ch := s.ch
		if ch == '\n' || ch < 0 {
			s.error(offs, "string literal not terminated")
			break
		}
		s.next()
		if ch == '"' {
			break
		}
		if ch == '\\' {
			s.scanEscape('"')
		}
	}
	return string(s.src[offs:s.offset])
}

func (s *Scanner) scanRune() string {
	// '\'' has already consumed
	offs := s.offset - 1
	valid := true
	n := 0
	for {
		ch := s.ch
		if ch == '\n' || ch < 0 {
			if valid {
				s.error(offs, "rune literal not terminated")
				valid = false
			}
			break
		}
		s.next()
		if ch == '\'' {
			break
		}
		n++
		if ch == '\\' {
			if !s.scanEscape('\'') {
				valid = false
			}
		}
	}

	if valid && n != 1 {
		s.error(offs, "illegal rune literal")
	}
	return string(s.src[offs:s.offset])
}

func stripCR(b []byte) []byte {
	c := make([]byte, len(b))
	i := 0
	for _, ch := range b {
		if ch != '\r' {
			c[i] = ch
			i++
		}
	}
	return c[:i]
}

func (s *Scanner) scanRawString() string {
	// '`' consumed
	offs := s.offset - 1

	hasCR := false
	for {
		ch := s.ch
		if ch < 0 {
			s.error(offs, "raw string literal not terminated")
			break
		}
		s.next()
		if ch == '`' {
			break
		}
		if ch == '\r' {
			hasCR = true
		}
	}
	lit := s.src[offs:s.offset]
	if hasCR {
		lit = stripCR(lit)
	}
	return string(lit)
}

// Helper functions for scanning multi-byte tokens such as >> += >>=
//
func (s *Scanner) switch2(tok0, tok1 mytoken.Token) mytoken.Token {
	if s.ch == '=' {
		s.next()
		return tok1
	}
	return tok0
}

func (s *Scanner) switch3(tok0, tok1 mytoken.Token, ch2 rune, tok2 mytoken.Token) mytoken.Token {
	if s.ch == '=' {
		s.next()
		return tok1
	}
	if s.ch == ch2 {
		s.next()
		return tok2
	}
	return tok0
}

func (s *Scanner) switch4(tok0, tok1 mytoken.Token, ch2 rune, tok2, tok3 mytoken.Token) mytoken.Token {
	if s.ch == '=' {
		s.next()
		return tok1
	}
	if s.ch == ch2 {
		s.next()
		if s.ch == '=' {
			s.next()
			return tok3
		}
		return tok2
	}
	return tok0
}

func (s *Scanner) Scan() (pos mytoken.Pos, tok mytoken.Token, lit string) {
	s.skipWhitespace()

	// current token start
	pos = s.file.Pos(s.offset)

	switch ch := s.ch; {
	case isLetter(ch):
		lit = s.scanIdentifier()
		if len(lit) > 1 {
			// keyword are longer than letter
			tok = mytoken.Lookup(lit)
		} else {
			tok = mytoken.IDENT
		}
	case isDigit(ch):
		tok, lit = s.scanNumber(false)
	default:
		s.next()
		switch ch {
		case -1:
			tok = mytoken.EOF
		case '\n':
			return pos, mytoken.SEMICOLON, "\n"
		case '"':
			tok = mytoken.STRING
			lit = s.scanString()
		case '\'':
			tok = mytoken.CHAR
			lit = s.scanRune()
		case '`':
			tok = mytoken.STRING
			lit = s.scanRawString()
		case ':':
			tok = s.switch2(mytoken.COLON, mytoken.DEFINE)
		case '.':
			if '0' <= s.ch && s.ch <= '9' {
				tok, lit = s.scanNumber(true)
			} else {
				tok = mytoken.PERIOD
			}
		case ',':
			tok = mytoken.COMMA
		case ';':
			tok = mytoken.SEMICOLON
		case '(':
			tok = mytoken.LPAREN
		case ')':
			tok = mytoken.RPAREN
		case '[':
			tok = mytoken.LBRACK
		case ']':
			tok = mytoken.RBRACK
		case '{':
			tok = mytoken.LBRACE
		case '}':
			tok = mytoken.RBRACE
		case '+':
			tok = s.switch3(mytoken.ADD, mytoken.ADD_ASSIGN, '+', mytoken.INC)
		case '-':
			tok = s.switch3(mytoken.SUB, mytoken.SUB_ASSIGN, '-', mytoken.DEC)
		case '*':
			tok = s.switch2(mytoken.MUL, mytoken.MUL_ASSIGN)
		case '/':
			if s.ch == '*' {
				lit = s.scanComment()
				tok = mytoken.COMMENT
			} else {
				tok = s.switch2(mytoken.QUO, mytoken.QUO_ASSIGN)
			}
		case '%':
			tok = s.switch2(mytoken.REM, mytoken.REM_ASSIGN)
		case '^':
			tok = s.switch2(mytoken.XOR, mytoken.XOR_ASSIGN)
		case '<':
			tok = s.switch4(mytoken.LSS, mytoken.LEQ, '<', mytoken.SHL, mytoken.SHL_ASSIGN)
		case '>':
			tok = s.switch4(mytoken.GTR, mytoken.GEQ, '>', mytoken.SHR, mytoken.SHR_ASSIGN)
		case '=':
			tok = s.switch2(mytoken.ASSIGN, mytoken.EQL)
		case '!':
			tok = s.switch2(mytoken.NOT, mytoken.NEQ)
		case '&':
			if s.ch == '^' {
				s.next()
				tok = s.switch2(mytoken.AND_NOT, mytoken.AND_NOT_ASSIGN)
			} else {
				tok = s.switch3(mytoken.AND, mytoken.AND_ASSIGN, '&', mytoken.LAND)
			}
		case '|':
			tok = s.switch3(mytoken.OR, mytoken.OR_ASSIGN, '|', mytoken.LOR)
		default:
			tok = mytoken.ILLEGAL
			lit = string(ch)
		}
	}
	return
}
