// scanner to parse terminals from input text.
package parsec

import (
	"regexp"
)

type SimpleScanner struct {
	buf      []byte                    // input buffer
	cursor   int                       // cursor within input buffer
	patterns map[string]*regexp.Regexp // cache of compiled regular expression
}

func NewScanner(text []byte) Scanner {
	return &SimpleScanner{
		buf:      text,
		cursor:   0,
		patterns: make(map[string]*regexp.Regexp),
	}
}

func (s *SimpleScanner) Clone() Scanner {
	return &SimpleScanner{
		buf:      s.buf,
		cursor:   s.cursor,
		patterns: s.patterns,
	}
}

func (s *SimpleScanner) GetCursor() int {
	return s.cursor
}

// Match current input with `pattern` regular expression.
func (s *SimpleScanner) Match(pattern string) ([]byte, Scanner) {
	var regc *regexp.Regexp
	var err error

	if pattern[0] != '^' {
		panic("match patterns must begin with `^`")
	}

	regc = s.patterns[pattern]
	if regc == nil {
		if regc, err = regexp.Compile(pattern); err == nil {
			s.patterns[pattern] = regc
		} else {
			panic(err.Error())
		}
	}
	if token := regc.Find(s.buf[s.cursor:]); token != nil {
		s.cursor += len(token)
		return token, s
	}
	return nil, s
}

func (s SimpleScanner) SkipWS() ([]byte, Scanner) {
	return s.Match(`^[ \t\r\n]+`)
}

func (s SimpleScanner) Endof() bool {
	if s.cursor >= len(s.buf) {
		return true
	}
	return false
}
