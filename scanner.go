// Copyright (c) 2013 Couchbase, Inc.

package parsec

import "regexp"

// SimpleScanner implements Scanner interface based on
// golang's regexp module.
type SimpleScanner struct {
	buf      []byte                    // input buffer
	cursor   int                       // cursor within input buffer
	patterns map[string]*regexp.Regexp // cache of compiled regular expression
}

// NewScanner creates and returns a reference to new instance
// of SimpleScanner object.
func NewScanner(text []byte) Scanner {
	return &SimpleScanner{
		buf:      text,
		cursor:   0,
		patterns: make(map[string]*regexp.Regexp),
	}
}

// Clone method receiver in Scanner{} interface.
func (s *SimpleScanner) Clone() Scanner {
	return &SimpleScanner{
		buf:      s.buf,
		cursor:   s.cursor,
		patterns: s.patterns,
	}
}

// GetCursor method receiver in Scanner{} interface.
func (s *SimpleScanner) GetCursor() int {
	return s.cursor
}

// Match method receiver in Scanner{} interface.
func (s *SimpleScanner) Match(pattern string) ([]byte, Scanner) {
	var err error
	regc := s.patterns[pattern]
	if regc == nil {
		regc, err = regexp.Compile(pattern)
		if err != nil {
			panic(err)
		}
		s.patterns[pattern] = regc
	}
	if token := regc.Find(s.buf[s.cursor:]); token != nil {
		s.cursor += len(token)
		return token, s
	}
	return nil, s
}

// SubmatchAll method receiver in Scanner{} interface.
func (s *SimpleScanner) SubmatchAll(pattern string) ([][]byte, Scanner) {
	var err error
	regc := s.patterns[pattern]
	if regc == nil {
		regc, err = regexp.Compile(pattern)
		if err != nil {
			panic(err.Error())
		}
		s.patterns[pattern] = regc
	}
	toks := regc.FindSubmatch(s.buf[s.cursor:])
	if toks != nil {
		s.cursor += len(toks[0])
		return toks, s
	}
	return nil, s
}

// SkipWS method receiver in Scanner{} interface.
func (s *SimpleScanner) SkipWS() ([]byte, Scanner) {
	return s.Match(`^[ \t\r\n]+`)
}

// Endof method receiver in Scanner{} interface.
func (s *SimpleScanner) Endof() bool {
	if s.cursor >= len(s.buf) {
		return true
	}
	return false
}
