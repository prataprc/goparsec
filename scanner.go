// Copyright (c) 2013 Goparsec AUTHORS. All rights reserved.
// Use of this source code is governed by LICENSE file.

package parsec

import "regexp"
import "reflect"
import "unsafe"
import "unicode"
import "bytes"
import "strings"

// Scanner interface defines necessary methods to match the input stream.
type Scanner interface {
	// SetWSPattern to configure white space pattern. Typically used as
	//		scanner := NewScanner(input).SetWSPattern(" ")
	SetWSPattern(pattern string) Scanner

	// TrackLineno as cursor moves forward, this can slow down parsing.
	// Useful when developing with parsec package.
	TrackLineno() Scanner

	// Clone will return new clone of the underlying scanner structure.
	// This will be used by combinators to _backtrack_.
	Clone() Scanner

	// GetCursor gets the current cursor position inside input text.
	GetCursor() int

	// Match the input stream with `pattern` and return matching string
	// after advancing the scanner's cursor.
	Match(pattern string) ([]byte, Scanner)

	// Match the input stream with a simple string, rather that a
	// pattern. It should be more efficient. Return a bool indicating
	// if the match was succesfull after advancing the scanner's cursor.
	MatchString(string) (bool, Scanner)

	// SubmatchAll the input stream with a choice of `patterns`
	// and return matching string and submatches, after advancing the
	// Scanner's cursor.
	SubmatchAll(pattern string) (map[string][]byte, Scanner)

	// SkipWs skips white space characters in the input stream.
	// Return skipped whitespaces as byte-slice after advance the
	// Scanner's cursor.
	SkipWS() ([]byte, Scanner)

	// SkipAny any occurrence of the elements of the slice.
	// Equivalent to Match(`(b[0]|b[1]|...|b[n])*`)
	// Returns Scanner after advancing its cursor.
	SkipAny(pattern string) ([]byte, Scanner)

	// Lineno return the current line-number of the cursor.
	Lineno() int

	// Endof detects whether end-of-file is reached in the input
	// stream and return a boolean indicating the same.
	Endof() bool
}

// SimpleScanner implements Scanner interface based on
// golang's regexp module.
type SimpleScanner struct {
	buf          []byte // input buffer
	cursor       int    // cursor within input buffer
	lineno       int
	patternCache map[string]*regexp.Regexp
	wsPattern    string // white space pattern used by SkipWS()
	// settings
	tracklineno bool
}

// NewScanner create and return a new instance of SimpleScanner object.
func NewScanner(text []byte) Scanner {
	return &SimpleScanner{
		buf:          text,
		cursor:       0,
		lineno:       1,
		patternCache: make(map[string]*regexp.Regexp),
		wsPattern:    `^[ \t\r\n]+`,
		tracklineno:  false,
	}
}

//---- Scanner{} interface.

// SetWSPattern implement Scanner{} interface.
func (s *SimpleScanner) SetWSPattern(pattern string) Scanner {
	s.wsPattern = pattern
	return s
}

// TrackLineno implement Scanner{} interface.
func (s *SimpleScanner) TrackLineno() Scanner {
	s.tracklineno = true
	return s
}

// Clone implement Scanner{} interface.
func (s *SimpleScanner) Clone() Scanner {
	return &SimpleScanner{
		buf:          s.buf,
		cursor:       s.cursor,
		lineno:       s.lineno,
		patternCache: s.patternCache,
		wsPattern:    s.wsPattern,
		tracklineno:  s.tracklineno,
	}
}

// GetCursor implement Scanner{} interface.
func (s *SimpleScanner) GetCursor() int {
	return s.cursor
}

// Match implement Scanner{} interface.
func (s *SimpleScanner) Match(pattern string) ([]byte, Scanner) {
	regc := s.getPattern(pattern)
	if token := regc.Find(s.buf[s.cursor:]); token != nil {
		if s.tracklineno && len(token) > 0 {
			s.lineno += len(bytes.Split(token, []byte{'\n'})) - 1
		}
		s.cursor += len(token)
		return token, s
	}
	return nil, s
}

// MatchString implement Scanner{} interface.
func (s *SimpleScanner) MatchString(str string) (bool, Scanner) {
	ln := len(str)
	if len(s.buf[s.cursor:]) < ln {
		return false, s
	} else if bytes.Compare(s.buf[s.cursor:s.cursor+ln], []byte(str)) != 0 {
		return false, s
	}
	if s.tracklineno && len(str) > 0 {
		s.lineno += len(strings.Split(str, "\n")) - 1
	}
	s.cursor += ln
	return true, s
}

// SubmatchAll implement Scanner{} interface.
func (s *SimpleScanner) SubmatchAll(patt string) (map[string][]byte, Scanner) {
	regc := s.getPattern(patt)
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
		if s.tracklineno && len(matches[0]) > 0 {
			s.lineno += len(bytes.Split(matches[0], []byte{'\n'})) - 1
		}
		s.cursor += len(matches[0])
		return captures, s
	}
	return nil, s
}

// SkipWS implement Scanner{} interface.
func (s *SimpleScanner) SkipWS() ([]byte, Scanner) {
	return s.SkipAny(s.wsPattern)
}

// SkipAny implement Scanner{} interface.
func (s *SimpleScanner) SkipAny(pattern string) ([]byte, Scanner) {
	if pattern[0] != '^' {
		pattern = "^" + pattern
	}
	return s.Match(pattern)
}

// Lineno implement Scanner{} interface.
func (s *SimpleScanner) Lineno() int {
	return s.lineno
}

// Endof implement Scanner{} interface.
func (s *SimpleScanner) Endof() bool {
	return s.cursor >= len(s.buf)
}

// SkipWSUnicode for looping through runes checking for whitespace.
func (s *SimpleScanner) SkipWSUnicode() ([]byte, Scanner) {
	for i, r := range bytes2str(s.buf[s.cursor:]) {
		if unicode.IsSpace(r) {
			if s.tracklineno && r == '\n' {
				s.lineno += 1
			}
			continue
		}
		token := s.buf[s.cursor : s.cursor+i]
		s.cursor += len(token)
		return token, s
	}
	token := s.buf[s.cursor:]
	s.cursor += len(token)
	return token, s
}

//---- local methods

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

func (s *SimpleScanner) resetcursor() {
	s.cursor = 0
}

func bytes2str(bytes []byte) string {
	if bytes == nil {
		return ""
	}
	sl := (*reflect.SliceHeader)(unsafe.Pointer(&bytes))
	st := &reflect.StringHeader{Data: sl.Data, Len: sl.Len}
	return *(*string)(unsafe.Pointer(st))
}
