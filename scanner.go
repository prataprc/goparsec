//  Copyright (c) 2013 Couchbase, Inc.
//  Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file
//  except in compliance with the License. You may obtain a copy of the License at
//    http://www.apache.org/licenses/LICENSE-2.0
//  Unless required by applicable law or agreed to in writing, software distributed under the
//  License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
//  either express or implied. See the License for the specific language governing permissions
//  and limitations under the License.

// scanner to parse terminals from input text.

package parsec

import (
    "regexp"
)

// SimpleScanner implements Scanner interface based on golang's regexp module.
type SimpleScanner struct {
    buf      []byte                    // input buffer
    cursor   int                       // cursor within input buffer
    patterns map[string]*regexp.Regexp // cache of compiled regular expression
}

// NewScanner creates and returns a reference to new instance of SimpleScanner
// object.
func NewScanner(text []byte) Scanner {
    return &SimpleScanner{
        buf:      text,
        cursor:   0,
        patterns: make(map[string]*regexp.Regexp),
    }
}

// Clone method receiver in Scanner interface.
func (s *SimpleScanner) Clone() Scanner {
    return &SimpleScanner{
        buf:      s.buf,
        cursor:   s.cursor,
        patterns: s.patterns,
    }
}

// GetCursor method receiver in Scanner interface.
func (s *SimpleScanner) GetCursor() int {
    return s.cursor
}

// Match method receiver in Scanner interface.
func (s *SimpleScanner) Match(pattern string) ([]byte, Scanner) {
    var regc *regexp.Regexp
    var err  error

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

// SubmatchAll method receiver in Scanner interface.
func (s *SimpleScanner) SubmatchAll(pattern string) ([][]byte, Scanner) {
    var regc *regexp.Regexp
    var err  error

    regc = s.patterns[pattern]
    if regc == nil {
        if regc, err = regexp.Compile(pattern); err == nil {
            s.patterns[pattern] = regc
        } else {
            panic(err.Error())
        }
    }
    toks := regc.FindAllSubmatch(s.buf[s.cursor:], 1)
    if len(toks) == 1 && toks[0] != nil && len(toks[0]) > 0 {
        s.cursor += len(toks[0][0])
        return toks[0], s
    }
    return nil, s
}

// SkipWS method receiver in Scanner interface.
func (s *SimpleScanner) SkipWS() ([]byte, Scanner) {
    return s.Match(`^[ \t\r\n]+`)
}

// Endof method receiver in Scanner interface.
func (s *SimpleScanner) Endof() bool {
    if s.cursor >= len(s.buf) {
        return true
    }
    return false
}
