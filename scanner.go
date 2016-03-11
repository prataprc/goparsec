// Copyright (c) 2013 Couchbase, Inc.

package parsec

import "regexp"
import "fmt"

var _ = fmt.Sprintf("dummy print")

// Scanner interface supplies necessary methods to match the
// input stream.
type Scanner interface {
	// Clone will return new clone of the underlying scanner structure.
	// This will be used by combinators to _backtrack_.
	Clone() Scanner

	// GetCursor gets the current cursor position inside input text.
	GetCursor() int

	// SetCursor to set the current cursor position inside input text,
	// return the old cursor position.
	SetCursor(cursor int) int

	// Match the input stream with `pattern` and return
	// matching string after advancing the cursor.
	Match(pattern string) ([]byte, Scanner)

	// SubmatchAll the input stream with a choice of `patterns`
	// and return matching string and submatches, after
	// advancing the cursor.
	SubmatchAll(pattern string) (map[string][]byte, Scanner)

	// SkipWs skips white space characters in the input stream.
	// Return skipped whitespaces as byte-slice and advance the cursor.
	SkipWS() ([]byte, Scanner)

	// Remaining returns the remaining un-parsed text as byte slice.
	Remaining() []byte

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

// SetCursor method receiver in Scanner{} interface.
func (s *SimpleScanner) SetCursor(cursor int) int {
	old := s.cursor
	s.cursor = cursor
	return old
}

// Match method receiver in Scanner{} interface.
func (s *SimpleScanner) Match(pattern string) ([]byte, Scanner) {
	var err error
	regc, ok := s.patternCache[pattern]
	if !ok {
		if regc, err = regexp.Compile(pattern); err != nil {
			panic(err)
		}
		s.patternCache[pattern] = regc
	}
	if token := regc.Find(s.buf[s.cursor:]); token != nil {
		s.cursor += len(token)
		return token, s
	}
	return nil, s
}

// SubmatchAll method receiver in Scanner{} interface.
func (s *SimpleScanner) SubmatchAll(
	pattern string) (map[string][]byte, Scanner) {

	var err error
	regc, ok := s.patternCache[pattern]
	if !ok {
		if regc, err = regexp.Compile(pattern); err != nil {
			panic(err)
		}
		s.patternCache[pattern] = regc
	}
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

// SkipWS method receiver in Scanner{} interface.
func (s *SimpleScanner) SkipWS() ([]byte, Scanner) {
	return s.Match(`^[ \t\r\n]+`)
}

// Remaining method receiver in Scanner{} interface.
func (s *SimpleScanner) Remaining() []byte {
	return s.buf[s.cursor:]
}

// Endof method receiver in Scanner{} interface.
func (s *SimpleScanner) Endof() bool {
	if s.cursor >= len(s.buf) {
		return true
	}
	return false
}
