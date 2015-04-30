// Copyright (c) 2013 Couchbase, Inc.

package parsec

import "regexp"

// Scanner interface supplies necessary methods to match the
// input stream.
type Scanner interface {
	// Clone will return new clone of the underlying scanner structure.
	// This will be used by combinators to _backtrack_.
	Clone() Scanner

	// GetCursor gets the current cursor position inside input text.
	GetCursor() int

	// Match the input stream with `pattern` and return
	// matching string after advancing the cursor.
	Match(pattern string) ([]byte, Scanner)

	// Match the input stream with a simple string,
	// rather that a pattern. It should be more efficient.
	// Returns a bool indicating if the match was succesfull
	// and the scanner
	MatchString(string) (bool, Scanner)

	// SubmatchAll the input stream with a choice of `patterns`
	// and return matching string and submatches, after
	// advancing the cursor.
	SubmatchAll(pattern string) (map[string][]byte, Scanner)

	// Skips any occurence of the elements of the slice.
	// Equivalent to Match(`(b[0]|b[1]|...|b[n])*`)
	// Returns Scanner and advances the cursor.
	SkipAny(b []byte) Scanner

	// Endof detects whether end-of-file is reached in the input
	// stream and return a boolean indicating the same.
	Endof() bool
}

// SimpleScanner implements Scanner interface based on
// golang's regexp module.
type SimpleScanner struct {
	buf          []byte // input buffer
	cursor       int    // cursor within input buffer
	patternCache map[string]*regexp.Regexp
}

// NewScanner creates and returns a reference to new instance
// of SimpleScanner object.
func NewScanner(text []byte) Scanner {
	return &SimpleScanner{
		buf:          text,
		cursor:       0,
		patternCache: make(map[string]*regexp.Regexp),
	}
}

// Clone method receiver in Scanner{} interface.
func (s *SimpleScanner) Clone() Scanner {
	return &SimpleScanner{
		buf:          s.buf,
		cursor:       s.cursor,
		patternCache: s.patternCache,
	}
}

// GetCursor method receiver in Scanner{} interface.
func (s *SimpleScanner) GetCursor() int {
	return s.cursor
}

func (s *SimpleScanner) getPattern(pattern string) *regexp.Regexp {
	regc, ok := s.patternCache[pattern]
	if !ok {
		var err error
		if regc, err = regexp.Compile(pattern); err != nil {
			panic(err)
		}
		s.patternCache[pattern] = regc
	}

	return regc
}

// Match method receiver in Scanner{} interface.
func (s *SimpleScanner) Match(pattern string) ([]byte, Scanner) {
	regc := s.getPattern(pattern)
	if token := regc.Find(s.buf[s.cursor:]); token != nil {
		s.cursor += len(token)
		return token, s
	}
	return nil, s
}

// MatchString method receiver in Scanner{} interface.
func (s *SimpleScanner) MatchString(str string) (bool, Scanner) {
	if len(s.buf[s.cursor:]) < len(str) {
		return false, s
	}

	for i, b := range []byte(str) {
		if s.buf[s.cursor+i] != b {
			return false, s
		}
	}

	s.cursor += len(str)
	return true, s
}

// SubmatchAll method receiver in Scanner{} interface.
func (s *SimpleScanner) SubmatchAll(
	pattern string) (map[string][]byte, Scanner) {

	regc := s.getPattern(pattern)
	matches := regc.FindSubmatch(s.buf[s.cursor:])

	if matches != nil {
		captures := make(map[string][]byte)
		names := regc.SubexpNames()
		for i, name := range names {
			if i == 0 || name == "" || matches[i] == nil {
				continue
			}
			captures[name] = matches[i]
		}
		s.cursor += len(matches[0])
		return captures, s
	}
	return nil, s
}

// SkipAny method receiver in Scanner{} interface.
func (s *SimpleScanner) SkipAny(bytes []byte) Scanner {
	matching := true

	if s.Endof() || bytes == nil {
		return s
	}

	for matching == true {
		matching = false
		for _, v := range bytes {
			if s.buf[s.cursor] == v {
				s.cursor++
				matching = true

				if s.Endof() {
					return s
				}
			}
		}
	}

	return s
}

// Endof method receiver in Scanner{} interface.
func (s *SimpleScanner) Endof() bool {
	return s.cursor >= len(s.buf)
}
